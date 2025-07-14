package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/database/protected/db"
	"github.com/donnyhardyanto/dxlib/database/protected/export"
	dbUtils "github.com/donnyhardyanto/dxlib/database2/utils"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/master_data"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"github.com/tealeg/xlsx"
)

func CustomerCreateBulk(aepr *api.DXAPIEndPointRequest) (err error) {
	// Get the request body stream
	bs := aepr.Request.Body
	if bs == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "FAILED_TO_GET_BODY_STREAM:%s", "CustomerCreateBulk")
	}
	defer bs.Close()

	// RequestRead the entire request body into a buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, bs)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "FAILED_TO_READ_REQUEST_BODY:%s=%v", "CustomerCreateBulk", err.Error())
	}

	// Determine the file type and parse accordingly
	contentType := aepr.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, "csv") {
		err = parseAndCreateCustomersFromCSV(&buf, aepr)
	} else if strings.Contains(contentType, "excel") || strings.Contains(contentType, "spreadsheetml") {
		err = parseAndCreateCustomersFromXLSX(&buf, aepr)
	} else {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnsupportedMediaType, "UNSUPPORTED_FILE_TYPE:%s", contentType)
	}

	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
	return nil
}

func CustomerCreateBulkBase64(aepr *api.DXAPIEndPointRequest) (err error) {
	_, fileContentBase64, err := aepr.GetParameterValueAsString("content_base64")
	if err != nil {
		return err
	}

	_, contentType, err := aepr.GetParameterValueAsString("content_type")
	if err != nil {
		return err
	}

	// Decode base64 content
	decodedBytes, err := base64.StdEncoding.DecodeString(fileContentBase64)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "INVALID_BASE64_CONTENT")
	}

	// Create a buffer with the decoded content
	var buf bytes.Buffer
	if _, err := buf.Write(decodedBytes); err != nil {
		return err
	}

	// Determine the file type and parse accordingly
	if strings.Contains(strings.ToLower(contentType), "csv") {
		err = parseAndCreateCustomersFromCSV(&buf, aepr)
	} else if strings.Contains(strings.ToLower(contentType), "excel") || strings.Contains(strings.ToLower(contentType), "spreadsheetml") {
		err = parseAndCreateCustomersFromXLSX(&buf, aepr)
	} else {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "UNSUPPORTED_FILE_FORMAT")
	}

	if err != nil {
		return errors.Wrap(err, "error occurred")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
	return nil
}

func parseAndCreateCustomersFromCSV(buf *bytes.Buffer, aepr *api.DXAPIEndPointRequest) error {
	// Create a new reader with comma as delimiter
	reader := csv.NewReader(buf)
	reader.Comma = ';'
	reader.LazyQuotes = true    // Handle quotes more flexibly
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read header row
	headers, err := reader.Read()
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity,
			"FAILED_TO_READ_CSV_HEADERS: %s", err.Error())
	}

	// Clean headers - trim spaces and empty fields
	cleanHeaders := make([]string, 0)
	for _, h := range headers {
		h = strings.TrimSpace(h)
		if h != "" {
			cleanHeaders = append(cleanHeaders, h)
		}
	}

	// Process each row
	lineNum := 1 // Keep track of line numbers for error reporting
	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "",
				"FAILED_TO_PARSE_CSV_LINE_%d: %s", lineNum, err.Error())
		}

		// Create customer data map
		customerData := make(map[string]interface{})
		for i, value := range record {
			if i >= len(cleanHeaders) {
				break
			}
			// Clean and validate the value
			value = strings.TrimSpace(value)
			if value != "" {
				customerData[cleanHeaders[i]] = value
			}
		}

		// Skip empty rows
		if len(customerData) == 0 {
			continue
		}

		// Create customer
		_, err = task_management.ModuleTaskManagement.DoCustomerCreate(&aepr.Log, customerData)
		if err != nil {
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "",
				"FAILED_TO_CREATE_CUSTOMER_LINE_%d: %s", lineNum, err.Error())
		}
	}

	return nil
}

