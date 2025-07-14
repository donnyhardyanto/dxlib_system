package handler

import (
	"fmt"
	"net/http"

	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/database/protected/db"
	dbUtils "github.com/donnyhardyanto/dxlib/database2/utils"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/external/relyon"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/master_data"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
)

func TaskList(aepr *api.DXAPIEndPointRequest) (err error) {
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

	t := task_management.ModuleTaskManagement.Task
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

	for i, row := range list {
		customerId, ok := row["customer_id"].(int64)
		if ok {
			_, customer, err := task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, customerId)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			list[i]["customer"] = customer
		}
	}

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

func TaskListPerConstructionArea(aepr *api.DXAPIEndPointRequest) (err error) {

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

	t := task_management.ModuleTaskManagement.Task
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

	for i, row := range list {
		customerId, ok := row["customer_id"].(int64)
		if ok {
			_, customer, err := task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, customerId)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			list[i]["customer"] = customer
		}
	}

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

func TaskSearch(aepr *api.DXAPIEndPointRequest) (err error) {
	isTaskTypeIdExist, taskTypeId, err := aepr.GetParameterValueAsInt64("task_type_id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isSearchExist, search, err := aepr.GetParameterValueAsString("search", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isLocationCodeExist, locationCode, err := aepr.GetParameterValueAsString("location_code", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isAreaCodeExist, areaCode, err := aepr.GetParameterValueAsString("area_code", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, sort, err := aepr.GetParameterValueAsString("sort", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isLatitudeExist, latitude, err := aepr.GetParameterValueAsFloat64("latitude")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isLongitudeExist, longitude, err := aepr.GetParameterValueAsFloat64("longitude")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, rowPerPage, err := aepr.GetParameterValueAsInt64("row_per_page")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, pageIndex, err := aepr.GetParameterValueAsInt64("page_index")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	t := task_management.ModuleTaskManagement.Task

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

	conditions := ""

	isStatusExist := true
	status := base.TaskStatusWaitingAssignment

	if isStatusExist {
		if conditions != "" {
			conditions += " AND "
		}
		conditions += fmt.Sprintf("status = '%s'", status)
	}

	if isTaskTypeIdExist {
		if conditions != "" {
			conditions += " AND "
		}
		conditions += fmt.Sprintf("task_type_id = %d", taskTypeId)
	}

	if isSearchExist && search != "" {
		if conditions != "" {
			conditions += " AND "
		}
		conditions += fmt.Sprintf(
			"(customer_fullname ILIKE '%%%s%%' OR customer_number ILIKE '%%%s%%')",
			search, search,
		)
	}

	if isLocationCodeExist && locationCode != "" {
		if conditions != "" {
			conditions += " AND "
		}
		conditions += fmt.Sprintf(
			"(customer_address_province_location_code = '%s' OR customer_address_kabupaten_location_code = '%s' OR customer_address_kecamatan_location_code = '%s' OR customer_address_kelurahan_location_code = '%s')",
			locationCode, locationCode, locationCode, locationCode,
		)
	}

	if isAreaCodeExist && areaCode != "" {
		_, area, err := master_data.ModuleMasterData.Area.ShouldGetByNameId(&aepr.Log, areaCode)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		_, areaChildren, err := master_data.ModuleMasterData.Area.Select(&aepr.Log, nil, utils.JSON{
			"parent_group": area["type"].(string),
			"parent_value": area["name"].(string),
		}, nil, nil, nil)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}

		childrenCodes := []string{}
		for _, areaChild := range areaChildren {
			childrenCodes = append(childrenCodes, areaChild["code"].(string))
		}

		if conditions != "" {
			conditions += " AND "
		}

		c := fmt.Sprintf("(customer_sales_area_code = '%s')", areaCode)
		if len(childrenCodes) > 0 {
			c1 := dbUtils.SQLBuildWhereInClause("customer_sales_area_code", childrenCodes)
			c += fmt.Sprintf("OR (%s)", c1)
		}

		conditions += c
	}

	var orderBy string
	var args utils.JSON
	switch sort {
	case "LATEST":
		orderBy = "created_at DESC, last_modified_at DESC"
	case "CLOSEST_TO_TODAY":
		if !isLatitudeExist || !isLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		args = utils.JSON{
			"longitude": longitude,
			"latitude":  latitude,
		}
		orderBy = "customer_geom <-> ST_SetSRID(ST_MakePoint(:longitude, :latitude), 4326)"
	case "CLOSEST_TO_THE_LOCATION":
		if !isLatitudeExist || !isLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		args = utils.JSON{
			"longitude": longitude,
			"latitude":  latitude,
		}
		orderBy = "customer_geom <-> ST_SetSRID(ST_MakePoint(:longitude, :latitude), 4326)"
	default:
		orderBy = ""
	}

	rowsInfo, list, totalRows, totalPage, _, err := db.NamedQueryPaging(
		t.Database.Connection, t.FieldTypeMapping,
		"",
		rowPerPage,
		pageIndex,
		"*",
		t.ListViewNameId,
		conditions,
		"",
		orderBy,
		args,
	)
	if err != nil {
		aepr.Log.Errorf(err, "Error At paging table %s (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}

	for i, row := range list {
		customerId, ok := row["customer_id"].(int64)
		if ok {
			_, customer, err := task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, customerId)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			list[i]["customer"] = customer
		}
	}

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

func TaskMultiSearch(aepr *api.DXAPIEndPointRequest) (err error) {
	isSearchExist, search, err := aepr.GetParameterValueAsString("search", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isTaskTypeIdsExist, taskTypeIds, err := aepr.GetParameterValueAsArrayOfInt64("task_type_ids")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isLocationCodesExist, locationCodes, err := aepr.GetParameterValueAsArrayOfString("location_codes")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isAreaCodesExist, areaCodes, err := aepr.GetParameterValueAsArrayOfString("area_codes")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isStatusesExist, statuses, err := aepr.GetParameterValueAsArrayOfString("statuses")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isSpecialFlagSKPrimerExist, specialFlagSKPrimer, err := aepr.GetParameterValueAsBool("special_flag_sk_primer")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, sort, err := aepr.GetParameterValueAsString("sort", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isLatitudeExist, latitude, err := aepr.GetParameterValueAsFloat64("latitude")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isLongitudeExist, longitude, err := aepr.GetParameterValueAsFloat64("longitude")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, rowPerPage, err := aepr.GetParameterValueAsInt64("row_per_page")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, pageIndex, err := aepr.GetParameterValueAsInt64("page_index")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	t := task_management.ModuleTaskManagement.Task

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

	conditions := ""

	if isStatusesExist {
		if conditions != "" {
			conditions += " AND "
		}
		if len(statuses) > 0 {
			c := dbUtils.SQLBuildWhereInClause("status", statuses)
			conditions += c
		}
	}

	if isTaskTypeIdsExist {
		if conditions != "" {
			conditions += " AND "
		}
		if len(taskTypeIds) > 0 {
			c := dbUtils.SQLBuildWhereInClauseInt64("task_type_id", taskTypeIds)
			conditions += c
		}
	}

	if isSpecialFlagSKPrimerExist {
		if conditions != "" {
			conditions += " AND "
		}

		c := dbUtils.SQLBuildWhereInClauseBool("special_flag_sk_primer", []bool{specialFlagSKPrimer})
		conditions += c

	}

	if isSearchExist && search != "" {
		if conditions != "" {
			conditions += " AND "
		}
		conditions += fmt.Sprintf(
			"(customer_fullname ILIKE '%%%s%%' OR customer_number ILIKE '%%%s%%')",
			search, search,
		)
	}

	if isLocationCodesExist {
		if conditions != "" {
			conditions += " AND "
		}
		if len(locationCodes) == 0 {
			c1 := dbUtils.SQLBuildWhereInClause("customer_address_province_location_code", locationCodes)
			c2 := dbUtils.SQLBuildWhereInClause("customer_address_kabupaten_location_code", locationCodes)
			c3 := dbUtils.SQLBuildWhereInClause("customer_address_kecamatan_location_code", locationCodes)
			c4 := dbUtils.SQLBuildWhereInClause("customer_address_kelurahan_location_code", locationCodes)
			conditions += fmt.Sprintf("(%s OR %s OR %s OR %s)", c1, c2, c3, c4)
		}
	}

	if isAreaCodesExist {
		for _, areaCode := range areaCodes {
			_, area, err := master_data.ModuleMasterData.Area.ShouldGetByNameId(&aepr.Log, areaCode)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			_, areaChildren, err := master_data.ModuleMasterData.Area.Select(&aepr.Log, nil, utils.JSON{
				"parent_group": area["type"].(string),
				"parent_value": area["name"].(string),
			}, nil, nil, nil)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}

			var childrenCodes []string
			for _, areaChild := range areaChildren {
				childrenCodes = append(childrenCodes, areaChild["code"].(string))
			}

			if conditions != "" {
				conditions += " AND "
			}

			c := fmt.Sprintf("(customer_sales_area_code = '%s')", areaCode)
			if len(childrenCodes) > 0 {
				c1 := dbUtils.SQLBuildWhereInClause("customer_sales_area_code", childrenCodes)
				c += fmt.Sprintf("OR (%s)", c1)
			}

			conditions += c
		}
	}

	var orderBy string
	var args utils.JSON
	switch sort {
	case "LATEST":
		orderBy = "created_at DESC, last_modified_at DESC"
	case "CLOSEST_TODAY":
		if !isLatitudeExist || !isLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		args = utils.JSON{
			"longitude": longitude,
			"latitude":  latitude,
		}
		orderBy = "customer_geom <-> ST_SetSRID(ST_MakePoint(:longitude, :latitude), 4326)"
	case "CLOSEST_TO_THE_LOCATION":
		if !isLatitudeExist || !isLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		args = utils.JSON{
			"longitude": longitude,
			"latitude":  latitude,
		}
		orderBy = "customer_geom <-> ST_SetSRID(ST_MakePoint(:longitude, :latitude), 4326)"
	default:
		orderBy = ""
	}

	rowsInfo, list, totalRows, totalPage, _, err := db.NamedQueryPaging(
		t.Database.Connection, t.FieldTypeMapping,
		"",
		rowPerPage,
		pageIndex,
		"*",
		t.ListViewNameId,
		conditions,
		"",
		orderBy,
		args,
	)
	if err != nil {
		aepr.Log.Errorf(err, "Error At paging table %s (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}

	for i, row := range list {
		customerId, ok := row["customer_id"].(int64)
		if ok {
			_, customer, err := task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, customerId)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			list[i]["customer"] = customer
		}
	}

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

func TaskMultiSearchPerConstructionArea(aepr *api.DXAPIEndPointRequest) (err error) {
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

	conditions := ""

	s1 := dbUtils.SQLBuildWhereInClauseInt64("task_type_id", taskTypeIds)
	s2 := dbUtils.SQLBuildWhereInClause("area_code", areaCodes)

	if len(taskTypeIds) > 0 {
		if conditions != "" {
			conditions = fmt.Sprintf("(%s) and ", conditions)
		}
		conditions = conditions + s1
	}
	if len(areaCodes) > 0 {
		if conditions != "" {
			conditions = fmt.Sprintf("(%s) and ", conditions)
		}
		conditions = conditions + s2
	}

	isSearchExist, search, err := aepr.GetParameterValueAsString("search", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isTaskTypeIdsExist, taskTypeIds, err := aepr.GetParameterValueAsArrayOfInt64("task_type_ids")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isLocationCodesExist, locationCodes, err := aepr.GetParameterValueAsArrayOfString("location_codes")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isAreaCodesExist, areaCodes, err := aepr.GetParameterValueAsArrayOfString("area_codes")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isStatusesExist, statuses, err := aepr.GetParameterValueAsArrayOfString("statuses")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isSpecialFlagSKPrimerExist, specialFlagSKPrimer, err := aepr.GetParameterValueAsBool("special_flag_sk_primer")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, sort, err := aepr.GetParameterValueAsString("sort", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isLatitudeExist, latitude, err := aepr.GetParameterValueAsFloat64("latitude")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isLongitudeExist, longitude, err := aepr.GetParameterValueAsFloat64("longitude")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, rowPerPage, err := aepr.GetParameterValueAsInt64("row_per_page")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, pageIndex, err := aepr.GetParameterValueAsInt64("page_index")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	t := task_management.ModuleTaskManagement.Task

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

	if isStatusesExist {
		if conditions != "" {
			conditions += " AND "
		}
		if len(statuses) > 0 {
			c := dbUtils.SQLBuildWhereInClause("status", statuses)
			conditions += c
		}
	}

	if isTaskTypeIdsExist {
		if conditions != "" {
			conditions += " AND "
		}
		if len(taskTypeIds) > 0 {
			c := dbUtils.SQLBuildWhereInClauseInt64("task_type_id", taskTypeIds)
			conditions += c
		}
	}

	if isSpecialFlagSKPrimerExist {
		if conditions != "" {
			conditions += " AND "
		}

		c := dbUtils.SQLBuildWhereInClauseBool("special_flag_sk_primer", []bool{specialFlagSKPrimer})
		conditions += c

	}

	if isSearchExist && search != "" {
		if conditions != "" {
			conditions += " AND "
		}
		conditions += fmt.Sprintf(
			"(customer_fullname ILIKE '%%%s%%' OR customer_number ILIKE '%%%s%%')",
			search, search,
		)
	}

	if isLocationCodesExist {
		if conditions != "" {
			conditions += " AND "
		}
		if len(locationCodes) == 0 {
			c1 := dbUtils.SQLBuildWhereInClause("customer_address_province_location_code", locationCodes)
			c2 := dbUtils.SQLBuildWhereInClause("customer_address_kabupaten_location_code", locationCodes)
			c3 := dbUtils.SQLBuildWhereInClause("customer_address_kecamatan_location_code", locationCodes)
			c4 := dbUtils.SQLBuildWhereInClause("customer_address_kelurahan_location_code", locationCodes)
			conditions += fmt.Sprintf("(%s OR %s OR %s OR %s)", c1, c2, c3, c4)
		}
	}

	if isAreaCodesExist {
		for _, areaCode := range areaCodes {
			_, area, err := master_data.ModuleMasterData.Area.ShouldGetByNameId(&aepr.Log, areaCode)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			_, areaChildren, err := master_data.ModuleMasterData.Area.Select(&aepr.Log, nil, utils.JSON{
				"parent_group": area["type"].(string),
				"parent_value": area["name"].(string),
			}, nil, nil, nil)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}

			var childrenCodes []string
			for _, areaChild := range areaChildren {
				childrenCodes = append(childrenCodes, areaChild["code"].(string))
			}

			if conditions != "" {
				conditions += " AND "
			}

			c := fmt.Sprintf("(customer_sales_area_code = '%s')", areaCode)
			if len(childrenCodes) > 0 {
				c1 := dbUtils.SQLBuildWhereInClause("customer_sales_area_code", childrenCodes)
				c += fmt.Sprintf("OR (%s)", c1)
			}

			conditions += c
		}
	}

	var orderBy string
	var args utils.JSON
	switch sort {
	case "LATEST":
		orderBy = "created_at DESC, last_modified_at DESC"
	case "CLOSEST_TODAY":
		if !isLatitudeExist || !isLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		args = utils.JSON{
			"longitude": longitude,
			"latitude":  latitude,
		}
		orderBy = "customer_geom <-> ST_SetSRID(ST_MakePoint(:longitude, :latitude), 4326)"
	case "CLOSEST_TO_THE_LOCATION":
		if !isLatitudeExist || !isLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		args = utils.JSON{
			"longitude": longitude,
			"latitude":  latitude,
		}
		orderBy = "customer_geom <-> ST_SetSRID(ST_MakePoint(:longitude, :latitude), 4326)"
	default:
		orderBy = ""
	}

	rowsInfo, list, totalRows, totalPage, _, err := db.NamedQueryPaging(
		t.Database.Connection, t.FieldTypeMapping,
		"",
		rowPerPage,
		pageIndex,
		"*",
		t.ListViewNameId,
		conditions,
		"",
		orderBy,
		args,
	)
	if err != nil {
		aepr.Log.Errorf(err, "Error At paging table %s (%s) ", t.NameId, err.Error())
		return errors.Wrap(err, "error occured")
	}

	for i, row := range list {
		customerId, ok := row["customer_id"].(int64)
		if ok {
			_, customer, err := task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, customerId)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			list[i]["customer"] = customer
		}
	}

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

func TaskCreate(aepr *api.DXAPIEndPointRequest) (err error) {
	code := aepr.ParameterValues["code"].Value.(string)
	taskTypeId := aepr.ParameterValues["task_type_id"].Value.(int64)
	customerId := aepr.ParameterValues["customer_id"].Value.(int64)
	data1 := aepr.ParameterValues["data1"].Value.(string)
	data2 := aepr.ParameterValues["data2"].Value.(string)

	_, _, err = task_management.ModuleTaskManagement.TaskType.ShouldGetById(&aepr.Log, taskTypeId)
	if err != nil {
		aepr.Log.Errorf(err, "ERROR: TaskType.ShouldGetById: %s", err.Error())
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "TASK_TYPE_NOT_FOUND:%d", taskTypeId)
	}

	_, _, err = task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, customerId)
	if err != nil {
		aepr.Log.Errorf(err, "ERROR: Customer.ShouldGetById: %s", err.Error())
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "CUSTOMER_NOT_FOUND:%d", taskTypeId)
	}

	err = task_management.ModuleTaskManagement.ValidationTaskCodeShouldNotAlreadyExist(&aepr.Log, code)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	var taskId int64

	db := database.Manager.Databases[task_management.ModuleTaskManagement.DatabaseNameId]
	err = db.Tx(&aepr.Log, database.LevelReadCommitted, func(tx *database.DXDatabaseTx) (err2 error) {
		switch taskTypeId {
		case base.TaskTypeIdConstruction:
			taskId, err = task_management.ModuleTaskManagement.TaskTxCreateConstruction(tx, code, customerId, data1, data2)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
		case base.TaskTypeIdTechnicalSupport:
			return errors.Errorf("TASK_TYPE_ID_UNDER_DEVELOPMENT:%d", taskTypeId)
		case base.TaskTypeIdDebtManagement:
			return errors.Errorf("TASK_TYPE_ID_UNDER_DEVELOPMENT:%d", taskTypeId)
		default:
			return errors.Errorf("TASK_TYPE_ID_NOT_SUPPORTED:%d", taskTypeId)
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_ = task_management.ModuleTaskManagement.DoNotifyTaskCreate(&aepr.Log, taskId)

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"task_id": taskId,
	}},
	)
	return nil
}

func TaskConstructionCreate(aepr *api.DXAPIEndPointRequest) (err error) {
	code := aepr.ParameterValues["code"].Value.(string)
	customerId := aepr.ParameterValues["customer_id"].Value.(int64)
	data1 := aepr.ParameterValues["data1"].Value.(string)
	data2 := aepr.ParameterValues["data2"].Value.(string)

	// Hardcode task type to construction
	taskTypeId := base.TaskTypeIdConstruction

	_, _, err = task_management.ModuleTaskManagement.TaskType.ShouldGetById(&aepr.Log, taskTypeId)
	if err != nil {
		aepr.Log.Errorf(err, "ERROR: TaskType.ShouldGetById: %s", err.Error())
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "TASK_TYPE_NOT_FOUND:%d", taskTypeId)
	}

	_, _, err = task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, customerId)
	if err != nil {
		aepr.Log.Errorf(err, "ERROR: Customer.ShouldGetById: %s", err.Error())
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "CUSTOMER_NOT_FOUND:%d", customerId)
	}

	err = task_management.ModuleTaskManagement.ValidationTaskCodeShouldNotAlreadyExist(&aepr.Log, code)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	var taskId int64

	db := database.Manager.Databases[task_management.ModuleTaskManagement.DatabaseNameId]
	err = db.Tx(&aepr.Log, database.LevelReadCommitted, func(tx *database.DXDatabaseTx) (err2 error) {
		taskId, err = task_management.ModuleTaskManagement.TaskTxCreateConstruction(tx, code, customerId, data1, data2)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_ = task_management.ModuleTaskManagement.DoNotifyTaskCreate(&aepr.Log, taskId)

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"task_id": taskId,
	}},
	)
	return nil
}

func TaskRelyOnSyncByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	_, uid, err := aepr.GetParameterValueAsString("uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, task, err := task_management.ModuleTaskManagement.Task.ShouldGetByUid(&aepr.Log, uid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	taskTypeId, ok := task["task_type_id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:TASK_ID_NOT_FOUND")
	}
	if taskTypeId != 1 {
		return errors.Errorf("TASK_TYPE_IS_NOT_CONSTRUCTION")
	}
	taskStatus, ok := task["status"].(string)
	if !ok {
		return errors.Errorf("TASK_STATUS_IS_NOT_COMPLETED")
	}
	if taskStatus != "COMPLETED" {
		return errors.Errorf("TASK_STATUS_IS_NOT_COMPLETED")
	}
	taskCustomerRegistrationNumber, ok := task["customer_registration_number"].(string)
	if !ok {
		return errors.Errorf("TASK_CUSTOMER_REGISTRATION_NUMBER_IS_NOT_FOUND")
	}

	_, session, err := relyon.Auth(&aepr.Log)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	responseStatusCode, r, err := relyon.RegisterInstallationUpdateStatus(&aepr.Log, session, taskCustomerRegistrationNumber)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if responseStatusCode != 200 {
		return errors.Errorf("RELYON_REGISTER_INSTALLATION_UPDATE_STATUS_FAILED")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": r})

	return nil
}

