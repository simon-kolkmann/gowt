package views

import (
	"gowt/i18n"
	"gowt/types"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Settings struct {
	hoursPerWeek textinput.Model
}

func NewSettings() Settings {
	hoursPerWeek := textinput.New()
	hoursPerWeek.Placeholder = "1h23m4s"
	hoursPerWeek.CharLimit = 10
	hoursPerWeek.Focus()

	return Settings{
		hoursPerWeek: hoursPerWeek,
	}
}

func (s *Settings) Init() tea.Cmd {
	return s.hoursPerWeek.Cursor.BlinkCmd()
}

func (s *Settings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	s.hoursPerWeek, cmd = s.hoursPerWeek.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return s, changeTargetDuration(s.hoursPerWeek.Value())
		}
	}

	return s, cmd
}

func (s *Settings) View() string {
	s.hoursPerWeek.Prompt = i18n.Strings().HOURS_PER_WEEK_LABEL + ":\n"

	return lipgloss.JoinVertical(
		lipgloss.Left,
		i18n.Strings().VIEW_CAPTION_SETTINGS+"\n",
		s.hoursPerWeek.View()+"\n",
	)
}

func changeTargetDuration(d string) tea.Cmd {
	duration, err := time.ParseDuration(d)

	if err != nil {
		return nil
	}

	return func() tea.Msg {
		return types.TargetDurationChangedMsg(duration)
	}
}