func parseAndCreateCustomersFromXLSX(buf *bytes.Buffer, aepr *api.DXAPIEndPointRequest) error {
	xlFile, err := xlsx.OpenBinary(buf.Bytes())
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "FAILED_TO_PARSE_XLSX: %s", err.Error())
	}

	for _, sheet := range xlFile.Sheets {
		if len(sheet.Rows) < 2 {
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "XLSX_FILE_MUST_HAVE_HEADER_AND_DATA")
		}

		// Validate and extract headers
		headers := make([]string, 0, len(sheet.Rows[0].Cells))
		for _, cell := range sheet.Rows[0].Cells {
			header := strings.TrimSpace(cell.String())
			if header == "" {
				return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "EMPTY_HEADER_NOT_ALLOWED")
			}
			headers = append(headers, header)
		}

		// Process data rows
		for rowIdx, row := range sheet.Rows[1:] {
			customerData := make(map[string]interface{}, len(headers))

			if len(row.Cells) == 0 {
				continue // Skip empty rows
			}

			// Map cell values to headers with type conversion
			for i, cell := range row.Cells {
				if i >= len(headers) {
					break
				}

				value := strings.TrimSpace(cell.String())
				if value == "" {
					continue // Skip empty values instead of adding them to customerData
				}

				// Try to convert numeric values
				if isNumericColumn(headers[i]) {
					if numVal, err := cell.Float(); err == nil {
						customerData[headers[i]] = numVal
					} else {
						return aepr.WriteResponseAndNewErrorf(
							http.StatusUnprocessableEntity, "",
							"INVALID_NUMERIC_VALUE_AT_ROW_%d_COLUMN_%s: %q",
							rowIdx+2,
							headers[i],
							value,
						)
					}
				} else {
					customerData[headers[i]] = value
				}
			}

			if len(customerData) == 0 {
				continue
			}

			if _, err = task_management.ModuleTaskManagement.DoCustomerCreate(&aepr.Log, customerData); err != nil {
				// Check for specific PostgreSQL errors
				if strings.Contains(err.Error(), "invalid input syntax for type double precision") {
					return aepr.WriteResponseAndNewErrorf(
						http.StatusUnprocessableEntity, "",
						"INVALID_NUMERIC_VALUE_AT_ROW_%d: Please ensure all numeric fields contain valid numbers",
						rowIdx+2,
					)
				}
				return aepr.WriteResponseAndNewErrorf(
					http.StatusUnprocessableEntity, "",
					"FAILED_TO_CREATE_CUSTOMER_AT_ROW_%d: %s",
					rowIdx+2,
					err.Error(),
				)
			}
		}
	}

	return nil
}

// Helper function to identify numeric columns
func isNumericColumn(header string) bool {
	// Add your numeric column names here
	numericColumns := map[string]bool{
		"amount":   true,
		"balance":  true,
		"price":    true,
		"quantity": true,
		// Add other numeric column names as needed
	}

	header = strings.ToLower(header)
	return numericColumns[header]
}

func CustomerNearestList(aepr *api.DXAPIEndPointRequest) (err error) {
	_, filterWhere, err := aepr.GetParameterValueAsString("filter_where")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, filterKeyValues, err := aepr.GetParameterValueAsJSON("filter_key_values")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, longitude, err := aepr.GetParameterValueAsFloat64("longitude")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, latitude, err := aepr.GetParameterValueAsFloat64("latitude")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, limit, err := aepr.GetParameterValueAsInt64("limit")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, isDeletedIncluded, err := aepr.GetParameterValueAsBool("is_deleted", false)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	d := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]

	if !isDeletedIncluded {
		if filterWhere != "" {
			filterWhere = fmt.Sprintf("(%s) and ", filterWhere)
		}

		switch d.DatabaseType.String() {
		case "sqlserver":
			filterWhere = filterWhere + "(is_deleted=0)"
		case "postgres":
			filterWhere = filterWhere + "(is_deleted=false)"
		default:
			filterWhere = filterWhere + "(is_deleted=0)"
		}
	}

	if !d.Connected {
		err := d.Connect()
		if err != nil {
			aepr.Log.Errorf(err, "DB_RECONNECT_ERROR:%s ", err.Error())
			return errors.Wrap(err, "error occured")
		}
	}

	query := `
		SELECT c.*,
			ST_Distance(c.geom, ST_SetSRID(ST_MakePoint(:longitude, :latitude), 4326)) AS distance
		FROM task_management.customer c`

	if filterWhere != "" {
		query += " WHERE " + filterWhere + " "
	}

	query = query + `
		ORDER BY c.geom <-> ST_SetSRID(ST_MakePoint(:longitude, :latitude), 4326)
		LIMIT :limit
	`

	filterKeyValues["longitude"] = longitude
	filterKeyValues["latitude"] = latitude
	filterKeyValues["limit"] = limit

	rowsInfo, list, err := db.NamedQueryRows(d.Connection, nil, query, filterKeyValues)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	data := utils.JSON{"data": utils.JSON{
		"list": utils.JSON{
			"rows":      list,
			"rows_info": rowsInfo,
		},
	}}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, data)

	return nil
}

