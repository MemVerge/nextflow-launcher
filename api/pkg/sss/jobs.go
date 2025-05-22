package sss

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/ChristianKniep/mv-launcher-api/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetJobs(cfg aws.Config, bucket string) (jobs types.Jobs, err error) {
	svc := s3.NewFromConfig(cfg)
	logrus.Infof("Check bucket: %s", bucket)

	result, err := svc.ListObjects(context.TODO(), &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return jobs, err
	}
	logrus.Infof("Found %d jobs", len(result.Contents))
	for _, item := range result.Contents {
		objInput := &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(*item.Key),
		}

		objResult, err := svc.GetObject(context.Background(), objInput)
		if err != nil {
			log.Fatalf("failed to get object, %v", err)
		}

		defer objResult.Body.Close()
		var job types.Job
		if err := json.NewDecoder(objResult.Body).Decode(&job); err != nil {
			log.Printf("failed to decode object %s, %v", *item.Key, err)
			continue
		}
		jobs = append(jobs, job)

	}
	return jobs, nil
}

func PutJob(cfg aws.Config, bucket string, job types.Job) error {
	svc := s3.NewFromConfig(cfg)
	job.Verify()
	jobData, err := json.Marshal(job)
	if err != nil {
		return err
	}
	if job.ID == "" {
		id := uuid.New()
		job.ID = id.String()
	}
	jobKey := fmt.Sprintf("job-%s.json", job.ID)
	if pout, err := svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(jobKey),
		Body:        bytes.NewReader(jobData),
		ContentType: aws.String("application/json"),
	}); err != nil {
		return err
	} else {
		log.Printf("Put job %s to %s", job.ID, *pout.ETag)
	}

	return nil
}
