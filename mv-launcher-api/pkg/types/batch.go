package types

type Queue struct {
	Name string `json:"name" example:"spot-xs"`
	ARN  string `json:"arn" example:"arn:aws:batch:us-west-2:123456789012:job-queue/spot-xs"`
}

type Queues []Queue
