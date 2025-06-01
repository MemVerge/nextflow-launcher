package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MemVerge/nf-launcher/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

func GetPipelines(s3Client *s3.Client, bucket string) (pipelines types.Pipelines, err error) {
	logrus.Infof("Checking pipeline bucket: %s", bucket)
	result, err := s3Client.ListObjects(context.TODO(), &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pipeline objects: %v", err)
	}
	logrus.Infof("Found %d pipelines", len(result.Contents))
	for _, item := range result.Contents {
		objInput := &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(*item.Key),
		}
		objResult, err := s3Client.GetObject(context.Background(), objInput)
		if err != nil {
			logrus.Warnf("Failed to get pipeline object %s: %v", *item.Key, err)
			continue
		}
		defer objResult.Body.Close()
		var pipeline types.Pipeline
		if err := json.NewDecoder(objResult.Body).Decode(&pipeline); err != nil {
			logrus.Warnf("Failed to decode pipeline object %s: %v", *item.Key, err)
			continue
		}
		pipelines = append(pipelines, pipeline)
	}
	return pipelines, nil
}
