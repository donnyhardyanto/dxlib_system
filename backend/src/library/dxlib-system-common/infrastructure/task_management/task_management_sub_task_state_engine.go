package task_management

import (
	"database/sql"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/database/protected/db"
	dbUtils "github.com/donnyhardyanto/dxlib/database2/utils"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/module/general"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"net/http"
	"slices"
	"time"
)

func validateFieldExecutor(aepr *api.DXAPIEndPointRequest, userId int64, subTask utils.JSON) (fieldExecutor utils.JSON, err error) {
	_, fieldExecutor, err = partner_management.ModulePartnerManagement.FieldExecutor.ShouldSelectOne(&aepr.Log, utils.JSON{
		"user_id":    userId,
		"is_deleted": false,
	}, nil, nil)
	if err != nil {
		return nil, aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "USER_IS_NOT_IN_FIELD_EXECUTOR_LIST")
	}
	subTaskId := subTask["id"].(int64)
	subTaskStatus := subTask["status"].(string)
	switch subTaskStatus {
	case base.SubTaskStatusWaitingAssignment, base.SubTaskStatusBlockingDependency:
	default:
		lastFieldExecutorUserId, ok := subTask["last_field_executor_user_id"].(int64)
		if !ok {
			return nil, aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "WRONG_FIELD_EXECUTOR_IS_NIL_FOR_SUB_TASK:%d", subTaskId)
		}
		if lastFieldExecutorUserId != userId {
			return nil, aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "WRONG_FIELD_EXECUTOR_FOR_SUB_TASK:(subTaskId=%d),(lastFieldExecutorUserId=%d)", subTaskId, lastFieldExecutorUserId)
		}
	}
	return fieldExecutor, nil
}

func validateCGP(aepr *api.DXAPIEndPointRequest, userId int64, subTask utils.JSON) (err error) {
	// check role
	return nil
}

func validateFieldSupervisor(aepr *api.DXAPIEndPointRequest, userId int64, subTask utils.JSON) (userRoleMembership utils.JSON, err error) {

	roleIdFieldSupervisor, err := general.ModuleGeneral.Property.GetAsInt64(&aepr.Log, base.ConfigRoleIdFieldSupervisor)
	if err != nil {
		return nil, err
	}

	_, userRoleMembership, err = partner_management.ModulePartnerManagement.FieldSupervisor.ShouldSelectOne(&aepr.Log, utils.JSON{
		"user_id": userId,
		"role_id": roleIdFieldSupervisor,
	}, nil, map[string]string{"organization_id": "asc"})
	if err != nil {
		return nil, aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "USER_IS_NOT_IN_FIELD_SUPERVISOR_LIST")
	}

	return userRoleMembership, nil
}

