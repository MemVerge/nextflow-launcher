#!/bin/bash

# Check if AWS CLI is installed
if ! command -v aws &> /dev/null; then
    echo "Error: AWS CLI is not installed"
    exit 1
fi

# Check if environment name is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <environment-name>"
    echo "Example: $0 dev"
    exit 1
fi

ENVIRONMENT_NAME=$1
NEXTFLOW_IMAGE="achyutha98/nextflow:latest"

echo "Creating CloudFormation stack for account setup..."

# Create the CloudFormation stack
aws cloudformation create-stack \
    --stack-name ${ENVIRONMENT_NAME}-nextflow-launcher-setup \
    --template-body file://../setup-account.yaml \
    --parameters \
        ParameterKey=EnvironmentName,ParameterValue=${ENVIRONMENT_NAME} \
        ParameterKey=NextflowImage,ParameterValue=${NEXTFLOW_IMAGE} \
    --capabilities CAPABILITY_NAMED_IAM

echo "Waiting for stack creation to complete..."
aws cloudformation wait stack-create-complete --stack-name ${ENVIRONMENT_NAME}-nextflow-launcher-setup

echo "Getting stack outputs and saving to .env.local..."

# Get stack outputs and save to .env.local
aws cloudformation describe-stacks \
    --stack-name ${ENVIRONMENT_NAME}-nextflow-launcher-setup \
    --query 'Stacks[0].Outputs[*].[OutputKey,OutputValue]' \
    --output text | while read -r key value; do
        case $key in
            PipelineBucket)
                echo "PIPELINE_BUCKET=$value" >> .env.local
                ;;
            JobBucket)
                echo "JOB_BUCKET=$value" >> .env.local
                ;;
            LogBucket)
                echo "LOG_BUCKET=$value" >> .env.local
                ;;
            NextflowHeadnodeJobDefinition)
                echo "NEXTFLOW_HEADNODE_JOB_DEFINITION=$value" >> .env.local
                ;;
        esac
    done

# Add AWS configuration
echo "AWS_REGION=$(aws configure get region)" >> .env.local
echo "AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id)" >> .env.local
echo "AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key)" >> .env.local

# Add server configuration
echo "PORT=8080" >> .env.local
echo "CORS_ALLOWED_ORIGINS=http://localhost:5173" >> .env.local

echo "Account setup completed successfully! Created .env.local with AWS resource values." 