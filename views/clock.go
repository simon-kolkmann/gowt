package views

import (
	last_clock_in "gowt/bubbles/last-clock-in"
	"gowt/bubbles/table"
	"gowt/messages"
	"gowt/store"
	"gowt/types"
	"gowt/util"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type Clock struct {
	now         string
	progress    progress.Model
	table       tea.Model
	lastClockIn tea.Model
}

func NewClock() Clock {
	return Clock{
		progress: progress.New(
			progress.WithSolidFill(types.Theme.Success),
			progress.WithWidth(50),
			progress.WithoutPercentage(),
		),
		table:       table.NewTable(),
		lastClockIn: last_clock_in.NewLastClockIn(),
	}
}

func clockIn(entry types.Entry) tea.Cmd {
	return func() tea.Msg {
		return messages.ClockInMsg{
			Entry: entry,
		}
	}
}

func clockOut() tea.Msg {
	return messages.ClockOutMsg{}
}

func (c Clock) Init() tea.Cmd {
	return nil
}

func (c Clock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, util.Keys.Enter):
			if store.LastClockIn().IsZero() {
				cmds = append(cmds, clockIn(types.Entry{
					Id:    uuid.NewString(),
					Start: time.Now(),
				}))
			} else {
				cmds = append(cmds, clockOut)
			}
		}

	case util.TimeTickMsg:
		c.now = string(msg)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := c.progress.Update(msg)
		c.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)

	case messages.ClockInMsg:
		c.progress.FullColor = types.Theme.Success
		cmds = append(cmds, store.AddEntry(msg.Entry))

	case messages.ClockOutMsg:
		c.progress.FullColor = types.Theme.Error
		entries := store.GetEntries()
		entries[len(entries)-1].End = time.Now()
		cmds = append(cmds, store.SetEntries(entries))
	}

	c.table, cmd = c.table.Update(msg)
	cmds = append(cmds, cmd)

	lastClockIn, cmd := c.lastClockIn.Update(msg)
	c.lastClockIn = lastClockIn
	cmds = append(cmds, cmd)

	// Return the updated model to the Bubble Tea runtime for processing.
	return c, tea.Batch(cmds...)
}

func (c Clock) View() string {
	row := lipgloss.NewStyle().Margin(0, 0, 1, 0).Render

	elapsed, percent := c.getElapsedTime()

	components := []string{}
	components = append(components,
		row(strings.Replace(store.Strings().CURRENT_TIME, "$time", c.now, 1)),
		row(c.lastClockIn.View()),
		row(c.progress.ViewAs(percent/100)),
		row(elapsed+" / "+store.GetHoursPerDay().String()+" ("+c.getRemainingTime()+", "+strconv.FormatFloat(percent, 'f', 2, 64)+"%)"),
	)

	if len(store.GetEntries()) > 0 {
		components = append(components, row(c.table.View()))
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		components...,
	)
}

func (c Clock) getElapsedTime() (string, float64) {
	var elapsed time.Duration

	for _, entry := range store.GetEntries() {
		elapsed += entry.Duration()
	}

	percent := elapsed.Seconds() / (store.GetHoursPerDay().Seconds() / 100)

	return elapsed.String(), percent
}

func (c Clock) getRemainingTime() string {
	var elapsed time.Duration

	for _, entry := range store.GetEntries() {
		elapsed += entry.Duration()
	}

	remaining := (store.GetHoursPerDay() - elapsed) * -1

	if remaining < 0 {
		return remaining.String()
	} else {
		return "+" + remaining.String()
	}
}
