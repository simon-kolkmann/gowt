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

type ViewChangedMsg types.View
