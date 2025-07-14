package master_data

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/master_data"
)

func DefineAPIEndPoints(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("Area.List.Download.CMS",
		"Retrieves a paginated list download of Area  with filtering and sorting capabilities. "+
			"Returns a structured list download of Area  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/area/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.Area.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"AREA.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Area.List.CMS",
		"Retrieves a paginated list of Area  with filtering and sorting capabilities. "+
			"Returns a structured list of Area  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/area/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, master_data.ModuleMasterData.Area.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"AREA.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("Area.Create.CMS",
		"Creates a new Area in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Area record with assigned unique identifier.",
		"/v1/area/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "type", Type: "string", Description: "Area type", IsMustExist: true},
			{NameId: "code", Type: "string", Description: "Area code", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Area name", IsMustExist: true},
			{NameId: "value", Type: "string", Description: "Area value", IsMustExist: true},
			{NameId: "description", Type: "string", Description: "Area description", IsMustExist: false},
			{NameId: "cost_center", Type: "string", Description: "Area cost_center", IsMustExist: true},
			{NameId: "parent_group", Type: "string", Description: "Area parent_group", IsMustExist: false},
			{NameId: "parent_value", Type: "string", Description: "Area parent_value", IsMustExist: false},
			{NameId: "status", Type: "string", Description: "Area status", IsMustExist: true},
			{NameId: "attribute1", Type: "string", Description: "Area attribute1", IsMustExist: false},
			{NameId: "attribute2", Type: "string", Description: "Area attribute2", IsMustExist: false},
			{NameId: "attribute3", Type: "string", Description: "Area attribute3", IsMustExist: false},
			{NameId: "attribute4", Type: "string", Description: "Area attribute4", IsMustExist: false},
			{NameId: "attribute5", Type: "string", Description: "Area attribute5", IsMustExist: false},
		}, master_data.ModuleMasterData.Area.RequestCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"AREA.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Area.Read.CMS",
		"Retrieves detailed information for a specific Area by ID. "+
			"Returns comprehensive Area data including specifications. "+
			"Essential for Area specification views and data verification.",
		"/v1/area/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.Area.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"AREA.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("Area.Edit.CMS",
		"Updates Area information with comprehensive data validation. "+
			"Allows modification of Area specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/area/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "type", Type: "string", Description: "Area type", IsMustExist: false},
				{NameId: "code", Type: "string", Description: "Area code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "Area name", IsMustExist: false},
				{NameId: "value", Type: "string", Description: "Area value", IsMustExist: false},
				{NameId: "description", Type: "string", Description: "Area description", IsMustExist: false},
				{NameId: "cost_center", Type: "string", Description: "Area cost_center", IsMustExist: false},
				{NameId: "parent_group", Type: "string", Description: "Area parent_group", IsMustExist: false},
				{NameId: "parent_value", Type: "string", Description: "Area parent_value", IsMustExist: false},
				{NameId: "status", Type: "string", Description: "Area status", IsMustExist: false},
				{NameId: "attribute1", Type: "string", Description: "Area attribute1", IsMustExist: false},
				{NameId: "attribute2", Type: "string", Description: "Area attribute2", IsMustExist: false},
				{NameId: "attribute3", Type: "string", Description: "Area attribute3", IsMustExist: false},
				{NameId: "attribute4", Type: "string", Description: "Area attribute4", IsMustExist: false},
				{NameId: "attribute5", Type: "string", Description: "Area attribute5", IsMustExist: false},
			}},
		}, master_data.ModuleMasterData.Area.RequestEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"AREA.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Area.Delete.CMS",
		"Permanently removes a Area record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/area/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.Area.RequestSoftDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"AREA.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("Location.List.Download.CMS",
		"Retrieves a paginated list download of Location  with filtering and sorting capabilities. "+
			"Returns a structured list download of Location  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/location/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.Location.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"LOCATION.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Location.List.CMS",
		"Retrieves a paginated list of Location  with filtering and sorting capabilities. "+
			"Returns a structured list of Location  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/location/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, master_data.ModuleMasterData.Location.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"LOCATION.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("Location.Create.CMS",
		"Creates a new Location in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Location record with assigned unique identifier.",
		"/v1/location/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "type", Type: "string", Description: "Location type", IsMustExist: true},
			{NameId: "code", Type: "string", Description: "Location code", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Location name", IsMustExist: true},
			{NameId: "value", Type: "string", Description: "Location value", IsMustExist: true},
			{NameId: "description", Type: "string", Description: "Location description", IsMustExist: false},
			{NameId: "parent_group", Type: "string", Description: "Location parent_group", IsMustExist: false},
			{NameId: "parent_value", Type: "string", Description: "Location parent_value", IsMustExist: false},
			{NameId: "status", Type: "string", Description: "Location status", IsMustExist: true},
			{NameId: "attribute1", Type: "string", Description: "Location attribute1", IsMustExist: false},
		}, master_data.ModuleMasterData.Location.RequestCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"LOCATION.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Location.Read.CMS",
		"Retrieves detailed information for a specific Location by ID. "+
			"Returns comprehensive Location data including specifications. "+
			"Essential for Location specification views and data verification.",
		"/v1/location/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.Location.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"LOCATION.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("Location.Edit.CMS",
		"Updates Location information with comprehensive data validation. "+
			"Allows modification of Location specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/location/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "type", Type: "string", Description: "Location type", IsMustExist: false},
				{NameId: "code", Type: "string", Description: "Location code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "Location name", IsMustExist: false},
				{NameId: "value", Type: "string", Description: "Location value", IsMustExist: false},
				{NameId: "description", Type: "string", Description: "Location description", IsMustExist: false},
				{NameId: "parent_group", Type: "string", Description: "Location parent_group", IsMustExist: false},
				{NameId: "parent_value", Type: "string", Description: "Location parent_value", IsMustExist: false},
				{NameId: "status", Type: "string", Description: "Location status", IsMustExist: false},
				{NameId: "attribute1", Type: "string", Description: "Location attribute1", IsMustExist: false},
			}},
		}, master_data.ModuleMasterData.Location.RequestEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"LOCATION.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Location.Delete.CMS",
		"Permanently removes a Location record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/location/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.Location.RequestSoftDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"LOCATION.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("CustomerRef.List.Download.CMS",
		"Retrieves a paginated list download of Customer Ref  with filtering and sorting capabilities. "+
			"Returns a structured list download of Customer Ref  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/customer_ref/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.CustomerRef.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER_REF.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("CustomerRef.List.CMS",
		"Retrieves a paginated list of Customer Ref  with filtering and sorting capabilities. "+
			"Returns a structured list of Customer Ref  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/customer_ref/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, master_data.ModuleMasterData.CustomerRef.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER_REF.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("CustomerRef.Create.CMS",
		"Creates a new Customer Ref in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Customer Ref record with assigned unique identifier.",
		"/v1/customer_ref/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "type", Type: "string", Description: "Customer Ref type", IsMustExist: true},
			{NameId: "code", Type: "string", Description: "Customer Ref code", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Customer Ref name", IsMustExist: true},
			{NameId: "value", Type: "string", Description: "Customer Ref value", IsMustExist: true},
			{NameId: "description", Type: "string", Description: "Customer Ref description", IsMustExist: false},
			{NameId: "parent_group", Type: "string", Description: "Customer Ref parent_group", IsMustExist: false},
			{NameId: "parent_value", Type: "string", Description: "Customer Ref parent_value", IsMustExist: false},
			{NameId: "status", Type: "string", Description: "Customer Ref status", IsMustExist: true},
		}, master_data.ModuleMasterData.CustomerRef.RequestCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER_REF.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("CustomerRef.Read.CMS",
		"Retrieves detailed information for a specific Customer Ref by ID. "+
			"Returns comprehensive Customer Ref data including specifications. "+
			"Essential for Customer Ref specification views and data verification.",
		"/v1/customer_ref/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.CustomerRef.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER_REF.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("CustomerRef.Edit.CMS",
		"Updates Customer Ref information with comprehensive data validation. "+
			"Allows modification of Customer Ref specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/customer_ref/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "type", Type: "string", Description: "Customer Ref type", IsMustExist: false},
				{NameId: "code", Type: "string", Description: "Customer Ref code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "Customer Ref name", IsMustExist: false},
				{NameId: "value", Type: "string", Description: "Customer Ref value", IsMustExist: false},
				{NameId: "description", Type: "string", Description: "Customer Ref description", IsMustExist: false},
				{NameId: "parent_group", Type: "string", Description: "Customer Ref parent_group", IsMustExist: false},
				{NameId: "parent_value", Type: "string", Description: "Customer Ref parent_value", IsMustExist: false},
				{NameId: "status", Type: "string", Description: "Customer Ref status", IsMustExist: false},
			}},
		}, master_data.ModuleMasterData.CustomerRef.RequestEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER_REF.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("CustomerRef.Delete.CMS",
		"Permanently removes a Customer Ref record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/customer_ref/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.CustomerRef.RequestSoftDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER_REF.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("GlobalLookup.List.Download.CMS",
		"Retrieves a paginated list download of Global Lookup  with filtering and sorting capabilities. "+
			"Returns a structured list download of Global Lookup  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/global_lookup/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.GlobalLookup.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GLOBAL_LOOKUP.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("GlobalLookup.List.CMS",
		"Retrieves a paginated list of Global Lookup  with filtering and sorting capabilities. "+
			"Returns a structured list of Global Lookup  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/global_lookup/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, master_data.ModuleMasterData.GlobalLookup.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GLOBAL_LOOKUP.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("GlobalLookup.Create.CMS",
		"Creates a new Global Lookup in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Global Lookup record with assigned unique identifier.",
		"/v1/global_lookup/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "type", Type: "string", Description: "Global Lookup type", IsMustExist: true},
			{NameId: "code", Type: "string", Description: "Global Lookup code", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Global Lookup name", IsMustExist: true},
			{NameId: "value", Type: "string", Description: "Global Lookup value", IsMustExist: true},
			{NameId: "description", Type: "string", Description: "Global Lookup description", IsMustExist: false},
			{NameId: "parent_group", Type: "string", Description: "Global Lookup parent_group", IsMustExist: false},
			{NameId: "parent_value", Type: "string", Description: "Global Lookup parent_value", IsMustExist: false},
			{NameId: "status", Type: "string", Description: "Global Lookup status", IsMustExist: true},
		}, master_data.ModuleMasterData.GlobalLookup.RequestCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GLOBAL_LOOKUP.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("GlobalLookup.Read.CMS",
		"Retrieves detailed information for a specific Global Lookup by ID. "+
			"Returns comprehensive Global Lookup data including specifications. "+
			"Essential for Global Lookup specification views and data verification.",
		"/v1/global_lookup/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.GlobalLookup.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GLOBAL_LOOKUP.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("GlobalLookup.Edit.CMS",
		"Updates Global Lookup information with comprehensive data validation. "+
			"Allows modification of Global Lookup specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/global_lookup/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "type", Type: "string", Description: "Global Lookup type", IsMustExist: false},
				{NameId: "code", Type: "string", Description: "Global Lookup code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "Global Lookup name", IsMustExist: false},
				{NameId: "value", Type: "string", Description: "Global Lookup value", IsMustExist: false},
				{NameId: "description", Type: "string", Description: "Global Lookup description", IsMustExist: false},
				{NameId: "parent_group", Type: "string", Description: "Global Lookup parent_group", IsMustExist: false},
				{NameId: "parent_value", Type: "string", Description: "Global Lookup parent_value", IsMustExist: false},
				{NameId: "status", Type: "string", Description: "Global Lookup status", IsMustExist: false},
			}},
		}, master_data.ModuleMasterData.GlobalLookup.RequestEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GLOBAL_LOOKUP.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("GlobalLookup.Delete.CMS",
		"Permanently removes a Global Lookup record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/global_lookup/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.GlobalLookup.RequestSoftDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GLOBAL_LOOKUP.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("CustomerSegment.List.Download.CMS",
		"Retrieves a paginated list download of Customer Segment  with filtering and sorting capabilities. "+
			"Returns a structured list download of Customer Segment  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/customer_segment/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.CustomerSegment.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER_SEGMENT.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("CustomerSegment.List.CMS",
		"Retrieves a paginated list of Customer Segment  with filtering and sorting capabilities. "+
			"Returns a structured list of Customer Segment  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/customer_segment/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		},
		master_data.ModuleMasterData.CustomerSegment.RequestList,
		nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER_SEGMENT.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("CustomerType.List.Download.CMS",
		"Retrieves a paginated list download of Customer Type  with filtering and sorting capabilities. "+
			"Returns a structured list download of Customer Type  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/customer_type/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.CustomerType.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER_TYPE.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("CustomerType.List.CMS",
		"Retrieves a paginated list of Customer Type  with filtering and sorting capabilities. "+
			"Returns a structured list of Customer Type  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/customer_type/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		},
		master_data.ModuleMasterData.CustomerType.RequestList, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER_TYPE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("PaymentSchema.List.Download.CMS",
		"Retrieves a paginated list download of Payment Schema  with filtering and sorting capabilities. "+
			"Returns a structured list download of Payment Schema  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/payment_schema/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.PaymentScheme.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"PAYMENT_SCHEMA.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("PaymentSchema.List.CMS",
		"Retrieves a paginated list of Payment Schema  with filtering and sorting capabilities. "+
			"Returns a structured list of Payment Schema  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/payment_schema/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, master_data.ModuleMasterData.PaymentScheme.RequestList,
		nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"PAYMENT_SCHEMA.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("RSCustomerSector.List.Download.CMS",
		"Retrieves a paginated list download of RS Customer Sector  with filtering and sorting capabilities. "+
			"Returns a structured list download of RS Customer Sector  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/rs_customer_sector/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.RsCustomerSector.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"RS_CUSTOMER_SECTOR.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("RSCustomerSector.List.CMS",
		"Retrieves a paginated list of RS Customer Sector  with filtering and sorting capabilities. "+
			"Returns a structured list of RS Customer Sector  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/rs_customer_sector/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, master_data.ModuleMasterData.RsCustomerSector.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"RS_CUSTOMER_SECTOR.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("RSCustomerSector.Create.CMS",
		"Creates a new RS Customer Sector in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created RS Customer Sector record with assigned unique identifier.",
		"/v1/rs_customer_sector/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "code", Type: "string", Description: "RS Customer Sector code", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "RS Customer Sector name", IsMustExist: true},
		}, master_data.ModuleMasterData.RsCustomerSector.RequestCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"RS_CUSTOMER_SECTOR.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("RSCustomerSector.Read.CMS",
		"Retrieves detailed information for a specific RS Customer Sector by ID. "+
			"Returns comprehensive RS Customer Sector data including specifications. "+
			"Essential for RS Customer Sector specification views and data verification.",
		"/v1/rs_customer_sector/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.RsCustomerSector.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"RS_CUSTOMER_SECTOR.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("RSCustomerSector.Edit.CMS",
		"Updates RS Customer Sector information with comprehensive data validation. "+
			"Allows modification of RS Customer Sector specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/rs_customer_sector/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "code", Type: "string", Description: "RS Customer Sector code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "RS Customer Sector name", IsMustExist: false},
			}},
		}, master_data.ModuleMasterData.RsCustomerSector.RequestEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"RS_CUSTOMER_SECTOR.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("RS Customer Sector Delete",
		"Permanently removes a RS Customer Sector record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/rs_customer_sector/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, master_data.ModuleMasterData.RsCustomerSector.RequestSoftDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"RS_CUSTOMER_SECTOR.DELETE"}, 0, "default",
	)
}
