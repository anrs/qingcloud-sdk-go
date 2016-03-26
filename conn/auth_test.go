package conn

import (
	"reflect"
	"testing"

	/*
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"*/
)

var accessKeyID string = "QYACCESSKEYIDEXAMPLE"
var secretAccessKey string = "SECRETACCESSKEY"
var signature string = "32bseYy39DOlatuewpeuW5vpmW51sD1A%2FJdGynqSpP8%3D"

func TestBuildQuery(t *testing.T) {
	tests := []struct {
		in  map[string]string
		out string
	}{
		{
			map[string]string{
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
				"access_key_id":     accessKeyID,
				"action":            "RunInstances",
				"time_stamp":        "2013-08-27T14:30:10Z",
			},
			"access_key_id=QYACCESSKEYIDEXAMPLE&action=RunInstances&count=1&image_id=centos64x86a&instance_name=demo&instance_type=small_b&login_mode=passwd&login_passwd=QingCloud20130712&signature_method=HmacSHA256&signature_version=1&time_stamp=2013-08-27T14%3A30%3A10Z&version=1&vxnets.1=vxnet-0&zone=pek1",
		},
	}

	for _, test := range tests {
		q := BuildRawQuery(test.in)
		if q != test.out {
			t.Errorf("query %v != %v", q, test.out)
		}
	}
}

type UtilsInterface interface {
	UTCTimestamp() string
}

func TestAuthorize(t *testing.T) {
	tests := []struct{
		in  map[string]string
		out map[string]string
	}{
		{
			map[string]string{
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
				"access_key_id":     accessKeyID,
				"action":            "RunInstances",
				"time_stamp":        "2013-08-27T14:30:10Z",
			},
			map[string]string{
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
				"access_key_id":     accessKeyID,
				"action":            "RunInstances",
				"time_stamp":        "2013-08-27T14:30:10Z",
				"signature":         "32bseYy39DOlatuewpeuW5vpmW51sD1A/JdGynqSpP8=",
			},
		},
	}

	a := NewQuerySignatureAuth(accessKeyID, secretAccessKey)

	for _, test := range tests {
		if err := a.Authorize(&test.in, "/iaas/"); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(test.in, test.out) {
			t.Errorf("signed %v != %v", test.in, test.out)
		}
	}
}

func TestAuthorizeNewUTCTimestamp(t *testing.T) {
	a := NewQuerySignatureAuth("", "")
	in := map[string]string{}

	if err := a.Authorize(&in, "/iaas/"); err != nil {
		t.Error(err)
	}

	if len(in["time_stamp"]) != 20 {
		t.Errorf("len(%v) != 20", in["time_stamp"])
	}
}
