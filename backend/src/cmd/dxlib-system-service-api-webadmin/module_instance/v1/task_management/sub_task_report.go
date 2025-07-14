package task_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	task_management_handler "github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management/handler"
)

func defineAPISubTaskReport(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("SubTask.Report.List.CMS",
		"Retrieves a paginated list of Sub Task Report with advanced filtering and sorting capabilities. "+
			"Returns a structured list of Sub Task Report records including basic information and status. "+
			"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
		"/v1/sub_task_report/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, task_management_handler.SubTaskReportList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_REPORT.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTask.Report.List.Per.Construction.Area.CMS",
		"Retrieves a paginated list of Sub Task Report with advanced filtering and sorting capabilities. "+
			"Returns a structured list of Sub Task Report records including basic information and status. "+
			"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
		"/v1/sub_task_report/list/per_construction_area", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, task_management_handler.SubTaskReportListPerConstructionArea, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_REPORT.LIST.PER_CONSTRUCTION_AREA"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskReport.Read.CMS",
		"Retrieves comprehensive information for a specific Sub Task Report by ID. "+
			"Returns detailed Sub Task Report data including personal information, contact details, "+
			"location data, and current status. Essential for Sub Task Report profile views and data verification.",
		"/v1/sub_task_report/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management_handler.SubTaskReportRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_REPORT.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskReport.FileGroup.List.CMS",
		"Retrieves a paginated list of Sub Task Report File Group with advanced filtering and sorting capabilities. "+
			"Returns a structured list of Sub Task Report File Group records including basic information and status. "+
			"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
		"/v1/sub_task_report_file_group/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, task_management.ModuleTaskManagement.SubTaskReportFileGroup.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_REPORT_FILE_GROUP.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskReport.Picture.Upload.CMS",
		"Upload suc task report  picture by system to a system storage",
		"/v1/sub_task_report/picture/update", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "sub_task_report_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "sub_task_report_file_group_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "at", Type: "time", Description: "", IsMustExist: false},
			{NameId: "longitude", Type: "float64", Description: "", IsMustExist: false},
			{NameId: "latitude", Type: "float64", Description: "", IsMustExist: false},
			{NameId: "content_base64", Type: "string", Description: "File content in Base64", IsMustExist: true},
		}, task_management_handler.SubTaskReportPictureUpdateFileContentBase64, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_REPORT_PICTURE.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Sub Task Report Picture Assign Existing",
		"Sub Task Report Picture Assign Existing",
		"/v1/sub_task_report/picture/assign", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "sub_task_report_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "sub_task_file_id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management_handler.SubTaskReportPictureAssignExisting, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_REPORT_PICTURE.ASSIGN_EXISTING"}, 0, "default",
	)

	anAPI.NewEndPoint("Sub Task Report Picture Download Source",
		"Sub Task Report Picture download source",
		"/v1/sub_task_report/picture/source", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management_handler.SubTaskReportPictureDownloadSource, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_REPORT_PICTURE.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskReport.Picture.DownloadSmall.CMS",
		"Download SubTask Report Picture with small size picture",
		"/v1/sub_task_report/picture/small", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management_handler.SubTaskReportPictureDownloadSmall, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_REPORT_PICTURE.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskReport.Picture.DownloadMedium.CMS",
		"Download SubTask Report Picture with medium size picture",
		"/v1/sub_task_report/picture/medium", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management_handler.SubTaskReportPictureDownloadMedium, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_REPORT_PICTURE.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskReport.Picture.DownloadBig.CMS",
		"Download SubTask Report Picture with big size picture",
		"/v1/sub_task_report/picture/big", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management_handler.SubTaskReportPictureDownloadBig, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_REPORT_PICTURE.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTaskReport.Picture.List.CMS",
		"Retrieves a paginated list of Sub Report Picture with advanced filtering and sorting capabilities. "+
			"Returns a structured list of Sub Report Picture records including basic information and status. "+
			"Supports complex filtering conditions and customizable result ordering for efficient data retrieval.",
		"/v1/sub_task_report/picture/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "sub_task_report_id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management_handler.SubTaskReportPictureListBySubTaskReportId, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"SUB_TASK_REPORT_PICTURE.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("SubTask.Report.GenerateBeritaAcara.CMS",
		"Sub Task Report Generate Berita Acara",
		"/v1/sub_task_report/berita_acara/create", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management_handler.GenerateBeritaAcara, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("SubTask.Report.DownloadBeritaAcara.CMS",
		"Sub Task Report Download Berita Acara",
		"/v1/sub_task_report/berita_acara/download", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "uid", Type: "string", Description: "", IsMustExist: true},
		}, task_management_handler.DownloadBeritaAcara, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)
}
