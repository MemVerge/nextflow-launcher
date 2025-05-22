# Nextflow Launcher

A web-based launcher for Nextflow pipelines on AWS Batch.

## Project Structure

```
mv-launcher/
├── api/                    # Backend API service
│   ├── cmd/               # Command-line entry points
│   ├── internal/          # Private application code
│   │   ├── api/          # API handlers
│   │   ├── config/       # Configuration
│   │   └── service/      # Business logic
│   ├── pkg/              # Public libraries
│   └── scripts/          # Build and utility scripts
├── frontend/              # Vue.js frontend
│   ├── src/              # Source code
│   ├── public/           # Static files
│   └── dist/             # Build output
├── docs/                  # Documentation
├── deploy/               # Deployment configurations
│   ├── docker/          # Docker-related files
│   └── kubernetes/      # K8s manifests
├── scripts/              # Project-wide scripts
└── test/                 # Test environments and data
```

## Prerequisites

- Go 1.21 or later
- Node.js 16 or later
- Docker
- AWS CLI configured with appropriate credentials

## Environment Variables

The following environment variables need to be set:

```bash
# AWS Credentials
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=your_region

# S3 Bucket Configuration
PIPELINE_BUCKET=your_pipeline_bucket
JOB_BUCKET=your_job_bucket
LOG_BUCKET=your_log_bucket

# AWS Batch Configuration
JOB_ROLE_ARN=your_job_role_arn
```

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/MemVerge/nextflow-launcher.git
   cd nextflow-launcher
   ```

2. Set up the backend:
   ```bash
   cd api
   go mod download
   go run cmd/main.go
   ```

3. Set up the frontend:
   ```bash
   cd frontend
   npm install
   npm run serve
   ```

## Building

1. Build the backend:
   ```bash
   cd api
   go build -o bin/server cmd/main.go
   ```

2. Build the frontend:
   ```bash
   cd frontend
   npm run build
   ```

## Docker Deployment

1. Build the images:
   ```bash
   docker build -t nextflow-launcher-api -f deploy/docker/api.Dockerfile .
   docker build -t nextflow-launcher-frontend -f deploy/docker/frontend.Dockerfile .
   ```

2. Run the containers:
   ```bash
   docker run -p 8080:8080 nextflow-launcher-api
   docker run -p 80:80 nextflow-launcher-frontend
   ```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 