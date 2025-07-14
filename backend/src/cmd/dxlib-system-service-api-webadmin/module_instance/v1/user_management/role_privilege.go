package user_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	utilsHttp "github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
)

func defineAPIRolePrivilege(anAPI *api.DXAPI) {
	anAPI.NewEndPoint("RolePrivilege.List.CMS",
		"Retrieves a paginated list of Role Privilege with filtering and sorting capabilities. "+
			"Returns a structured list of Role Privilege with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/role_privilege/list", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, user_management.ModuleUserManagement.RolePrivilegeList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ROLE_PRIVILEGE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("RolePrivilege.Create.CMS",
		"Creates a new Role Privilege in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Role Privilege record with assigned unique identifier.",
		"/v1/role_privilege/create", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "privilege_id", Type: "int64", Description: "Privilege user_id", IsMustExist: true},
			{NameId: "role_id", Type: "int64", Description: "Privilege role_id", IsMustExist: true},
		}, user_management.ModuleUserManagement.RolePrivilege.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ROLE_PRIVILEGE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("RolePrivilege.Delete.CMS",
		"Permanently removes a Role Privilege record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/role_privilege/delete", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.RolePrivilegeDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ROLE_PRIVILEGE.DELETE"}, 0, "default",
	)
}
