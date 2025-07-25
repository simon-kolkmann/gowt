package table

import (
	"gowt/messages"
	"gowt/store"
	"gowt/types"
	"gowt/util"
	"time"

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
		{Title: store.Strings().START, Width: 10},
		{Title: store.Strings().END, Width: 10},
		{Title: store.Strings().DURATION, Width: 10},
		{Title: store.Strings().SUM, Width: 10},
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
			entryIndex := len(store.GetEntries()) - 1 - cursor

			for i, entry := range store.GetEntries() {
				if i != entryIndex {
					entries = append(entries, entry)
				}
			}
			m.table.SetCursor(cursor - 1)
			cmds = append(cmds, store.SetEntries(entries))

		case key.Matches(msg, util.Keys.AltDelete):
			cmds = append(cmds, store.SetEntries(make([]types.Entry, 0)))

		case key.Matches(msg, util.Keys.Up, util.Keys.Down):
			m.cursor = m.table.Cursor()
			cmds = append(cmds, store.SetActiveEntry(m.getSelectedEntry()))
		}

	case util.TimeTickMsg, messages.ClockInMsg, messages.ClockOutMsg:
		m.calculateTableRows()

	case store.StoreChangedMsg:
		// TODO: more specific messages
		// language change
		m.table = createTable()
		m.entries = store.GetEntries()
		m.calculateTableRows()
		m.table.SetCursor(m.cursor)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func (m *Model) calculateTableRows() {
	rows := make([]table.Row, 0)

	totalWorkTime := time.Time{}

	// calculate total work time (all entries)
	for _, entry := range m.entries {
		totalWorkTime = totalWorkTime.Add(entry.Duration())
	}

	for i := len(m.entries) - 1; i >= 0; i-- {
		entry := m.entries[i]

		if entry.End.IsZero() {
			rows = append(rows, table.Row{
				entry.Start.Format(time.TimeOnly),
				"-",
				entry.Duration().String(),
				totalWorkTime.Format(time.TimeOnly),
			})
		} else {
			rows = append(rows, table.Row{
				entry.Start.Format(time.TimeOnly),
				entry.End.Format(time.TimeOnly),
				entry.Duration().String(),
				totalWorkTime.Format(time.TimeOnly),
			})
		}

		totalWorkTime = totalWorkTime.Add(entry.Duration() * -1)

	}

	m.table.SetRows(rows)

	const MAX_ROWS = 5

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
