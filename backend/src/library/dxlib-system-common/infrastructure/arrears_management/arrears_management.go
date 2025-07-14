package arrears_management

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/log"
	dxlibModule "github.com/donnyhardyanto/dxlib/module"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	_ "github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	_ "net/http"
	"time"
)

type SubTaskEventFunc func(aepr *api.DXAPIEndPointRequest, dtx *database.DXDatabaseTx, userId int64, subTaskId int64, subTask map[string]any, subTaskReportId int64, report map[string]any) (err error)

type ArrearsManagement struct {
	dxlibModule.DXModule

	TaskArrears      *table.DXTable
	UploadArrearsRow *table.DXTable
}

var ModuleArrearsManagement = ArrearsManagement{}

// UserValidationFunc defines the function signature for task validation
type UserValidationFunc func(aepr *api.DXAPIEndPointRequest, userId int64, subTask utils.JSON) error

// StatusUpdateFunc defines the function signature for status-specific updates
type StatusUpdateFunc func(oldStatus string, updateJSON utils.JSON) utils.JSON

type StateEngineStructSubTaskStatus struct {
	Aepr                   *api.DXAPIEndPointRequest
	UserType               base.UserType
	SubTaskStatusCondition []string
	NewSubTaskStatus       string
	OperationName          string
	Report                 utils.JSON
	OnExecute              SubTaskEventFunc
}

func (am *ArrearsManagement) Init(aDatabaseNameId string) {
	am.DatabaseNameId = aDatabaseNameId
	am.TaskArrears = table.Manager.NewTable(am.DatabaseNameId, "arrears_management.task_arrears", "arrears_management.task_arrears",
		"arrears_management.task_arrears", "id", "id", "uid", "data")
	am.UploadArrearsRow = table.Manager.NewTable(am.DatabaseNameId, "arrears_management.upload_arrears_row", "arrears_management.upload_arrears_row",
		"arrears_management.upload_arrears_row", "id", "id", "uid", "data")
}

func (am *ArrearsManagement) DoTaskArrearsCreate(l *log.DXLog, data utils.JSON) (err error, subTaskId int64) {
	fmt.Println("========= DoTaskArrearsCreate ==========")
	subTaskId = 0
	customerId, ok := data["customer_id"].(int64)
	if !ok {
		err = errors.Errorf("customer_id is blank")
		return err, subTaskId
	}
	_, _, err = task_management.ModuleTaskManagement.Customer.ShouldGetById(l, customerId)
	if err != nil {
		return err, subTaskId
	}

	spkNo, ok := data["spk_no"].(string)
	if !ok {
		err = errors.Errorf("spk_no is blank")
		return err, subTaskId
	}
	//code, ok := data["code"].(string)
	//if !ok {
	//	err = errors.Errorf("code is blank")
	//	return err, subTaskId
	//}
	code := spkNo
	_, r, _ := task_management.ModuleTaskManagement.Task.SelectOne(l, nil,
		utils.JSON{"code": code}, nil, nil)
	if r != nil {
		return errors.Errorf("Code already exist"), subTaskId
	}
	db := database.Manager.Databases[task_management.ModuleTaskManagement.DatabaseNameId]
	//var taskId int64
	subTaskTypeId, ok := data["sub_task_type_id"].(int64)
	if !ok {
		err = errors.Errorf("subTaskTypeId is blank")
		return err, subTaskId
	}
	fmt.Println("========= check subtask ==========")

	_, _, err = task_management.ModuleTaskManagement.SubTaskType.ShouldSelectOne(l, utils.JSON{
		"id":           subTaskTypeId,
		"task_type_id": base.TaskTypeIdDebtManagement,
	}, nil, nil)

	if err != nil {
		return err, subTaskId
	}
	fmt.Println("========= start db transaction ==========")
	txErr := db.Tx(l, database.LevelReadCommitted, func(tx *database.DXDatabaseTx) (err error) {
		taskId, err := task_management.ModuleTaskManagement.Task.TxInsert(tx,
			utils.JSON{
				"code":         code, //data["code"].(string),
				"spk_no":       spkNo,
				"task_type_id": base.TaskTypeIdDebtManagement,
				"customer_id":  customerId,
				"status":       base.TaskStatusWaitingAssignment,
			})
		if err != nil {
			fmt.Println("ada error ==========")
			// Create a dummy subTaskId for error case
			subTaskId = -1
			return err
		}
		fmt.Printf("id _____ %d", taskId)
		var subTaskIdLocal int64
		subTaskIdLocal, err4 := task_management.ModuleTaskManagement.SubTask.TxInsert(tx, utils.JSON{
			"task_id":          taskId,
			"sub_task_type_id": subTaskTypeId,
			"status":           base.SubTaskStatusWaitingAssignment,
		})
		if err4 != nil {
			// Keep the dummy subTaskId for error case
			subTaskId = -1
			return err4
		}
		subTaskId = subTaskIdLocal
		taskArrearsData := utils.JSON{}
		taskArrearsData["task_id"] = taskId
		taskArrearsData["amount_usage_bill"] = data["amount_usage_bill"]
		taskArrearsData["amount_fine"] = data["amount_fine"]
		taskArrearsData["amount_payment_guarantee"] = data["amount_payment_guarantee"]
		taskArrearsData["amount_reinstallation_cost"] = data["amount_reinstallation_cost"]
		taskArrearsData["amount_reflow_cost"] = data["amount_reflow_cost"]
		taskArrearsData["amount_bill_total"] = data["amount_bill_total"]
		taskArrearsData["period_begin"] = data["period_begin"]
		taskArrearsData["period_end"] = data["period_end"]
		//
		_, err5 := ModuleArrearsManagement.TaskArrears.TxInsert(tx, taskArrearsData)

		if err5 != nil {
			fmt.Printf("error:___ %v", err5.Error())
			// Keep the subTaskId for error case (it's already set to a valid value from SubTask.TxInsert)
			return err5
		} else {
			fmt.Println("===================")
		}
		return nil
	})

	if txErr != nil {
		fmt.Println("====error===")
		fmt.Printf("%s", txErr.Error())
		// subTaskId is already set appropriately (either to a valid value or -1)
		return txErr, subTaskId
	} else {
		fmt.Println("====no error===")
	}
	return nil, subTaskId
}

