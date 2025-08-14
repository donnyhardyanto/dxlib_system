package database

import (
	"database/sql"
	"github.com/donnyhardyanto/dxlib/log"
)

type DXDatabaseScript struct {
	Owner              *DXDatabaseManager
	NameId             string
	ManagementDatabase *DXDatabase
	Files              []string
}

func (dm *DXDatabaseManager) NewDatabaseScript(nameId string, files []string) *DXDatabaseScript {
	ds := DXDatabaseScript{
		Owner:  dm,
		NameId: nameId,
		Files:  files,
	}
	dm.Scripts[nameId] = &ds
	return &ds
}

func (ds *DXDatabaseScript) ExecuteFile(d *DXDatabase, filename string) (r sql.Result, err error) {
	log.Log.Infof("Executing SQL file %s... start", filename)
	r, err = d.ExecuteFile(filename)
	if err != nil {
		return nil, err
	}
	log.Log.Infof("Executing SQL file %s... done", filename)
	return r, nil
}

func (ds *DXDatabaseScript) Execute(d *DXDatabase) (rs []sql.Result, err error) {
	rs = []sql.Result{}
	for k, v := range ds.Files {
		r, err := ds.ExecuteFile(d, v)
		if err != nil {
			log.Log.Errorf(err, "Error executing file %d:'%s' (%s)", k, v, err.Error())
			return rs, err
		}
		rs = append(rs, r)
	}
	return rs, nil
}
