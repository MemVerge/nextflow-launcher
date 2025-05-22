package types

import "testing"

func TestJobVerify(t *testing.T) {
	j := Job{
		WorkDir:   "nextflow-workbucket",
		ResultDir: "nextflow-results",
	}
	j.Verify()
	if j.WorkDir != "s3://nextflow-workbucket" {
		t.Errorf("WorkDir should be prefixed with s3://, got %s", j.WorkDir)
	}
	if j.ResultDir != "s3://nextflow-results" {
		t.Errorf("ResultDir should be prefixed with s3://, got %s", j.ResultDir)
	}
}
