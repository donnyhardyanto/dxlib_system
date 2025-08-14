package external_system

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	externalsystem "github.com/donnyhardyanto/dxlib_module/module/external_system"
	"github.com/donnyhardyanto/dxlib_module/module/self"
)

func DefineAPIEndPoints(anAPI *api.DXAPI) {
	anAPI.NewEndPoint("ExternalSystem.List.CMS",
		"Retrieves a paginated list of External System with filtering and sorting capabilities. "+
			"Returns a structured list of External System with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/external_system/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, externalsystem.ModuleExternalSystem.ExternalSystemList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"EXTERNAL_SYSTEM.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("ExternalSystem.Create.CMS",
		"Creates a new External System  in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created External System record with assigned unique identifier.",
		"/v1/external_system/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "nameid", Type: "string", Description: "External system nameId", IsMustExist: true},
			{NameId: "type", Type: "string", Description: "External system type", IsMustExist: true},
			{NameId: "configuration", Type: "json-passthrough", Description: "External system configuration", IsMustExist: true},
		}, externalsystem.ModuleExternalSystem.ExternalSystemCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"EXTERNAL_SYSTEM.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("ExternalSystem.Read.CMS",
		"Retrieves detailed information for a specific External System by ID. "+
			"Returns comprehensive External System data including specifications. "+
			"Essential for External System specification views and data verification.",
		"/v1/external_system/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, externalsystem.ModuleExternalSystem.ExternalSystemRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"EXTERNAL_SYSTEM.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("ExternalSystem.Edit.CMS",
		"Updates External System information with comprehensive data validation. "+
			"Allows modification of External System specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/external_system/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "nameid", Type: "string", Description: "FCMApplication NameId", IsMustExist: false},
				{NameId: "configuration", Type: "json-passthrough", Description: "External system configuration", IsMustExist: false},
			}},
		}, externalsystem.ModuleExternalSystem.ExternalSystemEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"EXTERNAL_SYSTEM.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("External System.Delete",
		"Permanently removes a External System record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/external_system/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, externalsystem.ModuleExternalSystem.ExternalSystemDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"EXTERNAL_SYSTEM.DELETE"}, 0, "default",
	)
}
