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

type HTTPConnector interface {
	NextRequestID() string
	
	BuildRequest(
		method string,
		path string,
		params map[string]interface{},
		authpath string,
		headers map[string]string,
		host string,
		data string,
	) (*http.Request, error)

	Send(
		method string,
		path string,
		params map[string]interface{},
		authpath string,
		headers map[string]string,
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

func (c HTTPConnection) wrapInnerMap(prefix string, m interface{}, p map[string]string) error {
	if m, ok := m.(map[string]string); ok {
		for k, v := range m {
			p[fmt.Sprintf("%s.%s", prefix, k)] = v
		}
		return nil
	}

	// Jsonize non-string.
	if m, ok := m.(map[string]interface{}); ok {
		for vk, vv := range m {
			b, err := json.Marshal(vv)
			if err != nil {
				return err
			}
			p[fmt.Sprintf("%s.%s", prefix, vk)] = string(b)
		}
	}

	return nil
}

func (c HTTPConnection) WrapParams(p map[string]interface{}) (map[string]string, error) {
	params := make(map[string]string)

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
				c.wrapInnerMap(prefix, v, params)
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
