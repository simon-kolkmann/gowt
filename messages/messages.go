package messages

import (
	"gowt/types"
)

type ClockInMsg struct {
	Entry types.Entry
}

type ClockOutMsg struct {
	Entry types.Entry
}

type StartBreakMsg struct {
	Entry types.Entry
}

type EndBreakMsg struct {
	Entry types.Entry
}

type ViewChangedMsg types.View

type StoreMutatedMsg struct {
	Mutation types.Mutation
}
