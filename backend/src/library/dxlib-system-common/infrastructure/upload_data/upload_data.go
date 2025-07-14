package upload_data

import (
	dxlibModule "github.com/donnyhardyanto/dxlib/module"
	"github.com/donnyhardyanto/dxlib/table"
)

type UploadData struct {
	dxlibModule.DXModule
	Organization *table.DXTable
	User         *table.DXTable
	Customer     *table.DXTable
	Arrears      *table.DXTable
}

var ModuleUploadData = UploadData{}

func (am *UploadData) Init(aDatabaseNameId string) {
	am.DatabaseNameId = aDatabaseNameId
	am.Organization = table.Manager.NewTable(am.DatabaseNameId, "upload_data.organization", "upload_data.organization",
		"upload_data.organization", "id", "id", "uid", "data")

	am.User = table.Manager.NewTable(am.DatabaseNameId, "upload_data.user", "upload_data.user",
		"upload_data.user", "id", "id", "uid", "data")

	am.Customer = table.Manager.NewTable(am.DatabaseNameId, "upload_data.customer", "upload_data.customer",
		"upload_data.customer", "id", "id", "uid", "data")

	am.Arrears = table.Manager.NewTable(am.DatabaseNameId, "upload_data.arrears", "upload_data.arrears",
		"upload_data.arrears", "id", "id", "uid", "data")

}
