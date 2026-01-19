package store

import (
	"gowt/i18n"
	"gowt/types"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type StateMutatedMsg struct {
	Field  field
	Before state
	After  state
}

var internalState state = state{}

type state struct {
	ActiveView     types.View    `json:"-"`
	ActiveEntry    *types.Entry  `json:"-"`
	Date           time.Time     `json:"date"`
	HoursPerDay    time.Duration `json:"hoursPerDay"`
	DailySetupTime time.Duration `json:"dailySetupTime"`
	Entries        []types.Entry `json:"entries"`
	Language       i18n.Language `json:"language"`
}

type field string

const (
	FIELD_ACTIVE_VIEW      field = "ActiveView"
	FIELD_ACTIVE_ENTRY     field = "ActiveEntry"
	FIELD_DATE             field = "Date"
	FIELD_HOURS_PER_DAY    field = "HoursPerDay"
	FIELD_DAILY_SETUP_TIME field = "DailySetupTime"
	FIELD_ENTRIES          field = "Entries"
	FIELD_LANGUAGE         field = "Language"
)

type mutation = struct {
	Fields []field
	Mutate func(*state)
}

func Initialize() tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	cmds = append(cmds, internalState.loadFromFileOrUseDefaults())

	stateIsFromToday := internalState.Date.Format(time.DateOnly) == time.Now().Format(time.DateOnly)

	if !stateIsFromToday {
		cmds = append(cmds, Commit(
			SetEntries(make([]types.Entry, 0)),
			SetDate(time.Now()),
		))
	}

	return tea.Batch(cmds...)
}

func Commit(mutations ...mutation) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	for _, mutation := range mutations {
		before := internalState
		mutation.Mutate(&internalState)

		for _, field := range mutation.Fields {
			cmds = append(cmds, func() tea.Msg {
				return StateMutatedMsg{
					Field:  field,
					Before: before,
					After:  internalState,
				}
			})
		}
	}

	internalState.saveToFile()

	return tea.Batch(cmds...)
}

func State() state {
	return internalState
}