func CustomerCreate(aepr *api.DXAPIEndPointRequest) (err error) {
	_, customerSegmentCode, err := aepr.GetParameterValueAsString("customer_segment_code")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, customerTypeCode, err := aepr.GetParameterValueAsString("customer_type_code")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, customerType, err := master_data.ModuleMasterData.CustomerRef.ShouldSelectOne(&aepr.Log, utils.JSON{
		"type":   "CUSTOMER_TYPE",
		"code":   customerTypeCode,
		"status": "ACTIVE",
	}, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aCustomerSegmentCode := customerType["parent_value"].(string)

	if customerSegmentCode == "" {
		customerSegmentCode = aCustomerSegmentCode
	} else {
		if customerSegmentCode != aCustomerSegmentCode {
			return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "CUSTOMER_SEGMENT_CODE_NOT_MATCH_CUSTOMER_TYPE_CODE")
		}
	}

	_, _, err = master_data.ModuleMasterData.CustomerRef.ShouldSelectOne(&aepr.Log, utils.JSON{
		"type":   "CUSTOMER_SEGMENT",
		"code":   customerSegmentCode,
		"status": "ACTIVE",
	}, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	p := map[string]interface{}{}
	for k, v := range aepr.ParameterValues {
		p[k] = v.Value
	}

	if p["email"] == "" {
		p["email"] = nil
	}

	if p["phonenumber"] == "" {
		p["phonenumber"] = nil
	}

	id, err := task_management.ModuleTaskManagement.DoCustomerCreate(&aepr.Log, p)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"id": id,
	}})
	return nil
}

