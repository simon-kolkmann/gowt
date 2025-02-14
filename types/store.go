package types

import (
	"time"
)

type Store struct {
	Date         time.Time     `json:"date"`
	HoursPerWeek time.Duration `json:"hoursPerWeek"`
	Entries      []Entry       `json:"entries"`
}
