package upload_data

import (
	"bytes"
	"encoding/csv"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/configuration"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
	"io"
	"strconv"
	"strings"
	"time"
)

func SplitAndTrimStrings(rolesStr string) []string {
	if rolesStr == "" {
		return []string{}
	}

	items := strings.Split(rolesStr, ",")
	result := make([]string, 0, len(items))

	for _, item := range items {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func SplitAndParseInt64s(data interface{}) []int64 {
	dataStr, ok := data.(string)
	if !ok || dataStr == "" {
		return []int64{}
	}

	items := strings.Split(dataStr, ",")
	result := make([]int64, 0, len(items))

	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}

		if val, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
			result = append(result, val)
		}
		// Silently skip invalid numbers, or you could log the error
	}

	return result
}

func GetDataFromXLS(buf *bytes.Buffer, aepr *api.DXAPIEndPointRequest) (r []utils.JSON, err error) {
	data := []utils.JSON{}
	xlFile, err := xlsx.OpenBinary(buf.Bytes())
	if err != nil {
		//return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "FAILED_TO_PARSE_XLSX: %s", err.Error())
	}
	for _, sheet := range xlFile.Sheets {
		if len(sheet.Rows) < 2 {
			//return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "XLSX_FILE_MUST_HAVE_HEADER_AND_DATA")
		}

		// Validate and extract headers
		headers := make([]string, 0, len(sheet.Rows[0].Cells))
		for _, cell := range sheet.Rows[0].Cells {
			header := strings.TrimSpace(cell.String())
			if header == "" {
				//return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "EMPTY_HEADER_NOT_ALLOWED")
			}
			header = strings.ReplaceAll(header, " ", "_")
			headers = append(headers, strings.ToLower(header))
		}

		// Process data rows
		for rowIdx, row := range sheet.Rows[1:] {
			rowData := make(map[string]interface{}, len(headers))

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
				rowData["row_status"] = "READY"
				if isNumericColumn(headers[i]) {
					rowData[headers[i]] = value
					if _, err := cell.Float(); err != nil {
						// customerData[headers[i]] = numVal
						rowData["row_status"] = "ERROR"
					}
				} else if isDateColumn(headers[i]) {
					if numVal, err := cell.Float(); err == nil {
						excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
						// Add the number of days to get the actual date
						date := excelEpoch.AddDate(0, 0, int(numVal))
						rowData[headers[i]] = date.Format("2006-01-02")
						//customerData["row_status"] = "ERROR"
					}

					//layout := "2006-01-02"
					//if !isDate(value, layout) {
					//	customerData["row_status"] = "ERROR"
					//}
					//customerData[headers[i]] = value
				} else {
					rowData[headers[i]] = value
				}
				rowData["row_no"] = rowIdx + 1
			}

			if len(rowData) == 0 {
				continue
			}
			data = append(data, rowData)
		}
	}
	return data, nil
}

func GetDataFromCSV(buf *bytes.Buffer, aepr *api.DXAPIEndPointRequest) (r []utils.JSON, err error) {
	data := []utils.JSON{}
	// Create a new reader with comma as delimiter
	reader := csv.NewReader(buf)
	reader.Comma = ';'
	reader.LazyQuotes = true    // Handle quotes more flexibly
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read header row
	headers, err := reader.Read()
	if err != nil {
		//return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity,
		//	"FAILED_TO_READ_CSV_HEADERS: %s", err.Error())
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
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			//return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "",
			//	"FAILED_TO_PARSE_CSV_LINE_%d: %s", lineNum, err.Error())
		}

		// Create customer data map
		rowData := make(map[string]interface{})
		rowData["row_status"] = "READY"
		for i, value := range record {
			if i >= len(cleanHeaders) {
				break
			}
			// Clean and validate the value
			value = strings.TrimSpace(value)
			if value != "" {
				rowData[cleanHeaders[i]] = value
			}
		}
		rowData["row_no"] = lineNum

		// Skip empty rows
		if len(rowData) == 0 {
			continue
		}
		data = append(data, rowData)
		lineNum++
	}
	return data, nil
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

func isDateColumn(header string) bool {
	// Add your numeric column names here
	dateColumns := map[string]bool{
		"periode_awal_tunggakan":  true,
		"periode_akhir_tunggakan": true,
		"tgl_surat":               true,
		// Add other numeric column names as needed
	}

	header = strings.ToLower(header)
	return dateColumns[header]
}

func isDate(str string, layout string) bool {
	_, err := time.Parse(layout, str)
	return err == nil
}

// setErrorMessage sets the row_status to ERROR and appends or creates an error message
func setErrorMessage(data utils.JSON, message string) {
	data["row_status"] = "ERROR"
	if str, ok := data["row_message"].(string); ok {
		data["row_message"] = str + "\r" + message
	} else {
		data["row_message"] = message
	}
}

func GetMaxWorkers() (v int, err error) {
	configUpload := *configuration.Manager.Configurations["upload"].Data
	workerConfiguration, ok := configUpload["WORKER"].(utils.JSON)
	if !ok {
		return 0, errors.Errorf("GET_CONFIGURATION:UPLOAD_WORKER_CONFIG_NOT_FOUND")
	}
	processWorkerCountMax, ok := workerConfiguration["process_worker_count_max"].(int)
	if !ok {
		return 0, errors.Errorf("GET_CONFIGURATION:UPLOAD_WORKER_COUNT_MAX_NOT_FOUND")
	}
	return processWorkerCountMax, nil
}

func GetMinWorkers() (v int, err error) {
	configUpload := *configuration.Manager.Configurations["upload"].Data
	workerConfiguration, ok := configUpload["WORKER"].(utils.JSON)
	if !ok {
		return 0, errors.Errorf("GET_CONFIGURATION:UPLOAD_WORKER_CONFIG_NOT_FOUND")
	}
	processWorkerCountMin, ok := workerConfiguration["process_worker_count_min"].(int)
	if !ok {
		return 0, errors.Errorf("GET_CONFIGURATION:UPLOAD_WORKER_COUNT_MIN_NOT_FOUND")
	}
	return processWorkerCountMin, nil
}
