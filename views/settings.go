package views

import (
	"gowt/i18n"
	"gowt/messages"
	"gowt/store"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Settings struct {
	hoursPerDay textinput.Model
}

func NewSettings() Settings {
	hoursPerDay := textinput.New()
	hoursPerDay.Placeholder = "1h23m4s"
	hoursPerDay.CharLimit = 10
	hoursPerDay.Focus()

	return Settings{
		hoursPerDay: hoursPerDay,
	}
}

func (s *Settings) Init() tea.Cmd {
	return s.hoursPerDay.Cursor.BlinkCmd()
}

func (s *Settings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	s.hoursPerDay, cmd = s.hoursPerDay.Update(msg)

	switch msg.(type) {
	case tea.KeyMsg:
		return s, changeTargetDurationIfValid(s.hoursPerDay.Value())

	case messages.ViewChangedMsg:
		s.hoursPerDay.SetValue(store.GetHoursPerDay().String())
		s.hoursPerDay.CursorEnd()
	}

	return s, cmd
}

func (s *Settings) View() string {
	s.hoursPerDay.Prompt = i18n.Strings().HOURS_PER_DAY_LABEL + ":\n"

	return lipgloss.JoinVertical(
		lipgloss.Left,
		i18n.Strings().VIEW_CAPTION_SETTINGS+"\n",
		s.hoursPerDay.View()+"\n",
	)
}

func changeTargetDurationIfValid(d string) tea.Cmd {
	duration, err := time.ParseDuration(d)

	if err != nil {
		return nil
	}

	return store.SetHoursPerDay(duration)
}
