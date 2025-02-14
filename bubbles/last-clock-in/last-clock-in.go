package last_clock_in

import (
	"gowt/i18n"
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

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case types.ClockInMsg:
		m.lastClockIn = msg.Entry.Start

	case types.ClockOutMsg:
		m.lastClockIn = time.Time{}

	case types.StoreChangedMsg:
		m.lastClockIn = msg.Store.LastClockIn()

	}

	return m, nil
}

func (m *Model) View() string {
	style := lipgloss.NewStyle().Bold(true)

	if m.lastClockIn.IsZero() {
		s := i18n.Strings().CLOCKED_OUT
		return style.Foreground(lipgloss.Color(types.Theme.Error)).Render(s)
	}

	template := i18n.Strings().CLOCKED_IN
	lastClockIn := m.lastClockIn.Format(time.TimeOnly)
	s := strings.Replace(template, "$time", lastClockIn, 1)

	return style.Foreground(lipgloss.Color(types.Theme.Success)).Render(s)
}
