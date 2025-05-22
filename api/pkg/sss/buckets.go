package sss

import (
	"context"

	"github.com/ChristianKniep/mv-launcher-api/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sirupsen/logrus"
)

func ListBuckets(cfg aws.Config) (bs types.Buckets, err error) {
	svc := s3.NewFromConfig(cfg)
	logrus.Info("Listing all buckets")
	result, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		logrus.Errorf("Error listing buckets: %v", err)
		return bs, err
	}
	logrus.Infof("Found %d buckets", len(result.Buckets))

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
