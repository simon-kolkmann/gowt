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
	table table.Model
}

func NewTable() Model {
	model := Model{
		table: getTable(),
	}

	return model
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
			cmds = append(cmds, store.Commit(store.DeleteActiveEntry()))

		case key.Matches(msg, util.Keys.AltDelete):
			cmds = append(cmds, store.Commit(store.SetEntries(make([]types.Entry, 0))))

		case key.Matches(msg, util.Keys.Up):
			index := store.State().GetActiveEntryIndex()

			if index < len(store.State().Entries)-1 {
				cmds = append(cmds, store.Commit(
					store.SetActiveEntryByIndex(index+1),
				))
			}

		case key.Matches(msg, util.Keys.Down):
			index := store.State().GetActiveEntryIndex()

			if index > 0 {
				cmds = append(cmds, store.Commit(
					store.SetActiveEntryByIndex(index-1),
				))
			}
		}

	case util.TimeTickMsg, messages.ClockInMsg, messages.ClockOutMsg:
		m.UpdateTable()

	case store.StateMutatedMsg:
		switch msg.Field {
		case
			store.FIELD_ENTRIES,
			store.FIELD_ACTIVE_ENTRY,
			store.FIELD_ACTIVE_VIEW:
			m.UpdateTable()
		case store.FIELD_LANGUAGE:
			m.RecreateTable()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func (m *Model) RecreateTable() {
	m.table = getTable()
	m.table.SetCursor(getCursor())
}

func (m *Model) UpdateTable() {
	rows, height := getRows()

	m.table.SetRows(rows)
	m.table.SetHeight(height)
	m.table.SetCursor(getCursor())
}

func getTable() table.Model {
	t := table.New(
		table.WithColumns(getColumns()),
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

	rows, height := getRows()
	t.SetRows(rows)
	t.SetHeight(height)

	return t
}

func getColumns() []table.Column {
	return []table.Column{
		{Title: store.State().Strings().KIND, Width: 10},
		{Title: store.State().Strings().START, Width: 10},
		{Title: store.State().Strings().END, Width: 10},
		{Title: store.State().Strings().DURATION, Width: 10},
		{Title: store.State().Strings().SUM, Width: 10},
	}
}

func getRows() (rows []table.Row, height int) {
	rows = make([]table.Row, 0)

	entries := store.State().Entries
	totalWorkTime := store.State().GetElapsedWorkTime()

	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]

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

	const MAX_ROWS = 10

	height = 2
	if len(rows) > MAX_ROWS {
		height += MAX_ROWS
	} else {
		height += len(rows)
	}

	return rows, height
}

func getCursor() int {
	index := store.State().GetActiveEntryIndex()
	return len(store.State().Entries) - index - 1
}
