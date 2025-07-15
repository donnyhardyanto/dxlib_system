package main

import (
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure"
	"github.com/donnyhardyanto/dxlib/app"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib/utils/os"
	"github.com/donnyhardyanto/dxlib/vault"
)

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
		"PGN Partner API WebAdmin",
		"PGN Partner 2 API WebAmin",
		true,
		"SERVICE_PGN_PARTNER_API_WEBADMIN_DEBUG",
		"abc",
	)
	app.App.Version = VersionNumber + "+" + utils.GetBuildTime()
	app.App.OnDefineConfiguration = doOnDefineConfiguration
	app.App.OnDefineSetVariables = infrastructure.DoOnDefineSetVariables
	app.App.OnDefineAPIEndPoints = doOnDefineAPIEndPoints
	_ = app.App.Run()
}