// processSubTaskRaw is the combined implementation for all user types
func (tm *TaskManagement) ProcessSubTaskStateRaw(subTaskId int64, userId int64, at time.Time,
	d StateEngineStructSubTaskStatus, isSendResponse bool) (subTask utils.JSON, subTaskReportId int64, subTaskReportUid string, err error) {
	//	userId := d.Aepr.LocalData["user_id"].(int64)
	var ok bool
	updateSubTaskJSON := utils.JSON{}

	// Use the passed validation function
	_, subTask, err = tm.SubTask.ShouldGetById(&d.Aepr.Log, subTaskId)
	if err != nil {
		return nil, 0, "", err
	}

	subTaskUid := subTask["uid"].(string)

	userUid := ""
	userLoginId := "SYSTEM"
	userFullName := "SYSTEM"
	userPhoneNumber := ""
	organizationId := int64(0)
	organizationUid := ""
	organizationName := "SYSTEM"

	if userId != 0 {
		_, user, err := user_management.ModuleUserManagement.User.ShouldGetById(&d.Aepr.Log, userId)
		if err != nil {
			return nil, 0, "", err
		}
		if user["is_deleted"].(bool) {
			return nil, 0, "", d.Aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "USER_IS_DELETED")
		}
		if user["status"].(string) != user_management.UserStatusActive {
			return nil, 0, "", d.Aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "USER_IS_NOT_ACTIVE")
		}
		userUid = user["uid"].(string)
		userLoginId = user["loginid"].(string)
		userFullName = user["fullname"].(string)
	}
	switch d.UserType {
	case base.UserTypeCGP:
		err = validateCGP(d.Aepr, userId, subTask)
		if err != nil {
			return subTask, 0, "", err
		}
	case base.UserTypeFieldSupervisor:
		_, err = validateFieldSupervisor(d.Aepr, userId, subTask)
		if err != nil {
			return subTask, 0, "", err
		}
	case base.UserTypeFieldExecutor:
		fieldExecutor, err := validateFieldExecutor(d.Aepr, userId, subTask)
		if err != nil {
			return subTask, 0, "", err
		}

		userLoginId, ok = fieldExecutor["user_loginid"].(string)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:USER_LOGINID_NOT_FOUND_ON_FIELD_EXECUTOR")
		}
		userFullName, ok = fieldExecutor["user_fullname"].(string)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:USER_FULLNAME_NOT_FOUND_ON_FIELD_EXECUTOR")
		}
		userPhoneNumber, ok = fieldExecutor["user_phonenumber"].(string)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:USER_PHONENUMBER_NOT_FOUND_ON_FIELD_EXECUTOR")
		}
		organizationId, ok = fieldExecutor["organization_id"].(int64)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:ORGANIZATION_ID_NOT_FOUND_ON_FIELD_EXECUTOR")
		}
		organizationUid, ok = fieldExecutor["organization_uid"].(string)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:ORGANIZATION_UID_NOT_FOUND_ON_FIELD_EXECUTOR")
		}
		organizationName, ok = fieldExecutor["organization_name"].(string)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:ORGANIZATION_NAME_NOT_FOUND_ON_FIELD_EXECUTOR")
		}
	case base.UserTypeNone, base.UserTypeAny:
		// No validation required
	}

	isPickSubTask := d.OperationName == base.UserAsFieldExecutorOperationNameSubTaskPick
	isolationLevel := sql.LevelReadCommitted
	if isPickSubTask {
		isolationLevel = sql.LevelSerializable
	}
	// Begin transaction
	dbTaskDispatcher := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]

	err = dbTaskDispatcher.Tx(&d.Aepr.Log, isolationLevel, func(dtx *database.DXDatabaseTx) (err2 error) {
		// Get current task status
		p := utils.JSON{
			"id": subTaskId,
		}
		if d.SubTaskStatusCondition != nil {
			p["c1"] = db.SQLExpression{Expression: dbUtils.SQLBuildWhereInClause("status", d.SubTaskStatusCondition)}
		}

		var oldSubTask utils.JSON
		_, oldSubTask, err2 = tm.SubTask.TxShouldSelectOneForUpdate(dtx, p, nil)
		if err2 != nil {
			if isPickSubTask {
				d.Aepr.Log.Warnf("SUB_TASK_PICK_FAILED: %v", err2)
				d.Aepr.WriteResponseAsJSON(http.StatusConflict, nil, utils.JSON{
					"result": "FAILED",
				})
				return err2
			}
			_ = d.Aepr.WriteResponseAndNewErrorf(http.StatusConflict, "INVALID_SUB_TASK_STATUS_FOR_OPERATION", "INVALID_SUB_TASK_STATUS_FOR_OPERATION")
			return err2
		}

		oldSubTaskStatus := oldSubTask["status"].(string)

		// Handle Report creation
		if d.Report != nil {
			subTaskReportId, err2 = tm.TxSubTaskReportCreate(dtx, at,
				userId,
				userUid,
				userLoginId,
				userFullName,
				userPhoneNumber,
				organizationId,
				organizationUid,
				organizationName,
				subTaskId,
				subTaskUid,
				d.NewSubTaskStatus,
				d.Report)
			if err2 != nil {
				return d.Aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "SUB_TASK_REPORT_CREATION_FAILED:%w", err2)
			}
			_, subTaskReport, err := tm.SubTaskReport.TxShouldGetById(dtx, subTaskReportId)
			if err != nil {
				return d.Aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "SUB_TASK_GET_REPORT_FAILED:%w", err)
			}
			subTaskReportUid, ok = subTaskReport["uid"].(string)
			if !ok {
				return d.Aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "REPORT_UID_NOT_FOUND:%w", err)
			}

			d.Report["sub_task_report_id"] = subTaskReportId
			d.Report["sub_task_report_uid"] = subTaskReportUid
			updateSubTaskJSON["last_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_sub_task_report_uid"] = subTaskReportUid

		}

		// Set new status
		updateSubTaskJSON["status"] = d.NewSubTaskStatus

		switch d.UserType {
		case base.UserTypeCGP:
			updateSubTaskJSON["last_cgp_user_id"] = userId
			updateSubTaskJSON["last_cgp_user_uid"] = userUid
			updateSubTaskJSON["last_cgp_user_loginid"] = userLoginId
			updateSubTaskJSON["last_cgp_user_fullname"] = userFullName
		case base.UserTypeFieldExecutor:
			updateSubTaskJSON["last_field_executor_user_id"] = userId
			updateSubTaskJSON["last_field_executor_user_uid"] = userUid
			updateSubTaskJSON["last_field_executor_user_loginid"] = userLoginId
			updateSubTaskJSON["last_field_executor_user_fullname"] = userFullName
		case base.UserTypeFieldSupervisor:
			updateSubTaskJSON["last_field_supervisor_user_id"] = userId
			updateSubTaskJSON["last_field_supervisor_user_uid"] = userUid
			updateSubTaskJSON["last_field_supervisor_user_loginid"] = userLoginId
			updateSubTaskJSON["last_field_supervisor_user_fullname"] = userFullName
		case base.UserTypeNone, base.UserTypeAny:
		default:
			return d.Aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "USER_TYPE_NOT_SUPPORTED:%s", d.UserType)
		}

		switch oldSubTaskStatus {
		case base.SubTaskStatusWorking:
			updateSubTaskJSON["working_end_at"] = at
			if d.Report != nil {
				updateSubTaskJSON["last_working_end_sub_task_report_id"] = subTaskReportId
				updateSubTaskJSON["last_working_end_sub_task_report_uid"] = subTaskReportUid
			}
		case base.SubTaskStatusReworking:
			updateSubTaskJSON["last_reworking_end_at"] = at
			if d.Report != nil {
				updateSubTaskJSON["last_reworking_end_sub_task_report_id"] = subTaskReportId
				updateSubTaskJSON["last_reworking_end_sub_task_report_uid"] = subTaskReportUid
			}
		case base.SubTaskStatusFixing:
			updateSubTaskJSON["last_fixing_end_at"] = at
			if d.Report != nil {
				updateSubTaskJSON["last_fixing_end_sub_task_report_id"] = subTaskReportId
				updateSubTaskJSON["last_fixing_end_sub_task_report_uid"] = subTaskReportUid
			}
		case base.SubTaskStatusPaused:
			updateSubTaskJSON["last_end_pause_at"] = at
		}

		switch d.NewSubTaskStatus {
		case base.SubTaskStatusWaitingAssignment:
			updateSubTaskJSON["last_field_executor_user_id"] = nil
			updateSubTaskJSON["last_field_executor_user_uid"] = nil
			updateSubTaskJSON["last_field_executor_user_loginid"] = nil
			updateSubTaskJSON["last_field_executor_user_fullname"] = nil
			updateSubTaskJSON["last_field_supervisor_user_id"] = nil
			updateSubTaskJSON["last_field_supervisor_user_uid"] = nil
			updateSubTaskJSON["last_field_supervisor_user_loginid"] = nil
			updateSubTaskJSON["last_field_supervisor_user_fullname"] = nil
		case base.SubTaskStatusWorking:
			updateSubTaskJSON["working_start_at"] = at
		case base.SubTaskStatusWaitingVerification:
			updateSubTaskJSON["is_working_finish"] = true
			updateSubTaskJSON["last_form_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_form_sub_task_report_uid"] = subTaskReportUid
		case base.SubTaskStatusVerificationSuccess:
			updateSubTaskJSON["last_verification_end_at"] = at
			updateSubTaskJSON["is_verification_success"] = true
			if d.Report != nil {
				updateSubTaskJSON["last_verification_sub_task_report_id"] = subTaskReportId
				updateSubTaskJSON["last_verification_sub_task_report_uid"] = subTaskReportUid
			}
		case base.SubTaskStatusVerificationFail:
			updateSubTaskJSON["last_verification_end_at"] = at
			updateSubTaskJSON["is_verification_success"] = false
			if d.Report != nil {
				updateSubTaskJSON["last_verification_sub_task_report_id"] = subTaskReportId
				updateSubTaskJSON["last_verification_sub_task_report_uid"] = subTaskReportUid
			}
		case base.SubTaskStatusFixing:
			fixCount := oldSubTask["fix_count"].(int64)
			updateSubTaskJSON["fix_count"] = fixCount + 1
			if fixCount == 1 {
				updateSubTaskJSON["first_fixing_start_at"] = at
			}
		case base.SubTaskStatusPaused:
			updateSubTaskJSON["last_start_pause_at"] = at
			if d.Report != nil {
				updateSubTaskJSON["last_pause_sub_task_report_id"] = subTaskReportId
				updateSubTaskJSON["last_pause_sub_task_report_uid"] = subTaskReportUid
			}
		case base.SubTaskStatusCanceledByFieldExecutor:
			updateSubTaskJSON["last_canceled_by_field_executor_at"] = at
			if d.Report != nil {
				updateSubTaskJSON["last_canceled_by_field_executor_sub_task_report_id"] = subTaskReportId
				updateSubTaskJSON["last_canceled_by_field_executor_sub_task_report_uid"] = subTaskReportUid
			}
		case base.SubTaskStatusCanceledByCustomer:
			updateSubTaskJSON["last_canceled_by_customer_at"] = at
			updateSubTaskJSON["completed_at"] = at
			if d.Report != nil {
				updateSubTaskJSON["last_canceled_by_customer_sub_task_report_id"] = subTaskReportId
				updateSubTaskJSON["last_canceled_by_customer_sub_task_report_uid"] = subTaskReportUid
			}
		case base.SubTaskStatusCGPVerificationSuccess:
			updateSubTaskJSON["completed_at"] = at
			updateSubTaskJSON["last_cgp_verification_end_at"] = at
			updateSubTaskJSON["is_cgp_verification_success"] = true
			if d.Report != nil {
				updateSubTaskJSON["last_cgp_verification_sub_task_report_id"] = subTaskReportId
				updateSubTaskJSON["last_cgp_verification_sub_task_report_uid"] = subTaskReportUid
			}
		case base.SubTaskStatusCGPVerificationFail:
			updateSubTaskJSON["last_cgp_verification_end_at"] = at
			updateSubTaskJSON["is_cgp_verification_success"] = false
			if d.Report != nil {
				updateSubTaskJSON["last_cgp_verification_sub_task_report_id"] = subTaskReportId
				updateSubTaskJSON["last_cgp_verification_sub_task_report_uid"] = subTaskReportUid
			}
		}

		switch d.OperationName {
		case base.UserAsCGPOperationNameSubTaskCGPEditAfterVerifySuccess:
			updateSubTaskJSON["last_form_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_form_sub_task_report_uid"] = subTaskReportUid
		}

		// Update sub task
		_, err2 = tm.SubTask.TxUpdate(dtx, updateSubTaskJSON, utils.JSON{
			"id": subTaskId,
		})
		if err2 != nil {
			return d.Aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "SUB_TASK_UPDATE_FAILED: %w", err2)
		}

		// Get updated sub task
		_, subTask, err2 = tm.SubTask.TxShouldGetById(dtx, subTaskId)
		if err2 != nil {
			return d.Aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "GET_UPDATED_SUB_TASK_UPDATE_FAILED: %w", err2)
		}

		// RequestCreate history item
		h := utils.JSON{
			"sub_task_id":      subTaskId,
			"sub_task_uid":     subTaskUid,
			"from_status":      oldSubTaskStatus,
			"to_status":        d.NewSubTaskStatus,
			"timestamp":        at,
			"user_id":          userId,
			"user_uid":         userUid,
			"user_loginid":     userLoginId,
			"user_fullname":    userFullName,
			"operation_nameid": d.OperationName,
		}

		if d.Report != nil {
			h["sub_task_report_id"] = subTaskReportId
			h["sub_task_report_uid"] = subTaskReportUid
		}

		_, err2 = tm.SubTaskHistoryItem.TxInsert(dtx, h)
		if err2 != nil {
			return d.Aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "SUB_TASK_HISTORY_ITEM_CREATE_FAILED: %w", err2)
		}

		// Execute OnExecute callback if provided
		if d.OnExecute != nil {
			err2 = d.OnExecute(d.Aepr, dtx, userId, subTaskId, subTask, subTaskReportId, d.Report)
			if err2 != nil {
				return d.Aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "ONEXECUTE_CALLBACK_FAILED: %w", err2)
			}
		}

		return nil
	})
	if err != nil {
		return nil, 0, "", err
	}

	if isSendResponse {
		// Build response
		r := utils.JSON{"data": utils.JSON{
			"sub_task": subTask,
		}}

		if isPickSubTask {
			r["result"] = "SUCCESS"
		}
		if d.Report != nil {
			r["sub_task_report_id"] = subTaskReportId
			r["sub_task_report_uid"] = subTaskReportUid
		}

		d.Aepr.WriteResponseAsJSON(http.StatusOK, nil, r)
	}

	return subTask, subTaskReportId, subTaskReportUid, nil
}

