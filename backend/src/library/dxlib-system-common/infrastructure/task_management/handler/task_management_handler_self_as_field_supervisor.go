package handler

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"net/http"
	"slices"
	"time"
)

func ValidateFieldSupervisorWithSubTask(aepr *api.DXAPIEndPointRequest, userId int64, subTask utils.JSON) (err error) {
	_, fieldSupervisors, err := partner_management.ModulePartnerManagement.FieldSupervisor.Select(&aepr.Log, nil, utils.JSON{
		"user_id": userId,
	}, nil, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if len(fieldSupervisors) == 0 {
		return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:USER_IS_NOT_FIELD_SUPERVISOR")
	}

	subTaskTypeId := subTask["sub_task_type_id"].(int64)
	subTaskCustomerAreaCode := subTask["customer_sales_area_code"].(string)
	subTaskCustomerProvinceLocationCode := subTask["customer_address_province_location_code"].(string)
	subTaskCustomerKabupatenLocationCode := subTask["customer_address_kabupaten_location_code"].(string)
	subTaskCustomerKecamatanLocationCode := subTask["customer_address_kecamatan_location_code"].(string)
	subTaskCustomerKelurahanLocationCode := subTask["customer_address_kelurahan_location_code"].(string)

	for _, fieldSupervisor := range fieldSupervisors {
		userRoleMembershipId := fieldSupervisor["id"].(int64)

		effectiveExpertises, err := FieldSupervisorGetEffectiveExpertise(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return err
		}
		if len(effectiveExpertises) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_EXPERTISE")
		}

		if !slices.Contains(effectiveExpertises, subTaskTypeId) {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:SUB_TASK_TYPE_NOT_WITHIN_YOUR_EXPERTISE")
		}

		effectiveAreaCodes, err := FieldSupervisorGetEffectiveAreas(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return err
		}
		if len(effectiveAreaCodes) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_AREA")
		}

		if !slices.Contains(effectiveAreaCodes, subTaskCustomerAreaCode) {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:SUB_TASK_CUSTOMER_SALES_AREA_NOT_WITHIN_YOUR_AREA")
		}

		effectiveLocations, err := FieldSupervisorGetEffectiveLocations(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveLocations) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_LOCATION")
		}

		found := false
		if slices.Contains(effectiveLocations, subTaskCustomerProvinceLocationCode) {
			found = true
		}
		if slices.Contains(effectiveLocations, subTaskCustomerKabupatenLocationCode) {
			found = true
		}
		if slices.Contains(effectiveLocations, subTaskCustomerKecamatanLocationCode) {
			found = true
		}
		if slices.Contains(effectiveLocations, subTaskCustomerKelurahanLocationCode) {
			found = true
		}
		if !found {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:SUB_TASK_CUSTOMER_LOCATION_NOT_WITHIN_YOUR_LOCATION")
		}
	}

	return nil
}

func SelfAsFieldSupervisorUserSubTaskVerifySuccess(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return err
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return err
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	report := map[string]any{}
	userId := aepr.LocalData["user_id"].(int64)

	err = ValidateFieldSupervisorWithSubTask(aepr, userId, subTask)
	if err != nil {
		return err
	}

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldSupervisor,
		SubTaskStatusCondition: []string{base.SubTaskStatusWaitingVerification},
		NewSubTaskStatus:       base.SubTaskStatusVerificationSuccess,
		OperationName:          base.UserAsFieldSupervisorOperationNameSubTaskVerifySuccess,
		Report:                 report,
		OnExecute:              task_management.ModuleTaskManagement.OnSubTaskVerifySuccess,
	}, true)

	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func SelfAsFieldSupervisorUserSubTaskVerifyFail(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return err
	}
	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return err
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}

	at := aepr.ParameterValues["at"].Value.(time.Time)
	report := aepr.ParameterValues["report"].Value.(map[string]any)
	userId := aepr.LocalData["user_id"].(int64)

	err = ValidateFieldSupervisorWithSubTask(aepr, userId, subTask)
	if err != nil {
		return err
	}

	_, _, _, err = task_management.ModuleTaskManagement.ProcessSubTaskStateRaw(subTaskId, userId, at, task_management.StateEngineStructSubTaskStatus{
		Aepr:                   aepr,
		UserType:               base.UserTypeFieldSupervisor,
		SubTaskStatusCondition: []string{base.SubTaskStatusWaitingVerification},
		NewSubTaskStatus:       base.SubTaskStatusVerificationFail,
		OperationName:          base.UserAsFieldSupervisorOperationNameSubTaskVerifyFail,
		Report:                 report,
		OnExecute:              task_management.ModuleTaskManagement.OnSubTaskVerifyFail,
	}, true)

	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}
