package types

import "strings"

type Job struct {
	ID               string `json:"id,omitempty"`
	Name             string `json:"name" example:"job-1234"`
	Pipeline         string `json:"pipeline" example:"hello"`
	Profile          string `json:"profile,omitempty" example:"test"`
	HeadNodeQueue    string `json:"head_node_queue" example:"spot-xs"`
	TaskQueue        string `json:"task_queue" example:"spot-xs"`
	WorkDir          string `json:"work_dir" example:"s3://nextflow-workbucket"`
	ResultDir        string `json:"result_dir" example:"s3://nextflow-results"`
	InputDir         string `json:"input_dir,omitempty" example:"s3://nextflow-input"`
	LogBucket        string `json:"log_bucket" example:"nextflow-logs-memverge-launcher"`
	Memory           string `json:"memory"`
	MaxRetries       int    `json:"max_retries"`
	AdditionalConfig string `json:"additional_config"`
	AWSAccessKey     string `json:"aws_access_key"`
	AWSSecretKey     string `json:"aws_secret_key"`
}

type Jobs []Job

func (j *Job) Verify() {
	if !strings.HasPrefix(j.WorkDir, "s3://") {
		j.WorkDir = "s3://" + j.WorkDir
	}
	if !strings.HasPrefix(j.ResultDir, "s3://") {
		j.ResultDir = "s3://" + j.ResultDir
	}
	if !strings.HasPrefix(j.InputDir, "s3://") {
		j.InputDir = "s3://" + j.InputDir
	}
	if j.Profile == "none" {
		j.Profile = ""
	}
	if j.Memory == "" {
		j.Memory = "20G"
	}
	if j.MaxRetries == 0 {
		j.MaxRetries = 5
	}
}
