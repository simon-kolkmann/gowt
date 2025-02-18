package store

import (
	"encoding/json"
	"gowt/types"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var s store = store{}

type store struct {
	date        time.Time
	hoursPerDay time.Duration
	entries     []types.Entry
	language    types.Language
}

type storeJsonFile struct {
	Date        time.Time      `json:"date"`
	HoursPerDay time.Duration  `json:"hoursPerDay"`
	Entries     []types.Entry  `json:"entries"`
	Language    types.Language `json:"language"`
}

type StoreChangedMsg struct{}

func Init() tea.Cmd {
	loadFromFileOrUseDefaults()

	stateIsFromToday := s.date.Format(time.DateOnly) == time.Now().Format(time.DateOnly)

	if !stateIsFromToday {
		SetEntries(make([]types.Entry, 0))
		s.date = time.Now()
	}

	return saveAndSendStoreChangedMsg
}

// Returns the last time the user clocked in.
//
// If the user is currently clocked out, a zeroed
// time will be returned.
func LastClockIn() time.Time {
	hasNoEntries := len(GetEntries()) == 0

	if hasNoEntries {
		return time.Time{}
	}

	mostRecentEntry := GetEntries()[len(GetEntries())-1]

	if mostRecentEntry.End.IsZero() {
		// currently clocked in
		return mostRecentEntry.Start
	} else {
		// currently clocked out
		return time.Time{}
	}
}

func SetEntries(entries []types.Entry) tea.Cmd {
	s.entries = entries
	return saveAndSendStoreChangedMsg
}

func GetEntries() []types.Entry {
	return s.entries
}

func AddEntry(entry types.Entry) tea.Cmd {
	s.entries = append(s.entries, entry)
	return saveAndSendStoreChangedMsg
}

func SetHoursPerDay(hoursPerDay time.Duration) tea.Cmd {
	s.hoursPerDay = hoursPerDay
	return saveAndSendStoreChangedMsg
}

func GetHoursPerDay() time.Duration {
	return s.hoursPerDay
}

func SetLanguage(l types.Language) tea.Cmd {
	s.language = l
	return saveAndSendStoreChangedMsg
}

func GetLanguage() types.Language {
	return s.language
}

func saveAndSendStoreChangedMsg() tea.Msg {
	saveToFile(s)

	return StoreChangedMsg{}
}

func getFilePath() string {
	value, _ := os.UserConfigDir()
	path := filepath.Join(value, "gowt")
	_ = os.Mkdir(path, 0700)

	return filepath.Join(path, "state.json")
}

func loadFromFileOrUseDefaults() {
	file, err := os.ReadFile(getFilePath())

	if err != nil {
		s.date = time.Now()
		s.hoursPerDay = time.Duration(time.Hour * 8)
		s.entries = make([]types.Entry, 0)
		s.language = types.LANG_ENGLISH
	} else {
		s = jsonToStore(file)
	}
}

func saveToFile(s store) {
	b, _ := json.Marshal(storeToJson(s))
	os.WriteFile(getFilePath(), b, 0700)
}

func storeToJson(s store) storeJsonFile {
	return storeJsonFile{
		Date:        s.date,
		HoursPerDay: s.hoursPerDay,
		Entries:     s.entries,
		Language:    s.language,
	}
}

func jsonToStore(f []byte) store {
	sj := storeJsonFile{}
	json.Unmarshal(f, &sj)

	store := store{}
	store.date = sj.Date
	store.hoursPerDay = sj.HoursPerDay
	store.entries = sj.Entries
	store.language = sj.Language

	return store
}
