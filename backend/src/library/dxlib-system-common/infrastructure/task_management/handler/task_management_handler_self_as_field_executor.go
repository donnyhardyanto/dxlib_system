package handler

import (
	"github.com/donnyhardyanto/dxlib/database/protected/db"
	dbUtils "github.com/donnyhardyanto/dxlib/database2/utils"
	"github.com/donnyhardyanto/dxlib/utils/json"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/construction_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"net/http"
	"time"

	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
)

func ValidateFieldExecutorUserWithSubTask(aepr *api.DXAPIEndPointRequest, userId int64, subTask utils.JSON) (fieldExecutor utils.JSON, err error) {
	_, fieldExecutor, err = partner_management.ModulePartnerManagement.FieldExecutor.ShouldSelectOne(&aepr.Log, utils.JSON{
		"user_id":    userId,
		"is_deleted": false,
	}, nil, nil)
	if err != nil {
		return nil, aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "USER_IS_NOT_IN_FIELD_EXECUTOR_LIST")
	}
	subTaskId := subTask["id"].(int64)

	lastFieldExecutorUserId, ok := subTask["last_field_executor_user_id"].(int64)
	if !ok {
		return nil, aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "WRONG_FIELD_EXECUTOR_IS_NIL_FOR_SUB_TASK:%d", subTaskId)
	}
	if lastFieldExecutorUserId != userId {
		return nil, aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "WRONG_FIELD_EXECUTOR_FOR_SUB_TASK:(subTaskId=%d),(lastFieldExecutorUserId=%d)", subTaskId, lastFieldExecutorUserId)
	}
	return fieldExecutor, nil
}

