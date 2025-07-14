package infrastructure

import (
	"encoding/json"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/module/audit_log"
)

func DoOnAuditLogStart(oldAuditLogId int64, parameters *api.DXAPIAuditLogEntry) (newAuditLogId int64, err error) {
	log.Log.Infof("Audit Log Start")

	jsonData, err := json.Marshal(parameters)
	if err != nil {
		log.Log.Errorf(err, "DO_ON_AUDIT_LOG_START:MARSHAL_ERROR: %v", err)
		return
	}
	mapData := map[string]any{}
	_ = json.Unmarshal(jsonData, &mapData)

	newAuditLogId, err = audit_log.ModuleAuditLog.UserActivityLog.Insert(&log.Log, mapData)
	if err != nil {
		log.Log.Errorf(err, "DO_ON_AUDIT_LOG_START:INSERT_ERROR: %v", err)
		return
	}
	return newAuditLogId, nil
}

func DoOnAuditLogUserIdentified(oldAuditLogId int64, parameters *api.DXAPIAuditLogEntry) (newAuditLogId int64, err error) {
	log.Log.Infof("Audit Log User Identified")
	jsonData, err := json.Marshal(parameters)
	if err != nil {
		log.Log.Errorf(err, "DO_ON_AUDIT_LOG_USER_IDENTIFIED:MARSHAL_ERROR: %v", err)
		return
	}
	mapData := map[string]any{}
	json.Unmarshal(jsonData, &mapData)

	_, err = audit_log.ModuleAuditLog.UserActivityLog.Update(mapData, utils.JSON{
		"id": oldAuditLogId,
	})
	if err != nil {
		log.Log.Errorf(err, "DO_ON_AUDIT_LOG_USER_IDENTIFIED:UPDATE_ERROR: %v", err)
		return
	}
	return oldAuditLogId, nil
}

func DoOnAuditLogEnd(oldAuditLogId int64, parameters *api.DXAPIAuditLogEntry) (newAuditLogId int64, err error) {
	log.Log.Infof("Audit Log End")
	jsonData, err := json.Marshal(parameters)
	if err != nil {
		log.Log.Errorf(err, "DO_ON_AUDIT_LOG_START:MARSHAL_ERROR: %v", err)
		return
	}
	mapData := map[string]any{}
	json.Unmarshal(jsonData, &mapData)

	_, err = audit_log.ModuleAuditLog.UserActivityLog.Update(mapData, utils.JSON{
		"id": oldAuditLogId,
	})
	if err != nil {
		log.Log.Errorf(err, "DO_ON_AUDIT_LOG_START:UPDATE_ERROR: %v", err)
		return
	}
	return oldAuditLogId, nil
}
