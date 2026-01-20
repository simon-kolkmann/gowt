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

	case store.StateMutatedMsg:
		switch msg.Field {
		case store.FIELD_ACTIVE_VIEW, store.FIELD_ENTRIES:
			m.lastClockIn = store.State().GetLastClockIn()
		}
	}

	return m, nil
}

func (m Model) View() string {
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000")).Padding(1, 4).Width(50).Align(lipgloss.Center)

	if store.State().IsClockedIn() {
		template := store.State().Strings().CLOCKED_IN
		lastClockIn := m.lastClockIn.Format(time.TimeOnly)
		s := strings.Replace(template, "$time", lastClockIn, 1)
		return style.Background(lipgloss.Color(types.Theme.Success)).Render(s)
	} else if store.State().IsAtBreak() {
		s := store.State().Strings().AT_BREAK
		return style.Background(lipgloss.Color(types.Theme.Warn)).Render(s)
	} else {
		s := store.State().Strings().CLOCKED_OUT
		return style.Background(lipgloss.Color(types.Theme.Error)).Render(s)
	}
}
