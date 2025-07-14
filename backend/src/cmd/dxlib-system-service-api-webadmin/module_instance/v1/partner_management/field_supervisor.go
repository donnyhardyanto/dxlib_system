package partner_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
)

func defineAPIFieldSupervisor(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("FieldSupervisor.List.Download.CMS",
		"Retrieves list download of Field Supervisor  with filtering and sorting capabilities. "+
			"Returns a structured list download of Field Supervisor  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_supervisor/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisor.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisor.List.CMS",
		"Retrieves a paginated list of Field Supervisor  with filtering and sorting capabilities. "+
			"Returns a structured list of Field Supervisor  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_supervisor/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, partner_management.ModulePartnerManagement.FieldSupervisor.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisor.Read.CMS",
		"Creates a new Field Supervisor in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Field Supervisor record with assigned unique identifier.",
		"/v1/field_supervisor/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisor.RequestRead, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisor.EffectiveLocations.List",
		"Retrieves a paginated list of Field Supervisor Effective Locations  with filtering and sorting capabilities. "+
			"Returns a structured list of Field Supervisor Effective Locations  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_supervisor_effective_location/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorEffectiveLocation.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_EFFECTIVE_LOCATION.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisorLocations.List.CMS",
		"Retrieves a paginated list of Field Supervisor Locations with filtering and sorting capabilities. "+
			"Returns a structured list of Field Supervisor Locations  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_supervisor_location/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorLocation.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_LOCATION.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisorLocation.Create.CMS",
		"Creates a new Field Supervisor Location in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Field Supervisor Location record with assigned unique identifier.",
		"/v1/field_supervisor_location/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "user_role_membership_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "location_code", Type: "protected-string", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorLocation.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_LOCATION.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisorLocation.Read.CMS",
		"Retrieves detailed information for a specific Field Supervisor Location by ID. "+
			"Returns comprehensive Field Supervisor Location data including specifications. "+
			"Essential for Field Supervisor Location specification views and data verification.",
		"/v1/field_supervisor_location/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorLocation.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_LOCATION.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisorLocation.Delete.CMS",
		"Permanently removes a Field Supervisor Location record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/field_supervisor_location/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorLocation.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_LOCATION.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisorEffectiveAreas.List.CMS",
		"Retrieves a paginated list of Field Supervisor Effective Areas  with filtering and sorting capabilities. "+
			"Returns a structured list of Field Supervisor Effective Areas  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_supervisor_effective_area/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorEffectiveArea.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_EFFECTIVE_AREA.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisorAreas.List.CMS",
		"Retrieves a paginated list of Field Supervisor Areas with filtering and sorting capabilities. "+
			"Returns a structured list of Field Supervisor Areas with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_supervisor_area/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorArea.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_AREA.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisorArea.Create",
		"Creates a new Field Supervisor Area in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Field Supervisor Area record with assigned unique identifier.",
		"/v1/field_supervisor_area/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "user_role_membership_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "area_code", Type: "string", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorArea.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_AREA.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisorArea.Read.CMS",
		"Retrieves detailed information for a specific Field Supervisor Area by ID. "+
			"Returns comprehensive Field Supervisor Area data including specifications. "+
			"Essential for Field Supervisor Area specification views and data verification.",
		"/v1/field_supervisor_area/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorArea.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_AREA.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisorArea.Delete.CMS",
		"Permanently removes a Field Supervisor Area record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/field_supervisor_area/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorArea.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_AREA.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisor.EffectiveExpertise.List.CMS",
		"Retrieves a paginated list of Field Supervisor Effective Expertise  with filtering and sorting capabilities. "+
			"Returns a structured list of Field Supervisor Effective Expertise  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_supervisor_effective_expertise/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorEffectiveExpertise.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_EFFECTIVE_EXPERTISE.LIST"}, 0, "default", // Fixed permission string
	)

	anAPI.NewEndPoint("FieldSupervisorExpertise.List.CMS",
		"Retrieves a paginated list of Field Supervisor Expertise  with filtering and sorting capabilities. "+
			"Returns a structured list of Field Supervisor Expertise  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_supervisor_expertise/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorExpertise.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_EXPERTISE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisorExpertise.Create.CMS",
		"Creates a new Field Supervisor Expertise in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Field Supervisor Expertise record with assigned unique identifier.",
		"/v1/field_supervisor_expertise/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "user_role_membership_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "sub_task_type_id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorExpertise.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_EXPERTISE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisorExpertise.Read.CMS",
		"Retrieves detailed information for a specific Field Supervisor Expertise by ID. "+
			"Returns comprehensive Field Supervisor Expertise data including specifications. "+
			"Essential for Field Supervisor Expertise specification views and data verification.",
		"/v1/field_supervisor_expertise/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorExpertise.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_EXPERTISE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldSupervisorExpertise.Delete.CMS",
		"Permanently removes a Field Supervisor Expertise record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/field_supervisor_expertise/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldSupervisorExpertise.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_SUPERVISOR_EXPERTISE.DELETE"}, 0, "default",
	)
}
