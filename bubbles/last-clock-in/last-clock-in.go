package last_clock_in

import (
	"gowt/types"
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

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case types.ClockInMsg:
		m.lastClockIn = msg.Entry.Start

	case types.ClockOutMsg:
		m.lastClockIn = time.Time{}

	}

	return m, nil
}

func (m *Model) View() string {
	style := lipgloss.NewStyle().Bold(true)

	if m.lastClockIn.IsZero() {
		s := "Derzeit nicht eingestempelt."
		return style.Foreground(lipgloss.Color(types.Theme.Error)).Render(s)
	}

	s := "Eingestempelt seit " + m.lastClockIn.Format(time.TimeOnly) + " Uhr."
	return style.Foreground(lipgloss.Color(types.Theme.Success)).Render(s)
}
