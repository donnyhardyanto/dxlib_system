package handler

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/database/protected/db"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib/utils/json"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"net/http"
	"time"
)

/*func SimpleKeySecretAuthentication(aepr *api.DXAPIEndPointRequest) (err error) {
	_, key, err := aepr.GetParameterValueAsString("key")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, secret, err := aepr.GetParameterValueAsString("secret")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	username, secret, err := task_management.ModuleTaskManagement.SimpleKeySecretAuthentication(key, secret)
	if err != nil {
	}
}*/

func CustomerTaskConstructionCreate(aepr *api.DXAPIEndPointRequest) (err error) {
	p := aepr.GetParameterValues()
	taskCode := aepr.ParameterValues["task_code"].Value.(string)
	delete(p, "task_code")
	var taskId int64
	var customerId int64
	db := database.Manager.Databases[task_management.ModuleTaskManagement.DatabaseNameId]
	err = db.Tx(&aepr.Log, database.LevelReadCommitted, func(tx *database.DXDatabaseTx) (err2 error) {
		customerId, err2 = task_management.ModuleTaskManagement.Customer.TxInsert(tx, p)
		if err2 != nil {
			return err2
		}

		taskId, err2 = task_management.ModuleTaskManagement.TaskTxCreateConstruction(tx, taskCode, customerId, "", "")
		if err2 != nil {
			return err2
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_ = task_management.ModuleTaskManagement.DoNotifyTaskCreate(&aepr.Log, taskId)

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"task_id":     taskId,
		"customer_id": customerId,
	}})
	return nil
}

func TaskReadWithSubTask(aepr *api.DXAPIEndPointRequest) (err error) {
	isExistCustomerRegistrationNumber, customerRegistrationNumber, err := aepr.GetParameterValueAsString("registration_number")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	ifExistCustomerNumber, customerNumber, err := aepr.GetParameterValueAsString("customer_number")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	p := utils.JSON{}
	if isExistCustomerRegistrationNumber {
		p["registration_number"] = customerRegistrationNumber
	}
	if ifExistCustomerNumber {
		p["customer_number"] = customerNumber
	}
	if !(isExistCustomerRegistrationNumber || ifExistCustomerNumber) {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "NO_INPUT_PARAMETERS", "")
	}
	_, customer, err := task_management.ModuleTaskManagement.Customer.SelectOne(&aepr.Log, nil, p, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if customer == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "CUSTOMER_NOT_FOUND", "")
	}

	_, tasks, err := task_management.ModuleTaskManagement.Task.Select(&aepr.Log, nil, utils.JSON{
		"customer_id": customer["id"].(int64),
	}, nil,
		map[string]string{"id": "asc"}, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	for k, task := range tasks {
		taskId, err := json.GetInt64(task, "id")
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		_, subTasks, err := task_management.ModuleTaskManagement.SubTask.Select(&aepr.Log, nil, utils.JSON{
			"task_id": taskId,
		}, nil, map[string]string{"id": "asc"}, nil)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		for k1, subTask := range subTasks {
			subTaskId, err := json.GetInt64(subTask, "id")
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			_, subTaskHistoryItems, err := task_management.ModuleTaskManagement.SubTaskHistoryItem.Select(&aepr.Log, nil, utils.JSON{
				"sub_task_id": subTaskId,
			}, nil, map[string]string{"id": "asc"}, nil, nil)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			for k2, subTaskHistoryItem := range subTaskHistoryItems {
				subTaskReportId, ok := subTaskHistoryItem["sub_task_report_id"].(int64)
				if !ok {
					continue
				}
				_, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.GetById(&aepr.Log, subTaskReportId)
				if err != nil {
					return errors.Wrap(err, "error occured")
				}
				subTaskHistoryItems[k2]["sub_task_report"] = subTaskReport
			}
			_, subTaskReports, err := task_management.ModuleTaskManagement.SubTaskReport.Select(&aepr.Log, nil, utils.JSON{
				"sub_task_id": subTaskId,
			}, nil, map[string]string{"id": "asc"}, nil, nil)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			subTasks[k1]["sub_task_history_items"] = subTaskHistoryItems
			subTasks[k1]["sub_task_reports"] = subTaskReports
		}
		tasks[k]["sub_tasks"] = subTasks
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"tasks": tasks,
	}})
	return nil
}

