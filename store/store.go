package store

import (
	"encoding/json"
	"gowt/i18n"
	"gowt/messages"
	"gowt/types"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var s store = store{}

type store struct {
	activeView     types.View
	activeEntry    *types.Entry
	date           time.Time
	hoursPerDay    time.Duration
	dailySetupTime time.Duration
	entries        []types.Entry
	language       i18n.Language
}

type storeJsonFile struct {
	Date           time.Time     `json:"date"`
	HoursPerDay    time.Duration `json:"hoursPerDay"`
	DailySetupTime time.Duration `json:"dailySetupTime"`
	Entries        []types.Entry `json:"entries"`
	Language       i18n.Language `json:"language"`
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
	entries := GetEntries()

	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]

		if entry.Kind == types.EntryKindBreak {
			continue
		}

		if entry.End.IsZero() {
			// currently clocked in
			return entry.Start
		} else {
			// currently clocked out
			return time.Time{}
		}
	}

	// no entries
	return time.Time{}
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

func LastEntry() *types.Entry {
	if len(s.entries) > 0 {
		return &s.entries[len(s.entries)-1]
	} else {
		return nil
	}
}

func AddEntry(entry types.Entry) tea.Cmd {
	s.entries = append(s.entries, entry)
	SetActiveEntry(&s.entries[len(s.entries)-1])
	return saveAndSendStoreChangedMsg
}

func SetHoursPerDay(hoursPerDay time.Duration) tea.Cmd {
	s.hoursPerDay = hoursPerDay
	return saveAndSendStoreChangedMsg
}

func GetHoursPerDay() time.Duration {
	return s.hoursPerDay
}

func GetTotalBreakTime() time.Duration {
	var total time.Duration

	for _, entry := range s.entries {
		if entry.Kind == types.EntryKindBreak {
			total += entry.Duration()
		}
	}

	return total
}

func SetDailySetupTime(dailySetupTime time.Duration) tea.Cmd {
	s.dailySetupTime = dailySetupTime
	return saveAndSendStoreChangedMsg
}

func GetDailySetupTime() time.Duration {
	return s.dailySetupTime
}

func ToggleLanguage() tea.Cmd {
	if s.language == i18n.LANG_EN {
		s.language = i18n.LANG_DE
	} else {
		s.language = i18n.LANG_EN
	}

	return saveAndSendStoreChangedMsg
}

func SetActiveView(v types.View) tea.Cmd {
	s.activeView = v

	saveToFile(s)

	return func() tea.Msg {
		return messages.ViewChangedMsg(v)
	}
}

func GetActiveView() types.View {
	return s.activeView
}

func SetActiveEntry(v *types.Entry) tea.Cmd {
	s.activeEntry = v
	return saveAndSendStoreChangedMsg
}

// returns a copy of the active entry
// do not update this copy - use UpdateActiveEntry instead.
func GetActiveEntry() *types.Entry {
	if s.activeEntry != nil {
		copy := *s.activeEntry
		return &copy
	} else {
		return nil
	}
}

func UpdateActiveEntry(start, end time.Time) tea.Cmd {
	s.activeEntry.Start = start
	s.activeEntry.End = end

	return saveAndSendStoreChangedMsg
}

func Strings() i18n.Strings {
	return i18n.GetStringsFor(s.language)
}

func GetElapsedWorkTime() time.Duration {
	var elapsed time.Duration

	for _, entry := range s.entries {
		if entry.Kind == types.EntryKindWork {
			elapsed += entry.Duration()
		}
	}

	return elapsed
}

func GetElapsedWorkTimeInPercent() float64 {
	return GetElapsedWorkTime().Seconds() / (s.hoursPerDay.Seconds() / 100)
}

func GetRemainingWorkTime() time.Duration {
	return time.Duration(s.hoursPerDay - GetElapsedWorkTime())
}

func IsClockedIn() bool {
	if len(s.entries) == 0 {
		return false
	}

	current := s.entries[len(s.entries)-1]

	return current.Kind == types.EntryKindWork && current.End.IsZero()
}

func IsAtBreak() bool {
	if len(s.entries) == 0 {
		return false
	}

	current := s.entries[len(s.entries)-1]

	return current.Kind == types.EntryKindBreak && current.End.IsZero()
}

func KindAsString(kind types.EntryKind) string {
	switch kind {
	case types.EntryKindWork:
		return Strings().ENTRY_KIND_WORK
	case types.EntryKindBreak:
		return Strings().ENTRY_KIND_BREAK
	default:
		return "n/a"
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
		s.dailySetupTime = time.Duration(0)
		s.entries = make([]types.Entry, 0)
		s.language = i18n.LANG_EN
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
		Date:           s.date,
		HoursPerDay:    s.hoursPerDay,
		DailySetupTime: s.dailySetupTime,
		Entries:        s.entries,
		Language:       s.language,
	}
}

func loadFromJson(f []byte, s *store) {
	sj := storeJsonFile{}
	json.Unmarshal(f, &sj)

	s.date = sj.Date
	s.hoursPerDay = sj.HoursPerDay
	s.dailySetupTime = sj.DailySetupTime
	s.entries = sj.Entries
	s.language = sj.Language
}
