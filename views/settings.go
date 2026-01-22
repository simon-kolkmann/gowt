package views

import (
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
	hoursPerDay.Width = 10
	hoursPerDay.Validate = util.Validators.Time
	hoursPerDay.PlaceholderStyle = lipgloss.NewStyle().Faint(true)
	hoursPerDay.Cursor.Blink = true
	hoursPerDay.Focus()

	dailySetupTime := textinput.New()
	dailySetupTime.Placeholder = "5m"
	dailySetupTime.CharLimit = 10
	dailySetupTime.Width = 10
	dailySetupTime.Validate = util.Validators.Time
	dailySetupTime.PlaceholderStyle = lipgloss.NewStyle().Faint(true)

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

	case store.StateMutatedMsg:
		switch msg.Field {
		case store.FIELD_ACTIVE_VIEW:
			if store.State().HoursPerDay != 0 {
				view.hoursPerDay.SetValue(store.State().HoursPerDay.String())
			}

			if store.State().DailySetupTime != 0 {
				view.dailySetupTime.SetValue(store.State().DailySetupTime.String())
			}
			view.hoursPerDay.CursorEnd()
			view.hoursPerDay.Prompt = store.State().Strings().HOURS_PER_DAY_LABEL + ": "
			view.dailySetupTime.Prompt = store.State().Strings().DAILY_SETUP_TIME_LABEL + ": "
		case store.FIELD_LANGUAGE:
			view.hoursPerDay.Prompt = store.State().Strings().HOURS_PER_DAY_LABEL + ": "
			view.dailySetupTime.Prompt = store.State().Strings().DAILY_SETUP_TIME_LABEL + ": "
		}
	}

	return view, tea.Batch(cmds...)
}

func (view ViewSettings) View() string {
	box := lipgloss.
		NewStyle().Align(lipgloss.Center).
		Padding(0, 2, 0, 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#ffffff"))

	caption := lipgloss.NewStyle().Bold(true).Underline(true).PaddingBottom(1)

	return box.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			caption.Render(store.State().Strings().VIEW_SETTINGS),
			view.hoursPerDay.View(),
			view.dailySetupTime.View(),
		),
	)
}

func (view *ViewSettings) saveSettingsIfValid() tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	if view.hoursPerDay.Err != nil {
		hoursPerDay, _ := time.ParseDuration(view.hoursPerDay.Value())
		cmds = append(cmds, store.Commit(store.SetHoursPerDay(hoursPerDay)))
	}

	if view.dailySetupTime.Err != nil {
		dailySetupTime, _ := time.ParseDuration(view.dailySetupTime.Value())
		cmds = append(cmds, store.Commit(store.SetDailySetupTime(dailySetupTime)))
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
