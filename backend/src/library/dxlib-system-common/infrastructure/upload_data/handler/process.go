package handler

import (
	"context"
	"fmt"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/upload_data"
	"net/http"
	"sync"
	"time"
)

// UploadProcess now returns immediately with a job ID
func Process(aepr *api.DXAPIEndPointRequest) (err error) {
	_, dataType, err := aepr.GetParameterValueAsString("data_type")
	if err != nil {
		return err
	}
	// Generate unique job ID
	jobID := generateJobID("process_" + dataType)

	// Initial validation and count check
	rows := []utils.JSON{}

	switch dataType {
	case "organization":
		_, rows, err = upload_data.ModuleUploadData.Organization.Select(&aepr.Log,
			[]string{"id"},
			utils.JSON{"row_status": "READY"}, nil, map[string]string{"parent_code": "DESC"}, nil) // DESC --> NULL first
	case "user":
		_, rows, err = upload_data.ModuleUploadData.User.Select(&aepr.Log,
			[]string{"id"},
			utils.JSON{"row_status": "READY"}, nil, nil, nil)
	case "customer":
		_, rows, err = upload_data.ModuleUploadData.Customer.Select(&aepr.Log,
			[]string{"id"},
			utils.JSON{"row_status": "READY"}, nil, nil, nil)
	case "arrears":
		_, rows, err = upload_data.ModuleUploadData.Arrears.Select(&aepr.Log,
			[]string{"id"},
			utils.JSON{"row_status": "READY"}, nil, nil, nil)
	default:
		aepr.WriteResponseAsJSON(http.StatusBadRequest, nil, utils.JSON{
			"error":          "Invalid data type",
			"reason_message": fmt.Sprintf("Unsupported dataType: %s", dataType),
		})
		return fmt.Errorf("invalid data type: %s", dataType)
	}

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
	go processAddUserInBackground(jobID, aepr, dataType)

	return nil
}

// Background processing function
func processAddUserInBackground(jobID string, aepr *api.DXAPIEndPointRequest, dataType string) {
	// Update status to PROCESSING
	updateJobStatus(jobID, func(status *JobStatus) {
		status.Status = "PROCESSING"
	})

	// Add timeout context for the background operation
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()
	var err error
	// Update READY to QUEUED
	switch dataType {
	case "organization":
		_, err = upload_data.ModuleUploadData.Organization.Update(utils.JSON{
			"row_status": "QUEUED",
		}, utils.JSON{
			"row_status": "READY",
		})
	case "user":
		_, err = upload_data.ModuleUploadData.User.Update(utils.JSON{
			"row_status": "QUEUED",
		}, utils.JSON{
			"row_status": "READY",
		})
	case "customer":
		_, err = upload_data.ModuleUploadData.Customer.Update(utils.JSON{
			"row_status": "QUEUED",
		}, utils.JSON{
			"row_status": "READY",
		})
	case "arrears":
		_, err = upload_data.ModuleUploadData.Arrears.Update(utils.JSON{
			"row_status": "QUEUED",
		}, utils.JSON{
			"row_status": "READY",
		})
	}
	if err != nil {
		updateJobStatus(jobID, func(status *JobStatus) {
			status.Status = "FAILED"
			status.Errors = append(status.Errors, fmt.Sprintf("Failed to update initial status: %v", err))
			endTime := time.Now()
			status.EndTime = &endTime
		})
		return
	}

	// Get full data for processing
	rows := []utils.JSON{}
	condition := utils.JSON{"row_status": "QUEUED"}

	switch dataType {
	case "organization":
		_, rows, err = upload_data.ModuleUploadData.Organization.Select(&aepr.Log, nil, condition, nil, nil, nil)
	case "user":
		_, rows, err = upload_data.ModuleUploadData.User.Select(&aepr.Log, nil, condition, nil, nil, nil)
	case "customer":
		_, rows, err = upload_data.ModuleUploadData.Customer.Select(&aepr.Log, nil, condition, nil, nil, nil)
	case "arrears":
		_, rows, err = upload_data.ModuleUploadData.Arrears.Select(&aepr.Log, nil, condition, nil, nil, nil)
	default:
		err = fmt.Errorf("unknown data type: %s", dataType)
	}
	if err != nil {
		updateJobStatus(jobID, func(status *JobStatus) {
			status.Status = "FAILED"
			status.Errors = append(status.Errors, fmt.Sprintf("Failed to select rows: %v", err))
			endTime := time.Now()
			status.EndTime = &endTime
		})
		return
	}

	// Process with worker pool
	processProcessWithWorkerPool(ctx, jobID, aepr, dataType, rows)
}

