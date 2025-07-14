package handler

import (
	"fmt"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/database/protected/db"
	dbUtils "github.com/donnyhardyanto/dxlib/database2/utils"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/master_data"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"net/http"
	"strings"
)

// Note: fieldExpertiseId is userRoleMembershipId

func FieldExecutorGetEffectiveTaskTypeIds(log *log.DXLog, fieldExpertiseId int64) (effectiveExpertises []int64, err error) {
	_, fieldExecutorEffectiveExpertises, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveExpertise.Select(log, nil, utils.JSON{
		"user_role_membership_id": fieldExpertiseId,
	}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	effectiveExpertises, err = utils.GetMapValueFromArrayOfJSON[int64](fieldExecutorEffectiveExpertises, "task_type_id")
	if err != nil {
		return nil, err
	}
	effectiveExpertises = utils.RemoveDuplicates(effectiveExpertises)

	return effectiveExpertises, nil
}

func FieldSupervisorGetEffectiveTaskTypeIds(log *log.DXLog, fieldExpertiseId int64) (effectiveExpertises []int64, err error) {
	_, fieldSupervisorEffectiveExpertises, err := partner_management.ModulePartnerManagement.FieldSupervisorEffectiveExpertise.Select(log, nil, utils.JSON{
		"user_role_membership_id": fieldExpertiseId,
	}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	effectiveExpertises, err = utils.GetMapValueFromArrayOfJSON[int64](fieldSupervisorEffectiveExpertises, "task_type_id")
	if err != nil {
		return nil, err
	}
	effectiveExpertises = utils.RemoveDuplicates(effectiveExpertises)

	return effectiveExpertises, nil
}

func FieldExecutorGetEffectiveExpertise(log *log.DXLog, fieldExpertiseId int64) (effectiveExpertises []int64, err error) {
	_, fieldExecutorEffectiveExpertises, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveExpertise.Select(log, nil, utils.JSON{
		"user_role_membership_id": fieldExpertiseId,
	}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	effectiveExpertises, err = utils.GetMapValueFromArrayOfJSON[int64](fieldExecutorEffectiveExpertises, "sub_task_type_id")
	if err != nil {
		return nil, err
	}
	effectiveExpertises = utils.RemoveDuplicates(effectiveExpertises)

	return effectiveExpertises, nil
}

func FieldSupervisorGetEffectiveExpertise(log *log.DXLog, fieldExpertiseId int64) (effectiveExpertises []int64, err error) {
	_, fieldSupervisorEffectiveExpertises, err := partner_management.ModulePartnerManagement.FieldSupervisorEffectiveExpertise.Select(log, nil, utils.JSON{
		"user_role_membership_id": fieldExpertiseId,
	}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	effectiveExpertises, err = utils.GetMapValueFromArrayOfJSON[int64](fieldSupervisorEffectiveExpertises, "sub_task_type_id")
	if err != nil {
		return nil, err
	}
	effectiveExpertises = utils.RemoveDuplicates(effectiveExpertises)

	return effectiveExpertises, nil
}

func FieldExecutorGetEffectiveAreas(log *log.DXLog, fieldExecutorId int64) (effectiveAreas []string, err error) {
	_, fieldExecutorEffectiveAreas, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveArea.Select(log, nil, utils.JSON{
		"user_role_membership_id": fieldExecutorId,
	}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	effectiveAreaCodes, err := utils.GetMapValueFromArrayOfJSON[string](fieldExecutorEffectiveAreas, "area_code")
	if err != nil {
		return nil, err
	}

	effectiveAreaCodes = utils.RemoveDuplicates(effectiveAreaCodes)

	return effectiveAreaCodes, nil
}

func FieldSupervisorGetEffectiveAreas(log *log.DXLog, fieldSupervisorId int64) (effectiveAreas []string, err error) {
	_, fieldSupervisorEffectiveAreas, err := partner_management.ModulePartnerManagement.FieldSupervisorEffectiveArea.Select(log, nil, utils.JSON{
		"user_role_membership_id": fieldSupervisorId,
	}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	effectiveAreaCodes, err := utils.GetMapValueFromArrayOfJSON[string](fieldSupervisorEffectiveAreas, "area_code")
	if err != nil {
		return nil, err
	}

	effectiveAreaCodes = utils.RemoveDuplicates(effectiveAreaCodes)

	return effectiveAreaCodes, nil
}

func FieldExecutorGetEffectiveLocations(log *log.DXLog, fieldExecutorId int64) (effectiveLocations []string, err error) {
	_, fieldExecutorEffectiveLocations, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveLocation.Select(log, nil, utils.JSON{
		"user_role_membership_id": fieldExecutorId,
	}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	effectiveLocations, err = utils.GetMapValueFromArrayOfJSON[string](fieldExecutorEffectiveLocations, "location_code")
	if err != nil {
		return nil, err
	}
	return effectiveLocations, nil
}

func FieldSupervisorGetEffectiveLocations(log *log.DXLog, fieldSupervisorId int64) (effectiveLocations []string, err error) {
	_, fieldSupervisorEffectiveLocations, err := partner_management.ModulePartnerManagement.FieldSupervisorEffectiveLocation.Select(log, nil, utils.JSON{
		"user_role_membership_id": fieldSupervisorId,
	}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	effectiveLocations, err = utils.GetMapValueFromArrayOfJSON[string](fieldSupervisorEffectiveLocations, "location_code")
	if err != nil {
		return nil, err
	}
	return effectiveLocations, nil
}

func SelfFieldExecutorTaskSearch(aepr *api.DXAPIEndPointRequest) (err error) {
	isParameterSearchExist, search, err := aepr.GetParameterValueAsString("search", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isParameterTaskTypeIdsExist, taskTypeIds, err := aepr.GetParameterValueAsArrayOfInt64("task_type_ids")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isParameterLocationCodesExist, locationCodes, err := aepr.GetParameterValueAsArrayOfString("location_codes")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isParameterAreaCodesExist, areaCodes, err := aepr.GetParameterValueAsArrayOfString("area_codes")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, statuses, err := aepr.GetParameterValueAsArrayOfString("statuses")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, sort, err := aepr.GetParameterValueAsString("sort", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isParameterLatitudeExist, latitude, err := aepr.GetParameterValueAsFloat64("latitude")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isParameterLongitudeExist, longitude, err := aepr.GetParameterValueAsFloat64("longitude")
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

	userId := aepr.LocalData["user_id"].(int64)

	_, fieldExecutors, err := partner_management.ModulePartnerManagement.FieldExecutor.Select(&aepr.Log, nil, utils.JSON{
		"user_id": userId,
	}, nil, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if len(fieldExecutors) == 0 {
		return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NO_FIELD_EXECUTOR_FOUND")
	}

	conditions := ""
	parameters := []any{}
	driverName := t.Database.Connection.DriverName()

	for _, fieldExecutor := range fieldExecutors {
		userRoleMembershipId := fieldExecutor["id"].(int64)

		effectiveTaskTypeIds, err := FieldExecutorGetEffectiveTaskTypeIds(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveTaskTypeIds) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_EXPERTISE")
		}

		var includedTaskTypeIds []int64

		if isParameterTaskTypeIdsExist {
			if len(taskTypeIds) == 0 {
				includedTaskTypeIds = effectiveTaskTypeIds
			} else {
				var missingTaskTypeIds []int64
				includedTaskTypeIds, missingTaskTypeIds = utils.Diff[int64](taskTypeIds, effectiveTaskTypeIds)
				if len(missingTaskTypeIds) > 0 {
					//	return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "TASK_TYPE_NOT_WITHIN_YOUR_EXPERTISE: %s", strings.Join(utils.Int64SliceToStrings(missingTaskTypeIds), ", "))
				}
				if len(includedTaskTypeIds) == 0 {
					return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "INCLUDED_TASK_TYPE_IDS_IS_EMPTY")
				}
			}
		} else {
			includedTaskTypeIds = effectiveTaskTypeIds
		}

		effectiveAreaCodes, err := FieldExecutorGetEffectiveAreas(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveAreaCodes) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_AREA")
		}

		var includedAreaCodes []string

		if isParameterAreaCodesExist {
			if len(areaCodes) == 0 {
				includedAreaCodes = effectiveAreaCodes
			} else {
				var missingAreaCodes []string
				includedAreaCodes, missingAreaCodes = utils.Diff[string](areaCodes, effectiveAreaCodes)
				if len(missingAreaCodes) > 0 {
					//	return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "AREA_CODE_NOT_WITHIN_YOUR_AREA: %s", strings.Join(missingAreaCodes, ", "))
				}
				if len(includedAreaCodes) == 0 {
					return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "INCLUDED_AREA_CODES_IS_EMPTY")
				}
			}
		} else {
			includedAreaCodes = effectiveAreaCodes
		}

		effectiveLocations, err := FieldExecutorGetEffectiveLocations(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveLocations) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_LOCATION")
		}

		var includedLocationCodes []string

		if isParameterLocationCodesExist {
			if len(locationCodes) == 0 {
				includedLocationCodes = effectiveLocations
			} else {
				var missingLocationCodes []string
				includedLocationCodes, missingLocationCodes = utils.Diff[string](locationCodes, effectiveLocations)
				if len(missingLocationCodes) > 0 {
					//	return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "LOCATION_CODE_NOT_WITHIN_WITH_YOUR_LOCATION: %s", strings.Join(missingLocationCodes, ", "))
				}
				if len(includedLocationCodes) == 0 {
					return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "INCLUDED_LOCATION_CODES_IS_EMPTY")
				}
			}
		} else {
			includedLocationCodes = effectiveLocations
		}

		/* Get Status */

		if len(statuses) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			conditions += dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "status", statuses)
		}

		/* Customer */

		if isParameterSearchExist && search != "" {
			if conditions != "" {
				conditions += " AND "
			}
			likeSearch := "%" + search + "%"
			conditions += dbUtils.SQLBuildParameterizedWhereClause(driverName, &parameters, "((LOWER(customer_fullname) LIKE LOWER(?)) OR (LOWER(customer_number) LIKE LOWER(?)))", likeSearch, likeSearch)
		}

		if len(includedTaskTypeIds) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			conditions += dbUtils.SQLBuildParameterizedWhereInClauseInt64(&parameters, "task_type_id", includedTaskTypeIds)
		}

		if len(includedLocationCodes) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			c1 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_province_location_code", includedLocationCodes)
			c2 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_kabupaten_location_code", includedLocationCodes)
			c3 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_kecamatan_location_code", includedLocationCodes)
			c4 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_kelurahan_location_code", includedLocationCodes)
			conditions += "(" + strings.Join([]string{c1, c2, c3, c4}, " OR ") + ")"
		}

		if len(includedAreaCodes) > 0 {
			allAreaCodes := []string{}
			for _, areaCode := range includedAreaCodes {
				childrenCodes, err := master_data.ModuleMasterData.AreaCodeExpandParentTreeDown(&aepr.Log, areaCode)
				if err != nil {
					return err
				}
				allAreaCodes = append(allAreaCodes, childrenCodes...)
			}

			if conditions != "" {
				conditions += " AND "
			}

			c := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_sales_area_code", allAreaCodes)

			conditions += c
		}
	}

	var orderBy string
	switch sort {
	case "LATEST":
		orderBy = "last_modified_at DESC, id DESC"
	case "CLOSEST_TODAY":
		if !isParameterLatitudeExist || !isParameterLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		orderBy = fmt.Sprintf("customer_geom <-> ST_SetSRID(ST_MakePoint(%.6f, %.6f), 4326)", longitude, latitude)

	case "CLOSEST_TO_THE_LOCATION":
		if !isParameterLatitudeExist || !isParameterLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		orderBy = fmt.Sprintf("customer_geom <-> ST_SetSRID(ST_MakePoint(%.6f, %.6f), 4326)", longitude, latitude)

	default:
		orderBy = ""
	}
	rowsInfo, list, totalRows, totalPage, _, err := db.QueryPaging(
		t.Database.Connection, t.FieldTypeMapping,
		"",
		rowPerPage,
		pageIndex,
		"*",
		t.ListViewNameId,
		conditions,
		"",
		orderBy,
		parameters,
	)
	if err != nil {
		aepr.Log.Errorf(err, "Error At paging table %s (%s) sql_where=%v sql_parameters=%v", t.NameId, err.Error(), conditions, parameters)
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

func SelfFieldExecutorSubTaskSearch(aepr *api.DXAPIEndPointRequest) (err error) {
	isSearchExist, search, err := aepr.GetParameterValueAsString("search", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isTaskUidExist, taskUid, err := aepr.GetParameterValueAsString("task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isSubTaskTypeIdsExist, subTaskTypeIds, err := aepr.GetParameterValueAsArrayOfInt64("sub_task_type_ids")
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
	_, statuses, err := aepr.GetParameterValueAsArrayOfString("statuses")
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

	_, fieldExecutors, err := partner_management.ModulePartnerManagement.FieldExecutor.Select(&aepr.Log, nil, utils.JSON{
		"user_id": userId,
	}, nil, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if len(fieldExecutors) == 0 {
		return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NO_FIELD_EXECUTOR_FOUND")
	}

	conditions := ""
	parameters := []any{}
	driverName := t.Database.Connection.DriverName()

	for _, fieldExecutor := range fieldExecutors {
		userRoleMembershipId := fieldExecutor["id"].(int64)

		effectiveExpertises, err := FieldExecutorGetEffectiveExpertise(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveExpertises) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_EXPERTISE")
		}

		var includedSubTaskTypeIds []int64

		if isSubTaskTypeIdsExist {
			if len(subTaskTypeIds) == 0 {
				includedSubTaskTypeIds = effectiveExpertises
			} else {
				var missingSubTaskTypeIds []int64
				includedSubTaskTypeIds, missingSubTaskTypeIds = utils.Diff[int64](subTaskTypeIds, effectiveExpertises)
				if len(missingSubTaskTypeIds) > 0 {
					// return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_TYPE_NOT_WITHIN_YOUR_EXPERTISE: %s", strings.Join(utils.Int64SliceToStrings(missingSubTaskTypeIds), ", "))
				}
				if len(includedSubTaskTypeIds) == 0 {
					return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "INCLUDED_SUB_TASK_TYPE_IDS_IS_EMPTY")
				}
			}
		} else {
			includedSubTaskTypeIds = effectiveExpertises
		}

		effectiveAreaCodes, err := FieldExecutorGetEffectiveAreas(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveAreaCodes) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_AREA")
		}

		var includedAreaCodes []string

		if isAreaCodesExist {
			if len(areaCodes) == 0 {
				includedAreaCodes = effectiveAreaCodes
			} else {
				var missingAreaCodes []string
				includedAreaCodes, missingAreaCodes = utils.Diff[string](areaCodes, effectiveAreaCodes)
				if len(missingAreaCodes) > 0 {
					// return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "AREA_CODE_NOT_WITHIN_YOUR_AREA: %s", strings.Join(missingAreaCodes, ", "))
				}
				if len(includedAreaCodes) == 0 {
					return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "INCLUDED_AREA_CODES_IS_EMPTY")
				}
			}
		} else {
			includedAreaCodes = effectiveAreaCodes
		}

		effectiveLocations, err := FieldExecutorGetEffectiveLocations(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveLocations) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_LOCATION")
		}
		var includedLocationCodes []string

		if isLocationCodesExist {
			if len(locationCodes) == 0 {
				includedLocationCodes = effectiveLocations
			} else {
				var missingLocationCodes []string
				includedLocationCodes, missingLocationCodes = utils.Diff[string](locationCodes, effectiveLocations)
				if len(missingLocationCodes) > 0 {
					// return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "LOCATION_CODE_NOT_WITHIN_WITH_YOUR_LOCATION: %s", strings.Join(missingLocationCodes, ", "))
				}
				if len(includedLocationCodes) == 0 {
					return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "INCLUDED_LOCATION_CODES_IS_EMPTY")
				}
			}
		} else {
			includedLocationCodes = effectiveLocations
		}

		/* Get TaskUid */
		if isTaskUidExist {
			if conditions != "" {
				conditions += " AND "
			}
			conditions += dbUtils.SQLBuildParameterizedWhereClause(driverName, &parameters, "task_uid=?", taskUid)
		}

		/* Get Status */

		if len(statuses) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			conditions += dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "status", statuses)
		}

		/* Customer */

		if isSearchExist && search != "" {
			if conditions != "" {
				conditions += " AND "
			}
			likeSearch := "%" + search + "%"
			conditions += dbUtils.SQLBuildParameterizedWhereClause(driverName, &parameters, "((LOWER(customer_fullname) LIKE LOWER(?)) OR (LOWER(customer_number) LIKE LOWER(?)))", likeSearch, likeSearch)
		}

		if len(includedSubTaskTypeIds) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			conditions += dbUtils.SQLBuildParameterizedWhereInClauseInt64(&parameters, "sub_task_type_id", includedSubTaskTypeIds)
		}

		if len(includedLocationCodes) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			c1 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_province_location_code", includedLocationCodes)
			c2 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_kabupaten_location_code", includedLocationCodes)
			c3 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_kecamatan_location_code", includedLocationCodes)
			c4 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_kelurahan_location_code", includedLocationCodes)
			conditions += "(" + strings.Join([]string{c1, c2, c3, c4}, " OR ") + ")"
		}

		if len(includedAreaCodes) > 0 {
			allAreaCodes := []string{}
			for _, areaCode := range includedAreaCodes {
				childrenCodes, err := master_data.ModuleMasterData.AreaCodeExpandParentTreeDown(&aepr.Log, areaCode)
				if err != nil {
					return err
				}
				allAreaCodes = append(allAreaCodes, childrenCodes...)
			}

			if conditions != "" {
				conditions += " AND "
			}

			c := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_sales_area_code", allAreaCodes)

			conditions += c
		}
	}

	var orderBy string
	switch sort {
	case "LATEST":
		orderBy = "created_at DESC, last_modified_at DESC"
	case "CLOSEST_TODAY":
		if !isLatitudeExist || !isLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		orderBy = fmt.Sprintf("customer_geom <-> ST_SetSRID(ST_MakePoint(%.6f, %.6f), 4326)", longitude, latitude)
	case "CLOSEST_TO_THE_LOCATION":
		if !isLatitudeExist || !isLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		orderBy = fmt.Sprintf("customer_geom <-> ST_SetSRID(ST_MakePoint(%.6f, %.6f), 4326)", longitude, latitude)
	default:
		orderBy = ""
	}

	rowsInfo, list, totalRows, totalPage, _, err := db.QueryPaging(
		t.Database.Connection, t.FieldTypeMapping,
		"",
		rowPerPage,
		pageIndex,
		"*",
		t.ListViewNameId,
		conditions,
		"",
		orderBy,
		parameters,
	)
	if err != nil {
		aepr.Log.Errorf(err, "Error At paging table %s (%s) sql_where=%v sql_parameters=%v", t.NameId, err.Error(), conditions, parameters)
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

func SelfFieldSupervisorTaskSearch(aepr *api.DXAPIEndPointRequest) (err error) {
	isParameterSearchExist, search, err := aepr.GetParameterValueAsString("search", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isParameterTaskTypeIdsExist, taskTypeIds, err := aepr.GetParameterValueAsArrayOfInt64("task_type_ids")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isParameterLocationCodesExist, locationCodes, err := aepr.GetParameterValueAsArrayOfString("location_codes")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isParameterAreaCodesExist, areaCodes, err := aepr.GetParameterValueAsArrayOfString("area_codes")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, statuses, err := aepr.GetParameterValueAsArrayOfString("statuses")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, sort, err := aepr.GetParameterValueAsString("sort", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isParameterLatitudeExist, latitude, err := aepr.GetParameterValueAsFloat64("latitude")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isParameterLongitudeExist, longitude, err := aepr.GetParameterValueAsFloat64("longitude")
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

	userId := aepr.LocalData["user_id"].(int64)

	_, fieldSupervisors, err := partner_management.ModulePartnerManagement.FieldSupervisor.Select(&aepr.Log, nil, utils.JSON{
		"user_id": userId,
	}, nil, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if len(fieldSupervisors) == 0 {
		return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:USER_IS_NOT_FIELD_SUPERVISOR")
	}

	conditions := ""
	parameters := []any{}
	driverName := t.Database.Connection.DriverName()

	for _, fieldSupervisor := range fieldSupervisors {
		userRoleMembershipId := fieldSupervisor["id"].(int64)

		effectiveTaskTypeIds, err := FieldSupervisorGetEffectiveTaskTypeIds(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveTaskTypeIds) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_EXPERTISE")
		}

		var includedTaskTypeIds []int64

		if isParameterTaskTypeIdsExist {
			if len(taskTypeIds) == 0 {
				includedTaskTypeIds = effectiveTaskTypeIds
			} else {
				var missingTaskTypeIds []int64
				includedTaskTypeIds, missingTaskTypeIds = utils.Diff[int64](taskTypeIds, effectiveTaskTypeIds)
				if len(missingTaskTypeIds) > 0 {
					//	return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "TASK_TYPE_NOT_WITHIN_YOUR_EXPERTISE: %s", strings.Join(utils.Int64SliceToStrings(missingTaskTypeIds), ", "))
				}
				if len(includedTaskTypeIds) == 0 {
					return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "INCLUDED_TASK_TYPE_IDS_IS_EMPTY")
				}
			}
		} else {
			includedTaskTypeIds = effectiveTaskTypeIds
		}

		effectiveAreaCodes, err := FieldSupervisorGetEffectiveAreas(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveAreaCodes) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_AREA")
		}

		var includedAreaCodes []string

		if isParameterAreaCodesExist {
			if len(areaCodes) == 0 {
				includedAreaCodes = effectiveAreaCodes
			} else {
				var missingAreaCodes []string
				includedAreaCodes, missingAreaCodes = utils.Diff[string](areaCodes, effectiveAreaCodes)
				if len(missingAreaCodes) > 0 {
					//	return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "AREA_CODE_NOT_WITHIN_YOUR_AREA: %s", strings.Join(missingAreaCodes, ", "))
				}
				if len(includedAreaCodes) == 0 {
					return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "INCLUDED_AREA_CODES_IS_EMPTY")
				}
			}
		} else {
			includedAreaCodes = effectiveAreaCodes
		}

		effectiveLocations, err := FieldSupervisorGetEffectiveLocations(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveLocations) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_LOCATION")
		}

		var includedLocationCodes []string

		if isParameterLocationCodesExist {
			if len(locationCodes) == 0 {
				includedLocationCodes = effectiveLocations
			} else {
				var missingLocationCodes []string
				includedLocationCodes, missingLocationCodes = utils.Diff[string](locationCodes, effectiveLocations)
				if len(missingLocationCodes) > 0 {
					//	return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "LOCATION_CODE_NOT_WITHIN_WITH_YOUR_LOCATION: %s", strings.Join(missingLocationCodes, ", "))
				}
				if len(includedLocationCodes) == 0 {
					return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "INCLUDED_LOCATION_CODES_IS_EMPTY")
				}
			}
		} else {
			includedLocationCodes = effectiveLocations
		}

		/* Get Status */

		if len(statuses) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			conditions += dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "status", statuses)
		}

		/* Customer */

		if isParameterSearchExist && search != "" {
			if conditions != "" {
				conditions += " AND "
			}
			likeSearch := "%" + search + "%"
			conditions += dbUtils.SQLBuildParameterizedWhereClause(driverName, &parameters, "((LOWER(customer_fullname) LIKE LOWER(?)) OR (LOWER(customer_number) LIKE LOWER(?)))", likeSearch, likeSearch)
		}

		if len(includedTaskTypeIds) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			conditions += dbUtils.SQLBuildParameterizedWhereInClauseInt64(&parameters, "task_type_id", includedTaskTypeIds)
		}

		if len(includedLocationCodes) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			c1 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_province_location_code", includedLocationCodes)
			c2 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_kabupaten_location_code", includedLocationCodes)
			c3 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_kecamatan_location_code", includedLocationCodes)
			c4 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_kelurahan_location_code", includedLocationCodes)
			conditions += "(" + strings.Join([]string{c1, c2, c3, c4}, " OR ") + ")"
		}

		if len(includedAreaCodes) > 0 {
			allAreaCodes := []string{}
			for _, areaCode := range includedAreaCodes {
				childrenCodes, err := master_data.ModuleMasterData.AreaCodeExpandParentTreeDown(&aepr.Log, areaCode)
				if err != nil {
					return err
				}
				allAreaCodes = append(allAreaCodes, childrenCodes...)
			}

			if conditions != "" {
				conditions += " AND "
			}

			c := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_sales_area_code", allAreaCodes)

			conditions += c
		}
	}

	var orderBy string
	switch sort {
	case "LATEST":
		orderBy = "last_modified_at DESC, id DESC"
	case "CLOSEST_TODAY":
		if !isParameterLatitudeExist || !isParameterLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		orderBy = fmt.Sprintf("customer_geom <-> ST_SetSRID(ST_MakePoint(%.6f, %.6f), 4326)", longitude, latitude)

	case "CLOSEST_TO_THE_LOCATION":
		if !isParameterLatitudeExist || !isParameterLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		orderBy = fmt.Sprintf("customer_geom <-> ST_SetSRID(ST_MakePoint(%.6f, %.6f), 4326)", longitude, latitude)

	default:
		orderBy = ""
	}

	rowsInfo, list, totalRows, totalPage, _, err := db.QueryPaging(
		t.Database.Connection, t.FieldTypeMapping,
		"",
		rowPerPage,
		pageIndex,
		"*",
		t.ListViewNameId,
		conditions,
		"",
		orderBy,
		parameters,
	)
	if err != nil {
		aepr.Log.Errorf(err, "Error At paging table %s (%s) sql_where=%v sql_parameters=%v", t.NameId, err.Error(), conditions, parameters)
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

func SelfFieldExecutorTaskSummary(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.LocalData["user_id"].(int64)

	// Check if user is a field executor
	_, fieldExecutors, err := partner_management.ModulePartnerManagement.FieldExecutor.Select(&aepr.Log, nil, utils.JSON{
		"user_id": userId,
	}, nil, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if len(fieldExecutors) == 0 {
		return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NO_FIELD_EXECUTOR_FOUND")
	}

	// Initialize counters
	totalKonstruksi := int64(0)
	totalPenangananPiutang := int64(0)
	totalPengaduanGangguan := int64(0)
	totalLayananTeknis := int64(0)

	// Get the database connection
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

	// Build conditions for the current field executor
	// Count Construction tasks (TaskTypeIdConstruction = 1)
	constructionCondition := "last_field_executor_user_id = $1 AND task_type_id = $2"
	constructionParameters := []any{userId, base.TaskTypeIdConstruction}
	_, _, constructionCount, _, _, err := db.QueryPaging(
		t.Database.Connection,
		t.FieldTypeMapping,
		"",
		1,    // rowPerPage - we only need the count
		0,    // pageIndex
		"id", // fields
		t.ListViewNameId,
		constructionCondition,
		"",
		"",
		constructionParameters,
	)
	if err != nil {
		aepr.Log.Errorf(err, "Error counting construction tasks: %s", err.Error())
		return errors.Wrap(err, "error occured")
	}
	totalKonstruksi = constructionCount

	// Count Debt Management tasks (TaskTypeIdDebtManagement = 3)
	debtCondition := "last_field_executor_user_id = $1 AND task_type_id = $2"
	debtParameters := []any{userId, base.TaskTypeIdDebtManagement}

	_, _, debtCount, _, _, err := db.QueryPaging(
		t.Database.Connection,
		t.FieldTypeMapping,
		"",
		1,    // rowPerPage
		0,    // pageIndex
		"id", // fields
		t.ListViewNameId,
		debtCondition,
		"",
		"",
		debtParameters,
	)
	if err != nil {
		aepr.Log.Errorf(err, "Error counting debt management tasks: %s", err.Error())
		return errors.Wrap(err, "error occured")
	}
	totalPenangananPiutang = debtCount

	// Count Technical Support tasks (TaskTypeIdTechnicalSupport = 2)
	technicalCondition := "last_field_executor_user_id = $1 AND task_type_id = $2"
	technicalParameters := []any{userId, base.TaskTypeIdTechnicalSupport}

	_, _, technicalCount, _, _, err := db.QueryPaging(
		t.Database.Connection,
		t.FieldTypeMapping,
		"",
		1,    // rowPerPage
		0,    // pageIndex
		"id", // fields
		t.ListViewNameId,
		technicalCondition,
		"",
		"",
		technicalParameters,
	)
	if err != nil {
		aepr.Log.Errorf(err, "Error counting technical support tasks: %s", err.Error())
		return errors.Wrap(err, "error occured")
	}
	totalLayananTeknis = technicalCount

	// For now, we'll set pengaduan_gangguan to 0 since we don't have specific information
	// In a real implementation, you would need to identify the correct criteria
	totalPengaduanGangguan = 0

	// Return the summary
	data := utils.JSON{"data": utils.JSON{
		"total_konstruksi":         totalKonstruksi,
		"total_penanganan_piutang": totalPenangananPiutang,
		"total_pengaduan_gangguan": totalPengaduanGangguan,
		"total_layanan_teknis":     totalLayananTeknis,
	}}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, data)
	return nil
}

func SelfFieldSupervisorSubTaskSearch(aepr *api.DXAPIEndPointRequest) (err error) {
	isSearchExist, search, err := aepr.GetParameterValueAsString("search", "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isTaskUidExist, taskUid, err := aepr.GetParameterValueAsString("task_uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	isSubTaskTypeIdsExist, subTaskTypeIds, err := aepr.GetParameterValueAsArrayOfInt64("sub_task_type_ids")
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
	_, statuses, err := aepr.GetParameterValueAsArrayOfString("statuses")
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

	_, fieldSupervisors, err := partner_management.ModulePartnerManagement.FieldSupervisor.Select(&aepr.Log, nil, utils.JSON{
		"user_id": userId,
	}, nil, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if len(fieldSupervisors) == 0 {
		return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:USER_IS_NOT_FIELD_SUPERVISOR")
	}

	conditions := ""
	parameters := []any{}
	driverName := t.Database.Connection.DriverName()

	for _, fieldSupervisor := range fieldSupervisors {
		userRoleMembershipId := fieldSupervisor["id"].(int64)

		effectiveExpertises, err := FieldSupervisorGetEffectiveExpertise(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveExpertises) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_EXPERTISE")
		}

		var includedSubTaskTypeIds []int64

		if isSubTaskTypeIdsExist {
			if len(subTaskTypeIds) == 0 {
				includedSubTaskTypeIds = effectiveExpertises
			} else {
				var missingSubTaskTypeIds []int64
				includedSubTaskTypeIds, missingSubTaskTypeIds = utils.Diff[int64](subTaskTypeIds, effectiveExpertises)
				if len(missingSubTaskTypeIds) > 0 {
					// return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "SUB_TASK_TYPE_NOT_WITHIN_YOUR_EXPERTISE: %s", strings.Join(utils.Int64SliceToStrings(missingSubTaskTypeIds), ", "))
				}
				if len(includedSubTaskTypeIds) == 0 {
					return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "INCLUDED_SUB_TASK_TYPE_IDS_IS_EMPTY")
				}
			}
		} else {
			includedSubTaskTypeIds = effectiveExpertises
		}

		effectiveAreaCodes, err := FieldSupervisorGetEffectiveAreas(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveAreaCodes) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_AREA")
		}

		var includedAreaCodes []string

		if isAreaCodesExist {
			if len(areaCodes) == 0 {
				includedAreaCodes = effectiveAreaCodes
			} else {
				var missingAreaCodes []string
				includedAreaCodes, missingAreaCodes = utils.Diff[string](areaCodes, effectiveAreaCodes)
				if len(missingAreaCodes) > 0 {
					// return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "AREA_CODE_NOT_WITHIN_YOUR_AREA: %s", strings.Join(missingAreaCodes, ", "))
				}
				if len(includedAreaCodes) == 0 {
					return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "INCLUDED_AREA_CODES_IS_EMPTY")
				}
			}
		} else {
			includedAreaCodes = effectiveAreaCodes
		}

		effectiveLocations, err := FieldSupervisorGetEffectiveLocations(&aepr.Log, userRoleMembershipId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		if len(effectiveLocations) == 0 {
			return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "NOT_ERROR:YOU_DONT_HAVE_ANY_LOCATION")
		}

		var includedLocationCodes []string

		if isLocationCodesExist {
			if len(locationCodes) == 0 {
				includedLocationCodes = effectiveLocations
			} else {
				var missingLocationCodes []string
				includedLocationCodes, missingLocationCodes = utils.Diff[string](locationCodes, effectiveLocations)
				if len(missingLocationCodes) > 0 {
					//	return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "LOCATION_CODE_NOT_WITHIN_WITH_YOUR_LOCATION: %s", strings.Join(missingLocationCodes, ", "))
				}
				if len(includedLocationCodes) == 0 {
					return aepr.WriteResponseAndNewErrorf(http.StatusForbidden, "", "INCLUDED_LOCATION_CODES_IS_EMPTY")
				}
			}
		} else {
			includedLocationCodes = effectiveLocations
		}

		/* Get TaskUid */
		if isTaskUidExist {
			if conditions != "" {
				conditions += " AND "
			}
			conditions += dbUtils.SQLBuildParameterizedWhereClause(driverName, &parameters, "task_uid=?", taskUid)
		}

		/* Get Status */

		if len(statuses) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			conditions += dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "status", statuses)
		}

		/* Customer */

		if isSearchExist && search != "" {
			if conditions != "" {
				conditions += " AND "
			}
			likeSearch := "%" + search + "%"
			conditions += dbUtils.SQLBuildParameterizedWhereClause(driverName, &parameters, "((LOWER(customer_fullname) LIKE LOWER(?)) OR (LOWER(customer_number) LIKE LOWER(?)))", likeSearch, likeSearch)
		}

		if len(includedSubTaskTypeIds) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			conditions += dbUtils.SQLBuildParameterizedWhereInClauseInt64(&parameters, "sub_task_type_id", includedSubTaskTypeIds)
		}

		if len(includedLocationCodes) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			c1 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_province_location_code", includedLocationCodes)
			c2 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_kabupaten_location_code", includedLocationCodes)
			c3 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_kecamatan_location_code", includedLocationCodes)
			c4 := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_address_kelurahan_location_code", includedLocationCodes)
			conditions += "(" + strings.Join([]string{c1, c2, c3, c4}, " OR ") + ")"
		}

		/* Get Status */

		if len(statuses) > 0 {
			if conditions != "" {
				conditions += " AND "
			}
			conditions += dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "status", statuses)
		}

		if len(includedAreaCodes) > 0 {
			allAreaCodes := []string{}
			for _, areaCode := range includedAreaCodes {
				childrenCodes, err := master_data.ModuleMasterData.AreaCodeExpandParentTreeDown(&aepr.Log, areaCode)
				if err != nil {
					return err
				}
				allAreaCodes = append(allAreaCodes, childrenCodes...)
			}

			if conditions != "" {
				conditions += " AND "
			}

			c := dbUtils.SQLBuildParameterizedWhereInClause(driverName, &parameters, "customer_sales_area_code", allAreaCodes)

			conditions += c
		}
	}

	var orderBy string
	switch sort {
	case "LATEST":
		orderBy = "last_modified_at DESC, created_at DESC"
	case "CLOSEST_TODAY":
		if !isLatitudeExist || !isLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		orderBy = fmt.Sprintf("customer_geom <-> ST_SetSRID(ST_MakePoint(%.6f, %.6f), 4326)", longitude, latitude)
	case "CLOSEST_TO_THE_LOCATION":
		if !isLatitudeExist || !isLongitudeExist {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "LATITUDE_AND_LONGITUDE_MUST_BE_PROVIDED")
		}
		orderBy = fmt.Sprintf("customer_geom <-> ST_SetSRID(ST_MakePoint(%.6f, %.6f), 4326)", longitude, latitude)
	default:
		orderBy = ""
	}

	rowsInfo, list, totalRows, totalPage, _, err := db.QueryPaging(
		t.Database.Connection,
		t.FieldTypeMapping,
		"",
		rowPerPage,
		pageIndex,
		"*",
		t.ListViewNameId,
		conditions,
		"",
		orderBy,
		parameters,
	)
	if err != nil {
		aepr.Log.Errorf(err, "Error At paging table %s (%s) sql_where=%v sql_parameters=%v", t.NameId, err.Error(), conditions, parameters)
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
