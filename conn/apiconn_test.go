package conn

import (
	"reflect"
	"testing"
)

func TestNewAPIConnector(t *testing.T) {
	tests := []struct {
		in  APIConnection
		out APIConnection
	}{
		{
			APIConnection{
				HTTPConnection{},
				"pek3a",
			},
			APIConnection{
				HTTPConnection{
					Host:     "api.qingcloud.com",
					Port:     443,
					Protocol: "https",
					Timeout:  60,
					Debug:    false,
					AuthHandler: NewQuerySignatureAuth("", ""),
				},
				"pek3a",
			},
		},
		{
			APIConnection{
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
			},
			APIConnection{
				HTTPConnection{
					AccessKeyID:     "access",
					SecretAccessKey: "secret",
					Host:            "qingcloud.com",
					Port:            80,
					Protocol:        "http",
					Expires:         1024,
					Timeout:         10,
					Debug:           true,
					AuthHandler: NewQuerySignatureAuth("access", "secret"),
				},
				"pek2",
			},
		},
	}

	for _, test := range tests {
		c := NewAPIConnector(test.in)
		c, ok := c.(APIConnection)
		if !ok {
			t.Errorf("%v is not APIConnection", c)
		}

		if !reflect.DeepEqual(c, test.out) {
			t.Errorf("%v != %v", c, test.out)
		}
	}
}
