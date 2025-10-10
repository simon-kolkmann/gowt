package i18n

import (
	"gowt/types"
)

type Language string

const (
	LANG_DE Language = "de"
	LANG_EN Language = "en"
)

type Strings struct {
	START                    string
	END                      string
	DURATION                 string
	SUM                      string
	CURRENT_TIME             string
	CLOCKED_IN               string
	CLOCKED_OUT              string
	AT_BREAK                 string
	ESTIMATED_END_OF_WORKDAY string

	VIEW_CAPTION_SETTINGS  string
	HOURS_PER_DAY_LABEL    string
	DAILY_SETUP_TIME_LABEL string

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

func GetStringsFor(lang Language) Strings {
	switch lang {
	case LANG_DE:
		return de

	case LANG_EN:
		return en

	default:
		return en
	}
}
