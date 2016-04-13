package conn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type IaaSConnection struct {
	HTTPConnection

	Zone string
}

func NewIaaSConnection(httpConn HTTPConnection, zone string) IaaSConnection {
	if strings.Trim(httpConn.Host, " ") == "" {
		httpConn.Host = "api.qingcloud.com"
	}

	if httpConn.Port == 0 {
		httpConn.Port = 443
	}

	if strings.Trim(httpConn.Protocol, " ") == "" {
		httpConn.Protocol = "https"
	}

	if httpConn.Timeout == 0 {
		httpConn.Timeout = 60
	}

	httpConn.AuthHandler = NewQueryAuth(httpConn.AccessKeyID, httpConn.SecretAccessKey)

	c := IaaSConnection{
		HTTPConnection: httpConn,
		Zone:           zone,
	}

	var connector HTTPConnector = c
	c.connector = connector
	return c
}

func (c IaaSConnection) BuildRequest(
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

	if err := c.Auth(&wrappedParams, path, &headers, method); err != nil {
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

func (c IaaSConnection) SendRequest(
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
