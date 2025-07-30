package base

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"io"
	"net/http"
	"strings"
	"time"
)

type UserType string

var (
	RelyOnInboundCredentials      []utils.JSON
	RelyOnInboundRedisSessionName = "session"
	DbReady                       = false
)

const (
	DatabaseNameIdAuditLog = "auditlog"
	DatabaseNameIdDbBase   = "db_base"
	DatabaseNameIdConfig   = "config"
)

var SubTaskFormReportAPIEndPointParameter = []api.DXAPIEndPointParameter{
	{NameId: "sk", Type: "json", Description: "Report SK", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
		{NameId: "pipe_length", Type: "float64zp", Description: "Report SK pipe_length", IsMustExist: true},
		{NameId: "calculated_extra_pipe_length", Type: "float64zp", Description: "Report SK calculated_extra_pipe_length", IsMustExist: true},
		{NameId: "test_start_time", Type: "non-empty-string", Description: "Report SK test_start_time", IsMustExist: true},
		{NameId: "test_end_time", Type: "non-empty-string", Description: "Report SK test_end_time", IsMustExist: true},
		{NameId: "calculated_test_duration_minute", Type: "float64p", Description: "Report SK calculated_test_duration_minute", IsMustExist: true},
		{NameId: "test_pressure", Type: "float64zp", Description: "Report SK test_pressure", IsMustExist: true},
		{NameId: "finished_date", Type: "date", Description: "Date of finished work", IsMustExist: true},
		{NameId: "gas_appliances", Type: "array-json-template", Description: "Report SK gas Appliance", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64p", Description: "Report, SK, Gas, Appliance gas_appliance_id", IsMustExist: true},
			{NameId: "quantity", Type: "int64zp", Description: "Report SK Gas Appliance quantity", IsMustExist: true},
		}},
	}},
	{NameId: "sr", Type: "json", Description: "Report SK", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
		{NameId: "tapping_saddle_id", Type: "int64p", Description: "Report SR quantity", IsMustExist: true},
		{NameId: "tapping_saddle_custom", Type: "string", Description: "Report SR quantity", IsMustExist: false}, // if tapping_saddle_id is 0
		{NameId: "test_start_time", Type: "non-empty-string", Description: "Report SR sk_report", IsMustExist: true},
		{NameId: "test_end_time", Type: "non-empty-string", Description: "Report SR gas_appliance_id", IsMustExist: true},
		{NameId: "calculated_test_duration_minute", Type: "float64p", Description: "Report SR quantity", IsMustExist: true},
		{NameId: "test_pressure", Type: "float64zp", Description: "Report SR sk_report", IsMustExist: true},
		{NameId: "branch_pipe_availability", Type: "bool", Description: "Report SR gas_appliance_id", IsMustExist: true},
		{NameId: "finished_date", Type: "date", Description: "Date of finished work", IsMustExist: true},
	}},
	{NameId: "meter_installation", Type: "json", Description: "Report SK", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
		{NameId: "meter_id", Type: "int64p", Description: "Report Gas Meter Installation quantity", IsMustExist: true},
		{NameId: "meter_brand", Type: "non-empty-string", Description: "Report Gas Meter Installation sk_report", IsMustExist: true},
		{NameId: "sn_meter", Type: "non-empty-string", Description: "Report Gas Meter Installation gas_appliance_id", IsMustExist: true},
		{NameId: "g_size_id", Type: "int64p", Description: "Report Gas Meter Installation quantity", IsMustExist: true},
		{NameId: "qmin", Type: "float64zp", Description: "Report Gas Meter Installation sk_report", IsMustExist: true},
		{NameId: "qmax", Type: "float64zp", Description: "Report Gas Meter Installation gas_appliance_id", IsMustExist: true},
		{NameId: "pmax", Type: "float64zp", Description: "Report Gas Meter Installation quantity", IsMustExist: true},
		{NameId: "start_calibration_month", Type: "int64p", Description: "Report Gas Meter Installation start_calibration_month", IsMustExist: true},
		{NameId: "start_calibration_year", Type: "int64p", Description: "Report Gas Meter Installation start_calibration_year", IsMustExist: true},
		{NameId: "regulator_brand", Type: "non-empty-string", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "regulator_size_inch", Type: "float64", Description: "Report Gas In sk_report", IsMustExist: false},
	}},
	{NameId: "gas_in", Type: "json", Description: "Report Gas In", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
		{NameId: "meter_id", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "meter_brand", Type: "non-empty-string", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "sn_meter", Type: "non-empty-string", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "g_size_id", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "pmax", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "stand_meter_start_number", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "pressure_start", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "temperature_start", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "meter_location_longitude", Type: "float64", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "meter_location_latitude", Type: "float64", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "gas_in_date", Type: "date", Description: "Date of actual gas in", IsMustExist: true},
		{NameId: "gas_appliances", Type: "array-json-template", Description: "Report SK gas Appliance", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64p", Description: "Report Gas In Gas Appliance gas_appliance_id", IsMustExist: true},
			{NameId: "quantity", Type: "int64zp", Description: "Report  Gas In Gas Appliance quantity", IsMustExist: true},
		}},
	}},

	// penanganan piutang

	{NameId: "pp_tutup_aliran", Type: "json", Description: "Report Gas In", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
		{NameId: "meter_id", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "meter_brand", Type: "non-empty-string", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "sn_meter", Type: "non-empty-string", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "g_size_id", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "calibration_end_month", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "calibration_end_year", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "stand_meter_number", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "pressure_start", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "temperature_start", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "meter_location_longitude", Type: "float64", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "meter_location_latitude", Type: "float64", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "seal_no", Type: "string", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "condition", Type: "string", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "condition_notes", Type: "string", Description: "Report Gas In sk_report", IsMustExist: false},
	}},
	{NameId: "pp_buka_aliran", Type: "json", Description: "Report Gas In", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
		{NameId: "meter_id", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "meter_brand", Type: "non-empty-string", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "sn_meter", Type: "non-empty-string", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "g_size_id", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "calibration_end_month", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "calibration_end_year", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "stand_meter_number", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "pressure_start", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "temperature_start", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "meter_location_longitude", Type: "float64", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "meter_location_latitude", Type: "float64", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "seal_no", Type: "string", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "condition", Type: "string", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "condition_notes", Type: "string", Description: "Report Gas In sk_report", IsMustExist: false},
	}},

	{NameId: "pp_cabut_meter_gas", Type: "json", Description: "Report Gas In", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
		{NameId: "meter_id", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "meter_brand", Type: "non-empty-string", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "sn_meter", Type: "non-empty-string", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "g_size_id", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "calibration_month", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "calibration_year", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "stand_meter_number", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "pressure_start", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "temperature_start", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "meter_location_longitude", Type: "float64", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "meter_location_latitude", Type: "float64", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "seal_no", Type: "string", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "condition", Type: "string", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "condition_notes", Type: "string", Description: "Report Gas In sk_report", IsMustExist: false},
	}},
	{NameId: "pp_pasang_meter_gas", Type: "json", Description: "Report Gas In", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
		{NameId: "meter_id", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "meter_brand", Type: "string", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "sn_meter", Type: "string", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "g_size_id", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "calibration_month", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "calibration_year", Type: "int64p", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "stand_meter_number", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "pressure_start", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "temperature_start", Type: "float64zp", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "meter_location_longitude", Type: "float64", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "meter_location_latitude", Type: "float64", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "seal_no", Type: "string", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "condition", Type: "string", Description: "Report Gas In sk_report", IsMustExist: false},
		{NameId: "condition_notes", Type: "string", Description: "Report Gas In sk_report", IsMustExist: false},
	}},
}