// Worker pool processing
func processProcessWithWorkerPool(ctx context.Context, jobID string, aepr *api.DXAPIEndPointRequest, dataType string, rows []utils.JSON) (err error) {
	var processingErrors []interface{}
	var errorMutex sync.Mutex
	var wg sync.WaitGroup
	var workerCount int

	if dataType == "organization" {
		workerCount = 1

		// For organizations, process parent jobs first (where parent_code is null)
		// and then process child jobs
		var parentRows []utils.JSON
		var childRows []utils.JSON

		for _, row := range rows {
			parentCode, hasParentCode := row["parent_code"]
			if !hasParentCode || parentCode == nil || parentCode == "" {
				parentRows = append(parentRows, row)
				fmt.Printf("Added parent row with ID: %v\n", row["id"])
			} else {
				childRows = append(childRows, row)
				fmt.Printf("Added child row with ID: %v, parent_code: %v\n", row["id"], parentCode)
			}
		}

		fmt.Printf("Processing %d parent rows first, followed by %d child rows\n", len(parentRows), len(childRows))

		// Process parent jobs first
		processJobBatch(ctx, jobID, aepr, dataType, parentRows, &wg, workerCount, &processingErrors, &errorMutex)

		// Then process child jobs
		processJobBatch(ctx, jobID, aepr, dataType, childRows, &wg, workerCount, &processingErrors, &errorMutex)
	} else {
		// Worker count based on load
		workerCount, err = calculateOptimalWorkers(len(rows))
		if err != nil {
			return err
		}
		// For other data types, process all jobs together
		processJobBatch(ctx, jobID, aepr, dataType, rows, &wg, workerCount, &processingErrors, &errorMutex)
	}

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

// Helper function to process a batch of jobs
func processJobBatch(ctx context.Context, jobID string, aepr *api.DXAPIEndPointRequest, dataType string, rows []utils.JSON, wg *sync.WaitGroup, workerCount int, processingErrors *[]interface{}, errorMutex *sync.Mutex) {
	if len(rows) == 0 {
		fmt.Printf("Skipping empty batch for dataType: %s\n", dataType)
		return
	}

	// Log the start of batch processing
	fmt.Printf("Starting to process batch of %d rows for dataType: %s\n", len(rows), dataType)

	jobs := make(chan struct {
		dataType string
		rowData  utils.JSON
		rowId    int64
	}, len(rows))

	// Start workers
	for w := 1; w <= workerCount; w++ {
		go func(workerID int) {
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return
					}

					processJobWithStatusUpdate(ctx, aepr, job, jobID, processingErrors, errorMutex)
					wg.Done()

				case <-ctx.Done():
					return
				}
			}
		}(w)
	}

	// Send jobs
	for _, row := range rows {
		id, ok := row["id"].(int64)
		if !ok {
			continue
		}

		wg.Add(1)
		select {
		case jobs <- struct {
			dataType string
			rowData  utils.JSON
			rowId    int64
		}{
			dataType: dataType,
			rowData:  row,
			rowId:    id,
		}:
		case <-ctx.Done():
			wg.Done()
			break
		}
	}

	close(jobs)

	// Create a channel to wait for this batch to complete
	batchDone := make(chan struct{})
	go func() {
		// Create a separate WaitGroup for this batch
		var batchWg sync.WaitGroup
		batchWg.Add(1)

		// Start a goroutine to decrement the batch WaitGroup when all jobs are done
		go func() {
			// Wait for the main WaitGroup to reach its previous count
			// This is a bit of a hack, but we can't directly access the WaitGroup's counter
			time.Sleep(100 * time.Millisecond) // Give time for jobs to be added to the WaitGroup

			// Poll until all jobs in this batch are done
			for {
				select {
				case <-time.After(100 * time.Millisecond):
					// Check if all jobs have been processed
					// We can't directly check the WaitGroup, so we rely on the jobs channel being closed
					// and all workers having processed their jobs
					if len(jobs) == 0 {
						batchWg.Done()
						return
					}
				case <-ctx.Done():
					batchWg.Done()
					return
				}
			}
		}()

		batchWg.Wait()
		close(batchDone)
	}()

	// Wait for this batch to complete before returning
	select {
	case <-batchDone:
		// Batch completed
		fmt.Printf("Completed processing batch of %d rows for dataType: %s\n", len(rows), dataType)
		return
	case <-ctx.Done():
		// Context cancelled
		fmt.Printf("Context cancelled while processing batch for dataType: %s\n", dataType)
		return
	}
}

