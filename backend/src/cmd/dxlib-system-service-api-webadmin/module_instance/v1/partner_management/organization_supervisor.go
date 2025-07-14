package partner_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	partner_management_handler "github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management/handler"
)

func defineAPIOrganizationSupervisor(anAPI *api.DXAPI) {
	anAPI.NewEndPoint("OrganizationSupervisor.List.CMS",
		"Retrieves a paginated list of Organization Supervisor  with filtering and sorting capabilities. "+
			"Returns a structured list of Organization Supervisor  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/organization_supervisor/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management_handler.OrganizationSupervisorPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisor.Read.CMS",
		"Retrieves detailed information for a specific Organization Supervisor by ID. "+
			"Returns comprehensive Organization Supervisor data including specifications. "+
			"Essential for Announcement specification views and data verification.",
		"/v1/organization_supervisor/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisor.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisorLocations.List.CMS",
		"Retrieves a paginated list of Organization Supervisor Locations  with filtering and sorting capabilities. "+
			"Returns a structured list of Organization Supervisor Locations  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/organization_supervisor_location/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisorLocation.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR_LOCATION.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisorLocation.Create.CMS",
		"Creates a new Organization Supervisor Location in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Organization Supervisor Location record with assigned unique identifier.",
		"/v1/organization_supervisor_location/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "organization_role_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "location_code", Type: "protected-string", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisorLocation.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR_LOCATION.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisorLocation.Read.CMS",
		"Retrieves detailed information for a specific Organization Supervisor Location by ID. "+
			"Returns comprehensive Organization Supervisor Location data including specifications. "+
			"Essential for Organization Supervisor Location specification views and data verification.",
		"/v1/organization_supervisor_location/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisorLocation.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR_LOCATION.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisorLocation.Delete.CMS",
		"Permanently removes a Organization Supervisor Location record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/organization_supervisor_location/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisorLocation.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR_LOCATION.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisorAreas.List.CMS",
		"Retrieves a paginated list of Organization Supervisor Areas  with filtering and sorting capabilities. "+
			"Returns a structured list of Organization Supervisor Areas  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/organization_supervisor_area/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisorArea.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR_AREA.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisorArea.Create.CMS",
		"Creates a new Organization Supervisor Area in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Organization Supervisor Area record with assigned unique identifier.",
		"/v1/organization_supervisor_area/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "organization_role_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "area_code", Type: "string", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisorArea.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR_AREA.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisorArea.Read.CMS",
		"Retrieves detailed information for a specific Organization Supervisor Area by ID. "+
			"Returns comprehensive Organization Supervisor Area data including specifications. "+
			"Essential for Announcement specification views and data verification.",
		"/v1/organization_supervisor_area/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisorArea.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR_AREA.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisorArea.Delete.CMS",
		"Permanently removes a Organization Supervisor Area record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/organization_supervisor_area/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisorArea.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR_AREA.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisorExpertise.List.CMS",
		"Retrieves a paginated list of Organization Supervisor Expertise  with filtering and sorting capabilities. "+
			"Returns a structured list of Organization Supervisor Expertise  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/organization_supervisor_expertise/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisorExpertise.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR_EXPERTISE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisorExpertise.Create.CMS",
		"Creates a new Organization Supervisor Expertise in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Organization Supervisor Expertise record with assigned unique identifier.",
		"/v1/organization_supervisor_expertise/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "organization_role_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "sub_task_type_id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisorExpertise.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR_EXPERTISE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisorExpertise.Read.CMS",
		"Retrieves detailed information for a specific Organization Supervisor Expertise by ID. "+
			"Returns comprehensive Organization Supervisor Expertise data including specifications. "+
			"Essential for Organization Supervisor Expertise specification views and data verification.",
		"/v1/organization_supervisor_expertise/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisorExpertise.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR_EXPERTISE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationSupervisorExpertise.Delete.CMS",
		"Permanently removes a Organization Supervisor Expertise record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/organization_supervisor_expertise/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationSupervisorExpertise.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_SUPERVISOR_EXPERTISE.DELETE"}, 0, "default",
	)
}
