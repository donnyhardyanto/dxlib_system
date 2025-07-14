package task_management

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/configuration"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/database/protected/db"
	"github.com/donnyhardyanto/dxlib/log"
	dxlibModule "github.com/donnyhardyanto/dxlib/module"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/lib"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/configuration_settings"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
)

type SubTaskTransitionValidateFunc func(aepr *api.DXAPIEndPointRequest, dtx *database.DXDatabaseTx, userId int64, subTaskId int64, subTask map[string]any, subTaskReportId int64,
	report map[string]any) (err error)
type SubTaskEventFunc func(aepr *api.DXAPIEndPointRequest, dtx *database.DXDatabaseTx, userId int64, subTaskId int64, subTask map[string]any, subTaskReportId int64, report map[string]any) (err error)

type TaskManagement struct {
	dxlibModule.DXModule

	TaskType    *table.DXTable
	SubTaskType *table.DXTable

	Customer           *table.DXTable
	CustomerMeter      *table.DXRawTable
	Task               *table.DXTable
	SubTask            *table.DXTable
	SubTaskHistoryItem *table.DXRawTable
	//TaskReport               *table.DXTable
	SubTaskReport            *table.DXRawTable
	SubTaskReportFileGroup   *table.DXTable
	SubTaskFile              *table.DXRawTable
	SubTaskReportFile        *table.DXRawTable
	SubTaskReportPicture     *lib.ImageObjectStorage
	SubTaskReportBeritaAcara *lib.ImageObjectStorage

	OnSubTaskPick                  SubTaskEventFunc
	OnSubTaskWorkingFinish         SubTaskEventFunc
	OnSubTaskVerifySuccess         SubTaskEventFunc
	OnSubTaskVerifyFail            SubTaskEventFunc
	OnSubTaskCGPVerifySuccess      SubTaskEventFunc
	OnSubTaskCGPVerifyFail         SubTaskEventFunc
	OnSubTaskFixingFinish          SubTaskEventFunc
	OnSubTaskPause                 SubTaskEventFunc
	OnSubTaskCancelByFieldExecutor SubTaskEventFunc
	OnSubTaskCancelByCustomer      SubTaskEventFunc
	OnSubTaskCancel                SubTaskEventFunc
}

var ModuleTaskManagement = TaskManagement{}

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
	OnValidateTransition   SubTaskTransitionValidateFunc
	OnExecute              SubTaskEventFunc
}

func (tm *TaskManagement) Init(aDatabaseNameId string) {
	tm.DatabaseNameId = aDatabaseNameId
	tm.TaskType = table.Manager.NewTable(tm.DatabaseNameId, "task_management.task_type", "task_management.task_type",
		"task_management.task_type", "code", "id", "uid", "data")
	tm.SubTaskType = table.Manager.NewTable(tm.DatabaseNameId, "task_management.sub_task_type", "task_management.sub_task_type",
		"task_management.v_sub_task_type", "full_code", "id", "uid", "data")

	tm.Customer = table.Manager.NewTable(tm.DatabaseNameId, "task_management.customer", "task_management.customer",
		"task_management.v_customer", "registration_number", "id", "uid", "data")
	tm.CustomerMeter = table.Manager.NewRawTable(tm.DatabaseNameId, "task_management.customer_meter", "task_management.customer_meter",
		"task_management.customer_meter", "uid", "id", "uid", "data")
	tm.Task = table.Manager.NewTable(tm.DatabaseNameId, "task_management.task", "task_management.task",
		"task_management.v_task", "code", "id", "uid", "data")
	tm.SubTask = table.Manager.NewTable(tm.DatabaseNameId, "task_management.sub_task", "task_management.sub_task",
		"task_management.v_sub_task", "code", "id", "uid", "data")

	tm.SubTaskHistoryItem = table.Manager.NewRawTable(tm.DatabaseNameId, "task_management.sub_task_history_item", "task_management.sub_task_history_item",
		"task_management.sub_task_history_item", "id", "id", "uid", "data")

	tm.SubTaskReport = table.Manager.NewRawTable(tm.DatabaseNameId, "task_management.sub_task_report", "task_management.sub_task_report",
		"task_management.v_sub_task_report", "id", "id", "uid", "data")
	tm.SubTaskReport.FieldTypeMapping = map[string]string{
		"report": "json",
	}
	tm.SubTaskReportFileGroup = table.Manager.NewTable(tm.DatabaseNameId, "task_management.sub_task_report_file_group", "task_management.sub_task_report_file_group",
		"task_management.sub_task_report_file_group", "id", "id", "uid", "data")
	tm.SubTaskFile = table.Manager.NewRawTable(tm.DatabaseNameId, "task_management.sub_task_file", "task_management.sub_task_file",
		"task_management.v_sub_task_file", "id", "id", "uid", "data")
	tm.SubTaskReportFile = table.Manager.NewRawTable(tm.DatabaseNameId, "task_management.sub_task_report_file", "task_management.sub_task_report_file",
		"task_management.v_sub_task_report_file", "id", "id", "uid", "data")

	configObjectStorage := *configuration.Manager.Configurations["object_storage"].Data

	configObjectStorageSubTaskReportSource := configObjectStorage["sub-task-report-picture-source"].(utils.JSON)
	configObjectStorageSubTaskReportSmall := configObjectStorage["sub-task-report-picture-small"].(utils.JSON)
	configObjectStorageStorageSubTaskReportMedium := configObjectStorage["sub-task-report-picture-medium"].(utils.JSON)
	configObjectStorageStorageSubTaskReportBig := configObjectStorage["sub-task-report-picture-big"].(utils.JSON)

	configSecurity := *configuration.Manager.Configurations["security"].Data
	configSecurityImageUploader := configSecurity["image_uploader"].(utils.JSON)

	tm.SubTaskReportPicture = lib.NewImageObjectStorage(configObjectStorageSubTaskReportSource["nameid"].(string),
		configSecurityImageUploader["max_request_size"].(int64),
		configSecurityImageUploader["max_pixel_width"].(int64),
		configSecurityImageUploader["max_pixel_height"].(int64),
		configSecurityImageUploader["max_bytes_per_pixel"].(int64),
		configSecurityImageUploader["max_pixels"].(int64),
		map[string]lib.ProcessedImageObjectStorage{
			"small": {
				ObjectStorageNameId: configObjectStorageSubTaskReportSmall["nameid"].(string),
				Width:               configObjectStorageSubTaskReportSmall["file_image_width"].(int),
				Height:              configObjectStorageSubTaskReportSmall["file_image_height"].(int),
			},
			"medium": {
				ObjectStorageNameId: configObjectStorageStorageSubTaskReportMedium["nameid"].(string),
				Width:               configObjectStorageStorageSubTaskReportMedium["file_image_width"].(int),
				Height:              configObjectStorageStorageSubTaskReportMedium["file_image_height"].(int),
			},
			"big": {
				ObjectStorageNameId: configObjectStorageStorageSubTaskReportBig["nameid"].(string),
				Width:               configObjectStorageStorageSubTaskReportBig["file_image_width"].(int),
				Height:              configObjectStorageStorageSubTaskReportBig["file_image_height"].(int),
			},
		})

	tm.SubTaskReportBeritaAcara = &lib.ImageObjectStorage{
		ObjectStorageSourceNameId: "berita-acara",
		MaxRequestSize:            lib.MaxRequestSize,
		ProcessedImages:           nil,
	}

	tm.OnSubTaskPick = tm.DoOnSubTaskPick
	tm.OnSubTaskWorkingFinish = tm.DoOnSubTaskWorkingFinish
	tm.OnSubTaskVerifyFail = tm.DoOnSubTaskVerifyFail
	tm.OnSubTaskCGPVerifySuccess = tm.DoOnSubTaskCGPVerifySuccess
	tm.OnSubTaskCancelByCustomer = tm.DoOnSubTaskCancelByCustomer
	//	tm.OnSubTaskCancel = tm.DoOnSubTaskCancel
}