func TaskRelyOnSyncAll(aepr *api.DXAPIEndPointRequest) (err error) {
	_, tasks, err := task_management.ModuleTaskManagement.Task.Select(&aepr.Log, nil, utils.JSON{
		"status":       "COMPLETED",
		"task_type_id": 1,
		"c1":           db.SQLExpression{Expression: "last_relyon_sync_success_at is null"},
	}, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	if len(tasks) == 0 {
		return nil
	}

	_, session, err := relyon.Auth(&aepr.Log)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	go func() {
		for _, task := range tasks {
			_, _, _ = relyon.RegisterTaskInstallationUpdateStatus(&aepr.Log, session, task)
		}
	}()
	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)

	return nil
}

func TaskRead(aepr *api.DXAPIEndPointRequest) (err error) {
	t := task_management.ModuleTaskManagement.Task
	_, id, err := aepr.GetParameterValueAsInt64(t.FieldNameForRowId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, d, err := t.ShouldGetById(&aepr.Log, id)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	customerId := d["customer_id"].(int64)
	_, customer, err := task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, customerId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	d["customer"] = customer

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{t.ResultObjectName: d}})

	return nil
}

func TaskReadByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	_, uid, err := aepr.GetParameterValueAsString("uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, task, err := task_management.ModuleTaskManagement.Task.ShouldGetByUid(&aepr.Log, uid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	taskId, ok := task["id"].(int64)
	if !ok {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "TASK_ID_NOT_FOUND")
	}

	userId := aepr.LocalData["user_id"].(int64)
	userOrganizationId := aepr.LocalData["organization_id"].(int64)

	/* start filtering user as field_executor or field_supervisor */

	taskTypeId, ok := task["task_type_id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:TASK[TASK_TYPE_ID]==NOT_FOUND")
	}
	taskCustomerAddressKelurahanLocationCode, ok := task["customer_address_kelurahan_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:TASK[CUSTOMER_ADDRESS_KELURAHAN_LOCATION_CODE]==NOT_FOUND")
	}
	taskCustomerAddressKecamatanLocationCode, ok := task["customer_address_kecamatan_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:TASK[CUSTOMER_ADDRESS_KECAMATAN_LOCATION_CODE]==NOT_FOUND")
	}
	taskCustomerAddressKabupatenLocationCode, ok := task["customer_address_kabupaten_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:TASK[CUSTOMER_ADDRESS_KABUPATEN_LOCATION_CODE]==NOT_FOUND")
	}
	taskCustomerAddressProvinceLocationCode, ok := task["customer_address_province_location_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:TASK[CUSTOMER_ADDRESS_PROVINCE_LOCATION_CODE]==NOT_FOUND")
	}
	taskCustomerSalesAreaCode, ok := task["customer_sales_area_code"].(string)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:TASK[CUSTOMER_SALES_AREA_CODE]==NOT_FOUND")
	}

	areaCodes, err := master_data.ModuleMasterData.AreaCodeExpandParentTreeUp(&aepr.Log, taskCustomerSalesAreaCode)
	if err != nil {
		return err
	}

	isAllowed := false

	_, subTasks, err := task_management.ModuleTaskManagement.SubTask.Select(&aepr.Log, nil, utils.JSON{
		"task_id": taskId,
		"c1":      db.SQLExpression{Expression: fmt.Sprintf("last_field_executor_user_id=%d", userId)},
	}, nil, map[string]string{"id": "asc"}, nil)
	if err != nil {
		return err
	}
	if len(subTasks) >= 0 {
		isAllowed = true
	}

	_, subTaskTypes, err := task_management.ModuleTaskManagement.SubTaskType.Select(&aepr.Log, nil, utils.JSON{
		"task_type_id": taskTypeId,
	}, nil, nil, nil)
	if err != nil {
		return err
	}

	taskSubTaskTypeIds, err := utils.GetMapValueFromArrayOfJSON[int64](subTaskTypes, "id")
	if err != nil {
		return err
	}

	if !isAllowed {
		isAllowed = partner_management.ModulePartnerManagement.FieldSupervisorDoRequestIsUserAndHasEffectiveExpertiseAreaLocation(aepr, userId, userOrganizationId, taskSubTaskTypeIds,
			[]string{
				taskCustomerAddressKelurahanLocationCode,
				taskCustomerAddressKecamatanLocationCode,
				taskCustomerAddressKabupatenLocationCode,
				taskCustomerAddressProvinceLocationCode,
			}, areaCodes)
		if !isAllowed {
			return nil
		}
	}

	/* end filtering user as field_executor or field_supervisor */

	customerId := task["customer_id"].(int64)
	_, customer, err := task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, customerId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	task["customer"] = customer

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
		"data": utils.JSON{
			task_management.ModuleTaskManagement.Task.ResultObjectName: task,
		},
	})

	return nil
}
