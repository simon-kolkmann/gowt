package i18n

import "gowt/types"

var en Strings = Strings{
	START:                    "Start",
	END:                      "End",
	DURATION:                 "Duration",
	SUM:                      "Sum",
	CURRENT_TIME:             "It is $time.",
	CLOCKED_IN:               "Clocked in since $time.",
	CLOCKED_OUT:              "Currently not clocked in.",
	AT_BREAK:                 "Having a break.",
	ESTIMATED_END_OF_WORKDAY: "Estimated end of workday",

	VIEW_CAPTION_SETTINGS:  "Settings",
	HOURS_PER_DAY_LABEL:    "Daily work time",
	DAILY_SETUP_TIME_LABEL: "Daily set-up time",

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
