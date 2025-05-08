package util

import (
	"gowt/i18n"
	"gowt/types"

	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type KeyMap struct {
	Up        key.Binding
	Down      key.Binding
	CtrlLeft  key.Binding
	CtrlRight key.Binding
	Enter     key.Binding
	Quit      key.Binding
	CtrlL     key.Binding
	Delete    key.Binding
	AltDelete key.Binding
}

var Keys KeyMap = KeyMap{
	Up:        key.NewBinding(key.WithKeys("up")),
	Down:      key.NewBinding(key.WithKeys("down")),
	CtrlLeft:  key.NewBinding(key.WithKeys("ctrl+left")),
	CtrlRight: key.NewBinding(key.WithKeys("ctrl+right")),
	Enter:     key.NewBinding(key.WithKeys("enter")),
	Quit:      key.NewBinding(key.WithKeys("q"), key.WithKeys("ctrl+c")),
	CtrlL:     key.NewBinding(key.WithKeys("ctrl+l")),
	Delete:    key.NewBinding(key.WithKeys("delete")),
	AltDelete: key.NewBinding(key.WithKeys("alt+delete")),
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp(view types.View, strings i18n.Strings) [][]key.Binding {
	Keys.Up.SetHelp("↑", strings.HELP_MOVE_UP)
	Keys.Down.SetHelp("↓", strings.HELP_MOVE_DOWN)
	Keys.Quit.SetHelp(strings.HELP_QUIT_KEY, strings.HELP_QUIT)
	Keys.CtrlL.SetHelp(strings.HELP_CHANGE_LANG_KEY, strings.HELP_CHANGE_LANG)
	Keys.Delete.SetHelp(strings.HELP_DELETE_ENTRY_KEY, strings.HELP_DELETE_ENTRY)
	Keys.AltDelete.SetHelp(strings.HELP_DELETE_ALL_ENTRIES_KEY, strings.HELP_DELETE_ALL_ENTRIES)

	switch view {

	case types.ViewClock:
		Keys.Enter.SetHelp("enter", strings.HELP_CLOCK_IN_OUT)
		Keys.CtrlLeft.SetHelp(strings.HELP_PREV_VIEW_KEY, strings.HELP_VIEW_NAME(types.ViewSettings))
		Keys.CtrlRight.SetHelp(strings.HELP_NEXT_VIEW_KEY, strings.HELP_VIEW_NAME(types.ViewEdit))
		return [][]key.Binding{
			{k.CtrlLeft, k.CtrlRight, k.CtrlL, k.Quit},     // first column
			{k.Enter, k.Delete, k.AltDelete, k.Up, k.Down}, // second column
		}

	case types.ViewSettings:
		Keys.CtrlRight.SetHelp(strings.HELP_NEXT_VIEW_KEY, strings.HELP_VIEW_NAME(types.ViewClock))
		return [][]key.Binding{
			{k.CtrlRight, k.CtrlL, k.Quit}, // first column
			{},                             // second column
		}

	case types.ViewEdit:
		Keys.Enter.SetHelp(strings.HELP_SUBMIT_KEY, strings.HELP_SUBMIT)
		Keys.CtrlLeft.SetHelp(strings.HELP_PREV_VIEW_KEY, strings.HELP_VIEW_NAME(types.ViewClock))
		return [][]key.Binding{
			{k.Enter, k.CtrlLeft},
			{},
		}

	default:
		return [][]key.Binding{}

	}

}
