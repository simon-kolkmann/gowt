package views

import (
	"gowt/bubbles/help"
	last_clock_in "gowt/bubbles/last-clock-in"
	"gowt/bubbles/table"
	"gowt/i18n"
	"gowt/types"
	"gowt/util"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Clock struct {
	now            string
	entries        []types.Entry
	progress       progress.Model
	table          table.Model
	lastClockIn    last_clock_in.Model
	help           help.Model
	targetDuration time.Duration
}

func NewClock() Clock {
	return Clock{
		entries: []types.Entry{},
		progress: progress.New(
			progress.WithSolidFill(types.Theme.Success),
			progress.WithWidth(50),
			progress.WithoutPercentage(),
		),
		table:          table.NewTable(),
		lastClockIn:    last_clock_in.NewLastClockIn(),
		help:           help.NewHelp(),
		targetDuration: time.Duration((time.Hour * 7) + (time.Minute * 42)),
	}
}

func clockIn(entry types.Entry) tea.Cmd {
	return func() tea.Msg {
		return types.ClockInMsg{
			Entry: entry,
		}
	}
}

func clockOut() tea.Msg {
	return types.ClockOutMsg{}
}

func (c *Clock) Init() tea.Cmd {
	return nil
}

func (c *Clock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if len(c.entries) == 0 {
				cmds = append(cmds, clockIn(types.Entry{Start: time.Now()}))
			} else {
				current := c.entries[len(c.entries)-1]

				if current.End.IsZero() {
					cmds = append(cmds, clockOut)
				} else {
					cmds = append(cmds, clockIn(types.Entry{Start: time.Now()}))
				}
			}

		}

	case util.TimeTickMsg:
		c.now = string(msg)

	case types.StoreChangedMsg:
		c.targetDuration = msg.Store.HoursPerDay
		c.entries = util.Store.Entries
		c.table.SetEntries(&c.entries)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := c.progress.Update(msg)
		c.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)

	case types.ClockInMsg:
		c.progress.FullColor = types.Theme.Success
		c.entries = append(c.entries, msg.Entry)
		util.Store.Entries = c.entries
		cmds = append(cmds, util.SendStoreChangedMsg)

	case types.ClockOutMsg:
		c.progress.FullColor = types.Theme.Error
		c.entries[len(c.entries)-1].End = time.Now()
		util.Store.Entries = c.entries
		cmds = append(cmds, util.SendStoreChangedMsg)

	}

	_, cmd = c.table.Update(msg)
	cmds = append(cmds, cmd)

	_, cmd = c.lastClockIn.Update(msg)
	cmds = append(cmds, cmd)

	_, cmd = c.help.Update(msg)
	cmds = append(cmds, cmd)

	// Return the updated model to the Bubble Tea runtime for processing.
	return c, tea.Batch(cmds...)
}

func (c *Clock) View() string {
	row := lipgloss.NewStyle().Margin(0, 0, 1, 0).Render

	elapsed, percent := c.getElapsedTime()

	components := []string{}
	components = append(components,
		row(strings.Replace(i18n.Strings().CURRENT_TIME, "$time", c.now, 1)),
		row(c.lastClockIn.View()),
		row(c.progress.ViewAs(percent/100)),
		row(elapsed+" / "+c.targetDuration.String()+" ("+strconv.FormatFloat(percent, 'f', 2, 64)+"%)"),
	)

	if len(c.entries) > 0 {
		components = append(components, row(c.table.View()))
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		components...,
	)
}

func (c Clock) getElapsedTime() (string, float64) {
	var elapsed time.Duration

	for _, entry := range c.entries {
		elapsed += entry.Duration()
	}

	percent := elapsed.Seconds() / (c.targetDuration.Seconds() / 100)

	return elapsed.String(), percent
}
