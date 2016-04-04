package conn

import (
	"fmt"
	"testing"
)

type Test struct {
	cond    Condition
	params  *Dict
	ok      bool
}

func TestCheckRequired(t *testing.T) {
	tests := []Test{
		{
			Condition{Required: Keys{"offset"}},
			&Dict{"offset": 0},
			true,
		},
		{
			Condition{},
			&Dict{"verbose": 0},
			true,
		},
		{
			Condition{Required: Keys{"offset"}},
			&Dict{"verbose": 0},
			false,
		},
		{
			Condition{Integers: Keys{"offset"}},
			&Dict{"offset": 0},
			true,
		},
		{
			Condition{Integers: Keys{"offset"}},
			&Dict{},
			true,
		},
		{
			Condition{Integers: Keys{"offset"}},
			&Dict{"verbose": "on", "offset": 0},
			true,
		},
		{
			Condition{Integers: Keys{"offset"}},
			&Dict{"offset": "0"},
			false,
		},
		{
			Condition{Lists: Keys{"status"}},
			&Dict{"status": List{"active", "stopped"}},
			true,
		},
		{
			Condition{Lists: Keys{"status"}},
			&Dict{},
			true,
		},
		{
			Condition{Lists: Keys{"status"}},
			&Dict{"status": "active"},
			false,
		},
		{
			Condition{Required: Keys{"offset"}, Integers: Keys{"offset"}},
			&Dict{"offset": 0},
			true,
		},
		{
			Condition{Required: Keys{"offset"}, Integers: Keys{"offset"}},
			&Dict{"offset": "0"},
			false,
		},
	}

	for _, test := range tests {
		err := test.cond.Check(test.params)
		if test.ok && err != nil {
			t.Logf("Condition %v on %v: ", test.cond, test.params)
			t.Error(err)
		}
		if !test.ok && err == nil {
			msg := fmt.Sprintf("Condition %v wants an error on %v", test.cond, test.params)
			t.Error(msg)
		}
	}
}
