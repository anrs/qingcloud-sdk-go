package conn

import (
	"errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type APIConnection struct {
	HTTPConnection

	Zone string
}

func NewAPIConnector(c APIConnection) HTTPConnector {
	if strings.Trim(c.Host, " ") == "" {
		c.Host = "api.qingcloud.com"
	}

	if c.Port == 0 {
		c.Port = 443
	}

	if strings.Trim(c.Protocol, " ") == "" {
		c.Protocol = "https"
	}

	if c.Timeout == 0 {
		c.Timeout = 60
	}

	c.AuthHandler = NewQuerySignatureAuth(c.AccessKeyID, c.SecretAccessKey)

	var connector HTTPConnector = c
	return connector
}

func (c APIConnection) BuildRequest(
	method string,
	path string,
	params Dict,
	authpath string,
	headers Dict,
	host string,
	data string,
) (*http.Request, error) {
	wrappedParams, err := c.WrapParams(params)
	if err != nil {
		return nil, err
	}

	if strings.Trim(path, " ") == "" {
		path = c.Path
	}

	if err := c.AuthHandler.Authorize(&wrappedParams, path); err != nil {
		return nil, err
	}

	query, err := BuildRawQuery(wrappedParams)
	if err != nil {
		return nil, err
	}

	u := &url.URL{
		Scheme:   c.Protocol,
		Host:     host,
		Path:     path,
		RawQuery: query,
	}

	method = strings.Trim(strings.ToUpper(method), " ")
	req, err := http.NewRequest(method, u.String(), strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		v, ok := v.(string)
		if !ok {
			return nil, errors.New(fmt.Sprintf("%s value is not a string", k))
		}
		req.Header.Set(k, v)
	}

	if method == "POST" {
		req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "text/plain")
		req.Header.Set("Connection", "Keep-Alive")
	}

	return req, nil
}

func (c APIConnection) SendRequest(
	action string,
	params Dict,
	path string,
	method string,
	headers Dict,
	authpath string,
	data string,
) (interface{}, error) {
	params["action"] = action

	if _, ok := params["zone"]; !ok {
		params["zone"] = c.Zone
	}

	if c.Debug {
		bs, err := json.Marshal(params)
		if err != nil {
			return "", err
		}
		os.Stdout.Write(bs)
	}

	resp, err := c.Send(method, path, params, authpath, headers, c.Host, data)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != 200 {
		s := fmt.Sprintf("%v", resp)
		return nil, errors.New(s)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		os.Stdout.Write(bs)
	}

	var v interface{}
	if err := json.Unmarshal(bs, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func (c APIConnection) Send(
	method string,
	path string,
	params Dict,
	authpath string,
	headers Dict,
	host string,
	data string,
) (*http.Response, error) {
	if strings.Trim(host, " ") == "" {
		host = c.Host
	}

	req, err := c.BuildRequest(method, path, params, authpath, headers, host, data)
	if err != nil {
		return nil, err
	}

	resp, err := send(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type ReqArg struct {
	action   string
	params   Dict
	path     string
	method   string
	headers  Dict
	authpath string
	data     string
}

func (c APIConnection) preSendRequest(arg ReqArg) (interface{}, error) {
	if arg.params == nil {
		arg.params = make(Dict)
	}

	if arg.headers == nil {
		arg.headers = make(Dict)
	}

	return c.SendRequest(
		arg.action,
		arg.params,
		arg.path,
		arg.method,
		arg.headers,
		arg.authpath,
		arg.data,
	)
}

func (c APIConnection) DescribeZones() (interface{}, error) {
	arg := ReqArg{action: "DescribeZones"}
	return c.preSendRequest(arg)
}

func (c APIConnection) DescribeJobs(params Dict) (interface{}, error) {
	cond := Condition{
		Integers: Keys{"offset", "limit"},
		Lists:    Keys{"jobs"},
	}
	if err := cond.Check(&params); err != nil {
		return nil, err
	}

	arg := ReqArg{action: "DescribeJobs", params: params}
	return c.preSendRequest(arg)
}

func (c APIConnection) DescribeImages(params Dict) (interface{}, error) {
	cond := Condition{
		Integers: Keys{"offset", "limit", "verbose"},
		Lists:    Keys{"images"},
	}
	if err := cond.Check(&params); err != nil {
		return nil, err
	}

	arg := ReqArg{action: "DescribeImages", params: params}
	return c.preSendRequest(arg)
}
