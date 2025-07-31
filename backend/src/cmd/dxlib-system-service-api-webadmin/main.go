package main

import (
	moduleInstance "github.com/donnyhardyanto/dxlib-system/service-api-webadmin/module_instance"
	moduleInstanceV1AuditLog "github.com/donnyhardyanto/dxlib-system/service-api-webadmin/module_instance/v1/audit_log"
	moduleInstanceV1ExternalSystem "github.com/donnyhardyanto/dxlib-system/service-api-webadmin/module_instance/v1/external_system"
	moduleInstanceV1General "github.com/donnyhardyanto/dxlib-system/service-api-webadmin/module_instance/v1/general"
	moduleInstanceV1MasterData "github.com/donnyhardyanto/dxlib-system/service-api-webadmin/module_instance/v1/master_data"
	moduleInstanceV1PushNotification "github.com/donnyhardyanto/dxlib-system/service-api-webadmin/module_instance/v1/push_notification"
	moduleInstanceV1Self "github.com/donnyhardyanto/dxlib-system/service-api-webadmin/module_instance/v1/self"
	moduleInstanceV1UserManagement "github.com/donnyhardyanto/dxlib-system/service-api-webadmin/module_instance/v1/user_management"

	"github.com/donnyhardyanto/dxlib-system/common/infrastructure"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/app"
	"github.com/donnyhardyanto/dxlib/configuration"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib/utils/os"
	"github.com/donnyhardyanto/dxlib/vault"
	"github.com/donnyhardyanto/dxlib_module/module/oam"
)

var isAPISpec = false

func doOnDefineConfiguration() (err error) {
	infrastructure.DefineConfiguration()

	// API
	configuration.Manager.NewIfNotExistConfiguration("api", "api.json", "json", false, false, map[string]any{
		"oam": map[string]any{
			"nameid":  "oam",
			"address": os.GetEnvDefaultValue("SYSTEM_API_OAM_WEBADMIN_ADDRESS", "0.0.0.0:14000"),
		},
		"webadmin": map[string]any{
			"nameid":  "webadmin",
			"address": os.GetEnvDefaultValue("SYSTEM_API_WEBADMIN_ADDRESS", "0.0.0.0:15000"),
		},
	}, []string{})

	return nil
}

func doOnDefineAPIEndPoints() (err error) {
	err = oam.DefineAPIEndPoints(api.Manager.APIs["oam"])
	if err != nil {
		return err
	}

	apiWebadmin := api.Manager.APIs["webadmin"]
	apiWebadmin.Version = app.App.Version
	apiWebadmin.OnAuditLogStart = infrastructure.DoOnAuditLogStart
	apiWebadmin.OnAuditLogUserIdentified = infrastructure.DoOnAuditLogUserIdentified
	apiWebadmin.OnAuditLogEnd = infrastructure.DoOnAuditLogEnd

	if isAPISpec {
		apiWebadmin.NewEndPoint("PrintSpec",
			"Print the API Specification",
			"/spec", "GET", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, nil,
			apiWebadmin.APIHandlerPrintSpec, nil, nil, nil, nil,
			0, "",
		)
	}

	// Version endpoint
	apiWebadmin.NewEndPoint("Version",
		"Get API and client version information",
		"/version", "GET", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, nil,
		moduleInstance.VersionHandler, nil, nil, nil, nil, 0, "default",
	)
	moduleInstanceV1AuditLog.DefineAPIEndPoints(apiWebadmin)
	moduleInstanceV1Self.DefineAPIEndPoints(apiWebadmin)
	moduleInstanceV1General.DefineAPIEndPoints(apiWebadmin)
	moduleInstanceV1UserManagement.DefineAPIEndPoints(apiWebadmin)
	moduleInstanceV1ExternalSystem.DefineAPIEndPoints(apiWebadmin)
	//moduleInstanceV1Webapp.DefineAPIEndPoints(apiWebadmin)
	moduleInstanceV1MasterData.DefineAPIEndPoints(apiWebadmin)
	moduleInstanceV1PushNotification.DefineAPIEndPoints(apiWebadmin)
	return nil
}

var VersionNumber = "1.0.1"

func main() {
	isAPISpec = os.GetEnvDefaultValueAsBool("IS_API_SPEC", false)

	isLocal := os.GetEnvDefaultValue("IS_LOCAL", "false")
	if isLocal == "true" {
		log.SetFormatText()
		//	app.App.LocalData["user-create-no-send-email"] = true
	}
	app.App.InitVault = vault.NewHashiCorpVault(
		os.GetEnvDefaultValue("VAULT_ADDRESS", "http://127.0.0.1:8200/"),
		os.GetEnvDefaultValue("VAULT_TOKEN", "dev-vault-token"),
		"__VAULT__",
		os.GetEnvDefaultValue("VAULT_PATH", "dev-vault-path"),
	)
	app.Set("dxlib-system-service-api-webadmin",
		"DXLib System API WebAdmin",
		"DXLib System 2 API WebAmin",
		true,
		"DXLIB_SYSTEM_SERVICE_API_WEBADMIN_DEBUG",
		"abc",
	)
	app.App.Version = VersionNumber + "+" + utils.GetBuildTime()
	app.App.OnDefineConfiguration = doOnDefineConfiguration
	app.App.OnDefineSetVariables = infrastructure.DoOnDefineSetVariables
	app.App.OnDefineAPIEndPoints = doOnDefineAPIEndPoints
	_ = app.App.Run()
}
