package partner_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database/protected/db"
	"github.com/donnyhardyanto/dxlib/log"
	dxlibModule "github.com/donnyhardyanto/dxlib/module"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils"
	"gopkg.in/errgo.v2/fmt/errors"
	"net/http"
	"strings"
)

type PartnerManagement struct {
	dxlibModule.DXModule

	RoleArea           *table.DXRawTable
	RoleTaskType       *table.DXRawTable
	Role               *table.DXTable
	UserRoleMembership *table.DXTable

	OrganizationSupervisor          *table.DXRawTable
	OrganizationSupervisorArea      *table.DXRawTable
	OrganizationSupervisorLocation  *table.DXRawTable
	OrganizationSupervisorExpertise *table.DXRawTable

	FieldSupervisor          *table.DXRawTable
	FieldSupervisorLocation  *table.DXRawTable
	FieldSupervisorArea      *table.DXRawTable
	FieldSupervisorExpertise *table.DXRawTable

	FieldSupervisorSubTaskVerificationStats *table.DXRawTable
	FieldSupervisorEffectiveLocation        *table.DXRawTable
	FieldSupervisorEffectiveArea            *table.DXRawTable
	FieldSupervisorEffectiveExpertise       *table.DXRawTable

	OrganizationExecutor          *table.DXRawTable
	OrganizationExecutorArea      *table.DXRawTable
	OrganizationExecutorLocation  *table.DXRawTable
	OrganizationExecutorExpertise *table.DXRawTable

	FieldExecutor *table.DXRawTable

	FieldExecutorLocation  *table.DXRawTable
	FieldExecutorArea      *table.DXRawTable
	FieldExecutorExpertise *table.DXRawTable

	FieldExecutorSubTaskStatusStats *table.DXRawTable
	FieldExecutorEffectiveLocation  *table.DXRawTable
	FieldExecutorEffectiveArea      *table.DXRawTable
	FieldExecutorEffectiveExpertise *table.DXRawTable
}

var ModulePartnerManagement = PartnerManagement{}

func (pm *PartnerManagement) Init(databaseNameId string) {
	pm.DatabaseNameId = databaseNameId

	pm.RoleArea = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.role_area", "partner_management.role_area",
		"partner_management.v_role_area", "id", "id", "uid", "data")

	pm.RoleTaskType = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.role_task_type", "partner_management.role_task_type",
		"partner_management.v_role_task_type", "id", "id", "uid", "data")

	pm.Role = table.Manager.NewTable(pm.DatabaseNameId,
		"user_management.role", "user_management.role",
		"partner_management.v_role", "nameid", "id", "uid", "data")

	pm.Role.FieldTypeMapping = map[string]string{
		"organization_types": "array-string",
	}
	pm.UserRoleMembership = table.Manager.NewTable(pm.DatabaseNameId,
		"user_management.user_role_membership", "user_management.user_role_membership",
		"partner_management.v_user_role_membership", "id", "id", "uid", "data")

	pm.OrganizationSupervisor = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.v_organization_supervisor", "partner_management.organization_supervisor",
		"partner_management.v_organization_supervisor", "id", "id", "uid", "data")

	pm.OrganizationSupervisorArea = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.organization_supervisor_area", "partner_management.organization_supervisor_area",
		"partner_management.v_organization_supervisor_area", "id", "id", "uid", "data")

	pm.OrganizationSupervisorLocation = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.organization_supervisor_location", "partner_management.organization_supervisor_location",
		"partner_management.v_organization_supervisor_location", "id", "id", "uid", "data")

	pm.OrganizationSupervisorExpertise = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.organization_supervisor_expertise", "partner_management.organization_supervisor_expertise",
		"partner_management.v_organization_supervisor_expertise", "id", "id", "uid", "data")

	pm.FieldSupervisor = table.Manager.NewRawTable(pm.DatabaseNameId,
		"user_management.user_role_membership", "partner_management.field_supervisor",
		"partner_management.v_field_supervisor", "id", "id", "uid", "data")

	pm.FieldSupervisorLocation = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.field_supervisor_location", "partner_management.field_supervisor_location",
		"partner_management.v_field_supervisor_location", "id", "id", "uid", "data")
	pm.FieldSupervisorArea = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.field_supervisor_area", "partner_management.field_supervisor_area",
		"partner_management.v_field_supervisor_area", "id", "id", "uid", "data")
	pm.FieldSupervisorExpertise = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.field_supervisor_expertise", "partner_management.field_supervisor_expertise",
		"partner_management.v_field_supervisor_expertise", "id", "id", "uid", "data")

	pm.FieldSupervisorSubTaskVerificationStats = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.mv_field_supervisor_verification_stats", "partner_management.field_supervisor_verification_stats",
		"partner_management.mv_field_supervisor_verification_stats", "id", "id", "uid", "data")

	pm.FieldSupervisorEffectiveLocation = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.mv_field_supervisor_effective_location", "partner_management.field_supervisor_effective_location",
		"partner_management.mv_field_supervisor_effective_location", "", "", "", "data")

	pm.FieldSupervisorEffectiveArea = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.mv_field_supervisor_effective_area", "partner_management.field_supervisor_effective_area",
		"partner_management.mv_field_supervisor_effective_area", "", "", "", "data")

	pm.FieldSupervisorEffectiveExpertise = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.mv_field_supervisor_effective_expertise", "partner_management.field_supervisor_effective_expertise",
		"partner_management.mv_field_supervisor_effective_expertise", "", "", "", "data")

	pm.OrganizationExecutor = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.v_organization_executor", "partner_management.organization_executor",
		"partner_management.v_organization_executor", "id", "id", "uid", "data")

	pm.OrganizationExecutorArea = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.organization_executor_area", "partner_management.organization_executor_area",
		"partner_management.v_organization_executor_area", "id", "id", "uid", "data")

	pm.OrganizationExecutorLocation = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.organization_executor_location", "partner_management.organization_executor_location",
		"partner_management.v_organization_executor_location", "id", "id", "uid", "data")

	pm.OrganizationExecutorExpertise = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.organization_executor_expertise", "partner_management.organization_executor_expertise",
		"partner_management.v_organization_executor_expertise", "id", "id", "uid", "data")

	pm.FieldExecutor = table.Manager.NewRawTable(pm.DatabaseNameId,
		"user_management.user_role_membership", "partner_management.field_executor",
		"partner_management.v_field_executor", "id", "id", "uid", "data")

	pm.FieldExecutorLocation = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.field_executor_location", "partner_management.field_executor_location",
		"partner_management.v_field_executor_location", "id", "id", "uid", "data")

	pm.FieldExecutorArea = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.field_executor_area", "partner_management.field_executor_area",
		"partner_management.v_field_executor_area", "id", "id", "uid", "data")

	pm.FieldExecutorExpertise = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.field_executor_expertise", "partner_management.field_executor_expertise",
		"partner_management.v_field_executor_expertise", "id", "id", "uid", "data")

	pm.FieldExecutorSubTaskStatusStats = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.field_executor_sub_task_status_summary", "task_management.field_executor_sub_task_status_summary",
		"partner_management.mv_field_executor_sub_task_status_summary", "id", "id", "uid", "data")

	pm.FieldExecutorEffectiveLocation = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.mv_field_executor_effective_location", "partner_management.field_executor_effective_location",
		"partner_management.mv_field_executor_effective_location", "", "", "", "data")

	pm.FieldExecutorEffectiveArea = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.mv_field_executor_effective_area", "partner_management.field_executor_effective_area",
		"partner_management.mv_field_executor_effective_area", "", "", "", "data")

	pm.FieldExecutorEffectiveExpertise = table.Manager.NewRawTable(pm.DatabaseNameId,
		"partner_management.mv_field_executor_effective_expertise", "partner_management.field_executor_effective_expertise",
		"partner_management.mv_field_executor_effective_expertise", "", "", "", "data")

}

