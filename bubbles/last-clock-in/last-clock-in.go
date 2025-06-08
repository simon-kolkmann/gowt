package last_clock_in

import (
	"gowt/messages"
	"gowt/store"
	"gowt/types"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	lastClockIn time.Time
}

func NewLastClockIn() Model {
	return Model{
		lastClockIn: time.Time{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case messages.ClockInMsg:
		m.lastClockIn = msg.Entry.Start

	case messages.ClockOutMsg:
		m.lastClockIn = time.Time{}

	case store.StoreChangedMsg:
		m.lastClockIn = store.LastClockIn()

	}

	return m, nil
}

func (m Model) View() string {
	style := lipgloss.NewStyle().Bold(true)

	if m.lastClockIn.IsZero() {
		s := store.Strings().CLOCKED_OUT
		return style.Foreground(lipgloss.Color(types.Theme.Error)).Render(s)
	}

	if store.IsAtBreak() {
		s := store.Strings().AT_BREAK
		return style.Foreground(lipgloss.Color(types.Theme.Warn)).Render(s)
	}

	template := store.Strings().CLOCKED_IN
	lastClockIn := m.lastClockIn.Format(time.TimeOnly)
	s := strings.Replace(template, "$time", lastClockIn, 1)

	return style.Foreground(lipgloss.Color(types.Theme.Success)).Render(s)
}
