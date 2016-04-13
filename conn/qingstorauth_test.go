package conn

import (
	"testing"
)

func TestBuildQingStorRaw(t *testing.T) {
	tests := []struct {
		params   *Dict
		authpath string
		headers  *Dict
		method   string
		out      string
	}{
		{
			&Dict{},
			"/",
			&Dict{},
			"get",
			"GET\n\n\n\n/",
		},
		{
			&Dict{},
			"/auth",
			&Dict{},
			"post",
			"POST\n\n\n\n/auth",
		},
		{
			&Dict{},
			"/auth",
			&Dict{"X-QS-Date": "aaa"},
			"get",
			"GET\n\n\nx-qs-date:aaa\n/auth",
		},
		{
			&Dict{"X-QS-Date": "bbb"},
			"/auth",
			&Dict{},
			"get",
			"GET\n\n\n\n/auth?X-QS-Date=bbb",
		},
		{
			&Dict{"X-QS-Date": "bbb"},
			"/",
			&Dict{"X-QS-Date": "aaa"},
			"get",
			"GET\n\n\nx-qs-date:aaa\n/?X-QS-Date=bbb",
		},
		{
			&Dict{"X-QS-Date": "bbb"},
			"/",
			&Dict{"Date": "ccc"},
			"get",
			"GET\n\n\n\n/?X-QS-Date=bbb",
		},
		{
			&Dict{"X-QS-Date": "bbb"},
			"/",
			&Dict{"Date": "ccc", "X-QS-Date": "aaa"},
			"get",
			"GET\n\n\nx-qs-date:aaa\n/?X-QS-Date=bbb",
		},
		{
			&Dict{},
			"/",
			&Dict{"Date": "ccc"},
			"get",
			"GET\n\n\nccc\n/",
		},
		{
			&Dict{"X-QS-Date": "bbb", "Debug": "on", "debug": "on"},
			"/",
			&Dict{"Date": "ccc", "X-QS-Date": "aaa", "X-QS-Debug": "on"},
			"get",
			"GET\n\n\nx-qs-date:aaa\nx-qs-debug:on\n/?Debug=on&X-QS-Date=bbb&debug=on",
		},
	}

	a := &QingStorAuth{}

	for _, test := range tests {
		raw, err := a.buildQingStorRawQuery(
			test.params,
			test.authpath,
			test.headers,
			test.method,
		)
		if err != nil {
			t.Error(err)
		}
		if raw != test.out {
			t.Errorf("%v != %v", raw, test.out)
		}
	}
}
