package arrears_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	arrearsmanagementhandler "github.com/donnyhardyanto/dxlib-system/common/infrastructure/arrears_management/handler"
)

func DefineAPIEndPoints(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("PenangananPiutang.Create.CMS",
		"create new piutang "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Piutang record with assigned unique identifier.",
		"/v1/task/penanganan_piutang/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "code", Type: "string", Description: "", IsMustExist: false},
			{NameId: "customer_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "sub_task_type_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "spk_no", Type: "string", Description: "", IsMustExist: true},
			{NameId: "amount_usage_bill", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "amount_fine", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "amount_payment_guarantee", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "amount_reinstallation_cost", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "amount_reflow_cost", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "amount_bill_total", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "period_begin", Type: "date", Description: "", IsMustExist: false},
			{NameId: "period_end", Type: "date", Description: "", IsMustExist: true},
			{NameId: "status", Type: "string", Description: "", IsMustExist: true},
		},
		arrearsmanagementhandler.TaskPiutangCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.CREATE"}, 0, "default",
	)

}
