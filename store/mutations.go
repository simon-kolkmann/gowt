package store

import (
	"gowt/i18n"
	"gowt/types"
	"time"
)

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
