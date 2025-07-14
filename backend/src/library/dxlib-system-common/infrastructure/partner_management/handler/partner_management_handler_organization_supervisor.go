package handler

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
)

func OrganizationSupervisorPagingList(aepr *api.DXAPIEndPointRequest) (err error) {

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

	return partner_management.ModulePartnerManagement.OrganizationSupervisor.DoRequestPagingList(aepr, filterWhere, filterOrderBy, filterKeyValues, func(listRow utils.JSON) (r utils.JSON, err error) {
		_, organizationSupervisorArea, err := partner_management.ModulePartnerManagement.OrganizationSupervisorArea.Select(&aepr.Log, nil, utils.JSON{
			"organization_role_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["organization_supervisor_area"] = organizationSupervisorArea

		_, organizationSupervisorLocation, err := partner_management.ModulePartnerManagement.OrganizationSupervisorLocation.Select(&aepr.Log, nil, utils.JSON{
			"organization_role_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["organization_supervisor_location"] = organizationSupervisorLocation

		_, organizationSupervisorExpertise, err := partner_management.ModulePartnerManagement.OrganizationSupervisorExpertise.Select(&aepr.Log, nil, utils.JSON{
			"organization_role_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["organization_supervisor_expertise"] = organizationSupervisorExpertise

		_, fieldSupervisors, err := partner_management.ModulePartnerManagement.FieldSupervisor.Select(&aepr.Log, nil, utils.JSON{
			"organization_id": listRow["organization_id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_supervisors"] = fieldSupervisors

		return listRow, nil
	})
}
