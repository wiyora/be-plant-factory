package helper

import (
	"math"
	"time"
)

func GetRemainingSeconds(remaining time.Duration) int {
	return int(math.Ceil(math.Abs(remaining.Seconds())))
}

func FormatDateTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
