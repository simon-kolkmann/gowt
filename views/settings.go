package views

import (
	"gowt/messages"
	"gowt/store"
	"gowt/util"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Settings struct {
	hoursPerDay    textinput.Model
	dailySetupTime textinput.Model
}

func NewSettings() Settings {
	hoursPerDay := textinput.New()
	hoursPerDay.Placeholder = "1h23m4s"
	hoursPerDay.CharLimit = 10
	hoursPerDay.Prompt = store.Strings().HOURS_PER_DAY_LABEL + ":\n"
	hoursPerDay.Validate = util.Validators.Time
	hoursPerDay.Cursor.Blink = true
	hoursPerDay.Focus()

	dailySetupTime := textinput.New()
	dailySetupTime.Placeholder = "10m"
	dailySetupTime.CharLimit = 10
	dailySetupTime.Prompt = store.Strings().DAILY_SETUP_TIME_LABEL + ":\n"
	dailySetupTime.Validate = util.Validators.Time

	return Settings{
		hoursPerDay:    hoursPerDay,
		dailySetupTime: dailySetupTime,
	}
}

func (s Settings) Init() tea.Cmd {
	return s.hoursPerDay.Cursor.BlinkCmd()
}

func (s Settings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	s.hoursPerDay, cmd = s.hoursPerDay.Update(msg)
	cmds = append(cmds, cmd)

	s.dailySetupTime, cmd = s.dailySetupTime.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, util.Keys.Tab, util.Keys.ShiftTab):
			s.toggleFocus()
		default:
			return s, s.saveSettingsIfValid()
		}

	case messages.ViewChangedMsg:
		s.hoursPerDay.SetValue(store.GetHoursPerDay().String())
		s.dailySetupTime.SetValue(store.GetDailySetupTime().String())

		s.hoursPerDay.CursorEnd()
	}

	return s, tea.Batch(cmds...)
}

func (s Settings) View() string {
	box := lipgloss.
		NewStyle().Align(lipgloss.Center).
		Padding(1, 2, 2, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffffff"))

	return box.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			store.Strings().VIEW_CAPTION_SETTINGS+"\n",
			s.hoursPerDay.View()+"\n",
			s.dailySetupTime.View()+"\n",
		),
	)
}

func (s *Settings) saveSettingsIfValid() tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	if s.hoursPerDay.Err != nil {
		hoursPerDay, _ := time.ParseDuration(s.hoursPerDay.Value())
		cmds = append(cmds, store.SetHoursPerDay(hoursPerDay))
	}

	if s.dailySetupTime.Err != nil {
		dailySetupTime, _ := time.ParseDuration(s.dailySetupTime.Value())
		cmds = append(cmds, store.SetDailySetupTime(dailySetupTime))
	}

	return tea.Batch(cmds...)
}

func (s *Settings) toggleFocus() {
	if s.hoursPerDay.Focused() {
		s.hoursPerDay.Blur()
		s.dailySetupTime.Focus()
	} else {
		s.dailySetupTime.Blur()
		s.hoursPerDay.Focus()
	}
}
