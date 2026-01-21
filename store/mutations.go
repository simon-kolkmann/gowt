package store

import (
	"gowt/i18n"
	"gowt/types"
	"slices"
	"time"
)

func SetActiveView(view types.View) mutation {
	return func(state *state) []field {
		state.ActiveView = view
		return []field{FIELD_ACTIVE_VIEW}
	}
}

func SetEntries(entries []types.Entry) mutation {
	return func(state *state) []field {
		state.Entries = entries

		if len(entries) > 0 {
			state.ActiveEntry = &entries[len(entries)-1]
		} else {
			state.ActiveEntry = nil
		}

		return []field{FIELD_ENTRIES, FIELD_ACTIVE_ENTRY}
	}
}

func AddEntry(entry types.Entry) mutation {
	return func(state *state) []field {
		state.Entries = append(state.Entries, entry)
		state.ActiveEntry = &state.Entries[len(state.Entries)-1]
		return []field{FIELD_ENTRIES, FIELD_ACTIVE_ENTRY}
	}
}

func SetDate(date time.Time) mutation {
	return func(state *state) []field {
		state.Date = date
		return []field{FIELD_DATE}
	}
}

func SetHoursPerDay(hoursPerDay time.Duration) mutation {
	return func(state *state) []field {
		state.HoursPerDay = hoursPerDay
		return []field{FIELD_HOURS_PER_DAY}
	}
}

func SetDailySetupTime(dailySetupTime time.Duration) mutation {
	return func(state *state) []field {
		state.DailySetupTime = dailySetupTime
		return []field{FIELD_DAILY_SETUP_TIME}
	}
}

func SetLanguage(language i18n.Language) mutation {
	return func(state *state) []field {
		state.Language = language
		return []field{FIELD_LANGUAGE}
	}
}

func ToggleLanguage() mutation {
	return func(state *state) []field {
		if state.Language == i18n.LANG_EN {
			state.Language = i18n.LANG_DE
		} else {
			state.Language = i18n.LANG_EN
		}

		return []field{FIELD_LANGUAGE}
	}
}

func SetActiveEntry(v *types.Entry) mutation {
	return func(state *state) []field {
		internalState.ActiveEntry = v
		return []field{FIELD_ACTIVE_ENTRY}
	}
}

func SetActiveEntryByIndex(index int) mutation {
	return func(state *state) []field {
		internalState.ActiveEntry = &internalState.Entries[index]
		return []field{FIELD_ACTIVE_ENTRY}
	}
}

func ModifyActiveEntry(start, end time.Time) mutation {
	return func(state *state) []field {
		state.ActiveEntry.Start = start
		state.ActiveEntry.End = end
		return []field{FIELD_ACTIVE_ENTRY}
	}
}

func DeleteActiveEntry() mutation {
	return func(state *state) []field {
		idx := state.GetActiveEntryIndex()

		if idx == -1 {
			return []field{}
		}

		state.Entries = slices.Delete(state.Entries, idx, idx+1)
		if idx > 0 {
			state.ActiveEntry = &state.Entries[idx-1]
		} else if idx == 0 && len(state.Entries) > 0 {
			state.ActiveEntry = &state.Entries[0]
		} else {
			state.ActiveEntry = nil
		}

		return []field{FIELD_ACTIVE_ENTRY, FIELD_ENTRIES}
	}
}

func ModifyEntry(entry types.Entry) mutation {
	return func(state *state) []field {
		for i, candidate := range state.Entries {
			if candidate.Id == entry.Id {
				state.Entries[i] = entry
			}
		}

		return []field{FIELD_ENTRIES}
	}
}
