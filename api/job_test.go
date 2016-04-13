package api

import (
	"testing"

	"github.com/anrs/qingcloud-sdk-go/conn"
)

func NewTestJobAPI(t *testing.T) JobAPI {
	return NewJobAPI(NewTestIaaSConnection(t))
}

func TestDescribeJobs(t *testing.T) {
	a := NewTestJobAPI(t)
	args := conn.Dict{
		"limit": 10, "offset": 0,
	}

	jobs, err := a.DescribeJobs(args)
	if err != nil {
		t.Error(err)
	}

	CheckIaaSAPIResponse(t, jobs, "DescribeJobsResponse")
}
