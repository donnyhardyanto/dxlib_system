package database

import (
	"database/sql"
	"github.com/donnyhardyanto/dxlib/database/protected/db"
	"github.com/donnyhardyanto/dxlib/database/protected/dbtx"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DXDatabaseTxCallback func(dtx *DXDatabaseTx) (err error)

type DXDatabaseTxIsolationLevel = sql.IsolationLevel

const (
	LevelDefault DXDatabaseTxIsolationLevel = iota
	LevelReadUncommitted
	LevelReadCommitted
	LevelWriteCommitted
	LevelRepeatableRead
	LevelSnapshot
	LevelSerializable
	LevelLinearizable
)

type DXDatabaseTx struct {
	*sqlx.Tx
	Log *log.DXLog
}

func (dtx *DXDatabaseTx) Commit() (err error) {
	err = dtx.Tx.Commit()
	if err != nil {
		dtx.Log.Errorf(err, "TX_ERROR_IN_COMMIT: (%v)", err.Error())
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func (dtx *DXDatabaseTx) Rollback() (err error) {
	err = dtx.Tx.Rollback()
	if err != nil {
		dtx.Log.Errorf(err, "TX_ERROR_IN_ROLLBACK: (%v)", err.Error())
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func (dtx *DXDatabaseTx) Finish(log *log.DXLog, err error) {
	if err != nil {
		err2 := dtx.Rollback()
		if err2 != nil {
			log.Errorf(err, "ROLLBACK_ERROR:%+v", err2)
		}
	} else {
		err2 := dtx.Commit()
		if err2 != nil {
			log.Errorf(err2, "ROLLBACK_ERROR:%+v", err2)
		}
	}
}

func (dtx *DXDatabaseTx) Select(tableName string, fieldTypeMapping db.FieldTypeMapping, fieldNames []string, whereAndFieldNameValues utils.JSON, joinSQLPart any,
	orderbyFieldNameDirections db.FieldsOrderBy, limit any, forUpdatePart any) (rowsInfo *db.RowsInfo, r []utils.JSON, err error) {

	return dbtx.TxSelect(dtx.Log, false, dtx.Tx, tableName, fieldTypeMapping, fieldNames, whereAndFieldNameValues, joinSQLPart, orderbyFieldNameDirections, limit, forUpdatePart)
}

func (dtx *DXDatabaseTx) SelectOne(tableName string, fieldTypeMapping db.FieldTypeMapping, fieldNames []string, whereAndFieldNameValues utils.JSON, joinSQLPart any,
	orderbyFieldNameDirections db.FieldsOrderBy, forUpdatePart any) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {

	return dbtx.TxSelectOne(dtx.Log, false, dtx.Tx, tableName, fieldTypeMapping, fieldNames, whereAndFieldNameValues, joinSQLPart, orderbyFieldNameDirections, forUpdatePart)
}

func (dtx *DXDatabaseTx) ShouldSelectOne(tableName string, fieldTypeMapping db.FieldTypeMapping, fieldNames []string, whereAndFieldNameValues utils.JSON, joinSQLPart any,
	orderbyFieldNameDirections db.FieldsOrderBy, forUpdatePart any) (rowsInfo *db.RowsInfo, r utils.JSON, err error) {
	return dbtx.TxShouldSelectOne(dtx.Log, nil, false, dtx.Tx, tableName, fieldNames, whereAndFieldNameValues, joinSQLPart, orderbyFieldNameDirections, forUpdatePart)
}
func (dtx *DXDatabaseTx) Insert(tableName string, keyValues utils.JSON) (id int64, err error) {
	return dbtx.TxInsert(dtx.Log, false, dtx.Tx, tableName, keyValues)
}

/*func (dtx *DXDatabaseTx) UpdateOne(tableName string, setKeyValues utils.JSON, whereKeyValues utils.JSON) (result sql.Result, err error) {
	return dbtx.TxUpdateOne(dtx.Log, false, dtx.Tx, tableName, setKeyValues, whereKeyValues)
}*/

func (dtx *DXDatabaseTx) Update(tableName string, setKeyValues utils.JSON, whereKeyValues utils.JSON) (result sql.Result, err error) {
	return dbtx.TxUpdate(dtx.Log, false, dtx.Tx, tableName, setKeyValues, whereKeyValues)
}

/*
	func (dtx *DXDatabaseTx) RequestSoftDelete(tableName string, whereKeyValues utils.JSON) (result sql.Result, err error) {
		return dbtx.TxUpdate(dtx.Log, false, dtx.Tx, tableName, utils.JSON{
			"is_deleted": true,
		}, whereKeyValues)
	}
*/
func (dtx *DXDatabaseTx) Delete(tableName string, whereKeyValues utils.JSON) (result sql.Result, err error) {
	return dbtx.TxDelete(dtx.Log, false, dtx.Tx, tableName, whereKeyValues)
}
