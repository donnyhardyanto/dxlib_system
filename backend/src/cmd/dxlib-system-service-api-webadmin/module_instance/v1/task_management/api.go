package task_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
)

func DefineAPIEndPoints(anAPI *api.DXAPI) {

	defineAPICustomer(anAPI)
	defineAPITask(anAPI)
	defineAPISubTask(anAPI)
	defineAPISubTaskReport(anAPI)

	anAPI.NewEndPoint("TaskType.List.Download.CMS",
		"Retrieves a paginated list download of task types  with filtering and sorting capabilities. "+
			"Returns a structured list download of task types  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/task_type/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.TaskType.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK_TYPE.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("TaskType.List.CMS",
		"Retrieves a paginated list of Task Type with filtering and sorting capabilities. "+
			"Returns a structured list of Task Type with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/task_type/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		},
		task_management.ModuleTaskManagement.TaskType.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK_TYPE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("TaskType.Create.CMS",
		"Creates a new Task Type  in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Task Type record with assigned unique identifier.",
		"/v1/task_type/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "Task Type Id", IsMustExist: true},
			{NameId: "code", Type: "string", Description: "Task Type code", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Task Type name", IsMustExist: true},
		}, task_management.ModuleTaskManagement.TaskType.RequestCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK_TYPE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("TaskType.Read.CMS",
		"Retrieves detailed information for a specific Task Type by ID. "+
			"Returns comprehensive Task Type data including specifications and flow rate parameters. "+
			"Essential for Task Type specification views and data verification.",
		"/v1/task_type/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.TaskType.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK_TYPE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("TaskType.Edit.CMS",
		"Updates Task Type information with comprehensive data validation. "+
			"Allows modification of Task Type specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/task_type/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "code", Type: "string", Description: "Task Type code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "Task Type name", IsMustExist: false},
			}},
		}, task_management.ModuleTaskManagement.TaskType.RequestEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK_TYPE.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("TaskType.Delete.CMS",
		"Permanently removes a Task Type record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/task_type/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.TaskType.RequestSoftDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK_TYPE.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskType.List.Download.CMS",
		"Retrieves a paginated list download of sub task types  with filtering and sorting capabilities. "+
			"Returns a structured list download of sub task types  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/sub_task_type/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.SubTaskType.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_TYPE.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskType.List.CMS",
		"Retrieves a paginated list of Sub Task Type with filtering and sorting capabilities. "+
			"Returns a structured list of Sub Task Type with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/sub_task_type/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		},
		task_management.ModuleTaskManagement.SubTaskType.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_TYPE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskType.Create.CMS",
		"Creates a new Sub Task Type  in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Sub Task Type record with assigned unique identifier.",
		"/v1/sub_task_type/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "Sub Task Type Id", IsMustExist: true},
			{NameId: "task_type_id", Type: "int64", Description: "Sub Task Type code", IsMustExist: true},
			{NameId: "code", Type: "string", Description: "Sub Task Type code", IsMustExist: true},
			{NameId: "name", Type: "string", Description: "Sub Task Type name", IsMustExist: true},
		}, task_management.ModuleTaskManagement.SubTaskType.RequestCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_TYPE.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskType.Read.CMS",
		"Retrieves detailed information for a specific Sub Task Type by ID. "+
			"Returns comprehensive Sub Task Type data including specifications and flow rate parameters. "+
			"Essential for Sub Task Type specification views and data verification.",
		"/v1/sub_task_type/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.SubTaskType.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_TYPE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskType.Edit.CMS",
		"Updates Sub Task Type information with comprehensive data validation. "+
			"Allows modification of Sub Task Type specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/sub_task_type/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "task_type_id", Type: "int64", Description: "Sub Task Type code", IsMustExist: false},
				{NameId: "code", Type: "string", Description: "Sub Task Type code", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "Sub Task Type name", IsMustExist: false},
			}},
		}, task_management.ModuleTaskManagement.SubTaskType.RequestEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_TYPE.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskType.Delete.CMS",
		"Permanently removes a Sub Task Type record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/sub_task_type/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.SubTaskType.RequestSoftDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_TYPE.DELETE"}, 0, "default",
	)

}
