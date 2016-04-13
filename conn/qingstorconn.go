package conn

import (
	"net/http"
	"strings"

	"github.com/anrs/qingcloud-sdk-go/utils"
)

type QingStorConnection struct {
	HTTPConnection

	UserAgent   string
	Secure      bool
	StyleFormat StyleFormatter
}

func (c QingStorConnection) BuildRequest(
	method string,
	path string,
	params Dict,
	authpath string,
	headers Dict,
	host string,
	data string,
) (*http.Request, error) {
	return nil, nil
}

func (c QingStorConnection) Send(
	method string,
	path string,
	params Dict,
	authpath string,
	headers Dict,
	host string,
	data string,
) (*http.Response, error) {
	return nil, nil
}

func NewQingStorConnection(httpConn HTTPConnection, secure bool) QingStorConnection {
	if strings.Trim(httpConn.Host, " ") == "" {
		httpConn.Host = "qingstor.com"
	}

	if secure {
		httpConn.Protocol = "https"
	} else {
		httpConn.Protocol = "http"
	}

	httpConn.AuthHandler = NewQingStorAuth(httpConn.AccessKeyID, httpConn.SecretAccessKey)

	c := QingStorConnection{
		HTTPConnection: httpConn,
		UserAgent:      "QingStor SDK Go",
		Secure:         secure,
		StyleFormat:    NewVirtualHostStyleFormatter(),
	}

	var connector HTTPConnector = c
	c.connector = connector
	return c
}

func getContentLength(data interface{}) int64 {
	var len int64 = 1
	return len
}

func (c QingStorConnection) wrapHeaders(headers *Dict, host string, data interface{}) {
	if _, ok := (*headers)["Host"]; !ok {
		(*headers)["Host"] = host
	}

	if _, ok := (*headers)["Date"]; !ok {
		(*headers)["Date"] = utils.GMTime()
	}

	if _, ok := (*headers)["Content-Length"]; !ok {
		(*headers)["Content-Length"] = getContentLength(data)
	}
}

func (c QingStorConnection) SendRequest(
	method string,
	bucket string,
	key string,
	headers Dict,
	data interface{},
	params Dict,
) (interface{}, error) {
	host := c.StyleFormat.BuildHost(c.Host, bucket)
	//path := c.StyleFormat.BuildPathBase(bucket, key)
	//authpath := c.StyleFormat.BuildAuthPath(bucket, key)

	if headers == nil {
		headers = make(Dict)
	}
	c.wrapHeaders(&headers, host, data)

	return nil, nil
}
