package conn

import (
	"testing"
)

func TestNewIaaSConnection(t *testing.T) {
	tests := []struct {
		httpConn HTTPConnection
		zone     string
		out      IaaSConnection
	}{
		{
			HTTPConnection{},
			"pek3a",
			IaaSConnection{
				HTTPConnection{
					Host:        "api.qingcloud.com",
					Port:        443,
					Protocol:    "https",
					Timeout:     60,
					Debug:       false,
					AuthHandler: NewQueryAuth("", ""),
				},
				"pek3a",
			},
		},
		{
			HTTPConnection{
				AccessKeyID:     "access",
				SecretAccessKey: "secret",
				Host:            "qingcloud.com",
				Port:            80,
				Protocol:        "http",
				Expires:         1024,
				Timeout:         10,
				Debug:           true,
			},
			"pek2",
			IaaSConnection{
				HTTPConnection{
					AccessKeyID:     "access",
					SecretAccessKey: "secret",
					Host:            "qingcloud.com",
					Port:            80,
					Protocol:        "http",
					Expires:         1024,
					Timeout:         10,
					Debug:           true,
					AuthHandler:     NewQueryAuth("access", "secret"),
				},
				"pek2",
			},
		},
	}

	for _, test := range tests {
		c := NewIaaSConnection(test.httpConn, test.zone)
		CheckConnectionFields(t, c, test.out, "Zone")
	}
}