var SubTaskFormCancelAPIEndPointParameter = []api.DXAPIEndPointParameter{
	{NameId: "construction", Type: "json", Description: "Report SK", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
		{NameId: "pipe_length", Type: "float64zp", Description: "Report SK pipe_length", IsMustExist: true},
		{NameId: "calculated_extra_pipe_length", Type: "float64zp", Description: "Report SK calculated_extra_pipe_length", IsMustExist: true},
		{NameId: "test_start_time", Type: "string", Description: "Report SK test_start_time", IsMustExist: true},
		{NameId: "test_end_time", Type: "string", Description: "Report SK test_end_time", IsMustExist: true},
		{NameId: "calculated_test_duration_minute", Type: "float64p", Description: "Report SK calculated_test_duration_minute", IsMustExist: true},
		{NameId: "test_pressure", Type: "float64zp", Description: "Report SK test_pressure", IsMustExist: true},
		{NameId: "finished_date", Type: "date", Description: "Date of finished work", IsMustExist: true},
		{NameId: "gas_appliances", Type: "array-json-template", Description: "Report SK gas Appliance", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64p", Description: "Report, SK, Gas, Appliance gas_appliance_id", IsMustExist: true},
			{NameId: "quantity", Type: "int64zp", Description: "Report SK Gas Appliance quantity", IsMustExist: true},
		}},
	}},

	// penanganan piutang

	{NameId: "penanganan_piutang", Type: "json", Description: "Report Gas In", IsMustExist: false, Children: []api.DXAPIEndPointParameter{
		{NameId: "reason_type", Type: "string", Description: "Report Gas In sk_report", IsMustExist: true},
		{NameId: "reason_notes", Type: "string", Description: "Report Gas In sk_report", IsMustExist: true},
	}},
}

