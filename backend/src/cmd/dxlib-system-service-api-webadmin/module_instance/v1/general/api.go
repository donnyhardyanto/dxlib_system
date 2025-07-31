package general

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/general"
	"github.com/donnyhardyanto/dxlib_module/module/self"
)

func DefineAPIEndPoints(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("Announcement.List.Download.CMS",
		"Retrieves a paginated list download of Announcement  with filtering and sorting capabilities. "+
			"Returns a structured list download of Announcement  with their basic information "+
			"filtering conditions and flexible result ordering.",
		"/v1/announcement/list/download", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, general.ModuleGeneral.Announcement.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ANNOUNCEMENT.LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Announcement.List.CMS",
		"Retrieves a paginated list of Announcement  with filtering and sorting capabilities. "+
			"Returns a structured list of Announcement  with their basic information "+
			"filtering conditions and flexible result ordering.",
		"/v1/announcement/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64zp", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64zp", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, general.ModuleGeneral.AnnouncementList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ANNOUNCEMENT.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("Announcement.Create.CMS",
		"Creates a new Announcement in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Announcement record with assigned unique identifier.",
		"/v1/announcement/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "title", Type: "string", Description: " title", IsMustExist: true},
			{NameId: "content", Type: "string", Description: " content", IsMustExist: true},
		}, general.ModuleGeneral.AnnouncementCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ANNOUNCEMENT.CREATE"}, 0, "default",
	)
	anAPI.NewEndPoint("Announcement.Read.CMS",
		"Retrieves detailed information for a specific Announcement by ID. "+
			"Returns comprehensive Announcement data including specifications. "+
			"Essential for Announcement specification views and data verification.",
		"/v1/announcement/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64p", Description: "", IsMustExist: true},
		}, general.ModuleGeneral.AnnouncementRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ANNOUNCEMENT.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("Announcement.Edit.CMS",
		"Updates Announcement information with comprehensive data validation. "+
			"Allows modification of Announcement specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/announcement/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64p", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "title", Type: "string", Description: " title", IsMustExist: false},
				{NameId: "content", Type: "string", Description: " content", IsMustExist: false},
			}},
		}, general.ModuleGeneral.AnnouncementEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ANNOUNCEMENT.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Announcement.Delete.CMS",
		"Permanently removes a Announcement record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/announcement/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64p", Description: "", IsMustExist: true},
		}, general.ModuleGeneral.AnnouncementDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ANNOUNCEMENT.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("Announcement.Picture.Update.CMS",
		"Updates Announcement Picture information with comprehensive data validation. "+
			"Allows modification of Announcement Picture specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/announcement/picture/update", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64p", Description: "", IsMustExist: true},
			{NameId: "content_base64", Type: "string", Description: "File content in Base64", IsMustExist: true},
		}, general.ModuleGeneral.AnnouncementPictureUpdateFileContentBase64, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ANNOUNCEMENT.UPLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Announcement.Picture.DownloadSource.CMS",
		"Download an Announcement Picture with origin size",
		"/v1/announcement/picture/source", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64p", Description: "", IsMustExist: true},
		}, general.ModuleGeneral.AnnouncementPictureDownloadSource, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ANNOUNCEMENT.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Announcement.Picture.DownloadSmall.CMS",
		"Download an Announcement Picture with small size",
		"/v1/announcement/picture/small", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64p", Description: "", IsMustExist: true},
		}, general.ModuleGeneral.AnnouncementPictureDownloadSmall, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ANNOUNCEMENT.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Announcement.Picture.DownloadMedium.CMS",
		"Download an Announcement Picture with medium size",
		"/v1/announcement/picture/medium", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64p", Description: "", IsMustExist: true},
		}, general.ModuleGeneral.AnnouncementPictureDownloadMedium, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ANNOUNCEMENT.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Announcement.Picture.DownloadBig.CMS",
		"Download an Announcement Picture with big size",
		"/v1/announcement/picture/big", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64p", Description: "", IsMustExist: true},
		}, general.ModuleGeneral.AnnouncementPictureDownloadBig, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"ANNOUNCEMENT.DOWNLOAD"}, 0, "default",
	)
}
