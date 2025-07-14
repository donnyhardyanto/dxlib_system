package push_notification

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/push_notification"
	"github.com/donnyhardyanto/dxlib_module/module/self"
)

func DefineAPIEndPoints(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("FCMApplication.List.CMS",
		"Retrieves a paginated list of FCM Application  with filtering and sorting capabilities. "+
			"Returns a structured list of FCM Application  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/fcm_application/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, push_notification.ModulePushNotification.FCM.FCMApplication.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FCM_APPLICATION.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FCMApplication.Create.CMS",
		"Creates a new FCM Application in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created FCM Application record with assigned unique identifier.",
		"/v1/fcm_application/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "nameid", Type: "string", Description: "FCM Application nameid", IsMustExist: true},
			{NameId: "service_account_data", Type: "json-passthrough", Description: "FCM Application service account data", IsMustExist: true},
		}, push_notification.ModulePushNotification.FCM.ApplicationCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FCM_APPLICATION.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("FCMApplication.Read.CMS",
		"Retrieves detailed information for a specific FCM Application by ID. "+
			"Returns comprehensive FCM Application data including specifications. "+
			"Essential for FCM Application specification views and data verification.",
		"/v1/fcm_application/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, push_notification.ModulePushNotification.FCM.FCMApplication.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FCM_APPLICATION.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("FCMApplication.Edit.CMS",
		"Updates FCM Application information with comprehensive data validation. "+
			"Allows modification of FCM Application specifications while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/fcm_application/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "nameid", Type: "string", Description: "FCM Application nameid", IsMustExist: false},
				{NameId: "service_account_data", Type: "json", Description: "FCM Application service account data", IsMustExist: false},
			}},
		}, push_notification.ModulePushNotification.FCM.FCMApplication.RequestEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FCM_APPLICATION.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("FCMApplication.Delete.CMS",
		"Permanently removes a FCM Application record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/fcm_application/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, push_notification.ModulePushNotification.FCM.FCMApplication.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FCM_APPLICATION.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("FCMToken.List.CMS",
		"Retrieves a paginated list of FCM Token with filtering and sorting capabilities. "+
			"Returns a structured list of FCM Token with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/fcm_token/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, push_notification.ModulePushNotification.FCM.FCMUserToken.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FCM_TOKEN.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FCMToken.Read.CMS",
		"Retrieves detailed information for a specific FCM Token by ID. "+
			"Returns comprehensive FCM Token data including specifications. "+
			"Essential for FCM Token specification views and data verification.",
		"/v1/fcm_token/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, push_notification.ModulePushNotification.FCM.FCMUserToken.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FCM_TOKEN.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("FCMToken.Delete.CMS",
		"Permanently removes a FCM Token record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/fcm_token/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, push_notification.ModulePushNotification.FCM.FCMUserToken.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FCM_TOKEN.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("FCMMessage.List.CMS",
		"Retrieves a paginated list of FCM Message with filtering and sorting capabilities. "+
			"Returns a structured list of FCM Message with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/fcm_message/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, push_notification.ModulePushNotification.FCM.FCMMessage.RequestPagingList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FCM_MESSAGE.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("FCMMessage.Read.CMS",
		"Retrieves detailed information for a specific FCM Message by ID. "+
			"Returns comprehensive FCM Message data including specifications. "+
			"Essential for FCM Message specification views and data verification.",
		"/v1/fcm_message/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, push_notification.ModulePushNotification.FCM.FCMMessage.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FCM_MESSAGE.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("FCMMessage.Delete.CMS",
		"Permanently removes a FCM Message record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/fcm_message/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, push_notification.ModulePushNotification.FCM.FCMMessage.RequestHardDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FCM_MESSAGE.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("FCMMessage.CreateTestToUser.CMS",
		"",
		"/v1/fcm_message/create/test/user", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "application_nameid", Type: "string", Description: "", IsMustExist: true},
			{NameId: "user_id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "msg_title", Type: "string", Description: "", IsMustExist: true},
			{NameId: "msg_body", Type: "string", Description: "", IsMustExist: true},
			{NameId: "msg_data", Type: "json-passthrough", Description: "", IsMustExist: true},
		}, push_notification.ModulePushNotification.FCM.RequestCreateTestMessageToUser, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"FCM_MESSAGE.CREATE_TEST_MESSAGE_TO_USER"}, 0, "default",
	)
}
