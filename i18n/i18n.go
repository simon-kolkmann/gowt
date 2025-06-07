package i18n

import (
	"gowt/types"
)

type Language string

const (
	LANG_GERMAN  Language = "ger"
	LANG_ENGLISH Language = "eng"
)

type Strings struct {
	START                    string
	END                      string
	DURATION                 string
	SUM                      string
	CURRENT_TIME             string
	CLOCKED_IN               string
	CLOCKED_OUT              string
	ESTIMATED_END_OF_WORKDAY string

	VIEW_CAPTION_SETTINGS string
	HOURS_PER_DAY_LABEL   string

	EDIT_ENTRY         string
	ENTRY_SAVE_SUCCESS string
	ENTRY_SAVE_FAILED  string
	NO_ENTRY_SELECTED  string

	HELP_CLOCK_IN_OUT           string
	HELP_QUIT                   string
	HELP_QUIT_KEY               string
	HELP_MOVE_UP                string
	HELP_MOVE_DOWN              string
	HELP_NEXT_VIEW_KEY          string
	HELP_PREV_VIEW_KEY          string
	HELP_VIEW_NAME              func(v types.View) string
	HELP_CHANGE_LANG            string
	HELP_CHANGE_LANG_KEY        string
	HELP_DELETE_ENTRY           string
	HELP_DELETE_ENTRY_KEY       string
	HELP_DELETE_ALL_ENTRIES     string
	HELP_DELETE_ALL_ENTRIES_KEY string
	HELP_SUBMIT                 string
	HELP_SUBMIT_KEY             string
	HELP_RESET                  string
	HELP_RESET_KEY              string
}

var German Strings = Strings{
	START:                    "Beginn",
	END:                      "Ende",
	DURATION:                 "Dauer",
	SUM:                      "Saldo",
	CURRENT_TIME:             "Es ist $time Uhr.",
	CLOCKED_IN:               "Eingestempelt seit $time Uhr.",
	CLOCKED_OUT:              "Derzeit nicht eingestempelt.",
	ESTIMATED_END_OF_WORKDAY: "Voraussichtlicher Feierabend",

	VIEW_CAPTION_SETTINGS: "Einstellungen",
	HOURS_PER_DAY_LABEL:   "tägliche Arbeitszeit",

	EDIT_ENTRY:         "Eintrag bearbeiten",
	ENTRY_SAVE_SUCCESS: "Die Eingaben wurden gespeichert.",
	ENTRY_SAVE_FAILED:  "Mindestens eine Eingabe ist fehlerhaft und kann nicht gespeichert werden.",
	NO_ENTRY_SELECTED:  "Kein Eintrag ausgewählt.",

	HELP_CLOCK_IN_OUT:  "ein- und ausstempeln",
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
}

var English Strings = Strings{
	START:                    "Start",
	END:                      "End",
	DURATION:                 "Duration",
	SUM:                      "Sum",
	CURRENT_TIME:             "It is $time.",
	CLOCKED_IN:               "Clocked in since $time.",
	CLOCKED_OUT:              "Currently not clocked in.",
	ESTIMATED_END_OF_WORKDAY: "Estimated end of workday",

	VIEW_CAPTION_SETTINGS: "Settings",
	HOURS_PER_DAY_LABEL:   "Daily work time",

	EDIT_ENTRY:         "Edit entry",
	ENTRY_SAVE_SUCCESS: "Entry saved.",
	ENTRY_SAVE_FAILED:  "At least one value is invalid and cannot be saved.",
	NO_ENTRY_SELECTED:  "No entry selected.",

	HELP_CLOCK_IN_OUT:  "clock in/out",
	HELP_QUIT:          "quit",
	HELP_QUIT_KEY:      "q/ctrl+c",
	HELP_MOVE_UP:       "move up",
	HELP_MOVE_DOWN:     "move down",
	HELP_NEXT_VIEW_KEY: "ctrl+right",
	HELP_PREV_VIEW_KEY: "ctrl+left",
	HELP_VIEW_NAME: func(v types.View) string {
		switch v {
		case types.ViewClock:
			return "view: clock"

		case types.ViewSettings:
			return "view: settings"

		case types.ViewEdit:
			return "view: edit"

		default:
			return "view: n/a"

		}
	},
	HELP_CHANGE_LANG:            "change language",
	HELP_CHANGE_LANG_KEY:        "ctrl+l",
	HELP_DELETE_ENTRY:           "delete entry",
	HELP_DELETE_ENTRY_KEY:       "del",
	HELP_DELETE_ALL_ENTRIES:     "delete all entries",
	HELP_DELETE_ALL_ENTRIES_KEY: "alt+del",
	HELP_SUBMIT:                 "submit",
	HELP_SUBMIT_KEY:             "enter",
	HELP_RESET:                  "reset",
	HELP_RESET_KEY:              "ctrl+r",
}
