package handler

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/arrears_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	_ "github.com/tealeg/xlsx"
	"net/http"
	"time"
)

func TaskPiutangCreate(aepr *api.DXAPIEndPointRequest) (err error) {

	taskArrears := map[string]interface{}{}
	for k, v := range aepr.ParameterValues {
		taskArrears[k] = v.Value
	}
	err, subTaskId := arrears_management.ModuleArrearsManagement.DoTaskArrearsCreate(&aepr.Log, taskArrears)
	if err != nil {
		aepr.WriteResponseAsError(http.StatusBadRequest, err)
		return nil
	}
	_, data, err := task_management.ModuleTaskManagement.SubTask.ShouldGetById(&aepr.Log, subTaskId)
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": data})
	return nil

}

// TaskPiutangReadByTaskUid retrieves customer piutang records by customer ID
func TaskPiutangReadByTaskUid(aepr *api.DXAPIEndPointRequest) (err error) {
	_, taskUid, err := aepr.GetParameterValueAsString("task_uid")
	if err != nil {
		return errors.Wrap(err, "error occurred")
	}

	// Verify customer exists
	_, r, err := task_management.ModuleTaskManagement.Task.ShouldGetByUid(&aepr.Log, taskUid)
	if err != nil {
		return errors.Wrap(err, "error occurred")
	}

	_, list, err := arrears_management.ModuleArrearsManagement.DoTaskArrearsGetByTaskId(&aepr.Log, r["id"].(int64))
	if err != nil {
		return errors.Wrap(err, "error occurred")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": list})
	return nil
}

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
	//_, err = ValidateFieldExecutorUserWithSubTask(aepr, userId, subTask)
	//if err != nil {
	//	return err
	//}
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