func CustomerEdit(aepr *api.DXAPIEndPointRequest) (err error) {
	_, dataNew, err := aepr.GetParameterValueAsJSON("new")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	customerSegmentCode, ok := dataNew["customer_segment_code"].(string)
	if ok {
		_, _, err = master_data.ModuleMasterData.CustomerRef.ShouldSelectOne(&aepr.Log, utils.JSON{
			"type":   "CUSTOMER_SEGMENT",
			"code":   customerSegmentCode,
			"status": "ACTIVE",
		}, nil, nil)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	}
	customerTypeCode, ok := dataNew["customer_type_code"].(string)
	if ok {
		_, _, err = master_data.ModuleMasterData.CustomerRef.ShouldSelectOne(&aepr.Log, utils.JSON{
			"type":         "CUSTOMER_TYPE",
			"code":         customerTypeCode,
			"parent_value": customerSegmentCode,
			"parent_group": "CUSTOMER_SEGMENT",
			"status":       "ACTIVE",
		}, nil, nil)
	}
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	err = task_management.ModuleTaskManagement.Customer.RequestEdit(aepr)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func CustomerEditByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	_, dataNew, err := aepr.GetParameterValueAsJSON("new")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	customerSegmentCode, ok := dataNew["customer_segment_code"].(string)
	if ok {
		_, _, err = master_data.ModuleMasterData.CustomerRef.ShouldSelectOne(&aepr.Log, utils.JSON{
			"type":   "CUSTOMER_SEGMENT",
			"code":   customerSegmentCode,
			"status": "ACTIVE",
		}, nil, nil)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	}
	customerTypeCode, ok := dataNew["customer_type_code"].(string)
	if ok {
		_, _, err = master_data.ModuleMasterData.CustomerRef.ShouldSelectOne(&aepr.Log, utils.JSON{
			"type":         "CUSTOMER_TYPE",
			"code":         customerTypeCode,
			"parent_value": customerSegmentCode,
			"parent_group": "CUSTOMER_SEGMENT",
			"status":       "ACTIVE",
		}, nil, nil)
	}
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	err = task_management.ModuleTaskManagement.Customer.RequestEditByUid(aepr)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func FieldExecutorCustomerEditByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.LocalData["user_id"].(int64)

	_, customerUid, err := aepr.GetParameterValueAsString("uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, dataNew, err := aepr.GetParameterValueAsJSON("new")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	customerSegmentCode, ok := dataNew["customer_segment_code"].(string)
	if ok {
		_, _, err = master_data.ModuleMasterData.CustomerRef.ShouldSelectOne(&aepr.Log, utils.JSON{
			"type":   "CUSTOMER_SEGMENT",
			"code":   customerSegmentCode,
			"status": "ACTIVE",
		}, nil, nil)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	}
	customerTypeCode, ok := dataNew["customer_type_code"].(string)
	if ok {
		_, _, err = master_data.ModuleMasterData.CustomerRef.ShouldSelectOne(&aepr.Log, utils.JSON{
			"type":         "CUSTOMER_TYPE",
			"code":         customerTypeCode,
			"parent_value": customerSegmentCode,
			"parent_group": "CUSTOMER_SEGMENT",
			"status":       "ACTIVE",
		}, nil, nil)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	}

	_, subTask, err := task_management.ModuleTaskManagement.SubTask.SelectOne(&aepr.Log, nil, utils.JSON{
		"customer_uid":                customerUid,
		"sub_task_type_full_code":     base.SubTaskTypeFullCodeConstructionSK,
		"c1":                          db.SQLExpression{"status in ('ASSIGNED', 'WORKING', 'FIXING', 'REWORKING')"},
		"last_field_executor_user_id": userId,
	}, nil, map[string]string{"id": "DESC"})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	if subTask == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "NO_VALID_SUB_TASK_SK_FOUND_FOR_CUSTOMER_UID_FOR_FIELD_EXECUTOR")
	}

	err = task_management.ModuleTaskManagement.Customer.RequestEditByUid(aepr)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func FieldExecutorCustomerEditLocationByUid(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.LocalData["user_id"].(int64)

	_, customerUid, err := aepr.GetParameterValueAsString("uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, subTask, err := task_management.ModuleTaskManagement.SubTask.SelectOne(&aepr.Log, nil, utils.JSON{
		"customer_uid":                customerUid,
		"sub_task_type_full_code":     base.SubTaskTypeFullCodeConstructionGasIn,
		"c1":                          db.SQLExpression{"status in ('ASSIGNED', 'WORKING', 'FIXING', 'REWORKING')"},
		"last_field_executor_user_id": userId,
	}, nil, map[string]string{"id": "DESC"})
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	if subTask == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "NO_VALID_SUB_TASK_GAS_IN_FOUND_FOR_CUSTOMER_UID_FOR_FIELD_EXECUTOR")
	}

	err = task_management.ModuleTaskManagement.Customer.RequestEditByUid(aepr)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func CustomerList(aepr *api.DXAPIEndPointRequest) (err error) {
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

	t := task_management.ModuleTaskManagement.Customer
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
		customerId, ok := row["id"].(int64)
		if !ok {
			return aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "CUSTOMER_ID_NOT_FOUND_IN_CUSTOMER_LIST")
		}
		_, tasks, err := task_management.ModuleTaskManagement.Task.Select(&aepr.Log, nil, utils.JSON{
			"customer_id": customerId,
		}, nil, map[string]string{"id": "DESC"}, nil)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		for _, task := range tasks {
			taskId, ok := task["id"].(int64)
			if !ok {
				return aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "TASK_ID_NOT_FOUND_IN_TASK_LIST")
			}
			_, subTasks, err := task_management.ModuleTaskManagement.SubTask.Select(&aepr.Log, nil, utils.JSON{
				"task_id": taskId,
			}, nil, map[string]string{"id": "DESC"}, nil)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			task["sub_tasks"] = subTasks
		}
		list[i]["tasks"] = tasks
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

