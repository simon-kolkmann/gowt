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

var s store = store{}

type store struct {
	activeView  types.View
	activeEntry *types.Entry
	date        time.Time
	hoursPerDay time.Duration
	entries     []types.Entry
	language    i18n.Language
}

type storeJsonFile struct {
	Date        time.Time     `json:"date"`
	HoursPerDay time.Duration `json:"hoursPerDay"`
	Entries     []types.Entry `json:"entries"`
	Language    i18n.Language `json:"language"`
}

type StoreChangedMsg struct{}

func Init() tea.Cmd {
	s.activeView = types.ViewClock

	loadFromFileOrUseDefaults()

	stateIsFromToday := s.date.Format(time.DateOnly) == time.Now().Format(time.DateOnly)

	if !stateIsFromToday {
		SetEntries(make([]types.Entry, 0))
		s.date = time.Now()
	}

	// FIXME: I use this so that the active entry is being set right after app launch,
	// but it obviously sucks.
	SetEntries(s.entries)

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

	if len(entries) > 0 {
		SetActiveEntry(&entries[len(entries)-1])
	} else {
		SetActiveEntry(nil)
	}

	return saveAndSendStoreChangedMsg
}

func GetEntries() []types.Entry {
	return s.entries
}

func AddEntry(entry types.Entry) tea.Cmd {
	s.entries = append(s.entries, entry)
	SetActiveEntry(&entry)
	return saveAndSendStoreChangedMsg
}

func SetHoursPerDay(hoursPerDay time.Duration) tea.Cmd {
	s.hoursPerDay = hoursPerDay
	return saveAndSendStoreChangedMsg
}

func GetHoursPerDay() time.Duration {
	return s.hoursPerDay
}

func SetLanguage(l i18n.Language) tea.Cmd {
	s.language = l
	return saveAndSendStoreChangedMsg
}

func GetLanguage() i18n.Language {
	return s.language
}

func SetActiveView(v types.View) tea.Cmd {
	s.activeView = v
	return saveAndSendStoreChangedMsg
}

func GetActiveView() types.View {
	return s.activeView
}

func SetActiveEntry(v *types.Entry) tea.Cmd {
	s.activeEntry = v
	return saveAndSendStoreChangedMsg
}

func GetActiveEntry() *types.Entry {
	return s.activeEntry
}

func Strings() i18n.Strings {
	switch s.language {

	case i18n.LANG_GERMAN:
		return i18n.German

	case i18n.LANG_ENGLISH:
		return i18n.English

	default:
		return i18n.English

	}
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
		s.language = i18n.LANG_ENGLISH
	} else {
		loadFromJson(file, &s)
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

func loadFromJson(f []byte, s *store) {
	sj := storeJsonFile{}
	json.Unmarshal(f, &sj)

	s.date = sj.Date
	s.hoursPerDay = sj.HoursPerDay
	s.entries = sj.Entries
	s.language = sj.Language
}
