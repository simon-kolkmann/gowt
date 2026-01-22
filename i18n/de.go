package i18n

import "gowt/types"

var de Strings = Strings{
	VIEW_CLOCK:               "Uhr",
	KIND:                     "Typ",
	START:                    "Beginn",
	END:                      "Ende",
	DURATION:                 "Dauer",
	SUMMARIZED_WORK_TIME:     "Saldo (Arbeit)",
	SUMMARIZED_BREAK_TIME:    "Saldo (Pause)",
	CURRENT_TIME:             "Es ist $time Uhr.",
	CLOCKED_IN:               "Eingestempelt seit $time Uhr.",
	CLOCKED_OUT:              "Derzeit nicht eingestempelt.",
	AT_BREAK:                 "In Pause.",
	ESTIMATED_END_OF_WORKDAY: "Frühstmöglicher Feierabend: $time Uhr.",

	ENTRY_KIND: func(v types.EntryKind) string {
		switch v {
		case types.EntryKindWork:
			return "Arbeit"
		case types.EntryKindBreak:
			return "Pause"
		default:
			return "n/a"
		}
	},

	VIEW_SETTINGS:          "Einstellungen",
	HOURS_PER_DAY_LABEL:    "tägliche Arbeitszeit",
	DAILY_SETUP_TIME_LABEL: "tägliche Rüstzeit",

	VIEW_EDIT:          "Bearbeiten",
	EDIT_ENTRY:         "Eintrag bearbeiten",
	ENTRY_SAVE_SUCCESS: "Die Eingaben wurden gespeichert.",
	ENTRY_SAVE_FAILED:  "Mindestens eine Eingabe ist fehlerhaft und kann nicht gespeichert werden.",
	NO_ENTRY_SELECTED:  "Kein Eintrag ausgewählt.",

	HELP_CLOCK_IN_OUT:  "ein- und ausstempeln",
	HELP_BREAK:         "pause starten/beenden",
	HELP_QUIT:          "beenden",
	HELP_QUIT_KEY:      "q/strg+c",
	HELP_MOVE_UP:       "hoch",
	HELP_MOVE_DOWN:     "runter",
	HELP_NEXT_VIEW_KEY: "strg+rechts",
	HELP_PREV_VIEW_KEY: "strg+links",
	HELP_VIEW_NAME: func(v types.View) string {
		switch v {
		case types.ViewClock:
			return "ansicht: uhr"

		case types.ViewSettings:
			return "ansicht: einstellungen"

		case types.ViewEdit:
			return "ansicht: bearbeiten"

		default:
			return "ansicht: n/a"

		}
	},
	HELP_CHANGE_LANG:            "sprache wechseln",
	HELP_CHANGE_LANG_KEY:        "strg+l",
	HELP_DELETE_ENTRY:           "eintrag löschen",
	HELP_DELETE_ENTRY_KEY:       "entf",
	HELP_DELETE_ALL_ENTRIES:     "alle einträge löschen",
	HELP_DELETE_ALL_ENTRIES_KEY: "alt+entf",
	HELP_SUBMIT:                 "bestätigen",
	HELP_SUBMIT_KEY:             "enter",
	HELP_RESET:                  "zurücksetzen",
	HELP_RESET_KEY:              "strg+r",
	HELP_BREAK_KEY:              "alt+enter",
}
