package conn

import (
	"reflect"
	"testing"
)

func TestWrapParams(t *testing.T) {
	req_id := "8608eb56894f425b833251e2fd9955e3"

	tests := []struct {
		in  map[string]interface{}
		out map[string]string
	}{
		{
			map[string]interface{}{
				"action": "DescribeZones",
				"zone":   "pek1",
				"req_id": req_id,
				"ver":    1,
			},
			map[string]string{
				"action": "DescribeZones",
				"zone":   "pek1",
				"req_id": req_id,
				"ver":    "1",
			},
		},
		{
			map[string]interface{}{
				"action": "DescribeInstances",
				"zone":   "pek1",
				"req_id": req_id,
				"ver":    1,
				"status": []string{"running", "stopped"},
			},
			map[string]string{
				"action":   "DescribeInstances",
				"zone":     "pek1",
				"req_id":   req_id,
				"ver":      "1",
				"status.1": "running",
				"status.2": "stopped",
			},
		},
		{
			map[string]interface{}{
				"action": "CreateSomething",
				"zone":   "pek1",
				"req_id": req_id,
				"ver":    1,
				"ip": []interface{}{
					map[string]string{"master": "192.168.100.100"},
					map[string]string{"slave": "192.168.100.200"},
				},
			},
			map[string]string{
				"action":      "CreateSomething",
				"zone":        "pek1",
				"req_id":      req_id,
				"ver":         "1",
				"ip.1.master": "192.168.100.100",
				"ip.2.slave":  "192.168.100.200",
			},
		},
		{
			map[string]interface{}{
				"action": "UpdateSomething",
				"zone":   "pek1",
				"req_id": req_id,
				"ver":    1,
				"ip": []interface{}{
					map[string]interface{}{"master": []string{"192.168.100.100"}},
					map[string]interface{}{"slave": []string{"192.168.100.200"}},
				},
			},
			map[string]string{
				"action":      "UpdateSomething",
				"zone":        "pek1",
				"req_id":      req_id,
				"ver":         "1",
				"ip.1.master": "[\"192.168.100.100\"]",
				"ip.2.slave":  "[\"192.168.100.200\"]",
			},
		},
		{
			map[string]interface{}{
				"action": "UpdateSomething",
				"zone":   "pek1",
				"req_id": req_id,
				"ver":    1,
				"ip": []interface{}{
					map[string]interface{}{
						"master": map[string]string{"ip": "192.168.100.100"},
					},
					map[string]interface{}{
						"slave": map[string]string{"ip": "192.168.100.200"},
					},
				},
			},
			map[string]string{
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
