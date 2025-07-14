package relyon

import (
	"fmt"
	"github.com/donnyhardyanto/dxlib/configuration"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib/utils/http/client"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/construction_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"net/http"
	"strconv"
	"time"
)

func getConfiguration(l *log.DXLog) (c utils.JSON, err error) {
	configExternalSystem := *configuration.Manager.Configurations["external_system"].Data
	var ok bool
	c, ok = configExternalSystem["RELYON1"].(utils.JSON)
	if !ok {
		return nil, l.ErrorAndCreateErrorf("RELYON_CONFIG_NOT_FOUND")
	}
	return c, nil
}

func processResponse(l *log.DXLog, requestUrl string, response *client.HTTPResponse) (responseStatusCode int, responseData map[string]any, err error) {
	responseData, errBodyAsJSON := response.BodyAsJSON()
	info := ""
	responseStatusCode = response.StatusCode
	if responseStatusCode != http.StatusOK {
		if errBodyAsJSON == nil {
			info = fmt.Sprintf("INFO=%s", responseData["info"].(string))
		} else {
			info = response.BodyAsString()
		}
		return responseStatusCode, responseData, l.WarnAndCreateErrorf("RELYON_OUTBOUND_RESPONSE_STATUS_NOT_OK:%s=%d:%s", requestUrl, response.StatusCode, info)
	}

	if errBodyAsJSON != nil {
		return 0, responseData, l.WarnAndCreateErrorf("RELYON_OUTBOUND_RESPONSE_JSON_PARSE_ERROR:%s=%s", requestUrl, errBodyAsJSON.Error())
	}

	//info = responseData["info"].(string)
	//if info != "OK" {
	//	return 0, responseData, l.WarnAndCreateErrorf("RELYON_OUTBOUND_RESPONSE_200_NOT_SUCCESS:%s=%d:%s", requestUrl, response.StatusCode, info)
	//
	//}
	return responseStatusCode, responseData, nil
}

func Auth(l *log.DXLog) (responseStatusCode int, session string, err error) {
	relyOnConfig, err := getConfiguration(l)
	if err != nil {
		return 0, "", err
	}

	relyOnConfigAuthUrl, ok := relyOnConfig["outbound_auth_url"].(string)
	if !ok {
		return 0, "", errors.Errorf("RELYON_OUTBOUND_AUTH_URL_NOT_FOUND")
	}
	relyOnConfigAuthMethod, ok := relyOnConfig["outbound_auth_method"].(string)
	if !ok {
		return 0, "", errors.Errorf("RELYON_OUTBOUND_AUTH_METHOD_NOT_FOUND")
	}
	relyOnConfigHeaderAsString, ok := relyOnConfig["outbound_auth_header"].(string)
	if !ok {
		return 0, "", errors.Errorf("RELYON_OUTBOUND_AUTH_HEADER_NOT_FOUND")
	}
	relyOnConfigHeaderAsJSON, err := utils.StringToJSON(relyOnConfigHeaderAsString)
	if err != nil {
		return 0, "", err
	}
	relyOnConfigHeader := utils.JSONToMapStringString(relyOnConfigHeaderAsJSON)

	relyOnConfigRawBody, ok := relyOnConfig["outbound_auth_raw_body"].(string)
	if !ok {
		return 0, "", errors.Errorf("RELYON_OUTBOUND_AUTH_RAW_BODY_NOT_FOUND")
	}

	_, response, err := client.HTTPClientReadAll(relyOnConfigAuthMethod, relyOnConfigAuthUrl, relyOnConfigHeader, relyOnConfigRawBody)
	if err != nil {
		return 0, "", err
	}
	responseStatusCode, responseData, err := processResponse(l, relyOnConfigAuthUrl, response)
	if responseStatusCode != 200 {
		return responseStatusCode, "", l.WarnAndCreateErrorf("RELYON_OUTBOUND_AUTH_RESPONSE_STATUS_NOT_OK:%s=%d;%+v", relyOnConfigAuthUrl, responseStatusCode, responseData)
	}
	if err != nil {
		return 0, "", err
	}

	session, ok = responseData["data"].(map[string]interface{})["token"].(string)
	if !ok {
		return 0, "", l.WarnAndCreateErrorf("RELYON_OUTBOUND_AUTH_RESPONSE_ERROR_TOKEN_NOT_FOUND_IN_RESPONSE")
	}

	return responseStatusCode, session, nil
}

