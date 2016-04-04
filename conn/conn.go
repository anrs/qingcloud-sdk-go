package conn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/nu7hatch/gouuid"
)

type Dict map[string]interface{}

type List []interface{}

type HTTPConnector interface {
	NextRequestID() string
	
	BuildRequest(
		method string,
		path string,
		params Dict,
		authpath string,
		headers Dict,
		host string,
		data string,
	) (*http.Request, error)

	Send(
		method string,
		path string,
		params Dict,
		authpath string,
		headers Dict,
		host string,
		data string,
	) (*http.Response, error)
}

type HTTPConnection struct {	
	AccessKeyID     string
	SecretAccessKey string
	Host            string
	Path            string
	Port            int
	Protocol        string
	Expires         int
	Timeout         int
	Debug           bool
	AuthHandler     Authable
}

func (c HTTPConnection) NextRequestID() string {
	u, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return strings.Replace(u.String(), "-", "", -1)
}

func (c HTTPConnection) wrapInnerMap(prefix string, dict interface{}, params *Dict) error {
	for k, v := range dict.(Dict){
		newkey := fmt.Sprintf("%s.%s", prefix, k)

		switch v := v.(type) {
		case string:
			(*params)[newkey] = v

		default:
			bs, err := json.Marshal(v)
			if err != nil {
				return err
			}
			(*params)[newkey] = string(bs)
		}
	}

	return nil
}

func (c HTTPConnection) WrapParams(p Dict) (Dict, error) {
	params := make(Dict)

	for key, val := range p {
		if val == nil {
			continue
		}

		switch val := val.(type) {
		case string:
			params[key] = fmt.Sprintf("%s", val)

		case int:
			params[key] = strconv.Itoa(val)
			
		case []string:
			for i, v := range val {
				params[fmt.Sprintf("%s.%d", key, i+1)] = v
			}
			
		case []interface{}:
			for i, v := range val {
				if reflect.ValueOf(v).Kind() != reflect.Map {
					params[fmt.Sprintf("%s.%d", key, i+1)] = v.(string)
					continue
				}

				prefix := fmt.Sprintf("%s.%d", key, i+1)
				if err := c.wrapInnerMap(prefix, v, &params); err != nil {
					return nil, err
				}
			}
		}
	}

	if _, ok := params["req_id"]; !ok {
		req_id := c.NextRequestID()
		if req_id == "" {
			return params, nil
		}
		params["req_id"] = req_id
	}

	return params, nil
}

func send(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
