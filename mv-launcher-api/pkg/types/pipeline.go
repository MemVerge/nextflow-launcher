package types

type Pipeline struct {
	Name       string            `json:"name" example:"nextflow-pipeline"`
	Image      string            `json:"image" example:"registry.gitlab.com/qnib-pub-containers/qnib/nextflow-workflow-run:24.10.4-1"`
	Command    string            `json:"command" example:"start.sh Ref::pipeline Ref::work-dir Ref::result-dir"`
	Parameters map[string]string `json:"parameters" example:"{'pipeline': 'hello', 'work-dir': 'addme', 'result-dir': 'addme'}"`
	Memory     string            `json:"memory" example:"2048"`
	VCPUs      string            `json:"vcpus" example:"1"`
}

type Pipelines []Pipeline
