package types

type ClockInMsg struct {
	Entry Entry
}

type ClockOutMsg struct {
	Entry Entry
}

type LanguageChangedMsg Language

type ViewChangedMsg View

type StoreChangedMsg Store
