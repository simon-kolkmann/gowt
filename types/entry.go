package types

import (
	"time"
)

type Entry struct {
	Id    string    `json:"id"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (e *Entry) Duration() time.Duration {
	if e.End.IsZero() {
		return time.Since(e.Start).Round(time.Second)
	}

	return e.End.Sub(e.Start).Round(time.Second)
}
