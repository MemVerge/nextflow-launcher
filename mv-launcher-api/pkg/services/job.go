package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/MemVerge/nf-launcher/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetJobs retrieves all jobs from S3
func GetJobs(s3Client *s3.Client, bucket string) (types.Jobs, error) {
	result, err := s3Client.ListObjects(context.TODO(), &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("jobs/"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	jobs := make(types.Jobs, 0)
	for _, item := range result.Contents {
		if *item.Key == "jobs/" {
			continue
		}

		objInput := &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(*item.Key),
		}

		objResult, err := s3Client.GetObject(context.Background(), objInput)
		if err != nil {
			continue
		}

		defer objResult.Body.Close()
		var job types.Job
		if err := json.NewDecoder(objResult.Body).Decode(&job); err != nil {
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// GetJob retrieves a job from S3
func GetJob(s3Client *s3.Client, bucket string, jobID string) (*types.Job, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("jobs/%s/job.json", jobID)),
	}

	result, err := s3Client.GetObject(context.TODO(), getObjectInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get job from S3: %v", err)
	}
	defer result.Body.Close()

	var job types.Job
	if err := json.NewDecoder(result.Body).Decode(&job); err != nil {
		return nil, fmt.Errorf("failed to decode job: %v", err)
	}

	return &job, nil
}

// PutJob stores a job in S3
func PutJob(s3Client *s3.Client, bucket string, job types.Job) error {
	// Convert job to JSON
	jobJSON, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %v", err)
	}

	// Upload to S3
	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("jobs/%s/job.json", job.ID)),
		Body:   bytes.NewReader(jobJSON),
	}

	_, err = s3Client.PutObject(context.TODO(), putObjectInput)
	if err != nil {
		return fmt.Errorf("failed to put job in S3: %v", err)
	}

	return nil
}
