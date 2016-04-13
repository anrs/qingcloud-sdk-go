package conn

import (
	"testing"
)

func TestNewQingStorConnection(t *testing.T) {
	tests := []struct {
		httpConn HTTPConnection
		secure   bool
		out      QingStorConnection
	}{
		{
			HTTPConnection{},
			false,
			QingStorConnection{
				HTTPConnection: HTTPConnection{
					Host:        "qingstor.com",
					Protocol:    "http",
					AuthHandler: NewQingStorAuth("", ""),
				},
				Secure:      false,
				UserAgent:   "QingStor SDK Go",
				StyleFormat: NewVirtualHostStyleFormatter(),
			},
		},
		{
			HTTPConnection{
				Host:     "api.qingstor.com",
				Protocol: "https",
			},
			false,
			QingStorConnection{
				HTTPConnection: HTTPConnection{
					Host:        "api.qingstor.com",
					Protocol:    "http",
					AuthHandler: NewQingStorAuth("", ""),
				},
				Secure:      false,
				UserAgent:   "QingStor SDK Go",
				StyleFormat: NewVirtualHostStyleFormatter(),
			},
		},
		{
			HTTPConnection{
				Host:            "qingstor.com",
				Protocol:        "http",
				AccessKeyID:     "access",
				SecretAccessKey: "secret",
			},
			true,
			QingStorConnection{
				HTTPConnection: HTTPConnection{
					Host:            "qingstor.com",
					Protocol:        "https",
					AccessKeyID:     "access",
					SecretAccessKey: "secret",
					AuthHandler:     NewQingStorAuth("access", "secret"),
				},
				Secure:      true,
				UserAgent:   "QingStor SDK Go",
				StyleFormat: NewVirtualHostStyleFormatter(),
			},
		},
	}

	for _, test := range tests {
		c := NewQingStorConnection(test.httpConn, test.secure)
		CheckConnectionFields(t, c, test.out, "Secure", "UserAgent", "StyleFormat")
	}
}

func TestWrapHeaders(t *testing.T) {
}

func TestGetContentLength(t *testing.T) {
	tests := []struct {
		data string
		out  int64
	}{
		{"", 0},
		{"a", 1},
		{"abc", 3},
	}

	for _, test := range tests {
		if len := getContentLength(test.data); len != test.out {
			t.Errorf("lenght %d != %d", len, test.out)
		}
	}
}
