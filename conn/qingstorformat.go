package conn

import (
	"fmt"
	"net/url"
	"strings"
)

type StyleFormatter interface {
	BuildHost(server string, bucket string) string

	BuildAuthPath(bucket string, key string) string

	BuildPathBase(bucket string, key string) string
}

type VirtualHostStyleFormatter struct{}

func (f *VirtualHostStyleFormatter) BuildHost(server string, bucket string) string {
	if strings.Trim(bucket, " ") == "" {
		return server
	}
	return fmt.Sprintf("%s.%s", bucket, server)
}

func (f *VirtualHostStyleFormatter) BuildAuthPath(bucket string, key string) string {
	path := "/"

	bucket = strings.Trim(bucket, " ")
	if bucket != "" {
		path = fmt.Sprintf("%s%s", path, bucket)
	}

	key = strings.Trim(key, " ")
	if key != "" {
		path = fmt.Sprintf("%s/%s", path, key)
	}

	return path
}

func (f *VirtualHostStyleFormatter) BuildPathBase(bucket string, key string) string {
	path := "/"

	key = strings.Trim(key, " ")
	if key != "" {
		path = fmt.Sprintf("%s%s", path, url.QueryEscape(key))
	}

	return path
}

func NewVirtualHostStyleFormatter() *VirtualHostStyleFormatter {
	return &VirtualHostStyleFormatter{}
}
