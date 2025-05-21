package types

type Bucket struct {
	Name string            `json:"name" example:"nxf-pipelines"`
	ARN  string            `json:"arn" example:"arn:aws:s3:::nxf-pipelines"`
	Tags map[string]string `json:"tags,omitempty"`
}

type Buckets []Bucket
