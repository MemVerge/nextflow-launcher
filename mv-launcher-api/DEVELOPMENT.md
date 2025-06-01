# Local Development Setup

This guide will help you set up the Nextflow Launcher for local development.

## Prerequisites

1. Go 1.21 or later
2. Node.js 18 or later
3. AWS CLI configured
4. Docker (for running Nextflow containers)

## Setup Steps

1. Clone the repository:
```bash
git clone https://github.com/MemVerge/nf-launcher.git
cd nf-launcher
```

2. Set up AWS resources (one-time):
```bash
cd mv-launcher-api
./setup-account.sh dev
```
This will create S3 buckets, IAM roles, and a Batch job definition using CloudFormation, and save outputs to `.env.local`.

3. Install backend dependencies:
```bash
go mod download
```

4. Install frontend dependencies:
```bash
cd vue
npm install
```

5. Start the development servers:
```bash
cd mv-launcher-api
./dev.sh
```

The application will be available at:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

## API Endpoints

- `GET /health` - Health check
- `GET /v1/buckets` - List S3 buckets
- `GET /v1/pipelines` - List pipelines
- `GET /v1/jobs` - List jobs
- `POST /v1/jobs` - Submit a job
- `GET /v1/jobs/:id/logs` - Get job logs
- `GET /v1/jobs/:id/log-url` - Get presigned S3 log URL
- `GET /v1/batch/queues` - List AWS Batch queues

## Troubleshooting

1. If you see CORS errors:
   - The backend is configured to allow requests from http://localhost:5173
   - Check that the frontend is running on the correct port

2. If AWS API calls fail:
   - Verify your AWS credentials are correctly configured
   - Check that the required AWS resources exist
   - Ensure your IAM roles have the necessary permissions

3. If job submission fails:
   - Verify the AWS Batch job definition exists (should be `${ENVIRONMENT}-nextflow-headnode`)
   - Check that the job queue exists and is active
   - Ensure the container image is available in ECR

4. If the frontend can't connect to the backend:
   - Verify both servers are running (check the dev.sh output)
   - Check that the backend is running on port 8080
   - Ensure the frontend is configured to use the correct API URL 