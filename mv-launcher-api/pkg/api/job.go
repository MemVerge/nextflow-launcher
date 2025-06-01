package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MemVerge/nf-launcher/pkg/services"
	"github.com/MemVerge/nf-launcher/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	batchtypes "github.com/aws/aws-sdk-go-v2/service/batch/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// JobWithStatus represents a job with its AWS Batch status and timing information
type JobWithStatus struct {
	types.Job
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	StartedAt       time.Time `json:"started_at,omitempty"`
	StoppedAt       time.Time `json:"stopped_at,omitempty"`
	ExitCode        int32     `json:"exit_code,omitempty"`
	StatusReason    string    `json:"status_reason,omitempty"`
	BatchJobId      string    `json:"batch_job_id,omitempty"`
	JobDefinition   string    `json:"job_definition,omitempty"`
	JobQueue        string    `json:"job_queue,omitempty"`
	Attempts        int32     `json:"attempts,omitempty"`
	ContainerReason string    `json:"container_reason,omitempty"`
	LogStreamName   string    `json:"log_stream_name,omitempty"`
	Duration        int64     `json:"duration,omitempty"` // Duration in seconds
	Memory          int32     `json:"memory,omitempty"`   // Memory in MB
	Vcpus           int32     `json:"vcpus,omitempty"`
}

