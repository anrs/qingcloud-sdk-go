package api

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/anrs/qingcloud-sdk-go/conn"
)

var access string
var secret string

func NewTestIaaSConnection(t *testing.T) conn.IaaSConnection {
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

	c := conn.NewIaaSConnection(
		conn.HTTPConnection{
			Host:            "api.qingcloud.com",
			Path:            "/iaas/",
			Protocol:        "http",
			AccessKeyID:     access,
			SecretAccessKey: secret,
		},
		"pek2",
	)

	return c
}

func CheckIaaSAPIResponse(t *testing.T, resp interface{}, action string) {
	if resp, ok := resp.(map[string]interface{}); ok {
		if retcode, ok := resp["ret_code"]; !ok {
			t.Fatalf("action %s response has no ret_code", action)
		} else if retcode.(float64) != 0 {
			t.Fatalf("action %s ret_code %v != 0", action, retcode)
		}

		if len(resp) < 2 {
			t.Fatalf("len of %v is less than 2", resp)
		}

		if a, ok := resp["action"]; !ok {
			t.Fatal("there is not an action key")
		} else if a != action {
			t.Fatalf("%s != %s", a, action)
		}
	} else {
		t.Fatal("response %v is not a map", resp)
	}
}

func TestAPI(t *testing.T) {
}
