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
	"strconv"
	"strings"
	"time"
)

func SubTaskSchedule(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskUid, err := aepr.GetParameterValueAsString("sub_task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldGetByUid(&aepr.Log, subTaskUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	userId := aepr.LocalData["user_id"].(int64)
	_, err = ValidateFieldExecutorUserWithSubTask(aepr, userId, subTask)
	if err != nil {
		return err
	}
	subTaskId, ok := subTask["id"].(int64)
	if !ok {
		return errors.Errorf("subTask.id:NOT_INT64")
	}
	subTaskStatus := subTask["status"].(string)
	switch subTaskStatus {
	case
		base.SubTaskStatusBlockingDependency,
		base.SubTaskStatusWaitingAssignment,
		base.SubTaskStatusAssigned,
		base.SubTaskStatusWorking,
		base.SubTaskStatusPaused:
	default:
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "INVALID_SUB_TASK_STATUS_FOR_OPERATION:%s", subTaskStatus)
	}

	_, newValue, err := aepr.GetParameterValueAsJSON("new")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	startDate := newValue["start_date"].(time.Time)
	endDate := newValue["end_date"].(time.Time)

	_, err = task_management.ModuleTaskManagement.SubTask.UpdateOne(&aepr.Log, subTaskId, utils.JSON{
		"scheduled_start_date": startDate,
		"scheduled_end_date":   endDate,
	})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, subTask, err = task_management.ModuleTaskManagement.SubTask.ShouldGetById(&aepr.Log, subTaskId)
	if err != nil {

	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"sub_task": subTask,
	}})
	return nil
}

func SubTaskGroupByStatus(aepr *api.DXAPIEndPointRequest) (err error) {
	_, groupStatus, err := aepr.GetParameterValueAsString("group_status")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	userId := aepr.LocalData["user_id"].(int64)
	conditions := utils.JSON{
		"last_field_executor_user_id": userId,
	}
	switch groupStatus {
	case "ON_PROGRESS":
		conditions = utils.JSON{
			"last_field_executor_user_id": userId,
			"c1":                          db.SQLExpression{Expression: fmt.Sprintf("status IN ('%s', '%s', '%s')", base.SubTaskStatusAssigned, base.SubTaskStatusWorking, base.SubTaskStatusPaused)},
		}
	case "REVISION":
		conditions = utils.JSON{
			"last_field_executor_user_id": userId,
			"c1":                          db.SQLExpression{Expression: fmt.Sprintf("status IN ('%s', '%s')", base.SubTaskStatusVerificationFail, base.SubTaskStatusFixing)},
		}
	case "CANCELED":
		conditions = utils.JSON{
			"last_field_executor_user_id": userId,
			"c1":                          db.SQLExpression{Expression: fmt.Sprintf("status IN ('%s', '%s')", base.SubTaskStatusCanceledByFieldExecutor, base.SubTaskStatusCanceledByCustomer)},
		}
	case "DONE":
		conditions = utils.JSON{
			"last_field_executor_user_id": userId,
			"c1":                          db.SQLExpression{Expression: fmt.Sprintf("status IN ('%s', '%s')", base.SubTaskStatusWaitingVerification, base.SubTaskStatusVerificationSuccess)},

			//			"c1":                          db.SQLExpression{Expression: fmt.Sprintf("status IN ('%s', '%s', '%s')", SubTaskStatusWaitingVerification, SubTaskStatusVerificationSuccess, SubTaskStatusVerifying)},
		}
	}

	_, datas, err := task_management.ModuleTaskManagement.SubTask.Select(&aepr.Log, nil, conditions, nil, map[string]string{"id": "asc"}, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	for i, row := range datas {
		customerId, ok := row["customer_id"].(int64)
		if ok {
			_, customer, err := task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, customerId)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			datas[i]["customer"] = customer
		}
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
		"data": utils.JSON{"list": datas},
	})
	return nil
}

func SubTaskList(aepr *api.DXAPIEndPointRequest) (err error) {
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

	t := task_management.ModuleTaskManagement.SubTask
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

	data := utils.JSON{
		"data": utils.JSON{
			"list": utils.JSON{
				"rows":       list,
				"total_rows": totalRows,
				"total_page": totalPage,
				"rows_info":  rowsInfo,
			},
		},
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, data)
	return nil
}

