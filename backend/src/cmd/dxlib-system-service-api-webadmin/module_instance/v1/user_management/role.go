package user_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	utilsHttp "github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	partner_management_handler "github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management/handler"
)

func defineAPIRole(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("Role.List.Download.CMS",
		"Retrieves a paginated list download of Role  with filtering and sorting capabilities. "+
			"Returns a structured list download of Role  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/role/list/download", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.Role.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ROLE.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Role.List.CMS",
		"Retrieves a paginated list of Role with filtering and sorting capabilities. "+
			"Returns a structured list of Role with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/role/list", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, partner_management.ModulePartnerManagement.Role.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ROLE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("Role.Create.CMS",
		"Creates a new Role in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Role record with assigned unique identifier.",
		"/v1/role/create", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "nameid", Type: "string", Description: "Role nameId", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Role name", IsMustExist: true},
			{NameId: "description", Type: "string", Description: "Role description", IsMustExist: true},
			{NameId: "area_code", Type: "string", Description: "Role area code", IsMustExist: false},
			{NameId: "task_type_id", Type: "int64", Description: "Role task type id", IsMustExist: false},
		}, partner_management_handler.RoleCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ROLE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Role.Read.CMS",
		"Retrieves detailed information for a specific Role by ID. "+
			"Returns comprehensive Role data including specifications. "+
			"Essential for Role specification views and data verification.",
		"/v1/role/read", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.Role.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ROLE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("Role.Read.ByNameId.CMS",
		"Retrieves detailed information for a specific Role by Name Id. "+
			"Returns comprehensive Role data including specifications. "+
			"Essential for Role specification views and data verification.",
		"/v1/role/read/nameid", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "nameid", Type: "string", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.Role.RequestReadByNameId, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ROLE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("Role.Edit.CMS",
		"Updates Role information with comprehensive data validation. "+
			"Allows modification of Role specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/role/edit", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "nameid", Type: "string", Description: "Role NameId", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "Role name", IsMustExist: false},
				{NameId: "description", Type: "string", Description: "Role description", IsMustExist: false},
				{NameId: "area_code", Type: "string", Description: "Role area code", IsMustExist: false},
				{NameId: "task_type_id", Type: "int64", Description: "Role task type id", IsMustExist: false},
			}},
		}, partner_management_handler.RoleEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ROLE.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Role.Delete.CMS",
		"Permanently removes a Role record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/role/delete", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management_handler.RoleDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ROLE.DELETE"}, 0, "default",
	)

}
