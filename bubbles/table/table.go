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

func (c *Model) Init() tea.Cmd {
	return nil
}

func (c *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd = make([]tea.Cmd, 0)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, util.Keys.Delete):
			entries := make([]types.Entry, 0)
			cursor := c.table.Cursor()
			entryIndex := len(store.GetEntries()) - 1 - cursor

			for i, entry := range store.GetEntries() {
				if i != entryIndex {
					entries = append(entries, entry)
				}
			}
			c.table.SetCursor(cursor - 1)
			cmds = append(cmds, store.SetEntries(entries))

		case key.Matches(msg, util.Keys.AltDelete):
			cmds = append(cmds, store.SetEntries(make([]types.Entry, 0)))
		}

	case util.TimeTickMsg, messages.ClockInMsg, messages.ClockOutMsg:
		c.calculateTableRows()

	case store.StoreChangedMsg:
		// TODO: more specific messages
		// language change
		c.table = createTable()
		c.entries = store.GetEntries()
		c.calculateTableRows()
	}

	c.table, cmd = c.table.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c *Model) View() string {
	return c.table.View()
}

func (c *Model) calculateTableRows() {
	rows := make([]table.Row, 0)

	totalWorkTime := time.Time{}

	// calculate total work time (all entries)
	for _, entry := range c.entries {
		totalWorkTime = totalWorkTime.Add(entry.Duration())
	}

	for i := len(c.entries) - 1; i >= 0; i-- {
		entry := c.entries[i]

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

	c.table.SetRows(rows)

	const MAX_ROWS = 5

	if len(rows) > MAX_ROWS {
		c.table.SetHeight(MAX_ROWS + 2)
	} else {
		c.table.SetHeight(len(rows) + 2)
	}
}
