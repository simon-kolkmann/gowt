package types

import "gowt/i18n"

type ClockInMsg struct {
	Entry Entry
}

type ClockOutMsg struct {
	Entry Entry
}

type LanguageChangedMsg i18n.Language
