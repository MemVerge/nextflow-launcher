#!/bin/bash
set -e
job=$(cat "./job.json")
task_queue=$(echo "$job" | jq -r .task_queue)
memory=$(echo "$job" | jq -r .memory)
max_retries=$(echo "$job" | jq -r .max_retries)
memory=${memory:-"20G"}
max_retries=${max_retries:-5}
echo "Task Queue: $task_queue"
echo "Memory: $memory"
echo "Max Retries: $max_retries"
cat > aws.config <<EOF
plugins {
    id "nf-amazon"
}
process {
    executor = "awsbatch"
    queue = ""
    maxRetries = 
    memory = ""
}
process.containerOptions = "--env MMC_CHECKPOINT_DIAGNOSIS=true --env MMC_CHECKPOINT_IMAGE_SUBPATH=nextflow --env MMC_CHECKPOINT_INTERVAL=5m --env MMC_CHECKPOINT_MODE=true --env MMC_CHECKPOINT_IMAGE_PATH=/mmc-checkpoint"
aws {
    region = "us-west-2"
    client {
        maxConnections = 20
        connectionTimeout = 10000
        uploadStorageClass = "INTELLIGENT_TIERING"
        storageEncryption = "AES256"
    }
    batch {
        cliPath = "/nextflow_awscli/bin/aws"
        maxTransferAttempts = 3
        delayBetweenAttempts = "5 sec"
    }
}
EOF
echo -e "
Generated aws.config content:"
cat aws.config
