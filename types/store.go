package types

import (
	"time"
)

type Store struct {
	Date        time.Time     `json:"date"`
	HoursPerDay time.Duration `json:"hoursPerDay"`
	Entries     []Entry       `json:"entries"`
}
