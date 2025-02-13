package types

import "time"

type ClockInMsg struct {
	Entry Entry
}

type ClockOutMsg struct {
	Entry Entry
}

type LanguageChangedMsg Language

type ViewChangedMsg View

type TargetDurationChangedMsg time.Duration
