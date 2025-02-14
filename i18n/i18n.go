package i18n

import "gowt/types"

const (
	LANG_GERMAN  types.Language = "LANG_GERMAN"
	LANG_ENGLISH types.Language = "LANG_ENGLISH"
)

type strings struct {
	START        string
	END          string
	DURATION     string
	SUM          string
	CURRENT_TIME string
	CLOCKED_IN   string
	CLOCKED_OUT  string

	VIEW_CAPTION_SETTINGS string
	HOURS_PER_WEEK_LABEL  string

	HELP_CLOCK_IN_OUT    string
	HELP_QUIT            string
	HELP_QUIT_KEY        string
	HELP_MOVE_UP         string
	HELP_MOVE_DOWN       string
	HELP_NEXT_VIEW_KEY   string
	HELP_PREV_VIEW_KEY   string
	HELP_VIEW_NAME       func(v types.View) string
	HELP_CHANGE_LANG     string
	HELP_CHANGE_LANG_KEY string
}

var german strings = strings{
	START:        "Beginn",
	END:          "Ende",
	DURATION:     "Dauer",
	SUM:          "Saldo",
	CURRENT_TIME: "Es ist $time Uhr.",
	CLOCKED_IN:   "Eingestempelt seit $time Uhr.",
	CLOCKED_OUT:  "Derzeit nicht eingestempelt.",

	VIEW_CAPTION_SETTINGS: "Einstellungen",
	HOURS_PER_WEEK_LABEL:  "Wochenarbeitszeit",

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

		default:
			return "ansicht: n/a"

		}
	},
	HELP_CHANGE_LANG:     "sprache wechseln",
	HELP_CHANGE_LANG_KEY: "strg+l",
}

var english strings = strings{
	START:        "Start",
	END:          "End",
	DURATION:     "Duration",
	SUM:          "Sum",
	CURRENT_TIME: "It is $time.",
	CLOCKED_IN:   "Clocked in since $time.",
	CLOCKED_OUT:  "Currently not clocked in.",

	VIEW_CAPTION_SETTINGS: "Settings",
	HOURS_PER_WEEK_LABEL:  "Weekly work time",

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

		default:
			return "view: n/a"

		}
	},
	HELP_CHANGE_LANG:     "change language",
	HELP_CHANGE_LANG_KEY: "ctrl+l",
}

var Selected types.Language

func Strings() strings {
	if Selected == LANG_GERMAN {
		return german
	}

	return english
}
