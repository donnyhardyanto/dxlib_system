package construction_management

import (
	dxlibModule "github.com/donnyhardyanto/dxlib/module"
	"github.com/donnyhardyanto/dxlib/table"
)

type ConstructionManagement struct {
	dxlibModule.DXModule

	GasAppliance           *table.DXTable
	TappingSaddleAppliance *table.DXTable
	MeterApplianceType     *table.DXTable
	RegulatorAppliance     *table.DXTable
	GSize                  *table.DXTable
}

func (cm *ConstructionManagement) Init(databaseNameId string) {
	cm.DatabaseNameId = databaseNameId
	cm.GasAppliance = table.Manager.NewTable(cm.DatabaseNameId, "construction_management.gas_appliance", "construction_management.gas_appliance",
		"construction_management.gas_appliance", "code", "id", "uid", "data")
	cm.TappingSaddleAppliance = table.Manager.NewTable(cm.DatabaseNameId, "construction_management.tapping_saddle_appliance", "construction_management.tapping_saddle_appliance",
		"construction_management.tapping_saddle_appliance", "code", "id", "uid", "data")
	cm.MeterApplianceType = table.Manager.NewTable(cm.DatabaseNameId, "construction_management.meter_appliance_type", "construction_management.meter_appliance_type",
		"construction_management.meter_appliance_type", "code", "id", "uid", "data")
	cm.RegulatorAppliance = table.Manager.NewTable(cm.DatabaseNameId, "construction_management.regulator_appliance", "construction_management.regulator_appliance",
		"construction_management.regulator_appliance", "code", "id", "uid", "data")
	cm.GSize = table.Manager.NewTable(cm.DatabaseNameId, "construction_management.g_size", "construction_management.g_size",
		"construction_management.g_size", "code", "id", "uid", "data")
}

var ModuleConstructionManagement = ConstructionManagement{}
