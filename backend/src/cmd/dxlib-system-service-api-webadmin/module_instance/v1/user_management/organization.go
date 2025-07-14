package user_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
)

func defineAPIOrganization(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("Organization.Upload.CMS",
		"Upload file csv or Excel to creates some new Organizations in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Organization records with assigned unique identifiers.",
		"/v1/organization/create_bulk", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{},
		user_management.ModuleUserManagement.OrganizationCreateBulk, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION.UPLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Organization.Download.CMS",
		"Download Organization to CSV or Excel file based on given filters",
		"/v1/organization/list/download", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.Organization.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Organization.List.CMS",
		"Retrieves a paginated list of Organization with filtering and sorting capabilities. "+
			"Returns a structured list of Organization with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/organization/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, user_management.ModuleUserManagement.OrganizationList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("Organization.Create.CMS",
		"Creates a new Organization in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Organization record with assigned unique identifier.",
		"/v1/organization/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "code", Type: "protected-string", Description: "Organization code", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Organization name", IsMustExist: true},
			{NameId: "parent_id", Type: "int64", Description: "Organization parent_id", IsMustExist: true, IsNullable: true},
			{NameId: "type", Type: "protected-string", Description: "Organization type", IsMustExist: true},
			{NameId: "address", Type: "string", Description: "Organization address", IsMustExist: false},
			{NameId: "npwp", Type: "npwp", Description: "Organization NPWP", IsMustExist: false},
			{NameId: "email", Type: "email", Description: "Organization email", IsMustExist: false},
			{NameId: "phonenumber", Type: "phonenumber", Description: "Organization phonenumber", IsMustExist: false},
			{NameId: "attribute1", Type: "string", Description: "Organization attribute", IsMustExist: false, IsNullable: true},
			{NameId: "auth_source1", Type: "protected-string", Description: "Organization auth_source1", IsMustExist: false, IsNullable: true},
			{NameId: "attribute2", Type: "string", Description: "Organization attribute", IsMustExist: false, IsNullable: true},
			{NameId: "auth_source2", Type: "protected-string", Description: "Organization auth_source2", IsMustExist: false, IsNullable: true},
			{NameId: "utag", Type: "protected-string", Description: "Unique Tag for Organization", IsMustExist: false, IsNullable: true},
		}, user_management.ModuleUserManagement.OrganizationCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Organization.Read.ByName.CMS",
		"Retrieves detailed information for a specific Organization by ID. "+
			"Returns comprehensive Organization data including specifications and flow rate parameters. "+
			"Essential for Organization specification views and data verification.",
		"/v1/organization/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.OrganizationRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("Organization.ReadByName.CMS",
		"Retrieves detailed information for a specific Organization by Name. "+
			"Returns comprehensive Organization data including specifications"+
			"Essential for Organization specification views and data verification.",
		"/v1/organization/read/name", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "name", Type: "string", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.OrganizationReadByName, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION.READ_BY_NAME"}, 0, "default",
	)

	anAPI.NewEndPoint("Organization.ReadByUtag.CMS",
		"Retrieves detailed information for a specific Organization by Utag. "+
			"Returns comprehensive Organization data including specifications"+
			"Essential for Organization specification views and data verification.",
		"/v1/organization/read/utag", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "utag", Type: "string", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.Organization.RequestReadByUtag, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION.READ_BY_UTAG"}, 0, "default",
	)

	anAPI.NewEndPoint("Organization.Edit.CMS",
		"Updates Organization information with comprehensive data validation. "+
			"Allows modification of Organization specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/organization/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "code", Type: "protected-string", Description: "Organization code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "Organization name", IsMustExist: false},
				{NameId: "parent_id", Type: "int64", Description: "Organization parent_id", IsMustExist: false},
				{NameId: "type", Type: "protected-string", Description: "Organization type", IsMustExist: false},
				{NameId: "address", Type: "string", Description: "Organization address", IsMustExist: false},
				{NameId: "npwp", Type: "npwp", Description: "Organization NPWP", IsMustExist: false},
				{NameId: "email", Type: "email", Description: "Organization email", IsMustExist: false},
				{NameId: "phonenumber", Type: "phonenumber", Description: "Organization phonenumber", IsMustExist: false},
				{NameId: "status", Type: "protected-string", Description: "Organization status", IsMustExist: false},
				{NameId: "attribute1", Type: "string", Description: "Organization attribute", IsMustExist: false, IsNullable: true},
				{NameId: "auth_source1", Type: "protected-string", Description: "Organization auth_source1", IsMustExist: false, IsNullable: true},
				{NameId: "attribute2", Type: "string", Description: "Organization attribute", IsMustExist: false, IsNullable: true},
				{NameId: "auth_source2", Type: "stprotected-stringring", Description: "Organization auth_source2", IsMustExist: false, IsNullable: true},
			}},
		}, user_management.ModuleUserManagement.OrganizationEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Organization.Delete",
		"Permanently removes a Organization record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/organization/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.OrganizationDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ORGANIZATION.DELETE"}, 0, "default",
	)

}
