package api

import (
	_ "github.com/anrs/qingcloud-sdk-go/conn"
)

type ZoneAPI struct {
	IaaSAPI
}

func (a ZoneAPI) DescribeZones() (interface{}, error) {
	arg := ReqArg{action: "DescribeZones"}
	return a.SendAPIRequest(arg)
}