func (pm *PartnerManagement) FieldSupervisorCheckIsUser(log *log.DXLog, userId int64, organizationId int64) (isTrue bool, err error) {
	_, fieldSupervisor, err := pm.FieldSupervisor.SelectOne(log, nil, utils.JSON{
		"user_id":         userId,
		"organization_id": organizationId,
	}, nil, nil)
	if err != nil {
		return false, err
	}
	if fieldSupervisor == nil {
		return false, nil
	}
	return true, nil
}

func (pm *PartnerManagement) FieldSupervisorDoRequestIsUser(aepr *api.DXAPIEndPointRequest, userId int64, organizationId int64) (isTrue bool) {
	isTrue, err := pm.FieldSupervisorCheckIsUser(&aepr.Log, userId, organizationId)
	if !isTrue {
		err = errors.Newf("USER_NOT_FIELD_SUPERVISOR\n{userId=%d, organizationId=%d}", userId, organizationId)
		aepr.WriteResponseAndLogAsError(http.StatusForbidden, "", err)
		return isTrue
	}
	if err != nil {
		aepr.WriteResponseAndLogAsError(http.StatusBadRequest, "", err)
		return isTrue
	}
	return isTrue
}

func (pm *PartnerManagement) FieldSupervisorCheckHasEffectiveExpertise(log *log.DXLog, userId int64, organizationId int64, subTaskTypeIds []int64) (isTrue bool, err error) {
	subTaskTypeIdAsStrings := utils.ArrayInt64ToArrayString(subTaskTypeIds)
	conditions := "sub_task_type_id in (" + strings.Join(subTaskTypeIdAsStrings, " , ") + ")"

	_, effectiveExpertise, err := pm.FieldSupervisorEffectiveExpertise.SelectOne(log, nil, utils.JSON{
		"user_id":         userId,
		"organization_id": organizationId,
		"c1":              db.SQLExpression{Expression: conditions},
	}, nil, nil)
	if err != nil {
		return false, err
	}
	if effectiveExpertise == nil {
		return false, nil
	}
	return true, nil
}