func CustomerCancellationProgressStatusUpdate(aepr *api.DXAPIEndPointRequest) (err error) {
	isExistCustomerRegistrationNumber, customerRegistrationNumber, err := aepr.GetParameterValueAsString("registration_number")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	ifExistCustomerNumber, customerNumber, err := aepr.GetParameterValueAsString("customer_number")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, cancellationSubmissionStatus, err := aepr.GetParameterValueAsInt64("cancellation_submission_status")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	p := utils.JSON{}
	if isExistCustomerRegistrationNumber {
		p["registration_number"] = customerRegistrationNumber
	}
	if ifExistCustomerNumber {
		p["customer_number"] = customerNumber
	}
	if !(isExistCustomerRegistrationNumber || ifExistCustomerNumber) {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "NO_INPUT_PARAMETERS")
	}
	_, customer, err := task_management.ModuleTaskManagement.Customer.SelectOne(&aepr.Log, nil, p, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if customer == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "CUSTOMER_NOT_FOUND")
	}

	customerId, err := json.GetInt64(customer, "id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, err = task_management.ModuleTaskManagement.Customer.UpdateOne(&aepr.Log, customerId, utils.JSON{
		"cancellation_submission_status": cancellationSubmissionStatus,
	})

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"customer_id": customerId,
	}})
	return nil
}

func CustomerCancelByCustomerFromRelyOn(aepr *api.DXAPIEndPointRequest) (err error) {
	isExistCustomerRegistrationNumber, customerRegistrationNumber, err := aepr.GetParameterValueAsString("registration_number")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	ifExistCustomerNumber, customerNumber, err := aepr.GetParameterValueAsString("customer_number")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	p := utils.JSON{}
	if isExistCustomerRegistrationNumber {
		p["registration_number"] = customerRegistrationNumber
	}
	if ifExistCustomerNumber {
		p["customer_number"] = customerNumber
	}
	if !(isExistCustomerRegistrationNumber || ifExistCustomerNumber) {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "NO_INPUT_PARAMETERS")
	}
	_, customer, err := task_management.ModuleTaskManagement.Customer.SelectOne(&aepr.Log, nil, p, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if customer == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "CUSTOMER_NOT_FOUND")
	}

	customerId, err := json.GetInt64(customer, "id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldSelectOne(&aepr.Log, utils.JSON{
		"task_type_id": base.TaskTypeCodeConstruction,
		"customer_id":  customerId,
		"c1":           db.SQLExpression{Expression: "status NOT IN ('CANCELED_BY_CUSTOMER', 'CGP_VERIFICATION_SUCCESS')"},
	}, nil, map[string]string{"id": "asc"})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.Log.Debugf("CustomerCancelByCustomerFromRelyOn")
	subTaskId, err := json.GetInt64(subTask, "id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	report := aepr.ParameterValues["report"].Value.(map[string]any)

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, 0, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:     aepr,
		UserType: base.UserTypeNone,
		SubTaskStatusCondition: []string{
			base.SubTaskStatusBlockingDependency,
			base.SubTaskStatusWaitingAssignment,
			base.SubTaskStatusAssigned,
			base.SubTaskStatusWorking,
			base.SubTaskStatusPaused,
			base.SubTaskStatusReworking,
			base.SubTaskStatusFixing,
			base.SubTaskStatusVerificationFail,
			base.SubTaskStatusWaitingVerification},
		NewSubTaskStatus: base.SubTaskStatusCanceledByCustomer,
		OperationName:    base.UserAsFieldExecutorOperationNameSubTaskStatusCanceledByCustomer,
		Report:           report,
		OnExecute:        task_management.ModuleTaskManagement.OnSubTaskCancelByCustomer,
	}, true)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}
