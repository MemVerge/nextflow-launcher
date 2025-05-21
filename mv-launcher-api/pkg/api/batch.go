package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/ChristianKniep/mv-launcher-api/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	batchtypes "github.com/aws/aws-sdk-go-v2/service/batch/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/gin-gonic/gin"
)

// BatchConfig represents the configuration for AWS Batch setup
type BatchConfig struct {
	Region           string   `json:"region" binding:"required"`
	ComputeEnvName   string   `json:"compute_env_name"`
	JobQueueName     string   `json:"job_queue_name"`
	InstanceTypes    []string `json:"instance_types" binding:"required"`
	MinvCpus         int      `json:"min_vcpus" binding:"required"`
	MaxvCpus         int      `json:"max_vcpus" binding:"required"`
	DesiredvCpus     int      `json:"desired_vcpus" binding:"required"`
	SubnetId         string   `json:"subnet_id" binding:"required"`
	SecurityGroupId  string   `json:"security_group_id" binding:"required"`
	UseSpot          bool     `json:"use_spot"`
	EnableMultiQueue bool     `json:"enable_multi_queue"`
	UniquePrefix     string   `json:"unique_prefix" binding:"required"`
}

// @Summary List all queues
// @Description Returns a JSON blob with a list of all queues
// @Accept  json
// @Produce json
// @Success 200 {object} BatchQueues
// @Router /batch/queues [get]
func (a API) ListQueues(c *gin.Context) {
	out, err := a.BatchClient.DescribeJobQueues(c.Request.Context(), &batch.DescribeJobQueuesInput{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	qs := types.Queues{}
	for _, q := range out.JobQueues {
		qs = append(qs, types.Queue{
			Name: *q.JobQueueName,
			ARN:  *q.JobQueueArn,
		})
	}
	c.JSON(200, gin.H{"queues": qs})
}

// createBatchServiceRole creates the AWS Batch service role
func createBatchServiceRole(ctx context.Context, iamClient *iam.Client, roleName string) (*iam.CreateRoleOutput, error) {
	// Create the role with the AWS Batch service principal
	role, err := iamClient.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"batch.amazonaws.com"},"Action":"sts:AssumeRole"}]}`),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create service role: %v", err)
	}

	// Attach the AWS Batch service role policy
	_, err = iamClient.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AWSBatchServiceRole"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to attach service role policy: %v", err)
	}

	return role, nil
}

// createBatchInstanceRole creates the AWS Batch instance role
func createBatchInstanceRole(ctx context.Context, iamClient *iam.Client, roleName string) (*iam.CreateRoleOutput, error) {
	// Create the role with the EC2 service principal
	role, err := iamClient.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"ec2.amazonaws.com"},"Action":"sts:AssumeRole"}]}`),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create instance role: %v", err)
	}

	// Attach the EC2 container service role policy
	_, err = iamClient.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to attach instance role policy: %v", err)
	}

	return role, nil
}

// createBatchInstanceProfile creates the AWS Batch instance profile
func createBatchInstanceProfile(ctx context.Context, iamClient *iam.Client, profileName, roleName string) (*iam.CreateInstanceProfileOutput, error) {
	// Create the instance profile
	profile, err := iamClient.CreateInstanceProfile(ctx, &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create instance profile: %v", err)
	}

	// Add the role to the instance profile
	_, err = iamClient.AddRoleToInstanceProfile(ctx, &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
		RoleName:            aws.String(roleName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add role to instance profile: %v", err)
	}

	return profile, nil
}

// waitForComputeEnvironment waits for the compute environment to be ready
func waitForComputeEnvironment(ctx context.Context, batchClient *batch.Client, computeEnvName string) error {
	for {
		out, err := batchClient.DescribeComputeEnvironments(ctx, &batch.DescribeComputeEnvironmentsInput{
			ComputeEnvironments: []string{computeEnvName},
		})
		if err != nil {
			return fmt.Errorf("failed to describe compute environment: %v", err)
		}

		if len(out.ComputeEnvironments) == 0 {
			return fmt.Errorf("compute environment %s not found", computeEnvName)
		}

		status := out.ComputeEnvironments[0].Status
		if status == batchtypes.CEStatusValid {
			return nil
		}
		if status == batchtypes.CEStatusInvalid {
			return fmt.Errorf("compute environment %s is invalid", computeEnvName)
		}

		// Wait for 5 seconds before checking again
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			continue
		}
	}
}