func (pm *PartnerManagement) FieldSupervisorDoRequestHasEffectiveExpertise(aepr *api.DXAPIEndPointRequest, userId int64, organizationId int64, subTaskTypeIds []int64) (isTrue bool) {
	isTrue, err := pm.FieldSupervisorCheckHasEffectiveExpertise(&aepr.Log, userId, organizationId, subTaskTypeIds)
	if !isTrue {
		err = errors.Newf("USER_AS_FIELD_SUPERVISOR_DONT_HAVE_SPESIFIC_EFFECTIVE_EXPERTISE\n{userId=%d,organizationId=%d,subTaskTypeIds=%v}", userId, organizationId, subTaskTypeIds)
		aepr.WriteResponseAndLogAsError(http.StatusForbidden, "", err)
		return isTrue
	}
	if err != nil {
		aepr.WriteResponseAndLogAsError(http.StatusBadRequest, "", err)
		return isTrue
	}
	return isTrue
}

func (pm *PartnerManagement) FieldSupervisorCheckHasEffectiveLocation(log *log.DXLog, userId int64, organizationId int64, locationCodes []string) (isTrue bool, err error) {
	for k, v := range locationCodes {
		locationCodes[k] = "'" + v + "'"
	}
	conditions := "location_code in (" + strings.Join(locationCodes, " , ") + ")"

	_, effectiveExpertise, err := pm.FieldSupervisorEffectiveLocation.SelectOne(log, nil, utils.JSON{
		"user_id":         userId,
		"organization_id": organizationId,
		"c1":              db.SQLExpression{Expression: conditions},
	}, nil, nil)
	if err != nil {
		return false, err
	}
	if effectiveExpertise == nil {
		return false, nil
	}
	return true, nil
}

func (pm *PartnerManagement) FieldSupervisorDoRequestHasEffectiveLocation(aepr *api.DXAPIEndPointRequest, userId int64, organizationId int64, locationCodes []string) (isTrue bool) {
	isTrue, err := pm.FieldSupervisorCheckHasEffectiveLocation(&aepr.Log, userId, organizationId, locationCodes)
	if !isTrue {
		err = errors.Newf("USER_AS_FIELD_SUPERVISOR_DONT_HAVE_SPESIFIC_EFFECTIVE_LOCATION\n{userId=%d,organizationId=%d,locationCodes=%v}", userId, organizationId, locationCodes)
		aepr.WriteResponseAndLogAsError(http.StatusForbidden, "", err)
		return isTrue
	}
	if err != nil {
		aepr.WriteResponseAndLogAsError(http.StatusBadRequest, "", err)
		return isTrue
	}
	return isTrue
}

func (pm *PartnerManagement) FieldSupervisorCheckHasEffectiveArea(log *log.DXLog, userId int64, organizationId int64, areaCodes []string) (isTrue bool, err error) {
	for k, v := range areaCodes {
		areaCodes[k] = "'" + v + "'"
	}
	conditions := "area_code in (" + strings.Join(areaCodes, " , ") + ")"

	_, effectiveExpertise, err := pm.FieldSupervisorEffectiveArea.SelectOne(log, nil, utils.JSON{
		"user_id":         userId,
		"organization_id": organizationId,
		"c1":              db.SQLExpression{Expression: conditions},
	}, nil, nil)
	if err != nil {
		return false, err
	}
	if effectiveExpertise == nil {
		return false, nil
	}
	return true, nil
}

func (pm *PartnerManagement) FieldSupervisorDoRequestHasEffectiveArea(aepr *api.DXAPIEndPointRequest, userId int64, organizationId int64, areaCodes []string) (isTrue bool) {
	isTrue, err := pm.FieldSupervisorCheckHasEffectiveArea(&aepr.Log, userId, organizationId, areaCodes)
	if !isTrue {
		err = errors.Newf("USER_AS_FIELD_SUPERVISOR_DONT_HAVE_SPESIFIC_EFFECTIVE_AREA\n{userId=%d,organizationId=%d,areaCodes=%v}", userId, organizationId, areaCodes)
		aepr.WriteResponseAndLogAsError(http.StatusForbidden, "", err)
		return isTrue
	}
	if err != nil {
		aepr.WriteResponseAndLogAsError(http.StatusBadRequest, "", err)
		return isTrue
	}
	return isTrue
}

func (pm *PartnerManagement) FieldSupervisorDoRequestIsUserAndHasEffectiveExpertiseAreaLocation(aepr *api.DXAPIEndPointRequest, userId int64, organizationId int64, subTaskTypeIds []int64,
	locationCodes []string, areaCodes []string) (isTrue bool) {

	isFieldSupervisor := pm.FieldSupervisorDoRequestIsUser(aepr, userId, organizationId)
	if !isFieldSupervisor {
		return false
	}

	isHasExpertise := pm.FieldSupervisorDoRequestHasEffectiveExpertise(aepr, userId, organizationId, subTaskTypeIds)
	if !isHasExpertise {
		return false
	}

	isHasLocation := pm.FieldSupervisorDoRequestHasEffectiveLocation(aepr, userId, organizationId, locationCodes)
	if !isHasLocation {
		return false
	}

	isHasArea := pm.FieldSupervisorDoRequestHasEffectiveArea(aepr, userId, organizationId, areaCodes)
	if !isHasArea {
		return false
	}

	return true
}
