package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/upload_data"
	"net/http"
	"strings"
	"sync"
	"time"
)

func Upload(aepr *api.DXAPIEndPointRequest) (err error) {
	_, fileContentBase64, err := aepr.GetParameterValueAsString("content_base64")
	if err != nil {
		return err
	}

	_, contentType, err := aepr.GetParameterValueAsString("content_type")
	if err != nil {
		return err
	}

	_, dataType, err := aepr.GetParameterValueAsString("data_type")
	if err != nil {
		return err
	}

	_, option, err := aepr.GetParameterValueAsString("option")
	if err != nil {
		option = ""
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
	rows := []utils.JSON{}
	if strings.Contains(strings.ToLower(contentType), "csv") {
		rows, err = upload_data.GetDataFromCSV(&buf, aepr)
	} else if strings.Contains(strings.ToLower(contentType), "excel") || strings.Contains(strings.ToLower(contentType), "spreadsheetml") {
		rows, err = upload_data.GetDataFromXLS(&buf, aepr)
	}
	jobID := generateJobID("upload_" + dataType)

	// Initial validation and count check

	if err != nil {
		aepr.WriteResponseAsJSON(http.StatusInternalServerError, nil, utils.JSON{
			"error":          "Failed to check pending rows",
			"reason_message": err.Error(),
		})
		return err
	}

	if len(rows) == 0 {
		response := utils.JSON{
			"message":    "No rows to process",
			"job_id":     jobID,
			"status":     "COMPLETED",
			"total_rows": 0,
		}
		aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": response})
		return nil
	}

	// Initialize job status
	jobStatus := &JobStatus{
		JobID:     jobID,
		Status:    "STARTED",
		TotalRows: len(rows),
		StartTime: time.Now(),
		CreatedBy: getUserFromRequest(aepr), // Helper function to get user info
	}

	// Store job status
	jobMutex.Lock()
	jobTracker[jobID] = jobStatus
	jobMutex.Unlock()

	// Return immediately with job ID
	response := utils.JSON{
		"job_id":           jobID,
		"job_status":       "STARTED",
		"total_rows":       len(rows),
		"check_status_url": fmt.Sprintf("/api/upload-status/%s", jobID),
		"estimated_time":   fmt.Sprintf("%d minutes", estimateProcessingTime(len(rows))),
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
		"data":           response,
		"reason_message": "Upload process started in background",
	})
	// Start background processing
	go func() {
		err2 := processInBackground(jobID, aepr, dataType, rows, option)
		if err2 != nil {
			log.Log.Errorf(err, "Error in upload process background.")
		}
	}()

	if err != nil {
		return errors.Wrap(err, "error occurred")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
	return nil
}

// Background processing function
func processInBackground(jobID string, aepr *api.DXAPIEndPointRequest, dataType string, rows []utils.JSON, option string) (err error) {
	fmt.Println("----start background processing")
	// Update status to PROCESSING
	updateJobStatus(jobID, func(status *JobStatus) {
		status.Status = "PROCESSING"
	})

	// Add timeout context for the background operation
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	// Process with worker pool
	err = processWithWorkerPool(ctx, jobID, aepr, dataType, rows, option)
	return err
}

// Worker pool processing
func processWithWorkerPool(ctx context.Context, jobID string, aepr *api.DXAPIEndPointRequest, dataType string, rows []utils.JSON, option string) (err error) {
	fmt.Println("----start processUserWithWorkerPool")
	var processingErrors []interface{}
	var errorMutex sync.Mutex
	var wg sync.WaitGroup

	// Worker count based on load
	workerCount, err := calculateOptimalWorkers(len(rows))
	if err != nil {
		return err
	}
	jobs := make(chan struct {
		dataType string
		rowData  utils.JSON
		rowId    int64
		option   string
	}, len(rows))

	// Start workers
	for w := 1; w <= workerCount; w++ {
		go func(workerID int) {
			for {
				select {
				case job, ok := <-jobs:
					if !ok {

						fmt.Println("==== !ok ====")

						return

					}
					uploadJob(ctx, aepr, job, jobID, &processingErrors, &errorMutex)
					wg.Done()

				case <-ctx.Done():
					return
				}
			}
		}(w)
	}

	// Send jobs

	for _, row := range rows {
		id, ok := row["row_no"].(int)
		if !ok {
			continue
		}

		wg.Add(1)
		select {
		case jobs <- struct {
			dataType string
			rowData  utils.JSON
			rowId    int64
			option   string
		}{
			dataType: dataType,
			rowData:  row,
			rowId:    int64(id),
			option:   option,
		}:
		case <-ctx.Done():
			wg.Done()
			break
		}
	}

	close(jobs)

	// Wait for completion
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Completed successfully
		updateJobStatus(jobID, func(status *JobStatus) {
			status.Status = "COMPLETED"
			endTime := time.Now()
			status.EndTime = &endTime
			if len(processingErrors) > 0 {
				status.Errors = getSampleErrors(processingErrors, 10)
			}
		})
	case <-ctx.Done():
		// Timed out
		updateJobStatus(jobID, func(status *JobStatus) {
			status.Status = "FAILED"
			status.Errors = append(status.Errors, "Processing timed out")
			endTime := time.Now()
			status.EndTime = &endTime
		})
	}

	return nil
}

