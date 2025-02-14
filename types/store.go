package types

import (
	"time"
)

type Store struct {
	Date        time.Time     `json:"date"`
	HoursPerDay time.Duration `json:"hoursPerDay"`
	Entries     []Entry       `json:"entries"`
}

// Returns the last time the user clocked in.
//
// If the user is currently clocked out, a zeroed
// time will be returned.
func (s *Store) LastClockIn() time.Time {
	hasNoEntries := len(s.Entries) == 0

	if hasNoEntries {
		return time.Time{}
	}

	mostRecentEntry := s.Entries[len(s.Entries)-1]

	if mostRecentEntry.End.IsZero() {
		// currently clocked in
		return mostRecentEntry.Start
	} else {
		// currently clocked out
		return time.Time{}
	}
}
