package handler

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	pgn_partner_common "common"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"net/http"
)

func RoleCreate(aepr *api.DXAPIEndPointRequest) (err error) {
	p := aepr.GetParameterValues()

	roleId, err := pgn_partner_common.PartnerInstance.RoleCreate(&aepr.Log, p)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		pgn_partner_common.PartnerInstance.PartnerManagement.Role.FieldNameForRowId: roleId,
	}})
	return nil
}

func RoleEdit(aepr *api.DXAPIEndPointRequest) (err error) {
	_, roleId, err := aepr.GetParameterValueAsInt64("id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, p, err := aepr.GetParameterValueAsJSON("new")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	if p == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "EMPTY_DATA")
	}

	areaCode, isAreaCodeExist := p["area_code"].(string)
	if isAreaCodeExist {
		delete(p, "area_code")
	}

	taskTypeId, isTaskTypeIdExist := p["task_type_id"].(int64)
	if isTaskTypeIdExist {
		delete(p, "task_type_id")
	}

	t := partner_management.ModulePartnerManagement.Role

	db := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]
	err = db.Tx(&aepr.Log, database.LevelReadCommitted, func(tx *database.DXDatabaseTx) (err2 error) {
		if len(p) > 0 {
			_, err2 = t.TxUpdate(tx, p, utils.JSON{
				t.FieldNameForRowId: roleId,
			})
			if err2 != nil {
				return err2
			}
		}
		if isAreaCodeExist {
			_, err2 = partner_management.ModulePartnerManagement.RoleArea.TxUpdate(tx, utils.JSON{
				"area_code": areaCode,
			}, utils.JSON{
				"role_id": roleId,
			})
			if err2 != nil {
				return err2
			}
		}
		if isTaskTypeIdExist {
			_, err2 = partner_management.ModulePartnerManagement.RoleTaskType.TxUpdate(tx, utils.JSON{
				"task_type_id": taskTypeId,
			}, utils.JSON{
				"role_id": roleId,
			})
			if err2 != nil {
				return err2
			}
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
		"data": utils.JSON{
			t.FieldNameForRowId: roleId,
		},
	})

	return nil
}

func RoleDelete(aepr *api.DXAPIEndPointRequest) (err error) {
	_, roleId, err := aepr.GetParameterValueAsInt64("id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	t := partner_management.ModulePartnerManagement.Role

	db := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]
	err = db.Tx(&aepr.Log, database.LevelReadCommitted, func(tx *database.DXDatabaseTx) (err2 error) {
		_, err2 = t.TxHardDelete(tx, utils.JSON{
			t.FieldNameForRowId: roleId,
		})
		if err2 != nil {
			return err2
		}
		_, err2 = partner_management.ModulePartnerManagement.RoleArea.TxHardDelete(tx, utils.JSON{
			"role_id": roleId,
		})
		if err2 != nil {
			return err2
		}
		_, err2 = partner_management.ModulePartnerManagement.RoleTaskType.TxHardDelete(tx, utils.JSON{
			"role_id": roleId,
		})
		if err2 != nil {
			return err2
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
		"data": utils.JSON{
			t.FieldNameForRowId: roleId,
		},
	})
	return nil
}
