#!/bin/bash
set -e

echo "Downloading job spec from S3..."
job=$(aws s3 cp "${JOB_SPEC_PATH}" -)

# Parse fields
name=$(echo "$job" | jq -r .name)
pipeline=$(echo "$job" | jq -r .pipeline)
profile=$(echo "$job" | jq -r .profile)
job_queue=$(echo "$job" | jq -r .job_queue)
work_dir=$(echo "$job" | jq -r .work_dir)
result_dir=$(echo "$job" | jq -r .result_dir)
memory=$(echo "$job" | jq -r .memory)
max_retries=$(echo "$job" | jq -r .max_retries)
additional_config=$(echo "$job" | jq -r .additional_config)

# Fall back if values aren't set
memory=${memory:-20G}
max_retries=${max_retries:-5}

# Optional AWS keys
access_key="${AWS_ACCESS_KEY_ID:-}"
secret_key="${AWS_SECRET_ACCESS_KEY:-}"

echo "Generating aws.config dynamically..."

cat > aws.config <<EOF
plugins {
    id 'nf-amazon'
}
process {
    executor = 'awsbatch'
    queue = '${job_queue}'
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

# Append custom config block if present
if [ "$additional_config" != "null" ] && [ -n "$additional_config" ]; then
    echo -e "\n\n// ----- User-defined Config Start -----" >> aws.config
    echo -e "$additional_config" >> aws.config
    echo -e "// ----- User-defined Config End -----" >> aws.config
fi

echo "Launching pipeline: $pipeline"
nextflow run "$pipeline" \
  -profile "$profile" \
  -work-dir "$work_dir" \
  --outdir "$result_dir" \
  -c aws.config 