func RegisterInstallationUpdateType(l *log.DXLog, session string, registrationId string, installationTypeId int64) (responseStatusCode int, err error) {
	relyOnConfig, err := getConfiguration(l)
	if err != nil {
		return 0, err
	}
	relyOnConfigOutboundAPIURLRegisterInstallationUpdateType, ok := relyOnConfig["outbound_api_url_register_installation_update_type"].(string)
	if !ok {
		return 0, errors.Errorf("RELYON_OUTBOUND_API_URL_REGISTER_INSTALLATION_UPDATE_TYPE_NOT_FOUND")
	}
	_, response, err := client.HTTPClientReadAll("POST", relyOnConfigOutboundAPIURLRegisterInstallationUpdateType,
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", session)}, map[string]any{
			"data": map[string]any{
				"registration_id":      registrationId,
				"installation_type_id": installationTypeId,
			},
		})
	if err != nil {
		return 0, err
	}
	responseStatusCode, _, err = processResponse(l, relyOnConfigOutboundAPIURLRegisterInstallationUpdateType, response)
	if err != nil {
		return 0, err
	}

	return responseStatusCode, nil
}

func RegisterTaskInstallationUpdateStatus(l *log.DXLog, session string, task utils.JSON) (responseStatusCode int, r utils.JSON, err error) {
	taskId, ok := task["id"].(int64)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:TASK_ID_NOT_FOUND_IN_TASK:task_id=%s", taskId)
	}
	taskCustomerRegistrationNumber, ok := task["customer_registration_number"].(string)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:CUSTOMER_REGISTRATION_ID_NOT_FOUND_IN_TASK:task_id=%d", taskId)
	}
	taskStatus, ok := task["status"].(string)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:TASK_ID_NOT_FOUND_IN_TASK:task_id=%d", taskId)
	}
	if taskStatus != base.TaskStatusCompleted {
		return 0, nil, l.WarnAndCreateErrorf("CUSTOMER_TASK_TYPE_CONSTRUCTION_NOT_COMPLETED:task_id=%s", taskId)
	}

	_, subTaskReportSK, err := task_management.ModuleTaskManagement.SubTaskReport.SelectOne(l, nil, utils.JSON{
		"task_id":                 taskId,
		"sub_task_type_full_code": base.SubTaskTypeFullCodeConstructionSK,
		"is_deleted":              false,
	}, nil, map[string]string{"id": "desc"})
	if err != nil {
		return 0, nil, err
	}
	subTaskReportSKSubTaskId, ok := subTaskReportSK["sub_task_id"].(int64)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:SUB_TASK_ID_NOT_FOUND_IN_SUB_TASK_REPORT_SK:task_id=%d", taskId)
	}
	subTaskReportSKUserFullname, ok := subTaskReportSK["user_fullname"].(string)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:USER_FULLNAME_NOT_FOUND_IN_SUB_TASK_REPORT:sub_task_id=%d", subTaskReportSKSubTaskId)
	}
	subTaskReportSKTimestampAny, ok := subTaskReportSK["timestamp"]
	subTaskReportSKTimestamp, err := ConvertAnyToDateFormat(subTaskReportSKTimestampAny)

	//if err != nil {
	//	fmt.Println("Error:", err)
	//}
	fmt.Println(subTaskReportSKTimestampAny)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:TIMESTAMP_NOT_FOUND_IN_SUB_TASK_REPORT:sub_task_id=%d", subTaskReportSKSubTaskId)
	}
	_, subTaskReportSR, err := task_management.ModuleTaskManagement.SubTaskReport.SelectOne(l, nil, utils.JSON{
		"task_id":                 taskId,
		"sub_task_type_full_code": base.SubTaskTypeFullCodeConstructionSR,
		"is_deleted":              false,
	}, nil, map[string]string{"id": "desc"})
	if err != nil {
		return 0, nil, err
	}
	subTaskReportSRSubTaskId, ok := subTaskReportSR["sub_task_id"].(int64)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:SUB_TASK_ID_NOT_FOUND_IN_SUB_TASK_REPORT_SR:task_id=%d", taskId)
	}
	subTaskReportSRUserFullname, ok := subTaskReportSR["user_fullname"].(string)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:USER_FULLNAME_NOT_FOUND_IN_SUB_TASK_REPORT:sub_task_id=%d", subTaskReportSRSubTaskId)
	}
	subTaskReportSRTimestampAny, ok := subTaskReportSR["timestamp"]
	subTaskReportSRTimestamp, err := ConvertAnyToDateFormat(subTaskReportSRTimestampAny)

	//if err != nil {
	//	fmt.Println("Error:", err)
	//}

	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:TIMESTAMP_NOT_FOUND_IN_SUB_TASK_REPORT:sub_task_id=%d", subTaskReportSRSubTaskId)
	}
	_, subTaskReportPMG, err := task_management.ModuleTaskManagement.SubTaskReport.SelectOne(l, nil, utils.JSON{
		"task_id":                 taskId,
		"sub_task_type_full_code": base.SubTaskTypeFullCodeConstructionMeterInstallation,
		"is_deleted":              false,
	}, nil, map[string]string{"id": "desc"})
	if err != nil {
		return 0, nil, err
	}
	subTaskReportPMGSubTaskId, ok := subTaskReportPMG["sub_task_id"].(int64)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:SUB_TASK_ID_NOT_FOUND_IN_SUB_TASK_REPORT_PMG:task_id=%d", taskId)
	}
	subTaskReportPMGUserFullname, ok := subTaskReportPMG["user_fullname"].(string)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:USER_FULLNAME_NOT_FOUND_IN_SUB_TASK_REPORT:sub_task_id=%d", subTaskReportPMGSubTaskId)
	}
	subTaskReportPMGTimestampAny, ok := subTaskReportPMG["timestamp"]
	subTaskReportPMGTimestamp, err := ConvertAnyToDateFormat(subTaskReportPMGTimestampAny)

	//if err != nil {
	//	fmt.Println("Error:", err)
	//}

	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:TIMESTAMP_NOT_FOUND_IN_SUB_TASK_REPORT:sub_task_id=%d", subTaskReportPMGSubTaskId)
	}
	_, subTaskReportGasIn, err := task_management.ModuleTaskManagement.SubTaskReport.SelectOne(l, nil, utils.JSON{
		"task_id":                 taskId,
		"sub_task_type_full_code": base.SubTaskTypeFullCodeConstructionGasIn,
		"is_deleted":              false,
		"sub_task_status":         "WAITING_VERIFICATION",
	}, nil, map[string]string{"id": "desc"})
	if err != nil {
		return 0, nil, err
	}
	subTaskReportGasInSubTaskId, ok := subTaskReportGasIn["sub_task_id"].(int64)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:SUB_TASK_ID_NOT_FOUND_IN_SUB_TASK_REPORT_GasIn:task_id=%d", taskId)
	}
	subTaskReportGasInUserFullname, ok := subTaskReportGasIn["user_fullname"].(string)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:USER_FULLNAME_NOT_FOUND_IN_SUB_TASK_REPORT:sub_task_id=%d", subTaskReportGasInSubTaskId)
	}
	subTaskReportGasInFieldExecutorUserPhone, ok := subTaskReportGasIn["user_phonenumber"].(string)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:FIELD_EXECUTOR_USER_PHONENUMBER_NOT_FOUND_IN_SUB_TASK_REPORT:sub_task_id=%d", subTaskReportPMGSubTaskId)
	}
	subTaskReportGasInFieldExecutorOrganizationName, ok := subTaskReportGasIn["organization_name"].(string)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:FIELD_EXECUTOR_ORGANIZATION_NAME_NOT_FOUND_IN_SUB_TASK_REPORT:sub_task_id=%d", subTaskReportGasInSubTaskId)
	}
	subTaskReportGasInTimestampAny, ok := subTaskReportGasIn["timestamp"]
	subTaskReportGasInTimestamp, err := ConvertAnyToDateFormat(subTaskReportGasInTimestampAny)

	//if err != nil {
	//	fmt.Println("Error:", err)
	//}

	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:TIMESTAMP_NOT_FOUND_IN_SUB_TASK_REPORT:sub_task_id=%d", subTaskReportGasInSubTaskId)
	}
	subTaskReportGasInReport, ok := subTaskReportGasIn["report"].(utils.JSON)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:REPORT_NOT_FOUND_IN_SUB_TASK_REPORT:sub_task_id=%d", subTaskReportGasInSubTaskId)
	}
	subTaskReportGasInReportStandMeterStartNumber, ok := subTaskReportGasInReport["stand_meter_start_number"].(float64)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:STAND_METER_START_NUMBER_NOT_FOUND_IN_SUB_TASK_REPORT_REPORT:sub_task_id=%d", subTaskReportGasInSubTaskId)
	}
	subTaskReportGasInReportGasInDateAny, ok := subTaskReportGasInReport["gas_in_date"]

	subTaskReportGasInReportGasInDate, err := ConvertAnyToDateFormat(subTaskReportGasInReportGasInDateAny)

	//if err != nil {
	//	fmt.Println("Error:", err)
	//}

	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:GAS_IN_DATE_NOT_FOUND_IN_SUB_TASK_REPORT_REPORT:sub_task_id=%d", subTaskReportGasInSubTaskId)
	}
	subTaskReportGasInReportSNMeter, ok := subTaskReportGasInReport["sn_meter"].(string)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:SN_METER_NOT_FOUND_IN_SUB_TASK_REPORT_REPORT:sub_task_id=%d", subTaskReportGasInSubTaskId)
	}
	subTaskReportGasInReportMeterBrand, ok := subTaskReportGasInReport["meter_brand"].(string)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:METER_BRAND_NOT_FOUND_IN_SUB_TASK_REPORT_REPORT:sub_task_id=%d", subTaskReportGasInSubTaskId)
	}
	subTaskReportGasInReportPressureStartAny, ok := subTaskReportGasInReport["pressure_start"].(float64)
	subTaskReportGasInReportPressureStart := strconv.FormatFloat(subTaskReportGasInReportPressureStartAny, 'f', -1, 64)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:PRESSURE_START_NOT_FOUND_IN_SUB_TASK_REPORT_REPORT:sub_task_id=%d", subTaskReportGasInSubTaskId)
	}
	subTaskReportGasInReportTemperatureStart, ok := subTaskReportGasInReport["temperature_start"].(float64)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:TEMPERATURE_START_NOT_FOUND_IN_SUB_TASK_REPORT_REPORT:sub_task_id=%d", subTaskReportGasInSubTaskId)
	}

	subTaskReportGasInReportGSizeIdAny, ok := subTaskReportGasInReport["g_size_id"]
	subTaskReportGasInReportGSizeId, err := anyToInt64(subTaskReportGasInReportGSizeIdAny)
	//if err != nil {
	//	fmt.Printf("Error converting value: %v\n", err)
	//}
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:G_SIZE_ID_NOT_FOUND_IN_SUB_TASK_REPORT_REPORT:sub_task_id=%d", subTaskReportGasInSubTaskId)
	}
	_, gSize, err := construction_management.ModuleConstructionManagement.GSize.GetById(l, subTaskReportGasInReportGSizeId)
	if err != nil {
		return 0, nil, err
	}
	qMinFloat, ok := gSize["qmin"].(float64)
	qMin := strconv.FormatFloat(qMinFloat, 'f', -1, 64)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:QMIN_NOT_FOUND_IN_G_SIZE:sub_task_id=%d:g_size_id=%d", subTaskReportGasInSubTaskId, subTaskReportGasInReportGSizeId)
	}
	qMaxFloat, ok := gSize["qmax"].(float64)
	qMax := strconv.FormatFloat(qMaxFloat, 'f', -1, 64)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:QMAX_NOT_FOUND_IN_G_SIZE:sub_task_id=%d:g_size_id=%d", subTaskReportGasInSubTaskId, subTaskReportGasInReportGSizeId)
	}
	subTaskReportGasInReportMeterIdAny, ok := subTaskReportGasInReport["meter_id"]
	subTaskReportGasInReportMeterId, err := anyToInt64(subTaskReportGasInReportMeterIdAny)
	//if err != nil {
	//	fmt.Printf("Error converting value: %v\n", err)
	//}

	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:METER_ID_NOT_FOUND_IN_SUB_TASK_REPORT_REPORT:sub_task_id=%d", subTaskReportGasInSubTaskId)
	}
	_, meterApplianceType, err := construction_management.ModuleConstructionManagement.MeterApplianceType.GetById(l, subTaskReportGasInReportMeterId)
	if err != nil {
		return 0, nil, l.WarnAndCreateErrorf("ERROR:FAILED_TO_GET_METER_APPLIANCE_TYPE_BY_ID:sub_task_id=%d:meter_id=%d", subTaskReportGasInSubTaskId, subTaskReportGasInReportMeterId)
	}
	meterApplianceTypeName, ok := meterApplianceType["name"].(string)
	if !ok {
		return 0, nil, l.WarnAndCreateErrorf("IMPOSSIBLE:NAME_NOT_FOUND_IN_METER_APPLIANCE_TYPE:sub_task_id=%d:meter_id=%d", subTaskReportGasInSubTaskId, subTaskReportGasInReportMeterId)
	}

	relyOnConfig, err := getConfiguration(l)
	if err != nil {
		return 0, nil, err
	}
	relyOnConfigOutboundAPIURLRegisterInstallationUpdateStatus, ok := relyOnConfig["outbound_api_url_register_installation_update_status"].(string)
	if !ok {
		return 0, nil, errors.Errorf("RELYON_OUTBOUND_API_URL_REGISTER_INSTALLATION_UPDATE_STATUS_NOT_FOUND")
	}

	responseMessage := ""
	n := time.Now()

	_, response, err := client.HTTPClientReadAll("POST", relyOnConfigOutboundAPIURLRegisterInstallationUpdateStatus,
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", session)}, map[string]any{
			"data": map[string]any{
				"registration_id": taskCustomerRegistrationNumber,
				"remark":          "Instalasi selesai dilakukan",
				"detail": map[string]any{
					"sk_installation_by":        subTaskReportSKUserFullname,
					"sk_installation_date":      subTaskReportSKTimestamp,
					"sr_installation_by":        subTaskReportSRUserFullname,
					"sr_installation_date":      subTaskReportSRTimestamp,
					"meter_installation_by":     subTaskReportPMGUserFullname,
					"meter_installation_date":   subTaskReportPMGTimestamp,
					"gasin_by":                  subTaskReportGasInUserFullname,
					"gasin_finish_date":         subTaskReportGasInTimestamp,
					"gasin_initial_stand":       subTaskReportGasInReportStandMeterStartNumber,
					"gasin_pickup_date":         subTaskReportGasInReportGasInDate,
					"gasin_installatur_name":    subTaskReportGasInFieldExecutorOrganizationName,
					"gasin_meter_serial_number": subTaskReportGasInReportSNMeter,
					"gasin_by_phone":            subTaskReportGasInFieldExecutorUserPhone,
					"gasin_meter_merk":          subTaskReportGasInReportMeterBrand,
					"gasin_qmin":                qMin,
					"gasin_qmax":                qMax,
					"gasin_pressure":            subTaskReportGasInReportPressureStart,
					"gasin_temperature":         subTaskReportGasInReportTemperatureStart,
					"gasin_initial_calibration": "",
					"gasin_meter_type":          meterApplianceTypeName,
					"file_bagi_url":             ".",
				},
			},
		})

	responseStatusCode, _, err = processResponse(l, relyOnConfigOutboundAPIURLRegisterInstallationUpdateStatus, response)

	if responseStatusCode != 200 {
		responseMessage = http.StatusText(responseStatusCode)
	}

	if err != nil {
		responseMessage = err.Error()
	}

	p := utils.JSON{
		"last_relyon_sync_at":          n,
		"last_relyon_sync_status_code": responseStatusCode,
		"last_relyon_sync_message":     responseMessage,
	}

	r = utils.JSON{
		"sync_status":  responseStatusCode,
		"sync_message": responseMessage,
	}

	if responseStatusCode == 200 {
		p["last_relyon_sync_success_at"] = n
		r["last_sync_success_at"] = n
	}

	_, err = task_management.ModuleTaskManagement.Task.UpdateOne(l, taskId, p)
	if err != nil {
		return 0, r, err
	}
	return responseStatusCode, r, nil
}

