package handler

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
)

func OrganizationExecutorPagingList(aepr *api.DXAPIEndPointRequest) (err error) {

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

	return partner_management.ModulePartnerManagement.OrganizationExecutor.DoRequestPagingList(aepr, filterWhere, filterOrderBy, filterKeyValues, func(listRow utils.JSON) (r utils.JSON, err error) {
		_, organizationExecutorArea, err := partner_management.ModulePartnerManagement.OrganizationExecutorArea.Select(&aepr.Log, nil, utils.JSON{
			"organization_role_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["organization_executor_area"] = organizationExecutorArea

		_, organizationExecutorLocation, err := partner_management.ModulePartnerManagement.OrganizationExecutorLocation.Select(&aepr.Log, nil, utils.JSON{
			"organization_role_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["organization_executor_location"] = organizationExecutorLocation

		_, organizationExecutorExpertise, err := partner_management.ModulePartnerManagement.OrganizationExecutorExpertise.Select(&aepr.Log, nil, utils.JSON{
			"organization_role_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["organization_executor_expertise"] = organizationExecutorExpertise

		_, fieldExecutors, err := partner_management.ModulePartnerManagement.FieldExecutor.Select(&aepr.Log, nil, utils.JSON{
			"organization_id": listRow["organization_id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_executors"] = fieldExecutors

		return listRow, nil
	})
}
