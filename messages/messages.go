package messages

import "gowt/types"

type ClockInMsg struct {
	Entry types.Entry
}

type ClockOutMsg struct {
	Entry types.Entry
}

type LanguageChangedMsg types.Language

type ViewChangedMsg types.View
