package conn

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func BuildRawQuery(params Dict) (string, error) {
	v := url.Values{}

	for k, p := range params {
		p, ok := p.(string)
		if !ok {
			return "", errors.New(fmt.Sprintf("%s value is not a string", k))
		}
		v.Add(k, p)
	}

	q := v.Encode()
	return strings.Replace(q, "+", "%20", -1), nil
}

type Authable interface {
	Authorize(params *Dict, authpath string, headers *Dict, method string) error
}

type Auth struct {
	AccessKeyID     string
	SecretAccessKey string
}

func (a *Auth) Sign(raw string) (string, error) {
	hm := hmac.New(sha256.New, []byte(a.SecretAccessKey))
	n, err := hm.Write([]byte(raw))
	if err != nil {
		return "", err
	}
	if n != len(raw) {
		return "", errors.New("sign failed")
	}

	encoded := base64.StdEncoding.EncodeToString(hm.Sum(nil))
	return encoded, nil
}
