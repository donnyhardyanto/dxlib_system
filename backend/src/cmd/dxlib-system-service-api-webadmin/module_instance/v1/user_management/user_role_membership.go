package user_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	utilsHttp "github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
)

func defineAPIUserRoleMembership(anAPI *api.DXAPI) {
	anAPI.NewEndPoint("UserRole.List.CMS",
		"Retrieves a paginated list of User Role with filtering and sorting capabilities. "+
			"Returns a structured list of User Role with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/user_role_membership/list", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, partner_management.ModulePartnerManagement.UserRoleMembership.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER_ROLE_MEMBERSHIP.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("UserRole.Create.CMS",
		"Creates a new User Role in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created User Role record with assigned unique identifier.",
		"/v1/user_role_membership/create", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "user_id", Type: "int64", Description: "Privilege user_id", IsMustExist: true},
			{NameId: "organization_id", Type: "int64", Description: "Privilege organization_id", IsMustExist: true},
			{NameId: "role_id", Type: "int64", Description: "Privilege role_id", IsMustExist: true},
		}, user_management.ModuleUserManagement.UserRoleMembershipCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER_ROLE_MEMBERSHIP.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("UserRole.Delete.CMS",
		"Permanently removes a User Role record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/user_role_membership/delete", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.UserRoleMembershipHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER_ROLE_MEMBERSHIP.DELETE"}, 0, "default",
	)
}
