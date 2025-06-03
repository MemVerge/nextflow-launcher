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
    --template-body file://setup-account.yaml \
    --parameters \
        ParameterKey=EnvironmentName,ParameterValue=${ENVIRONMENT_NAME} \
        ParameterKey=NextflowImage,ParameterValue=${NEXTFLOW_IMAGE} \
    --capabilities CAPABILITY_NAMED_IAM

echo "Waiting for stack creation to complete..."
aws cloudformation wait stack-create-complete --stack-name ${ENVIRONMENT_NAME}-nextflow-launcher-setup

echo "Getting stack outputs and saving to .env.local..."

# Create a fresh .env.local file
cat > .env.local << EOF
# AWS Configuration
AWS_REGION=$(aws configure get region)
AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id)
AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key)

# Server Configuration
PORT=8080
CORS_ALLOWED_ORIGINS=http://localhost:5173

EOF

# Get stack outputs and append to .env.local
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

echo "Account setup completed successfully! Created .env.local with AWS resource values." 