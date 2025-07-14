package handler

import (
	"database/sql"
	"fmt"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/database/protected/db"
	dbUtils "github.com/donnyhardyanto/dxlib/database2/utils"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/master_data"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"net/http"
)

func SubTaskReportList(aepr *api.DXAPIEndPointRequest) (err error) {
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

	_, rowPerPage, err := aepr.GetParameterValueAsInt64("row_per_page")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, pageIndex, err := aepr.GetParameterValueAsInt64("page_index")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, isDeletedIncluded, err := aepr.GetParameterValueAsBool("is_deleted", false)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	t := task_management.ModuleTaskManagement.SubTaskReport

	if !isDeletedIncluded {
		if filterWhere != "" {
			filterWhere = fmt.Sprintf("(%s) and ", filterWhere)
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

	rowsInfo, list, totalRows, totalPage, _, err := db.NamedQueryPaging(t.Database.Connection, t.FieldTypeMapping, "", rowPerPage, pageIndex, "*", t.ListViewNameId,
		filterWhere, "", filterOrderBy, filterKeyValues)
	if err != nil {
		aepr.Log.Errorf(err, "Error At paging table %s (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}

	/*	for i := range list {
		report, err := utils.GetJSONFromKV(list[i], "report")
		if err != nil {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "JSON_FIELD_REPORT_MARSHALL_ERROR:%s", err.Error())
		}
		list[i]["report"] = report
	}*/

	data := utils.JSON{"data": utils.JSON{
		"list": utils.JSON{
			"rows":       list,
			"total_rows": totalRows,
			"total_page": totalPage,
			"rows_info":  rowsInfo,
		},
	}}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, data)

	return nil
}

func SubTaskReportListPerConstructionArea(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.LocalData["user_id"].(int64)

	_, userRoleMemberships, err := partner_management.ModulePartnerManagement.UserRoleMembership.Select(&aepr.Log, nil, utils.JSON{
		"user_id": userId,
	}, nil, map[string]string{"id": "asc"}, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	taskTypeIds := []int64{}
	areaCodes := []string{}

	for _, userRoleMembership := range userRoleMemberships {
		task_type_id, ok := userRoleMembership["task_type_id"].(int64)
		if ok {
			taskTypeIds = append(taskTypeIds, task_type_id)
		}
		area_code, ok := userRoleMembership["area_code"].(string)
		if ok {
			areaCodes = append(areaCodes, area_code)
		}
	}

	isExistFilterWhere, filterWhere, err := aepr.GetParameterValueAsString("filter_where")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if !isExistFilterWhere {
		filterWhere = ""
	}

	s1 := dbUtils.SQLBuildWhereInClauseInt64("task_type_id", taskTypeIds)
	s2 := dbUtils.SQLBuildWhereInClause("area_code", areaCodes)

	if len(taskTypeIds) > 0 {
		if filterWhere != "" {
			filterWhere = fmt.Sprintf("(%s) and ", filterWhere)
		}
		filterWhere = filterWhere + s1
	}
	if len(areaCodes) > 0 {
		if filterWhere != "" {
			filterWhere = fmt.Sprintf("(%s) and ", filterWhere)
		}
		filterWhere = filterWhere + s2
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

	_, rowPerPage, err := aepr.GetParameterValueAsInt64("row_per_page")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, pageIndex, err := aepr.GetParameterValueAsInt64("page_index")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, isDeletedIncluded, err := aepr.GetParameterValueAsBool("is_deleted", false)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	t := task_management.ModuleTaskManagement.SubTaskReport

	if !isDeletedIncluded {
		if filterWhere != "" {
			filterWhere = fmt.Sprintf("(%s) and ", filterWhere)
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

	rowsInfo, list, totalRows, totalPage, _, err := db.NamedQueryPaging(t.Database.Connection, t.FieldTypeMapping, "", rowPerPage, pageIndex, "*", t.ListViewNameId,
		filterWhere, "", filterOrderBy, filterKeyValues)
	if err != nil {
		aepr.Log.Errorf(err, "Error At paging table %s (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}

	/*	for i := range list {
			report, err := utils.GetJSONFromKV(list[i], "report")
			if err != nil {
				return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "JSON_FIELD_REPORT_MARSHALL_ERROR:%s", err.Error())
			}
			list[i]["report"] = report
		}
	*/
	data := utils.JSON{"data": utils.JSON{
		"list": utils.JSON{
			"rows":       list,
			"total_rows": totalRows,
			"total_page": totalPage,
			"rows_info":  rowsInfo,
		},
	}}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, data)

	return nil
}

func SubTaskReportReadByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	_, uid, err := aepr.GetParameterValueAsString("uid")
	if err != nil {
		return err
	}
	rowsInfo, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.ShouldGetByUid(&aepr.Log, uid)
	if err != nil {
		return errors.Wrapf(err, "PARAMETER_VALUE_ERROR | SubTaskReportReadByUid | SUB_TASK_REPORT[UID=%s]==NOT_FOUND", uid)
	}

	subTaskReportLastFieldExecutorUserId, ok := subTaskReport["sub_task_last_field_executor_user_id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_LAST_FIELD_EXECUTOR_USER_ID_NOT_FOUND")
	}

	subTaskReportSubTaskTypeId, ok := subTaskReport["sub_task_type_id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[SUB_TASK_TYPE_ID]==NOT_FOUND")
	}
	subTaskReportCustomerAddressKelurahanLocationCode, ok := subTaskReport["customer_address_kelurahan_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[CUSTOMER_ADDRESS_KELURAHAN_LOCATION_CODE]==NOT_FOUND")
	}
	subTaskReportCustomerAddressKecamatanLocationCode, ok := subTaskReport["customer_address_kecamatan_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[CUSTOMER_ADDRESS_KECAMATAN_LOCATION_CODE]==NOT_FOUND")
	}
	subTaskReportCustomerAddressKabupatenLocationCode, ok := subTaskReport["customer_address_kabupaten_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[CUSTOMER_ADDRESS_KABUPATEN_LOCATION_CODE]==NOT_FOUND")
	}
	subTaskReportCustomerAddressProvinceLocationCode, ok := subTaskReport["customer_address_province_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[CUSTOMER_ADDRESS_PROVINCE_LOCATION_CODE]==NOT_FOUND")
	}
	subTaskReportCustomerSalesAreaCode, ok := subTaskReport["customer_sales_area_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[CUSTOMER_SALES_AREA_CODE]==NOT_FOUND")
	}

	areaCodes, err := master_data.ModuleMasterData.AreaCodeExpandParentTreeUp(&aepr.Log, subTaskReportCustomerSalesAreaCode)
	if err != nil {
		return err
	}

	userId := aepr.LocalData["user_id"].(int64)
	userOrganizationId := aepr.LocalData["organization_id"].(int64)
	isAllowed := false

	if userId == subTaskReportLastFieldExecutorUserId {
		isAllowed = true
	}

	if !isAllowed {
		isAllowed = partner_management.ModulePartnerManagement.FieldSupervisorDoRequestIsUserAndHasEffectiveExpertiseAreaLocation(aepr, userId, userOrganizationId,
			[]int64{subTaskReportSubTaskTypeId},
			[]string{
				subTaskReportCustomerAddressKelurahanLocationCode,
				subTaskReportCustomerAddressKecamatanLocationCode,
				subTaskReportCustomerAddressKabupatenLocationCode,
				subTaskReportCustomerAddressProvinceLocationCode,
			}, areaCodes)
		if !isAllowed {
			return nil
		}
	}
	t := task_management.ModuleTaskManagement.SubTaskReport

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{t.ResultObjectName: subTaskReport, "rows_info": rowsInfo}})

	return nil
}

