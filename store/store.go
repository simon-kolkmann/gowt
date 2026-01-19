package store

import (
	"encoding/json"
	"gowt/i18n"
	"gowt/types"
	"os"
	"path/filepath"
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

// Mutations

func SetActiveView(view types.View) mutation {
	return mutation{
		Fields: []field{FIELD_ACTIVE_VIEW},
		Mutate: func(state *state) {
			state.ActiveView = view
		},
	}
}

func SetEntries(entries []types.Entry) mutation {
	return mutation{
		Fields: []field{FIELD_ENTRIES, FIELD_ACTIVE_ENTRY},
		Mutate: func(state *state) {
			state.Entries = entries

			if len(entries) > 0 {
				state.ActiveEntry = &entries[len(entries)-1]
			} else {
				state.ActiveEntry = nil
			}
		},
	}
}

func AddEntry(entry types.Entry) mutation {
	return mutation{
		Fields: []field{FIELD_ENTRIES, FIELD_ACTIVE_ENTRY},
		Mutate: func(state *state) {
			state.Entries = append(state.Entries, entry)
			state.ActiveEntry = &state.Entries[len(state.Entries)-1]
		},
	}
}

func SetDate(date time.Time) mutation {
	return mutation{
		Fields: []field{FIELD_DATE},
		Mutate: func(state *state) {
			state.Date = date
		},
	}
}

func SetHoursPerDay(hoursPerDay time.Duration) mutation {
	return mutation{
		Fields: []field{FIELD_HOURS_PER_DAY},
		Mutate: func(state *state) {
			state.HoursPerDay = hoursPerDay
		},
	}
}

func SetDailySetupTime(dailySetupTime time.Duration) mutation {
	return mutation{
		Fields: []field{FIELD_DAILY_SETUP_TIME},
		Mutate: func(state *state) {
			state.DailySetupTime = dailySetupTime
		},
	}
}

func SetLanguage(language i18n.Language) mutation {
	return mutation{
		Fields: []field{FIELD_LANGUAGE},
		Mutate: func(state *state) {
			state.Language = language
		},
	}
}

func ToggleLanguage() mutation {
	return mutation{
		Fields: []field{FIELD_LANGUAGE},
		Mutate: func(state *state) {
			if state.Language == i18n.LANG_EN {
				state.Language = i18n.LANG_DE
			} else {
				state.Language = i18n.LANG_EN
			}
		},
	}
}

func SetActiveEntry(v *types.Entry) mutation {
	return mutation{
		Fields: []field{FIELD_ACTIVE_ENTRY},
		Mutate: func(state *state) {
			internalState.ActiveEntry = v
		},
	}
}

func ModifyActiveEntry(start, end time.Time) mutation {
	return mutation{
		Fields: []field{FIELD_ACTIVE_ENTRY, FIELD_ENTRIES},
		Mutate: func(state *state) {
			state.ActiveEntry.Start = start
			state.ActiveEntry.End = end
		},
	}
}

func ModifyEntry(entry types.Entry) mutation {
	return mutation{
		Fields: []field{FIELD_ENTRIES},
		Mutate: func(state *state) {
			for i, candidate := range state.Entries {
				if candidate.Id == entry.Id {
					state.Entries[i] = entry
				}
			}
		},
	}
}

// Getters

// Returns the last time the user clocked in.
//
// If the user is currently clocked out, a zeroed
// time will be returned.
func (state state) GetLastClockIn() time.Time {
	entries := state.Entries

	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]

		if entry.Kind == types.EntryKindBreak {
			continue
		}

		if entry.End.IsZero() {
			// currently clocked in
			return entry.Start
		} else {
			// currently clocked out
			return time.Time{}
		}
	}

	// no entries
	return time.Time{}
}

func (state state) GetLastEntry() (last types.Entry, ok bool) {
	if len(state.Entries) > 0 {
		return state.Entries[len(state.Entries)-1], true
	} else {
		return types.Entry{}, false
	}
}

func (state state) GetElapsedWorkTime() time.Duration {
	var elapsed time.Duration

	for _, entry := range internalState.Entries {
		if entry.Kind == types.EntryKindWork {
			elapsed += entry.Duration()
		}
	}

	return elapsed
}

func (state state) GetRemainingWorkTime() time.Duration {
	return time.Duration(state.HoursPerDay - state.GetElapsedWorkTime())
}

func (state state) IsClockedIn() bool {
	if len(state.Entries) == 0 {
		return false
	}

	current := state.Entries[len(state.Entries)-1]

	return current.Kind == types.EntryKindWork && current.End.IsZero()
}

func (state state) IsAtBreak() bool {
	if len(state.Entries) == 0 {
		return false
	}

	current := state.Entries[len(state.Entries)-1]

	return current.Kind == types.EntryKindBreak && current.End.IsZero()
}

func (state state) Strings() i18n.Strings {
	return i18n.GetStringsFor(state.Language)
}

// Internal helpers

func (state state) getFilePath() string {
	value, _ := os.UserConfigDir()
	path := filepath.Join(value, "gowt")
	_ = os.Mkdir(path, 0700)

	return filepath.Join(path, "state.json")
}

func (state state) saveToFile() {
	data, _ := json.Marshal(state)
	os.WriteFile(state.getFilePath(), data, 0700)
}

func (s *state) loadFromFileOrUseDefaults() tea.Cmd {
	// these mutations should always happen every time the app starts.
	useFixed := func() tea.Cmd {
		return Commit(
			SetActiveView(types.ViewClock),
			SetActiveEntry(nil),
		)
	}

	// if the state can't be read from the state file, these defaults should
	// be used.
	useDefaults := func() tea.Cmd {
		return tea.Batch(
			useFixed(),
			Commit(
				SetDate(time.Now()),
				SetHoursPerDay(time.Duration(time.Hour*8)),
				SetDailySetupTime(time.Duration(0)),
				SetEntries(make([]types.Entry, 0)),
				SetLanguage(i18n.LANG_EN),
			),
		)
	}

	file, err := os.ReadFile(s.getFilePath())

	if err != nil {
		return useDefaults()
	}

	stateFromFile := state{}
	err = json.Unmarshal(file, &stateFromFile)

	if err != nil {
		return useDefaults()
	}

	return tea.Batch(
		useFixed(),
		Commit(
			SetDate(stateFromFile.Date),
			SetHoursPerDay(stateFromFile.HoursPerDay),
			SetDailySetupTime(stateFromFile.DailySetupTime),
			SetEntries(stateFromFile.Entries),
			SetLanguage(stateFromFile.Language),
		),
	)
}
