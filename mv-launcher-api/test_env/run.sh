#!/bin/bash
set -e

# Get the job ID from environment variable
if [ -z "$JOB_ID" ]; then
    echo "Error: JOB_ID environment variable not set"
    exit 1
fi

# Set the job spec path
JOB_SPEC_PATH="s3://usw2-nextflow-pipeline-workflow/job-${JOB_ID}.json"
echo "Downloading job spec from S3: $JOB_SPEC_PATH"

# Check if the file exists in S3 before trying to download
if ! aws s3 ls "$JOB_SPEC_PATH" > /dev/null 2>&1; then
    echo "Error: Job spec file not found in S3: $JOB_SPEC_PATH"
    exit 1
fi

# Download the job spec to a temporary file
TEMP_FILE=$(mktemp)
aws s3 cp "$JOB_SPEC_PATH" "$TEMP_FILE"
if [ $? -ne 0 ]; then
    echo "Error: Failed to download job spec from S3"
    exit 1
fi

# Read the job spec from the temporary file
job=$(cat "$TEMP_FILE")
rm "$TEMP_FILE"

# Parse fields from JSON with error handling
name=$(echo "$job" | jq -r .name)
pipeline=$(echo "$job" | jq -r .pipeline)
profile=$(echo "$job" | jq -r .profile)
head_node_queue=$(echo "$job" | jq -r .head_node_queue)
task_queue=$(echo "$job" | jq -r .task_queue)
work_dir=$(echo "$job" | jq -r .work_dir)
result_dir=$(echo "$job" | jq -r .result_dir)
memory=$(echo "$job" | jq -r .memory)
max_retries=$(echo "$job" | jq -r .max_retries)

# Set defaults if not provided
memory=${memory:-"20G"}
max_retries=${max_retries:-5}

# Validate required fields
if [ -z "$pipeline" ] || [ -z "$profile" ] || [ -z "$work_dir" ] || [ -z "$result_dir" ] || [ -z "$task_queue" ]; then
    echo "Error: Missing required fields in job spec"
    echo "Pipeline: $pipeline"
    echo "Profile: $profile"
    echo "Work Dir: $work_dir"
    echo "Result Dir: $result_dir"
    echo "Task Queue: $task_queue"
    exit 1
fi

# Pull credentials from env (if present)
access_key="${AWS_ACCESS_KEY_ID:-}"
secret_key="${AWS_SECRET_ACCESS_KEY:-}"

echo "Generating aws.config with dynamic credentials..."

cat > aws.config <<EOF
plugins {
    id 'nf-amazon'
}
process {
    executor = 'awsbatch'
    queue = '${task_queue}'
    maxRetries = ${max_retries}
    memory = '${memory}'
}
process.containerOptions = '--env MMC_CHECKPOINT_DIAGNOSIS=true --env MMC_CHECKPOINT_IMAGE_SUBPATH=nextflow --env MMC_CHECKPOINT_INTERVAL=5m --env MMC_CHECKPOINT_MODE=true --env MMC_CHECKPOINT_IMAGE_PATH=/mmc-checkpoint'

aws {
    ${access_key:+accessKey = '${access_key}'}
    ${secret_key:+secretKey = '${secret_key}'}
    region = 'us-west-2'
    client {
        maxConnections = 20
        connectionTimeout = 10000
        uploadStorageClass = 'INTELLIGENT_TIERING'
        storageEncryption = 'AES256'
    }
    batch {
        cliPath = '/nextflow_awscli/bin/aws'
        maxTransferAttempts = 3
        delayBetweenAttempts = '5 sec'
    }
}
EOF

# Create log directory if it doesn't exist
mkdir -p /var/log/nextflow

# Create work directory if it doesn't exist
mkdir -p /workspace/work

echo "Running Nextflow pipeline: $pipeline"
nextflow -log "$NEXTFLOW_LOG_PATH" run "$pipeline" \
    -profile "$profile" \
    -work-dir "$work_dir" \
    --outdir "$result_dir" \
    -c aws.config \
    -ansi-log false

# Upload Nextflow logs to S3
if [ -f "$NEXTFLOW_LOG_PATH" ]; then
    echo "Uploading Nextflow logs to S3..."
    # Create a unique S3 path using the job ID
    s3_log_path="s3://${LOG_BUCKET}/jobs/${JOB_ID}/nextflow.log"
    aws s3 cp "$NEXTFLOW_LOG_PATH" "$s3_log_path"
    if [ $? -ne 0 ]; then
        echo "Warning: Failed to upload Nextflow logs to S3"
    else
        echo "Nextflow logs uploaded successfully to $s3_log_path"
    fi
else
    echo "Warning: Nextflow log file not found at $NEXTFLOW_LOG_PATH"
fi

# Upload trace, report, and timeline files if they exist
for file in trace.txt report.html timeline.html; do
    if [ -f "$file" ]; then
        s3_path="s3://${LOG_BUCKET}/jobs/${JOB_ID}/$file"
        aws s3 cp "$file" "$s3_path"
        if [ $? -eq 0 ]; then
            echo "Uploaded $file to $s3_path"
        fi
    fi
done

# Upload DAG file if it exists
if [ -f "pipeline_dag.html" ]; then
    s3_path="s3://${LOG_BUCKET}/jobs/${JOB_ID}/pipeline_dag.html"
    aws s3 cp "pipeline_dag.html" "$s3_path"
    if [ $? -eq 0 ]; then
        echo "Uploaded pipeline_dag.html to $s3_path"
    fi
fi
