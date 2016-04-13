package api

import (
	"github.com/anrs/qingcloud-sdk-go/conn"
)

type IaaSAPI struct {
	API
}

func (a IaaSAPI) SendAPIRequest(arg ReqArg) (interface{}, error) {
	if arg.params == nil {
		arg.params = make(conn.Dict)
	}

	if arg.headers == nil {
		arg.headers = make(conn.Dict)
	}

	if arg.method == "" {
		arg.method = "GET"
	}

	return a.connector.(conn.IaaSConnection).SendRequest(
		arg.action,
		arg.params,
		arg.path,
		arg.method,
		arg.headers,
		arg.authpath,
		arg.data,
	)
}