func SubTaskListPerConstructionArea(aepr *api.DXAPIEndPointRequest) (err error) {
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

	t := task_management.ModuleTaskManagement.SubTask
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
		}},
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, data)
	return nil
}

func SubTaskSearch(aepr *api.DXAPIEndPointRequest) (err error) {
	_, search, err := aepr.GetParameterValueAsString("search", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, locationCode, err := aepr.GetParameterValueAsString("location_code", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, types, err := aepr.GetParameterValueAsArrayOfAny("types")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	strArr := make([]string, len(types))
	for i, v := range types {
		strArr[i] = strconv.Itoa(int(v.(float64)))
	}
	subTaskTypes := fmt.Sprintf("(%s)", strings.Join(strArr, ","))

	_, rowPerPage, err := aepr.GetParameterValueAsInt64("row_per_page")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, pageIndex, err := aepr.GetParameterValueAsInt64("page_index")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	t := task_management.ModuleTaskManagement.SubTask

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
	userId := aepr.LocalData["user_id"].(int64)

	conditions := fmt.Sprintf(
		"(customer_fullname ILIKE '%%%s%%' OR customer_number ILIKE '%%%s%%') AND (customer_address_province_location_code = '%s' OR customer_address_kabupaten_location_code = '%s' OR customer_address_kecamatan_location_code = '%s' OR customer_address_kelurahan_location_code = '%s') AND last_field_executor_user_id = '%d' AND sub_task_type_id IN %s",
		search,
		search,
		locationCode,
		locationCode,
		locationCode,
		locationCode,
		userId,
		subTaskTypes,
	)
	if locationCode == "" {

		conditions = fmt.Sprintf(
			"(customer_fullname ILIKE '%%%s%%' OR customer_number ILIKE '%%%s%%') AND last_field_executor_user_id = '%d' AND sub_task_type_id IN %s",
			search,
			search,
			userId,
			subTaskTypes,
		)
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
		"",
		utils.JSON{},
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

func SubTaskListByLastFieldExecutorSelfAndStatusesAssignedScheduledWorkingPausedFixingVerificationFailOrderByScheduledStartDateASC(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.LocalData["user_id"].(int64)
	filterWhere := fmt.Sprintf("last_field_executor_user_id = :user_id AND (status = 'ASSIGNED' OR status = 'SCHEDULED' OR status = 'WORKING' OR status = 'PAUSED' OR status = 'FIXING' OR status = 'VERIFICATION_FAIL')", userId)
	filterOrderBy := "scheduled_start_date ASC"

	_, rowPerPage, err := aepr.GetParameterValueAsInt64("row_per_page")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, pageIndex, err := aepr.GetParameterValueAsInt64("page_index")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	isDeletedIncluded := false

	t := task_management.ModuleTaskManagement.SubTask
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
		filterWhere, "", filterOrderBy, map[string]any{
			"user_id": userId,
		})
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

func SubTaskListByLastFieldExecutorSelfAndStatusesNotCanceledByFieldExecutorOrCustomerAndCustomerFullnameOrNumberLike(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.LocalData["user_id"].(int64)
	filterWhere := fmt.Sprintf("last_field_executor_user_id = :user_id AND (status != 'CANCELED_BY_FIELD_EXECUTOR' OR status != 'CANCELED_BY_CUSTOMER') AND (LOWER(customer_fullname) LIKE '%:customer_fullname%' OR LOWER(customer_number) LIKE '%last_field_executor_user_id = ${authUser.id} AND (status != 'CANCELED_BY_FIELD_EXECUTOR' OR status != 'CANCELED_BY_CUSTOMER') AND (LOWER(customer_fullname) LIKE '%${searchController.text.toLowerCase()}%' OR LOWER(customer_number) LIKE '%:customer_number%')", userId)
	filterOrderBy := "id asc"

	_, customerFullname, err := aepr.GetParameterValueAsString("customer_fullname", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, customerNumber, err := aepr.GetParameterValueAsString("customer_number", "")
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

	isDeletedIncluded := false

	t := task_management.ModuleTaskManagement.SubTask
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
		filterWhere, "", filterOrderBy, map[string]any{
			"user_id":           userId,
			"customer_fullname": customerFullname,
			"customer_number":   customerNumber,
		})
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

func SubTaskListByLastFieldExecutorSelfAndStatusesAssignedScheduledWorkingPausedFixingVerificationFailOrderByLastModifiedAtDesc(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.LocalData["user_id"].(int64)
	filterWhere := fmt.Sprintf("last_field_executor_user_id = :user_id AND (status = 'ASSIGNED' OR status = 'SCHEDULED' OR status = 'WORKING' OR status = 'PAUSED' OR status = 'FIXING' OR status = 'VERIFICATION_FAIL')", userId)
	filterOrderBy := "last_modified_at DESC"

	_, rowPerPage, err := aepr.GetParameterValueAsInt64("row_per_page")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, pageIndex, err := aepr.GetParameterValueAsInt64("page_index")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	isDeletedIncluded := false

	t := task_management.ModuleTaskManagement.SubTask
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
		filterWhere, "", filterOrderBy, map[string]any{
			"user_id": userId,
		})
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

func SubTaskListByTaskId(aepr *api.DXAPIEndPointRequest) (err error) {
	_, taskId, err := aepr.GetParameterValueAsInt64("task_id")

	filterWhere := fmt.Sprintf("task_id = :task_id", taskId)
	filterOrderBy := "scheduled_start_date ASC"

	_, rowPerPage, err := aepr.GetParameterValueAsInt64("row_per_page")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, pageIndex, err := aepr.GetParameterValueAsInt64("page_index")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	isDeletedIncluded := false

	t := task_management.ModuleTaskManagement.SubTask
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
		filterWhere, "", filterOrderBy, map[string]any{
			"task_id": taskId,
		})
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

func SubTaskRead(aepr *api.DXAPIEndPointRequest) (err error) {
	_, id, err := aepr.GetParameterValueAsInt64("id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	t := task_management.ModuleTaskManagement.SubTask
	_, subTask, err := t.ShouldGetById(&aepr.Log, id)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	customerId := subTask["customer_id"].(int64)

	_, customer, err := task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, customerId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTask["customer"] = customer
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{t.ResultObjectName: subTask}})
	return nil
}

func SubTaskReportPictureUpdate(aepr *api.DXAPIEndPointRequest) (err error) {
	subTaskReportId := aepr.ParameterValues["sub_task_report_id"].Value.(int64)
	subTaskReportFileGroupId := aepr.ParameterValues["sub_task_report_file_group_id"].Value.(int64)
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

	_, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.ShouldGetById(&aepr.Log, subTaskReportId)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_REPORT_NOT_FOUND:%s", err.Error())
	}
	subTaskReportUid, ok := subTaskReport["uid"].(string)
	if !ok {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_UID_NOT_FOUND")
	}

	_, subTaskReportFileGroup, err := task_management.ModuleTaskManagement.SubTaskReportFileGroup.ShouldGetById(&aepr.Log, subTaskReportFileGroupId)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_REPORT_FILE_GROUP_NOT_FOUND:%s", err.Error())
	}
	subTaskReportFileGroupUid, ok := subTaskReportFileGroup["uid"].(string)
	if !ok {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_FILE_GROUP_UID_NOT_FOUND")
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
		subTaskReportFileId, err2 = task_management.ModuleTaskManagement.SubTaskReportFile.TxInsert(dtx, utils.JSON{
			"sub_task_file_id":    subTaskFileId,
			"sub_task_file_uid":   subTaskFileUid,
			"sub_task_report_id":  subTaskReportId,
			"sub_task_report_uid": subTaskReportUid,
		})
		if err2 != nil {
			return err2
		}

		_, subTaskReportFile, err2 := task_management.ModuleTaskManagement.SubTaskReportFile.TxShouldGetById(dtx, subTaskReportFileId)
		if err2 != nil {
			return err2
		}
		subTaskReportFileUid, ok = subTaskReportFile["uid"].(string)
		if !ok {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_FILE_UID_NOT_FOUND")
		}

		filename := subTaskFileGetPath(taskId, subTaskId, subTaskFileId, subTaskReportFileGroupId)

		err2 = task_management.ModuleTaskManagement.SubTaskReportPicture.Update(aepr, filename, "")
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

func SubTaskReportPictureUpdateFileContentBase64(aepr *api.DXAPIEndPointRequest) (err error) {
	subTaskReportId := aepr.ParameterValues["sub_task_report_id"].Value.(int64)
	subTaskReportFileGroupId := aepr.ParameterValues["sub_task_report_file_group_id"].Value.(int64)
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

	_, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.ShouldGetById(&aepr.Log, subTaskReportId)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_REPORT_NOT_FOUND:%s", err.Error())
	}
	subTaskReportUid, ok := subTaskReport["uid"].(string)
	if !ok {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_UID_NOT_FOUND")
	}

	_, subTaskReportFileGroup, err := task_management.ModuleTaskManagement.SubTaskReportFileGroup.ShouldGetById(&aepr.Log, subTaskReportFileGroupId)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_REPORT_FILE_GROUP_NOT_FOUND:%s", err.Error())
	}
	subTaskReportFileGroupUid, ok := subTaskReportFileGroup["uid"].(string)
	if !ok {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_FILE_GROUP_UID_NOT_FOUND")
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
		subTaskReportFileId, err2 = task_management.ModuleTaskManagement.SubTaskReportFile.TxInsert(dtx, utils.JSON{
			"sub_task_file_id":    subTaskFileId,
			"sub_task_file_uid":   subTaskFileUid,
			"sub_task_report_id":  subTaskReportId,
			"sub_task_report_uid": subTaskReportUid,
		})
		if err2 != nil {
			return err2
		}

		_, subTaskReportFile, err2 := task_management.ModuleTaskManagement.SubTaskReportFile.TxShouldGetById(dtx, subTaskReportFileId)
		if err2 != nil {
			return err2
		}
		subTaskReportFileUid, ok = subTaskReportFile["uid"].(string)
		if !ok {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_FILE_UID_NOT_FOUND")
		}

		filename := subTaskFileGetPath(taskId, subTaskId, subTaskFileId, subTaskReportFileGroupId)

		_, fileContentBase64, err := aepr.GetParameterValueAsString("content_base64")
		if err != nil {
			return err
		}

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

func SubTaskReportPictureUpdateByUid(aepr *api.DXAPIEndPointRequest) (err error) {
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

		err2 = task_management.ModuleTaskManagement.SubTaskReportPicture.Update(aepr, filename, "")
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

func SubTaskReportPictureAssignExisting(aepr *api.DXAPIEndPointRequest) (err error) {
	subTaskReportId := aepr.ParameterValues["sub_task_report_id"].Value.(int64)
	subTaskFileId := aepr.ParameterValues["sub_task_file_id"].Value.(int64)

	_, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.ShouldGetById(&aepr.Log, subTaskReportId)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_REPORT_NOT_FOUND:%s", err.Error())
	}
	subTaskReportUid, ok := subTaskReport["uid"].(string)
	if !ok {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_UID_NOT_FOUND")
	}

	_, subTaskFile, err := task_management.ModuleTaskManagement.SubTaskFile.ShouldGetById(&aepr.Log, subTaskFileId)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_FILE_NOT_FOUND:%s", err.Error())
	}
	subTaskFileUid, ok := subTaskFile["uid"].(string)
	if !ok {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_FILE_UID_NOT_FOUND")
	}
	var subTaskReportFileId int64
	var subTaskReportFileUid string
	dbTaskDispatcher := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]
	err = dbTaskDispatcher.Tx(&aepr.Log, sql.LevelReadCommitted, func(dtx *database.DXDatabaseTx) (err2 error) {
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

func SubTaskReportPictureAssignExistingByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	subTaskReportUid := aepr.ParameterValues["sub_task_report_uid"].Value.(string)
	subTaskFileUid := aepr.ParameterValues["sub_task_file_uid"].Value.(string)

	_, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.ShouldGetByUid(&aepr.Log, subTaskReportUid)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusNotFound, "", "SUB_TASK_REPORT==NOT_FOUND:%s", err.Error())
	}
	subTaskReportId, ok := subTaskReport["id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[ID]==NOT_FOUND")
	}

	subTaskReportSubTaskTypeId, ok := subTaskReport["sub_task_type_id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[SUB_TASK_TYPE_ID]==NOT_FOUND")
	}

	subTaskReportLastFieldExecutorUserId, ok := subTaskReport["sub_task_last_field_executor_user_id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_REPORT[SUB_TASK_LAST_FIELD_EXECUTOR_USER_ID]==NOT_FOUND")
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
			[]int64{subTaskReportSubTaskTypeId}, []string{
				subTaskReportCustomerAddressKelurahanLocationCode,
				subTaskReportCustomerAddressKecamatanLocationCode,
				subTaskReportCustomerAddressKabupatenLocationCode,
				subTaskReportCustomerAddressProvinceLocationCode,
			}, areaCodes)
		if !isAllowed {
			return nil
		}
	}

	_, subTaskFile, err := task_management.ModuleTaskManagement.SubTaskFile.ShouldGetByUid(&aepr.Log, subTaskFileUid)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_FILE_NOT_FOUND:%s", err.Error())
	}
	subTaskFileId, ok := subTaskFile["id"].(int64)
	if !ok {
		return errors.Errorf("IMPOSSIBLE:SUB_TASK_FILE_ID_NOT_FOUND")
	}

	_, _, err = task_management.ModuleTaskManagement.SubTaskReport.ShouldGetById(&aepr.Log, subTaskReportId)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_REPORT_NOT_FOUND:%s", err.Error())
	}

	var subTaskReportFileId int64
	var subTaskReportFileUid string
	dbTaskDispatcher := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]
	err = dbTaskDispatcher.Tx(&aepr.Log, sql.LevelReadCommitted, func(dtx *database.DXDatabaseTx) (err2 error) {

		subTaskReportFileId, err2 = task_management.ModuleTaskManagement.SubTaskReportFile.TxInsert(dtx, utils.JSON{
			"sub_task_file_id":    subTaskFileId,
			"sub_task_file_uid":   subTaskFileUid,
			"sub_task_report_id":  subTaskReportId,
			"sub_task_report_uid": subTaskReportUid,
		})
		if err2 != nil {
			return err2
		}

		_, subTaskReportFile, err2 := task_management.ModuleTaskManagement.SubTaskReportFile.TxShouldGetById(dtx, subTaskReportFileId)
		if err2 != nil {
			return err2
		}
		subTaskReportFileUid, ok = subTaskReportFile["uid"].(string)
		if !ok {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_FILE_UID_NOT_FOUND")
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

func subTaskFileGetPath(taskId, subTaskId, subTaskFileId, subTaskReportFileGroupId int64) string {
	return fmt.Sprintf("%d/%d/%d/%d.png", taskId, subTaskId, subTaskReportFileGroupId, subTaskFileId)
}

func subTaskReportFileGetPath(aepr *api.DXAPIEndPointRequest, subTaskReportFileId int64) (filename string, err error) {
	_, subTaskReportFile, err := task_management.ModuleTaskManagement.SubTaskReportFile.ShouldGetById(&aepr.Log, subTaskReportFileId)
	if err != nil {
		return "", err
	}

	taskId, ok := subTaskReportFile["task_id"].(int64)
	if !ok {
		return "", aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "TASK_ID_NOT_FOUND")
	}
	subTaskId, ok := subTaskReportFile["sub_task_id"].(int64)
	if !ok {
		return "", aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_ID_NOT_FOUND")
	}
	subTaskReportFileGroupId, ok := subTaskReportFile["sub_task_report_file_group_id"].(int64)
	if !ok {
		return "", aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_FILE_GROUP_ID_NOT_FOUND")
	}
	subTaskFileId, ok := subTaskReportFile["sub_task_file_id"].(int64)
	if !ok {
		return "", aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_FILE_ID_NOT_FOUND")
	}
	filename = subTaskFileGetPath(taskId, subTaskId, subTaskFileId, subTaskReportFileGroupId)
	return filename, err
}

func SubTaskReportPictureDownloadSource(aepr *api.DXAPIEndPointRequest) (err error) {
	subTaskReportFileId := aepr.ParameterValues["id"].Value.(int64)

	filename, err := subTaskReportFileGetPath(aepr, subTaskReportFileId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	err = task_management.ModuleTaskManagement.SubTaskReportPicture.DownloadSource(aepr, filename)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func SubTaskReportPictureDownloadSourceByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	subTaskReportFileUid := aepr.ParameterValues["uid"].Value.(string)

	_, subTaskReportFile, err := task_management.ModuleTaskManagement.SubTaskReportFile.ShouldGetByUid(&aepr.Log, subTaskReportFileUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskReportFileId, ok := subTaskReportFile["id"].(int64)
	if !ok {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_FILE_ID_NOT_FOUND")
	}

	subTaskReportId, ok := subTaskReportFile["sub_task_report_id"].(int64)
	if !ok {
		if !ok {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_ID_NOT_FOUND")
		}
	}
	_, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.ShouldGetById(&aepr.Log, subTaskReportId)
	if err != nil {
		return errors.Wrap(err, "error occured")
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
			[]int64{subTaskReportSubTaskTypeId}, []string{
				subTaskReportCustomerAddressKelurahanLocationCode,
				subTaskReportCustomerAddressKecamatanLocationCode,
				subTaskReportCustomerAddressKabupatenLocationCode,
				subTaskReportCustomerAddressProvinceLocationCode,
			}, areaCodes)
		if !isAllowed {
			return nil
		}
	}

	filename, err := subTaskReportFileGetPath(aepr, subTaskReportFileId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	err = task_management.ModuleTaskManagement.SubTaskReportPicture.DownloadSource(aepr, filename)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func SubTaskReportPictureDownloadSmallByProcessedImageNameId(aepr *api.DXAPIEndPointRequest, processedImageNameId string) (err error) {
	subTaskReportFileId := aepr.ParameterValues["id"].Value.(int64)

	filename, err := subTaskReportFileGetPath(aepr, subTaskReportFileId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	err = task_management.ModuleTaskManagement.SubTaskReportPicture.DownloadProcessedImage(aepr, processedImageNameId, filename)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func SubTaskReportPictureDownloadSmall(aepr *api.DXAPIEndPointRequest) (err error) {
	return SubTaskReportPictureDownloadSmallByProcessedImageNameId(aepr, "small")
}

func SubTaskReportPictureDownloadMedium(aepr *api.DXAPIEndPointRequest) (err error) {
	return SubTaskReportPictureDownloadSmallByProcessedImageNameId(aepr, "medium")
}

func SubTaskReportPictureDownloadBig(aepr *api.DXAPIEndPointRequest) (err error) {
	return SubTaskReportPictureDownloadSmallByProcessedImageNameId(aepr, "big")
}
func SubTaskReportPictureDownloadSmallByProcessedImageNameIdByUid(aepr *api.DXAPIEndPointRequest, processedImageNameId string) (err error) {
	subTaskReportFileUId := aepr.ParameterValues["uid"].Value.(string)
	_, subTaskReportFile, err := task_management.ModuleTaskManagement.SubTaskReportFile.ShouldGetByUid(&aepr.Log, subTaskReportFileUId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	subTaskReportFileId, ok := subTaskReportFile["id"].(int64)
	if !ok {
		return errors.Errorf("SUB_TASK_REPORT_FILE_ID_NOT_FOUND")
	}

	subTaskReportId, ok := subTaskReportFile["sub_task_report_id"].(int64)
	if !ok {
		if !ok {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "SUB_TASK_REPORT_ID_NOT_FOUND")
		}
	}
	_, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.ShouldGetById(&aepr.Log, subTaskReportId)
	if err != nil {
		return errors.Wrap(err, "error occured")
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
			[]int64{subTaskReportSubTaskTypeId}, []string{
				subTaskReportCustomerAddressKelurahanLocationCode,
				subTaskReportCustomerAddressKecamatanLocationCode,
				subTaskReportCustomerAddressKabupatenLocationCode,
				subTaskReportCustomerAddressProvinceLocationCode,
			}, areaCodes)
		if !isAllowed {
			return nil
		}
	}

	filename, err := subTaskReportFileGetPath(aepr, subTaskReportFileId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	err = task_management.ModuleTaskManagement.SubTaskReportPicture.DownloadProcessedImage(aepr, processedImageNameId, filename)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func SubTaskReportPictureDownloadSmallByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	return SubTaskReportPictureDownloadSmallByProcessedImageNameIdByUid(aepr, "small")
}

func SubTaskReportPictureDownloadMediumByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	return SubTaskReportPictureDownloadSmallByProcessedImageNameIdByUid(aepr, "medium")
}

func SubTaskReportPictureDownloadBigByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	return SubTaskReportPictureDownloadSmallByProcessedImageNameIdByUid(aepr, "big")
}
