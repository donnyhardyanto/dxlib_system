package table

import (
	"database/sql"
	"fmt"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/database/protected/db"
	"github.com/donnyhardyanto/dxlib/database/protected/export"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	utilsJson "github.com/donnyhardyanto/dxlib/utils/json"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
	_ "time/tzdata"
)

type DXTable struct {
	DatabaseNameId             string
	Database                   *database.DXDatabase
	NameId                     string
	ResultObjectName           string
	ListViewNameId             string
	FieldNameForRowId          string
	FieldNameForRowNameId      string
	FieldNameForRowUid         string
	FieldTypeMapping           db.FieldTypeMapping
	ResponseEnvelopeObjectName string
}

func (t *DXTable) DoInsert(aepr *api.DXAPIEndPointRequest, newKeyValues utils.JSON) (newId int64, err error) {
	newKeyValues["is_deleted"] = false

	tt := time.Now().UTC()
	newKeyValues["created_at"] = tt
	_, ok := newKeyValues["created_by_user_id"]
	if !ok {
		if aepr.CurrentUser.Id != "" {
			newKeyValues["created_by_user_id"] = aepr.CurrentUser.Id
			newKeyValues["created_by_user_nameid"] = aepr.CurrentUser.LoginId
		} else {
			newKeyValues["created_by_user_id"] = "0"
			newKeyValues["created_by_user_nameid"] = "SYSTEM"
		}
		newKeyValues["last_modified_at"] = tt
		if aepr.CurrentUser.Id != "" {
			newKeyValues["last_modified_by_user_id"] = aepr.CurrentUser.Id
			newKeyValues["last_modified_by_user_nameid"] = aepr.CurrentUser.LoginId
		} else {
			newKeyValues["last_modified_by_user_id"] = "0"
			newKeyValues["last_modified_by_user_nameid"] = "SYSTEM"
		}
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	newId, err = t.Database.Insert(t.NameId, t.FieldNameForRowId, newKeyValues)
	if err != nil {
		return 0, err
	}

	p := utils.JSON{
		t.FieldNameForRowId: newId,
	}

	if t.FieldNameForRowUid != "" {
		_, n, err := t.Database.SelectOne(t.ListViewNameId, nil, nil, utils.JSON{
			"id": newId,
		}, nil, nil)
		if err != nil {
			return 0, err
		}
		uid, ok := n[t.FieldNameForRowUid].(string)
		if !ok {
			return 0, errors.New("IMPOSSIBLE:UID")
		}
		p[t.FieldNameForRowUid] = uid
	}

	data := utilsJson.Encapsulate(t.ResponseEnvelopeObjectName, p)
	aepr.WriteResponseAsJSON(http.StatusOK, nil, data)

	return newId, nil
}

func (t *DXTable) DoCreate(aepr *api.DXAPIEndPointRequest, newKeyValues utils.JSON) (newId int64, err error) {
	newId, err = t.DoInsert(aepr, newKeyValues)
	if err != nil {
		aepr.WriteResponseAsError(http.StatusConflict, err)
		return 0, nil
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utilsJson.Encapsulate(t.ResponseEnvelopeObjectName, utils.JSON{
		t.FieldNameForRowId: newId,
	},
	))

	return newId, nil
}

func (t *DXTable) GetById(log *log.DXLog, id int64) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {
	rowsInfo, r, err = t.SelectOne(log, nil, utils.JSON{
		t.FieldNameForRowId: id,
		"is_deleted":        false,
	}, nil, map[string]string{t.FieldNameForRowId: "asc"})
	return rowsInfo, r, err
}

func (t *DXTable) ShouldGetById(log *log.DXLog, id int64) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {
	rowsInfo, r, err = t.ShouldSelectOne(log, utils.JSON{
		t.FieldNameForRowId: id,
		"is_deleted":        false,
	}, nil, map[string]string{t.FieldNameForRowId: "asc"})
	return rowsInfo, r, err
}

func (t *DXTable) ShouldGetByUid(log *log.DXLog, uid string) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {
	rowsInfo, r, err = t.ShouldSelectOne(log, utils.JSON{
		t.FieldNameForRowUid: uid,
		"is_deleted":         false,
	}, nil, map[string]string{t.FieldNameForRowId: "asc"})
	return rowsInfo, r, err
}
func (t *DXTable) ShouldGetByUtag(log *log.DXLog, utag string) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {
	rowsInfo, r, err = t.ShouldSelectOne(log, utils.JSON{
		"utag":       utag,
		"is_deleted": false,
	}, nil, map[string]string{t.FieldNameForRowId: "asc"})
	return rowsInfo, r, err
}

func (t *DXTable) GetByNameId(log *log.DXLog, nameid string) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {
	rowsInfo, r, err = t.SelectOne(log, nil, utils.JSON{
		t.FieldNameForRowNameId: nameid,
		"is_deleted":            false,
	}, nil, map[string]string{t.FieldNameForRowNameId: "asc"})
	return rowsInfo, r, err
}

func (t *DXTable) ShouldGetByNameId(log *log.DXLog, nameid string) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {
	rowsInfo, r, err = t.ShouldSelectOne(log, utils.JSON{
		t.FieldNameForRowNameId: nameid,
		"is_deleted":            false,
	}, nil, map[string]string{t.FieldNameForRowNameId: "asc"})
	return rowsInfo, r, err
}

func (t *DXTable) TxShouldGetById(tx *database.DXDatabaseTx, id int64) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {
	rowsInfo, r, err = tx.ShouldSelectOne(t.ListViewNameId, t.FieldTypeMapping, nil, utils.JSON{
		t.FieldNameForRowId: id,
		"is_deleted":        false,
	}, nil, nil, nil)
	return rowsInfo, r, err
}

func (t *DXTable) TxGetByNameId(tx *database.DXDatabaseTx, nameId string) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {
	rowsInfo, r, err = tx.SelectOne(t.ListViewNameId, t.FieldTypeMapping, nil, utils.JSON{
		t.FieldNameForRowNameId: nameId,
		"is_deleted":            false,
	}, nil, nil, nil)
	return rowsInfo, r, err
}

func (t *DXTable) TxShouldGetByNameId(tx *database.DXDatabaseTx, nameId string) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {
	rowsInfo, r, err = tx.ShouldSelectOne(t.ListViewNameId, t.FieldTypeMapping, nil, utils.JSON{
		t.FieldNameForRowNameId: nameId,
		"is_deleted":            false,
	}, nil, nil, nil)
	return rowsInfo, r, err
}

func (t *DXTable) TxInsert(tx *database.DXDatabaseTx, newKeyValues utils.JSON) (newId int64, err error) {
	//n := utils.NowAsString()
	tt := time.Now().UTC()
	newKeyValues["is_deleted"] = false
	//newKeyValues["created_at"] = n
	newKeyValues["created_at"] = tt

	_, ok := newKeyValues["created_by_user_id"]
	if !ok {
		newKeyValues["created_by_user_id"] = "0"
		newKeyValues["created_by_user_nameid"] = "SYSTEM"
		newKeyValues["last_modified_by_user_id"] = "0"
		newKeyValues["last_modified_by_user_nameid"] = "SYSTEM"
	}

	newId, err = tx.Insert(t.NameId, newKeyValues)
	return newId, err
}

func (t *DXTable) InRequestTxInsert(aepr *api.DXAPIEndPointRequest, tx *database.DXDatabaseTx, newKeyValues utils.JSON) (newId int64, err error) {
	n := utils.NowAsString()
	newKeyValues["is_deleted"] = false
	newKeyValues["created_at"] = n
	_, ok := newKeyValues["created_by_user_id"]
	if !ok {
		if aepr.CurrentUser.Id != "" {
			newKeyValues["created_by_user_id"] = aepr.CurrentUser.Id
			newKeyValues["created_by_user_nameid"] = aepr.CurrentUser.LoginId
		} else {
			newKeyValues["created_by_user_id"] = "0"
			newKeyValues["created_by_user_nameid"] = "SYSTEM"
		}
		newKeyValues["last_modified_at"] = n
		if aepr.CurrentUser.Id != "" {
			newKeyValues["last_modified_by_user_id"] = aepr.CurrentUser.Id
			newKeyValues["last_modified_by_user_nameid"] = aepr.CurrentUser.LoginId
		} else {
			newKeyValues["last_modified_by_user_id"] = "0"
			newKeyValues["last_modified_by_user_nameid"] = "SYSTEM"
		}
	}

	newId, err = tx.Insert(t.NameId, newKeyValues)
	return newId, err
}

func (t *DXTable) Insert(log *log.DXLog, newKeyValues utils.JSON) (newId int64, err error) {
	tt := time.Now().UTC()
	newKeyValues["created_at"] = tt
	newKeyValues["last_modified_at"] = tt
	newKeyValues["is_deleted"] = false
	_, ok := newKeyValues["created_by_user_id"]
	if !ok {
		newKeyValues["created_by_user_id"] = "0"
		newKeyValues["created_by_user_nameid"] = "SYSTEM"
		newKeyValues["last_modified_by_user_id"] = "0"
		newKeyValues["last_modified_by_user_nameid"] = "SYSTEM"
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	newId, err = t.Database.Insert(t.NameId, t.FieldNameForRowId, newKeyValues)
	return newId, err
}

func (t *DXTable) Upsert(setKeyValues utils.JSON, whereAndFieldNameValues utils.JSON) (result sql.Result, newId int64, err error) {
	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	_, r, err := t.Database.SelectOne(t.NameId, nil, nil, whereAndFieldNameValues, nil, nil)
	if err != nil {
		return nil, 0, err
	}
	if r == nil {
		newSetKeyValues := utilsJson.DeepMerge2(setKeyValues, whereAndFieldNameValues)

		newId, err = t.Database.Insert(t.NameId, t.FieldNameForRowId, newSetKeyValues)
		return nil, newId, err
	} else {
		result, err = t.Database.Update(t.NameId, setKeyValues, whereAndFieldNameValues)
		return result, 0, err
	}
}

func (t *DXTable) Update(setKeyValues utils.JSON, whereAndFieldNameValues utils.JSON) (result sql.Result, err error) {
	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}
	whereAndFieldNameValues["is_deleted"] = false

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	return t.Database.Update(t.NameId, setKeyValues, whereAndFieldNameValues)
}

func (t *DXTable) UpdateOne(l *log.DXLog, FieldValueForId int64, setKeyValues utils.JSON) (result sql.Result, err error) {
	_, _, err = t.ShouldGetById(l, FieldValueForId)
	if err != nil {
		return nil, err
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	return t.Database.Update(t.NameId, setKeyValues, utils.JSON{
		t.FieldNameForRowId: FieldValueForId,
	})
}

func (t *DXTable) InRequestInsert(aepr *api.DXAPIEndPointRequest, newKeyValues utils.JSON) (newId int64, err error) {
	n := utils.NowAsString()
	newKeyValues["is_deleted"] = false
	newKeyValues["created_at"] = n
	_, ok := newKeyValues["created_by_user_id"]
	if !ok {
		if aepr.CurrentUser.Id != "" {
			newKeyValues["created_by_user_id"] = aepr.CurrentUser.Id
			newKeyValues["created_by_user_nameid"] = aepr.CurrentUser.LoginId
		} else {
			newKeyValues["created_by_user_id"] = "0"
			newKeyValues["created_by_user_nameid"] = "SYSTEM"
		}
		newKeyValues["last_modified_at"] = n
		if aepr.CurrentUser.Id != "" {
			newKeyValues["last_modified_by_user_id"] = aepr.CurrentUser.Id
			newKeyValues["last_modified_by_user_nameid"] = aepr.CurrentUser.LoginId
		} else {
			newKeyValues["last_modified_by_user_id"] = "0"
			newKeyValues["last_modified_by_user_nameid"] = "SYSTEM"
		}
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	newId, err = t.Database.Insert(t.NameId, t.FieldNameForRowId, newKeyValues)
	return newId, err
}

func (t *DXTable) RequestRead(aepr *api.DXAPIEndPointRequest) (err error) {
	_, id, err := aepr.GetParameterValueAsInt64(t.FieldNameForRowId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	rowsInfo, d, err := t.ShouldGetById(&aepr.Log, id)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utilsJson.Encapsulate(t.ResponseEnvelopeObjectName, utils.JSON{t.ResultObjectName: d, "rows_info": rowsInfo}))

	return nil
}

func (t *DXTable) RequestReadByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	_, uid, err := aepr.GetParameterValueAsString(t.FieldNameForRowUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	rowsInfo, d, err := t.ShouldGetByUid(&aepr.Log, uid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utilsJson.Encapsulate(t.ResponseEnvelopeObjectName, utils.JSON{t.ResultObjectName: d, "rows_info": rowsInfo}))

	return nil
}

func (t *DXTable) RequestReadByNameId(aepr *api.DXAPIEndPointRequest) (err error) {
	_, nameid, err := aepr.GetParameterValueAsString(t.FieldNameForRowNameId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	rowsInfo, d, err := t.ShouldGetByNameId(&aepr.Log, nameid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utilsJson.Encapsulate(t.ResponseEnvelopeObjectName, utils.JSON{t.ResultObjectName: d, "rows_info": rowsInfo}))

	return nil
}

func (t *DXTable) RequestReadByUtag(aepr *api.DXAPIEndPointRequest) (err error) {
	_, utag, err := aepr.GetParameterValueAsString("utag")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	rowsInfo, d, err := t.ShouldGetByUtag(&aepr.Log, utag)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utilsJson.Encapsulate(t.ResponseEnvelopeObjectName, utils.JSON{t.ResultObjectName: d, "rows_info": rowsInfo}))

	return nil
}

func (t *DXTable) DoEdit(aepr *api.DXAPIEndPointRequest, id int64, newKeyValues utils.JSON) (err error) {
	_, _, err = t.ShouldGetById(&aepr.Log, id)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	tt := time.Now().UTC()
	newKeyValues["last_modified_at"] = tt

	_, ok := newKeyValues["last_modified_by_user_id"]
	if !ok {
		if aepr.CurrentUser.Id != "" {
			newKeyValues["last_modified_by_user_id"] = aepr.CurrentUser.Id
			newKeyValues["last_modified_by_user_nameid"] = aepr.CurrentUser.LoginId
		} else {
			newKeyValues["last_modified_by_user_id"] = "0"
			newKeyValues["last_modified_by_user_nameid"] = "SYSTEM"
		}
	}

	for k, v := range newKeyValues {
		if v == nil {
			delete(newKeyValues, k)
		}
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	_, err = db.Update(t.Database.Connection, t.NameId, newKeyValues, utils.JSON{
		t.FieldNameForRowId: id,
		"is_deleted":        false,
	})
	if err != nil {
		aepr.Log.Errorf(err, "Error at %s.DoEdit (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utilsJson.Encapsulate(t.ResponseEnvelopeObjectName, utils.JSON{
		t.FieldNameForRowId: id,
	},
	))
	return nil
}

func (t *DXTable) DoEditByUid(aepr *api.DXAPIEndPointRequest, uid string, newKeyValues utils.JSON) (err error) {
	_, _, err = t.ShouldGetByUid(&aepr.Log, uid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	tt := time.Now().UTC()
	newKeyValues["last_modified_at"] = tt

	_, ok := newKeyValues["last_modified_by_user_id"]
	if !ok {
		if aepr.CurrentUser.Id != "" {
			newKeyValues["last_modified_by_user_id"] = aepr.CurrentUser.Id
			newKeyValues["last_modified_by_user_nameid"] = aepr.CurrentUser.LoginId
		} else {
			newKeyValues["last_modified_by_user_id"] = "0"
			newKeyValues["last_modified_by_user_nameid"] = "SYSTEM"
		}
	}

	for k, v := range newKeyValues {
		if v == nil {
			delete(newKeyValues, k)
		}
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	_, err = db.Update(t.Database.Connection, t.NameId, newKeyValues, utils.JSON{
		t.FieldNameForRowUid: uid,
		"is_deleted":         false,
	})
	if err != nil {
		aepr.Log.Errorf(err, "Error at %s.DoEdit (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utilsJson.Encapsulate(t.ResponseEnvelopeObjectName, utils.JSON{
		t.FieldNameForRowUid: uid,
	},
	))
	return nil
}

func (t *DXTable) RequestEdit(aepr *api.DXAPIEndPointRequest) (err error) {
	_, id, err := aepr.GetParameterValueAsInt64(t.FieldNameForRowId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, newFieldValues, err := aepr.GetParameterValueAsJSON("new")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	err = t.DoEdit(aepr, id, newFieldValues)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func (t *DXTable) RequestEditByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	_, uid, err := aepr.GetParameterValueAsString(t.FieldNameForRowUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, newFieldValues, err := aepr.GetParameterValueAsJSON("new")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	err = t.DoEditByUid(aepr, uid, newFieldValues)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func (t *DXTable) DoDelete(aepr *api.DXAPIEndPointRequest, id int64) (err error) {
	_, _, err = t.ShouldGetById(&aepr.Log, id)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	_, err = db.Delete(t.Database.Connection, t.NameId, utils.JSON{
		t.FieldNameForRowId: id,
	})
	if err != nil {
		aepr.Log.Errorf(err, "Error at %s.DoDelete (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
	return nil
}

func (t *DXTable) DoDeleteByUid(aepr *api.DXAPIEndPointRequest, uid string) (err error) {
	_, _, err = t.ShouldGetByUid(&aepr.Log, uid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	_, err = db.Delete(t.Database.Connection, t.NameId, utils.JSON{
		t.FieldNameForRowUid: uid,
	})
	if err != nil {
		aepr.Log.Errorf(err, "Error at %s.DoDelete (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
	return nil
}

func (t *DXTable) SoftDelete(aepr *api.DXAPIEndPointRequest, id int64) (err error) {
	_, _, err = t.ShouldGetById(&aepr.Log, id)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	_, err = db.Update(t.Database.Connection, t.NameId, utils.JSON{
		"is_deleted": true,
	}, utils.JSON{
		t.FieldNameForRowId: id,
	})
	if err != nil {
		aepr.Log.Errorf(err, "Error at %s.DoDelete (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
	return nil
}

func (t *DXTable) RequestSoftDelete(aepr *api.DXAPIEndPointRequest) (err error) {
	_, id, err := aepr.GetParameterValueAsInt64(t.FieldNameForRowId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	err = t.DoEdit(aepr, id, utils.JSON{
		"is_deleted": true,
	})
	if err != nil {
		aepr.Log.Errorf(err, "Error at %s.RequestSoftDelete (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func (t *DXTable) RequestSoftDeleteByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	_, uid, err := aepr.GetParameterValueAsString(t.FieldNameForRowUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	err = t.DoEditByUid(aepr, uid, utils.JSON{
		"is_deleted": true,
	})
	if err != nil {
		aepr.Log.Errorf(err, "Error at %s.RequestSoftDelete (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func (t *DXTable) RequestHardDelete(aepr *api.DXAPIEndPointRequest) (err error) {
	_, id, err := aepr.GetParameterValueAsInt64(t.FieldNameForRowId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	err = t.DoDelete(aepr, id)
	if err != nil {
		aepr.Log.Errorf(err, "Error at %s.RequestHardDelete (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func (t *DXTable) RequestHardDeleteByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	_, uid, err := aepr.GetParameterValueAsString(t.FieldNameForRowUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	err = t.DoDeleteByUid(aepr, uid)
	if err != nil {
		aepr.Log.Errorf(err, "Error at %s.RequestHardDelete (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func (t *DXTable) SelectAll(log *log.DXLog) (rowsInfo *db.RowsInfo, r []utils.JSON, err error) {
	return t.Select(log, nil, nil, nil, map[string]string{t.FieldNameForRowId: "asc"}, nil)
}

func (t *DXTable) Count(log *log.DXLog, summaryCalcFieldsPart string, whereAndFieldNameValues utils.JSON, joinSQLPart any) (totalRows int64, summaryCalcRow utils.JSON, err error) {
	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{
			"is_deleted": false,
		}
		if t.Database == nil {
			t.Database = database.Manager.Databases[t.DatabaseNameId]
		}
		if t.Database.DatabaseType.String() == "sqlserver" {
			whereAndFieldNameValues["is_deleted"] = 0
		}
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	totalRows, summaryCalcRow, err = t.Database.CountOne(t.ListViewNameId, summaryCalcFieldsPart, whereAndFieldNameValues, joinSQLPart)
	return totalRows, summaryCalcRow, err
}

/*
	func (t *DXTable) TxSelectCount(tx *database.DXDatabaseTx, summaryCalcFieldsPart string, whereAndFieldNameValues utils.JSON) (totalRows int64, summaryCalcRow utils.JSON, err error) {
		if whereAndFieldNameValues == nil {
			whereAndFieldNameValues = utils.JSON{
				"is_deleted": false,
			}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

			if t.Database.DatabaseType.String() == "sqlserver" {
				whereAndFieldNameValues["is_deleted"] = 0
			}
		}

		totalRows, summaryCalcRow, err = tx.ShouldCount(t.ListViewNameId, summaryCalcFieldsPart, whereAndFieldNameValues)
		return totalRows, summaryCalcRow, err
	}
*/
func (t *DXTable) Select(log *log.DXLog, fieldNames []string, whereAndFieldNameValues utils.JSON, joinSQLPart any,
	orderByFieldNameDirections db.FieldsOrderBy, limit any) (rowsInfo *db.RowsInfo, r []utils.JSON, err error) {

	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{
			"is_deleted": false,
		}

		if t.Database == nil {
			t.Database = database.Manager.Databases[t.DatabaseNameId]
		}

		if t.Database.DatabaseType.String() == "sqlserver" {
			whereAndFieldNameValues["is_deleted"] = 0
		}
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	rowsInfo, r, err = t.Database.Select(t.ListViewNameId, t.FieldTypeMapping, fieldNames, whereAndFieldNameValues,
		joinSQLPart, orderByFieldNameDirections, limit, nil)
	if err != nil {
		return rowsInfo, nil, err
	}

	return rowsInfo, r, err
}

func (t *DXTable) ShouldSelectOne(log *log.DXLog, whereAndFieldNameValues utils.JSON, joinSQLPart any,
	orderbyFieldNameDirections db.FieldsOrderBy) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {

	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}
	whereAndFieldNameValues["is_deleted"] = false

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	return t.Database.ShouldSelectOne(t.ListViewNameId, t.FieldTypeMapping, nil, whereAndFieldNameValues, joinSQLPart, orderbyFieldNameDirections)
}

func (t *DXTable) TxShouldSelectOne(tx *database.DXDatabaseTx, whereAndFieldNameValues utils.JSON,
	orderbyFieldNameDirections db.FieldsOrderBy) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {

	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}
	whereAndFieldNameValues["is_deleted"] = false

	return tx.ShouldSelectOne(t.ListViewNameId, t.FieldTypeMapping, nil, whereAndFieldNameValues, nil, orderbyFieldNameDirections, nil)
}

func (t *DXTable) TxShouldSelectOneForUpdate(tx *database.DXDatabaseTx, whereAndFieldNameValues utils.JSON,
	orderbyFieldNameDirections db.FieldsOrderBy) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {

	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}
	whereAndFieldNameValues["is_deleted"] = false

	return tx.ShouldSelectOne(t.NameId, t.FieldTypeMapping, nil, whereAndFieldNameValues, nil, orderbyFieldNameDirections, true)
}

func (t *DXTable) TxSelect(tx *database.DXDatabaseTx, whereAndFieldNameValues utils.JSON,
	orderbyFieldNameDirections db.FieldsOrderBy, limit any) (rowsInfo *db.RowsInfo, r []utils.JSON, err error) {

	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}
	whereAndFieldNameValues["is_deleted"] = false

	return tx.Select(t.ListViewNameId, t.FieldTypeMapping, nil, whereAndFieldNameValues, nil, orderbyFieldNameDirections, limit, false)
}

func (t *DXTable) TxSelectOne(tx *database.DXDatabaseTx, whereAndFieldNameValues utils.JSON,
	orderbyFieldNameDirections db.FieldsOrderBy) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {

	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}
	whereAndFieldNameValues["is_deleted"] = false

	return tx.SelectOne(t.ListViewNameId, t.FieldTypeMapping, nil, whereAndFieldNameValues, nil, orderbyFieldNameDirections, false)
}

func (t *DXTable) TxSelectOneForUpdate(tx *database.DXDatabaseTx, whereAndFieldNameValues utils.JSON,
	orderbyFieldNameDirections db.FieldsOrderBy) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {

	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}
	whereAndFieldNameValues["is_deleted"] = false

	return tx.SelectOne(t.NameId, t.FieldTypeMapping, nil, whereAndFieldNameValues, nil, orderbyFieldNameDirections, true)
}

func (t *DXTable) TxUpsert(tx *database.DXDatabaseTx, setKeyValues utils.JSON, whereAndFieldNameValues utils.JSON) (result sql.Result, newId int64, err error) {
	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}

	_, r, err := tx.SelectOne(t.NameId, nil, nil, whereAndFieldNameValues, nil, nil, nil)
	if err != nil {
		return nil, 0, err
	}
	if r == nil {
		newSetKeyValues := utilsJson.DeepMerge2(setKeyValues, whereAndFieldNameValues)

		newId, err = tx.Insert(t.NameId, newSetKeyValues)
		return nil, newId, err
	} else {
		result, err = tx.Update(t.NameId, setKeyValues, whereAndFieldNameValues)
		return result, 0, err
	}
}

func (t *DXTable) TxUpdate(tx *database.DXDatabaseTx, setKeyValues utils.JSON, whereAndFieldNameValues utils.JSON) (result sql.Result, err error) {
	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}
	whereAndFieldNameValues["is_deleted"] = false

	return tx.Update(t.NameId, setKeyValues, whereAndFieldNameValues)
}

func (t *DXTable) TxSoftDelete(tx *database.DXDatabaseTx, whereAndFieldNameValues utils.JSON) (result sql.Result, err error) {
	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}

	return tx.Update(t.NameId, map[string]any{
		"is_deleted": true,
	}, whereAndFieldNameValues)
}

func (t *DXTable) TxHardDelete(tx *database.DXDatabaseTx, whereAndFieldNameValues utils.JSON) (r sql.Result, err error) {
	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}

	return tx.Delete(t.NameId, whereAndFieldNameValues)
}

func (t *DXTable) DoRequestList(aepr *api.DXAPIEndPointRequest, filterWhere string, filterOrderBy string, filterKeyValues utils.JSON, onResultList OnResultList) (err error) {
	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	if !t.Database.Connected {
		err := t.Database.Connect()
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	}

	rowsInfo, list, err := db.NamedQueryList(t.Database.Connection, t.FieldTypeMapping, "*", t.ListViewNameId,
		filterWhere, "", filterOrderBy, filterKeyValues)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	for i := range list {

		if onResultList != nil {
			aListRow, err := onResultList(list[i])
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			list[i] = aListRow
		}

	}

	data := utilsJson.Encapsulate(t.ResponseEnvelopeObjectName, utils.JSON{
		"list": utils.JSON{
			"rows":      list,
			"rows_info": rowsInfo,
		},
	})

	aepr.WriteResponseAsJSON(http.StatusOK, nil, data)

	return nil
}

func (t *DXTable) DoRequestPagingList(aepr *api.DXAPIEndPointRequest, filterWhere string, filterOrderBy string, filterKeyValues utils.JSON, onResultList OnResultList) (err error) {
	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	if !t.Database.Connected {
		err := t.Database.Connect()
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	}

	_, rowPerPage, err := aepr.GetParameterValueAsInt64("row_per_page")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, pageIndex, err := aepr.GetParameterValueAsInt64("page_index")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	rowsInfo, list, totalRows, totalPage, _, err := db.NamedQueryPaging(t.Database.Connection, t.FieldTypeMapping, "", rowPerPage, pageIndex, "*", t.ListViewNameId,
		filterWhere, "", filterOrderBy, filterKeyValues)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	for i := range list {

		if onResultList != nil {
			aListRow, err := onResultList(list[i])
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			list[i] = aListRow
		}

	}

	data := utils.JSON{
		"data": utils.JSON{
			"list": utils.JSON{
				"rows":       list,
				"total_rows": totalRows,
				"total_page": totalPage,
				"rows_info":  rowsInfo,
			},
		},
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, data)

	return nil
}

func (t *DXTable) RequestListAll(aepr *api.DXAPIEndPointRequest) (err error) {
	return t.DoRequestList(aepr, "", "", nil, nil)
}

func (t *DXTable) RequestList(aepr *api.DXAPIEndPointRequest) (err error) {
	isExistFilterWhere, filterWhere, err := aepr.GetParameterValueAsString("filter_where")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterWhere {
		filterWhere = ""
	}
	isExistFilterOrderBy, filterOrderBy, err := aepr.GetParameterValueAsString("filter_order_by")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterOrderBy {
		filterOrderBy = ""
	}

	isExistFilterKeyValues, filterKeyValues, err := aepr.GetParameterValueAsJSON("filter_key_values")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterKeyValues {
		filterKeyValues = nil
	}

	_, isDeletedIncluded, err := aepr.GetParameterValueAsBool("is_deleted", false)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	if !isDeletedIncluded {
		if filterWhere != "" {
			filterWhere = fmt.Sprintf("(%s) and ", filterWhere)
		}

		if t.Database == nil {
			t.Database = database.Manager.Databases[t.DatabaseNameId]
		}

		switch t.Database.DatabaseType.String() {
		case "sqlserver":
			filterWhere = filterWhere + "(is_deleted=0)"
		case "postgres":
			filterWhere = filterWhere + "(is_deleted=false)"
		default:
			filterWhere = filterWhere + "(is_deleted=0)"
		}
	}

	return t.DoRequestList(aepr, filterWhere, filterOrderBy, filterKeyValues, nil)
}

func (t *DXTable) RequestPagingList(aepr *api.DXAPIEndPointRequest) (err error) {
	isExistFilterWhere, filterWhere, err := aepr.GetParameterValueAsString("filter_where")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterWhere {
		filterWhere = ""
	}
	isExistFilterOrderBy, filterOrderBy, err := aepr.GetParameterValueAsString("filter_order_by")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterOrderBy {
		filterOrderBy = ""
	}

	isExistFilterKeyValues, filterKeyValues, err := aepr.GetParameterValueAsJSON("filter_key_values")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterKeyValues {
		filterKeyValues = nil
	}

	_, isDeletedIncluded, err := aepr.GetParameterValueAsBool("is_deleted", false)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	if !isDeletedIncluded {
		if filterWhere != "" {
			filterWhere = fmt.Sprintf("(%s) and ", filterWhere)
		}

		if t.Database == nil {
			t.Database = database.Manager.Databases[t.DatabaseNameId]
		}

		switch t.Database.DatabaseType.String() {
		case "sqlserver":
			filterWhere = filterWhere + "(is_deleted=0)"
		case "postgres":
			filterWhere = filterWhere + "(is_deleted=false)"
		default:
			filterWhere = filterWhere + "(is_deleted=0)"
		}
	}

	return t.DoRequestPagingList(aepr, filterWhere, filterOrderBy, filterKeyValues, nil)
}

func (t *DXTable) RequestPagingListAll(aepr *api.DXAPIEndPointRequest) (err error) {
	filterWhere := ""
	filterOrderBy := ""
	isDeletedIncluded := false

	if !isDeletedIncluded {
		if filterWhere != "" {
			filterWhere = fmt.Sprintf("(%s) and ", filterWhere)
		}

		if t.Database == nil {
			t.Database = database.Manager.Databases[t.DatabaseNameId]
		}

		switch t.Database.DatabaseType.String() {
		case "sqlserver":
			filterWhere = filterWhere + "(is_deleted=0)"
		case "postgres":
			filterWhere = filterWhere + "(is_deleted=false)"
		default:
			filterWhere = filterWhere + "(is_deleted=0)"
		}
	}

	return t.DoRequestPagingList(aepr, filterWhere, filterOrderBy, nil, nil)
}

func (t *DXTable) SelectOne(log *log.DXLog, fieldNames []string, whereAndFieldNameValues utils.JSON, joinSQLPart any, orderbyFieldNameDirections db.FieldsOrderBy) (
	rowsInfo *db.RowsInfo, r utils.JSON, err error) {

	if whereAndFieldNameValues == nil {
		whereAndFieldNameValues = utils.JSON{}
	}
	whereAndFieldNameValues["is_deleted"] = false

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	return t.Database.SelectOne(t.ListViewNameId, t.FieldTypeMapping, fieldNames, whereAndFieldNameValues, joinSQLPart, orderbyFieldNameDirections)
}

func (t *DXTable) IsFieldValueExistAsString(log *log.DXLog, fieldName string, fieldValue string) (bool, error) {
	_, r, err := t.SelectOne(log, nil, utils.JSON{
		fieldName: fieldValue,
	}, nil, nil)
	if err != nil {
		return false, err
	}
	if r == nil {
		return false, nil
	}
	return true, nil
}

func (t *DXTable) RequestCreate(aepr *api.DXAPIEndPointRequest) (err error) {
	p := map[string]interface{}{}
	for k, v := range aepr.ParameterValues {
		p[k] = v.Value
	}
	_, err = t.DoCreate(aepr, p)
	if err != nil {
		return err
	}
	return nil
}

func (t *DXTable) RequestListDownload(aepr *api.DXAPIEndPointRequest) (err error) {
	isExistFilterWhere, filterWhere, err := aepr.GetParameterValueAsString("filter_where")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterWhere {
		filterWhere = ""
	}
	isExistFilterOrderBy, filterOrderBy, err := aepr.GetParameterValueAsString("filter_order_by")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterOrderBy {
		filterOrderBy = ""
	}

	isExistFilterKeyValues, filterKeyValues, err := aepr.GetParameterValueAsJSON("filter_key_values")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterKeyValues {
		filterKeyValues = nil
	}

	_, format, err := aepr.GetParameterValueAsString("format")
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "FORMAT_PARAMETER_ERROR:%s", err.Error())
	}

	format = strings.ToLower(format)

	isDeletedIncluded := false
	if !isDeletedIncluded {
		if filterWhere != "" {
			filterWhere = fmt.Sprintf("(%s) and ", filterWhere)
		}

		if t.Database == nil {
			t.Database = database.Manager.Databases[t.DatabaseNameId]
		}

		switch t.Database.DatabaseType.String() {
		case "sqlserver":
			filterWhere = filterWhere + "(is_deleted=0)"
		case "postgres":
			filterWhere = filterWhere + "(is_deleted=false)"
		default:
			filterWhere = filterWhere + "(is_deleted=0)"
		}
	}

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	if !t.Database.Connected {
		err := t.Database.Connect()
		if err != nil {
			aepr.Log.Errorf(err, "error At reconnect db At table %s list (%s) ", t.NameId, err.Error())
			return errors.Wrap(err, "error occured")
		}
	}

	rowsInfo, list, err := db.NamedQueryList(t.Database.Connection, t.FieldTypeMapping, "*", t.ListViewNameId,
		filterWhere, "", filterOrderBy, filterKeyValues)

	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	// Set export options
	opts := export.ExportOptions{
		Format:     export.ExportFormat(format),
		SheetName:  "Sheet1",
		DateFormat: "2006-01-02 15:04:05",
	}

	// Get file as stream
	data, contentType, err := export.ExportToStream(rowsInfo, list, opts)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	// Set response headers
	filename := fmt.Sprintf("export_%s_%s.%s", t.NameId, time.Now().Format("20060102_150405"), format)

	responseWriter := *aepr.GetResponseWriter()
	responseWriter.Header().Set("Content-Type", contentType)
	responseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	responseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	responseWriter.WriteHeader(http.StatusOK)
	aepr.ResponseStatusCode = http.StatusOK

	_, err = responseWriter.Write(data)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.ResponseHeaderSent = true
	aepr.ResponseBodySent = true

	return nil
}
