package api

import (
	"github.com/anrs/qingcloud-sdk-go/conn"
)

type JobAPI struct {
	IaaSAPI
}

func NewJobAPI(c conn.IaaSConnection) JobAPI {
	j := JobAPI{IaaSAPI{
		API{c},
	}}
	j.connector = c
	return j
}

func (a JobAPI) DescribeJobs(params conn.Dict) (interface{}, error) {
	cond := conn.Condition{
		Integers: conn.Keys{"offset", "limit"},
		Lists:    conn.Keys{"jobs"},
	}
	if err := cond.Check(&params); err != nil {
		return nil, err
	}

	arg := ReqArg{action: "DescribeJobs", params: params}
	return a.SendAPIRequest(arg)
}
