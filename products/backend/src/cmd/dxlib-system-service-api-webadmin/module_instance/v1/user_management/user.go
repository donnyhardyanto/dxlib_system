package user_management

import (
	user_management2 "github.com/donnyhardyanto/dxlib-system/common/infrastructure/user_management/handler"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	utilsHttp "github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
)

func defineAPIUser(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("User.Upload.CMS",
		"Upload file csv or Excel to creates some new Users in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created User records with assigned unique identifiers.",
		"/v1/user/create_bulk", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{},
		user_management.ModuleUserManagement.UserCreateBulk, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER.UPLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("User.Download.CMS",
		"Download User to CSV or Excel file based on given filters",
		"/v1/user/list/download", "POST", api.EndPointTypeHTTPDownloadStream, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.User.RequestListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER_LIST.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("User.List.CMS",
		"Retrieves a paginated list of User with filtering and sorting capabilities. "+
			"Returns a structured list of User with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/user/list", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, user_management.ModuleUserManagement.UserList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("User.Create.CMS",
		"Creates a new User in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created User record with assigned unique identifier.",
		"/v1/user/create", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "organization_id", Type: "int64", Description: "Organization that user belong to at first time create", IsMustExist: true},
			{NameId: "role_id", Type: "int64", Description: "Role that user belong to at first time create", IsMustExist: true},
			{NameId: "loginid", Type: "string", Description: "Loginid", IsMustExist: true},
			{NameId: "email", Type: "email", Description: "Email", IsMustExist: true},
			{NameId: "fullname", Type: "string", Description: "Fullname", IsMustExist: true},
			{NameId: "phonenumber", Type: "phonenumber", Description: "Phonenumber", IsMustExist: true},
			{NameId: "attribute", Type: "string", Description: "Attribute", IsMustExist: false},
			{NameId: "identity_number", Type: "string", Description: "identity_number", IsMustExist: false},
			{NameId: "identity_type", Type: "string", Description: "identity_type", IsMustExist: false},
			{NameId: "gender", Type: "string", Description: "gender", IsMustExist: false},
			{NameId: "address_on_identity_card", Type: "string", Description: "address_on_identity_card", IsMustExist: false},
			{NameId: "membership_number", Type: "string", Description: "Attribute", IsMustExist: false},
			{NameId: "password_i", Type: "string", Description: "Password block", IsMustExist: true},
			{NameId: "password_d", Type: "string", Description: "Password block", IsMustExist: true},
		}, user_management.ModuleUserManagement.UserCreate, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("User.Read.CMS",
		"Retrieves detailed information for a specific User by ID. "+
			"Returns comprehensive User data including specifications. "+
			"Essential for User specification views and data verification.",
		"/v1/user/read", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.UserRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("User.Activate.CMS",
		"User Activation",
		"/v1/user/activate", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.UserActivate, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER.ACTIVATE"}, 0, "default",
	)

	anAPI.NewEndPoint("User.Suspend.CMS",
		"User Suspend",
		"/v1/user/suspend", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.UserSuspend, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER.SUSPEND"}, 0, "default",
	)

	anAPI.NewEndPoint("User.Edit.CMS",
		"Updates User information with comprehensive data validation. "+
			"Allows modification of User specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/user/edit", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "loginid", Type: "string", Description: "Loginid", IsMustExist: false},
				{NameId: "email", Type: "email", Description: "Email", IsMustExist: false},
				{NameId: "fullname", Type: "string", Description: "Fullname", IsMustExist: false},
				{NameId: "phonenumber", Type: "phonenumber", Description: "Phonenumber", IsMustExist: false},
				{NameId: "attribute", Type: "string", Description: "Attribute", IsMustExist: false},
				{NameId: "identity_number", Type: "string", Description: "identity_number", IsMustExist: false},
				{NameId: "identity_type", Type: "string", Description: "identity_type", IsMustExist: false},
				{NameId: "gender", Type: "string", Description: "gender", IsMustExist: false},
				{NameId: "address_on_identity_card", Type: "string", Description: "address_on_identity_card", IsMustExist: false},
				{NameId: "membership_number", Type: "string", Description: "Attribute", IsMustExist: false},
			}},
		}, user_management.ModuleUserManagement.UserEdit, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("User.Delete.CMS",
		"Permanently removes a User record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/user/delete", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.UserDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("User.ResetPassword.CMS",
		"User Reset Password",
		"/v1/user/password/reset", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "user_id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.UserResetPassword, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER.RESET_PASSWORD"}, 0, "default",
	)

	anAPI.NewEndPoint("User.IdentityCard.Update.CMS",
		"Self Identity Card  update",
		"/v1/user/identity_card/update", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "user_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "content_base64", Type: "string", Description: "File content in Base64", IsMustExist: true},
		}, user_management2.UserIdentityCardUpdateFileContentBase64, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER.ID_CARD.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("UserIdentityCard.DownloadSource.CMS",
		"User identity card download Source",
		"/v1/user/identity_card/source", "POST", api.EndPointTypeHTTPDownloadStream, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "user_id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management2.UserIdentityCardDownloadSource, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER.ID_CARD.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("User.IdentityCard.DownloadBig.CMS",
		"User identity card download Big",
		"/v1/user/identity_card/big", "POST", api.EndPointTypeHTTPDownloadStream, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "user_id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management2.UserIdentityCardDownloadBig, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER.ID_CARD.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("UserMessage.List.CMS",
		"User Message list",
		"/v1/user_message/list", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, user_management.ModuleUserManagement.UserMessage.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER_MESSAGE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("UserMessage.RequestRead.CMS",
		"User Message RequestRead",
		"/v1/user_message/read", "POST", api.EndPointTypeHTTPJSON, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, user_management.ModuleUserManagement.UserMessage.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"USER_MESSAGE.READ"}, 0, "default",
	)

}
