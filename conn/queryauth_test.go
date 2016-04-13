package conn

import (
	"reflect"
	"testing"
)

var queryAccessKeyID string = "QYACCESSKEYIDEXAMPLE"
var querySecretAccessKey string = "SECRETACCESSKEY"

func TestQueryAuthorize(t *testing.T) {
	tests := []struct {
		in  Dict
		out Dict
	}{
		{
			Dict{
				"count":             "1",
				"vxnets.1":          "vxnet-0",
				"zone":              "pek1",
				"instance_type":     "small_b",
				"signature_version": "1",
				"signature_method":  "HmacSHA256",
				"instance_name":     "demo",
				"image_id":          "centos64x86a",
				"login_mode":        "passwd",
				"login_passwd":      "QingCloud20130712",
				"version":           "1",
				"access_key_id":     queryAccessKeyID,
				"action":            "RunInstances",
				"time_stamp":        "2013-08-27T14:30:10Z",
			},
			Dict{
				"count":             "1",
				"vxnets.1":          "vxnet-0",
				"zone":              "pek1",
				"instance_type":     "small_b",
				"signature_version": "1",
				"signature_method":  "HmacSHA256",
				"instance_name":     "demo",
				"image_id":          "centos64x86a",
				"login_mode":        "passwd",
				"login_passwd":      "QingCloud20130712",
				"version":           "1",
				"access_key_id":     queryAccessKeyID,
				"action":            "RunInstances",
				"time_stamp":        "2013-08-27T14:30:10Z",
				"signature":         "32bseYy39DOlatuewpeuW5vpmW51sD1A/JdGynqSpP8=",
			},
		},
	}

	a := NewQueryAuth(queryAccessKeyID, querySecretAccessKey)

	for _, test := range tests {
		if err := a.Authorize(&test.in, "/iaas/", nil, "get"); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(test.in, test.out) {
			t.Errorf("signed %v != %v", test.in, test.out)
		}
	}
}