func (am *ArrearsManagement) DoUploadRowCreate(l *log.DXLog, data utils.JSON) (id int64, err error) {
	switch data["action"] {
	case "SPK SEGEL":
		data["sub_task_type_id"] = base.SubTaskTypeIdArrearsStopGasFlow
	case "SPK CABUT":
		data["sub_task_type_id"] = base.SubTaskTypeIdArrearsRemoveGasMeter
	case "SPK ALIR":
		data["sub_task_type_id"] = base.SubTaskTypeIdArrearsOpenGasFlow
	case "SPK PASANG":
		data["sub_task_type_id"] = base.SubTaskTypeIdArrearsReinstallGasMeter
	default:
		data["row_status"] = "ERROR"
	}
	if _, row, err := task_management.ModuleTaskManagement.Customer.ShouldSelectOne(l, utils.JSON{
		"customer_number": data["id_pelanggan"],
	}, nil, nil); err == nil {
		data["customer_id"] = row["id"]
	} else {
		data["row_status"] = "ERROR"
	}
	data["spk_no"] = data["nomor_surat"]
	layout := "2006-01-02"
	if dateStr, ok := data["periode_awal_tunggakan"].(string); ok {
		if date, err := time.Parse(layout, dateStr); err == nil {
			data["period_begin"] = date
		}
	}
	if dateStr, ok := data["periode_akhir_tunggakan"].(string); ok {
		if date, err := time.Parse(layout, dateStr); err == nil {
			data["period_end"] = date
		}
	}
	if dateStr, ok := data["tgl_surat"].(string); ok {
		if date, err := time.Parse(layout, dateStr); err == nil {
			data["date_issued"] = date
		}
	}
	id, err = am.UploadArrearsRow.Insert(l, data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (am *ArrearsManagement) DoUploadRowUpdate(l *log.DXLog, id int64, data utils.JSON) (err error) {
	_, err = am.UploadArrearsRow.UpdateOne(l, id, data)
	if err != nil {
		return err
	}
	return nil
}

func (am *ArrearsManagement) DoUploadProcess(l *log.DXLog, id int64, data utils.JSON) (err error) {
	// _, upload, err := am.UploadArrearsRow.Select(l, nil, nil, nil, nil, nil)
	// if err != nil {
	//	return err
	// }
	//for _, element := range upload {
	//	fcmApplicationId := element["id"].(int64)
	//	index is the index where we are
	//	// element is the element from someSlice for where we are
	//}
	return nil
}

// DoTaskArrearsUpdate updates an existing customer piutang record
func (am *ArrearsManagement) DoTaskArrearsUpdate(l *log.DXLog, id int64, customerPiutangData utils.JSON) (err error) {
	_, err = am.TaskArrears.UpdateOne(l, id, customerPiutangData)
	if err != nil {
		return err
	}

	return nil
}

// DoTaskArrearsGetById retrieves a customer piutang record by ID
func (am *ArrearsManagement) DoTaskArrearsGetById(l *log.DXLog, id int64) (exists bool, customerPiutang utils.JSON, err error) {
	_, customerPiutang, err = am.TaskArrears.ShouldSelectOne(l, utils.JSON{
		"id": id,
	}, nil, nil)
	if err != nil {
		return false, nil, err
	}

	return exists, customerPiutang, nil
}

// DoTaskArrearsGetByTaskId retrieves customer piutang records by customer ID
func (am *ArrearsManagement) DoTaskArrearsGetByTaskId(l *log.DXLog, taskId int64) (exists bool, taskPiutang utils.JSON, err error) {
	_, taskPiutang, err = am.TaskArrears.ShouldSelectOne(l, utils.JSON{
		"task_id": taskId,
	}, nil, nil)

	if err != nil {
		return false, nil, err
	}
	// Convert to map
	var data map[string]interface{}
	bytes, err := utils.JSONToBytes(taskPiutang)
	if err := json.Unmarshal(bytes, &data); err != nil {
		return false, nil, err
	}

	// Decode base64 data field if it exists
	if encodedData, ok := data["data"].(string); ok && encodedData != "" {
		decodedBytes, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			return exists, nil, fmt.Errorf("failed to decode base64: %v", err)
		}

		var jsonData interface{}
		if err := json.Unmarshal(decodedBytes, &jsonData); err != nil {
			return exists, nil, fmt.Errorf("failed to parse JSON: %v", err)
		}

		delete(data, "data")
		data["piutang"] = jsonData
	}

	return exists, data, nil

}
