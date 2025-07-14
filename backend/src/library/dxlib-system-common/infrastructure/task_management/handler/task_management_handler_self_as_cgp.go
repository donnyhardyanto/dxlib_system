package handler

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"time"
)

func SelfAsCGPUserSubTaskCGPVerifySuccess(aepr *api.DXAPIEndPointRequest) (err error) {
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
	report := map[string]any{}
	userId := aepr.LocalData["user_id"].(int64)

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeCGP,
		SubTaskStatusCondition: []string{base.SubTaskStatusWaitingVerification, base.SubTaskStatusVerificationSuccess},
		NewSubTaskStatus:       base.SubTaskStatusCGPVerificationSuccess,
		OperationName:          base.UserAsCGPOperationNameSubTaskCGPVerifySuccess,
		Report:                 report,
		OnExecute:              task_management.ModuleTaskManagement.OnSubTaskCGPVerifySuccess,
	}, true)

	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func SelfAsCGPUserSubTaskCGPVerifyFail(aepr *api.DXAPIEndPointRequest) (err error) {
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
		UserType:               base.UserTypeCGP,
		SubTaskStatusCondition: []string{base.SubTaskStatusWaitingVerification, base.SubTaskStatusVerificationSuccess},
		NewSubTaskStatus:       base.SubTaskStatusCGPVerificationFail,
		OperationName:          base.UserAsCGPOperationNameSubTaskCGPVerifyFail,
		Report:                 report,
		OnExecute:              task_management.ModuleTaskManagement.OnSubTaskCGPVerifyFail,
	}, true)

	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func SelfAsCGPUserSubTaskCGPEditAfterCGPVerifySuccess(aepr *api.DXAPIEndPointRequest) (err error) {
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
		UserType:               base.UserTypeCGP,
		SubTaskStatusCondition: []string{base.SubTaskStatusCGPVerificationSuccess},
		NewSubTaskStatus:       base.SubTaskStatusCGPVerificationSuccess,
		OperationName:          base.UserAsCGPOperationNameSubTaskCGPEditAfterVerifySuccess,
		Report:                 effectiveReport,
		OnExecute:              task_management.ModuleTaskManagement.OnSubTaskCGPVerifySuccess,
	}, true)

	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}