func SelfAsFieldExecutorUserSubTaskPick(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.LocalData["user_id"].(int64)
	userUid := aepr.LocalData["user_uid"].(string)

	isExistAt, at, err := aepr.GetParameterValueAsTime("at")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistAt {
		at = time.Now()
	}

	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUBTASK_ID_NOT_FOUND")
	}
	subTaskStatus, ok := subTask["status"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUBTASK_STATUS_NOT_FOUND")
	}
	if subTaskStatus != base.SubTaskStatusWaitingAssignment {
		return aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "SUB_TASK_STATUS_NOT_WAITING_ASSIGNMENT:subTaskStatus=%s", subTaskStatus)
	}
	subTaskCustomerId, ok := subTask["customer_id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUBTASK_CUSTOMER_ID_NOT_FOUND")
	}
	subTaskTypeId, ok := subTask["sub_task_type_id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUBTASK_TYPE_ID_NOT_FOUND")
	}

	_, userFieldExecutorEffectiveExpertise, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveExpertise.SelectOne(&aepr.Log, nil, utils.JSON{
		"user_uid":         userUid,
		"sub_task_type_id": subTaskTypeId,
	}, nil, map[string]string{"id": "asc"})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if userFieldExecutorEffectiveExpertise == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "USER_FIELD_EXECUTOR_EFFECTIVE_EXPERTISE_NOT_MATCH_SUB_TASK:subTaskId=%d,userUid=%s", subTaskId, userUid)
	}

	_, customer, err := task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, subTaskCustomerId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	customerSalesAreaCode, ok := customer["sales_area_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:CUSTOMER_SALES_AREA_CODE_NOT_FOUND")
	}

	locationCodesToSearch := []string{}

	customerAddressProvinceLocationCode, ok := customer["address_province_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:CUSTOMER_ADDRESS_PROVINCE_LOCATION_CODE_NOT_FOUND")
	}
	locationCodesToSearch = append(locationCodesToSearch, customerAddressProvinceLocationCode)

	customerAddressKabupatenLocationCode, ok := customer["address_kabupaten_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:CUSTOMER_ADDRESS_KABUPATEN_LOCATION_CODE_NOT_FOUND")
	}
	locationCodesToSearch = append(locationCodesToSearch, customerAddressKabupatenLocationCode)

	customerAddressKecamatanLocationCode, ok := customer["address_kecamatan_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:CUSTOMER_ADDRESS_KECAMATAN_LOCATION_CODE_NOT_FOUND")
	}
	locationCodesToSearch = append(locationCodesToSearch, customerAddressKecamatanLocationCode)

	customerAddressKelurahanLocationCode, ok := customer["address_kelurahan_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:CUSTOMER_ADDRESS_KELURAHAN_LOCATION_CODE_NOT_FOUND")
	}
	locationCodesToSearch = append(locationCodesToSearch, customerAddressKelurahanLocationCode)

	_, userFieldExecutorEffectiveArea, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveArea.SelectOne(&aepr.Log, nil, utils.JSON{
		"user_uid":  userUid,
		"area_code": customerSalesAreaCode,
	}, nil, map[string]string{"id": "asc"})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if userFieldExecutorEffectiveArea == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "USER_FIELD_EXECUTOR_EFFECTIVE_AREA_NOT_MATCH_SUB_TASK:subTaskId=%d,userUid=%s", subTaskId, userUid)
	}

	c1 := dbUtils.SQLBuildWhereInClause("location_code", locationCodesToSearch)
	_, userFieldExecutorEffectiveLocation, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveLocation.SelectOne(&aepr.Log, nil, utils.JSON{
		"user_uid": userUid,
		"c1":       db.SQLExpression{Expression: c1},
	}, nil, map[string]string{"id": "asc"})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if userFieldExecutorEffectiveLocation == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "USER_FIELD_EXECUTOR_EFFECTIVE_LOCATION_NOT_MATCH_SUB_TASK:subTaskId=%d,userUid=%s", subTaskId, userUid)
	}

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldExecutor,
		SubTaskStatusCondition: []string{base.SubTaskStatusWaitingAssignment},
		NewSubTaskStatus:       base.SubTaskStatusAssigned,
		OperationName:          base.UserAsFieldExecutorOperationNameSubTaskPick,
		Report:                 nil,
		OnExecute:              task_management.ModuleTaskManagement.OnSubTaskPick,
	}, true)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func SelfAsFieldExecutorUserSubTaskWorkingStart(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	userId := aepr.LocalData["user_id"].(int64)

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldExecutor,
		SubTaskStatusCondition: []string{base.SubTaskStatusAssigned},
		NewSubTaskStatus:       base.SubTaskStatusWorking,
		OperationName:          base.UserAsFieldExecutorOperationNameSubTaskWorkingStart,
		Report:                 nil,
		OnExecute:              nil,
	}, true)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func CheckAndGetGasApplianceName(aepr *api.DXAPIEndPointRequest, effectiveReport utils.JSON) (resultEffectiveReport utils.JSON, err error) {
	reportGasAppliances, ok := effectiveReport["gas_appliances"].([]map[string]any)
	if !ok {
		return nil, aepr.WriteResponseAndNewErrorf(http.StatusExpectationFailed, "", "IMPOSSIBLE_CONDITION:report=gas_appliances_NOT_FOUND")
	}
	for i, reportGasAppliance := range reportGasAppliances {
		//	gasApplianceMap := gasAppliance
		//	if !ok {
		//		return aepr.WriteResponseAndNewErrorf(http.StatusExpectationFailed, "", "IMPOSSIBLE_CONDITION:report=gas_appliances[%d]=NOT_MAP", i)
		//	}
		gasApplianceId, err := json.GetInt64(reportGasAppliance, "id")
		if err != nil {
			return nil, aepr.WriteResponseAndNewErrorf(http.StatusExpectationFailed, "", "IMPOSSIBLE_CONDITION:report=gas_appliances[%d].id:%w", err)
		}
		_, gasAppliance, err := construction_management.ModuleConstructionManagement.GasAppliance.ShouldGetById(&aepr.Log, gasApplianceId)
		if err != nil {
			return nil, aepr.WriteResponseAndNewErrorf(http.StatusExpectationFailed, "", "IMPOSSIBLE_CONDITION:report=gas_appliances[%d].id=%d_NOT_FOUND", i, gasApplianceId)
		}
		gasApplianceName, ok := gasAppliance["name"].(string)
		if !ok {
			return nil, aepr.WriteResponseAndNewErrorf(http.StatusExpectationFailed, "", "IMPOSSIBLE_CONDITION:report=gas_appliances[%d].name=NOT_STRING", i)
		}
		reportGasAppliances[i]["name"] = gasApplianceName
	}
	effectiveReport["gas_appliances"] = reportGasAppliances
	return effectiveReport, nil
}
func GetEffectiveReport(aepr *api.DXAPIEndPointRequest, subTaskId int64, report map[string]any) (effectiveReport utils.JSON, err error) {
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetById(&aepr.Log, subTaskId)
	if err != nil {
		return nil, err
	}

	subTaskTypeFullCode := subTask["sub_task_type_full_code"].(string)
	var ok bool
	switch subTaskTypeFullCode {
	case base.SubTaskTypeFullCodeConstructionGasIn:
		effectiveReport, ok = report["gas_in"].(map[string]any)
		if !ok {
			return nil, aepr.WriteResponseAndNewErrorf(http.StatusExpectationFailed, "IMPOSSIBLE_CONDITION:subTaskTypeFullCode=%s", subTaskTypeFullCode)
		}
		effectiveReport, err = CheckAndGetGasApplianceName(aepr, effectiveReport)
		if err != nil {
			return nil, aepr.WriteResponseAndNewErrorf(http.StatusExpectationFailed, "CHECK_REPORT_GAS_APPLIANCE:subTaskTypeFullCode=%s:%w", subTaskTypeFullCode, err)
		}
	case base.SubTaskTypeFullCodeConstructionMeterInstallation:
		effectiveReport, ok = report["meter_installation"].(map[string]any)
		if !ok {
			return nil, aepr.WriteResponseAndNewErrorf(http.StatusExpectationFailed, "IMPOSSIBLE_CONDITION:subTaskTypeFullCode=%s", subTaskTypeFullCode)
		}
	case base.SubTaskTypeFullCodeConstructionSK:
		effectiveReport, ok = report["sk"].(map[string]any)
		if !ok {
			return nil, aepr.WriteResponseAndNewErrorf(http.StatusExpectationFailed, "IMPOSSIBLE_CONDITION:subTaskTypeFullCode=%s", subTaskTypeFullCode)
		}
		effectiveReport, err = CheckAndGetGasApplianceName(aepr, effectiveReport)
		if err != nil {
			return nil, aepr.WriteResponseAndNewErrorf(http.StatusExpectationFailed, "CHECK_REPORT_GAS_APPLIANCE:subTaskTypeFullCode=%s:%w", subTaskTypeFullCode, err)
		}
	case base.SubTaskTypeFullCodeConstructionSR:
		effectiveReport, ok = report["sr"].(map[string]any)
		if !ok {
			return nil, aepr.WriteResponseAndNewErrorf(http.StatusExpectationFailed, "IMPOSSIBLE_CONDITION:subTaskTypeFullCode=%s", subTaskTypeFullCode)
		}
	default:
		return nil, aepr.WriteResponseAndNewErrorf(http.StatusExpectationFailed, "IMPOSSIBLE_CONDITION:subTaskTypeFullCode=%s", subTaskTypeFullCode)
	}

	return effectiveReport, nil
}