// Process individual job with status updates
func uploadJob(ctx context.Context, aepr *api.DXAPIEndPointRequest, job struct {
	dataType string
	rowData  utils.JSON
	rowId    int64
	option   string
}, jobID string, processingErrors *[]interface{}, errorMutex *sync.Mutex) {

	fmt.Printf("%v", job.dataType)
	select {
	case <-ctx.Done():
		return
	default:
	}

	// Update job statistics
	updateJobStatus(jobID, func(status *JobStatus) {
		status.Processed++
	})
	fmt.Println("==================")
	fmt.Println(job.dataType)
	// Business logic
	var err2 error
	switch job.dataType {
	case "organization":
		_, err2 = upload_data.ModuleUploadData.DoUploadOrganizationCreate(&aepr.Log, job.rowData)
	case "user":
		_, err2 = upload_data.ModuleUploadData.DoUploadUserCreate(&aepr.Log, job.rowData)
	case "customer":
		_, err2 = upload_data.ModuleUploadData.DoUploadCustomerCreate(&aepr.Log, job.rowData, job.option)
	case "arrears":
		_, err2 = upload_data.ModuleUploadData.DoUploadArrearsCreate(&aepr.Log, job.rowData)
	default:
		err2 = fmt.Errorf("unsupported data type: %s", job.dataType)
	}
	if err2 != nil {
		errorMutex.Lock()
		//*processingErrors = append(*processingErrors, fmt.Errorf("business logic failed for row %d: %w", job.rowId, err2))
		errorMsg := err2.Error()
		if len(errorMsg) > 100 {
			errorMsg = errorMsg[:97] + "..."
		}
		*processingErrors = append(*processingErrors, utils.JSON{
			"row":     job.rowId,
			"error":   "business logic failed",
			"message": errorMsg,
		})
		errorMutex.Unlock()

		updateJobStatus(jobID, func(status *JobStatus) {
			status.Failed++
		})

	} else {
		updateJobStatus(jobID, func(status *JobStatus) {
			status.Succeeded++
		})

	}
}

func UploadList(aepr *api.DXAPIEndPointRequest) (err error) {
	_, dataType, err := aepr.GetParameterValueAsString("data_type")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	switch dataType {
	case "organization":
		err = upload_data.ModuleUploadData.Organization.RequestPagingList(aepr)
	case "user":
		err = upload_data.ModuleUploadData.User.RequestPagingList(aepr)
	case "customer":
		err = upload_data.ModuleUploadData.Customer.RequestPagingList(aepr)
	default:
		err = fmt.Errorf("unsupported data type: %s", dataType)
	}
	if err != nil {
		return err
	}
	return nil
}

func UploadEdit(aepr *api.DXAPIEndPointRequest) (err error) {
	_, dataType, err := aepr.GetParameterValueAsString("data_type")
	_, id, err := aepr.GetParameterValueAsInt64("id")
	_, dataNew, err := aepr.GetParameterValueAsJSON("new")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	// Remove nil values from dataNew
	cleanData := make(utils.JSON)
	for key, value := range dataNew {
		if value != nil {
			cleanData[key] = value
		}
	}
	switch dataType {
	case "organization":
		err = upload_data.ModuleUploadData.DoUploadOrganizationUpdate(&aepr.Log, id, cleanData)
	case "user":
		err = upload_data.ModuleUploadData.DoUploadUserUpdate(&aepr.Log, id, cleanData)
	case "customer":
		err = upload_data.ModuleUploadData.DoUploadCustomerUpdate(&aepr.Log, id, cleanData)
	default:
		err = fmt.Errorf("unsupported data type: %s", dataType)
	}
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, data, err4 := upload_data.ModuleUploadData.User.ShouldGetById(&aepr.Log, id)
	if err4 != nil {
		return errors.Wrap(err4, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": data})
	return nil
}
