package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ChristianKniep/mv-launcher-api/pkg/sss"
	"github.com/ChristianKniep/mv-launcher-api/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	batchtypes "github.com/aws/aws-sdk-go-v2/service/batch/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
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

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AwsRegion))
	if err != nil {
		log.Printf("Error loading AWS config: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Store job in S3
	log.Printf("Storing job in S3 bucket: %s", JobBucket)
	// Generate job ID if not provided
	if pJob.ID == "" {
		id := uuid.New()
		pJob.ID = id.String()
	}
	err = sss.PutJob(cfg, JobBucket, pJob)
	if err != nil {
		log.Printf("Error storing job in S3: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Job %s stored successfully in S3", pJob.ID)

	// Check if job definition exists, create if it doesn't
	jobDefName := "nextflow-headnode-launcher"
	log.Printf("Using job definition: %s", jobDefName)

	describeInput := &batch.DescribeJobDefinitionsInput{
		JobDefinitionName: aws.String(jobDefName),
		Status:            aws.String("ACTIVE"),
	}

	describeOutput, err := a.BatchClient.DescribeJobDefinitions(context.TODO(), describeInput)
	if err != nil {
		log.Printf("Error describing job definitions: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(describeOutput.JobDefinitions) == 0 {
		log.Printf("Job definition not found, creating new one")
		registerInput := &batch.RegisterJobDefinitionInput{
			JobDefinitionName: aws.String(jobDefName),
			Type:              batchtypes.JobDefinitionTypeContainer,
			ContainerProperties: &batchtypes.ContainerProperties{
				Image:  aws.String("achyutha98/nextflow:latest"),
				Vcpus:  aws.Int32(4),
				Memory: aws.Int32(16384),
				Command: []string{
					"nextflow",
					"run",
					"Ref::pipeline",
					"-work-dir",
					"Ref::work_dir",
					"-resume",
					"-profile",
					"Ref::profile",
					"--outdir",
					"Ref::outdir",
					"--input",
					"Ref::input_dir",
				},
				JobRoleArn: aws.String(JobRoleArn),
				Environment: []batchtypes.KeyValuePair{
					{
						Name:  aws.String("NXF_ANSI_LOG"),
						Value: aws.String("false"),
					},
					{
						Name:  aws.String("NXF_OPTS"),
						Value: aws.String("-Xms1g -Xmx12g"),
					},
				},
			},
		}

		_, err = a.BatchClient.RegisterJobDefinition(context.TODO(), registerInput)
		if err != nil {
			log.Printf("Error creating job definition: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Job definition created successfully")
	} else {
		log.Printf("Using existing job definition: %s", *describeOutput.JobDefinitions[0].JobDefinitionArn)
		jobDefName = *describeOutput.JobDefinitions[0].JobDefinitionArn
	}

	// Submit job to AWS Batch
	log.Printf("Submitting job to queue: %s", pJob.HeadNodeQueue)
	submitInput := &batch.SubmitJobInput{
		JobName:       &pJob.Name,
		JobQueue:      &pJob.HeadNodeQueue,
		JobDefinition: aws.String("nextflow-headnode-launcher"),
		ContainerOverrides: &batchtypes.ContainerOverrides{
			Environment: []batchtypes.KeyValuePair{
				{
					Name:  aws.String("JOB_SPEC_PATH"),
					Value: aws.String(fmt.Sprintf("s3://%s/job-%s.json", JobBucket, pJob.ID)),
				},
				{
					Name:  aws.String("AWS_ACCESS_KEY_ID"),
					Value: aws.String(pJob.AWSAccessKey),
				},
				{
					Name:  aws.String("AWS_SECRET_ACCESS_KEY"),
					Value: aws.String(pJob.AWSSecretKey),
				},
				{
					Name:  aws.String("AWS_REGION"),
					Value: aws.String(AwsRegion),
				},
				{
					Name:  aws.String("JOB_ID"),
					Value: aws.String(pJob.ID),
				},
				{
					Name:  aws.String("NEXTFLOW_LOG_PATH"),
					Value: aws.String("/var/log/nextflow/nextflow.log"),
				},
				{
					Name:  aws.String("LOG_BUCKET"),
					Value: aws.String(pJob.LogBucket),
				},
				{
					Name:  aws.String("NXF_WORK"),
					Value: aws.String("/workspace/work"),
				},
			},
		},
	}

	submitOutput, err := a.BatchClient.SubmitJob(context.TODO(), submitInput)
	if err != nil {
		log.Printf("Error submitting job to AWS Batch: %v", err)
		c.JSON(500, gin.H{"error": "Failed to submit job: " + err.Error()})
		return
	}
	log.Printf("Job submitted successfully with ID: %s", *submitOutput.JobId)

	c.JSON(201, gin.H{
		"message": "Job submitted successfully",
		"jobId":   submitOutput.JobId,
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
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AwsRegion))
	if err != nil {
		log.Printf("Error loading AWS config: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	jobSpecs, err := sss.GetJobs(cfg, JobBucket)
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
		listOutput, err := a.BatchClient.ListJobs(context.TODO(), listInput)
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
		describeOutput, err := a.BatchClient.DescribeJobs(context.TODO(), describeInput)
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

// @Summary Get job logs
// @Description Get CloudWatch logs and nextflow.log for a specific job
// @Accept  json
// @Produce json
// @Param   id path string true "Job ID"
// @Success 200 {object} JobLogs
// @Router /jobs/{id}/logs [get]
func (a API) GetJobLogs(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(400, gin.H{"error": "Job ID is required"})
		return
	}

	// Get job details directly using DescribeJobs
	describeInput := &batch.DescribeJobsInput{
		Jobs: []string{jobID},
	}

	describeOutput, err := a.BatchClient.DescribeJobs(context.TODO(), describeInput)
	if err != nil {
		log.Printf("Error describing job: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(describeOutput.Jobs) == 0 {
		// Try to find the job in the list of all jobs
		listInput := &batch.ListJobsInput{
			JobQueue: aws.String("launcher-test-on-demand-MM-Batch-JobQueue"),
		}

		listOutput, err := a.BatchClient.ListJobs(context.TODO(), listInput)
		if err != nil {
			log.Printf("Error listing jobs: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Search for the job by name
		var foundJob *batchtypes.JobSummary
		for _, job := range listOutput.JobSummaryList {
			if *job.JobName == jobID {
				foundJob = &job
				break
			}
		}

		if foundJob == nil {
			c.JSON(404, gin.H{"error": "Job not found"})
			return
		}

		// Get job details using the found job ID
		describeInput = &batch.DescribeJobsInput{
			Jobs: []string{*foundJob.JobId},
		}

		describeOutput, err = a.BatchClient.DescribeJobs(context.TODO(), describeInput)
		if err != nil {
			log.Printf("Error describing job: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if len(describeOutput.Jobs) == 0 {
			c.JSON(404, gin.H{"error": "Job not found"})
			return
		}
	}

	jobDetail := describeOutput.Jobs[0]

	type LogEntry struct {
		Timestamp time.Time `json:"timestamp"`
		Message   string    `json:"message"`
	}

	type JobLogs struct {
		CloudWatchLogs []LogEntry `json:"cloudwatch_logs"`
		NextflowLog    string     `json:"nextflow_log"`
	}

	var logs JobLogs

	// Get container logs from CloudWatch
	if jobDetail.Container != nil && jobDetail.Container.LogStreamName != nil {
		logStreamName := *jobDetail.Container.LogStreamName
		logGroupName := "/aws/batch/job"

		// Create CloudWatch Logs client
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AwsRegion))
		if err != nil {
			log.Printf("Error loading AWS config: %v", err)
		} else {
			cloudwatchLogs := cloudwatchlogs.NewFromConfig(cfg)

			// Get log events
			getLogEventsInput := &cloudwatchlogs.GetLogEventsInput{
				LogGroupName:  aws.String(logGroupName),
				LogStreamName: aws.String(logStreamName),
				StartFromHead: aws.Bool(true),
			}

			for {
				output, err := cloudwatchLogs.GetLogEvents(context.TODO(), getLogEventsInput)
				if err != nil {
					log.Printf("Error getting log events: %v", err)
					break
				}

				for _, event := range output.Events {
					logs.CloudWatchLogs = append(logs.CloudWatchLogs, LogEntry{
						Timestamp: time.UnixMilli(*event.Timestamp),
						Message:   *event.Message,
					})
				}

				// If there are no more events or we've reached the end, break
				if output.NextForwardToken == nil {
					break
				}
				if getLogEventsInput.NextToken != nil && *output.NextForwardToken == *getLogEventsInput.NextToken {
					break
				}

				// Update the token for the next iteration
				getLogEventsInput.NextToken = output.NextForwardToken
			}
		}
	}

	// Try to get Nextflow log from S3
	logKey := fmt.Sprintf("jobs/%s/nextflow.log", jobID)
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(LogBucket),
		Key:    aws.String(logKey),
	}

	result, err := a.S3Client.GetObject(context.TODO(), getObjectInput)
	if err == nil {
		defer result.Body.Close()
		body, err := io.ReadAll(result.Body)
		if err == nil {
			logs.NextflowLog = string(body)
		}
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

	logBucket := LogBucket
	logKey := fmt.Sprintf("jobs/%s/nextflow.log", jobID)
	presignClient := s3.NewPresignClient(a.S3Client)
	presignInput := &s3.GetObjectInput{
		Bucket: aws.String(logBucket),
		Key:    aws.String(logKey),
	}
	presignResult, err := presignClient.PresignGetObject(context.TODO(), presignInput, s3.WithPresignExpires(5*time.Minute))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate presigned URL"})
		return
	}
	c.JSON(200, gin.H{"url": presignResult.URL})
}