func SubTaskReportRead(aepr *api.DXAPIEndPointRequest) (err error) {
	_, id, err := aepr.GetParameterValueAsInt64("id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	rowsInfo, report, err := task_management.ModuleTaskManagement.SubTaskReport.GetById(&aepr.Log, id)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	t := task_management.ModuleTaskManagement.SubTaskReport

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{t.ResultObjectName: report, "rows_info": rowsInfo}})

	return nil
}

func SubTaskReportPictureListBySubTaskReportId(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskReportId, err := aepr.GetParameterValueAsInt64("sub_task_report_id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTaskReportFiles, err := task_management.ModuleTaskManagement.SubTaskReportFile.Select(&aepr.Log, nil, utils.JSON{
		"sub_task_report_id": subTaskReportId,
	}, nil, map[string]string{"id": "asc"}, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{"list": subTaskReportFiles}})
	return nil
}

func SubTaskReportPictureListBySubTaskReportUid(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskReportUid, err := aepr.GetParameterValueAsString("sub_task_report_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.ShouldGetByUid(&aepr.Log, subTaskReportUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskReportId, ok := subTaskReport["id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT.ID:NOT_INT64")
	}

	subTaskReportLastFieldExecutorUserId, ok := subTaskReport["sub_task_last_field_executor_user_id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_LAST_FIELD_EXECUTOR_USER_ID_NOT_FOUND")
	}

	subTaskReportSubTaskTypeId, ok := subTaskReport["sub_task_type_id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[SUB_TASK_TYPE_ID]==NOT_FOUND")
	}

	subTaskReportCustomerAddressKelurahanLocationCode, ok := subTaskReport["customer_address_kelurahan_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[CUSTOMER_ADDRESS_KELURAHAN_LOCATION_CODE]==NOT_FOUND")
	}
	subTaskReportCustomerAddressKecamatanLocationCode, ok := subTaskReport["customer_address_kecamatan_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[CUSTOMER_ADDRESS_KECAMATAN_LOCATION_CODE]==NOT_FOUND")
	}
	subTaskReportCustomerAddressKabupatenLocationCode, ok := subTaskReport["customer_address_kabupaten_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[CUSTOMER_ADDRESS_KABUPATEN_LOCATION_CODE]==NOT_FOUND")
	}
	subTaskReportCustomerAddressProvinceLocationCode, ok := subTaskReport["customer_address_province_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[CUSTOMER_ADDRESS_PROVINCE_LOCATION_CODE]==NOT_FOUND")
	}
	subTaskReportCustomerSalesAreaCode, ok := subTaskReport["customer_sales_area_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[CUSTOMER_SALES_AREA_CODE]==NOT_FOUND")
	}

	areaCodes, err := master_data.ModuleMasterData.AreaCodeExpandParentTreeUp(&aepr.Log, subTaskReportCustomerSalesAreaCode)
	if err != nil {
		return err
	}

	userId := aepr.LocalData["user_id"].(int64)
	userOrganizationId := aepr.LocalData["organization_id"].(int64)
	isAllowed := false

	if userId == subTaskReportLastFieldExecutorUserId {
		isAllowed = true
	}

	if !isAllowed {
		isAllowed = partner_management.ModulePartnerManagement.FieldSupervisorDoRequestIsUserAndHasEffectiveExpertiseAreaLocation(aepr, userId, userOrganizationId,
			[]int64{subTaskReportSubTaskTypeId},
			[]string{
				subTaskReportCustomerAddressKelurahanLocationCode,
				subTaskReportCustomerAddressKecamatanLocationCode,
				subTaskReportCustomerAddressKabupatenLocationCode,
				subTaskReportCustomerAddressProvinceLocationCode,
			}, areaCodes)
		if !isAllowed {
			return nil
		}
	}
	_, subTaskReportFiles, err := task_management.ModuleTaskManagement.SubTaskReportFile.Select(&aepr.Log, nil, utils.JSON{
		"sub_task_report_id": subTaskReportId,
	}, nil, map[string]string{"id": "asc"}, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{"list": subTaskReportFiles}})
	return nil
}

