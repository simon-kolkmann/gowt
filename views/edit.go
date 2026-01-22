package views

import (
	"gowt/bubbles/time_input"
	"gowt/store"
	"gowt/types"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewEdit struct {
	entry       *types.Entry
	start       time_input.Model
	end         time_input.Model
	message     string
	showMessage bool
}

func NewEdit() ViewEdit {
	return ViewEdit{
		start: time_input.New(store.State().Strings().START + ": "),
		end:   time_input.New(store.State().Strings().END + ": "),
	}
}

func (view ViewEdit) Init() tea.Cmd {
	return nil
}

func (view ViewEdit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 3)

	view.start, cmd = view.start.Update(msg)
	cmds = append(cmds, cmd)

	view.end, cmd = view.end.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "tab", "shift+tab":
			view.focusNext()

		case
			"0",
			"1",
			"2",
			"3",
			"4",
			"5",
			"6",
			"7",
			"8",
			"9",
			"left",
			"right",
			"delete",
			"backspace":
			view.message = ""
			view.showMessage = false

		case "enter":
			start := view.start.GetTime()
			end := view.end.GetTime()
			cmds = append(cmds, store.Commit(store.ModifyActiveEntry(start, end)))
			view.showMessage = true

		case "ctrl+r":
			view.SetEntry(store.State().ActiveEntry)
		}

	case store.StateMutatedMsg:
		switch msg.Field {
		case store.FIELD_ACTIVE_VIEW:
			view.end.Input.Blur()
			view.start.Input.CursorEnd()
			cmds = append(cmds, view.start.Input.Focus())
			view.SetEntry(store.State().ActiveEntry)
			view.showMessage = false
		}
	}

	return view, tea.Batch(cmds...)
}

func (view ViewEdit) View() string {
	if view.showMessage {
		if view.hasError() {
			view.message = "❌" + store.State().Strings().ENTRY_SAVE_FAILED
		} else {
			view.message = store.State().Strings().ENTRY_SAVE_SUCCESS
		}
	}

	box := lipgloss.
		NewStyle().
		Padding(0, 2, 0, 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#ffffff"))

	caption := lipgloss.NewStyle().Bold(true).Underline(true).PaddingBottom(1)
	message := lipgloss.NewStyle().Bold(true)

	if view.hasError() {
		message = message.Foreground(lipgloss.Color(types.Theme.Error))
	} else {
		message = message.Foreground(lipgloss.Color(types.Theme.Success))
	}

	if view.entry == nil {
		return box.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				caption.Render(store.State().Strings().VIEW_EDIT),
				store.State().Strings().NO_ENTRY_SELECTED,
			),
		)
	}

	return box.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			caption.Render(store.State().Strings().VIEW_EDIT),
			store.State().Strings().KIND+": "+store.State().Strings().ENTRY_KIND(view.entry.Kind),
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				view.start.View(),
				"   ",
				view.end.View(),
			),
			"",
			message.Render(view.message),
		),
	)
}

func (view *ViewEdit) focusNext() {
	if view.start.Input.Focused() {
		view.start.Input.Blur()
		view.end.Input.Focus()
	} else {
		view.end.Input.Blur()
		view.start.Input.Focus()
	}
}

func (view *ViewEdit) hasError() bool {
	return view.start.Input.Err != nil || view.end.Input.Err != nil
}

func (view *ViewEdit) SetEntry(entry *types.Entry) {
	if entry == nil {
		view.entry = nil
		return
	}

	if entry.Start.IsZero() {
		view.start.Input.SetValue("")
	} else {
		view.start.Input.SetValue(entry.Start.Format(time.TimeOnly))
	}

	if entry.End.IsZero() {
		view.end.Input.SetValue("")
	} else {
		view.end.Input.SetValue(entry.End.Format(time.TimeOnly))
	}

	view.entry = entry
}
