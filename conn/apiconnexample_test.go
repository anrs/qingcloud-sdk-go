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
		access, secret = ret[0], ret[1]

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

func TestDescribeZones(t *testing.T) {
	c := newConnection(t)
	zones, err := c.DescribeZones()
	if err != nil {
		t.Error(err)
	}

	var ok bool
	if zones, ok = zones.(map[string]interface{}); ok {
		if len(zones.(map[string]interface{})) < 1 {
			t.Errorf("len of %v less than 1", zones)
		}
	} else {
		t.Log(zones)
	}
}
