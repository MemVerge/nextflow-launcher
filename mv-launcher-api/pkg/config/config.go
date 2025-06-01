package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration values
type Config struct {
	// AWS Configuration
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string

	// Environment
	Environment string

	// S3 Buckets
	PipelineBucket string
	JobBucket      string
	LogBucket      string

	// Job Configuration
	NextflowImage   string
	NextflowVCPUs   int32
	NextflowMemory  int32
	NextflowWorkDir string
	NextflowLogPath string

	// Server Configuration
	Port               int
	CORSAllowedOrigins []string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		// AWS Configuration
		AWSRegion:          getEnvOrDefault("AWS_REGION", "us-west-2"),
		AWSAccessKeyID:     getEnvOrDefault("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnvOrDefault("AWS_SECRET_ACCESS_KEY", ""),

		// Environment
		Environment: getEnvOrDefault("ENVIRONMENT", "dev"),

		// S3 Buckets
		PipelineBucket: getEnvOrDefault("PIPELINE_BUCKET", ""),
		JobBucket:      getEnvOrDefault("JOB_BUCKET", ""),
		LogBucket:      getEnvOrDefault("LOG_BUCKET", ""),

		// Job Configuration
		NextflowImage:   getEnvOrDefault("NEXTFLOW_IMAGE", "achyutha98/nextflow:latest"),
		NextflowVCPUs:   getEnvInt32OrDefault("NEXTFLOW_VCPUS", 4),
		NextflowMemory:  getEnvInt32OrDefault("NEXTFLOW_MEMORY", 16384),
		NextflowWorkDir: getEnvOrDefault("NEXTFLOW_WORK_DIR", "/workspace/work"),
		NextflowLogPath: getEnvOrDefault("NEXTFLOW_LOG_PATH", "/var/log/nextflow/nextflow.log"),

		// Server Configuration
		Port:               getEnvIntOrDefault("PORT", 8080),
		CORSAllowedOrigins: getEnvStringSliceOrDefault("CORS_ALLOWED_ORIGINS", []string{"http://localhost:5173"}),
	}

	// Validate required fields
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate checks if all required configuration values are set
func (c *Config) validate() error {
	required := map[string]string{
		"AWS_ACCESS_KEY_ID":     c.AWSAccessKeyID,
		"AWS_SECRET_ACCESS_KEY": c.AWSSecretAccessKey,
		"PIPELINE_BUCKET":       c.PipelineBucket,
		"JOB_BUCKET":            c.JobBucket,
		"LOG_BUCKET":            c.LogBucket,
	}

	for name, value := range required {
		if value == "" {
			return fmt.Errorf("required environment variable %s is not set", name)
		}
	}

	return nil
}

// Helper functions for environment variables
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt32OrDefault(key string, defaultValue int32) int32 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 32); err == nil {
			return int32(intValue)
		}
	}
	return defaultValue
}

func getEnvStringSliceOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return []string{value} // For now, just split by comma if needed
	}
	return defaultValue
}
