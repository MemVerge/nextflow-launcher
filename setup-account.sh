#!/bin/bash

# Check if AWS CLI is installed
if ! command -v aws &> /dev/null; then
    echo "AWS CLI is not installed. Please install it first."
    exit 1
fi

# Check if required parameters are provided
if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <environment-name> [nextflow-image]"
    echo "Example: $0 dev achyutha98/nextflow:latest"
    exit 1
fi

ENVIRONMENT_NAME=$1
NEXTFLOW_IMAGE=${2:-"achyutha98/nextflow:latest"}

# Create CloudFormation stack
echo "Creating CloudFormation stack for account setup..."
aws cloudformation create-stack \
    --stack-name ${ENVIRONMENT_NAME}-nextflow-launcher-setup \
    --template-body file://setup-account.yaml \
    --parameters \
        ParameterKey=EnvironmentName,ParameterValue=${ENVIRONMENT_NAME} \
        ParameterKey=NextflowImage,ParameterValue=${NEXTFLOW_IMAGE} \
    --capabilities CAPABILITY_IAM

# Wait for stack creation to complete
echo "Waiting for stack creation to complete..."
aws cloudformation wait stack-create-complete \
    --stack-name ${ENVIRONMENT_NAME}-nextflow-launcher-setup

# Get stack outputs
echo "Getting stack outputs..."
aws cloudformation describe-stacks \
    --stack-name ${ENVIRONMENT_NAME}-nextflow-launcher-setup \
    --query 'Stacks[0].Outputs' \
    --output table

echo "Account setup completed successfully!" 