func (tm *TaskManagement) TxProcessSubTaskStateRaw(dtx *database.DXDatabaseTx, subTaskId int64, userId int64, at time.Time, d StateEngineStructSubTaskStatus) (subTask utils.JSON, subTaskReportId int64, subTaskReportUid string, err error) {
	//	userId := d.Aepr.LocalData["user_id"].(int64)
	var ok bool
	updateSubTaskJSON := utils.JSON{}

	// Use the passed validation function
	_, subTask, err = tm.SubTask.ShouldGetById(&d.Aepr.Log, subTaskId)
	if err != nil {
		return nil, 0, "", err
	}
	subTaskUid := subTask["uid"].(string)

	userUid := ""
	userLoginId := "SYSTEM"
	userFullName := "SYSTEM"
	userPhoneNumber := ""
	organizationId := int64(0)
	organizationUid := ""
	organizationName := "SYSTEM"

	if userId != 0 {
		_, user, err := user_management.ModuleUserManagement.User.ShouldGetById(&d.Aepr.Log, userId)
		if err != nil {
			return nil, 0, "", err
		}
		if user["is_deleted"].(bool) {
			return nil, 0, "", d.Aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "USER_IS_DELETED")
		}
		if user["status"].(string) != user_management.UserStatusActive {
			return nil, 0, "", d.Aepr.WriteResponseAndNewErrorf(http.StatusConflict, "", "USER_IS_NOT_ACTIVE")
		}
		userUid, ok = user["uid"].(string)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:USER_UID_NOT_FOUND_ON_USER")
		}
	}

	switch d.UserType {
	case base.UserTypeCGP:
		err = validateCGP(d.Aepr, userId, subTask)
		if err != nil {
			return subTask, 0, "", err
		}
	case base.UserTypeFieldSupervisor:
		_, err = validateFieldSupervisor(d.Aepr, userId, subTask)
		if err != nil {
			return subTask, 0, "", err
		}
	case base.UserTypeFieldExecutor:
		fieldExecutor, err := validateFieldExecutor(d.Aepr, userId, subTask)
		if err != nil {
			return subTask, 0, "", err
		}

		userLoginId, ok = fieldExecutor["user_loginid"].(string)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:USER_LOGINID_NOT_FOUND_ON_FIELD_EXECUTOR")
		}
		userFullName, ok = fieldExecutor["user_fullname"].(string)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:USER_FULLNAME_NOT_FOUND_ON_FIELD_EXECUTOR")
		}
		userPhoneNumber, ok = fieldExecutor["user_phonenumber"].(string)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:USER_PHONENUMBER_NOT_FOUND_ON_FIELD_EXECUTOR")
		}
		organizationId, ok = fieldExecutor["organization_id"].(int64)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:ORGANIZATION_ID_NOT_FOUND_ON_FIELD_EXECUTOR")
		}
		organizationUid, ok = fieldExecutor["organization_uid"].(string)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:ORGANIZATION_UID_NOT_FOUND_ON_FIELD_EXECUTOR")
		}
		organizationName, ok = fieldExecutor["organization_name"].(string)
		if !ok {
			return subTask, 0, "", errors.Errorf("IMPOSSIBLE:ORGANIZATION_NAME_NOT_FOUND_ON_FIELD_EXECUTOR")
		}
	case base.UserTypeNone, base.UserTypeAny:
		// No validation required
	}

	// Get current task status
	p := utils.JSON{
		"id": subTaskId,
	}
	if d.SubTaskStatusCondition != nil {
		p["c1"] = db.SQLExpression{Expression: dbUtils.SQLBuildWhereInClause("status", d.SubTaskStatusCondition)}
	}
	isPickSubTask := slices.Contains(d.SubTaskStatusCondition, base.SubTaskStatusWaitingAssignment)

	var oldSubTask utils.JSON
	_, oldSubTask, err = tm.SubTask.TxShouldSelectOneForUpdate(dtx, p, nil)
	if err != nil {
		if isPickSubTask {
			d.Aepr.Log.Warnf("SUB_TASK_PICK_FAILED: %v", err)
			d.Aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
				"result": "FAILED",
			})
			return subTask, 0, "", nil
		}
		d.Aepr.WriteResponseAsError(http.StatusBadRequest, errors.New("INVALID_SUB_TASK_STATUS_FOR_OPERATION"))
		return subTask, 0, "", err
	}

	oldSubTaskStatus := oldSubTask["status"].(string)

	// Handle Report creation
	if d.Report != nil {
		subTaskReportId, err = tm.TxSubTaskReportCreate(dtx, at,
			userId,
			userUid,
			userLoginId,
			userFullName,
			userPhoneNumber,
			organizationId,
			organizationUid,
			organizationName,
			subTaskId,
			subTaskUid,
			d.NewSubTaskStatus, d.Report)
		if err != nil {
			return subTask, 0, "", errors.Errorf("SUB_TASK_REPORT_CREATION_FAILED:%w", err)
		}

		_, subTaskReport, err := tm.SubTaskReport.TxShouldGetById(dtx, subTaskReportId)
		if err != nil {
			return subTask, 0, "", errors.Errorf("SUB_TASK_REPORT_GET_FAILED:%w", err)
		}
		subTaskReportUid = subTaskReport["uid"].(string)

		d.Report["sub_task_report_id"] = subTaskReportId
		d.Report["sub_task_report_uid"] = subTaskReportUid
		updateSubTaskJSON["last_sub_task_report_id"] = subTaskReportId
		updateSubTaskJSON["last_sub_task_report_uid"] = subTaskReportUid
	}

	// Set new status
	updateSubTaskJSON["status"] = d.NewSubTaskStatus

	switch d.UserType {
	case base.UserTypeCGP:
		updateSubTaskJSON["last_cgp_user_id"] = userId
		updateSubTaskJSON["last_cgp_user_uid"] = userUid
	case base.UserTypeFieldExecutor:
		updateSubTaskJSON["last_field_executor_user_id"] = userId
		updateSubTaskJSON["last_field_executor_user_uid"] = userUid
	case base.UserTypeFieldSupervisor:
		updateSubTaskJSON["last_field_supervisor_user_id"] = userId
		updateSubTaskJSON["last_field_supervisor_user_uid"] = userUid
	case base.UserTypeNone, base.UserTypeAny:
	default:
		return subTask, 0, "", d.Aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "USER_TYPE_NOT_SUPPORTED:%s", d.UserType)
	}

	switch oldSubTaskStatus {
	case base.SubTaskStatusWorking:
		updateSubTaskJSON["working_end_at"] = at
		if d.Report != nil {
			updateSubTaskJSON["last_working_end_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_working_end_sub_task_report_uid"] = subTaskReportUid
		}
	case base.SubTaskStatusReworking:
		updateSubTaskJSON["last_reworking_end_at"] = at
		if d.Report != nil {
			updateSubTaskJSON["last_reworking_end_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_reworking_end_sub_task_report_uid"] = subTaskReportUid
		}
	case base.SubTaskStatusFixing:
		updateSubTaskJSON["last_fixing_end_at"] = at
		if d.Report != nil {
			updateSubTaskJSON["last_fixing_end_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_fixing_end_sub_task_report_uid"] = subTaskReportUid
		}
	case base.SubTaskStatusPaused:
		updateSubTaskJSON["last_end_pause_at"] = at
	}

	switch d.NewSubTaskStatus {
	case base.SubTaskStatusWaitingAssignment:
		updateSubTaskJSON["last_field_executor_user_id"] = nil
		updateSubTaskJSON["last_field_executor_user_uid"] = nil
		updateSubTaskJSON["last_field_executor_user_loginid"] = nil
		updateSubTaskJSON["last_field_executor_user_fullname"] = nil
		updateSubTaskJSON["last_field_supervisor_user_id"] = nil
		updateSubTaskJSON["last_field_supervisor_user_uid"] = nil
		updateSubTaskJSON["last_field_supervisor_user_loginid"] = nil
		updateSubTaskJSON["last_field_supervisor_user_fullname"] = nil
	case base.SubTaskStatusWorking:
		updateSubTaskJSON["working_start_at"] = at
	case base.SubTaskStatusWaitingVerification:
		updateSubTaskJSON["is_working_finish"] = true
	case base.SubTaskStatusVerificationSuccess:
		updateSubTaskJSON["last_verification_end_at"] = at
		updateSubTaskJSON["is_verification_success"] = true
		if d.Report != nil {
			updateSubTaskJSON["last_verification_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_verification_sub_task_report_uid"] = subTaskReportUid
		}

	case base.SubTaskStatusVerificationFail:
		updateSubTaskJSON["last_verification_end_at"] = at
		updateSubTaskJSON["is_verification_success"] = false
		if d.Report != nil {
			updateSubTaskJSON["last_verification_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_verification_sub_task_report_uid"] = subTaskReportUid
		}

	case base.SubTaskStatusFixing:
		fixCount := oldSubTask["fix_count"].(int64)
		updateSubTaskJSON["fix_count"] = fixCount + 1
		if fixCount == 1 {
			updateSubTaskJSON["first_fixing_start_at"] = at
		}
	case base.SubTaskStatusPaused:
		updateSubTaskJSON["last_start_pause_at"] = at
		if d.Report != nil {
			updateSubTaskJSON["last_pause_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_pause_sub_task_report_uid"] = subTaskReportUid
		}
	case base.SubTaskStatusCanceledByFieldExecutor:
		updateSubTaskJSON["last_canceled_by_field_executor_at"] = at
		if d.Report != nil {
			updateSubTaskJSON["last_canceled_by_field_executor_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_canceled_by_field_executor_sub_task_report_uid"] = subTaskReportUid
		}
	case base.SubTaskStatusCanceledByCustomer:
		updateSubTaskJSON["last_canceled_by_customer_at"] = at
		updateSubTaskJSON["completed_at"] = at
		if d.Report != nil {
			updateSubTaskJSON["last_canceled_by_customer_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_canceled_by_customer_sub_task_report_uid"] = subTaskReportUid
		}
	case base.SubTaskStatusCGPVerificationSuccess:
		updateSubTaskJSON["completed_at"] = at
		updateSubTaskJSON["last_cgp_verification_end_at"] = at
		updateSubTaskJSON["is_cgp_verification_success"] = true
		if d.Report != nil {
			updateSubTaskJSON["last_cgp_verification_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_cgp_verification_sub_task_report_uid"] = subTaskReportUid
		}
	case base.SubTaskStatusCGPVerificationFail:
		updateSubTaskJSON["last_cgp_verification_end_at"] = at
		updateSubTaskJSON["is_cgp_verification_success"] = false
		if d.Report != nil {
			updateSubTaskJSON["last_cgp_verification_sub_task_report_id"] = subTaskReportId
			updateSubTaskJSON["last_cgp_verification_sub_task_report_uid"] = subTaskReportUid
		}
	}

	// Update sub task
	_, err = tm.SubTask.TxUpdate(dtx, updateSubTaskJSON, utils.JSON{
		"id": subTaskId,
	})
	if err != nil {
		return subTask, 0, "", errors.Errorf("sub task update failed: %w", err)
	}

	// Get updated sub task
	_, subTask, err = tm.SubTask.TxShouldGetById(dtx, subTaskId)
	if err != nil {
		return subTask, 0, "", errors.Errorf("get updated sub task failed: %w", err)
	}

	// RequestCreate history item
	h := utils.JSON{
		"sub_task_id":      subTaskId,
		"sub_task_uid":     subTaskUid,
		"from_status":      oldSubTaskStatus,
		"to_status":        d.NewSubTaskStatus,
		"timestamp":        at,
		"user_id":          userId,
		"user_uid":         userUid,
		"user_loginid":     userLoginId,
		"user_fullname":    userFullName,
		"operation_nameid": d.OperationName,
	}

	if d.Report != nil {
		h["sub_task_report_id"] = subTaskReportId
		h["sub_task_report_uid"] = subTaskReportUid
	}

	_, err = tm.SubTaskHistoryItem.TxInsert(dtx, h)
	if err != nil {
		return subTask, 0, "", errors.Errorf("history item creation failed: %w", err)
	}

	// Execute OnExecute callback if provided
	if d.OnExecute != nil {
		err = d.OnExecute(d.Aepr, dtx, userId, subTaskId, subTask, subTaskReportId, d.Report)
		if err != nil {
			return subTask, 0, "", errors.Errorf("OnExecute callback failed: %w", err)
		}
	}

	return subTask, subTaskReportId, "", nil
}
