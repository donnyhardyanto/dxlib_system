package handler

import (
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/configuration_settings"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/app"
	"github.com/donnyhardyanto/dxlib/configuration"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/utils"
	"time"
)

func DoOnUserAfterCreate(aepr *api.DXAPIEndPointRequest, dtx *database.DXDatabaseTx, user utils.JSON, userPassword string) (err error) {
	isNoSendEmail, ok := app.App.LocalData["user-create-no-send-email"].(bool)
	if ok && isNoSendEmail {
		return nil
	}
	configExternalSystem := *configuration.Manager.Configurations["external_system"].Data

	go func() {
		smtpConfiguration, ok := configExternalSystem["SMTP1"].(utils.JSON)
		if !ok {
			aepr.Log.Warn("USER_AFTER_CREATE:SMTP_CONFIG_NOT_FOUND")
			return
		}

		_, emailTemplate, err := configuration_settings.ModuleConfigurationSettings.EMailTemplate.ShouldGetByNameId(&aepr.Log, "USER_REGISTRATION")
		if err != nil {
			aepr.Log.Warnf("USER_AFTER_CREATE:USER_REGISTRATION_EMAIL_TEMPLATE_NOT_FOUND:%s", err.Error())
			return
		}
		emailTemplateContentType := emailTemplate["content_type"].(string)
		emailTemplateTitle := emailTemplate["subject"].(string)
		emailTemplateBody := emailTemplate["body"].(string)

		aUserLoginId := user["loginid"].(string)
		aUserEmail := user["email"].(string)
		aUserFullname := user["fullname"].(string)

		data := utils.JSON{
			"fullname": aUserFullname,
			"loginid":  aUserLoginId,
			"password": userPassword,
		}

		err = base.EmailSend(data, emailTemplateContentType, emailTemplateTitle, emailTemplateBody, smtpConfiguration, aUserEmail)
		if err != nil {
			aepr.Log.Warnf("USER_AFTER_CREATE:SEND_MAIL_ERROR:%s", err.Error())
			return
		}
		time.Sleep(5 * time.Second)
	}()

	go func() {
		smsConfiguration, ok := configExternalSystem["SMS1"].(utils.JSON)
		if !ok {
			aepr.Log.Warn("USER_AFTER_CREATE:SMS_CONFIG_NOT_FOUND")
			return
		}
		if !smsConfiguration["enabled"].(bool) {
			return
		}
		if !smsConfiguration["is_send_when_user_created"].(bool) {
			return
		}
		_, smsTemplate, err := configuration_settings.ModuleConfigurationSettings.SMSTemplate.ShouldGetByNameId(&aepr.Log, "USER_REGISTRATION")
		if err != nil {
			aepr.Log.Warnf("USER_AFTER_CREATE:USER_REGISTRATION_SMS_TEMPLATE_NOT_FOUND:%s", err.Error())
			return
		}

		smsTemplateBody := smsTemplate["body"].(string)

		aUserLoginId := user["loginid"].(string)
		aUserPhoneNumber := user["phonenumber"].(string)
		aUserFullname := user["fullname"].(string)

		data := utils.JSON{
			"fullname": aUserFullname,
			"loginid":  aUserLoginId,
			"password": userPassword,
		}

		err = base.SMSSend(aUserPhoneNumber, data, smsTemplateBody, smsConfiguration)
		if err != nil {
			aepr.Log.Warnf("USER_AFTER_CREATE:SEND_MAIL_ERROR:%s", err.Error())
			return
		}
	}()

	return nil
}

func DoOnUserResetPassword(aepr *api.DXAPIEndPointRequest, dtx *database.DXDatabaseTx, user utils.JSON, userPasswordNew string) (err error) {
	configExternalSystem := *configuration.Manager.Configurations["external_system"].Data

	go func() {
		smtpConfiguration, ok := configExternalSystem["SMTP1"].(utils.JSON)
		if !ok {
			aepr.Log.Warn("USER_AFTER_CREATE:SMTP_CONFIG_NOT_FOUND")
			return
		}

		_, emailTemplate, err := configuration_settings.ModuleConfigurationSettings.EMailTemplate.ShouldGetByNameId(&aepr.Log, "USER_RESET_PASSWORD")
		if err != nil {
			aepr.Log.Warnf("USER_RESET_PASSWORD:USER_RESET_PASSWORD_EMAIL_TEMPLATE_NOT_FOUND:%s", err.Error())
			return
		}
		emailTemplateContentType := emailTemplate["content_type"].(string)
		emailTemplateTitle := emailTemplate["subject"].(string)
		emailTemplateBody := emailTemplate["body"].(string)

		aUserLoginId := user["loginid"].(string)
		aUserEmail := user["email"].(string)
		aUserFullname := user["fullname"].(string)

		data := utils.JSON{
			"fullname": aUserFullname,
			"loginid":  aUserLoginId,
			"password": userPasswordNew,
		}

		err = base.EmailSend(data, emailTemplateContentType, emailTemplateTitle, emailTemplateBody, smtpConfiguration, aUserEmail)
		if err != nil {
			aepr.Log.Warnf("USER_RESET_PASSWORD:SEND_MAIL_ERROR:%s", err.Error())
			return
		}
	}()

	go func() {
		smsConfiguration, ok := configExternalSystem["SMS1"].(utils.JSON)
		if !ok {
			aepr.Log.Errorf(err, "USER_ON_RESET_PASSWORD:SMS_CONFIG_NOT_FOUND")
		}
		if !smsConfiguration["enabled"].(bool) {
			return
		}
		if !smsConfiguration["is_send_when_user_reset_password"].(bool) {
			return
		}
		_, smsTemplate, err := configuration_settings.ModuleConfigurationSettings.SMSTemplate.ShouldGetByNameId(&aepr.Log, "USER_RESET_PASSWORD")
		if err != nil {
			aepr.Log.Warnf("USER_ON_RESET_PASSWORD:USER_RESET_PASSWORD_SMS_TEMPLATE_NOT_FOUND:%s", err.Error())
			return
		}
		smsTemplateBody := smsTemplate["body"].(string)

		aUserLoginId := user["loginid"].(string)
		aUserPhoneNumber := user["phonenumber"].(string)
		aUserFullname := user["fullname"].(string)

		data := utils.JSON{
			"fullname": aUserFullname,
			"loginid":  aUserLoginId,
			"password": userPasswordNew,
		}

		err = base.SMSSend(aUserPhoneNumber, data, smsTemplateBody, smsConfiguration)
		if err != nil {
			aepr.Log.Warnf("USER_ON_RESET_PASSWORD:SEND_MAIL_ERROR:%s", err.Error())
			return
		}
	}()

	return nil
}