// Process individual job with status updates
func processJobWithStatusUpdate(ctx context.Context, aepr *api.DXAPIEndPointRequest, job struct {
	dataType string
	rowData  utils.JSON
	rowId    int64
}, jobID string, processingErrors *[]interface{}, errorMutex *sync.Mutex) {

	select {
	case <-ctx.Done():
		return
	default:
	}

	// Update job statistics
	updateJobStatus(jobID, func(status *JobStatus) {
		status.Processed++
	})

	// Update to PROCESSING
	var err1 error

	switch job.dataType {
	case "user":
		_, err1 = upload_data.ModuleUploadData.User.UpdateOne(&aepr.Log, job.rowId, utils.JSON{
			"row_status": "PROCESSING",
		})
	case "organization":
		_, err1 = upload_data.ModuleUploadData.Organization.UpdateOne(&aepr.Log, job.rowId, utils.JSON{
			"row_status": "PROCESSING",
		})
	case "customer":
		_, err1 = upload_data.ModuleUploadData.Customer.UpdateOne(&aepr.Log, job.rowId, utils.JSON{
			"row_status": "PROCESSING",
		})
	case "arrears":
		_, err1 = upload_data.ModuleUploadData.Arrears.UpdateOne(&aepr.Log, job.rowId, utils.JSON{
			"row_status": "PROCESSING",
		})
	default:
		err1 = fmt.Errorf("unsupported data type: %s", job.dataType)
	}

	if err1 != nil {
		errorMutex.Lock()
		//*processingErrors = append(*processingErrors, fmt.Errorf("failed to update status for row %d: %w", job.rowId, err1))
		*processingErrors = append(*processingErrors, utils.JSON{
			"error":   "business logic failed",
			"row":     job.rowId,
			"message": err1.Error(),
		})
		errorMutex.Unlock()

		updateJobStatus(jobID, func(status *JobStatus) {
			status.Failed++
		})
		return
	}

	// Business logic
	//err2, subTaskId := arrears_management.ModuleArrearsManagement.DoTaskArrearsCreate(&aepr.Log, job.rowData)
	var err error
	if job.dataType == "organization" {
		err, _ = upload_data.ModuleUploadData.DoOrganizationCreate(&aepr.Log, job.rowData)
	}
	if job.dataType == "user" {
		err, _ = upload_data.ModuleUploadData.DoUserCreate(&aepr.Log, job.rowData)
	}
	if job.dataType == "customer" {
		err, _ = upload_data.ModuleUploadData.DoCustomerCreate(&aepr.Log, job.rowData)
	}
	if job.dataType == "arrears" {
		err, _ = upload_data.ModuleUploadData.DoArrearsCreate(&aepr.Log, job.rowData)
	}

	if err != nil {
		fmt.Println("=== error create data ===")
		errorMutex.Lock()
		//*processingErrors = append(*processingErrors, fmt.Errorf("business logic failed for row %d: %w", job.rowId, err))
		*processingErrors = append(*processingErrors, utils.JSON{
			"row":     job.rowId,
			"error":   "business logic failed",
			"message": err.Error(),
		})
		errorMutex.Unlock()

		updateJobStatus(jobID, func(status *JobStatus) {
			status.Failed++
		})
		originalError := err.Error()

		switch job.dataType {
		case "organization":
			_, err = upload_data.ModuleUploadData.Organization.UpdateOne(&aepr.Log, job.rowId, utils.JSON{
				"row_status":  "FAILED",
				"row_message": originalError,
			})
		case "user":
			_, err = upload_data.ModuleUploadData.User.UpdateOne(&aepr.Log, job.rowId, utils.JSON{
				"row_status":  "FAILED",
				"row_message": originalError,
			})
		case "customer":
			_, err = upload_data.ModuleUploadData.Customer.UpdateOne(&aepr.Log, job.rowId, utils.JSON{
				"row_status":  "FAILED",
				"row_message": originalError,
			})
		case "arrears":
			_, err = upload_data.ModuleUploadData.Arrears.UpdateOne(&aepr.Log, job.rowId, utils.JSON{
				"row_status":  "FAILED",
				"row_message": originalError,
			})
		default:
			err = fmt.Errorf("unsupported data type: %s", job.dataType)
		}
		if err != nil {
			return
		}
	} else {
		updateJobStatus(jobID, func(status *JobStatus) {
			status.Succeeded++
		})
		switch job.dataType {
		case "organization":
			_, err = upload_data.ModuleUploadData.Organization.UpdateOne(&aepr.Log, job.rowId, utils.JSON{
				"row_status": "UPLOADED",
			})
		case "user":
			_, err = upload_data.ModuleUploadData.User.UpdateOne(&aepr.Log, job.rowId, utils.JSON{
				"row_status": "UPLOADED",
			})
		case "customer":
			_, err = upload_data.ModuleUploadData.Customer.UpdateOne(&aepr.Log, job.rowId, utils.JSON{
				"row_status": "UPLOADED",
			})
		case "arrears":
			_, err = upload_data.ModuleUploadData.Arrears.UpdateOne(&aepr.Log, job.rowId, utils.JSON{
				"row_status": "UPLOADED",
			})
		default:
			err = fmt.Errorf("unsupported data type: %s", job.dataType)
		}
		if err != nil {
			return
		}
	}
}
