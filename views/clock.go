package views

import (
	"gowt/bubbles/last_clock_in"
	"gowt/bubbles/table"
	"gowt/bubbles/time_progress"
	"gowt/messages"
	"gowt/store"
	"gowt/types"
	"gowt/util"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type ViewClock struct {
	now          string
	progressWork time_progress.Model
	table        tea.Model
	lastClockIn  tea.Model
}

func NewClock() ViewClock {
	return ViewClock{
		progressWork: time_progress.NewTimeProgress(),
		table:        table.NewTable(),
		lastClockIn:  last_clock_in.NewLastClockIn(),
	}
}

func clockIn(entry types.Entry) tea.Cmd {
	// if this is the first entry of the day, subtract the daily setup time
	if len(store.State().Entries) == 0 {
		entry.Start = entry.Start.Add(store.State().DailySetupTime * -1)
	}

	return func() tea.Msg {
		return messages.ClockInMsg{
			Entry: entry,
		}
	}
}

func clockOut() tea.Msg {
	return messages.ClockOutMsg{}
}

func startBreak(entry types.Entry) tea.Cmd {
	return func() tea.Msg {
		return messages.StartBreakMsg{
			Entry: entry,
		}
	}
}

func endBreak() tea.Msg {
	return messages.EndBreakMsg{}
}

func (view ViewClock) Init() tea.Cmd {
	return nil
}

func (view ViewClock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, util.Keys.Enter):
			if store.State().GetLastClockIn().IsZero() {
				cmds = append(cmds, clockIn(types.Entry{
					Id:    uuid.NewString(),
					Kind:  types.EntryKindWork,
					Start: time.Now(),
				}))
			} else {
				cmds = append(cmds, clockOut)
			}
		case key.Matches(msg, util.Keys.AltEnter):
			if !store.State().IsAtBreak() {
				cmds = append(cmds, startBreak(types.Entry{
					Id:    uuid.NewString(),
					Kind:  types.EntryKindBreak,
					Start: time.Now(),
				}))
			} else {
				cmds = append(cmds, endBreak)
			}

		}

	case util.TimeTickMsg:
		view.now = string(msg)

	case messages.ClockInMsg:
		if last, ok := store.State().GetLastEntry(); ok {
			if last.End.IsZero() {
				last.End = time.Now()
				cmds = append(cmds, store.Commit(store.ModifyEntry(last)))
			}
		}

		cmds = append(cmds, store.Commit(store.AddEntry(msg.Entry)))

	case messages.ClockOutMsg:
		entries := store.State().Entries
		entries[len(entries)-1].End = time.Now()
		cmds = append(cmds, store.Commit(store.SetEntries(entries)))

	case messages.StartBreakMsg:
		if last, ok := store.State().GetLastEntry(); ok {
			if last.End.IsZero() {
				last.End = time.Now()
				cmds = append(cmds, store.Commit(store.ModifyEntry(last)))
			}
		}

		cmds = append(cmds, store.Commit(store.AddEntry(msg.Entry)))

	case messages.EndBreakMsg:
		entries := store.State().Entries
		entries[len(entries)-1].End = time.Now()
		cmds = append(cmds, store.Commit(store.SetEntries(entries)))
	}

	view.table, cmd = view.table.Update(msg)
	cmds = append(cmds, cmd)

	view.lastClockIn, cmd = view.lastClockIn.Update(msg)
	cmds = append(cmds, cmd)

	// Return the updated model to the Bubble Tea runtime for processing.
	return view, tea.Batch(cmds...)
}

func (view ViewClock) View() string {
	row := lipgloss.NewStyle().Margin(0, 0, 1, 0).Width(50).Render

	estimatedEndOfWorkday := view.getEstimatedEndOfWorkdayAsString()

	if store.State().IsClockedIn() {
		view.progressWork.Color = types.Theme.Success
	} else if store.State().IsAtBreak() {
		view.progressWork.Color = types.Theme.Warn
	} else {
		view.progressWork.Color = types.Theme.Error
	}

	view.progressWork.Elapsed = store.State().GetElapsedWorkTime()
	view.progressWork.Target = store.State().HoursPerDay

	components := []string{}
	components = append(components,
		row(strings.Replace(store.State().Strings().CURRENT_TIME, "$time", view.now, 1)),
		row(strings.Replace(store.State().Strings().ESTIMATED_END_OF_WORKDAY, "$time", estimatedEndOfWorkday, 1)),
		row(view.lastClockIn.View()),
		row(view.progressWork.View()),
	)

	if len(store.State().Entries) > 0 {
		components = append(components, row(view.table.View()))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		components...,
	)
}

func (view ViewClock) getRemainingTimeAsString() string {
	remaining := store.State().GetRemainingWorkTime() * -1

	if remaining < 0 {
		return remaining.String()
	} else {
		return "+" + remaining.String()
	}
}

func (view ViewClock) getEstimatedEndOfWorkday() time.Time {
	remaining := store.State().GetRemainingWorkTime()
	estimatedEndOfWorkday := time.Now().Add(remaining)
	return estimatedEndOfWorkday
}

func (view ViewClock) getEstimatedEndOfWorkdayAsString() string {
	estimatedEndOfWorkday := view.getEstimatedEndOfWorkday()
	return estimatedEndOfWorkday.Format(time.TimeOnly)
}
