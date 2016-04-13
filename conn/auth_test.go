package conn

import (
	"testing"
)

func TestSign(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"a", "lhWpXUozYRjENbnNVMXoZEq5VrVzqikmJ0oSgLZnRxM="},
		{"", "thNnmggU2ex3L5XXeMNfxf8Wl8STcVZTxscSFEKSxa0="},
	}

	a := NewQueryAuth("", "")

	for _, test := range tests {
		signature, err := a.Sign(test.in)
		if err != nil {
			t.Error(err)
		}
		if signature != test.out {
			t.Errorf("%s != %s", signature, test.out)
		}
	}
}

func TestBuildQuery(t *testing.T) {
	tests := []struct {
		in  Dict
		out string
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
				"access_key_id":     "QYACCESSKEYIDEXAMPLE",
				"action":            "RunInstances",
				"time_stamp":        "2013-08-27T14:30:10Z",
			},
			"access_key_id=QYACCESSKEYIDEXAMPLE&action=RunInstances&count=1&image_id=centos64x86a&instance_name=demo&instance_type=small_b&login_mode=passwd&login_passwd=QingCloud20130712&signature_method=HmacSHA256&signature_version=1&time_stamp=2013-08-27T14%3A30%3A10Z&version=1&vxnets.1=vxnet-0&zone=pek1",
		},
	}

	for _, test := range tests {
		q, err := BuildRawQuery(test.in)
		if err != nil {
			t.Error(err)
		}
		if q != test.out {
			t.Errorf("query %v != %v", q, test.out)
		}
	}
}

func TestAuthorizeNewUTCTimestamp(t *testing.T) {
	a := NewQueryAuth("", "")
	in := Dict{}

	if err := a.Authorize(&in, "/iaas/", nil, "get"); err != nil {
		t.Error(err)
	}

	if len(in["time_stamp"].(string)) != 20 {
		t.Errorf("len(%v) != 20", in["time_stamp"])
	}
}
