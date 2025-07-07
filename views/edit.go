package views

import (
	"gowt/bubbles/time_input"
	"gowt/messages"
	"gowt/store"
	"gowt/types"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Edit struct {
	entry       *types.Entry
	start       time_input.Model
	end         time_input.Model
	message     string
	showMessage bool
}

func NewEdit() Edit {
	return Edit{
		start: time_input.New(store.Strings().START + ": "),
		end:   time_input.New(store.Strings().END + ": "),
	}
}

func (e Edit) Init() tea.Cmd {
	return nil
}

func (e Edit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 3)

	e.start, cmd = e.start.Update(msg)
	cmds = append(cmds, cmd)

	e.end, cmd = e.end.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "tab", "shift+tab":
			e.focusNext()

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
			e.message = ""
			e.showMessage = false

		case "enter":
			cmds = append(cmds, store.UpdateActiveEntry(e.start.Time, e.end.Time))

			e.showMessage = true

		case "ctrl+r":
			e.SetEntry(store.GetActiveEntry())
		}

	case messages.ViewChangedMsg:
		e.end.Input.Blur()
		e.start.Input.CursorEnd()
		cmds = append(cmds, e.start.Input.Focus())
		e.SetEntry(store.GetActiveEntry())
		e.showMessage = false
	}

	return e, tea.Batch(cmds...)
}

func (e Edit) View() string {
	if e.showMessage {
		if e.hasError() {
			e.message = "❌" + store.Strings().ENTRY_SAVE_FAILED
		} else {
			e.message = store.Strings().ENTRY_SAVE_SUCCESS
		}
	}

	box := lipgloss.
		NewStyle().
		Padding(1, 2, 1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffffff"))

	caption := lipgloss.NewStyle().Bold(true).Underline(true)
	message := lipgloss.NewStyle().Bold(true)

	if e.hasError() {
		message = message.Foreground(lipgloss.Color(types.Theme.Error))
	} else {
		message = message.Foreground(lipgloss.Color(types.Theme.Success))
	}

	if e.entry == nil {
		return box.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				caption.Render(store.Strings().EDIT_ENTRY+"\n"),
				store.Strings().NO_ENTRY_SELECTED,
			),
		)
	}

	return box.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			caption.Render(store.Strings().EDIT_ENTRY+"\n"),
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				e.start.View(),
				"   ",
				e.end.View(),
			),
			"",
			message.Render(e.message),
		),
	)
}

func (e *Edit) focusNext() {
	if e.start.Input.Focused() {
		e.start.Input.Blur()
		e.end.Input.Focus()
	} else {
		e.end.Input.Blur()
		e.start.Input.Focus()
	}
}

func (e *Edit) hasError() bool {
	return e.start.Input.Err != nil || e.end.Input.Err != nil
}

func (e *Edit) SetEntry(entry *types.Entry) {
	if entry == nil {
		e.entry = nil
		return
	}

	if entry.Start.IsZero() {
		e.start.Input.SetValue("")
	} else {
		e.start.Input.SetValue(entry.Start.Format(time.TimeOnly))
	}

	if entry.End.IsZero() {
		e.end.Input.SetValue("")
	} else {
		e.end.Input.SetValue(entry.End.Format(time.TimeOnly))
	}

	e.entry = entry
}
