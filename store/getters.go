package store

import (
	"gowt/i18n"
	"gowt/types"
	"time"
)

func (state state) GetActiveEntryIndex() int {
	if state.ActiveEntry == nil {
		return -1
	}

	for i, entry := range state.Entries {
		if entry.Id == state.ActiveEntry.Id {
			return i
		}
	}

	return -1
}

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
