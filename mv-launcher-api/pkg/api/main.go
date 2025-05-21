package api

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	PipelineBucket = getEnvOrDefault("PIPELINE_BUCKET", "nextflow-pipelines-memverge-launcher")
	JobBucket      = getEnvOrDefault("JOB_BUCKET", "usw2-nextflow-pipeline-workflow")
	LogBucket      = getEnvOrDefault("LOG_BUCKET", "aharish-memverge-nextflow")
	AwsRegion      = getEnvOrDefault("AWS_REGION", "us-west-2")
	JobRoleArn     = getEnvOrDefault("JOB_ROLE_ARN", "")
)

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

type API struct {
	BatchClient *batch.Client
	S3Client    *s3.Client
}

func NewAPI() *API {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(AwsRegion))
	if err != nil {
		log.Printf("Warning: unable to load SDK config, %v", err)
		// Return a minimal API instance without AWS clients
		return &API{}
	}

	// Log AWS configuration details
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		log.Printf("Warning: Unable to retrieve AWS credentials: %v", err)
	} else {
		log.Printf("Using AWS credentials for account: %s in region: %s", creds.AccessKeyID, AwsRegion)
	}

	batchClient := batch.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)
	return &API{
		BatchClient: batchClient,
		S3Client:    s3Client,
	}
}