func CustomerMeterList(aepr *api.DXAPIEndPointRequest) (err error) {
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

	t := task_management.ModuleTaskManagement.CustomerMeter

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

func CustomerPerConstructionAreaList(aepr *api.DXAPIEndPointRequest) (err error) {
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

	t := task_management.ModuleTaskManagement.Customer
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
		customerId, ok := row["id"].(int64)
		if !ok {
			return aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "CUSTOMER_ID_NOT_FOUND_IN_CUSTOMER_LIST")
		}
		_, tasks, err := task_management.ModuleTaskManagement.Task.Select(&aepr.Log, nil, utils.JSON{
			"customer_id": customerId,
		}, nil, map[string]string{"id": "DESC"}, nil)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		for _, task := range tasks {
			taskId, ok := task["id"].(int64)
			if !ok {
				return aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "TASK_ID_NOT_FOUND_IN_TASK_LIST")
			}
			_, subTasks, err := task_management.ModuleTaskManagement.SubTask.Select(&aepr.Log, nil, utils.JSON{
				"task_id": taskId,
			}, nil, map[string]string{"id": "DESC"}, nil)
			if err != nil {
				return errors.Wrap(err, "error occured")
			}
			task["sub_tasks"] = subTasks
		}
		list[i]["tasks"] = tasks
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

func CustomerListDownload(aepr *api.DXAPIEndPointRequest) (err error) {
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
	t := task_management.ModuleTaskManagement.Customer
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

	rowsInfo, list, err := db.NamedQueryList(t.Database.Connection, t.FieldTypeMapping, "*", t.ListViewNameId,
		filterWhere, "", filterOrderBy, filterKeyValues)

	if err != nil {
		return errors.Wrap(err, "error occurred")
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

func CustomerSoftDelete(aepr *api.DXAPIEndPointRequest) (err error) {
	_, customerId, err := aepr.GetParameterValueAsInt64("customer_id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	_, _, err = task_management.ModuleTaskManagement.Customer.ShouldGetById(&aepr.Log, customerId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	taskTotalRows, _, err := task_management.ModuleTaskManagement.Task.Count(&aepr.Log, "", utils.JSON{
		"customer_id": customerId,
	}, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if taskTotalRows > 0 {
		err = aepr.WriteResponseAndNewErrorf(402, "", "CUSTOMER_HAS_TASK:%d:%d", customerId, taskTotalRows)
		return errors.Wrap(err, "error occured")
	}

	err = task_management.ModuleTaskManagement.Customer.SoftDelete(aepr, customerId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
	return nil
}

func CustomerConstructionTaskCreateBulk(aepr *api.DXAPIEndPointRequest) (err error) {
	// Get the request body stream
	bs := aepr.Request.Body
	if bs == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "FAILED_TO_GET_BODY_STREAM:%s", "CustomerConstructionTaskCreateBulk")
	}
	defer bs.Close()

	// Read the entire request body into a buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, bs)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "FAILED_TO_READ_REQUEST_BODY:%s=%v", "CustomerConstructionTaskCreateBulk", err.Error())
	}

	// Determine the file type and parse accordingly
	contentType := aepr.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, "csv") {
		err = parseAndCreateConstructionTasksFromCSV(&buf, aepr)
	} else if strings.Contains(contentType, "excel") || strings.Contains(contentType, "spreadsheetml") {
		err = parseAndCreateConstructionTasksFromXLSX(&buf, aepr)
	} else {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnsupportedMediaType, "UNSUPPORTED_FILE_TYPE:%s", contentType)
	}

	if err != nil {
		return errors.Wrap(err, "error occurred")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
	return nil
}

func CustomerConstructionTaskCreateBulkBase64(aepr *api.DXAPIEndPointRequest) (err error) {
	_, fileContentBase64, err := aepr.GetParameterValueAsString("content_base64")
	if err != nil {
		return err
	}

	_, contentType, err := aepr.GetParameterValueAsString("content_type")
	if err != nil {
		return err
	}

	// Decode base64 content
	decodedBytes, err := base64.StdEncoding.DecodeString(fileContentBase64)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "INVALID_BASE64_CONTENT")
	}

	// Create a buffer with the decoded content
	var buf bytes.Buffer
	if _, err := buf.Write(decodedBytes); err != nil {
		return err
	}

	// Determine the file type and parse accordingly
	if strings.Contains(strings.ToLower(contentType), "csv") {
		err = parseAndCreateConstructionTasksFromCSV(&buf, aepr)
	} else if strings.Contains(strings.ToLower(contentType), "excel") || strings.Contains(strings.ToLower(contentType), "spreadsheetml") {
		err = parseAndCreateConstructionTasksFromXLSX(&buf, aepr)
	} else {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "UNSUPPORTED_FILE_FORMAT")
	}

	if err != nil {
		return errors.Wrap(err, "error occurred")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
	return nil
}

