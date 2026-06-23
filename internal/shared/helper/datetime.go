package helper

import (
	"math"
	"time"
)

func GetRemainingSeconds(remaining time.Duration) int {
	return int(math.Ceil(math.Abs(remaining.Seconds())))
}
