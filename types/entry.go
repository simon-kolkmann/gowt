package types

import (
	"time"
)

type EntryKind int

const (
	EntryKindWork EntryKind = iota
	EntryKindBreak
)

type Entry struct {
	Id    string    `json:"id"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Kind  EntryKind `json:"kind"`
}

func (e *Entry) Duration() time.Duration {
	if e.End.IsZero() {
		return time.Since(e.Start).Round(time.Second)
	}

	return e.End.Sub(e.Start).Round(time.Second)
}

func (e Entry) ToString() (start, end, duration string) {
	start = e.Start.Format(time.TimeOnly)

	if e.End.IsZero() {
		end = "-"
	} else {
		end = e.End.Format(time.TimeOnly)
	}

	duration = e.Duration().String()

	return start, end, duration
}
