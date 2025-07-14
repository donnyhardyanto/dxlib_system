package construction_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/construction_management"
)

func DefineAPIEndPoints(anAPI *api.DXAPI) {
	// Gas Appliance Management
	/*	anAPI.NewEndPoint("GasAppliance.List.Mobile",
		"Retrieves a paginated list of gas appliances with advanced filtering and sorting capabilities. "+
			"Returns a structured list of gas appliance records including basic information and status. "+
			"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
		"/v1/gas_appliance/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "SQL-like where clause for filtering", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "Sort order specification", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "Values for the filter conditions", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "Number of items per page", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "Zero-based page number", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "Include deleted items flag", IsMustExist: false},
		}, construction_management.ModuleConstructionManagement.GasAppliance.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GAS_APPLIANCE.LIST"},0,
	)*/

	anAPI.NewEndPoint("GasAppliance.ListAll.Mobile",
		"Retrieves a paginated list of gas appliances with advanced filtering and sorting capabilities. "+
			"Returns a structured list of gas appliance records including basic information and status. "+
			"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
		"/v1/gas_appliance/list_all", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "row_per_page", Type: "int64", Description: "Number of items per page", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "Zero-based page number", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.GasAppliance.RequestPagingListAll, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GAS_APPLIANCE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("GasAppliance.Create.Mobile",
		"Creates a new gas appliance record in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created gas appliance record with assigned unique identifier.",
		"/v1/gas_appliance/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "code", Type: "string", Description: "Unique identifier code for the gas appliance", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Display name of the gas appliance", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.GasAppliance.RequestCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GAS_APPLIANCE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("GasAppliance.Read.Mobile",
		"Retrieves detailed information for a specific gas appliance by ID. "+
			"Returns comprehensive gas appliance data including specifications and current status. "+
			"Essential for gas appliance profile views and data verification.",
		"/v1/gas_appliance/read_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the gas appliance", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.GasAppliance.RequestReadByUid, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GAS_APPLIANCE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("GasAppliance.Read.Mobile",
		"Retrieves detailed information for a specific gas appliance by ID. "+
			"Returns comprehensive gas appliance data including specifications and current status. "+
			"Essential for gas appliance profile views and data verification.",
		"/v1/gas_appliance/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "Unique identifier of the gas appliance", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.GasAppliance.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GAS_APPLIANCE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("GasAppliance.Edit.Mobile",
		"Updates gas appliance information with comprehensive data validation. "+
			"Allows modification of gas appliance details while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/gas_appliance/edit_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the gas appliance to update", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "Updated gas appliance information", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "code", Type: "string", Description: "New unique identifier code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "New display name", IsMustExist: false},
			}},
		}, construction_management.ModuleConstructionManagement.GasAppliance.RequestEditByUid, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GAS_APPLIANCE.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("GasAppliance.Delete.Mobile",
		"Permanently removes a gas appliance record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/gas_appliance/delete_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the gas appliance to delete_by_uid", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.GasAppliance.RequestHardDeleteByUid, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"GAS_APPLIANCE.DELETE"}, 0, "default",
	)

	/*	// Tapping Saddle Appliance Management
		anAPI.NewEndPoint("TappingSaddle.List.Mobile",
			"Retrieves a paginated list of tapping saddle appliances with advanced filtering and sorting capabilities. "+
				"Returns a structured list of tapping saddle records including basic information and status. "+
				"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
			"/v1/tapping_saddle_appliance/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
				{NameId: "filter_where", Type: "string", Description: "SQL-like where clause for filtering", IsMustExist: true},
				{NameId: "filter_order_by", Type: "string", Description: "Sort order specification", IsMustExist: true},
				{NameId: "filter_key_values", Type: "json-passthrough", Description: "Values for the filter conditions", IsMustExist: true},
				{NameId: "row_per_page", Type: "int64", Description: "Number of items per page", IsMustExist: true},
				{NameId: "page_index", Type: "int64", Description: "Zero-based page number", IsMustExist: true},
				{NameId: "is_deleted", Type: "bool", Description: "Include deleted items flag", IsMustExist: false},
			}, construction_management.ModuleConstructionManagement.TappingSaddleAppliance.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
				self.ModuleSelf.MiddlewareUserLogged,
			}, []string{"TAPPING_SADDLE_APPLIANCE.LIST"},0,
		)*/

	// Tapping Saddle Appliance Management
	anAPI.NewEndPoint("TappingSaddle.ListAll.Mobile",
		"Retrieves a paginated list of tapping saddle appliances with advanced filtering and sorting capabilities. "+
			"Returns a structured list of tapping saddle records including basic information and status. "+
			"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
		"/v1/tapping_saddle_appliance/list_all", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "row_per_page", Type: "int64", Description: "Number of items per page", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "Zero-based page number", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.TappingSaddleAppliance.RequestPagingListAll, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TAPPING_SADDLE_APPLIANCE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("TappingSaddle.Create.Mobile",
		"Creates a new tapping saddle appliance record in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created tapping saddle record with assigned unique identifier.",
		"/v1/tapping_saddle_appliance/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "code", Type: "string", Description: "Unique identifier code for the tapping saddle", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Display name of the tapping saddle", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.TappingSaddleAppliance.RequestCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TAPPING_SADDLE_APPLIANCE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("TappingSaddle.Read.Mobile",
		"Retrieves detailed information for a specific tapping saddle by ID. "+
			"Returns comprehensive tapping saddle data including specifications and current status. "+
			"Essential for tapping saddle profile views and data verification.",
		"/v1/tapping_saddle_appliance/read_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the tapping saddle", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.TappingSaddleAppliance.RequestReadByUid, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TAPPING_SADDLE_APPLIANCE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("TappingSaddle.Read.Mobile",
		"Retrieves detailed information for a specific tapping saddle by ID. "+
			"Returns comprehensive tapping saddle data including specifications and current status. "+
			"Essential for tapping saddle profile views and data verification.",
		"/v1/tapping_saddle_appliance/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "Unique identifier of the tapping saddle", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.TappingSaddleAppliance.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TAPPING_SADDLE_APPLIANCE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("TappingSaddle.Edit.Mobile",
		"Updates tapping saddle information with comprehensive data validation. "+
			"Allows modification of tapping saddle details while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/tapping_saddle_appliance/edit_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the tapping saddle to update", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "Updated tapping saddle information", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "code", Type: "string", Description: "New unique identifier code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "New display name", IsMustExist: false},
			}},
		}, construction_management.ModuleConstructionManagement.TappingSaddleAppliance.RequestEditByUid, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TAPPING_SADDLE_APPLIANCE.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("TappingSaddle.Delete.Mobile",
		"Permanently removes a tapping saddle record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/tapping_saddle_appliance/delete_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the tapping saddle to delete_by_uid", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.TappingSaddleAppliance.RequestHardDeleteByUid, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TAPPING_SADDLE_APPLIANCE.DELETE"}, 0, "default",
	)

	/*	// Meter Appliance Management
		anAPI.NewEndPoint("MeterApplianceType.List.Mobile",
			"Retrieves a paginated list of meter appliances type with advanced filtering and sorting capabilities. "+
				"Returns a structured list of meter appliance type records including basic information and status. "+
				"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
			"/v1/meter_appliance_type/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
				{NameId: "filter_where", Type: "string", Description: "SQL-like where clause for filtering", IsMustExist: true},
				{NameId: "filter_order_by", Type: "string", Description: "Sort order specification", IsMustExist: true},
				{NameId: "filter_key_values", Type: "json-passthrough", Description: "Values for the filter conditions", IsMustExist: true},
				{NameId: "row_per_page", Type: "int64", Description: "Number of items per page", IsMustExist: true},
				{NameId: "page_index", Type: "int64", Description: "Zero-based page number", IsMustExist: true},
				{NameId: "is_deleted", Type: "bool", Description: "Include deleted items flag", IsMustExist: false},
			}, construction_management.ModuleConstructionManagement.MeterApplianceType.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
				self.ModuleSelf.MiddlewareUserLogged,
			}, []string{"METER_APPLIANCE_TYPE.LIST"},0,
		)*/

	// Meter Appliance Management
	anAPI.NewEndPoint("MeterApplianceType.ListAll.Mobile",
		"Retrieves a paginated list of meter appliances type with advanced filtering and sorting capabilities. "+
			"Returns a structured list of meter appliance type records including basic information and status. "+
			"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
		"/v1/meter_appliance_type/list_all", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "row_per_page", Type: "int64", Description: "Number of items per page", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "Zero-based page number", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.MeterApplianceType.RequestPagingListAll, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"METER_APPLIANCE_TYPE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("MeterApplianceType.Create.Mobile",
		"Creates a new meter appliance type record in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created meter appliance type record with assigned unique identifier.",
		"/v1/meter_appliance_type/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "code", Type: "string", Description: "Unique identifier code for the meter", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Display name of the meter", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.MeterApplianceType.RequestCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"METER_APPLIANCE_TYPE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("MeterApplianceType.Read.Mobile",
		"Retrieves detailed information for a specific meter appliance type by ID. "+
			"Returns comprehensive meter appliancetype  data including specifications and current status. "+
			"Essential for meter appliance ptype rofile views and data verification.",
		"/v1/meter_appliance_type/read_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the meter", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.MeterApplianceType.RequestReadByUid, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"METER_APPLIANCE_TYPE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("MeterApplianceType.Read.Mobile",
		"Retrieves detailed information for a specific meter appliance type by ID. "+
			"Returns comprehensive meter appliancetype  data including specifications and current status. "+
			"Essential for meter appliance ptype rofile views and data verification.",
		"/v1/meter_appliance_type/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "Unique identifier of the meter", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.MeterApplianceType.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"METER_APPLIANCE_TYPE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("MeterApplianceType.Edit.Mobile",
		"Updates meter appliance type information with comprehensive data validation. "+
			"Allows modification of meter appliance type details while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/meter_appliance_type/edit_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the meter to update", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "Updated meter information", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "code", Type: "string", Description: "New unique identifier code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "New display name", IsMustExist: false},
			}},
		}, construction_management.ModuleConstructionManagement.MeterApplianceType.RequestEditByUid, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"METER_APPLIANCE_TYPE.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("MeterApplianceType.Delete.Mobile",
		"Permanently removes a meter appliance type record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/meter_appliance_type/delete_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the meter to delete_by_uid", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.MeterApplianceType.RequestHardDeleteByUid, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"METER_APPLIANCE_TYPE.DELETE"}, 0, "default",
	)

	/*	// Regulator Appliance Management
		anAPI.NewEndPoint("RegulatorAppliance.List.Mobile",
			"Retrieves a paginated list of regulator appliances with advanced filtering and sorting capabilities. "+
				"Returns a structured list of regulator appliance records including basic information and status. "+
				"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
			"/v1/regulator_appliance/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
				{NameId: "filter_where", Type: "string", Description: "SQL-like where clause for filtering", IsMustExist: true},
				{NameId: "filter_order_by", Type: "string", Description: "Sort order specification", IsMustExist: true},
				{NameId: "filter_key_values", Type: "json-passthrough", Description: "Values for the filter conditions", IsMustExist: true},
				{NameId: "row_per_page", Type: "int64", Description: "Number of items per page", IsMustExist: true},
				{NameId: "page_index", Type: "int64", Description: "Zero-based page number", IsMustExist: true},
				{NameId: "is_deleted", Type: "bool", Description: "Include deleted items flag", IsMustExist: false},
			}, construction_management.ModuleConstructionManagement.RegulatorAppliance.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
				self.ModuleSelf.MiddlewareUserLogged,
			}, []string{"REGULATOR_APPLIANCE.LIST"},0,
		)*/

	anAPI.NewEndPoint("RegulatorAppliance.Create.Mobile",
		"Creates a new regulator appliance record in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created regulator appliance record with assigned unique identifier.",
		"/v1/regulator_appliance/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "code", Type: "string", Description: "Unique identifier code for the regulator", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Display name of the regulator", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.RegulatorAppliance.RequestCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"REGULATOR_APPLIANCE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("RegulatorAppliance.Read.Mobile",
		"Retrieves detailed information for a specific regulator appliance by ID. "+
			"Returns comprehensive regulator appliance data including specifications and current status. "+
			"Essential for regulator appliance profile views and data verification.",
		"/v1/regulator_appliance/read_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the regulator", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.RegulatorAppliance.RequestReadByUid, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"REGULATOR_APPLIANCE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("RegulatorAppliance.Read.Mobile",
		"Retrieves detailed information for a specific regulator appliance by ID. "+
			"Returns comprehensive regulator appliance data including specifications and current status. "+
			"Essential for regulator appliance profile views and data verification.",
		"/v1/regulator_appliance/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "Unique identifier of the regulator", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.RegulatorAppliance.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"REGULATOR_APPLIANCE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("RegulatorAppliance.Edit.Mobile",
		"Updates regulator appliance information with comprehensive data validation. "+
			"Allows modification of regulator appliance details while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/regulator_appliance/edit_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the regulator to update", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "Updated regulator information", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "code", Type: "string", Description: "New unique identifier code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "New display name", IsMustExist: false},
			}},
		}, construction_management.ModuleConstructionManagement.RegulatorAppliance.RequestEditByUid, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"REGULATOR_APPLIANCE.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("RegulatorAppliance.Delete.Mobile",
		"Permanently removes a regulator appliance record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/regulator_appliance/delete_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the regulator to delete_by_uid", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.RegulatorAppliance.RequestHardDeleteByUid, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"REGULATOR_APPLIANCE.DELETE"}, 0, "default",
	)

	// G Size Management
	/*	anAPI.NewEndPoint("GSize.List.Mobile",
		"Retrieves a paginated list of G sizes with advanced filtering and sorting capabilities. "+
			"Returns a structured list of G size records including basic information and specifications. "+
			"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
		"/v1/g_size/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "SQL-like where clause for filtering", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "Sort order specification", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "Values for the filter conditions", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "Number of items per page", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "Zero-based page number", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "Include deleted items flag", IsMustExist: false},
		}, construction_management.ModuleConstructionManagement.GSize.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"G_SIZE.LIST"},0,
	)*/

	anAPI.NewEndPoint("GSize.ListAll.Mobile",
		"Retrieves a paginated list of G sizes with advanced filtering and sorting capabilities. "+
			"Returns a structured list of G size records including basic information and specifications. "+
			"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
		"/v1/g_size/list_all", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "row_per_page", Type: "int64", Description: "Number of items per page", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "Zero-based page number", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.GSize.RequestPagingListAll, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"G_SIZE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("GSize.Create.Mobile",
		"Creates a new G size specification in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created G size record with assigned unique identifier.",
		"/v1/g_size/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "code", Type: "string", Description: "Unique identifier code for the G size", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Display name of the G size", IsMustExist: true},
			{NameId: "qmin", Type: "float64", Description: "Minimum flow rate (m³/h)", IsMustExist: true},
			{NameId: "qmax", Type: "float64", Description: "Maximum flow rate (m³/h)", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.GSize.RequestCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"G_SIZE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("GSize.Read.Mobile",
		"Retrieves detailed information for a specific G size by ID. "+
			"Returns comprehensive G size data including specifications and flow rate parameters. "+
			"Essential for G size specification views and data verification.",
		"/v1/g_size/read_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the G size", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.GSize.RequestReadByUid, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"G_SIZE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("GSize.Read.Mobile",
		"Retrieves detailed information for a specific G size by ID. "+
			"Returns comprehensive G size data including specifications and flow rate parameters. "+
			"Essential for G size specification views and data verification.",
		"/v1/g_size/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "Unique identifier of the G size", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.GSize.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"G_SIZE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("GSize.Edit.Mobile",
		"Updates G size information with comprehensive data validation. "+
			"Allows modification of G size specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/g_size/edit_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the G size to update", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "Updated G size information", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "code", Type: "string", Description: "New unique identifier code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "New display name", IsMustExist: false},
				{NameId: "qmin", Type: "float64", Description: "New minimum flow rate (m³/h)", IsMustExist: false},
				{NameId: "qmax", Type: "float64", Description: "New maximum flow rate (m³/h)", IsMustExist: false},
			}},
		}, construction_management.ModuleConstructionManagement.GSize.RequestEditByUid, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"G_SIZE.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("GSize.Delete.Mobile",
		"Permanently removes a G size record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/g_size/delete_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the G size to delete_by_uid", IsMustExist: true},
		}, construction_management.ModuleConstructionManagement.GSize.RequestHardDeleteByUid, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"G_SIZE.DELETE"}, 0, "default",
	)

}
