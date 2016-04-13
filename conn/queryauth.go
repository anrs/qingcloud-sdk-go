package conn

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/anrs/qingcloud-sdk-go/utils"
)

type QueryAuth struct {
	Auth
	SignatureVersion int
	APIVersion       int
}

func NewQueryAuth(access string, secret string) *QueryAuth {
	return &QueryAuth{
		Auth{access, secret}, 1, 1,
	}
}

func (a *QueryAuth) Authorize(
	params *Dict,
	authpath string,
	headers *Dict,
	method string,
) error {
	(*params)["access_key_id"] = a.AccessKeyID
	(*params)["signature_version"] = strconv.Itoa(a.SignatureVersion)
	(*params)["signature_method"] = "HmacSHA256"
	(*params)["version"] = strconv.Itoa(a.APIVersion)

	if _, ok := (*params)["time_stamp"]; !ok {
		(*params)["time_stamp"] = utils.UTCTimestamp()
	}

	raw, err := BuildRawQuery(*params)
	if err != nil {
		return err
	}

	raw = fmt.Sprintf("%s\n%s\n%s", strings.ToUpper(method), authpath, raw)
	signature, err := a.Sign(raw)
	if err != nil {
		return err
	}

	(*params)["signature"] = signature
	return nil
}
