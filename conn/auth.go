package conn

import (
	"encoding/base64"
	"errors"
	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"net/url"
	"strconv"
	"strings"

	"github.com/anrs/qingcloud-sdk-go/utils"
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
	Authorize(s *Dict, path string) error
}

type Auth struct {
	AccessKeyID     string
	SecretAccessKey string
}

type QuerySignatureAuth struct {
	Auth
	SignatureVersion int
	APIVersion       int
}

func NewQuerySignatureAuth(access string, secret string) *QuerySignatureAuth {
	return &QuerySignatureAuth{
		Auth{access, secret}, 1, 1,
	}
}

func (a *QuerySignatureAuth) sign(params Dict, path string) (string, error) {
	query, err := BuildRawQuery(params)
	if err != nil {
		return "", err
	}

	raw := fmt.Sprintf("GET\n%s\n%s", path, query)
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

func (a *QuerySignatureAuth) Authorize(params *Dict, path string) error {
	(*params)["access_key_id"] = a.AccessKeyID
	(*params)["signature_version"] = strconv.Itoa(a.SignatureVersion)
	(*params)["signature_method"] = "HmacSHA256"
	(*params)["version"] = strconv.Itoa(a.APIVersion)

	if _, ok := (*params)["time_stamp"]; !ok {
		(*params)["time_stamp"] = utils.UTCTimestamp()
	}

	signature, err := a.sign(*params, path)
	if err != nil {
		return err
	}

	(*params)["signature"] = signature

	return nil
}
