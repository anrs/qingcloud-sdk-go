package utils

import (
	"strings"
	"time"
)

func UTCTimestamp() string {
	t := time.Now().UTC()
	return t.Format(time.RFC3339)
}

func UTCTimeRFC1123() string {
	t := time.Now().UTC()
	return t.Format(time.RFC1123)
}

func GMTime() string {
	ts := UTCTimeRFC1123()
	return strings.Replace(ts, " UTC", " GMT", -1)
}
