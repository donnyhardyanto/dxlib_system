package task_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	task_management_handler "github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management/handler"
)

func defineAPISubTask(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("SubTask.Download.CMS",
		"Download Task to CSV or Excel file based on given filters",
		"/v1/sub_task/list/download", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.SubTask.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTask.List.CMS",
		"Retrieves a paginated list of Sub Task  with filtering and sorting capabilities. "+
			"Returns a structured list of Sub Task  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/sub_task/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, task_management_handler.SubTaskList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTask.List.Per.Construction.Area.CMS",
		"Retrieves a paginated list of Sub Task  with filtering and sorting capabilities. "+
			"Returns a structured list of Sub Task  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/sub_task/list/per_construction_area", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, task_management_handler.SubTaskListPerConstructionArea, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.LIST.PER_CONSTRUCTION_AREA"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTask.Read.CMS",
		"Retrieves detailed information for a specific Sub Task by ID. "+
			"Returns comprehensive Sub Task data including specifications. "+
			"Essential for Sub Task specification views and data verification.",
		"/v1/sub_task/read_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.SubTask.RequestReadByUid, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTask.Edit.CMS",
		"Updates Sub Task information with comprehensive data validation. "+
			"Allows modification of Sub Task specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/sub_task/edit_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "task_id", Type: "int64", Description: "Sub Task task_id", IsMustExist: false},
				{NameId: "sub_task_type_id", Type: "int64", Description: "Sub Task sub_task_type_id", IsMustExist: false},
				{NameId: "code", Type: "string", Description: "Sub Task code", IsMustExist: false},
				{NameId: "status", Type: "string", Description: "Sub Task sub_task_status", IsMustExist: false},
			}},
		}, task_management.ModuleTaskManagement.SubTask.RequestEditByUid, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Sub.Task.Delete.CMS",
		"Permanently removes a Sub Task record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/sub_task/delete_by_uid", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.SubTask.RequestSoftDeleteByUid, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTask.Schedule.CMS",
		"Update schedule for a specific sub task",
		"/v1/sub_task/schedule", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "sub_task_uid", Type: "string", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "Updated schedule information", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "start_date", Type: "date", Description: "Scheduled start date", IsMustExist: true},
				{NameId: "end_date", Type: "date", Description: "Scheduled end date", IsMustExist: true},
			}},
		}, task_management_handler.SubTaskSchedule, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.SCHEDULE"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTask.Assign.CMS",
		"Assign Sub Task to an Field Executor",
		"/v1/sub_task/assign", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "sub_task_uid", Type: "string", Description: "", IsMustExist: true},
			{NameId: "field_executor_user_id", Type: "int64", Description: "Updated schedule information", IsMustExist: true},
			{NameId: "at", Type: "iso8601", Description: "Updated schedule information", IsMustExist: true},
		}, task_management_handler.SubTaskAssign, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.ASSIGN"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTask.Replace.FieldExecutor",
		"Sub Task Replace Field Executor",
		"/v1/sub_task/replace_field_executor", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "sub_task_uid", Type: "string", Description: "", IsMustExist: true},
			{NameId: "at", Type: "iso8601", Description: "", IsMustExist: true},
			{NameId: "new_field_executor_user_id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management_handler.SubTaskReplaceFieldExecutor, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.REPLACE"}, 0, "default",
	)

	/* For Role CGP Operation */

	anAPI.NewEndPoint("SubTask.Finish.CGPVerifySuccess",
		"User Sub Task Finish CGP Verify Success",
		"/v1/sub_task/cgp_verify_success", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "sub_task_uid", Type: "string", Description: "", IsMustExist: true},
			{NameId: "at", Type: "iso8601", Description: "", IsMustExist: true},
		}, task_management_handler.SelfAsCGPUserSubTaskCGPVerifySuccess, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER_SUB_TASK.CGP_VERIFY_SUCCESS"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTask.Finish.CGPVerifyFail",
		"User Sub Task Finish CGP Verify Fail",
		"/v1/sub_task/cgp_verify_fail", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "sub_task_uid", Type: "string", Description: "", IsMustExist: true},
			{NameId: "at", Type: "iso8601", Description: "", IsMustExist: true},
			{NameId: "report", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "reason", Type: "string", Description: "Report Fail reason", IsMustExist: true},
				{NameId: "remedial_action", Type: "string", Description: "Report Fail Remedial action", IsMustExist: true},
				{NameId: "sub_task_report_file_ids", Type: "array-int64", Description: "Sub Task Report File Ids", IsMustExist: true},
			}},
		}, task_management_handler.SelfAsCGPUserSubTaskCGPVerifyFail, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER_SUB_TASK.CGP_VERIFY_FAIL"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTask.Finish.CGPEditAfterVerifySuccess",
		"User Sub Task Finish CGP Edit After Verify Success",
		"/v1/sub_task/cgp_edit_after_cgp_verify_success", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "sub_task_uid", Type: "string", Description: "", IsMustExist: true},
			{NameId: "at", Type: "iso8601", Description: "", IsMustExist: true},
			{NameId: "report", Type: "json", Description: "", IsMustExist: true, Children: base.SubTaskFormReportAPIEndPointParameter},
		}, task_management_handler.SelfAsCGPUserSubTaskCGPEditAfterCGPVerifySuccess, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER_SUB_TASK.CGP_EDIT_AFTER_CGP_VERIFY_SUCCESS"}, 0, "default",
	)

}