func SelfAsFieldExecutorUserSubTaskWorkingFinish(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	report := aepr.ParameterValues["report"].Value.(map[string]any)
	userId := aepr.LocalData["user_id"].(int64)

	effectiveReport, err := GetEffectiveReport(aepr, subTaskId, report)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldExecutor,
		SubTaskStatusCondition: []string{base.SubTaskStatusWorking},
		NewSubTaskStatus:       base.SubTaskStatusWaitingVerification,
		OperationName:          base.UserAsFieldExecutorOperationNameSubTaskWorkingFinish,
		Report:                 effectiveReport,
		OnExecute:              task_management.ModuleTaskManagement.OnSubTaskWorkingFinish,
	}, true)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func SelfAsFieldExecutorSubTaskReworkStart(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	report := aepr.ParameterValues["report"].Value.(map[string]any)
	userId := aepr.LocalData["user_id"].(int64)

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldExecutor,
		SubTaskStatusCondition: []string{base.SubTaskStatusWaitingVerification},
		NewSubTaskStatus:       base.SubTaskStatusReworking,
		OperationName:          base.UserAsFieldExecutorOperationNameSubTaskReworkingStart,
		Report:                 report,
		OnExecute:              nil,
	}, true)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func SelfAsFieldExecutorSubTaskReworkFinish(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	report := aepr.ParameterValues["report"].Value.(map[string]any)
	userId := aepr.LocalData["user_id"].(int64)

	effectiveReport, err := GetEffectiveReport(aepr, subTaskId, report)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldExecutor,
		SubTaskStatusCondition: []string{base.SubTaskStatusReworking},
		NewSubTaskStatus:       base.SubTaskStatusWaitingVerification,
		OperationName:          base.UserAsFieldExecutorOperationNameSubTaskReworkingFinish,
		Report:                 effectiveReport,
		OnExecute:              nil,
	}, true)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func SelfAsFieldExecutorSubTaskReworkCancel(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	report := aepr.ParameterValues["report"].Value.(map[string]any)
	userId := aepr.LocalData["user_id"].(int64)

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldExecutor,
		SubTaskStatusCondition: []string{base.SubTaskStatusReworking},
		NewSubTaskStatus:       base.SubTaskStatusWaitingVerification,
		OperationName:          base.UserAsFieldExecutorOperationNameSubTaskReworkingCancel,
		Report:                 report,
		OnExecute:              nil,
	}, true)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func SelfAsFieldExecutorUserSubTaskFixingStart(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	userId := aepr.LocalData["user_id"].(int64)

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldExecutor,
		SubTaskStatusCondition: []string{base.SubTaskStatusVerificationFail, base.SubTaskStatusCGPVerificationFail},
		NewSubTaskStatus:       base.SubTaskStatusFixing,
		OperationName:          base.UserAsFieldExecutorOperationNameSubTaskFixingStart,
		Report:                 nil,
		OnExecute:              nil,
	}, true)
	return nil
}

