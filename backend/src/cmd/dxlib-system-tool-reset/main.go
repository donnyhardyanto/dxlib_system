package main

import (
	"github.com/donnyhardyanto/dxlib-system/tool-reset/seed"
	"bufio"
	"fmt"
	"github.com/donnyhardyanto/dxlib/app"
	"github.com/donnyhardyanto/dxlib/configuration"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	utilsOs "github.com/donnyhardyanto/dxlib/utils/os"
	"github.com/donnyhardyanto/dxlib/vault"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure"
	"os"
	"strings"
)

const (
	ConfirmationKey1 = "sudah makan?"
	ConfirmationKey2 = "pluto"
)

var (
	bypassConfirmation = false
	deleteAndCreateDb  = false
)

func doOnDefineConfiguration() (err error) {
	t1 := utilsOs.GetEnvDefaultValue("IS_PGN_PARTNER_RESET_DELETE_AND_CREATE_DB", "false")
	if t1 == "true" {
		deleteAndCreateDb = true
	}

	t2 := utilsOs.GetEnvDefaultValueAsInt("PGN_PARTNER_RESET_BYPASS_CONFIRMATION", 0)
	if t2 == 1 {
		bypassConfirmation = true
	}
	infrastructure.DefineConfiguration()
	configStorage := *configuration.Manager.Configurations["storage"].Data

	configStorageDbConfig := configStorage["config"].(utils.JSON)
	configStorageDbConfig["is_connect_at_start"] = false
	configStorageDbConfig["must_connected"] = false

	configStorageDbTaskDispatcher := configStorage["task-dispatcher"].(utils.JSON)
	configStorageDbTaskDispatcher["is_connect_at_start"] = false
	configStorageDbTaskDispatcher["must_connected"] = false

	configStorageDbAuditLog := configStorage["auditlog"].(utils.JSON)
	configStorageDbAuditLog["is_connect_at_start"] = false
	configStorageDbAuditLog["must_connected"] = false

	configuration.Manager.NewIfNotExistConfiguration("storage", "storage.json", "json", false, false, map[string]any{
		"postgres": map[string]any{
			"nameid":              "postgres",
			"database_type":       "postgres",
			"address":             app.App.InitVault.GetStringOrDefault("DB_POSTGRES_ADDRESS", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("DB_POSTGRES_USER_NAME", ""),
			"user_password":       app.App.InitVault.GetStringOrDefault("DB_POSTGRES_USER_PASSWORD", ""),
			"database_name":       "postgres",
			"connection_options":  "sslmode=disable",
			"must_connected":      false,
			"is_connect_at_start": false,
		}}, []string{"postgres.user_name", "postgres.user_password"})

	configRedis := *configuration.Manager.Configurations["redis"].Data
	for k := range configRedis {
		configRedis[k].(utils.JSON)["must_connected"] = false
		configRedis[k].(utils.JSON)["is_connect_at_start"] = false
	}

	return nil

}

// Function to kill all connections to a specific database
func killConnections(db *sqlx.DB, dbName string) error {
	query := fmt.Sprintf(`
        SELECT pg_terminate_backend(pg_stat_activity.pid)
        FROM pg_stat_activity
        WHERE pg_stat_activity.datname = '%s'
          AND pid <> pg_backend_pid();
    `, dbName)
	_, err := db.Exec(query)
	if err != nil {
		return errors.Errorf("failed to kill connections: %w", err)
	}
	return nil
}

func dropDatabase(db *sqlx.DB, dbName string) (err error) {
	defer func() {
		if err != nil {
			log.Log.Warnf("Error drop database %s:%s", dbName, err.Error())
		}
	}()

	// Kill all connections to the target database
	err = killConnections(db, dbName)
	if err != nil {
		log.Log.Errorf(err, "Failed to kill connections: %s", err.Error())
		return err
	}

	query := fmt.Sprintf(`DROP DATABASE "%s"`, dbName)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func createDatabase(db *sqlx.DB, dbName string) error {
	query := fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)
	_, err := db.Exec(query)
	if err != nil {
		return errors.Errorf("failed to create database: %w", err)
	}
	return nil
}

func doOnAfterConfigurationStartAll() (err error) {

	if !bypassConfirmation {
		reader := bufio.NewReader(os.Stdin)

		log.Log.Warnf("Input confirmation key 1?")
		userInputConfirmationKey1, err := reader.ReadString('\n')
		if err != nil {
			log.Log.Errorf(err, "Failed to input confirmation key 1: %s", err.Error())
			return err
		}
		userInputConfirmationKey1 = strings.TrimSpace(userInputConfirmationKey1)

		log.Log.Warnf("Input the input confirmation key 2 to confirm:")
		userInputConfirmationKey2, err := reader.ReadString('\n')
		if err != nil {
			log.Log.Errorf(err, "Failed to input confirmation key 2: %s", err.Error())
			return err
		}
		userInputConfirmationKey2 = strings.TrimSpace(userInputConfirmationKey2)

		if userInputConfirmationKey1 != ConfirmationKey1 {
			err := log.Log.ErrorAndCreateErrorf("Confirmation key mismatch")
			return err
		}
		if userInputConfirmationKey2 != ConfirmationKey2 {
			err := log.Log.ErrorAndCreateErrorf("Confirmation key mismatch")
			return err
		}
	}
	log.Log.Warn("Executing wipe... START")

	dbPartnerAuditLog := database.Manager.Databases["auditlog"]
	dbPartnerTaskDispatcher := database.Manager.Databases["task-dispatcher"]
	dbPartnerConfig := database.Manager.Databases["config"]

	if deleteAndCreateDb {
		dbPostgres := database.Manager.Databases["postgres"]
		_ = dbPostgres.Connect()
		_ = dropDatabase(dbPostgres.Connection, dbPartnerTaskDispatcher.DatabaseName)
		_ = dropDatabase(dbPostgres.Connection, dbPartnerAuditLog.DatabaseName)
		_ = dropDatabase(dbPostgres.Connection, dbPartnerConfig.DatabaseName)

		_ = createDatabase(dbPostgres.Connection, dbPartnerConfig.DatabaseName)
		_ = createDatabase(dbPostgres.Connection, dbPartnerAuditLog.DatabaseName)
		_ = createDatabase(dbPostgres.Connection, dbPartnerTaskDispatcher.DatabaseName)
	}

	_, err = dbPartnerConfig.ExecuteCreateScripts()
	if err != nil {
		log.Log.Errorf(err, "Failed to connect/execute to database %s: %s", dbPartnerConfig.DatabaseName, err.Error())
		return err
	}

	_, err = dbPartnerAuditLog.ExecuteCreateScripts()
	if err != nil {
		log.Log.Errorf(err, "Failed to connect/execute to database %s: %s", dbPartnerAuditLog.DatabaseName, err.Error())
		return err

	}

	_, err = dbPartnerTaskDispatcher.ExecuteCreateScripts()
	if err != nil {
		log.Log.Errorf(err, "Failed to connect/execute to database %s: %s", dbPartnerTaskDispatcher.DatabaseName, err.Error())
		return err

	}
	log.Log.Warn("Executing wipe... DONE")

	return nil
}

func doOnExecute() (err error) {
	log.Log.Warn("Executing seed... START")

	err = seed.Seed()
	if err != nil {
		log.Log.Errorf(err, "Failed to seed database: %s", err.Error())
		return err
	}

	log.Log.Warn("Executing seed... DONE")
	return nil
}

func main() {
	//log.SetFormatJSON()
	log.SetFormatText()
	app.App.InitVault = vault.NewHashiCorpVault(
		utilsOs.GetEnvDefaultValue("VAULT_ADDRESS", "http://127.0.0.1:8200/"),
		utilsOs.GetEnvDefaultValue("VAULT_TOKEN", " dev-vault-token"),
		"__VAULT__",
		utilsOs.GetEnvDefaultValue("VAULT_PATH", "kv/data/pgn-partner-dev"),
	)
	app.Set("pgn-partner-reset",
		"PGN Partner Reset CLI",
		"PGN Partner Reset CLI",
		false,
		"pgn-partner-reset-debug",
		"abc",
	)
	app.App.OnDefineConfiguration = doOnDefineConfiguration
	app.App.OnAfterConfigurationStartAll = doOnAfterConfigurationStartAll
	app.App.OnDefineSetVariables = infrastructure.DoOnDefineSetVariables
	app.App.OnExecute = doOnExecute
	_ = app.App.Run()
}
