package partner_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	partner_management_handler "github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management/handler"
)

func defineAPIOrganizationExecutor(anAPI *api.DXAPI) {
	anAPI.NewEndPoint("OrganizationExecutor.List.CMS",
		"Retrieves a paginated list of Organization Executor  with filtering and sorting capabilities. "+
			"Returns a structured list of Organization Executor  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/organization_executor/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management_handler.OrganizationExecutorPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutor.Read.CMS",
		"Retrieves detailed information for a specific Organization Executor by ID. "+
			"Returns comprehensive Organization Executor data including specifications. "+
			"Essential for Organization Executor specification views and data verification.",
		"/v1/organization_executor/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutor.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutorLocations.List.CMS",
		"Retrieves a paginated list of Organization Executor Locations  with filtering and sorting capabilities. "+
			"Returns a structured list of Organization Executor Locations with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/organization_executor_location/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutorLocation.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR_LOCATION.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutorLocation.Create.CMS",
		"Creates a new Organization Executor Location in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Organization Executor Location record with assigned unique identifier.",
		"/v1/organization_executor_location/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "organization_role_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "location_code", Type: "protected-string", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutorLocation.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR_LOCATION.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutorLocation.Read.CMS",
		"Retrieves detailed information for a specific Organization Executor Location by ID. "+
			"Returns comprehensive Organization Executor Location data including specifications. "+
			"Essential for Organization Executor Location specification views and data verification.",
		"/v1/organization_executor_location/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutorLocation.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR_LOCATION.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutorLocation.Delete.CMS",
		"Permanently removes a Organization Executor Location record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/organization_executor_location/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutorLocation.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR_LOCATION.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutorAreas.List.CMS",
		"Retrieves a paginated list of Organization Executor Areas  with filtering and sorting capabilities. "+
			"Returns a structured list of Organization Executor Areas  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/organization_executor_area/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutorArea.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR_AREA.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutorArea.Create.CMS",
		"Creates a new Organization Executor Area in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Organization Executor Area record with assigned unique identifier.",
		"/v1/organization_executor_area/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "organization_role_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "area_code", Type: "string", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutorArea.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR_AREA.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutorArea.Read.CMS",
		"Retrieves detailed information for a specific Organization Executor Area by ID. "+
			"Returns comprehensive Organization Executor Area data including specifications. "+
			"Essential for Organization Executor Area specification views and data verification.",
		"/v1/organization_executor_area/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutorArea.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR_AREA.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutorArea.Delete",
		"Permanently removes a Organization Executor Area record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/organization_executor_area/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutorArea.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR_AREA.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutorExpertise.List.CMS",
		"Retrieves a paginated list of Organization Executor Expertise  with filtering and sorting capabilities. "+
			"Returns a structured list of Organization Executor Expertise  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/organization_executor_expertise/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutorExpertise.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR_EXPERTISE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutorExpertise.Create.CMS",
		"Creates a new Organization Executor Expertise in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Organization Executor Expertise record with assigned unique identifier.",
		"/v1/organization_executor_expertise/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "organization_role_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "sub_task_type_id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutorExpertise.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR_EXPERTISE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutorExpertise.Read.CMS",
		"Retrieves detailed information for a specific Organization Executor Expertise by ID. "+
			"Returns comprehensive Organization Executor Expertise data including specifications. "+
			"Essential for Organization Executor Expertise specification views and data verification.",
		"/v1/organization_executor_expertise/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutorExpertise.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR_EXPERTISE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationExecutorExpertise.Delete.CMS",
		"Permanently removes a Organization Executor Expertise record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/organization_executor_expertise/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.OrganizationExecutorExpertise.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION_EXECUTOR_EXPERTISE.DELETE"}, 0, "default",
	)
}
