package types

import "time"

type MandatoryBreak struct {
	After    time.Duration
	Duration time.Duration
}
