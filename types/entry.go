package types

import (
	"time"
)

type Entry struct {
	Start time.Time
	End   time.Time
}

func (e *Entry) Duration() time.Duration {
	if e.End.IsZero() {
		return time.Since(e.Start).Round(time.Second)
	}

	return e.End.Sub(e.Start).Round(time.Second)
}