func SelfAsFieldExecutorSubTaskFixingDone(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	report := aepr.ParameterValues["report"].Value.(map[string]any)
	userId := aepr.LocalData["user_id"].(int64)

	effectiveReport, err := GetEffectiveReport(aepr, subTaskId, report)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldExecutor,
		SubTaskStatusCondition: []string{base.SubTaskStatusFixing},
		NewSubTaskStatus:       base.SubTaskStatusWaitingVerification,
		OperationName:          base.UserAsFieldExecutorOperationNameSubTaskFixingFinish,
		Report:                 effectiveReport,
		OnExecute:              task_management.ModuleTaskManagement.OnSubTaskFixingFinish,
	}, true)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func SelfAsFieldExecutorSubTaskPause(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	report := aepr.ParameterValues["report"].Value.(map[string]any)
	userId := aepr.LocalData["user_id"].(int64)

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldExecutor,
		SubTaskStatusCondition: []string{base.SubTaskStatusWorking},
		NewSubTaskStatus:       base.SubTaskStatusPaused,
		OperationName:          base.UserAsFieldExecutorOperationNameSubTaskPause,
		Report:                 report,
		OnExecute:              task_management.ModuleTaskManagement.OnSubTaskPause,
	}, true)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func SelfAsFieldExecutorSubTaskResume(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	userId := aepr.LocalData["user_id"].(int64)

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldExecutor,
		SubTaskStatusCondition: []string{base.SubTaskStatusPaused},
		NewSubTaskStatus:       base.SubTaskStatusWorking,
		OperationName:          base.UserAsFieldExecutorOperationNameSubTaskResume,
		Report:                 nil,
		OnExecute:              nil,
	}, true)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func SelfAsFieldExecutorUserSubTaskCancelByFieldExecutor(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	report := aepr.ParameterValues["report"].Value.(map[string]any)
	userId := aepr.LocalData["user_id"].(int64)

	_, subTaskReportId, subTaskReportUid, err := task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldExecutor,
		SubTaskStatusCondition: []string{base.SubTaskStatusAssigned, base.SubTaskStatusWorking, base.SubTaskStatusPaused},
		NewSubTaskStatus:       base.SubTaskStatusCanceledByFieldExecutor,
		OperationName:          base.UserAsFieldExecutorOperationNameSubTaskCancelByFieldExecutor,
		Report:                 report,
		OnExecute:              task_management.ModuleTaskManagement.OnSubTaskCancelByFieldExecutor,
	}, false)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeNone,
		SubTaskStatusCondition: []string{base.SubTaskStatusCanceledByFieldExecutor},
		NewSubTaskStatus:       base.SubTaskStatusWaitingAssignment,
		OperationName:          base.AutoOperationNameSubTaskWaitingAssignment,
		Report:                 nil,
		OnExecute:              nil,
	}, false)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"sub_task_report_id":  subTaskReportId,
		"sub_task_report_uid": subTaskReportUid,
	}})

	return nil
}

