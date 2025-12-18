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

type ViewSettings struct {
	hoursPerDay    textinput.Model
	dailySetupTime textinput.Model
}

func NewSettings() ViewSettings {
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

	return ViewSettings{
		hoursPerDay:    hoursPerDay,
		dailySetupTime: dailySetupTime,
	}
}

func (view ViewSettings) Init() tea.Cmd {
	return view.hoursPerDay.Cursor.BlinkCmd()
}

func (view ViewSettings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	view.hoursPerDay, cmd = view.hoursPerDay.Update(msg)
	cmds = append(cmds, cmd)

	view.dailySetupTime, cmd = view.dailySetupTime.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, util.Keys.Tab, util.Keys.ShiftTab):
			view.toggleFocus()
		default:
			return view, view.saveSettingsIfValid()
		}

	case messages.ViewChangedMsg:
		view.hoursPerDay.SetValue(store.GetHoursPerDay().String())
		view.dailySetupTime.SetValue(store.GetDailySetupTime().String())

		view.hoursPerDay.CursorEnd()
	}

	return view, tea.Batch(cmds...)
}

func (view ViewSettings) View() string {
	box := lipgloss.
		NewStyle().Align(lipgloss.Center).
		Padding(1, 2, 2, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffffff"))

	return box.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			store.Strings().VIEW_CAPTION_SETTINGS+"\n",
			view.hoursPerDay.View()+"\n",
			view.dailySetupTime.View()+"\n",
		),
	)
}

func (view *ViewSettings) saveSettingsIfValid() tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	if view.hoursPerDay.Err != nil {
		hoursPerDay, _ := time.ParseDuration(view.hoursPerDay.Value())
		cmds = append(cmds, store.SetHoursPerDay(hoursPerDay))
	}

	if view.dailySetupTime.Err != nil {
		dailySetupTime, _ := time.ParseDuration(view.dailySetupTime.Value())
		cmds = append(cmds, store.SetDailySetupTime(dailySetupTime))
	}

	return tea.Batch(cmds...)
}

func (view *ViewSettings) toggleFocus() {
	if view.hoursPerDay.Focused() {
		view.hoursPerDay.Blur()
		view.dailySetupTime.Focus()
	} else {
		view.dailySetupTime.Blur()
		view.hoursPerDay.Focus()
	}
}
