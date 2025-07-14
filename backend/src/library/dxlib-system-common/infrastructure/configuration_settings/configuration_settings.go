package configuration_settings

import (
	"github.com/donnyhardyanto/dxlib/log"
	dxlibModule "github.com/donnyhardyanto/dxlib/module"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
)

type ConfigurationSettings struct {
	dxlibModule.DXModule
	EMailTemplate   *table.DXTable
	SMSTemplate     *table.DXTable
	GeneralTemplate *table.DXTable
}

var ModuleConfigurationSettings = ConfigurationSettings{}

func (cs *ConfigurationSettings) Init(aDatabaseNameId string) {
	cs.DatabaseNameId = aDatabaseNameId
	cs.EMailTemplate = table.Manager.NewTable(cs.DatabaseNameId,
		"settings.email_template", "settings.email_template",
		"settings.email_template", "nameid", "id", "uid", "data")
	cs.SMSTemplate = table.Manager.NewTable(cs.DatabaseNameId,
		"settings.sms_template", "settings.sms_template",
		"settings.sms_template", "nameid", "id", "uid", "data")
	cs.GeneralTemplate = table.Manager.NewTable(cs.DatabaseNameId,
		"settings.general_template", "settings.general_template",
		"settings.general_template", "nameid", "id", "uid", "data")
}

func (cs *ConfigurationSettings) GeneralTemplateGetByNameId(l *log.DXLog, nameId string) (gt utils.JSON, templateTitle string, templateBody string, err error) {
	_, templateMessage, err := ModuleConfigurationSettings.GeneralTemplate.ShouldGetByNameId(l, nameId)
	if err != nil {
		return nil, "", "", err
	}
	templateTitle, ok := templateMessage["title"].(string)
	if !ok {
		return nil, "", "", errors.New("INVALID_TEMPLATE_TITLE")
	}
	templateBody, ok = templateMessage["body"].(string)
	if !ok {
		return nil, "", "", errors.New("INVALID_TEMPLATE_BODY")
	}
	return templateMessage, templateTitle, templateBody, nil
}
