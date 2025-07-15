package common

import (
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib/app"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib_module/module/audit_log"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
)

type Partner struct {
	DatabaseNameIdAuditLog       string
	DatabaseNameIdTaskDispatcher string
	DatabaseNameIdConfig         string
	UserManagement               *user_management.DxmUserManagement
}

var PartnerInstance Partner

func init() {
	app.App.OnStartStorageReady = func() (err error) {
		log.OnError = audit_log.ModuleAuditLog.DoError
		return nil
	}

	PartnerInstance = Partner{
		DatabaseNameIdAuditLog:       base.DatabaseNameIdAuditLog,
		DatabaseNameIdTaskDispatcher: base.DatabaseNameIdTaskDispatcher,
		DatabaseNameIdConfig:         base.DatabaseNameIdConfig,
		UserManagement:               &user_management.ModuleUserManagement,
	}
}
