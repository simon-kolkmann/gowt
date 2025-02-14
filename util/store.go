package util

import (
	"encoding/json"
	"gowt/types"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var Store types.Store = types.Store{}

func InitStore() {
	loadFromFileOrUseDefaults(&Store)

	stateIsFromToday := Store.Date.Format(time.DateOnly) == time.Now().Format(time.DateOnly)

	if !stateIsFromToday {
		Store.Entries = make([]types.Entry, 0)
		Store.Date = time.Now()
	}
}

func SendStoreChangedMsg() tea.Msg {
	saveToFile(&Store)
	return types.StoreChangedMsg{
		Store: Store,
	}
}

func getFilePath() string {
	value, _ := os.UserConfigDir()
	path := filepath.Join(value, "gowt")
	_ = os.Mkdir(path, 0700)

	return filepath.Join(path, "state.json")
}

func loadFromFileOrUseDefaults(s *types.Store) {
	f, err := os.ReadFile(getFilePath())

	if err != nil {
		s.Date = time.Now()
		s.HoursPerDay = time.Duration(time.Hour * 8)
		s.Entries = make([]types.Entry, 0)
	}

	json.Unmarshal(f, s)
}

func saveToFile(s *types.Store) {
	b, _ := json.Marshal(s)
	os.WriteFile(getFilePath(), b, 0700)
}
