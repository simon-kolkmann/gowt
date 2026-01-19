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
