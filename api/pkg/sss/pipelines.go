package sss

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ChristianKniep/mv-launcher-api/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetPipelines(cfg aws.Config, bucket string) (pipelines types.Pipelines, err error) {
	svc := s3.NewFromConfig(cfg)
	result, err := svc.ListObjects(context.TODO(), &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return pipelines, err
	}

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
		var pipeline types.Pipeline
		if err := json.NewDecoder(objResult.Body).Decode(&pipeline); err != nil {
			log.Printf("failed to decode object %s, %v", *item.Key, err)
			continue
		}
		pipelines = append(pipelines, pipeline)

	}
	return pipelines, nil
}