type ConstructionSubTasksStatusResult struct {
	CustomerId                           int64
	SubTaskSK                            utils.JSON
	SubTaskSR                            utils.JSON
	SubTaskPMG                           utils.JSON
	SubTaskGasIn                         utils.JSON
	SubTaskSKId                          int64
	SubTaskSKStatus                      string
	SubTaskSKIsWorkingFinish             bool
	SubTaskSKIsVerificationSuccess       bool
	SubTaskSKIsCGPVerificationSuccess    bool
	SubTaskSKLastFormSubTaskReportId     *int64
	SubTaskSRId                          int64
	SubTaskSRStatus                      string
	SubTaskSRIsWorkingFinish             bool
	SubTaskSRIsVerificationSuccess       bool
	SubTaskSRIsCGPVerificationSuccess    bool
	SubTaskSRLastFormSubTaskReportId     *int64
	SubTaskPMGId                         int64
	SubTaskPMGStatus                     string
	SubTaskPMGIsWorkingFinish            bool
	SubTaskPMGIsVerificationSuccess      bool
	SubTaskPMGIsCGPVerificationSuccess   bool
	SubTaskPMGLastFormSubTaskReportId    *int64
	SubTaskGasInId                       int64
	SubTaskGasInStatus                   string
	SubTaskGasInIsWorkingFinish          bool
	SubTaskGasInIsVerificationSuccess    bool
	SubTaskGasInIsCGPVerificationSuccess bool
	SubTaskGasInLastFormSubTaskReportId  *int64
}

