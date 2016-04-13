package conn

import (
	"testing"
)

func TestVirtualHostStyleFormatterBuildHost(t *testing.T) {
	tests := []struct {
		host   string
		bucket string
		out    string
	}{
		{"", "", ""},
		{"", "bucket", "bucket."},
		{"qingstor.com", "", "qingstor.com"},
		{"qingstor.com", "bucket", "bucket.qingstor.com"},
	}

	f := VirtualHostStyleFormatter{}

	for _, test := range tests {
		if host := f.BuildHost(test.host, test.bucket); host != test.out {
			t.Errorf("%s != %s", host, test.out)
		}
	}
}

func TestVirtualHostStyleFormatterBuildPathBase(t *testing.T) {
	tests := []struct {
		bucket string
		key    string
		out    string
	}{
		{"", "", "/"},
		{"bucket", "", "/"},
		{"bucket", "key", "/key"},
		{"bucket", "a+b", "/a%2Bb"},
		{"bucket", "a b", "/a+b"},
	}

	f := VirtualHostStyleFormatter{}

	for _, test := range tests {
		if base := f.BuildPathBase(test.bucket, test.key); base != test.out {
			t.Errorf("%s != %s", base, test.out)
		}
	}
}

func TestVirtualHostStyleFormatterBuildAuthPath(t *testing.T) {
	tests := []struct {
		bucket string
		key    string
		out    string
	}{
		{"", "", "/"},
		{"bucket", "", "/bucket"},
		{"bucket", "key", "/bucket/key"},
		{"bucket", "a+b", "/bucket/a+b"},
		{"bucket", "a b", "/bucket/a b"},
	}

	f := VirtualHostStyleFormatter{}

	for _, test := range tests {
		if base := f.BuildAuthPath(test.bucket, test.key); base != test.out {
			t.Errorf("%s != %s", base, test.out)
		}
	}
}
