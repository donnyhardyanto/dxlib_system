package partner_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	partner_management_handler "github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management/handler"
)

func defineAPIFieldExecutor(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("FieldExecutor.List.Download.CMS",
		"Retrieves a paginated list download of FieldExecutor  with filtering and sorting capabilities. "+
			"Returns a structured list download of FieldExecutor  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_executor/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutor.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutor.List.CMS",
		"Retrieves a paginated list of Field Executor  with filtering and sorting capabilities. "+
			"Returns a structured list of Field Executor  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_executor/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, partner_management_handler.FieldExecutorPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutor.Read.CMS",
		"Retrieves detailed information for a specific Field Executor by ID. "+
			"Returns comprehensive Field Executor data including specifications. "+
			"Essential for Field Executor specification views and data verification.",
		"/v1/field_executor/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutor.RequestRead, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutor.EffectiveLocations.List.CMS",
		"Retrieves a paginated list of Field Executor Effective Locations  with filtering and sorting capabilities. "+
			"Returns a structured list of Field Executor Effective Locations  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_executor_effective_location/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorEffectiveLocation.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_EFFECTIVE_LOCATION.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorLocations.List.CMS",
		"Retrieves a paginated list of Field Executor Locations  with filtering and sorting capabilities. "+
			"Returns a structured list of Field Executor Locations  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_executor_location/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorLocation.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_LOCATION.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorLocation.Create.CMS",
		"Creates a new Field Executor Location in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Field Executor Location record with assigned unique identifier.",
		"/v1/field_executor_location/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "user_role_membership_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "location_code", Type: "protected-string", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorLocation.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_LOCATION.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorLocation.Read.CMS",
		"Retrieves detailed information for a specific Field Executor Location by ID. "+
			"Returns comprehensive Field Executor Location data including specifications. "+
			"Essential for Field Executor Location specification views and data verification.",
		"/v1/field_executor_location/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorLocation.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_LOCATION.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorLocation.Delete.CMS",
		"Permanently removes a Field Executor Location record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/field_executor_location/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorLocation.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_LOCATION.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorEffectiveAreas.List.CMS",
		"Retrieves a paginated list of Field Executor Effective Areas  with filtering and sorting capabilities. "+
			"Returns a structured list of Field Executor Effective Areas  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_executor_effective_area/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorEffectiveArea.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_EFFECTIVE_AREA.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorAreas.List.CMS",
		"Retrieves a paginated list of Field Executor Areas  with filtering and sorting capabilities. "+
			"Returns a structured list of Field Executor Areas  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_executor_area/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorArea.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_AREA.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorArea.Create.CMS",
		"Creates a new Field Executor Area in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Field Executor Area record with assigned unique identifier.",
		"/v1/field_executor_area/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "user_role_membership_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "area_code", Type: "string", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorArea.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_AREA.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorArea.READ.CMS",
		"Retrieves detailed information for a specific Field Executor Area by ID. "+
			"Returns comprehensive Field Executor Area data including specifications. "+
			"Essential for Field Executor Area specification views and data verification.",
		"/v1/field_executor_area/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorArea.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_AREA.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorArea.DELETE.CMS",
		"Permanently removes a Field Executor Area record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/field_executor_area/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorArea.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_AREA.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorEffectiveExpertise.LIST.CMS",
		"Retrieves a paginated list of Field Executor Effective Expertise  with filtering and sorting capabilities. "+
			"Returns a structured list of Field Executor Effective Expertise  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_executor_effective_expertise/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorEffectiveExpertise.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_EFFECTIVE_EXPERTISE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorExpertise.LIST.CMS",
		"Retrieves a paginated list of Field Executor Expertise with filtering and sorting capabilities. "+
			"Returns a structured list of Field Executor Expertise with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/field_executor_expertise/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorExpertise.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_EXPERTISE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorExpertise.CREATE.CMS",
		"Creates a new Field Executor Expertise in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Field Executor Expertise record with assigned unique identifier.",
		"/v1/field_executor_expertise/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "user_role_membership_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "sub_task_type_id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorExpertise.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_EXPERTISE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorExpertise.READ.CMS",
		"Retrieves detailed information for a specific Field Executor Expertise by ID. "+
			"Returns comprehensive Field Executor Expertise data including specifications. "+
			"Essential for Field Executor Expertise specification views and data verification.",
		"/v1/field_executor_expertise/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorExpertise.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_EXPERTISE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("FieldExecutorExpertise.DELETE.CMS",
		"Permanently removes a Field Executor Expertise record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/field_executor_expertise/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.FieldExecutorExpertise.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FIELD_EXECUTOR_EXPERTISE.DELETE"}, 0, "default",
	)
}
