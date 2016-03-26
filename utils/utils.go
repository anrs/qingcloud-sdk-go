package utils

import (
	"time"
)

func UTCTimestamp() string {
	t := time.Now().UTC()
	return t.Format(time.RFC3339)
}