func (tm *TaskManagement) GetConstructionSubTasksStatus(dtx *database.DXDatabaseTx, taskId int64) (r ConstructionSubTasksStatusResult, err error) {
	_, r.SubTaskSK, err = tm.SubTask.TxSelectOne(dtx, utils.JSON{
		"task_id":                 taskId,
		"sub_task_type_full_code": base.SubTaskTypeFullCodeConstructionSK,
	}, nil)
	if err != nil {
		return r, errors.Wrap(err, "error occured at GetConstructionSubTasksStatus")
	}
	r.CustomerId = r.SubTaskSK["customer_id"].(int64)

	r.SubTaskSKId = r.SubTaskSK["id"].(int64)
	r.SubTaskSKStatus = r.SubTaskSK["status"].(string)
	r.SubTaskSKIsWorkingFinish = r.SubTaskSK["is_working_finish"].(bool)
	r.SubTaskSKIsVerificationSuccess = r.SubTaskSK["is_verification_success"].(bool)
	r.SubTaskSKIsCGPVerificationSuccess = r.SubTaskSK["is_cgp_verification_success"].(bool)
	if val, ok := r.SubTaskSK["last_form_sub_task_report_id"].(int64); ok {
		r.SubTaskSKLastFormSubTaskReportId = &val
	} else {
		r.SubTaskSKLastFormSubTaskReportId = nil
	}

	_, r.SubTaskSR, err = tm.SubTask.TxSelectOne(dtx, utils.JSON{
		"task_id":                 taskId,
		"sub_task_type_full_code": base.SubTaskTypeFullCodeConstructionSR,
	}, nil)
	if err != nil {
		return r, errors.Wrap(err, "error occured at GetConstructionSubTasksStatus")
	}
	r.SubTaskSRId = r.SubTaskSR["id"].(int64)
	r.SubTaskSRStatus = r.SubTaskSR["status"].(string)
	r.SubTaskSRIsWorkingFinish = r.SubTaskSR["is_working_finish"].(bool)
	r.SubTaskSRIsVerificationSuccess = r.SubTaskSR["is_verification_success"].(bool)
	r.SubTaskSRIsCGPVerificationSuccess = r.SubTaskSR["is_cgp_verification_success"].(bool)
	if val, ok := r.SubTaskSR["last_form_sub_task_report_id"].(int64); ok {
		r.SubTaskSRLastFormSubTaskReportId = &val
	} else {
		r.SubTaskSRLastFormSubTaskReportId = nil
	}

	_, r.SubTaskPMG, err = tm.SubTask.TxSelectOne(dtx, utils.JSON{
		"task_id":                 taskId,
		"sub_task_type_full_code": base.SubTaskTypeFullCodeConstructionMeterInstallation,
	}, nil)
	if err != nil {
		return r, errors.Wrap(err, "error occured at GetConstructionSubTasksStatus")
	}
	r.SubTaskPMGId = r.SubTaskPMG["id"].(int64)
	r.SubTaskPMGStatus = r.SubTaskPMG["status"].(string)
	r.SubTaskPMGIsWorkingFinish = r.SubTaskPMG["is_working_finish"].(bool)
	r.SubTaskPMGIsVerificationSuccess = r.SubTaskPMG["is_verification_success"].(bool)
	r.SubTaskPMGIsCGPVerificationSuccess = r.SubTaskPMG["is_cgp_verification_success"].(bool)
	if val, ok := r.SubTaskPMG["last_form_sub_task_report_id"].(int64); ok {
		r.SubTaskPMGLastFormSubTaskReportId = &val
	} else {
		r.SubTaskPMGLastFormSubTaskReportId = nil
	}

	_, r.SubTaskGasIn, err = tm.SubTask.TxSelectOne(dtx, utils.JSON{
		"task_id":                 taskId,
		"sub_task_type_full_code": base.SubTaskTypeFullCodeConstructionGasIn,
	}, nil)
	if err != nil {
		return r, errors.Wrap(err, "error occured at GetConstructionSubTasksStatus")
	}
	r.SubTaskGasInId = r.SubTaskGasIn["id"].(int64)
	r.SubTaskGasInStatus = r.SubTaskGasIn["status"].(string)
	r.SubTaskGasInIsWorkingFinish = r.SubTaskGasIn["is_working_finish"].(bool)
	r.SubTaskGasInIsVerificationSuccess = r.SubTaskGasIn["is_verification_success"].(bool)
	r.SubTaskGasInIsCGPVerificationSuccess = r.SubTaskGasIn["is_cgp_verification_success"].(bool)
	if val, ok := r.SubTaskGasIn["last_form_sub_task_report_id"].(int64); ok {
		r.SubTaskGasInLastFormSubTaskReportId = &val
	} else {
		r.SubTaskGasInLastFormSubTaskReportId = nil
	}

	return r, nil
}
func (tm *TaskManagement) CheckTaskConstructionSubTasksStatus(aepr *api.DXAPIEndPointRequest, dtx *database.DXDatabaseTx, taskId int64, currentSubTaskId int64, currentSubTask utils.JSON) (err error) {
	subTasksStatus, err := tm.GetConstructionSubTasksStatus(dtx, taskId)
	if err != nil {
		return errors.Wrap(err, "error occurred in CheckTaskConstructionSubTasksStatus ")
	}

	if subTasksStatus.SubTaskSKIsWorkingFinish || subTasksStatus.SubTaskSRIsWorkingFinish {
		if subTasksStatus.SubTaskPMGStatus == base.SubTaskStatusBlockingDependency {
			_, err = tm.SubTask.TxUpdate(dtx, utils.JSON{
				"status": base.SubTaskStatusWaitingAssignment,
			}, utils.JSON{
				"id": subTasksStatus.SubTaskPMGId,
			})
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
		}
	}

	if subTasksStatus.SubTaskSKIsWorkingFinish && subTasksStatus.SubTaskSRIsWorkingFinish && subTasksStatus.SubTaskPMGIsWorkingFinish {
		if subTasksStatus.SubTaskGasInStatus == base.SubTaskStatusBlockingDependency {
			_, err = tm.SubTask.TxUpdate(dtx, utils.JSON{
				"status": base.SubTaskStatusWaitingAssignment,
			}, utils.JSON{
				"id": subTasksStatus.SubTaskGasInId,
			})
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
		}
	}

	if currentSubTask["sub_task_type_full_code"] == base.SubTaskTypeFullCodeConstructionMeterInstallation {

		if subTasksStatus.SubTaskPMGStatus == base.SubTaskStatusWaitingVerification {
			lastFormSubTaskReport := utils.JSON{}
			if subTasksStatus.SubTaskPMGLastFormSubTaskReportId != nil {
				_, lastFormSubTaskReport, err = tm.SubTaskReport.TxShouldGetById(dtx, *subTasksStatus.SubTaskPMGLastFormSubTaskReportId)
				if err != nil {
					return errors.Wrap(err, "error occurred on DoOnSubTaskCGPVerifySuccess::get last form sub task report")
				}
				// Continue with the rest of the logic...
			} else {
				// Handle the case where the ID is nil
				// Either return an error or skip this processing
				return errors.New("SubTaskPMGLastFormSubTaskReportId is nil")
			}

			lastFormSubTaskReportReport := utils.JSON{}
			err = json.Unmarshal(lastFormSubTaskReport["report"].([]byte), &lastFormSubTaskReportReport)
			if err != nil {
				return errors.Wrap(err, "error occurred on DoOnSubTaskCGPVerifySuccess::unmarshal last form sub task report")
			}

			// Extract values from maps with proper error handling

			meterId, err := utils.ConvertToInt64(lastFormSubTaskReportReport["meter_id"])
			if err != nil {
				return errors.Wrap(err, "error getting meter_id")
			}

			gSizeId, err := utils.ConvertToInt64(lastFormSubTaskReportReport["g_size_id"])
			if err != nil {
				return errors.Wrap(err, "error getting g_size_id")
			}

			startCalibrationMonth, err := utils.ConvertToInt64(lastFormSubTaskReportReport["start_calibration_month"])
			if err != nil {
				return errors.Wrap(err, "error getting start_calibration_month")
			}

			startCalibrationYear, err := utils.ConvertToInt64(lastFormSubTaskReportReport["start_calibration_year"])
			if err != nil {
				return errors.Wrap(err, "error getting start_calibration_year")
			}

			t := lastFormSubTaskReport["timestamp"].(time.Time)

			_, _, err = tm.CustomerMeter.TxUpsert(dtx, utils.JSON{
				"meter_appliance_type_id": meterId,
				"meter_brand":             lastFormSubTaskReportReport["meter_brand"].(string),
				"sn_meter":                lastFormSubTaskReportReport["sn_meter"].(string),
				"g_size_id":               gSizeId,
				"qmin":                    lastFormSubTaskReportReport["qmin"].(float64),
				"qmax":                    lastFormSubTaskReportReport["qmax"].(float64),
				"start_calibration_month": startCalibrationMonth,
				"start_calibration_year":  startCalibrationYear,
				"register_timestamp":      t,
			}, utils.JSON{
				"customer_id": subTasksStatus.CustomerId,
			})
			if err != nil {
				return errors.Wrap(err, "error insert customer meter data")
			}

		}
	}

	if currentSubTask["sub_task_type_full_code"] == base.SubTaskTypeFullCodeConstructionGasIn {
		if subTasksStatus.SubTaskGasInStatus == base.SubTaskStatusWaitingVerification {
			lastFormSubTaskReport := utils.JSON{}
			if subTasksStatus.SubTaskGasInLastFormSubTaskReportId != nil {
				_, lastFormSubTaskReport, err = tm.SubTaskReport.TxShouldGetById(dtx, *subTasksStatus.SubTaskGasInLastFormSubTaskReportId)
				if err != nil {
					return errors.Wrap(err, "error occurred on DoOnSubTaskCGPVerifySuccess::get last form sub task report")
				}
				// Continue with the rest of the logic...
			} else {
				// Handle the case where the ID is nil
				// Either return an error or skip this processing
				return errors.New("SubTaskGasInLastFormSubTaskReportId is nil")
			}

			lastFormSubTaskReportReport := utils.JSON{}
			err = json.Unmarshal(lastFormSubTaskReport["report"].([]byte), &lastFormSubTaskReportReport)
			if err != nil {
				return errors.Wrap(err, "error occurred on DoOnSubTaskCGPVerifySuccess::unmarshal last form sub task report")
			}

			// Extract values from maps with proper error handling

			meterId, err := utils.ConvertToInt64(lastFormSubTaskReportReport["meter_id"])
			if err != nil {
				return errors.Wrap(err, "error getting meter_id")
			}

			gSizeId, err := utils.ConvertToInt64(lastFormSubTaskReportReport["g_size_id"])
			if err != nil {
				return errors.Wrap(err, "error getting g_size_id")
			}

			// Safely parse the date string
			gasInDateStr, ok := lastFormSubTaskReportReport["gas_in_date"].(string)
			if !ok {
				return errors.New("invalid type for gas_in_date, expected a string")
			}
			gasInDateTime, err := time.Parse(time.RFC3339, gasInDateStr)
			if err != nil {
				return errors.Wrap(err, "failed to parse gas_in_date")
			}

			_, _, err = tm.CustomerMeter.TxUpsert(dtx, utils.JSON{
				"meter_appliance_type_id":  meterId,
				"meter_brand":              lastFormSubTaskReportReport["meter_brand"].(string),
				"sn_meter":                 lastFormSubTaskReportReport["sn_meter"].(string),
				"g_size_id":                gSizeId,
				"meter_location_longitude": lastFormSubTaskReportReport["meter_location_longitude"].(float64),
				"meter_location_latitude":  lastFormSubTaskReportReport["meter_location_latitude"].(float64),
				"gas_in_date":              gasInDateTime.Format("2006-01-02"),
			}, utils.JSON{
				"customer_id": subTasksStatus.CustomerId,
			})

			if err != nil {
				return errors.Wrap(err, "error update customer meter data")
			}
		}
	}

	if (subTasksStatus.SubTaskGasInStatus == base.SubTaskStatusCGPVerificationSuccess) &&
		(subTasksStatus.SubTaskPMGStatus == base.SubTaskStatusCGPVerificationSuccess) &&
		(subTasksStatus.SubTaskSRStatus == base.SubTaskStatusCGPVerificationSuccess) &&
		(subTasksStatus.SubTaskSKStatus == base.SubTaskStatusCGPVerificationSuccess) {
		_, err = tm.Task.TxUpdate(dtx, utils.JSON{
			"status": base.TaskStatusCompleted,
		}, utils.JSON{
			"id": taskId,
		})
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	}

	return nil
}

