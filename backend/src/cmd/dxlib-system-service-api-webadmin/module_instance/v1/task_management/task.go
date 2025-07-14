package task_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	task_management_handler "github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management/handler"
)

func defineAPITask(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("Task.Download.CMS",
		"Download Task to CSV or Excel file based on given filters",
		"/v1/task/list/download", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.Task.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Task.MultiSearch.CMS",
		"Task Search Multi area, multi location",
		"/v1/task/search", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "task_type_ids", Type: "array-int64", Description: "", IsNullable: true, IsMustExist: false},
			{NameId: "search", Type: "string", Description: "By Fullname or Customer number", IsMustExist: false},
			{NameId: "location_codes", Type: "array-string", Description: "By Provinsi, Kabupaten, Kecamatan Or Kelurahan code", IsMustExist: false}, // array of string
			{NameId: "area_codes", Type: "array-string", Description: "By SOR Area code", IsMustExist: false},                                        //
			{NameId: "statuses", Type: "array-string", Description: "By SOR Area code", IsMustExist: false},                                          //
			{NameId: "sort", Type: "string", Description: "LATEST, CLOSEST_TO_TODAY, CLOSEST_TO_THE_LOCATION", IsMustExist: false},
			{NameId: "latitude", Type: "float64", Description: "Reference point latitude", IsMustExist: false},
			{NameId: "longitude", Type: "float64", Description: "Reference point longitude", IsMustExist: false},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, task_management_handler.TaskMultiSearch, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.SEARCH"}, 0, "default",
	)

	anAPI.NewEndPoint("Task.MultiSearch.Per.Construction.Area.CMS",
		"Task Search Multi area, multi location",
		"/v1/task/search/per_construction_area", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "task_type_ids", Type: "array-int64", Description: "", IsNullable: true, IsMustExist: false},
			{NameId: "search", Type: "string", Description: "By Fullname or Customer number", IsMustExist: false},
			{NameId: "location_codes", Type: "array-string", Description: "By Provinsi, Kabupaten, Kecamatan Or Kelurahan code", IsMustExist: false}, // array of string
			{NameId: "area_codes", Type: "array-string", Description: "By SOR Area code", IsMustExist: false},                                        //
			{NameId: "statuses", Type: "array-string", Description: "By SOR Area code", IsMustExist: false},                                          //
			{NameId: "sort", Type: "string", Description: "LATEST, CLOSEST_TO_TODAY, CLOSEST_TO_THE_LOCATION", IsMustExist: false},
			{NameId: "latitude", Type: "float64", Description: "Reference point latitude", IsMustExist: false},
			{NameId: "longitude", Type: "float64", Description: "Reference point longitude", IsMustExist: false},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
		}, task_management_handler.TaskMultiSearchPerConstructionArea, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.SEARCH.PER_CONSTRUCTION_AREA"}, 0, "default",
	)

	anAPI.NewEndPoint("Task.List.CMS",
		"Retrieves a paginated list of Task  with filtering and sorting capabilities. "+
			"Returns a structured list of Task  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/task/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, task_management_handler.TaskList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("Task.List.Per.Construction.Area.CMS",
		"Retrieves a paginated list of Task  with filtering and sorting capabilities. "+
			"Returns a structured list of Task  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/task/list/per_construction_area", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, task_management_handler.TaskListPerConstructionArea, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.LIST.PER_CONSTRUCTION_AREA"}, 0, "default",
	)

	anAPI.NewEndPoint("Task.Create.CMS",
		"Creates a new Task in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Task record with assigned unique identifier.",
		"/v1/task/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "code", Type: "string", Description: "Task Code", IsMustExist: true},
			{NameId: "task_type_id", Type: "int64", Description: "Task Task Type Id", IsMustExist: true},
			//	{NameId: "status", Type: "string", Description: "Task Status", IsMustExist: false},
			{NameId: "customer_id", Type: "int64", Description: "Task Customer Id", IsMustExist: true},
			{NameId: "data1", Type: "string", Description: "Task Data1", IsMustExist: false},
			{NameId: "data2", Type: "string", Description: "Task Data2", IsMustExist: false},
		}, task_management_handler.TaskCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Task.Create.CMS",
		"Creates a new Task in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Task record with assigned unique identifier.",
		"/v1/task/construction/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "code", Type: "string", Description: "Task Code", IsMustExist: true},
			{NameId: "customer_id", Type: "int64", Description: "Task Customer Id", IsMustExist: true},
			{NameId: "data1", Type: "string", Description: "Task Data1", IsMustExist: false},
			{NameId: "data2", Type: "string", Description: "Task Data2", IsMustExist: false},
		}, task_management_handler.TaskConstructionCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.CONSTRUCTION.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Task.Read.CMS",
		"Retrieves detailed information for a specific Task by ID. "+
			"Returns comprehensive Task data including specifications. "+
			"Essential for Task specification views and data verification.",
		"/v1/task/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.Task.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("Task.Edit.CMS",
		"Updates Task information with comprehensive data validation. "+
			"Allows modification of Task specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/task/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "code", Type: "string", Description: "Task Code", IsMustExist: false},
				{NameId: "status", Type: "string", Description: "Task Status", IsMustExist: false},
				{NameId: "customer_id", Type: "int64", Description: "Task Customer Id", IsMustExist: false},
				{NameId: "data1", Type: "string", Description: "Task Data1", IsMustExist: false},
				{NameId: "data2", Type: "string", Description: "Task Data2", IsMustExist: false},
			}},
		}, task_management.ModuleTaskManagement.Task.RequestEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Task.Delete.CMS",
		"Permanently removes a Task record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/task/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.Task.RequestSoftDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("Task.RelyOnSyncByUid",
		"Sync to RelyOn. The task status must be COMPLETED",
		"/v1/task/relyon_sync_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "Unique identifier of the task to retrieve", IsMustExist: true},
		}, task_management_handler.TaskRelyOnSyncByUid, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.SYNC"}, 0, "default",
	)
	anAPI.NewEndPoint("Task.TaskRelyOnSyncAll",
		"Sync All to RelyOn. The task statuses must be COMPLETED",
		"/v1/task/relyon_sync_all", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{},
		task_management_handler.TaskRelyOnSyncAll, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.SYNC"}, 0, "default",
	)

	anAPI.NewEndPoint("Customer.ConstructionTask.Upload.CMS",
		"Upload file csv or Excel to creates some new Construction Tasks in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Construction Task records with assigned unique identifiers.",
		"/v1/task/construction/create_bulk", "POST", api.EndPointTypeHTTPUploadStream, http.ContentTypeApplicationOctetStream, []api.DXAPIEndPointParameter{},
		task_management_handler.CustomerConstructionTaskCreateBulk, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.CONSTRUCTION.UPLOAD"}, 0, "default",
	)
}
