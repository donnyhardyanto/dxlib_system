package handler

import (
	"fmt"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/database/protected/db"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"net/http"
)

//func CustomerCreateBulkBase64(aepr *api.DXAPIEndPointRequest) (err error) {
//	_, fileContentBase64, err := aepr.GetParameterValueAsString("content_base64")
//	if err != nil {
//		return err
//	}
//
//	_, contentType, err := aepr.GetParameterValueAsString("content_type")
//	if err != nil {
//		return err
//	}
//
//	// Decode base64 content
//	decodedBytes, err := base64.StdEncoding.DecodeString(fileContentBase64)
//	if err != nil {
//		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "INVALID_BASE64_CONTENT")
//	}
//
//	// Create a buffer with the decoded content
//	var buf bytes.Buffer
//	if _, err := buf.Write(decodedBytes); err != nil {
//		return err
//	}
//
//	// Determine the file type and parse accordingly
//	if strings.Contains(strings.ToLower(contentType), "csv") {
//		err = parseAndCreateCustomersFromCSV(&buf, aepr)
//	} else if strings.Contains(strings.ToLower(contentType), "excel") || strings.Contains(strings.ToLower(contentType), "spreadsheetml") {
//		err = parseAndCreateCustomersFromXLSX(&buf, aepr)
//	} else {
//		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "UNSUPPORTED_FILE_FORMAT")
//	}
//
//	if err != nil {
//		return errors.Wrap(err, "error occurred")
//	}
//
//	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
//	return nil
//}
//
//func parseAndCreateCustomersFromCSV(buf *bytes.Buffer, aepr *api.DXAPIEndPointRequest) error {
//	// Create a new reader with comma as delimiter
//	reader := csv.NewReader(buf)
//	reader.Comma = ';'
//	reader.LazyQuotes = true    // Handle quotes more flexibly
//	reader.FieldsPerRecord = -1 // Allow variable number of fields
//
//	// Read header row
//	headers, err := reader.Read()
//	if err != nil {
//		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity,
//			"FAILED_TO_READ_CSV_HEADERS: %s", err.Error())
//	}
//
//	// Clean headers - trim spaces and empty fields
//	cleanHeaders := make([]string, 0)
//	for _, h := range headers {
//		h = strings.TrimSpace(h)
//		if h != "" {
//			cleanHeaders = append(cleanHeaders, h)
//		}
//	}
//
//	// Process each row
//	lineNum := 1 // Keep track of line numbers for error reporting
//	for {
//		lineNum++
//		record, err := reader.Read()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "",
//				"FAILED_TO_PARSE_CSV_LINE_%d: %s", lineNum, err.Error())
//		}
//
//		// Create customer data map
//		customerData := make(map[string]interface{})
//		for i, value := range record {
//			if i >= len(cleanHeaders) {
//				break
//			}
//			// Clean and validate the value
//			value = strings.TrimSpace(value)
//			if value != "" {
//				customerData[cleanHeaders[i]] = value
//			}
//		}
//
//		// Skip empty rows
//		if len(customerData) == 0 {
//			continue
//		}
//
//		// Create customer
//		_, err = task_management.ModuleTaskManagement.DoCustomerCreate(&aepr.Log, customerData)
//		if err != nil {
//			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "",
//				"FAILED_TO_CREATE_CUSTOMER_LINE_%d: %s", lineNum, err.Error())
//		}
//	}
//
//	return nil
//}
//
//func parseAndCreateCustomersFromXLSX(buf *bytes.Buffer, aepr *api.DXAPIEndPointRequest) error {
//	xlFile, err := xlsx.OpenBinary(buf.Bytes())
//	if err != nil {
//		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "FAILED_TO_PARSE_XLSX: %s", err.Error())
//	}
//
//	for _, sheet := range xlFile.Sheets {
//		if len(sheet.Rows) < 2 {
//			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "XLSX_FILE_MUST_HAVE_HEADER_AND_DATA")
//		}
//
//		// Validate and extract headers
//		headers := make([]string, 0, len(sheet.Rows[0].Cells))
//		for _, cell := range sheet.Rows[0].Cells {
//			header := strings.TrimSpace(cell.String())
//			if header == "" {
//				return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "EMPTY_HEADER_NOT_ALLOWED")
//			}
//			headers = append(headers, header)
//		}
//
//		// Process data rows
//		for rowIdx, row := range sheet.Rows[1:] {
//			customerData := make(map[string]interface{}, len(headers))
//
//			if len(row.Cells) == 0 {
//				continue // Skip empty rows
//			}
//
//			// Map cell values to headers with type conversion
//			for i, cell := range row.Cells {
//				if i >= len(headers) {
//					break
//				}
//
//				value := strings.TrimSpace(cell.String())
//				if value == "" {
//					continue // Skip empty values instead of adding them to customerData
//				}
//
//				// Try to convert numeric values
//				if upload_data.isNumericColumn(headers[i]) {
//					if numVal, err := cell.Float(); err == nil {
//						customerData[headers[i]] = numVal
//					} else {
//						return aepr.WriteResponseAndNewErrorf(
//							http.StatusUnprocessableEntity, "",
//							"INVALID_NUMERIC_VALUE_AT_ROW_%d_COLUMN_%s: %q",
//							rowIdx+2,
//							headers[i],
//							value,
//						)
//					}
//				} else {
//					customerData[headers[i]] = value
//				}
//			}
//
//			if len(customerData) == 0 {
//				continue
//			}
//
//			if _, err = task_management.ModuleTaskManagement.DoCustomerCreate(&aepr.Log, customerData); err != nil {
//				// Check for specific PostgreSQL errors
//				if strings.Contains(err.Error(), "invalid input syntax for type double precision") {
//					return aepr.WriteResponseAndNewErrorf(
//						http.StatusUnprocessableEntity, "",
//						"INVALID_NUMERIC_VALUE_AT_ROW_%d: Please ensure all numeric fields contain valid numbers",
//						rowIdx+2,
//					)
//				}
//				return aepr.WriteResponseAndNewErrorf(
//					http.StatusUnprocessableEntity, "",
//					"FAILED_TO_CREATE_CUSTOMER_AT_ROW_%d: %s",
//					rowIdx+2,
//					err.Error(),
//				)
//			}
//		}
//	}
//
//	return nil
//}

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
