package types

import "time"

type Job struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Pipeline      string            `json:"pipeline"`
	Parameters    map[string]string `json:"parameters"`
	Status        string            `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	BatchJobId    string            `json:"batch_job_id"`
	Memory        string            `json:"memory"`
	MaxRetries    int               `json:"max_retries"`
	HeadNodeQueue string            `json:"head_node_queue"`
	TaskQueue     string            `json:"task_queue"`
	WorkDir       string            `json:"work_dir"`
	ResultDir     string            `json:"result_dir"`
	LogBucket     string            `json:"log_bucket"`
	AWSAccessKey  string            `json:"aws_access_key"`
	AWSSecretKey  string            `json:"aws_secret_key"`
}

type Jobs []Job