func parseAndCreateConstructionTasksFromCSV(buf *bytes.Buffer, aepr *api.DXAPIEndPointRequest) error {
	// Create a new reader with comma as delimiter
	reader := csv.NewReader(buf)
	reader.Comma = ';'
	reader.LazyQuotes = true    // Handle quotes more flexibly
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read header row
	headers, err := reader.Read()
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity,
			"FAILED_TO_READ_CSV_HEADERS: %s", err.Error())
	}

	// Clean headers - trim spaces and empty fields
	cleanHeaders := make([]string, 0)
	for _, h := range headers {
		h = strings.TrimSpace(h)
		if h != "" {
			cleanHeaders = append(cleanHeaders, h)
		}
	}

	// Process each row
	lineNum := 1 // Keep track of line numbers for error reporting
	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "",
				"FAILED_TO_PARSE_CSV_LINE_%d: %s", lineNum, err.Error())
		}

		// Create task data map
		taskData := make(map[string]interface{})
		for i, value := range record {
			if i >= len(cleanHeaders) {
				break
			}
			// Clean and validate the value
			value = strings.TrimSpace(value)
			if value != "" {
				taskData[cleanHeaders[i]] = value
			}
		}

		// Skip empty rows
		if len(taskData) == 0 {
			continue
		}

		// Ensure registration_number is present
		if _, ok := taskData["registration_number"]; !ok {
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "",
				"REGISTRATION_NUMBER_REQUIRED_AT_LINE_%d", lineNum)
		}

		// Create construction task
		_, err = task_management.ModuleTaskManagement.DoCustomerConstructionTaskCreate(&aepr.Log, taskData)
		if err != nil {
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "",
				"FAILED_TO_CREATE_CONSTRUCTION_TASK_LINE_%d: %s", lineNum, err.Error())
		}
	}

	return nil
}

func parseAndCreateConstructionTasksFromXLSX(buf *bytes.Buffer, aepr *api.DXAPIEndPointRequest) error {
	// Read the XLSX file
	xlsxFile, err := xlsx.OpenBinary(buf.Bytes())
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity,
			"FAILED_TO_READ_XLSX_FILE: %s", err.Error())
	}

	// Process each sheet
	for _, sheet := range xlsxFile.Sheets {
		if len(sheet.Rows) < 2 {
			continue // Skip empty sheets
		}

		// Get headers from first row
		headers := make([]string, 0)
		for _, cell := range sheet.Rows[0].Cells {
			header := strings.TrimSpace(cell.String())
			if header != "" {
				headers = append(headers, header)
			}
		}

		// Process each row
		for rowIdx, row := range sheet.Rows[1:] {
			taskData := make(map[string]interface{}, len(headers))

			if len(row.Cells) == 0 {
				continue // Skip empty rows
			}

			// Map cell values to headers
			for i, cell := range row.Cells {
				if i >= len(headers) {
					break
				}

				value := strings.TrimSpace(cell.String())
				if value == "" {
					continue // Skip empty values
				}

				// Try to convert numeric values
				if isNumericColumn(headers[i]) {
					if numVal, err := cell.Float(); err == nil {
						taskData[headers[i]] = numVal
					} else {
						return aepr.WriteResponseAndNewErrorf(
							http.StatusUnprocessableEntity, "",
							"INVALID_NUMERIC_VALUE_AT_ROW_%d_COLUMN_%s: %q",
							rowIdx+2,
							headers[i],
							value,
						)
					}
				} else {
					taskData[headers[i]] = value
				}
			}

			if len(taskData) == 0 {
				continue
			}

			// Ensure registration_number is present
			if _, ok := taskData["registration_number"]; !ok {
				return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "",
					"REGISTRATION_NUMBER_REQUIRED_AT_ROW_%d", rowIdx+2)
			}

			// Create construction task
			if _, err = task_management.ModuleTaskManagement.DoCustomerConstructionTaskCreate(&aepr.Log, taskData); err != nil {
				return aepr.WriteResponseAndNewErrorf(
					http.StatusUnprocessableEntity, "",
					"FAILED_TO_CREATE_CONSTRUCTION_TASK_AT_ROW_%d: %s",
					rowIdx+2,
					err.Error(),
				)
			}
		}
	}

	return nil
}