func EmailSend(templateData map[string]any, templateContentType, templateTitle string, templateBody string, smtpConfig map[string]any, email string) error {
	// Replace all template keywords in templateBody using templateData
	for key, value := range templateData {
		placeholder := fmt.Sprintf("<%s>", key)
		aValue := fmt.Sprint("%v", value)
		templateBody = strings.ReplaceAll(templateBody, placeholder, aValue)
		templateTitle = strings.ReplaceAll(templateTitle, placeholder, aValue)
	}

	var emailToFullname string
	var ok bool
	emailToFullname, ok = templateData["fullname"].(string)
	if !ok {
		emailToFullname = email
	}
	emailBody := templateBody
	emailTitle := templateTitle

	smtpServer, ok := smtpConfig["host"].(string)
	if !ok {
		return errors.New("SMTP_HOST_NOT_FOUND_IN_CONFIG")
	}
	smtpUsername, ok := smtpConfig["username"].(string)
	if !ok {
		return errors.New("SMTP_USERNAME_NOT_FOUND_IN_CONFIG")
	}
	smtpPassword, ok := smtpConfig["password"].(string)
	if !ok {
		return errors.New("SMTP_PASSWORD_NOT_FOUND_IN_CONFIG")
	}
	smtpPortAsFloat, ok := smtpConfig["port"].(int)
	if !ok {
		return errors.New("SMTP_PORT_NOT_FOUND_IN_CONFIG")
	}
	smtpPort := int(smtpPortAsFloat)
	smtpSenderEmail, ok := smtpConfig["sender_email"].(string)
	if !ok {
		return errors.New("SMTP_SENDER_EMAIL_NOT_FOUND_IN_CONFIG")
	}
	smtpSSL, ok := smtpConfig["ssl"].(bool)
	if !ok {
		return errors.New("SMTP_SSL_NOT_FOUND_IN_CONFIG")
	}

	d := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)
	d.SSL = smtpSSL
	s, err := d.Dial()
	if err != nil {
		return errors.New(fmt.Sprintf("SMTP_DIAL_ERROR:%v", err.Error()))
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpSenderEmail)
	m.SetAddressHeader("To", email, emailToFullname)
	// Use subject from templateData if available
	m.SetHeader("Subject", emailTitle)
	m.SetBody(templateContentType, emailBody)

	err = gomail.Send(s, m)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func SMSSend(recipient string, templateData map[string]any, templateBody string, smsConfig map[string]any) error {
	// Replace all template keywords in templateBody using templateData
	for key, value := range templateData {
		placeholder := fmt.Sprintf("<%s>", key)
		aValue := fmt.Sprint("%v", value)
		templateBody = strings.ReplaceAll(templateBody, placeholder, aValue)
	}

	var ok bool
	smsBody := templateBody

	smsEnabled, ok := smsConfig["enabled"].(bool)
	if !ok {
		return errors.New("SMS_TYPE_NOT_FOUND_IN_CONFIG")
	}
	if !smsEnabled {
		return errors.New("SMS_NOT_ENABLED")
	}

	smsType, ok := smsConfig["type"].(string)
	if !ok {
		return errors.New("SMS_TYPE_NOT_FOUND_IN_CONFIG")
	}
	if smsType != "sms-md-media" {
		return errors.New("SMS_TYPE_NOT_MD_MEDIA")
	}

	smsAddress, ok := smsConfig["address"].(string)
	if !ok {
		return errors.New("SMS_ADDRESS_NOT_FOUND_IN_CONFIG")
	}

	smsUsername, ok := smsConfig["username"].(string)
	if !ok {
		return errors.New("SMS_USERNAME_NOT_FOUND_IN_CONFIG")
	}
	smsPassword, ok := smsConfig["password"].(string)
	if !ok {
		return errors.New("SMS_PASSWORD_NOT_FOUND_IN_CONFIG")
	}
	smsSenderId, ok := smsConfig["sender_id"].(string)
	if !ok {
		return errors.New("SMS_SENDER_ID_NOT_FOUND_IN_CONFIG")
	}

	token := base64.StdEncoding.EncodeToString([]byte(smsUsername + ":" + smsPassword))

	// Prepare JSON data
	jsonData := map[string]string{
		"sender":  smsSenderId,
		"msisdn":  recipient,
		"message": smsBody,
	}

	// Convert to JSON
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return errors.Wrap(err, "error in json.Marshall")
	}

	// Create a new request
	req, err := http.NewRequest("POST", smsAddress, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return errors.Wrap(err, "error in  http.NewRequest")
	}

	// Set headers
	req.Header.Add("Authorization", "Basic "+token)
	req.Header.Add("Content-Type", "application/json")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error in send request")
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error in io.ReadAll response")
	}

	responseInJSON := utils.JSON{}
	err = json.Unmarshal(body, &responseInJSON)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error in json.Unmarshal response=%v", responseInJSON))
	}
	code, ok := responseInJSON["code"].(int)
	if !ok {
		return errors.Errorf("response body code not found in %v", responseInJSON)
	}
	if code == 1 {
		log.Log.Info("SEND_SMS_SUCCESS")
		return nil
	}
	message, ok := responseInJSON["message"].(string)
	if !ok {
		return errors.Errorf("response body message not found in %v", responseInJSON)
	}
	return errors.Errorf("SMS SENT FAIL:response body code is 0, message=%v", message)
}
