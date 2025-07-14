package handler

import (
	"database/sql"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"net/http"
	"time"
)

func SubTaskAssign(aepr *api.DXAPIEndPointRequest) (err error) {
	adminUserId := aepr.LocalData["user_id"].(int64)
	adminUser := aepr.LocalData["user"].(utils.JSON)

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

	fieldExecutorUserId := aepr.ParameterValues["field_executor_user_id"].Value.(int64)
	at := aepr.ParameterValues["at"].Value.(time.Time)

	report := map[string]any{
		"assigned_by_user_id":       adminUserId,
		"assigned_by_user_loginid":  adminUser["loginid"],
		"assigned_by_user_fullname": adminUser["fullname"],
		"assigned_to_user_id":       fieldExecutorUserId,
	}

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, fieldExecutorUserId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldExecutor,
		SubTaskStatusCondition: []string{base.SubTaskStatusWaitingAssignment},
		NewSubTaskStatus:       base.SubTaskStatusAssigned,
		OperationName:          base.AdminOperationNameSubTaskAssign,
		Report:                 report,
		OnExecute:              task_management.ModuleTaskManagement.OnSubTaskPick,
	}, true)

	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func SubTaskReplaceFieldExecutor(aepr *api.DXAPIEndPointRequest) (err error) {
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

	_, newFieldExecutorUserId, err := aepr.GetParameterValueAsInt64("new_field_executor_user_id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, at, err := aepr.GetParameterValueAsTime("at")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	userId := aepr.LocalData["user_id"].(int64)
	userUid := aepr.LocalData["user_uid"].(string)
	_, user, err := user_management.ModuleUserManagement.User.ShouldGetById(&aepr.Log, userId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	userLoginId, ok := user["loginid"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:USER_LOGINID_NOT_FOUND")
	}

	userFullName, ok := user["fullname"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:USER_FULLNAME_NOT_FOUND")
	}

	userPhoneNumber, ok := user["phonenumber"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:USER_PHONENUMBER_NOT_FOUND")
	}
	userOrganizationId := aepr.LocalData["organization_id"].(int64)
	userOrganizationUid := aepr.LocalData["organization_uid"].(string)
	userOrganizationName := aepr.LocalData["organization_name"].(string)

	inputParameters := utils.JSON{
		"sub_task_id":                subTaskId,
		"new_field_executor_user_id": newFieldExecutorUserId,
		"at":                         at,
	}

	dbTaskDispatcher := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]
	var dtx *database.DXDatabaseTx
	dtx, err = dbTaskDispatcher.TransactionBegin(sql.LevelReadCommitted)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	defer dtx.Finish(&aepr.Log, err)

	subTaskStatus := subTask["status"].(string)
	oldFieldExecutorUserId := subTask["last_field_executor_user_id"].(int64)

	switch subTaskStatus {
	case
		base.SubTaskStatusBlockingDependency,
		base.SubTaskStatusWaitingAssignment,
		base.SubTaskStatusWaitingVerification,
		base.SubTaskStatusVerificationSuccess:
		err = aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "INVALID_SUB_TASK_STATUS_FOR_OPERATION: %s", subTaskStatus)
		return errors.Wrap(err, "error occured")
	}

	_, newFieldExecutor, err := partner_management.ModulePartnerManagement.FieldExecutor.TxShouldSelectOne(dtx, utils.JSON{
		"user_id": newFieldExecutorUserId,
	}, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	newFieldExecutorUserUid, ok := newFieldExecutor["user_uid"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:USER_UID_NOT_FOUND_ON_FIELD_EXECUTOR")
	}
	newFieldExecutorUserLoginId, ok := newFieldExecutor["user_loginid"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:USER_LOGINID_NOT_FOUND_ON_FIELD_EXECUTOR")
	}
	newFieldExecutorUserFullName, ok := newFieldExecutor["user_fullname"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:USER_FULLNAME_NOT_FOUND_ON_FIELD_EXECUTOR")
	}
	newFieldExecutorUserPhoneNumber, ok := newFieldExecutor["user_phonenumber"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:USER_PHONENUMBER_NOT_FOUND_ON_FIELD_EXECUTOR")
	}
	newFieldExecutorOrganizationId, ok := newFieldExecutor["organization_id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:ORGANIZATION_ID_NOT_FOUND_ON_FIELD_EXECUTOR")
	}
	newFieldExecutorOrganizationUid, ok := newFieldExecutor["organization_uid"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:ORGANIZATION_UID_NOT_FOUND_ON_FIELD_EXECUTOR")
	}
	newFieldExecutorOrganizationName, ok := newFieldExecutor["organization_name"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:ORGANIZATION_NAME_NOT_FOUND_ON_FIELD_EXECUTOR")
	}

	report := utils.JSON{
		"sub_task_id":                          subTaskId,
		"sub_task_uid":                         subTaskUid,
		"sub_task_status":                      subTaskStatus,
		"old_field_executor_user_id":           oldFieldExecutorUserId,
		"new_field_executor_user_id":           newFieldExecutorUserId,
		"new_field_executor_user_uid":          newFieldExecutorUserUid,
		"new_field_executor_user_loginid":      newFieldExecutorUserLoginId,
		"new_field_executor_user_fullname":     newFieldExecutorUserFullName,
		"new_field_executor_user_phonenumber":  newFieldExecutorUserPhoneNumber,
		"new_field_executor_organization_id":   newFieldExecutorOrganizationId,
		"new_field_executor_organization_uid":  newFieldExecutorOrganizationUid,
		"new_field_executor_organization_name": newFieldExecutorOrganizationName,
		"operation_by_user_id":                 userId,
	}

	subTaskReportId, err := task_management.ModuleTaskManagement.TxSubTaskReportCreate(dtx, at,
		userId,
		userUid,
		userLoginId,
		userFullName,
		userPhoneNumber,
		userOrganizationId,
		userOrganizationUid,
		userOrganizationName,
		subTaskId,
		subTaskUid,
		subTaskStatus, report)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.TxShouldGetById(dtx, subTaskReportId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskReportUid := subTaskReport["uid"].(string)

	report["sub_task_report_id"] = subTaskReportId
	report["sub_task_report_uid"] = subTaskReportUid

	_, err = task_management.ModuleTaskManagement.SubTask.TxUpdate(dtx, utils.JSON{
		"last_sub_task_report_id":      subTaskReportId,
		"last_sub_task_report_uid":     subTaskReportUid,
		"status":                       subTaskStatus,
		"last_field_executor_user_id":  newFieldExecutorUserId,
		"last_field_executor_user_uid": newFieldExecutorUserUid,
	}, utils.JSON{
		"id": subTaskId,
	})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, subTask, err = task_management.ModuleTaskManagement.SubTask.TxShouldGetById(dtx, subTaskId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	inputParametersAsString, err := utils.JSONToString(inputParameters)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, err = task_management.ModuleTaskManagement.SubTaskHistoryItem.TxInsert(dtx, utils.JSON{
		"sub_task_id":         subTaskId,
		"sub_task_uid":        subTaskUid,
		"from_status":         subTaskStatus,
		"to_status":           subTaskStatus,
		"timestamp":           at,
		"user_id":             userId,
		"user_uid":            userUid,
		"operation_nameid":    base.AdminOperationNameSubTaskReplaceFieldExecutor,
		"input_parameters":    inputParametersAsString,
		"sub_task_report_id":  subTaskReportId,
		"sub_task_report_uid": subTaskReportUid,
	})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	r := utils.JSON{"data": utils.JSON{
		"status":              http.StatusText(http.StatusOK),
		"sub_task":            subTask,
		"sub_task_report_id":  subTaskReportId,
		"sub_task_report_uid": subTaskReportUid,
	}}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, r)
	return nil

}