func (tm *TaskManagement) SetConstructionSubTasksStatusToCanceledByCustomer(aepr *api.DXAPIEndPointRequest, dtx *database.DXDatabaseTx, userId int64, taskId int64, subTaskId int64, subTask map[string]any, subTaskReportId int64, report map[string]any) (err error) {
	aepr.Log.Debugf("SetConstructionSubTasksStatusToCanceledByCustomer:taskId:%d:subTaskId:%d", taskId, subTaskId)

	at := aepr.ParameterValues["at"].Value.(time.Time)

	_, subTaskSK, err := tm.SubTask.TxSelectOne(dtx, utils.JSON{
		"c1":                      db.SQLExpression{Expression: fmt.Sprintf("id!=%d", subTaskId)},
		"task_id":                 taskId,
		"sub_task_type_full_code": base.SubTaskTypeFullCodeConstructionSK,
	}, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTaskSR, err := tm.SubTask.TxSelectOne(dtx, utils.JSON{
		"c1":                      db.SQLExpression{Expression: fmt.Sprintf("id!=%d", subTaskId)},
		"task_id":                 taskId,
		"sub_task_type_full_code": base.SubTaskTypeFullCodeConstructionSR,
	}, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTaskMeterInstallation, err := tm.SubTask.TxSelectOne(dtx, utils.JSON{
		"c1":                      db.SQLExpression{Expression: fmt.Sprintf("id!=%d", subTaskId)},
		"task_id":                 taskId,
		"sub_task_type_full_code": base.SubTaskTypeFullCodeConstructionMeterInstallation,
	}, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTaskGasIn, err := tm.SubTask.TxSelectOne(dtx, utils.JSON{
		"c1":                      db.SQLExpression{Expression: fmt.Sprintf("id!=%d", subTaskId)},
		"task_id":                 taskId,
		"sub_task_type_full_code": base.SubTaskTypeFullCodeConstructionGasIn,
	}, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	if subTaskSK != nil {
		subTaskStatusSK := subTaskSK["status"].(string)
		if subTaskStatusSK != base.SubTaskStatusCanceledByCustomer {
			subTaskIdSK := subTaskSK["id"].(int64)
			_, _, _, err = tm.TxProcessSubTaskStateRaw(dtx, subTaskIdSK, userId, at, StateEngineStructSubTaskStatus{
				Aepr:                   aepr,
				UserType:               base.UserTypeNone,
				SubTaskStatusCondition: nil,
				NewSubTaskStatus:       base.SubTaskStatusCanceledByCustomer,
				OperationName:          base.UserAsFieldExecutorOperationNameSubTaskStatusCanceledByCustomer,
				Report:                 report,
				OnExecute:              nil,
			})
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
		}
	}

	if subTaskSR != nil {
		subTaskStatusSR := subTaskSR["status"].(string)
		if subTaskStatusSR != base.SubTaskStatusCanceledByCustomer {
			subTaskIdSR := subTaskSR["id"].(int64)
			_, _, _, err = tm.TxProcessSubTaskStateRaw(dtx, subTaskIdSR, userId, at, StateEngineStructSubTaskStatus{
				Aepr:                   aepr,
				UserType:               base.UserTypeNone,
				SubTaskStatusCondition: nil,
				NewSubTaskStatus:       base.SubTaskStatusCanceledByCustomer,
				OperationName:          base.UserAsFieldExecutorOperationNameSubTaskStatusCanceledByCustomer,
				Report:                 report,
				OnExecute:              nil,
			})
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
		}
	}

	if subTaskMeterInstallation != nil {
		subTaskStatusMeterInstallation := subTaskMeterInstallation["status"].(string)
		if subTaskStatusMeterInstallation != base.SubTaskStatusCanceledByCustomer {
			subTaskIdMeterInstallation := subTaskMeterInstallation["id"].(int64)
			_, _, _, err = tm.TxProcessSubTaskStateRaw(dtx, subTaskIdMeterInstallation, userId, at, StateEngineStructSubTaskStatus{
				Aepr:                   aepr,
				UserType:               base.UserTypeNone,
				SubTaskStatusCondition: nil,
				NewSubTaskStatus:       base.SubTaskStatusCanceledByCustomer,
				OperationName:          base.UserAsFieldExecutorOperationNameSubTaskStatusCanceledByCustomer,
				Report:                 report,
				OnExecute:              nil,
			})
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
		}
	}

	if subTaskGasIn != nil {
		subTaskStatusGasIn := subTaskGasIn["status"].(string)
		if subTaskStatusGasIn != base.SubTaskStatusCanceledByCustomer {
			subTaskIdGasIn := subTaskGasIn["id"].(int64)
			_, _, _, err = tm.TxProcessSubTaskStateRaw(dtx, subTaskIdGasIn, userId, at, StateEngineStructSubTaskStatus{
				Aepr:                   aepr,
				UserType:               base.UserTypeNone,
				SubTaskStatusCondition: nil,
				NewSubTaskStatus:       base.SubTaskStatusCanceledByCustomer,
				OperationName:          base.UserAsFieldExecutorOperationNameSubTaskStatusCanceledByCustomer,
				Report:                 report,
				OnExecute:              nil,
			})
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
		}
	}

	_, err = tm.Task.TxUpdate(dtx, utils.JSON{
		"status": base.TaskStatusCanceledByCustomer,
	}, utils.JSON{
		"id": taskId,
	})

	return nil
}

func (tm *TaskManagement) DoOnSubTaskPick(aepr *api.DXAPIEndPointRequest, dtx *database.DXDatabaseTx, userId int64, subTaskId int64, subTask map[string]any, subTaskReportId int64, report map[string]any) (err error) {
	taskId := subTask["task_id"].(int64)

	_, task, err := tm.Task.ShouldGetById(&aepr.Log, taskId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	taskTypeCode := task["task_type_code"].(string)
	taskStatus := task["status"].(string)

	switch taskTypeCode {
	case base.TaskTypeCodeConstruction:
		err = tm.CheckTaskConstructionSubTasksStatus(aepr, dtx, taskId, subTaskId, subTask)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	case base.TaskTypeCodeDebtManagement:
		return errors.Errorf("TASK_TYPE_CODE_NOT_IMPLEMENTED:%s", taskTypeCode)
	case base.TaskTypeCodeTechnicalSupport:
		return errors.Errorf("TASK_TYPE_CODE_NOT_IMPLEMENTED:%s", taskTypeCode)
	default:
		return errors.Errorf("TASK_TYPE_CODE_NOT_SUPPORTED:%s", taskTypeCode)
	}

	if taskStatus == base.TaskStatusWaitingAssignment {
		_, err = tm.Task.TxUpdate(dtx, utils.JSON{
			"status": base.TaskStatusInProgress,
		}, utils.JSON{
			"id": taskId,
		})
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	}

	return nil
}

func (tm *TaskManagement) DoOnSubTaskWorkingFinish(aepr *api.DXAPIEndPointRequest, dtx *database.DXDatabaseTx, userId int64, subTaskId int64, subTask map[string]any, subTaskReportId int64, report map[string]any) (err error) {
	taskId := subTask["task_id"].(int64)

	_, task, err := tm.Task.ShouldGetById(&aepr.Log, taskId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	taskTypeCode := task["task_type_code"].(string)

	switch taskTypeCode {
	case base.TaskTypeCodeConstruction:
		// if SK atau Gas In, create Berita Acara di sini, dan filenya dimasukkan ke minio bucket
		err = tm.CheckTaskConstructionSubTasksStatus(aepr, dtx, taskId, subTaskId, subTask)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	case base.TaskTypeCodeDebtManagement:
		return errors.Errorf("TASK_TYPE_CODE_NOT_IMPLEMENTED:%s", taskTypeCode)
	case base.TaskTypeCodeTechnicalSupport:
		return errors.Errorf("TASK_TYPE_CODE_NOT_IMPLEMENTED:%s", taskTypeCode)
	default:
		return errors.Errorf("TASK_TYPE_CODE_NOT_SUPPORTED:%s", taskTypeCode)
	}

	return nil
}

func (tm *TaskManagement) DoOnSubTaskVerifyFail(aepr *api.DXAPIEndPointRequest, dtx *database.DXDatabaseTx, userId int64, subTaskId int64, subTask map[string]any, subTaskReportId int64, report map[string]any) (err error) {
	_ = tm.DoNotifySubTaskVerificationFailed(&aepr.Log, subTaskId)

	return nil
}

func (tm *TaskManagement) DoOnSubTaskCGPVerifySuccess(aepr *api.DXAPIEndPointRequest, dtx *database.DXDatabaseTx, userId int64, subTaskId int64, subTask map[string]any, subTaskReportId int64, report map[string]any) (err error) {
	taskId := subTask["task_id"].(int64)

	_, task, err := tm.Task.ShouldGetById(&aepr.Log, taskId)
	if err != nil {
		return errors.Wrap(err, "error occurred on DoOnSubTaskCGPVerifySuccess::get task")
	}

	taskTypeCode := task["task_type_code"].(string)

	switch taskTypeCode {
	case base.TaskTypeCodeConstruction:

		err = tm.CheckTaskConstructionSubTasksStatus(aepr, dtx, taskId, subTaskId, subTask)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	case base.TaskTypeCodeDebtManagement:
		return errors.Errorf("TASK_TYPE_CODE_NOT_IMPLEMENTED:%s", taskTypeCode)
	case base.TaskTypeCodeTechnicalSupport:
		return errors.Errorf("TASK_TYPE_CODE_NOT_IMPLEMENTED:%s", taskTypeCode)
	default:
		return errors.Errorf("TASK_TYPE_CODE_NOT_SUPPORTED:%s", taskTypeCode)
	}

	return nil
}

func (tm *TaskManagement) DoOnSubTaskCancelByCustomer(aepr *api.DXAPIEndPointRequest, dtx *database.DXDatabaseTx, userId int64, subTaskId int64, subTask map[string]any, subTaskReportId int64, report map[string]any) (err error) {
	aepr.Log.Debugf("DoOnSubTaskCancelByCustomer:subTaskId:%d", subTaskId)
	taskId := subTask["task_id"].(int64)

	_, task, err := tm.Task.ShouldGetById(&aepr.Log, taskId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	taskTypeCode := task["task_type_code"].(string)

	switch taskTypeCode {
	case base.TaskTypeCodeConstruction:
		err = tm.SetConstructionSubTasksStatusToCanceledByCustomer(aepr, dtx, userId, taskId, subTaskId, subTask, subTaskReportId, report)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	case base.TaskTypeCodeDebtManagement:
		/*	err = tm.SetDebtManagementSubTasksStatusToCanceled(aepr, dtx, userId, taskId, subTaskId, subTask, subTaskReportId, report)
			 if err != nil {
				return errors.Wrap(err, "error occured")
			}
		*/
		return errors.Errorf("TASK_TYPE_CODE_NOT_IMPLEMENTED:%s", taskTypeCode)
	case base.TaskTypeCodeTechnicalSupport:
		return errors.Errorf("TASK_TYPE_CODE_NOT_IMPLEMENTED:%s", taskTypeCode)
	default:
		return errors.Errorf("TASK_TYPE_CODE_NOT_SUPPORTED:%s", taskTypeCode)
	}

	_ = tm.DoNotifyTaskCancelByCustomer(&aepr.Log, taskId)

	return nil
}

func (tm *TaskManagement) DoNotifyTaskCreate(l *log.DXLog, taskId int64) (err error) {
	var ok bool
	_, task, err := tm.Task.ShouldGetById(l, taskId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	customerSalesAreaCode, ok := task["customer_sales_area_code"].(string)
	if !ok {
		return errors.New("INVALID_SALES_AREA_CODE")
	}
	taskTypeId, ok := task["task_type_id"].(int64)
	if !ok {
		return errors.New("INVALID_TASK_TYPE_ID")
	}
	taskCode, ok := task["task_code"].(string)
	if !ok {
		return errors.New("INVALID_TASK_CODE")
	}

	_, fieldExecutorEffectiveAreas, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveArea.Select(l, nil, utils.JSON{
		"area_code": customerSalesAreaCode,
	}, nil, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	var templateTitle, templateBody string
	switch taskTypeId {
	case base.TaskTypeIdConstruction:
		_, templateTitle, templateBody, err = configuration_settings.ModuleConfigurationSettings.GeneralTemplateGetByNameId(l, "NEW_CONSTRUCTION_TASK_AVAILABLE")
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	default:
		return errors.New(fmt.Sprintf("INVALID_TASK_TYPE_ID:%d", taskTypeId))
	}

	for _, fieldExecutorEffectiveArea := range fieldExecutorEffectiveAreas {
		userId, ok := fieldExecutorEffectiveArea["user_id"].(int64)
		if !ok {
			return errors.New("INVALID_USER_ID")
		}
		_ = user_management.ModuleUserManagement.UserMessageCreateAllApplication(l, userId, templateTitle, templateBody, utils.JSON{
			"task_code": taskCode,
		}, map[string]string{
			"task_id": fmt.Sprintf("%d", taskId),
		})
	}

	return nil
}

func (tm *TaskManagement) DoNotifySubTaskVerificationFailed(l *log.DXLog, subTaskId int64) (err error) {
	_, subTasks, err := ModuleTaskManagement.SubTask.Select(l, nil, utils.JSON{
		"id": subTaskId,
	}, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, templateTitle, templateBody, err := configuration_settings.ModuleConfigurationSettings.GeneralTemplateGetByNameId(l, "SUB_TASK_FIELD_SUPERVISOR_VERIFICATION_FAILED")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	for _, subTask := range subTasks {
		subTaskCode, ok := subTask["code"].(string)
		if !ok {
			return errors.New("INVALID_SUB_TASK_CODE")
		}
		subTaskLastFieldExecutorUserId, ok := subTask["last_field_executor_user_id"].(int64)
		if !ok {
			return errors.New("INVALID_LAST_FIELD_EXECUTOR_USER_ID")
		}

		_ = user_management.ModuleUserManagement.UserMessageCreateAllApplication(l, subTaskLastFieldExecutorUserId, templateTitle, templateBody, utils.JSON{
			"sub_task_code": subTaskCode,
		}, map[string]string{
			"sub_task_id": fmt.Sprintf("%d", subTaskId),
		})
	}

	return nil
}

func (tm *TaskManagement) DoNotifyTaskCancelByCustomer(l *log.DXLog, taskId int64) (err error) {
	_, subTasks, err := ModuleTaskManagement.SubTask.Select(l, nil, utils.JSON{
		"task_id": taskId,
	}, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, templateTitle, templateBody, err := configuration_settings.ModuleConfigurationSettings.GeneralTemplateGetByNameId(l, "TASK_CANCEL_BY_CUSTOMER")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	for _, subTask := range subTasks {
		subTaskCode, ok := subTask["code"].(string)
		if !ok {
			return errors.New("INVALID_SUB_TASK_CODE")
		}
		subTaskLastFieldExecutorUserId, ok := subTask["last_field_executor_user_id"].(int64)
		if !ok {
			return errors.New("INVALID_LAST_FIELD_EXECUTOR_USER_ID")
		}

		_ = user_management.ModuleUserManagement.UserMessageCreateAllApplication(l, subTaskLastFieldExecutorUserId, templateTitle, templateBody, utils.JSON{
			"sub_task_code": subTaskCode,
		}, map[string]string{
			"task_id": fmt.Sprintf("%d", taskId),
		})
	}

	return nil
}

func (tm *TaskManagement) ValidationTaskCodeShouldNotAlreadyExist(l *log.DXLog, newCode string) (err error) {
	_, task, err := tm.Task.GetByNameId(l, newCode)
	if err != nil {
		return errors.New(fmt.Sprintf("TASK_CODE_VALIDATION_FAILED:%s", newCode))
	}
	if task != nil {
		return errors.New(fmt.Sprintf("TASK_CODE_ALREADY_EXISTS:%s", newCode))
	}
	return nil
}

func (tm *TaskManagement) TxValidationTaskCodeShouldNotAlreadyExist(tx *database.DXDatabaseTx, newCode string) (err error) {
	_, task, err := tm.Task.TxGetByNameId(tx, newCode)
	if err != nil {
		return errors.New(fmt.Sprintf("TASK_CODE_VALIDATION_FAILED:%s", newCode))
	}
	if task != nil {
		return errors.New(fmt.Sprintf("TASK_CODE_ALREADY_EXISTS:%s", newCode))
	}
	return nil
}

func (tm *TaskManagement) ValidationTaskCodeShouldNotAlreadyExistExceptSelf(aepr *api.DXAPIEndPointRequest, newCode string, taskId int64) (err error) {
	_, task, err := tm.Task.ShouldGetByNameId(&aepr.Log, newCode)
	if err != nil {
		aepr.Log.Errorf(err, "ERROR: Task.GetById: %s", err.Error())
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "TASK_CODE_VALIDATION_FAILED:%s", newCode)
	}
	if task != nil {
		if task["id"].(int64) != taskId {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "TASK_CODE_ALREADY_EXISTS:%s", newCode)
		}
	}
	return nil
}

func (tm *TaskManagement) TaskTxCreateConstruction(tx *database.DXDatabaseTx, code string, customerId int64, data1, data2 string) (taskId int64, err error) {
	taskTypeId := base.TaskTypeIdConstruction
	u := utils.JSON{}

	// Only validate code if it's not empty
	if code != "" {
		err = tm.TxValidationTaskCodeShouldNotAlreadyExist(tx, code)
		if err != nil {
			return 0, err
		}
		u["code"] = code
	}

	_, subTaskTypeConstructionSK, err := ModuleTaskManagement.SubTaskType.TxShouldGetByNameId(tx, base.SubTaskTypeFullCodeConstructionSK)
	if err != nil {
		return 0, err
	}
	subTaskTypeIdConstructionSK := subTaskTypeConstructionSK["id"].(int64)

	_, subTaskTypeConstructionSR, err := ModuleTaskManagement.SubTaskType.TxShouldGetByNameId(tx, base.SubTaskTypeFullCodeConstructionSR)
	if err != nil {
		return 0, err
	}
	subTaskTypeIdConstructionSR := subTaskTypeConstructionSR["id"].(int64)

	_, subTaskTypeConstructionMeterInstallation, err := ModuleTaskManagement.SubTaskType.TxShouldGetByNameId(tx, base.SubTaskTypeFullCodeConstructionMeterInstallation)
	if err != nil {
		return 0, err
	}
	subTaskTypeIdConstructionMeterInstallation := subTaskTypeConstructionMeterInstallation["id"].(int64)

	_, subTaskTypeConstructionGasIn, err := ModuleTaskManagement.SubTaskType.TxShouldGetByNameId(tx, base.SubTaskTypeFullCodeConstructionGasIn)
	if err != nil {
		return 0, err
	}
	subTaskTypeIdConstructionGasIn := subTaskTypeConstructionGasIn["id"].(int64)

	u["task_type_id"] = taskTypeId
	u["status"] = base.TaskStatusWaitingAssignment
	u["customer_id"] = customerId
	u["data1"] = data1
	u["data2"] = data2

	taskId, err = ModuleTaskManagement.Task.TxInsert(tx, u)
	if err != nil {
		return 0, err
	}

	_, err = ModuleTaskManagement.SubTask.TxInsert(tx, utils.JSON{
		"task_id":          taskId,
		"sub_task_type_id": subTaskTypeIdConstructionSK,
		"status":           base.SubTaskStatusWaitingAssignment,
	})
	if err != nil {
		return 0, err
	}

	_, err = ModuleTaskManagement.SubTask.TxInsert(tx, utils.JSON{
		"task_id":          taskId,
		"sub_task_type_id": subTaskTypeIdConstructionSR,
		"status":           base.SubTaskStatusWaitingAssignment,
	})
	if err != nil {
		return 0, err
	}

	_, err = ModuleTaskManagement.SubTask.TxInsert(tx, utils.JSON{
		"task_id":          taskId,
		"sub_task_type_id": subTaskTypeIdConstructionMeterInstallation,
		"status":           base.SubTaskStatusBlockingDependency,
	})
	if err != nil {
		return 0, err
	}

	_, err = ModuleTaskManagement.SubTask.TxInsert(tx, utils.JSON{
		"task_id":          taskId,
		"sub_task_type_id": subTaskTypeIdConstructionGasIn,
		"status":           base.SubTaskStatusBlockingDependency,
	})
	if err != nil {
		return 0, err
	}

	return taskId, nil
}

func (tm *TaskManagement) TxSubTaskReportCreate(dtx *database.DXDatabaseTx, at time.Time,
	userId int64,
	userUid string,
	userLoginId string,
	userFullname string,
	userPhoneNumber string,
	organizationId int64,
	organizationUid string,
	organizationName string,
	subTaskId int64,
	subTaskUid string,
	subTaskStatus string,
	report map[string]any) (subTaskReportId int64, err error) {

	reportAsString, err := utils.JSONToString(report)
	if err != nil {
		return 0, err
	}

	subTaskReportId, err = tm.SubTaskReport.TxInsert(dtx, utils.JSON{
		"sub_task_id":       subTaskId,
		"sub_task_uid":      subTaskUid,
		"sub_task_status":   subTaskStatus,
		"timestamp":         at,
		"user_id":           userId,
		"user_uid":          userUid,
		"user_loginid":      userLoginId,
		"user_fullname":     userFullname,
		"user_phonenumber":  userPhoneNumber,
		"organization_id":   organizationId,
		"organization_uid":  organizationUid,
		"organization_name": organizationName,
		"report":            reportAsString,
	})
	if err != nil {
		return 0, err
	}

	formattedDate := at.Format("20060102")
	subTaskReportCode := fmt.Sprintf("%04d%s", subTaskReportId, formattedDate)

	_, err = tm.SubTaskReport.TxUpdate(dtx, utils.JSON{
		"code": subTaskReportCode,
	}, utils.JSON{
		"id": subTaskReportId,
	})
	if err != nil {
		return subTaskReportId, err
	}

	return subTaskReportId, nil
}

func (tm *TaskManagement) DoCustomerConstructionTaskCreate(l *log.DXLog, data map[string]interface{}) (taskId int64, err error) {
	// Extract required fields
	taskCode, _ := data["task_code"].(string) // Make task code optional

	registrationNumber, ok := data["registration_number"].(string)
	if !ok {
		return 0, errors.New("REGISTRATION_NUMBER_REQUIRED")
	}

	// Get customer by registration number
	_, customer, err := tm.Customer.ShouldGetByNameId(l, registrationNumber)
	if err != nil {
		return 0, errors.Wrap(err, "CUSTOMER_NOT_FOUND")
	}
	customerId := customer["id"].(int64)

	// Validate task code doesn't exist only if provided
	if taskCode != "" {
		err = tm.ValidationTaskCodeShouldNotAlreadyExist(l, taskCode)
		if err != nil {
			return 0, err
		}
	}

	// Create task in transaction
	db := database.Manager.Databases[tm.DatabaseNameId]
	err = db.Tx(l, database.LevelReadCommitted, func(tx *database.DXDatabaseTx) (err2 error) {
		taskId, err2 = tm.TaskTxCreateConstruction(tx, taskCode, customerId, "", "")
		if err2 != nil {
			return errors.Wrap(err2, "error occurred")
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return taskId, nil
}
