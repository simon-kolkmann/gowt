package help

import (
	"gowt/messages"
	"gowt/store"
	"gowt/types"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up             key.Binding
	Down           key.Binding
	CtrlLeft       key.Binding
	CtrlRight      key.Binding
	Enter          key.Binding
	Quit           key.Binding
	ChangeLanguage key.Binding
	Delete         key.Binding
	AltDelete      key.Binding
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CtrlLeft, k.CtrlRight, k.ChangeLanguage, k.Quit}, // first column
		{k.Enter, k.Delete, k.AltDelete, k.Up, k.Down},      // second column
	}
}

func createKeyMap(view types.View) keyMap {
	m := keyMap{
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp(store.Strings().HELP_QUIT_KEY, store.Strings().HELP_QUIT),
		),
		ChangeLanguage: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp(store.Strings().HELP_CHANGE_LANG_KEY, store.Strings().HELP_CHANGE_LANG),
		),
	}

	if view == types.ViewSettings {
		m.CtrlRight = key.NewBinding(
			key.WithKeys("ctrl+right"),
			key.WithHelp(store.Strings().HELP_NEXT_VIEW_KEY, store.Strings().HELP_VIEW_NAME(types.ViewClock)),
		)
	}

	if view == types.ViewClock {
		m.Up = key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", store.Strings().HELP_MOVE_UP),
		)

		m.Down = key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", store.Strings().HELP_MOVE_DOWN),
		)

		m.Enter = key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", store.Strings().HELP_CLOCK_IN_OUT),
		)

		m.Delete = key.NewBinding(
			key.WithKeys("delete"),
			key.WithHelp(store.Strings().HELP_DELETE_ENTRY_KEY, store.Strings().HELP_DELETE_ENTRY),
		)

		m.AltDelete = key.NewBinding(
			key.WithKeys("alt+delete"),
			key.WithHelp(store.Strings().HELP_DELETE_ALL_ENTRIES_KEY, store.Strings().HELP_DELETE_ALL_ENTRIES),
		)

		m.CtrlLeft = key.NewBinding(
			key.WithKeys("ctrl+left"),
			key.WithHelp(store.Strings().HELP_PREV_VIEW_KEY, store.Strings().HELP_VIEW_NAME(types.ViewSettings)),
		)
	}

	return m
}

type Model struct {
	activeView types.View
	keys       keyMap
	help       help.Model
}

func NewHelp() Model {
	initialView := types.ViewClock

	return Model{
		activeView: initialView,
		keys:       createKeyMap(initialView),
		help:       help.New(),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width

	case store.StoreChangedMsg:
		// language change
		// TODO: more specific messages
		m.keys = createKeyMap(m.activeView)

	case messages.ViewChangedMsg:
		m.activeView = types.View(msg)
		m.keys = createKeyMap(m.activeView)
	}

	return m, nil
}

func (m *Model) View() string {
	return m.help.FullHelpView(m.keys.FullHelp())
}
