package main

import (
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure"
	"github.com/donnyhardyanto/dxlib/app"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils/os"
	"github.com/donnyhardyanto/dxlib/vault"
	"github.com/donnyhardyanto/dxlib_module/module/push_notification"
	"time"
)

func doOnDefineConfiguration() (err error) {
	infrastructure.DefineConfiguration()
	return nil
}

func doOnExecute() error {
	app.App.RuntimeErrorGroup.Go(func() error {
		log.Log.Info("Starting push notification execution")
		for {
			select {
			case <-app.App.RuntimeErrorGroupContext.Done():
				log.Log.Info("Context done, stopping push notification execution")
				return nil
			default:
				if err := push_notification.ModulePushNotification.FCM.Execute(); err != nil {
					log.Log.Errorf(err, "Error executing push notification: %+v", err)
					// Depending on your error handling strategy, you might want to return the error here
					// return err
				}
				// Sleep for 10 seconds after each execution
				time.Sleep(10 * time.Second)
			}
		}
	})

	return nil
}

func main() {
	log.SetFormatText()
	app.App.InitVault = vault.NewHashiCorpVault(
		os.GetEnvDefaultValue("VAULT_ADDRESS", "http://127.0.0.1:8200/"),
		os.GetEnvDefaultValue("VAULT_TOKEN", "dev-vault-token"),
		"__VAULT__",
		os.GetEnvDefaultValue("VAULT_PATH", "dev-vault-path"),
	)
	app.Set("dxlib-system-service-push-notification",
		"PGN Partner Push Notification",
		"PGN Partner Push Notification",
		true,
		"SERVICE_PGN_PARTNER_PUSH_NOTIFICATION_DEBUG",
		"abc",
	)
	app.App.OnDefineConfiguration = doOnDefineConfiguration
	app.App.OnDefineSetVariables = infrastructure.DoOnDefineSetVariables
	app.App.OnExecute = doOnExecute
	_ = app.App.Run()
}
