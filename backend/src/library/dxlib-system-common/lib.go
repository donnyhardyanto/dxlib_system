package common

import (
	"github.com/donnyhardyanto/dxlib/app"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/module/audit_log"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/construction_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/master_data"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
)

type Partner struct {
	DatabaseNameIdAuditLog       string
	DatabaseNameIdTaskDispatcher string
	DatabaseNameIdConfig         string
	MasterData                   *master_data.MasterData
	PartnerManagement            *partner_management.PartnerManagement
	TaskManagement               *task_management.TaskManagement
	UserManagement               *user_management.DxmUserManagement
	ConstructionManagement       *construction_management.ConstructionManagement
}

func (p *Partner) RoleCreate(l *log.DXLog, pr utils.JSON) (roleId int64, err error) {
	for k, v := range pr {
		if v == nil {
			delete(pr, k)
		}
	}
	db := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]
	err = db.Tx(l, database.LevelReadCommitted, func(tx *database.DXDatabaseTx) (err2 error) {
		roleId, err2 = p.RoleTxCreate(tx, pr)
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return roleId, nil
}

func (p *Partner) RoleTxCreate(tx *database.DXDatabaseTx, pr utils.JSON) (roleId int64, err error) {
	isAreaCodeExist, areaCode, err := utils.ExtractMapValue[string](&pr, "area_code")
	if err != nil {
		return 0, err
	}

	isTaskTypeIdExist, taskTypeId, err := utils.ExtractMapValue[int64](&pr, "task_type_id")
	if err != nil {
		return 0, err
	}

	roleId, err = partner_management.ModulePartnerManagement.Role.TxInsert(tx, pr)
	if err != nil {
		return 0, err
	}
	if isAreaCodeExist {
		_, err = partner_management.ModulePartnerManagement.RoleArea.TxInsert(tx, map[string]any{
			"role_id":   roleId,
			"area_code": areaCode,
		})
		if err != nil {
			return roleId, err
		}
	}
	if isTaskTypeIdExist {
		_, err = partner_management.ModulePartnerManagement.RoleTaskType.TxInsert(tx, map[string]any{
			"role_id":      roleId,
			"task_type_id": taskTypeId,
		})
		if err != nil {
			return roleId, err
		}
	}

	return roleId, nil
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
		MasterData:                   &master_data.ModuleMasterData,
		PartnerManagement:            &partner_management.ModulePartnerManagement,
		TaskManagement:               &task_management.ModuleTaskManagement,
		UserManagement:               &user_management.ModuleUserManagement,
		ConstructionManagement:       &construction_management.ModuleConstructionManagement,
	}
}
