package master_data

import (
	"github.com/donnyhardyanto/dxlib/log"
	dxlibModule "github.com/donnyhardyanto/dxlib/module"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
)

type MasterData struct {
	dxlibModule.DXModule

	Area             *table.DXTable
	Location         *table.DXTable
	CustomerRef      *table.DXTable
	GlobalLookup     *table.DXTable
	RsCustomerSector *table.DXTable
	CustomerSegment  *table.DXTable
	CustomerType     *table.DXTable
	PaymentScheme    *table.DXTable
}

var ModuleMasterData = MasterData{}

func (md *MasterData) Init(databaseNameId string) {
	md.DatabaseNameId = databaseNameId

	md.Area = table.Manager.NewTable(base.DatabaseNameIdTaskDispatcher, "master_data.area", "master_data.area",
		"master_data.area", "code", "id", "uid", "data")
	md.Location = table.Manager.NewTable(base.DatabaseNameIdTaskDispatcher, "master_data.location", "master_data.location",
		"master_data.location", "code", "id", "uid", "data")
	md.CustomerRef = table.Manager.NewTable(base.DatabaseNameIdTaskDispatcher, "master_data.customer_ref", "master_data.customer_ref",
		"master_data.customer_ref", "code", "id", "uid", "data")
	md.GlobalLookup = table.Manager.NewTable(base.DatabaseNameIdTaskDispatcher, "master_data.global_lookup", "master_data.global_lookup",
		"master_data.global_lookup", "code", "id", "uid", "data")
	md.RsCustomerSector = table.Manager.NewTable(base.DatabaseNameIdTaskDispatcher, "master_data.rs_customer_sector", "master_data.rs_customer_sector",
		"master_data.rs_customer_sector", "code", "id", "uid", "data")
	md.CustomerSegment = table.Manager.NewTable(base.DatabaseNameIdTaskDispatcher, "master_data.vm_customer_segment", "master_data.vw_customer_segment",
		"master_data.vw_customer_segment", "code", "id", "uid", "data")
	md.CustomerType = table.Manager.NewTable(base.DatabaseNameIdTaskDispatcher, "master_data.vw_customer_type", "master_data.vw_customer_type",
		"master_data.vw_customer_type", "code", "id", "uid", "data")
	md.PaymentScheme = table.Manager.NewTable(base.DatabaseNameIdTaskDispatcher, "master_data.vw_payment_scheme", "master_data.vw_payment_scheme",
		"master_data.vw_payment_scheme", "code", "id", "uid", "data")
}

/* Only 1 level up */
func (md *MasterData) AreaCodeExpandParentTreeUp(log *log.DXLog, code string) (areaCodes []string, err error) {
	_, area, err := md.Area.ShouldGetByNameId(log, code)
	if err != nil {
		return nil, errors.Errorf("IMPOSSIBLE:AREA[CODE]==NOT_FOUND")
	}
	var parentAreaCode string
	areaCodes = []string{code}

	areaParentValue, ok := area["parent_value"]
	if !ok {
		return nil, errors.Errorf("IMPOSSIBLE:AREA[PARENT_VALUE]==NOT_FOUND")
	}
	if areaParentValue != nil {
		_, parentArea, err := md.Area.SelectOne(log, nil, utils.JSON{
			"name": areaParentValue.(string),
		}, nil, nil)
		if err != nil {
			return nil, err
		}
		parentAreaCode, ok = parentArea["code"].(string)
		if !ok {
			return nil, errors.Errorf("IMPOSSIBLE:PARENT_AREA[CODE]==NOT_FOUND")
		}
		areaCodes = append(areaCodes, parentAreaCode)
	}
	return areaCodes, nil
}

/* Only 1 level down */
func (md *MasterData) AreaCodeExpandParentTreeDown(log *log.DXLog, code string) (areaCodes []string, err error) {
	_, area, err := md.Area.ShouldGetByNameId(log, code)
	if err != nil {
		return nil, errors.Errorf("IMPOSSIBLE:AREA[CODE]==NOT_FOUND")
	}
	_, areaChildren, err := md.Area.Select(log, nil, utils.JSON{
		"parent_group": area["type"].(string),
		"parent_value": area["name"].(string),
	}, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	areaCodes = []string{code}
	for _, areaChild := range areaChildren {
		childAreaCode, ok := areaChild["code"].(string)
		if !ok {
			return nil, errors.Errorf("IMPOSSIBLE:CHILD_AREA[CODE]==NOT_FOUND")
		}
		areaCodes = append(areaCodes, childAreaCode)
	}
	return areaCodes, nil
}
