package user_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	utilsHttp "github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
)

func defineAPIOrganizationRoles(anAPI *api.DXAPI) {
	anAPI.NewEndPoint("OrganizationRoles.List.CMS",
		"Retrieves a paginated list of Organization Roles with filtering and sorting capabilities. "+
			"Returns a structured list of Organization Roles with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/organization_role/list", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, user_management.ModuleUserManagement.OrganizationRoles.RequestList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER_ROLE_MEMBERSHIP.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("OrganizationRoles.Create.CMS",
		"Creates a new Organization Roles in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Organization Roles record with assigned unique identifier.",
		"/v1/organization_role/create", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "organization_id", Type: "int64", Description: "Privilege organization_id", IsMustExist: true},
			{NameId: "role_id", Type: "int64", Description: "Privilege role_id", IsMustExist: true},
		}, user_management.ModuleUserManagement.OrganizationRoles.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER_ROLE_MEMBERSHIP.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Organization Roles.Delete",
		"Permanently removes a Organization Roles record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/organization_role/delete", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.OrganizationRoles.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER_ROLE_MEMBERSHIP.DELETE"}, 0, "default",
	)
}