func RegisterInstallationUpdateStatus(l *log.DXLog, session string, registrationId string) (responseStatusCode int, r utils.JSON, err error) {
	p := utils.JSON{
		"registration_number": registrationId,
	}
	_, customer, err := task_management.ModuleTaskManagement.Customer.SelectOne(l, nil, p, nil, nil)
	if err != nil {
		return 0, nil, err
	}
	if customer == nil {
		return 0, nil, l.WarnAndCreateErrorf("CUSTOMER_NOT_FOUND:%s", registrationId)
	}

	_, tasks, err := task_management.ModuleTaskManagement.Task.Select(l, nil, utils.JSON{
		"customer_id":  customer["id"].(int64),
		"task_type_id": base.TaskTypeIdConstruction,
	}, nil,
		map[string]string{"id": "desc"}, nil)
	if err != nil {
		return 0, nil, err
	}

	if tasks == nil {
		return 0, nil, l.WarnAndCreateErrorf("CUSTOMER_TASK_TYPE_CONSTRUCTION_NOT_FOUND:%s", registrationId)
	}

	if len(tasks) == 0 {
		return 0, nil, l.WarnAndCreateErrorf("CUSTOMER_TASK_TYPE_CONSTRUCTION_NOT_FOUND:%s", registrationId)
	}

	task := tasks[0]

	responseStatusCode, r, err = RegisterTaskInstallationUpdateStatus(l, session, task)
	if err != nil {
		return responseStatusCode, r, err
	}

	return responseStatusCode, r, nil
}

