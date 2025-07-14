package self

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
)

func DefineAPIEndPoints(anAPI *api.DXAPI) {
	anAPI.NewEndPoint("Prepare User Shared Key",
		"Shared key preparation",
		"/v1/self/prekey", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "a0", Type: "string", Description: "Public key A for verification", IsMustExist: true},
			{NameId: "a1", Type: "string", Description: "Public key A1 for ECDH", IsMustExist: true},
			{NameId: "a2", Type: "string", Description: "Public key A2 for ECDH", IsMustExist: true},
		}, self.ModuleSelf.SelfPrelogin, nil, nil, nil, nil,
		0, "/api-webadmin/login",
	)

	anAPI.NewEndPoint("Prepare User Shared Key With Captcha",
		"Shared key preparation",
		"/v1/self/prekey_captcha", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "a0", Type: "string", Description: "Public key A for verification", IsMustExist: true},
			{NameId: "a1", Type: "string", Description: "Public key A1 for ECDH", IsMustExist: true},
			{NameId: "a2", Type: "string", Description: "Public key A2 for ECDH", IsMustExist: true},
		}, self.ModuleSelf.SelfPreloginCaptcha, nil, nil, nil, nil,
		0, "/api-webadmin/login",
	)

	anAPI.NewEndPoint("User Login",
		"User login",
		"/v1/self/login", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "i", Type: "string", Description: "Pre-key index", IsMustExist: true},
			{NameId: "d", Type: "string", Description: "Login data", IsMustExist: true},
		}, self.ModuleSelf.SelfLogin, nil, nil, nil, []string{
			"ACCESS.WEB_CMS",
		}, 0, "/api-webadmin/login",
	)

	anAPI.NewEndPoint("User Login With Captcha",
		"User login with captcha",
		"/v1/self/login_captcha", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "i", Type: "string", Description: "Pre-key index", IsMustExist: true},
			{NameId: "d", Type: "string", Description: "Login data", IsMustExist: true},
		}, self.ModuleSelf.SelfLoginCaptcha, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
		}, []string{"ACCESS.WEB_CMS"}, 0, "/api-webadmin/login",
	)

	anAPI.NewEndPoint("User Logout",
		"User logout",
		"/v1/self/logout", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, nil,
		self.ModuleSelf.SelfLogout, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("Self User Change Password",
		"User change password",
		"/v1/self/password/change", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "i", Type: "string", Description: "Pre-key index", IsMustExist: true},
			{NameId: "d", Type: "string", Description: "Login data", IsMustExist: true},
		}, self.ModuleSelf.SelfPasswordChange, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("Self Token Detail",
		"Self token detail",
		"/v1/self/detail", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, nil,
		self.ModuleSelf.SelfLoginToken, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("Self Avatar Update",
		"Self avatar update",
		"/v1/self/avatar/update", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "content_base64", Type: "string", Description: "File content in Base64", IsMustExist: true},
		},
		self.ModuleSelf.SelfAvatarUpdateFileContentBase64, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("Self Avatar Download Source",
		"Self avatar download source",
		"/v1/self/avatar/source", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, nil,
		self.ModuleSelf.SelfAvatarDownloadSource, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("Self Avatar Download Small",
		"Self avatar download small",
		"/v1/self/avatar/small", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, nil,
		self.ModuleSelf.SelfAvatarDownloadSmall, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("Self Avatar Download Medium",
		"Self avatar download medium",
		"/v1/self/avatar/medium", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, nil,
		self.ModuleSelf.SelfAvatarDownloadMedium, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	/*	anAPI.NewEndPoint("Self Avatar Download Big",
		"Self avatar download big",
		"/v1/self/avatar/big", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, nil,
		self.ModuleSelf.SelfAvatarDownloadBig, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil,0,
	)*/

	anAPI.NewEndPoint("Self Profile",
		"Self Profile",
		"/v1/self/profile", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, nil,
		self.ModuleSelf.SelfProfile, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("Self Profile RequestEdit",
		"Self Profile RequestEdit",
		"/v1/self/profile/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "email", Type: "email", Description: "Self Profile RequestEdit email", IsMustExist: false},
				{NameId: "fullname", Type: "string", Description: "Self Profile RequestEdit fullname", IsMustExist: false},
				{NameId: "phonenumber", Type: "phonenumber", Description: "Self Profile RequestEdit phonenumber", IsMustExist: false},
			}},
		}, self.ModuleSelf.SelfProfileEdit, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)

	anAPI.NewEndPoint("Self Identity Card Download Big",
		"Self identity card download big",
		"/v1/self/identity_card/big", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, nil,
		self.ModuleSelf.SelfAvatarDownloadBig, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, nil, 0, "default",
	)
}
