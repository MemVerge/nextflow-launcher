#!/bin/bash
set -e

# Read and parse the job spec
job=$(cat "./job.json")

# Extract fields
name=$(echo "$job" | jq -r .name)
pipeline=$(echo "$job" | jq -r .pipeline)
profile=$(echo "$job" | jq -r .profile)
head_node_queue=$(echo "$job" | jq -r .head_node_queue)
task_queue=$(echo "$job" | jq -r .task_queue)
work_dir=$(echo "$job" | jq -r .work_dir)
result_dir=$(echo "$job" | jq -r .result_dir)
memory=$(echo "$job" | jq -r .memory)
max_retries=$(echo "$job" | jq -r .max_retries)

# Print parsed values
echo "Parsed values from job spec:"
echo "Name: $name"
echo "Pipeline: $pipeline"
echo "Profile: $profile"
echo "Head Node Queue: $head_node_queue"
echo "Task Queue: $task_queue"
echo "Work Dir: $work_dir"
echo "Result Dir: $result_dir"
echo "Memory: $memory"
echo "Max Retries: $max_retries"

# Generate aws.config
echo "Generating aws.config..."

cat > aws.config << 'EOF'
plugins {
    id "nf-amazon"
}
process {
    executor = "awsbatch"
    queue = "__TASK_QUEUE__"
    maxRetries = __MAX_RETRIES__
    memory = "__MEMORY__"
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

# Replace placeholders
sed -i '' "s/__TASK_QUEUE__/$task_queue/g" aws.config
sed -i '' "s/__MAX_RETRIES__/$max_retries/g" aws.config
sed -i '' "s/__MEMORY__/$memory/g" aws.config

echo -e "\nGenerated aws.config content:"
cat aws.config

echo -e "\nNextflow command that would be executed:"
echo "nextflow -log /var/log/nextflow/nextflow.log run $pipeline -profile $profile -work-dir $work_dir --outdir $result_dir -c aws.config -ansi-log false" 