// @Summary Create a new job
// @Description Create a new job by uploading a JSON specification
// @Accept  json
// @Produce json
// @Param   job body types.Job true "Job Specification"
// @Success 201 {object} types.Job
// @Router /jobs [post]
func (a API) CreateJob(c *gin.Context) {
	var pJob types.Job
	if err := c.ShouldBindJSON(&pJob); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Received job request: %+v", pJob)

	// Set default values if not provided
	if pJob.Memory == "" {
		pJob.Memory = "20G"
	}
	if pJob.MaxRetries == 0 {
		pJob.MaxRetries = 5
	}

	// Store job in S3
	log.Printf("Storing job in S3 bucket: %s", a.config.JobBucket)
	// Generate job ID if not provided
	if pJob.ID == "" {
		id := uuid.New()
		pJob.ID = id.String()
	}
	err := services.PutJob(a.s3Client, a.config.JobBucket, pJob)
	if err != nil {
		log.Printf("Error storing job in S3: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Job %s stored successfully in S3", pJob.ID)

	// Submit job to AWS Batch
	jobDefinition := fmt.Sprintf("%s-nextflow-headnode", a.config.Environment)
	log.Printf("Using job definition: %s", jobDefinition)

	// Create job submission
	jobInput := &batch.SubmitJobInput{
		JobName:       aws.String(pJob.ID),
		JobQueue:      aws.String(pJob.HeadNodeQueue),
		JobDefinition: aws.String(jobDefinition),
		ContainerOverrides: &batchtypes.ContainerOverrides{
			Environment: []batchtypes.KeyValuePair{
				{
					Name:  aws.String("JOB_ID"),
					Value: aws.String(pJob.ID),
				},
				{
					Name:  aws.String("PIPELINE"),
					Value: aws.String(pJob.Pipeline),
				},
				{
					Name:  aws.String("WORK_DIR"),
					Value: aws.String(pJob.WorkDir),
				},
				{
					Name:  aws.String("RESULT_DIR"),
					Value: aws.String(pJob.ResultDir),
				},
				{
					Name:  aws.String("LOG_BUCKET"),
					Value: aws.String(pJob.LogBucket),
				},
			},
		},
	}

	// Submit the job
	result, err := a.batchClient.SubmitJob(context.TODO(), jobInput)
	if err != nil {
		log.Printf("Error submitting job to AWS Batch: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Return job details
	c.JSON(200, gin.H{
		"id":     pJob.ID,
		"name":   pJob.Name,
		"status": "SUBMITTED",
		"arn":    result.JobArn,
	})
}

// @Summary List all jobs
// @Description Returns a JSON blob with a list of all jobs
// @Accept  json
// @Produce json
// @Success 200 {object} types.Jobs
// @Router /jobs [get]
func (a API) ListJobs(c *gin.Context) {
	// Get queue from query parameter
	queue := c.Query("queue")
	if queue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Queue parameter is required"})
		return
	}

	// List jobs from the specified queue for each valid status
	validStatuses := []batchtypes.JobStatus{
		batchtypes.JobStatusSubmitted,
		batchtypes.JobStatusPending,
		batchtypes.JobStatusRunnable,
		batchtypes.JobStatusStarting,
		batchtypes.JobStatusRunning,
		batchtypes.JobStatusSucceeded,
		batchtypes.JobStatusFailed,
	}

	// Fetch all job specs from S3
	jobSpecs, err := services.GetJobs(a.s3Client, a.config.JobBucket)
	if err != nil {
		log.Printf("Error fetching job specs from S3: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Build a map from job name to job spec
	jobSpecMap := make(map[string]types.Job)
	for _, spec := range jobSpecs {
		jobSpecMap[spec.Name] = spec
	}

	// Get job IDs for detailed information
	jobIds := make([]string, 0)
	for _, status := range validStatuses {
		listInput := &batch.ListJobsInput{
			JobQueue:  aws.String(queue),
			JobStatus: status,
		}

		log.Printf("Listing jobs with status %s from queue: %s", status, queue)
		listOutput, err := a.batchClient.ListJobs(context.TODO(), listInput)
		if err != nil {
			log.Printf("Error listing jobs with status %s from queue %s: %v", status, queue, err)
			continue
		}

		log.Printf("Found %d jobs with status %s in queue %s", len(listOutput.JobSummaryList), status, queue)
		for _, job := range listOutput.JobSummaryList {
			jobIds = append(jobIds, *job.JobId)
			log.Printf("Found job: %s (Name: %s, Status: %s)", *job.JobId, *job.JobName, job.Status)
		}
	}

	// Describe jobs in batches of 100 (AWS Batch limit)
	jobsWithStatus := make([]JobWithStatus, 0)
	for i := 0; i < len(jobIds); i += 100 {
		end := i + 100
		if end > len(jobIds) {
			end = len(jobIds)
		}

		describeInput := &batch.DescribeJobsInput{
			Jobs: jobIds[i:end],
		}
		describeOutput, err := a.batchClient.DescribeJobs(context.TODO(), describeInput)
		if err != nil {
			log.Printf("Error describing jobs: %v", err)
			continue
		}

		for _, job := range describeOutput.Jobs {
			// Convert timestamps from milliseconds to time.Time
			createdAt := time.Time{}
			if job.CreatedAt != nil {
				createdAt = time.UnixMilli(*job.CreatedAt)
			}
			startedAt := time.Time{}
			if job.StartedAt != nil {
				startedAt = time.UnixMilli(*job.StartedAt)
			}
			stoppedAt := time.Time{}
			if job.StoppedAt != nil {
				stoppedAt = time.UnixMilli(*job.StoppedAt)
			}

			// Get exit code if available
			exitCode := int32(0)
			if job.Container != nil && job.Container.ExitCode != nil {
				exitCode = *job.Container.ExitCode
			}

			// Get status reason if available
			statusReason := ""
			if job.StatusReason != nil {
				statusReason = *job.StatusReason
			}

			// Get job details
			jobDefinition := ""
			if job.JobDefinition != nil {
				jobDefinition = *job.JobDefinition
			}

			jobQueue := ""
			if job.JobQueue != nil {
				jobQueue = *job.JobQueue
			}

			attempts := int32(0)
			if job.Attempts != nil {
				attempts = int32(len(job.Attempts))
			}

			containerReason := ""
			if job.Container != nil && job.Container.Reason != nil {
				containerReason = *job.Container.Reason
			}

			logStreamName := ""
			if job.Container != nil && job.Container.LogStreamName != nil {
				logStreamName = *job.Container.LogStreamName
			}

			// Calculate duration in seconds
			duration := int64(0)
			if !startedAt.IsZero() && !stoppedAt.IsZero() {
				duration = int64(stoppedAt.Sub(startedAt).Seconds())
			}

			// Get resource requirements
			memory := int32(0)
			vcpus := int32(0)
			if job.Container != nil && job.Container.ResourceRequirements != nil {
				for _, req := range job.Container.ResourceRequirements {
					if req.Type == batchtypes.ResourceTypeMemory && req.Value != nil {
						if value, err := strconv.ParseInt(*req.Value, 10, 32); err == nil {
							memory = int32(value)
						}
					}
					if req.Type == batchtypes.ResourceTypeVcpu && req.Value != nil {
						if value, err := strconv.ParseInt(*req.Value, 10, 32); err == nil {
							vcpus = int32(value)
						}
					}
				}
			}

			// Find the job spec by name
			jobSpec, ok := jobSpecMap[*job.JobName]
			jobID := *job.JobId // fallback to Batch job ID if not found
			if ok {
				jobID = jobSpec.ID
			}

			jobWithStatus := JobWithStatus{
				Job: types.Job{
					ID:   jobID,
					Name: *job.JobName,
				},
				Status:          string(job.Status),
				CreatedAt:       createdAt,
				StartedAt:       startedAt,
				StoppedAt:       stoppedAt,
				ExitCode:        exitCode,
				StatusReason:    statusReason,
				BatchJobId:      *job.JobId,
				JobDefinition:   jobDefinition,
				JobQueue:        jobQueue,
				Attempts:        attempts,
				ContainerReason: containerReason,
				LogStreamName:   logStreamName,
				Duration:        duration,
				Memory:          memory,
				Vcpus:           vcpus,
			}

			jobsWithStatus = append(jobsWithStatus, jobWithStatus)
			log.Printf("Job details - Name: %s, Status: %s, Created: %s, Started: %s, Stopped: %s, ExitCode: %d, Duration: %ds, Memory: %dMB, vCPUs: %d",
				*job.JobName, job.Status, createdAt.Format(time.RFC3339), startedAt.Format(time.RFC3339), stoppedAt.Format(time.RFC3339), exitCode, duration, memory, vcpus)
		}
	}

	c.JSON(200, jobsWithStatus)
}

type JobLogs struct {
	NextflowLog string `json:"nextflow_log"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	JobName     string `json:"job_name,omitempty"`
	BatchJobId  string `json:"batch_job_id,omitempty"`
}

// @Summary Get job logs
// @Description Get CloudWatch logs and nextflow.log for a specific job
// @Accept  json
// @Produce json
// @Param   id path string true "Job ID"
// @Success 200 {object} JobLogs
// @Router /jobs/{id}/logs [get]
func (a *API) GetJobLogs(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(400, gin.H{"error": "Job ID is required"})
		return
	}

	log.Printf("Fetching logs for job ID: %s", jobID)

	// Try to get job spec from S3
	jobSpec, err := services.GetJob(a.s3Client, a.config.JobBucket, jobID)
	if err != nil {
		log.Printf("Error getting job spec from S3: %v", err)
		// Continue with job ID as name, as it might be the actual job name
	}

	// Use job name from spec if available, otherwise use the ID
	jobName := jobID
	if jobSpec != nil {
		jobName = jobSpec.Name
		log.Printf("Found job spec in S3, using job name: %s", jobName)
	}

	// Get job details from AWS Batch
	listInput := &batch.ListJobsInput{
		JobQueue: aws.String("launcher-test-on-demand-MM-Batch-JobQueue"),
	}

	// Try all possible job statuses
	validStatuses := []batchtypes.JobStatus{
		batchtypes.JobStatusSubmitted,
		batchtypes.JobStatusPending,
		batchtypes.JobStatusRunnable,
		batchtypes.JobStatusStarting,
		batchtypes.JobStatusRunning,
		batchtypes.JobStatusSucceeded,
		batchtypes.JobStatusFailed,
	}

	var foundJob *batchtypes.JobSummary
	for _, status := range validStatuses {
		listInput.JobStatus = status
		listOutput, err := a.batchClient.ListJobs(context.TODO(), listInput)
		if err != nil {
			log.Printf("Error listing jobs with status %s: %v", status, err)
			continue
		}

		// Search for the job by name
		for _, job := range listOutput.JobSummaryList {
			if *job.JobName == jobName {
				foundJob = &job
				log.Printf("Found job with name %s in status %s", jobName, status)
				break
			}
		}
		if foundJob != nil {
			break
		}
	}

	if foundJob == nil {
		log.Printf("Job not found with name: %s", jobName)
		c.JSON(404, gin.H{"error": "Job not found"})
		return
	}

	// Get job details from AWS Batch
	describeInput := &batch.DescribeJobsInput{
		Jobs: []string{*foundJob.JobId},
	}

	describeOutput, err := a.batchClient.DescribeJobs(context.TODO(), describeInput)
	if err != nil {
		log.Printf("Error describing job: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(describeOutput.Jobs) == 0 {
		log.Printf("No job details found for job ID: %s", *foundJob.JobId)
		c.JSON(404, gin.H{"error": "Job not found"})
		return
	}

	jobDetail := describeOutput.Jobs[0]
	log.Printf("Found job details - Name: %s, Status: %s", *jobDetail.JobName, jobDetail.Status)

	var logs JobLogs
	logs.Status = string(jobDetail.Status)
	logs.JobName = *jobDetail.JobName
	logs.BatchJobId = *jobDetail.JobId

	// Try to get Nextflow log from S3 if job is completed
	if jobDetail.Status == batchtypes.JobStatusSucceeded || jobDetail.Status == batchtypes.JobStatusFailed {
		logKey := fmt.Sprintf("jobs/%s/nextflow.log", jobID)
		getObjectInput := &s3.GetObjectInput{
			Bucket: aws.String(a.config.LogBucket),
			Key:    aws.String(logKey),
		}

		result, err := a.s3Client.GetObject(context.TODO(), getObjectInput)
		if err == nil {
			defer result.Body.Close()
			body, err := io.ReadAll(result.Body)
			if err == nil {
				logs.NextflowLog = string(body)
				log.Printf("Successfully retrieved Nextflow log from S3 for job: %s", *jobDetail.JobName)
			} else {
				log.Printf("Error reading S3 log file: %v", err)
				logs.Message = fmt.Sprintf("Error reading S3 log file: %v", err)
			}
		} else {
			log.Printf("Error getting S3 log file: %v", err)
			logs.Message = fmt.Sprintf("Nextflow log not available in S3 yet: %v", err)
		}
	} else {
		log.Printf("Job %s is not completed yet, Nextflow log will be available after completion", *jobDetail.JobName)
		logs.Message = "Nextflow log will be available after job completion"
	}

	c.JSON(200, logs)
}

// @Summary Get a pre-signed S3 URL for the nextflow.log file
// @Description Returns a pre-signed URL to download nextflow.log for a specific job
// @Accept  json
// @Produce json
// @Param   id path string true "Job ID"
// @Success 200 {object} map[string]string
// @Router /jobs/{id}/log-url [get]
func (a *API) GetJobLogPresignedURL(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(400, gin.H{"error": "Job ID is required"})
		return
	}

	log.Printf("Fetching presigned URL for job ID: %s", jobID)

	// Get job logs from S3
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(a.config.LogBucket),
		Key:    aws.String(fmt.Sprintf("jobs/%s/nextflow.log", jobID)),
	}

	result, err := a.s3Client.GetObject(context.TODO(), getObjectInput)
	if err != nil {
		log.Printf("Error getting job logs from S3: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer result.Body.Close()

	// Get presigned URL for the log file
	presignClient := s3.NewPresignClient(a.s3Client)
	presignedURL, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(a.config.LogBucket),
		Key:    aws.String(fmt.Sprintf("jobs/%s/nextflow.log", jobID)),
	})
	if err != nil {
		log.Printf("Error generating presigned URL: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"url": presignedURL.URL})
}

// GetJob retrieves a job by ID
func (a *API) GetJob(jobID string) (*types.Job, error) {
	return services.GetJob(a.s3Client, a.config.JobBucket, jobID)
}
