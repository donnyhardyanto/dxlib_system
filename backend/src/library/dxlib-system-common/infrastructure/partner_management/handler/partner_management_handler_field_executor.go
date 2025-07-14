package handler

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"net/http"
)

func FieldExecutorPagingList(aepr *api.DXAPIEndPointRequest) (err error) {

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

	return partner_management.ModulePartnerManagement.FieldExecutor.DoRequestPagingList(aepr, filterWhere, filterOrderBy, filterKeyValues, func(listRow utils.JSON) (r utils.JSON, err error) {
		_, fieldExecutorArea, err := partner_management.ModulePartnerManagement.FieldExecutorArea.Select(&aepr.Log, nil, utils.JSON{
			"user_role_membership_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_executor_area"] = fieldExecutorArea

		_, fieldExecutorLocation, err := partner_management.ModulePartnerManagement.FieldExecutorLocation.Select(&aepr.Log, nil, utils.JSON{
			"user_role_membership_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_executor_location"] = fieldExecutorLocation

		_, fieldExecutorExpertise, err := partner_management.ModulePartnerManagement.FieldExecutorExpertise.Select(&aepr.Log, nil, utils.JSON{
			"user_role_membership_id": listRow["id"],
		}, nil, map[string]string{"id": "asc"}, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_executor_expertise"] = fieldExecutorExpertise

		_, fieldExecutorEffectiveArea, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveArea.Select(&aepr.Log, nil, utils.JSON{
			"user_role_membership_id": listRow["id"],
		}, nil, nil, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_executor_effective_area"] = fieldExecutorEffectiveArea

		_, fieldExecutorEffectiveLocation, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveLocation.Select(&aepr.Log, nil, utils.JSON{
			"user_role_membership_id": listRow["id"],
		}, nil, nil, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_executor_effective_location"] = fieldExecutorEffectiveLocation

		_, fieldExecutorEffectiveExpertise, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveExpertise.Select(&aepr.Log, nil, utils.JSON{
			"user_role_membership_id": listRow["id"],
		}, nil, nil, nil, nil)
		if err != nil {
			return listRow, err
		}

		listRow["field_executor_effective_expertise"] = fieldExecutorEffectiveExpertise

		return listRow, nil
	})
}

func SelfFieldExecutorSubTaskStatusStats(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.LocalData["user_id"].(int64)
	_, data, err := partner_management.ModulePartnerManagement.FieldExecutorSubTaskStatusStats.SelectOne(&aepr.Log, nil, utils.JSON{
		"user_id": userId,
	}, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
		"data": data,
	})
	return nil
}

/*func FieldExecutorListDownload(aepr *api.DXAPIEndPointRequest) (err error) {
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

	_, format, err := aepr.GetParameterValueAsString("format")
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "FORMAT_PARAMETER_ERROR:%s", err.Error())
	}

	format = strings.ToLower(format)

	isDeletedIncluded := false
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

	t := partner_management.ModulePartnerManagement.FieldExecutor

	if t.Database == nil {
		t.Database = database.Manager.Databases[t.DatabaseNameId]
	}

	if !t.Database.Connected {
		err := t.Database.Connect()
		if err != nil {
			aepr.Log.Errorf("error At reconnect db At table %s list (%s) ", t.NameId, err.Error())
			return errors.Wrap(err, "error occured")
		}
	}

	rowsInfo, list, err := db.NamedQueryList(t.Database.Connection, t.FieldTypeMapping, "*", t.ListViewNameId,
		filterWhere, "", filterOrderBy, filterKeyValues)

	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	// Set export options
	opts := export.ExportOptions{
		Format:     export.ExportFormat(format),
		SheetName:  "Sheet1",
		DateFormat: "2006-01-02 15:04:05",
	}

	// Get file as stream
	data, contentType, err := export.ExportToStream(rowsInfo, list, opts)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	// Set response headers
	filename := fmt.Sprintf("export_%s_%s.%s", t.NameId, time.Now().Format("20060102_150405"), format)

	responseWriter := *aepr.GetResponseWriter()
	responseWriter.Header().Set("Content-Type", contentType)
	responseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	responseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	responseWriter.WriteHeader(http.StatusOK)
	aepr.ResponseStatusCode = http.StatusOK

	_, err = responseWriter.Write(data)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.ResponseHeaderSent = true
	aepr.ResponseBodySent = true

	return nil
}
*/
