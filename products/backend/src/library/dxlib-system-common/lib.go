package common

import (
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib/app"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib_module/module/audit_log"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
)

type System struct {
	DatabaseNameIdAuditLog       string
	DatabaseNameIdTaskDispatcher string
	DatabaseNameIdConfig         string
	UserManagement               *user_management.DxmUserManagement
}

var SystemInstance System

func init() {
	app.App.OnStartStorageReady = func() (err error) {
		log.OnError = audit_log.ModuleAuditLog.DoError
		return nil
	}

	SystemInstance = System{
		DatabaseNameIdAuditLog:       base.DatabaseNameIdAuditLog,
		DatabaseNameIdTaskDispatcher: base.DatabaseNameIdDbBase,
		DatabaseNameIdConfig:         base.DatabaseNameIdConfig,
		UserManagement:               &user_management.ModuleUserManagement,
	}
}