func SelfAsFieldExecutorUserSubTaskCancelByCustomer(aepr *api.DXAPIEndPointRequest) (err error) {
	aepr.Log.Debugf("SelfAsFieldExecutorUserSubTaskCancelByCustomer")

	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	report := aepr.ParameterValues["report"].Value.(map[string]any)

	userId := aepr.LocalData["user_id"].(int64)
	_, err = ValidateFieldExecutorUserWithSubTask(aepr, userId, subTask)
	if err != nil {
		return err
	}

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:     aepr,
		UserType: base.UserTypeAny,
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

// penanganan piutang

func SelfAsFieldExecutorUserSubTaskCancel(aepr *api.DXAPIEndPointRequest) (err error) {
	aepr.Log.Debugf("SelfAsFieldExecutorUserSubTaskCancel")
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}
	at := aepr.ParameterValues["at"].Value.(time.Time)
	report := aepr.ParameterValues["report"].Value.(map[string]any)
	rpt, rpt_ok := report["penanganan_piutang"].(map[string]any)
	NewSubTaskStatus := ""
	if rpt_ok {
		reasonType := rpt["reason_type"].(string)
		switch reasonType {
		case "01":
			NewSubTaskStatus = base.SubTaskStatusCanceledByPaid
		case "02":
			NewSubTaskStatus = base.SubTaskStatusCanceledByCustomer
		case "03":
			NewSubTaskStatus = base.SubTaskStatusCanceledByForceMajeure
		case "04":
			NewSubTaskStatus = base.SubTaskStatusCanceledByOther

		}
	}
	userId := aepr.LocalData["user_id"].(int64)
	_, err = ValidateFieldExecutorUserWithSubTask(aepr, userId, subTask)
	if err != nil {
		return err
	}
	SubTaskStatusCondition := []string{}
	SubTaskTypeId, ok := subTask["sub_task_type_id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}
	switch SubTaskTypeId {
	case 5, 6, 7, 8:
		SubTaskStatusCondition = []string{
			base.SubTaskStatusAssigned,
			base.SubTaskStatusWorking,
		}
	}

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeAny,
		SubTaskStatusCondition: SubTaskStatusCondition,
		NewSubTaskStatus:       NewSubTaskStatus,
		OperationName:          base.UserAsFieldExecutorOperationNameSubTaskCancel,
		Report:                 report,
		OnExecute:              task_management.ModuleTaskManagement.OnSubTaskCancel,
	}, true)

	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}
