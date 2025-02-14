package help

import (
	"gowt/i18n"
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
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CtrlLeft, k.CtrlRight, k.ChangeLanguage, k.Quit}, // first column
		{k.Enter, k.Up, k.Down},                             // second column
	}
}

func createKeyMap(view types.View) keyMap {
	m := keyMap{
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp(i18n.Strings().HELP_QUIT_KEY, i18n.Strings().HELP_QUIT),
		),
		ChangeLanguage: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp(i18n.Strings().HELP_CHANGE_LANG_KEY, i18n.Strings().HELP_CHANGE_LANG),
		),
	}

	if view == types.ViewSettings {
		m.CtrlRight = key.NewBinding(
			key.WithKeys("ctrl+right"),
			key.WithHelp(i18n.Strings().HELP_NEXT_VIEW_KEY, i18n.Strings().HELP_VIEW_NAME(types.ViewClock)),
		)
	}

	if view == types.ViewClock {
		m.Up = key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", i18n.Strings().HELP_MOVE_UP),
		)

		m.Down = key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", i18n.Strings().HELP_MOVE_DOWN),
		)

		m.Enter = key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", i18n.Strings().HELP_CLOCK_IN_OUT),
		)

		m.CtrlLeft = key.NewBinding(
			key.WithKeys("ctrl+left"),
			key.WithHelp(i18n.Strings().HELP_PREV_VIEW_KEY, i18n.Strings().HELP_VIEW_NAME(types.ViewSettings)),
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

	case types.LanguageChangedMsg:
		m.keys = createKeyMap(m.activeView)

	case types.ViewChangedMsg:
		m.activeView = types.View(msg)
		m.keys = createKeyMap(m.activeView)
	}

	return m, nil
}

func (m *Model) View() string {
	return m.help.FullHelpView(m.keys.FullHelp())
}
