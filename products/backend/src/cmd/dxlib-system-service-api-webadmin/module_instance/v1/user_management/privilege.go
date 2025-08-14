package user_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	utilsHttp "github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
)

func defineAPIPrivilege(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("Privilege.Download.CMS",
		"Download Privilege to CSV or Excel file based on given filters",
		"/v1/privilege/list/download", "POST", api.EndPointTypeHTTPDownloadStream, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.PrivilegeListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"PRIVILEGE_LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Privilege.List.CMS",
		"Retrieves a paginated list of Privilege with filtering and sorting capabilities. "+
			"Returns a structured list of Privilege with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/privilege/list", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, user_management.ModuleUserManagement.PrivilegeList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"PRIVILEGE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("Privilege.Create.CMS",
		"Creates a new Privilege in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Privilege record with assigned unique identifier.",
		"/v1/privilege/create", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "nameid", Type: "string", Description: "Privilege nameId", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Privilege name", IsMustExist: true},
			{NameId: "description", Type: "string", Description: "Privilege description", IsMustExist: true},
		}, user_management.ModuleUserManagement.PrivilegeCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"PRIVILEGE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Privilege.Read.CMS",
		"Retrieves detailed information for a specific Privilege by ID. "+
			"Returns comprehensive Privilege data including specifications. "+
			"Essential for Privilege specification views and data verification.",
		"/v1/privilege/read", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.PrivilegeRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"PRIVILEGE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("Privilege.Edit.CMS",
		"Updates Privilege information with comprehensive data validation. "+
			"Allows modification of Privilege specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/privilege/edit", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "nameid", Type: "string", Description: "Privilege nameid", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "Privilege name", IsMustExist: false},
				{NameId: "description", Type: "string", Description: "Privilege description", IsMustExist: false},
			}},
		}, user_management.ModuleUserManagement.PrivilegeEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"PRIVILEGE.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Privilege.Delete",
		"Permanently removes a Privilege record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/privilege/delete", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.PrivilegeDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"PRIVILEGE.DELETE"}, 0, "default",
	)
}
