package services

import (
	"context"
	"fmt"

	"github.com/MemVerge/nf-launcher/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sirupsen/logrus"
)

func ListBuckets(s3Client *s3.Client) ([]string, error) {
	logrus.Info("Listing S3 buckets")
	result, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %v", err)
	}
	buckets := make([]string, 0, len(result.Buckets))
	for _, bucket := range result.Buckets {
		buckets = append(buckets, *bucket.Name)
	}
	logrus.Infof("Found %d buckets", len(buckets))
	return buckets, nil
}

func ListBucketsDetailed(cfg aws.Config) (types.Buckets, error) {
	svc := s3.NewFromConfig(cfg)
	logrus.Info("Listing all buckets")
	result, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		logrus.Errorf("Error listing buckets: %v", err)
		return nil, err
	}
	logrus.Infof("Found %d buckets", len(result.Buckets))

	bs := make(types.Buckets, 0, len(result.Buckets))
	for _, bucket := range result.Buckets {
		logrus.Infof("Processing bucket: %s", *bucket.Name)
		bucketArn := "arn:aws:s3:::" + *bucket.Name
		tagsResult, err := svc.GetBucketTagging(context.TODO(), &s3.GetBucketTaggingInput{
			Bucket: bucket.Name,
		})
		if err != nil {
			logrus.Warnf("No tags found for bucket %s: %v", *bucket.Name, err)
			tagsResult = &s3.GetBucketTaggingOutput{TagSet: []s3types.Tag{}}
		}

		tags := make(map[string]string)
		for _, tag := range tagsResult.TagSet {
			tags[*tag.Key] = *tag.Value
		}

		bs = append(bs, types.Bucket{
			Name: *bucket.Name,
			ARN:  bucketArn,
			Tags: tags,
		})
	}
	return bs, nil
}