func CancelSubscriptionMileStone(l *log.DXLog, session string, registrationId string) (milestones any, err error) {
	relyOnConfig, err := getConfiguration(l)
	if err != nil {
		return nil, err
	}
	relyOnConfigOutboundAPIURLCancelSubscriptionMileStone, ok := relyOnConfig["outbound_api_url_cancel_subscription_milestone"].(string)
	if !ok {
		return nil, errors.Errorf("RELYON_OUTBOUND_API_URL_CANCEL_SUBSCRIPTION_MILESTONE_NOT_FOUND")
	}
	_, response, err := client.HTTPClientReadAll("POST", relyOnConfigOutboundAPIURLCancelSubscriptionMileStone,
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", session)}, map[string]any{
			"data": map[string]any{
				"registration_id": registrationId,
			},
		})
	fmt.Println("Response", response)
	if err != nil {
		return nil, err
	}
	_, responseData, err := processResponse(l, relyOnConfigOutboundAPIURLCancelSubscriptionMileStone, response)
	if err != nil {
		return nil, err
	}
	fmt.Println("Response Data", responseData)

	data, ok := responseData["data"].(map[string]any)
	if !ok {
		return nil, errors.Errorf("RELYON_OUTBOUND_API_RESPONSE_CANCEL_SUBSCRIPTION_MILESTONE_NO_DATA_PROPERTY")
	}
	fmt.Println("Data", data["milestones"])
	milestones = data["milestones"]
	//if !ok {
	//	return nil, errors.Errorf("RELYON_OUTBOUND_API_RESPONSE_CANCEL_SUBSCRIPTION_MILESTONE_NO_DATA_MILESTONES_PROPERTY2")
	//}
	return milestones, nil
}

