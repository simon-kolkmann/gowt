package table

import (
	"gowt/messages"
	"gowt/store"
	"gowt/types"
	"gowt/util"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	table   table.Model
	cursor  int
	entries []types.Entry
}

func NewTable() Model {
	model := Model{
		table: createTable(),
	}

	return model
}

func createTable() table.Model {
	columns := []table.Column{
		{Title: store.State().Strings().KIND, Width: 10},
		{Title: store.State().Strings().START, Width: 10},
		{Title: store.State().Strings().END, Width: 10},
		{Title: store.State().Strings().DURATION, Width: 10},
		{Title: store.State().Strings().SUM, Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(types.Theme.Primary)).
		Foreground(lipgloss.Color(types.Theme.Text)).
		BorderBottom(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(types.Theme.Text)).
		Background(lipgloss.Color(types.Theme.Primary)).
		Bold(false)
	t.SetStyles(s)

	return t
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd = make([]tea.Cmd, 0)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, util.Keys.Delete):
			entries := make([]types.Entry, 0)
			cursor := m.table.Cursor()
			entryIndex := len(store.State().Entries) - 1 - cursor

			for i, entry := range store.State().Entries {
				if i != entryIndex {
					entries = append(entries, entry)
				}
			}
			m.table.SetCursor(cursor - 1)
			cmds = append(cmds, store.Commit(store.SetEntries(entries)))

		case key.Matches(msg, util.Keys.AltDelete):
			cmds = append(cmds, store.Commit(store.SetEntries(make([]types.Entry, 0))))

		case key.Matches(msg, util.Keys.Up, util.Keys.Down):
			m.cursor = m.table.Cursor()
			cmds = append(cmds, store.Commit(store.SetActiveEntry(m.getSelectedEntry())))
		}

	case util.TimeTickMsg, messages.ClockInMsg, messages.ClockOutMsg:
		m.calculateTableRows()

	case store.StateMutatedMsg:
		switch msg.Field {
		case store.FIELD_ENTRIES, store.FIELD_LANGUAGE, store.FIELD_ACTIVE_VIEW:
			m.table = createTable()
			m.entries = store.State().Entries
			m.calculateTableRows()
			m.table.SetCursor(m.cursor)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func (m *Model) calculateTableRows() {
	rows := make([]table.Row, 0)

	totalWorkTime := store.State().GetElapsedWorkTime()

	for i := len(m.entries) - 1; i >= 0; i-- {
		entry := m.entries[i]

		var kind string
		var start string
		var end string
		var duration string
		var sum string

		kind = store.State().Strings().ENTRY_KIND(entry.Kind)
		start, end, duration = entry.ToString()
		sum = totalWorkTime.String()

		if entry.Kind == types.EntryKindBreak {
			sum = "-"
		}

		rows = append(rows, table.Row{
			kind,
			start,
			end,
			duration,
			sum,
		})

		if entry.Kind == types.EntryKindWork {
			totalWorkTime = totalWorkTime + entry.Duration()*-1
		}

	}

	m.table.SetRows(rows)

	const MAX_ROWS = 10

	if len(rows) > MAX_ROWS {
		m.table.SetHeight(MAX_ROWS + 2)
	} else {
		m.table.SetHeight(len(rows) + 2)
	}
}

func (m *Model) getSelectedEntry() *types.Entry {
	cursor := m.table.Cursor()
	return &m.entries[len(m.entries)-cursor-1]
}