func SubTaskReportPictureUpdateByUidFileContentBase64(aepr *api.DXAPIEndPointRequest) (err error) {
	subTaskReportByUid := aepr.ParameterValues["sub_task_report_uid"].Value.(string)
	subTaskReportFileGroupUid := aepr.ParameterValues["sub_task_report_file_group_uid"].Value.(string)

	isAtExist, at, err := aepr.GetParameterValueAsTime("at")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isLongitudeExist, longitude, err := aepr.GetParameterValueAsFloat64("longitude")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isLatitudeExist, latitude, err := aepr.GetParameterValueAsFloat64("latitude")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, fileContentBase64, err := aepr.GetParameterValueAsString("content_base64")
	if err != nil {
		return errors.Wrap(err, "error getting content_base64 parameter")
	}

	_, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.ShouldGetByUid(&aepr.Log, subTaskReportByUid)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_REPORT_NOT_FOUND:%s", err.Error())
	}
	subTaskReportId, ok := subTaskReport["id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT_ID_NOT_FOUND")
	}
	subTaskReportUid, ok := subTaskReport["uid"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT_UID_NOT_FOUND")
	}

	userId := aepr.LocalData["user_id"].(int64)
	isAllowed := false

	subTaskReportLastFieldExecutorUserId, ok := subTaskReport["sub_task_last_field_executor_user_id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_LAST_FIELD_EXECUTOR_USER_ID_NOT_FOUND")
	}
	if userId == subTaskReportLastFieldExecutorUserId {
		isAllowed = true
	}

	if !isAllowed {
		return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "USER_NOT_LAST_FIELD_EXECUTOR")
	}

	_, subTaskReportFileGroup, err := task_management.ModuleTaskManagement.SubTaskReportFileGroup.ShouldGetByUid(&aepr.Log, subTaskReportFileGroupUid)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_REPORT_FILE_NOT_FOUND:%s", err.Error())
	}
	subTaskReportFileGroupId, ok := subTaskReportFileGroup["id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT_FILE_GROUP_ID_NOT_FOUND")
	}

	_, _, err = task_management.ModuleTaskManagement.SubTaskReportFileGroup.ShouldGetById(&aepr.Log, subTaskReportFileGroupId)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_REPORT_FILE_GROUP_NOT_FOUND:%s", err.Error())
	}

	subTaskId, ok := subTaskReport["sub_task_id"].(int64)
	if !ok {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_ID_NOT_FOUND")
	}
	subTaskUid, ok := subTaskReport["sub_task_uid"].(string)
	if !ok {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_UID_NOT_FOUND")
	}
	taskId, ok := subTaskReport["task_id"].(int64)
	if !ok {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "TASK_ID_NOT_FOUND")
	}

	p := utils.JSON{
		"sub_task_id":                    subTaskId,
		"sub_task_uid":                   subTaskUid,
		"sub_task_report_file_group_id":  subTaskReportFileGroupId,
		"sub_task_report_file_group_uid": subTaskReportFileGroupUid,
	}
	if isAtExist {
		p["at"] = at
	}
	if isLongitudeExist {
		p["longitude"] = longitude
	}
	if isLatitudeExist {
		p["latitude"] = latitude
	}

	var subTaskFile utils.JSON
	var subTaskFileId int64
	var subTaskFileUid string
	var subTaskReportFileId int64
	var subTaskReportFileUid string

	dbTaskDispatcher := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]
	err = dbTaskDispatcher.Tx(&aepr.Log, sql.LevelReadCommitted, func(dtx *database.DXDatabaseTx) (err2 error) {
		_, subTaskFile, err2 = task_management.ModuleTaskManagement.SubTaskFile.TxSelectOne(dtx, p, nil)
		if err2 != nil {
			return err2
		}
		if subTaskFile == nil {
			subTaskFileId, err2 = task_management.ModuleTaskManagement.SubTaskFile.TxInsert(dtx, p)
			if err2 != nil {
				return err2
			}
			_, subTaskFile, err2 = task_management.ModuleTaskManagement.SubTaskFile.TxShouldGetById(dtx, subTaskFileId)
			if err2 != nil {
				return err2
			}
		} else {
			subTaskFileId = subTaskFile["id"].(int64)
		}

		subTaskFileUid, ok := subTaskFile["uid"].(string)
		if !ok {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_FILE_UID_NOT_FOUND")
		}

		// Check if a SubTaskReportFile already exists with the same sub_task_report_id and sub_task_file_id
		_, existingSubTaskReportFile, err2 := task_management.ModuleTaskManagement.SubTaskReportFile.TxSelectOne(dtx, utils.JSON{
			"sub_task_file_id":   subTaskFileId,
			"sub_task_report_id": subTaskReportId,
		}, nil)
		if err2 != nil {
			return err2
		}

		var subTaskReportFile utils.JSON
		if existingSubTaskReportFile != nil {
			// Use the existing SubTaskReportFile
			subTaskReportFile = existingSubTaskReportFile
			subTaskReportFileId = subTaskReportFile["id"].(int64)
		} else {
			// Create a new SubTaskReportFile
			subTaskReportFileId, err2 = task_management.ModuleTaskManagement.SubTaskReportFile.TxInsert(dtx, utils.JSON{
				"sub_task_file_id":    subTaskFileId,
				"sub_task_file_uid":   subTaskFileUid,
				"sub_task_report_id":  subTaskReportId,
				"sub_task_report_uid": subTaskReportUid,
			})
			if err2 != nil {
				return err2
			}

			_, subTaskReportFile, err2 = task_management.ModuleTaskManagement.SubTaskReportFile.TxShouldGetById(dtx, subTaskReportFileId)
			if err2 != nil {
				return err2
			}
		}

		subTaskReportFileUid, ok = subTaskReportFile["uid"].(string)
		if !ok {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_FILE_UID_NOT_FOUND")
		}
		filename := subTaskFileGetPath(taskId, subTaskId, subTaskFileId, subTaskReportFileGroupId)

		// Pass the base64 content to Update instead of using request body
		err2 = task_management.ModuleTaskManagement.SubTaskReportPicture.Update(aepr, filename, fileContentBase64)
		if err2 != nil {
			return err2
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"sub_task_file_id":         subTaskFileId,
		"sub_task_file_uid":        subTaskFileUid,
		"sub_task_report_file_id":  subTaskReportFileId,
		"sub_task_report_file_uid": subTaskReportFileUid,
	}})
	return nil
}
