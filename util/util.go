package util

import "time"

func ParseTimeUnsafe(s string) time.Time {
	t, _ := time.Parse(
		time.DateTime,
		time.Now().Format(time.DateOnly)+" "+s,
	)

	return t
}
