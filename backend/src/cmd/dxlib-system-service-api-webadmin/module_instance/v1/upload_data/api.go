package upload_data

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	upload_data_handler "github.com/donnyhardyanto/dxlib-system/common/infrastructure/upload_data/handler"
)

func DefineAPIEndPoints(anAPI *api.DXAPI) {

	//anAPI.NewEndPoint("Customer.CreateBulkBase64.CMS",
	//	"Upload base64 encoded file (csv or Excel) to create multiple customers in the system with validated information.",
	//	"/v1/upload_data/customer", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
	//		{NameId: "content_base64", Type: "string", Description: "File content in Base64", IsMustExist: true},
	//		{NameId: "content_type", Type: "string", Description: "Content type of the file (csv or excel)", IsMustExist: true},
	//	}, upload_data_handler.CustomerCreateBulkBase64, nil, nil, []api.DXAPIEndPointExecuteFunc{
	//		self.ModuleSelf.MiddlewareRequestRateLimitCheck,
	//		self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
	//	}, []string{"CUSTOMER.UPLOAD"}, 0, "default",
	//)

	anAPI.NewEndPoint("Data.Upload.CMS",
		"Upload base64 encoded file (csv or Excel) into temporary table.",
		"/v1/upload_data/upload", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "content_base64", Type: "string", Description: "File content in Base64", IsMustExist: true},
			{NameId: "content_type", Type: "string", Description: "Content type of the file (csv or excel)", IsMustExist: true},
			{NameId: "data_type", Type: "string", Description: "Content type of the file (csv or excel)", IsMustExist: true},
			{NameId: "option", Type: "string", Description: "Content type of the file (csv or excel)", IsMustExist: false},
		}, upload_data_handler.Upload, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.UPLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Data.Upload.CMS",
		"Upload base64 encoded file (csv or Excel) into temporary table.",
		"/v1/upload_data/status", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "job_id", Type: "string", Description: "", IsMustExist: false},
		}, upload_data_handler.GetJobStatus, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.UPLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Data.Upload.List.CMS",
		"List of uploaded data (csv or Excel) to create multiple organizations in the system with validated information.",
		"/v1/upload_data/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "data_type", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, upload_data_handler.UploadList, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.UPLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Data.Upload.Edit.CMS",
		"Edit uploaded data (csv or Excel) before create multiple organizations in the system with validated information.",
		"/v1/upload_data/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "data_type", Type: "string", Description: "", IsMustExist: true},
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "loginid", Type: "string", Description: "", IsMustExist: false},
				{NameId: "identity_number", Type: "string", Description: "", IsMustExist: false},
				{NameId: "fullname", Type: "string", Description: "", IsMustExist: false},
				{NameId: "name", Type: "string", Description: "", IsMustExist: false},
				{NameId: "email", Type: "string", Description: "", IsMustExist: false},
				{NameId: "phonenumber", Type: "string", Description: "", IsMustExist: false},
			}},
		}, upload_data_handler.UploadEdit, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.UPLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Data.Upload.Process.CMS",
		"Process uploaded data to create multiple organizations in the system with validated information.",
		"/v1/upload_data/process", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "data_type", Type: "string", Description: "", IsMustExist: true},
		}, upload_data_handler.Process, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.UPLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Job.Upload.CMS",
		"check active jobs. ",
		"/v1/upload_data/jobs/active", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			//			{NameId: "job_id", Type: "string", Description: "", IsMustExist: false},
		},
		upload_data_handler.GetActiveJobs, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.LIST"}, 0, "default",
	)
	anAPI.NewEndPoint("PenangananPiutang.Upload.CMS",
		"Upload file csv or Excel to creates some new Customer in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Customer record with assigned unique identifier.",
		"/v1/upload_data/jobs/all", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "status", Type: "string", Description: "", IsMustExist: false},
			{NameId: "user", Type: "string", Description: "", IsMustExist: false},
			{NameId: "limit", Type: "string", Description: "", IsMustExist: false},
			{NameId: "offset", Type: "string", Description: "", IsMustExist: false},
		},
		upload_data_handler.GetAllJobs, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("PenangananPiutang.Upload.CMS",
		"Upload file csv or Excel to creates some new Customer in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Customer record with assigned unique identifier.",
		"/v1/upload_data/jobs/stats", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "job_id", Type: "string", Description: "", IsMustExist: false},
		},
		upload_data_handler.GetJobStats, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.LIST"}, 0, "default",
	)
	anAPI.NewEndPoint("Job.Cancel.Upload.CMS",
		"cancel jobs",
		"/v1/upload_data/jobs/cancel", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "job_id", Type: "string", Description: "", IsMustExist: false},
		},
		upload_data_handler.CancelJob, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK.LIST"}, 0, "default",
	)

}
