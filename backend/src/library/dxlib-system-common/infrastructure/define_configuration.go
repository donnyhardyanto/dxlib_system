package infrastructure

import (
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib/app"
	"github.com/donnyhardyanto/dxlib/configuration"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib/utils/os"
)

func DefineConfiguration() {
	createScriptFileFolder := os.GetEnvDefaultValue("CREATE_SCRIPT_FILE_FOLDER", "./../sql")
	configuration.Manager.NewIfNotExistConfiguration("upload", "upload.json", "json", false, false, map[string]any{
		"WORKER": map[string]any{
			"nameid":                   "WORKER",
			"process_worker_count_min": app.App.InitVault.GetIntOrDefault("UPLOAD_WORKER_PROCESS_WORKER_COUNT_MIN", 1),
			"process_worker_count_max": app.App.InitVault.GetIntOrDefault("UPLOAD_WORKER_PROCESS_WORKER_COUNT_MAX", 4),
		},
	}, nil)

	d := configuration.Manager.NewIfNotExistConfiguration("external_system", "external_system.json", "json", false, false, map[string]any{

		"MOBILE_APP1": map[string]any{
			"nameid":             "MOBILE_APP1",
			"type":               "configuration-mobile_app",
			"api_key_google_map": app.App.InitVault.GetStringOrDefault("MOBILE_APP_API_KEY_GOOGLE_MAP", ""),
			"api_key_firebase":   app.App.InitVault.GetStringOrDefault("MOBILE_APP_API_KEY_FIREBASE", ""),
		},
		"SMTP1": map[string]any{
			"nameid":       "SMTP1",
			"type":         "smtp",
			"host":         app.App.InitVault.GetStringOrDefault("SMTP_HOST", ""),
			"port":         app.App.InitVault.GetIntOrDefault("SMTP_PORT", 587),
			"username":     app.App.InitVault.GetStringOrDefault("SMTP_USERNAME", ""),
			"password":     app.App.InitVault.GetStringOrDefault("SMTP_PASSWORD", ""),
			"sender_email": app.App.InitVault.GetStringOrDefault("SMTP_SENDER_EMAIL", ""),
			"ssl":          app.App.InitVault.GetBoolOrDefault("SMTP_SSL", false),
		},
		"LDAP1": map[string]any{
			"nameid":  "LDAP1",
			"type":    "ldap",
			"address": app.App.InitVault.GetStringOrDefault("LDAP_ADDRESS", ""),
		},
		"SMS1": map[string]any{
			"nameid":                           "SMS1",
			"enabled":                          app.App.InitVault.GetBoolOrDefault("SMS_ENABLED", false),
			"type":                             "sms-md_media",
			"address":                          app.App.InitVault.GetStringOrDefault("SMS_SERVER", ""),
			"username":                         app.App.InitVault.GetStringOrDefault("SMS_USERNAME", ""),
			"password":                         app.App.InitVault.GetStringOrDefault("SMS_PASSWORD", ""),
			"is_send_when_user_created":        app.App.InitVault.GetBoolOrDefault("SMS_IS_SEND_WHEN_USER_CREATED", false),
			"is_send_when_user_reset_password": app.App.InitVault.GetBoolOrDefault("SMS_IS_SEND_WHEN_USER_RESET_PASSWORD", false),
		},
		"RELYON1": map[string]any{
			"nameid":                 "RELYON1",
			"type":                   "relyon-inbound-outbound",
			"inbound_api_key":        app.App.InitVault.GetStringOrDefault("RELYON_INBOUND_API_KEY", ""),
			"inbound_api_secret":     app.App.InitVault.GetStringOrDefault("RELYON_INBOUND_API_SECRET", ""),
			"outbound_auth_url":      app.App.InitVault.GetStringOrDefault("RELYON_OUTBOUND_AUTH_URL", ""),
			"outbound_auth_method":   app.App.InitVault.GetStringOrDefault("RELYON_OUTBOUND_AUTH_METHOD", ""),
			"outbound_auth_header":   app.App.InitVault.GetStringOrDefault("RELYON_OUTBOUND_AUTH_HEADER", ""),
			"outbound_auth_raw_body": app.App.InitVault.GetStringOrDefault("RELYON_OUTBOUND_AUTH_RAW_BODY", ""),
			"outbound_api_url_register_installation_update_type":   app.App.InitVault.GetStringOrDefault("RELYON_OUTBOUND_API_URL_REGISTER_INSTALLATION_UPDATE_TYPE", ""),
			"outbound_api_url_register_installation_update_status": app.App.InitVault.GetStringOrDefault("RELYON_OUTBOUND_API_URL_REGISTER_INSTALLATION_UPDATE_STATUS", ""),
			"outbound_api_url_cancel_subscription_milestone":       app.App.InitVault.GetStringOrDefault("RELYON_OUTBOUND_API_URL_CANCEL_SUBSCRIPTION_MILESTONE", ""),
			"outbound_api_url_cancel_subscription_create":          app.App.InitVault.GetStringOrDefault("RELYON_OUTBOUND_API_URL_CANCEL_SUBSCRIPTION_CREATE", ""),
		},
	}, []string{
		"MOBILE_APP1.api_key_google_map", "MOBILE_APP1.api_key_firebase",
		"SMTP1.username", "SMTP1.password",
		"SMS1.username", "SMS1.password",
		"RELYON.api_key", "RELYON.api_secret",
	})

	r := (*d.Data)["RELYON1"].(map[string]any)
	base.RelyOnInboundCredentials = []utils.JSON{{
		"key":    r["inbound_api_key"],
		"secret": r["inbound_api_secret"],
	}}

	configuration.Manager.NewIfNotExistConfiguration("security", "security.json", "json", false, false, map[string]any{
		"image_uploader": map[string]any{
			"max_request_size":            app.App.InitVault.GetInt64OrDefault("MAX_REQUEST_SIZE", 100*1024*1024), // 100MB
			"max_pixel_width":             app.App.InitVault.GetInt64OrDefault("MAX_PIXEL_WIDTH", 4096),
			"max_pixel_height":            app.App.InitVault.GetInt64OrDefault("MAX_PIXEL_HEIGHT", 4096),
			"max_bytes_per_pixel":         app.App.InitVault.GetInt64OrDefault("MAX_BYTES_PER_PIXEL", 10),
			"max_pixels":                  app.App.InitVault.GetInt64OrDefault("MAX_PIXELS", 40000000), // ~40MP
			"image_process_limit_seconds": app.App.InitVault.GetInt64OrDefault("IMAGE_PROCESS_LIMIT_SECOND", 5),
		},
	}, nil)

	configuration.Manager.NewIfNotExistConfiguration("storage", "storage.json", "json", false, false, map[string]any{
		"config": map[string]any{
			"nameid":              "config",
			"database_type":       app.App.InitVault.GetStringOrDefault("DB_CONFIG_DATABASE_TYPE", ""),
			"address":             app.App.InitVault.GetStringOrDefault("DB_CONFIG_ADDRESS", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("DB_CONFIG_USER_NAME", ""),
			"user_password":       app.App.InitVault.GetStringOrDefault("DB_CONFIG_USER_PASSWORD", ""),
			"database_name":       app.App.InitVault.GetStringOrDefault("DB_CONFIG_DATABASE_NAME", ""),
			"connection_options":  app.App.InitVault.GetStringOrDefault("DB_CONFIG_CONNECTION_OPTIONS", ""),
			"must_connected":      true,
			"is_connect_at_start": true,
			"create_script_files": []string{createScriptFileFolder + "/db_configuration.sql"},
		},
		"db_base": map[string]any{
			"nameid":              "db_base",
			"database_type":       app.App.InitVault.GetStringOrDefault("DB_BASE_DATABASE_TYPE", ""),
			"address":             app.App.InitVault.GetStringOrDefault("DB_BASE_ADDRESS", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("DB_BASE_USER_NAME", ""),
			"user_password":       app.App.InitVault.GetStringOrDefault("DB_BASE_USER_PASSWORD", ""),
			"database_name":       app.App.InitVault.GetStringOrDefault("DB_BASE_DATABASE_NAME", ""),
			"connection_options":  app.App.InitVault.GetStringOrDefault("DB_BASE_CONNECTION_OPTIONS", ""),
			"must_connected":      true,
			"is_connect_at_start": true,
			"create_script_files": []string{
				createScriptFileFolder + "/db_base.general.sql",
				createScriptFileFolder + "/db_base.user_management.sql",
				createScriptFileFolder + "/db_base.push_notification.sql",
				createScriptFileFolder + "/db_base.user_management.init-data.sql",
			},
		},
		"auditlog": map[string]any{
			"nameid":              "auditlog",
			"database_type":       app.App.InitVault.GetStringOrDefault("DB_AUDITLOG_DATABASE_TYPE", ""),
			"address":             app.App.InitVault.GetStringOrDefault("DB_AUDITLOG_ADDRESS", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("DB_AUDITLOG_USER_NAME", ""),
			"user_password":       app.App.InitVault.GetStringOrDefault("DB_AUDITLOG_USER_PASSWORD", ""),
			"database_name":       app.App.InitVault.GetStringOrDefault("DB_AUDITLOG_DATABASE_NAME", ""),
			"connection_options":  app.App.InitVault.GetStringOrDefault("DB_AUDITLOG_CONNECTION_OPTIONS", ""),
			"must_connected":      true,
			"is_connect_at_start": true,
			"create_script_files": []string{createScriptFileFolder + "/db_auditlog.sql"},
		},
	}, []string{
		"system.user_name", "system.user_password",
		"task-dispatcher.user_name", "task-dispatcher.user_password",
		"auditlog.user_name", "auditlog.user_password",
	})

	// Transient Session Manager
	d = configuration.Manager.NewIfNotExistConfiguration("redis", "redis.json", "json", false, false, map[string]any{
		"session": map[string]any{
			"nameid":              "session",
			"address":             app.App.InitVault.GetStringOrDefault("REDIS_SESSION_ADDRESS", ""),
			"database_index":      app.App.InitVault.GetIntOrDefault("REDIS_SESSION_DB_INDEX", 0),
			"user_name":           app.App.InitVault.GetStringOrDefault("REDIS_SESSION_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("REDIS_SESSION_USER_PASSWORD", ""),
			"must_connected":      true,
			"is_connect_at_start": true,
		},
		"prekey": map[string]any{
			"nameid":              "prekey",
			"address":             app.App.InitVault.GetStringOrDefault("REDIS_PREKEY_ADDRESS", ""),
			"database_index":      app.App.InitVault.GetIntOrDefault("REDIS_PREKEY_DB_INDEX", 0),
			"user_name":           app.App.InitVault.GetStringOrDefault("REDIS_PREKEY_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("REDIS_PREKEY_USER_PASSWORD", ""),
			"must_connected":      true,
			"is_connect_at_start": true,
		},
		"rate_limit": map[string]any{
			"nameid":              "rate_limit",
			"address":             app.App.InitVault.GetStringOrDefault("REDIS_RATE_LIMIT_ADDRESS", ""),
			"database_index":      app.App.InitVault.GetIntOrDefault("REDIS_RATE_LIMIT_INDEX", 1),
			"user_name":           app.App.InitVault.GetStringOrDefault("REDIS_RATE_LIMIT_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("REDIS_RATE_LIMIT_USER_PASSWORD", ""),
			"must_connected":      true,
			"is_connect_at_start": true,
		},
	}, []string{"session.user_name", "session.password", "prekey.user_name", "prekey.password", "rate_limit.user_name", "rate_limit.password"})

	base.RelyOnInboundRedisSessionName = "session"

	configuration.Manager.NewIfNotExistConfiguration("object_storage", "object_storage.json", "json", false, false, map[string]any{
		/*"test_bucket": map[string]any{
			"nameid":              "test_bucket",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_TEST_BUCKET_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_TEST_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_TEST_BUCKET_USER_NAME",""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_TEST_BUCKET_USER_PASSWORD",""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_TEST_BUCKET_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_TEST_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
		},*/
		"user-avatar-source": map[string]any{
			"nameid":              "user-avatar-source",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_SOURCE_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_SOURCE_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_SOURCE_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_SOURCE_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_SOURCE_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_AVATAR_SOURCE_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
		},
		"user-avatar-small": map[string]any{
			"nameid":              "user-avatar-small",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_SMALL_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_SMALL_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_SMALL_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_SMALL_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_SMALL_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_AVATAR_SMALL_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
			"file_image_width":    app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_AVATAR_SMALL_FILE_IMAGE_WIDTH_IN_PIXEL", 128),
			"file_image_height":   app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_AVATAR_SMALL_FILE_IMAGE_HEIGHT_IN_PIXEL", 128),
		},
		"user-avatar-medium": map[string]any{
			"nameid":              "user-avatar-medium",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_MEDIUM_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_MEDIUM_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_MEDIUM_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_MEDIUM_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_AVATAR_MEDIUM_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_AVATAR_MEDIUM_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
			"file_image_width":    app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_AVATAR_MEDIUM_FILE_IMAGE_WIDTH_IN_PIXEL", 256),
			"file_image_height":   app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_AVATAR_MEDIUM_FILE_IMAGE_HEIGHT_IN_PIXEL", 256),
		},

		"announcement-picture-source": map[string]any{
			"nameid":              "announcement-picture-source",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_SOURCE_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_SOURCE_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_SOURCE_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_SOURCE_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_SOURCE_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_SOURCE_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
		},
		"announcement-picture-small": map[string]any{
			"nameid":              "announcement-picture-small",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_SMALL_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_SMALL_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_SMALL_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_SMALL_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_SMALL_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_SMALL_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
			"file_image_width":    app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_PICTURE_SMALL_FILE_IMAGE_WIDTH_IN_PIXEL", 256),
			"file_image_height":   app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_PICTURE_SMALL_FILE_IMAGE_HEIGHT_IN_PIXEL", 256),
		},
		"announcement-picture-medium": map[string]any{
			"nameid":              "announcement-picture-medium",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_MEDIUM_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_MEDIUM_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_MEDIUM_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_MEDIUM_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_MEDIUM_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_MEDIUM_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
			"file_image_width":    app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_PICTURE_MEDIUM_FILE_IMAGE_WIDTH_IN_PIXEL", 512),
			"file_image_height":   app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_PICTURE_MEDIUM_FILE_IMAGE_HEIGHT_IN_PIXEL", 512),
		},
		"announcement-picture-big": map[string]any{
			"nameid":              "announcement-picture-big",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_BIG_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_BIG_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_BIG_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_BIG_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_BIG_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_BIG_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
			"file_image_width":    app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_PICTURE_BIG_FILE_IMAGE_WIDTH_IN_PIXEL", 1024),
			"file_image_height":   app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_ANNOUNCEMENT_PICTURE_BIG_FILE_IMAGE_HEIGHT_IN_PIXEL", 1024),
		},
		"sub-task-report-picture-source": map[string]any{
			"nameid":              "sub-task-report-picture-source",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SOURCE_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SOURCE_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SOURCE_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SOURCE_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SOURCE_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SOURCE_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
		},
		"sub-task-report-picture-small": map[string]any{
			"nameid":              "sub-task-report-picture-small",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SMALL_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SMALL_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SMALL_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SMALL_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SMALL_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SMALL_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
			"file_image_width":    app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SMALL_FILE_IMAGE_WIDTH_IN_PIXEL", 256),
			"file_image_height":   app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_SMALL_FILE_IMAGE_HEIGHT_IN_PIXEL", 256),
		},
		"sub-task-report-picture-medium": map[string]any{
			"nameid":              "sub-task-report-picture-medium",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_MEDIUM_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_MEDIUM_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_MEDIUM_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_MEDIUM_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_MEDIUM_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_MEDIUM_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
			"file_image_width":    app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_MEDIUM_FILE_IMAGE_WIDTH_IN_PIXEL", 256),
			"file_image_height":   app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_MEDIUM_FILE_IMAGE_HEIGHT_IN_PIXEL", 256),
		},
		"sub-task-report-picture-big": map[string]any{
			"nameid":              "sub-task-report-picture-big",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_BIG_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_BIG_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_BIG_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_BIG_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_BIG_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_BIG_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
			"file_image_width":    app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_BIG_FILE_IMAGE_WIDTH_IN_PIXEL", 256),
			"file_image_height":   app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_SUB_TASK_REPORT_BIG_FILE_IMAGE_HEIGHT_IN_PIXEL", 256),
		},
		"user-identity-card-source": map[string]any{
			"nameid":              "user-identity-card-source",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_SOURCE_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_SOURCE_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_SOURCE_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_SOURCE_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_SOURCE_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_SOURCE_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
		},
		"user-identity-card-big": map[string]any{
			"nameid":              "user-identity-card-big",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_BIG_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_BIG_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_BIG_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_BIG_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_BIG_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_BIG_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
			"file_image_width":    app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_BIG_FILE_IMAGE_WIDTH_IN_PIXEL", 2048),
			"file_image_height":   app.App.InitVault.GetIntOrDefault("OBJECT_STORAGE_USER_IDENTITY_CARD_BIG_FILE_IMAGE_HEIGHT_IN_PIXEL", 1024),
		},
		"berita-acara": map[string]any{
			"nameid":              "berita-acara",
			"address":             app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_BERITA_ACARA_ADDRESS", ""),
			"bucket_name":         app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_BERITA_ACARA_BUCKET_NAME", ""),
			"user_name":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_BERITA_ACARA_USER_NAME", ""),
			"password":            app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_BERITA_ACARA_USER_PASSWORD", ""),
			"base_path":           app.App.InitVault.GetStringOrDefault("OBJECT_STORAGE_BERITA_ACARA_BASE_PATH", ""),
			"use_ssl":             app.App.InitVault.GetBoolOrDefault("OBJECT_STORAGE_BERITA_ACARA_USE_SSL", false),
			"must_connected":      true,
			"is_connect_at_start": true,
		},
	}, []string{
		//"test_bucket.user_name", "test_bucket.password",
		"user-avatar-source.user_name", "user-avatar-source.password",
		"user-avatar-small.user_name", "user-avatar-small.password",
		"user-avatar-medium.user_name", "user-avatar-medium.password",
		"user-avatar-big.user_name", "user-avatar-big.password",
		"announcement-picture-source.user_name", "announcement-picture-source.password",
		"announcement-picture-small.user_name", "announcement-picture-small.password",
		"announcement-picture-medium.user_name", "announcement-picture-medium.password",
		"announcement-picture-big.user_name", "announcement-picture-big.password",
		"sub-task-report-picture-source.user_name", "sub-task-report-picture-source.password",
		"sub-task-report-picture-small.user_name", "sub-task-report-picture-small.password",
		"sub-task-report-picture-medium.user_name", "sub-task-report-picture-medium.password",
		"sub-task-report-picture-big.user_name", "sub-task-report-picture-big.password",
		"user-identity-card-source.user_name", "user-identity-card-source.password",
		"user-identity-card-big.user_name", "user-identity-card-big.password",
	})
	configuration.Manager.NewIfNotExistConfiguration("api_endpoint_rate_limiter", "api_endpoint_rate_limiter.json", "json", false, false, map[string]any{
		"default": map[string]any{
			"nameid":                    "default",
			"max_attempts":              app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_DEFAULT_MAX_ATTEMPTS", 500),         // Default 100 requests
			"time_window_in_minutes":    app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_DEFAULT_TIME_WINDOW_IN_MINUTES", 1), // Per minute
			"block_duration_in_minutes": app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_DEFAULT_BLOCK_DURATION_IN_MINUTES", 5),
		},
		"/api-webadmin/login": utils.JSON{
			"nameid":                    "/api-webadmin/login",
			"max_attempts":              app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_API_WEBADMIN_LOGIN_MAX_ATTEMPTS", 200),         // Default 100 requests
			"time_window_in_minutes":    app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_API_WEBADMIN_LOGIN_TIME_WINDOW_IN_MINUTES", 1), // Per minute
			"block_duration_in_minutes": app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_API_WEBADMIN_LOGIN_BLOCK_DURATION_IN_MINUTES", 5),
		},
		"/api-mobile/login": utils.JSON{
			"nameid":                    "/api-mobile/login",
			"max_attempts":              app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_API_MOBILE_LOGIN_MAX_ATTEMPTS", 200),         // Default 100 requests
			"time_window_in_minutes":    app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_API_MOBILE_LOGIN_TIME_WINDOW_IN_MINUTES", 1), // Per minute
			"block_duration_in_minutes": app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_API_MOBILE_LOGIN_BLOCK_DURATION_IN_MINUTES", 5),
		},
		"/api-mobile/field_executor/picture/upload": utils.JSON{
			"nameid":                    "/mobile/field_executor/picture/upload",
			"max_attempts":              app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_API_MOBILE_FIELD_EXECUTOR_PICTURE_UPLOAD_MAX_ATTEMPTS", 200),         // Default 100 requests
			"time_window_in_minutes":    app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_API_MOBILE_FIELD_EXECUTOR_PICTURE_UPLOAD_TIME_WINDOW_IN_MINUTES", 1), // Per minute
			"block_duration_in_minutes": app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_API_MOBILE_FIELD_EXECUTOR_PICTURE_UPLOAD_BLOCK_DURATION_IN_MINUTES", 5),
		},
		"/api-mobile/*/picture/download": utils.JSON{
			"nameid":                    "/mobile/field_executor/picture/download",
			"max_attempts":              app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_API_MOBILE_FIELD_EXECUTOR_PICTURE_DOWNLOAD_MAX_ATTEMPTS", 200),         // Default 100 requests
			"time_window_in_minutes":    app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_API_MOBILE_FIELD_EXECUTOR_PICTURE_DOWNLOAD_TIME_WINDOW_IN_MINUTES", 1), // Per minute
			"block_duration_in_minutes": app.App.InitVault.GetIntOrDefault("API_ENDPOINT_RATE_LIMITER_API_MOBILE_FIELD_EXECUTOR_PICTURE_DOWNLOAD_BLOCK_DURATION_IN_MINUTES", 5),
		},
	}, nil)

}