// SetupAWSBatch handles the AWS Batch setup request
func (api *API) SetupAWSBatch(c *gin.Context) {
	var config BatchConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Received config: %+v", config) // Debug log

	// Validate unique prefix
	if config.UniquePrefix == "" {
		log.Printf("Unique prefix is empty in config: %+v", config) // Debug log
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unique prefix is required"})
		return
	}

	// Validate unique prefix format
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(config.UniquePrefix) {
		log.Printf("Invalid unique prefix format: %s", config.UniquePrefix) // Debug log
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unique prefix must contain only lowercase letters, numbers, and hyphens"})
		return
	}

	ctx := c.Request.Context()

	// Load AWS config
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(config.Region))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to load AWS config: %v", err)})
		return
	}

	// 1. Create IAM roles
	iamClient := iam.NewFromConfig(cfg)

	// Create Batch Service Role
	batchServiceRoleName := fmt.Sprintf("%s-batch-service-role", config.UniquePrefix)
	batchServiceRole, err := createBatchServiceRole(ctx, iamClient, batchServiceRoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create Batch Service Role: %v", err)})
		return
	}

	// Create Batch Instance Role
	batchInstanceRoleName := fmt.Sprintf("%s-batch-instance-role", config.UniquePrefix)
	batchInstanceRole, err := createBatchInstanceRole(ctx, iamClient, batchInstanceRoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create Batch Instance Role: %v", err)})
		return
	}

	// Create Batch Instance Profile
	batchInstanceProfileName := fmt.Sprintf("%s-batch-instance-profile", config.UniquePrefix)
	batchInstanceProfile, err := createBatchInstanceProfile(ctx, iamClient, batchInstanceProfileName, batchInstanceRoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create Batch Instance Profile: %v", err)})
		return
	}

	// 2. Create Compute Environment
	computeEnvName := config.ComputeEnvName
	if computeEnvName == "" {
		computeEnvName = fmt.Sprintf("mm-batch-compute-env-%s", config.UniquePrefix)
	}

	computeEnv, err := api.BatchClient.CreateComputeEnvironment(ctx, &batch.CreateComputeEnvironmentInput{
		ComputeEnvironmentName: aws.String(computeEnvName),
		Type:                   batchtypes.CETypeManaged,
		ComputeResources: &batchtypes.ComputeResource{
			Type:               batchtypes.CRTypeEc2,
			MinvCpus:           aws.Int32(int32(config.MinvCpus)),
			MaxvCpus:           aws.Int32(int32(config.MaxvCpus)),
			DesiredvCpus:       aws.Int32(int32(config.DesiredvCpus)),
			InstanceTypes:      config.InstanceTypes,
			Subnets:            []string{config.SubnetId},
			SecurityGroupIds:   []string{config.SecurityGroupId},
			InstanceRole:       aws.String(*batchInstanceProfile.InstanceProfile.Arn),
			AllocationStrategy: batchtypes.CRAllocationStrategyBestFit,
		},
		ServiceRole: aws.String(*batchServiceRole.Role.Arn),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create Compute Environment: %v", err)})
		return
	}

	// Wait for compute environment to be ready
	if err := waitForComputeEnvironment(ctx, api.BatchClient, computeEnvName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to wait for compute environment: %v", err)})
		return
	}

	// 3. Create Job Queue
	jobQueueName := config.JobQueueName
	if jobQueueName == "" {
		jobQueueName = fmt.Sprintf("mm-batch-job-queue-%s", config.UniquePrefix)
	}

	jobQueue, err := api.BatchClient.CreateJobQueue(ctx, &batch.CreateJobQueueInput{
		JobQueueName: aws.String(jobQueueName),
		State:        batchtypes.JQStateEnabled,
		Priority:     aws.Int32(1),
		ComputeEnvironmentOrder: []batchtypes.ComputeEnvironmentOrder{
			{
				Order:              aws.Int32(1),
				ComputeEnvironment: computeEnv.ComputeEnvironmentArn,
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create Job Queue: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "AWS Batch setup completed successfully",
		"resources": gin.H{
			"compute_environment": gin.H{
				"name": computeEnvName,
				"arn":  *computeEnv.ComputeEnvironmentArn,
			},
			"job_queue": gin.H{
				"name": jobQueueName,
				"arn":  *jobQueue.JobQueueArn,
			},
			"iam_roles": gin.H{
				"service_role": gin.H{
					"name": batchServiceRoleName,
					"arn":  *batchServiceRole.Role.Arn,
				},
				"instance_role": gin.H{
					"name": batchInstanceRoleName,
					"arn":  *batchInstanceRole.Role.Arn,
				},
				"instance_profile": gin.H{
					"name": batchInstanceProfileName,
					"arn":  *batchInstanceProfile.InstanceProfile.Arn,
				},
			},
		},
	})
}
