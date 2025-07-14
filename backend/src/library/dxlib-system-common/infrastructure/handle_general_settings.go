package infrastructure

import (
	"database/sql"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/module/general"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"net/http"
)

func GeneralSettingsTaskConstructionUpdate(aepr *api.DXAPIEndPointRequest) (err error) {
	_, value, err := aepr.GetParameterValueAsJSON("new")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	// Begin transaction
	dbTaskDispatcher := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]

	err = dbTaskDispatcher.Tx(&aepr.Log, sql.LevelReadCommitted, func(tx *database.DXDatabaseTx) (err error) {
		err = general.ModuleGeneral.Property.TxSetAsJSON(tx, base.SettingsTaskConstruction, value)

		if err != nil {
			aepr.Log.Warnf("GENERAL_SETTINGS_NOT_FOUND:%s", err.Error())
			err = aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "GENERAL_SETTINGS_NOT_FOUND")
			return errors.Wrap(err, "error occured")
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
	return nil
}

func GeneralSettingsTaskConstructionRead(aepr *api.DXAPIEndPointRequest) (err error) {
	value, err := general.ModuleGeneral.Property.GetAsJSON(&aepr.Log, base.SettingsTaskConstruction)
	if err != nil {
		aepr.Log.Warnf("GENERAL_SETTINGS_NOT_FOUND:%s", err.Error())
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "GENERAL_SETTINGS_NOT_FOUND")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"settings": value,
	}})
	return nil
}
