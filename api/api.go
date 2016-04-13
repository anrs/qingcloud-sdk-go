package api

import (
	"github.com/anrs/qingcloud-sdk-go/conn"
)

type API struct {
	connector conn.HTTPConnector
}

type ReqArg struct {
	action   string
	params   conn.Dict
	path     string
	method   string
	headers  conn.Dict
	authpath string
	data     string
}
