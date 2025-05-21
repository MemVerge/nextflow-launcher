package api

type BatchQueue struct {
	Name   string `json:"name"`
	State  string `json:"state"`
	Status string `json:"status"`
	ARN    string `json:"arn"`
}

type BatchQueues []BatchQueue

type AWSBatchConfig struct {
	Region           string   `json:"region"`
	ComputeEnvName   string   `json:"compute_env_name"`
	JobQueueName     string   `json:"job_queue_name"`
	InstanceTypes    []string `json:"instance_types"`
	MinvCpus         int      `json:"min_vcpus"`
	MaxvCpus         int      `json:"max_vcpus"`
	DesiredvCpus     int      `json:"desired_vcpus"`
	SubnetId         string   `json:"subnet_id"`
	SecurityGroupId  string   `json:"security_group_id"`
	UseSpot          bool     `json:"use_spot"`
	EnableMultiQueue bool     `json:"enable_multi_queue"`
	UniquePrefix     string   `json:"unique_prefix"`
}
