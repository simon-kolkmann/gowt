package help

import (
	"gowt/messages"
	"gowt/store"
	"gowt/types"
	"gowt/util"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	activeView types.View
	help       help.Model
}

func NewHelp() Model {
	initialView := types.ViewClock

	return Model{
		activeView: initialView,
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

	case messages.ViewChangedMsg:
		m.activeView = types.View(msg)
	}

	return m, nil
}

func (m *Model) View() string {
	return m.help.FullHelpView(util.Keys.FullHelp(m.activeView, store.Strings()))
}
