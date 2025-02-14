package types

import "time"

type Store struct {
	HoursPerWeek time.Duration
	Entries      []Entry
}
