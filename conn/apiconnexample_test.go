package conn

import (
	"io/ioutil"
	"strings"
	"testing"
)

var access string
var secret string

func newConnection(t *testing.T) APIConnection {
	if access == "" || secret == "" {
		dat, err := ioutil.ReadFile("access.key")
		if err != nil {
			t.Fatalf("access.key does not exist")
		}

		ret := strings.Split(string(dat), "\n")
		if len(ret) >= 2 {
			access, secret = ret[0], ret[1]
		} else {
			t.Fatalf("invalid access.key content")
		}

		if access == "" || secret == "" {
			t.Fatalf("invalid access.key content")
		}
	}

	c := NewAPIConnector(APIConnection{
		HTTPConnection{
			Host:            "api.qingcloud.com",
			Path:            "/iaas/",
			Protocol:        "http",
			AccessKeyID:     access,
			SecretAccessKey: secret,
		},
		"pek2",
	})

	v, ok := c.(APIConnection)
	if !ok {
		t.Fatalf("c is not an APIConnection")
	}

	return v
}

func check(t *testing.T, resp interface{}, action string) {
	if resp, ok := resp.(Dict); ok {
		if len(resp) < 2 {
			t.Fatalf("len of %v is less than 2", resp)
		}

		if a, ok := resp["action"]; !ok {
			t.Fatal("there is not an action key")
		} else if a != action {
			t.Fatalf("%s != %s", a, action)
		}
	}
}

func TestDescribeZones(t *testing.T) {
	c := newConnection(t)
	zones, err := c.DescribeZones()
	if err != nil {
		t.Error(err)
	}
	check(t, zones, "DescribeZonesResponse")
}

func TestDescribeJobs(t *testing.T) {
	c := newConnection(t)
	args := Dict{
		"limit": 10, "offset": 0,
	}
	jobs, err := c.DescribeJobs(args)
	if err != nil {
		t.Error(err)
	}
	check(t, jobs, "DescribeJobsResponse")
}
