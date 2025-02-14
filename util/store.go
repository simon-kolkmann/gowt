package util

import (
	"gowt/types"

	tea "github.com/charmbracelet/bubbletea"
)

var Store types.Store = types.Store{}

func SendStoreChangedMsg() tea.Msg {
	return types.StoreChangedMsg(Store)
}
