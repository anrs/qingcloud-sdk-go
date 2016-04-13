package conn

import (
	"fmt"
	"sort"
	"strings"
)

type QingStorAuth struct {
	Auth
}

func NewQingStorAuth(access string, secret string) *QingStorAuth {
	return &QingStorAuth{
		Auth{access, secret},
	}
}

func get(d Dict, key string) interface{} {
	if v, ok := d[key]; ok {
		return v
	}
	return ""
}

func (a *QingStorAuth) buildQingStorRawQuery(
	params *Dict,
	authpath string,
	headers *Dict,
	method string,
) (string, error) {
	md5 := get(*headers, "Content-MD5")
	ctype := get(*headers, "Content-Type")
	raw := fmt.Sprintf("%s\n%s\n%s", strings.ToUpper(method), md5, ctype)

	date := get(*headers, "X-QS-Date")
	if date == "" {
		if date = get(*params, "X-QS-Date"); date != "" {
			raw += "\n"
		}
	}
	if date == "" {
		date = get(*headers, "Date")
		raw += fmt.Sprintf("\n%s", date)
	}

	var keys []string
	for k, _ := range *headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// X-QS headers for signaturing.
	for _, k := range keys {
		v := (*headers)[k]
		k = strings.ToLower(k)
		if strings.HasPrefix(k, "x-qs-") {
			raw += fmt.Sprintf("\n%s:%s", k, v)
		}
	}

	raw += fmt.Sprintf("\n%s", authpath)

	if query, err := BuildRawQuery(*params); err != nil {
		return "", err
	} else {
		if query != "" {
			raw += fmt.Sprintf("?%s", query)
		}
	}

	return raw, nil
}

func (a *QingStorAuth) Authorize(
	params *Dict,
	authpath string,
	headers *Dict,
	method string,
) error {
	raw, err := a.buildQingStorRawQuery(params, authpath, headers, method)
	if err != nil {
		return err
	}

	signature, err := a.Sign(raw)
	if err != nil {
		return err
	}

	(*headers)["Authorization"] = fmt.Sprintf(
		"QS-HMAC-SHA256 %s:%s", a.AccessKeyID, signature,
	)
	return nil
}
