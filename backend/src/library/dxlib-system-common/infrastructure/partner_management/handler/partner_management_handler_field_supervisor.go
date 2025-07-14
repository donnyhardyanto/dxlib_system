package handler

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"net/http"
)

/*
	func SelfFieldSupervisorLocationListForFilter(aepr *api.DXAPIEndPointRequest) (err error) {
		userId := aepr.LocalData["user_id"].(int64)
		dbTaskDispatcher := database.Manager.Databases["task-dispatcher"]

		query := "
	        WITH RECURSIVE parent_location AS (
	            SELECT
	                l.name,
	                l.code,
	                l.type,
	                l.parent_value,
	                fe.user_id
	            FROM master_data.location l
	            RIGHT JOIN partner_management.field_supervisor_location fel ON fel.location_code = l.code
	            JOIN partner_management.field_supervisor fe ON fe.id = fel.user_role_membership_id
	            WHERE fe.user_id = $1
	            UNION ALL
	            SELECT
	                l.name,
	                l.code,
	                l.type,
	                l.parent_value,
	                null
	            FROM master_data.location l
	            INNER JOIN parent_location pl ON l.code = pl.parent_value
	        )
	        SELECT
	            MIN(name) AS name,
	            MIN(type) AS type,
	            code,
	            MIN(parent_value) AS parent_value,
	            MIN(user_id::text) AS user_id
	        FROM parent_location
	        GROUP BY code"

		var data []struct {
			Code        string  "db:"code" json:"code""
			Type        string  "db:"type" json:"type""
			Name        string  "db:"name" json:"name""
			ParentValue string  "db:"parent_value" json:"parent_value""
			UserId      *string "db:"user_id" json:"user_id""
		}
		err = dbTaskDispatcher.Connection.Select(&data, query, userId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
			"list": data,
		})
		return nil
	}
*/
func FieldSupervisorPagingList(aepr *api.DXAPIEndPointRequest) (err error) {

	isExistFilterWhere, filterWhere, err := aepr.GetParameterValueAsString("filter_where")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterWhere {
		filterWhere = ""
	}
	isExistFilterOrderBy, filterOrderBy, err := aepr.GetParameterValueAsString("filter_order_by")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterOrderBy {
		filterOrderBy = ""
	}

	isExistFilterKeyValues, filterKeyValues, err := aepr.GetParameterValueAsJSON("filter_key_values")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterKeyValues {
		filterKeyValues = nil
	}

	return partner_management.ModulePartnerManagement.FieldSupervisor.DoRequestPagingList(aepr, filterWhere, filterOrderBy, filterKeyValues, func(listRow utils.JSON) (r utils.JSON, err error) {
		_, fieldSupervisorArea, err := partner_management.ModulePartnerManagement.FieldSupervisorArea.Select(&aepr.Log, nil, utils.JSON{
			"user_role_membership_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_supervisor_area"] = fieldSupervisorArea

		_, fieldSupervisorLocation, err := partner_management.ModulePartnerManagement.FieldSupervisorLocation.Select(&aepr.Log, nil, utils.JSON{
			"user_role_membership_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_supervisor_location"] = fieldSupervisorLocation

		_, fieldSupervisorExpertise, err := partner_management.ModulePartnerManagement.FieldSupervisorExpertise.Select(&aepr.Log, nil, utils.JSON{
			"user_role_membership_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_supervisor_expertise"] = fieldSupervisorExpertise

		_, fieldSupervisorEffectiveArea, err := partner_management.ModulePartnerManagement.FieldSupervisorEffectiveArea.Select(&aepr.Log, nil, utils.JSON{
			"user_role_membership_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_supervisor_effective_area"] = fieldSupervisorEffectiveArea

		_, fieldSupervisorEffectiveLocation, err := partner_management.ModulePartnerManagement.FieldSupervisorEffectiveLocation.Select(&aepr.Log, nil, utils.JSON{
			"user_role_membership_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_supervisor_effective_location"] = fieldSupervisorEffectiveLocation

		_, fieldSupervisorEffectiveExpertise, err := partner_management.ModulePartnerManagement.FieldSupervisorEffectiveExpertise.Select(&aepr.Log, nil, utils.JSON{
			"user_role_membership_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_supervisor_effective_expertise"] = fieldSupervisorEffectiveExpertise

		return listRow, nil
	})
}

func SelfFieldSupervisorSubTaskVerificationStats(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.LocalData["user_id"].(int64)
	_, data, err := partner_management.ModulePartnerManagement.FieldSupervisorSubTaskVerificationStats.SelectOne(&aepr.Log, nil, utils.JSON{
		"user_id": userId,
	}, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
		"data": data,
	})
	return nil
}