func CancelSubscriptionCreate(l *log.DXLog, session string, registrationId string, reason []string, channel string, cancelBy string) (err error) {
	relyOnConfig, err := getConfiguration(l)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	relyOnConfigOutboundAPIURLCancelSubscriptionCreate, ok := relyOnConfig["outbound_api_url_cancel_subscription_create"].(string)
	if !ok {
		return errors.Errorf("RELYON_OUTBOUND_API_URL_CANCEL_SUBSCRIPTION_CREATE_NOT_FOUND")
	}
	_, response, err := client.HTTPClientReadAll("POST", relyOnConfigOutboundAPIURLCancelSubscriptionCreate,
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", session)}, map[string]any{
			"data": map[string]any{
				"registration_id": registrationId,
				"reason":          reason,
				"channel":         channel,
				"cancel_by":       cancelBy,
			},
		})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, _, err = processResponse(l, relyOnConfigOutboundAPIURLCancelSubscriptionCreate, response)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func anyToInt64(value any) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		if v > 9223372036854775807 {
			return 0, fmt.Errorf("uint64 value %d overflows int64", v)
		}
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", value)
	}
}

func ConvertAnyToDateFormat(value any) (string, error) {
	// Try to assert the value as string
	dateStr, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("value is not a string, got %T", value)
	}

	// Try to parse the string as RFC3339 format (which includes "T" and "Z")
	parsedTime, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// If that fails, try a few other common formats
		layouts := []string{
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}

		var parseErr error
		for _, layout := range layouts {
			parsedTime, parseErr = time.Parse(layout, dateStr)
			if parseErr == nil {
				break
			}
		}

		if parseErr != nil {
			return "", fmt.Errorf("failed to parse date string: %v", dateStr)
		}
	}

	// Format to the desired format
	return parsedTime.Format("2006-01-02"), nil
}
