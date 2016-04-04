package conn

import (
	"reflect"
	"testing"
)

func TestWrapParams(t *testing.T) {
	req_id := "8608eb56894f425b833251e2fd9955e3"

	tests := []struct {
		in  Dict
		out Dict
	}{
		{
			Dict{
				"action": "DescribeZones",
				"zone":   "pek1",
				"req_id": req_id,
				"ver":    1,
			},
			Dict{
				"action": "DescribeZones",
				"zone":   "pek1",
				"req_id": req_id,
				"ver":    "1",
			},
		},
		{
			Dict{
				"action": "DescribeInstances",
				"zone":   "pek1",
				"req_id": req_id,
				"ver":    1,
				"status": []string{"running", "stopped"},
			},
			Dict{
				"action":   "DescribeInstances",
				"zone":     "pek1",
				"req_id":   req_id,
				"ver":      "1",
				"status.1": "running",
				"status.2": "stopped",
			},
		},
		{
			Dict{
				"action": "CreateSomething",
				"zone":   "pek1",
				"req_id": req_id,
				"ver":    1,
				"ip": []interface{}{
					Dict{"master": "192.168.100.100"},
					Dict{"slave": "192.168.100.200"},
				},
			},
			Dict{
				"action":      "CreateSomething",
				"zone":        "pek1",
				"req_id":      req_id,
				"ver":         "1",
				"ip.1.master": "192.168.100.100",
				"ip.2.slave":  "192.168.100.200",
			},
		},
		{
			Dict{
				"action": "UpdateSomething",
				"zone":   "pek1",
				"req_id": req_id,
				"ver":    1,
				"ip": []interface{}{
					Dict{"master": []string{"192.168.100.100"}},
					Dict{"slave": []string{"192.168.100.200"}},
				},
			},
			Dict{
				"action":      "UpdateSomething",
				"zone":        "pek1",
				"req_id":      req_id,
				"ver":         "1",
				"ip.1.master": "[\"192.168.100.100\"]",
				"ip.2.slave":  "[\"192.168.100.200\"]",
			},
		},
		{
			Dict{
				"action": "UpdateSomething",
				"zone":   "pek1",
				"req_id": req_id,
				"ver":    1,
				"ip": []interface{}{
					Dict{
						"master": Dict{"ip": "192.168.100.100"},
					},
					Dict{
						"slave": Dict{"ip": "192.168.100.200"},
					},
				},
			},
			Dict{
				"action":      "UpdateSomething",
				"zone":        "pek1",
				"req_id":      req_id,
				"ver":         "1",
				"ip.1.master": "{\"ip\":\"192.168.100.100\"}",
				"ip.2.slave":  "{\"ip\":\"192.168.100.200\"}",
			},
		},
	}

	c := &HTTPConnection{}

	for _, test := range tests {
		params, err := c.WrapParams(test.in)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(params, test.out) {
			t.Errorf("%v != %v", params, test.out)
		}
	}
}
