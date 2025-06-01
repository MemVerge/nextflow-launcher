#!/bin/bash

# Check if .env.local exists
if [ ! -f .env.local ]; then
    echo "Error: .env.local not found. Please run setup-account.sh first."
    exit 1
fi

# Load environment variables
echo "Loading environment variables from .env.local"
export $(cat .env.local | grep -v '^#' | xargs)

# Check if AWS credentials are configured
if [ -z "$AWS_ACCESS_KEY_ID" ] || [ -z "$AWS_SECRET_ACCESS_KEY" ]; then
    echo "Error: AWS credentials not found in .env.local"
    exit 1
fi

# Check if required AWS resources are configured
if [ -z "$PIPELINE_BUCKET" ] || [ -z "$JOB_BUCKET" ] || [ -z "$LOG_BUCKET" ] || [ -z "$NEXTFLOW_HEADNODE_JOB_DEFINITION" ]; then
    echo "Error: Required AWS resources not found in .env.local"
    echo "Please run setup-account.sh to create the required resources"
    exit 1
fi

# Start the backend in the background
echo "Starting backend server..."
go run main.go &
BACKEND_PID=$!

# Wait for backend to be ready
echo "Waiting for backend to be ready..."
for i in {1..30}; do
    if curl -s http://localhost:${PORT}/health > /dev/null; then
        echo "Backend is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "Error: Backend failed to start"
        kill $BACKEND_PID
        exit 1
    fi
    sleep 1
done

# Start the frontend
echo "Starting frontend development server..."
cd vue
npm run dev &
FRONTEND_PID=$!

# Function to handle script termination
cleanup() {
    echo "Shutting down..."
    kill $BACKEND_PID
    kill $FRONTEND_PID
    exit 0
}

# Set up trap to catch termination signal
trap cleanup SIGINT SIGTERM

# Wait for both processes
wait $BACKEND_PID $FRONTEND_PID 