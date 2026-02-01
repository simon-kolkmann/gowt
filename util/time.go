package util

import (
	"strconv"
	"time"
)

func ElapsedInPercent(elapsed, target time.Duration) (percent float64, percentAsString string) {
	percent = elapsed.Seconds() / (target.Seconds() / 100)
	percentAsString = strconv.FormatFloat(percent, 'f', 2, 64) + "%"

	return percent, percentAsString
}
