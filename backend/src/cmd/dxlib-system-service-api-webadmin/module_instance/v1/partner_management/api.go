package partner_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	partner_management_handler "github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management/handler"
)

func DefineAPIEndPoints(anAPI *api.DXAPI) {
	defineAPIFieldExecutor(anAPI)
	defineAPIFieldSupervisor(anAPI)
	defineAPIOrganizationExecutor(anAPI)
	defineAPIOrganizationSupervisor(anAPI)

	anAPI.NewEndPoint("TaskStatusSummary.Read.CMS",
		"Retrieve current count of tasks by status",
		"/v1/dashboard/stats_task_type_status", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, nil,
		partner_management_handler.CMSDashboardStatsTaskTypeStatus, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("TaskStatusSummary.Read.CMS",
		"Retrieve current count of tasks by status",
		"/v1/dashboard/stats_field_executor_status", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, nil,
		partner_management_handler.CMSDashboardStatsFieldExecutorStatus, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)
	anAPI.NewEndPoint("TaskStatusSummary.Read.CMS",
		"Retrieve current count of tasks by status",
		"/v1/dashboard/stats_task_location_distribution", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, nil,
		partner_management_handler.CMSDashboardStatsTaskLocationDistribution, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("TaskStatusSummary.Read.CMS",
		"Retrieve current count of tasks by status",
		"/v1/dashboard/stats_task_status_time_series", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, nil,
		partner_management_handler.CMSDashboardStatsTaskTimeSeries, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("Role Area List.CMS",
		"",
		"/v1/role_area/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, partner_management.ModulePartnerManagement.RoleArea.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("Role Area.Create.CMS",
		"",
		"/v1/role_area/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "role_id", Type: "int64", Description: "Role Id", IsMustExist: true},
			{NameId: "area_code", Type: "string", Description: "Area code", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.RoleArea.RequestCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ROLE_AREA.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Role Area.Delete",
		"",
		"/v1/role_area/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, partner_management.ModulePartnerManagement.RoleArea.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ROLE_AREA.DELETE"}, 0, "default",
	)
}
