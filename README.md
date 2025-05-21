# Nextflow Launcher

A web application for launching and managing Nextflow pipelines on AWS Batch.

## Architecture

The application consists of:
- Frontend: Vue.js application
- Backend: Go API server
- AWS Resources: Batch, S3, IAM roles

## Prerequisites

- AWS Account with appropriate permissions
- AWS CLI configured
- Docker installed
- Node.js and npm installed
- Go 1.21 or later installed

## AWS Resource Setup

### 1. Create Required S3 Buckets

```bash
# Create buckets for different purposes
aws s3 mb s3://your-pipeline-bucket --region your-region
aws s3 mb s3://your-job-bucket --region your-region
aws s3 mb s3://your-log-bucket --region your-region
```

### 2. Create IAM Roles

#### Batch Job Execution Role
```bash
# Create the role
aws iam create-role --role-name batchJobExecutionRole \
    --assume-role-policy-document '{
        "Version": "2012-10-17",
        "Statement": [{
            "Effect": "Allow",
            "Principal": {
                "Service": "ecs-tasks.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }]
    }'

# Attach required policies
aws iam attach-role-policy --role-name batchJobExecutionRole \
    --policy-arn arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy
aws iam attach-role-policy --role-name batchJobExecutionRole \
    --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess
```

### 3. Create AWS Batch Resources

Use the AWS Batch Setup component in the application to create:
- Compute Environment
- Job Queue
- Job Definition

Or use AWS CLI:

```bash
# Create compute environment
aws batch create-compute-environment \
    --compute-environment-name nextflow-compute-env \
    --type MANAGED \
    --compute-resources type=EC2,minvCpus=0,maxvCpus=256,desiredvCpus=0,instanceTypes=m5.large,subnets=subnet-xxxxx,securityGroupIds=sg-xxxxx,instanceRole=ecsInstanceRole

# Create job queue
aws batch create-job-queue \
    --job-queue-name nextflow-job-queue \
    --priority 1 \
    --compute-environment-order order=1,computeEnvironment=nextflow-compute-env
```

## Backend Deployment

### 1. Build the Docker Image

```bash
# Build the backend image
cd mv-launcher-api
docker build -t nextflow-launcher-api:latest .
```

### 2. Deploy to ECS/EKS or EC2

#### Option A: Deploy to ECS
```bash
# Create ECS cluster and service
aws ecs create-cluster --cluster-name nextflow-launcher
aws ecs create-service \
    --cluster nextflow-launcher \
    --service-name api-service \
    --task-definition nextflow-launcher-api \
    --desired-count 1
```

#### Option B: Deploy to EC2
```bash
# Launch EC2 instance and run container
aws ec2 run-instances \
    --image-id ami-xxxxx \
    --instance-type t2.micro \
    --user-data '#!/bin/bash
        docker run -d \
        -p 8080:8080 \
        -e AWS_REGION=your-region \
        -e PIPELINE_BUCKET=your-pipeline-bucket \
        -e JOB_BUCKET=your-job-bucket \
        -e LOG_BUCKET=your-log-bucket \
        -e JOB_ROLE_ARN=arn:aws:iam::your-account:role/batchJobExecutionRole \
        nextflow-launcher-api:latest'
```

## Frontend Deployment

### 1. Build the Frontend

```bash
# Build the Vue.js application
cd vue
npm install
npm run build
```

### 2. Deploy to S3

```bash
# Create S3 bucket for frontend
aws s3 mb s3://your-frontend-bucket --region your-region

# Upload built files
aws s3 sync dist/ s3://your-frontend-bucket

# Enable static website hosting
aws s3 website s3://your-frontend-bucket \
    --index-document index.html \
    --error-document index.html
```

### 3. Configure CloudFront (Optional)

```bash
# Create CloudFront distribution
aws cloudfront create-distribution \
    --origin-domain-name your-frontend-bucket.s3.amazonaws.com \
    --default-root-object index.html
```

## Environment Variables

### Backend Environment Variables
```bash
AWS_REGION=your-region
PIPELINE_BUCKET=your-pipeline-bucket
JOB_BUCKET=your-job-bucket
LOG_BUCKET=your-log-bucket
JOB_ROLE_ARN=arn:aws:iam::your-account:role/batchJobExecutionRole
```

### Frontend Environment Variables
Create a `.env` file in the `vue` directory:
```
VITE_API_URL=http://your-backend-url:8080
```

## Updating Configuration

### Update Bucket Names
1. Update environment variables in your deployment
2. Update default values in `mv-launcher-api/pkg/api/main.go`:
```go
PipelineBucket = getEnvOrDefault("PIPELINE_BUCKET", "your-pipeline-bucket")
JobBucket      = getEnvOrDefault("JOB_BUCKET", "your-job-bucket")
LogBucket      = getEnvOrDefault("LOG_BUCKET", "your-log-bucket")
```

### Update Job Definition
The job definition is created automatically when the first job is submitted. To update it:
1. Go to AWS Batch console
2. Find the job definition "nextflow-headnode-launcher"
3. Create a new revision with updated parameters

## Monitoring and Logs

- Backend logs: CloudWatch Logs
- Frontend logs: Browser console
- Job logs: S3 bucket specified in LOG_BUCKET
- AWS Batch logs: CloudWatch Logs

## Security Considerations

1. Use AWS Secrets Manager for sensitive credentials
2. Implement proper IAM roles and policies
3. Enable HTTPS for all endpoints
4. Regular security updates and patches
5. Monitor AWS CloudTrail for suspicious activities

## Troubleshooting

### Common Issues

1. **Job Submission Fails**
   - Check IAM permissions
   - Verify job definition exists
   - Check compute environment status

2. **Frontend Can't Connect to Backend**
   - Verify CORS settings
   - Check network security groups
   - Verify API endpoint configuration

3. **S3 Access Issues**
   - Verify bucket policies
   - Check IAM role permissions
   - Verify bucket names and regions

## Support

For issues and support, please create an issue in the GitHub repository. 