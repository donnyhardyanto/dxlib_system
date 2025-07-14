package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/upload_data"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

// JobStatus represents the status of a background job
type JobStatus struct {
	JobID     string     `json:"job_id"`
	Status    string     `json:"status"` // STARTED, PROCESSING, COMPLETED, FAILED
	TotalRows int        `json:"total_rows"`
	Processed int64      `json:"processed"`
	Succeeded int64      `json:"succeeded"`
	Failed    int64      `json:"failed"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Errors    []string   `json:"errors,omitempty"`
	CreatedBy string     `json:"created_by,omitempty"`
}

// JobSummary represents a summary view of job status
type JobSummary struct {
	JobID        string    `json:"job_id"`
	Status       string    `json:"status"`
	Progress     string    `json:"progress"`
	TotalRows    int       `json:"total_rows"`
	Processed    int64     `json:"processed"`
	Succeeded    int64     `json:"succeeded"`
	Failed       int64     `json:"failed"`
	StartTime    time.Time `json:"start_time"`
	Duration     string    `json:"duration"`
	CreatedBy    string    `json:"created_by"`
	EstimatedETA string    `json:"estimated_eta,omitempty"`
}

// Global job tracker - in production, use Redis or database
var (
	jobTracker = make(map[string]*JobStatus)
	jobMutex   sync.RWMutex
)

// Helper functions
func generateJobID(prefix ...string) string {
	bytes := make([]byte, 8)
	rand.Read(bytes)

	// Use default prefix if none provided
	prefixStr := "upload"
	if len(prefix) > 0 && prefix[0] != "" {
		prefixStr = prefix[0]
	}
	return fmt.Sprintf("%s_%s_%d", prefixStr, hex.EncodeToString(bytes), time.Now().Unix())
}

func updateJobStatus(jobID string, updateFunc func(*JobStatus)) {
	jobMutex.Lock()
	defer jobMutex.Unlock()

	if status, exists := jobTracker[jobID]; exists {
		updateFunc(status)
	}
}

// Default maximum number of workers for small workloads

func calculateOptimalWorkers(rowCount int) (int, error) {
	if rowCount <= 50 {
		return upload_data.GetMinWorkers()
	} else if rowCount <= 500 {
		return 50, nil
	}
	return upload_data.GetMaxWorkers()
}

func estimateProcessingTime(rowCount int) int {
	// Rough estimate: 100 rows per minute
	return (rowCount / 100) + 1
}

func getUserFromRequest(aepr *api.DXAPIEndPointRequest) string {
	// Implement based on your auth system
	return "system" // placeholder
}

func getSampleErrors(errors []interface{}, maxCount int) []string {
	var samples []string
	count := len(errors)
	if count > maxCount {
		count = maxCount
	}

	for i := 0; i < count; i++ {
		// Convert the error object to JSON without escaping
		if jsonBytes, err := json.Marshal(errors[i]); err == nil {
			samples = append(samples, string(jsonBytes))
		} else {
			samples = append(samples, fmt.Sprintf("%v", errors[i]))
		}
	}

	return samples
}

// GetActiveJobs returns all currently active (non-completed) jobs
func GetActiveJobs(aepr *api.DXAPIEndPointRequest) error {
	jobMutex.RLock()
	defer jobMutex.RUnlock()

	var activeJobs []JobSummary

	for _, status := range jobTracker {
		// Only include active jobs (not completed or failed)
		if status.Status == "STARTED" || status.Status == "PROCESSING" {
			progress := 0.0
			if status.TotalRows > 0 {
				progress = float64(status.Processed) / float64(status.TotalRows) * 100
			}

			// Calculate estimated ETA for processing jobs
			estimatedETA := ""
			if status.Status == "PROCESSING" && status.Processed > 0 {
				elapsed := time.Since(status.StartTime)
				avgTimePerRow := elapsed / time.Duration(status.Processed)
				remainingRows := int64(status.TotalRows) - status.Processed
				estimatedRemaining := time.Duration(remainingRows) * avgTimePerRow
				estimatedETA = estimatedRemaining.Round(time.Minute).String()
			}

			activeJobs = append(activeJobs, JobSummary{
				JobID:        status.JobID,
				Status:       status.Status,
				Progress:     fmt.Sprintf("%.1f%%", progress),
				TotalRows:    status.TotalRows,
				Processed:    status.Processed,
				Succeeded:    status.Succeeded,
				Failed:       status.Failed,
				StartTime:    status.StartTime,
				Duration:     time.Since(status.StartTime).Round(time.Second).String(),
				CreatedBy:    status.CreatedBy,
				EstimatedETA: estimatedETA,
			})
		}
	}

	// Sort by start time (newest first)
	sort.Slice(activeJobs, func(i, j int) bool {
		return activeJobs[i].StartTime.After(activeJobs[j].StartTime)
	})

	response := utils.JSON{
		"active_jobs": activeJobs,
		"count":       len(activeJobs),
		"timestamp":   time.Now(),
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": response})
	return nil
}

// Get job status endpoint
func GetJobStatus(aepr *api.DXAPIEndPointRequest) (err error) {
	_, jobID, err := aepr.GetParameterValueAsString("job_id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	jobMutex.RLock()
	status, exists := jobTracker[jobID]
	jobMutex.RUnlock()

	if !exists {
		aepr.WriteResponseAsJSON(http.StatusNotFound, nil, utils.JSON{
			"data": utils.JSON{
				"job_id": jobID,
			},
			"reason_message": "Job not found",
		})
		return nil
	}

	// Calculate progress percentage
	progress := 0.0
	if status.TotalRows > 0 {
		progress = float64(status.Processed) / float64(status.TotalRows) * 100
	}

	response := utils.JSON{
		"job_id":     status.JobID,
		"status":     status.Status,
		"progress":   fmt.Sprintf("%.1f%%", progress),
		"total_rows": status.TotalRows,
		"processed":  status.Processed,
		"succeeded":  status.Succeeded,
		"failed":     status.Failed,
		"start_time": status.StartTime,
		"duration":   time.Since(status.StartTime).String(),
	}

	if status.EndTime != nil {
		response["end_time"] = *status.EndTime
		response["total_duration"] = status.EndTime.Sub(status.StartTime).String()
	}

	if len(status.Errors) > 0 {
		response["errors"] = status.Errors
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": response})
	return nil
}

// GetAllJobs returns all jobs with optional filtering
func GetAllJobs(aepr *api.DXAPIEndPointRequest) (err error) {
	// Get query parameters for filtering
	_, statusFilter, err := aepr.GetParameterValueAsString("status") // e.g., "PROCESSING", "COMPLETED", "FAILED"
	_, userFilter, err := aepr.GetParameterValueAsString("user")
	_, limitStr, err := aepr.GetParameterValueAsString("limit")
	_, offsetStr, err := aepr.GetParameterValueAsString("offset")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	// Parse pagination parameters
	limit := 50 // default limit
	offset := 0 // default offset
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	jobMutex.RLock()
	defer jobMutex.RUnlock()

	var filteredJobs []JobSummary

	for _, status := range jobTracker {
		// Apply filters
		if statusFilter != "" && status.Status != statusFilter {
			continue
		}
		if userFilter != "" && status.CreatedBy != userFilter {
			continue
		}

		progress := 0.0
		if status.TotalRows > 0 {
			progress = float64(status.Processed) / float64(status.TotalRows) * 100
		}

		duration := time.Since(status.StartTime)
		if status.EndTime != nil {
			duration = status.EndTime.Sub(status.StartTime)
		}

		// Calculate ETA for active jobs
		estimatedETA := ""
		if (status.Status == "PROCESSING") && status.Processed > 0 {
			elapsed := time.Since(status.StartTime)
			avgTimePerRow := elapsed / time.Duration(status.Processed)
			remainingRows := int64(status.TotalRows) - status.Processed
			estimatedRemaining := time.Duration(remainingRows) * avgTimePerRow
			estimatedETA = estimatedRemaining.Round(time.Minute).String()
		}

		filteredJobs = append(filteredJobs, JobSummary{
			JobID:        status.JobID,
			Status:       status.Status,
			Progress:     fmt.Sprintf("%.1f%%", progress),
			TotalRows:    status.TotalRows,
			Processed:    status.Processed,
			Succeeded:    status.Succeeded,
			Failed:       status.Failed,
			StartTime:    status.StartTime,
			Duration:     duration.Round(time.Second).String(),
			CreatedBy:    status.CreatedBy,
			EstimatedETA: estimatedETA,
		})
	}

	// Sort by start time (newest first)
	sort.Slice(filteredJobs, func(i, j int) bool {
		return filteredJobs[i].StartTime.After(filteredJobs[j].StartTime)
	})

	// Apply pagination
	total := len(filteredJobs)
	start := offset
	end := offset + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedJobs := filteredJobs[start:end]

	response := utils.JSON{
		"list": utils.JSON{
			"rows":       paginatedJobs,
			"total_page": total,
			"total_row":  len(paginatedJobs),
			"total":      total,
			"limit":      limit,
			"offset":     offset,
			"count":      len(paginatedJobs),
		},
		"filters": utils.JSON{
			"status": statusFilter,
			"user":   userFilter,
		},
		"timestamp": time.Now(),
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
		"data": response})
	return nil
}

// GetJobStats returns aggregate statistics for all jobs
func GetJobStats(aepr *api.DXAPIEndPointRequest) error {
	jobMutex.RLock()
	defer jobMutex.RUnlock()

	stats := struct {
		Total     int            `json:"total"`
		ByStatus  map[string]int `json:"by_status"`
		ByUser    map[string]int `json:"by_user"`
		TotalRows struct {
			Processed int64 `json:"processed"`
			Succeeded int64 `json:"succeeded"`
			Failed    int64 `json:"failed"`
		} `json:"total_rows"`
		ActiveJobs  int       `json:"active_jobs"`
		AvgDuration string    `json:"avg_duration"`
		SuccessRate string    `json:"success_rate"`
		Timestamp   time.Time `json:"timestamp"`
	}{
		ByStatus: make(map[string]int),
		ByUser:   make(map[string]int),
	}

	var totalDuration time.Duration
	var completedJobs int
	var totalProcessedRows int64
	var totalSucceededRows int64
	var totalFailedRows int64

	for _, status := range jobTracker {
		stats.Total++
		stats.ByStatus[status.Status]++
		stats.ByUser[status.CreatedBy]++

		// Count active jobs
		if status.Status == "STARTED" || status.Status == "PROCESSING" {
			stats.ActiveJobs++
		}

		// Calculate duration for completed jobs
		if status.EndTime != nil {
			totalDuration += status.EndTime.Sub(status.StartTime)
			completedJobs++
		}

		// Aggregate row statistics
		totalProcessedRows += status.Processed
		totalSucceededRows += status.Succeeded
		totalFailedRows += status.Failed
	}

	stats.TotalRows.Processed = totalProcessedRows
	stats.TotalRows.Succeeded = totalSucceededRows
	stats.TotalRows.Failed = totalFailedRows

	// Calculate average duration
	if completedJobs > 0 {
		avgDuration := totalDuration / time.Duration(completedJobs)
		stats.AvgDuration = avgDuration.Round(time.Second).String()
	} else {
		stats.AvgDuration = "N/A"
	}

	// Calculate success rate
	if totalProcessedRows > 0 {
		successRate := float64(totalSucceededRows) / float64(totalProcessedRows) * 100
		stats.SuccessRate = fmt.Sprintf("%.1f%%", successRate)
	} else {
		stats.SuccessRate = "N/A"
	}

	stats.Timestamp = time.Now()

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": stats})
	return nil
}

// CancelJob cancels an active job
func CancelJob(aepr *api.DXAPIEndPointRequest) error {
	_, jobID, err := aepr.GetParameterValueAsString("job_id")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	if jobID == "" {
		aepr.WriteResponseAsJSON(http.StatusBadRequest, nil, utils.JSON{
			"error": "job_id is required",
		})
		return nil
	}

	jobMutex.Lock()
	defer jobMutex.Unlock()

	status, exists := jobTracker[jobID]
	if !exists {
		aepr.WriteResponseAsJSON(http.StatusNotFound, nil, utils.JSON{
			"error":  "Job not found",
			"job_id": jobID,
		})
		return nil
	}

	// Check if job can be cancelled
	if status.Status != "STARTED" && status.Status != "PROCESSING" {
		aepr.WriteResponseAsJSON(http.StatusBadRequest, nil, utils.JSON{
			"data": utils.JSON{
				"job_id":         jobID,
				"current_status": status.Status,
			},
			"reason":         "Job cannot be cancelled",
			"reason_message": "Only STARTED or PROCESSING jobs can be cancelled",
		})
		return nil
	}

	// Update job status to cancelled
	status.Status = "CANCELLED"
	endTime := time.Now()
	status.EndTime = &endTime
	status.Errors = append(status.Errors, "Job was cancelled by user")

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
		"message":      "Job cancelled successfully",
		"job_id":       jobID,
		"status":       "CANCELLED",
		"cancelled_at": endTime,
	})

	return nil
}

// Additional utility function to get jobs by specific criteria
func GetJobsByStatus(status string) []*JobStatus {
	jobMutex.RLock()
	defer jobMutex.RUnlock()

	var jobs []*JobStatus
	for _, job := range jobTracker {
		if job.Status == status {
			jobs = append(jobs, job)
		}
	}

	return jobs
}

// Get jobs for a specific user
func GetJobsByUser(user string) []*JobStatus {
	jobMutex.RLock()
	defer jobMutex.RUnlock()

	var jobs []*JobStatus
	for _, job := range jobTracker {
		if job.CreatedBy == user {
			jobs = append(jobs, job)
		}
	}

	return jobs
}
