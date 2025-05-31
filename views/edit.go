package views

import (
	"gowt/messages"
	"gowt/store"
	"gowt/types"
	"gowt/util"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Edit struct {
	entry       *types.Entry
	start       textinput.Model
	end         textinput.Model
	message     string
	showMessage bool
}

func NewEdit() Edit {
	begin := textinput.New()
	begin.Placeholder = "hh:mm:ss"
	begin.CharLimit = 8
	begin.Width = 8
	begin.Validate = util.Validators.Time

	end := textinput.New()
	end.Placeholder = "hh:mm:ss"
	end.CharLimit = 8
	end.Width = 8
	end.Validate = util.Validators.Time

	begin.Focus()

	return Edit{
		start: begin,
		end:   end,
	}
}

func (e Edit) Init() tea.Cmd {
	return nil
}

func (e Edit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 3)

	start, cmd := e.start.Update(msg)
	e.start = start
	cmds = append(cmds, cmd)

	end, cmd := e.end.Update(msg)
	e.end = end
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
			input, cmd := e.getFocusedTextInput().Update(msg)
			e.autoFormatValue(&input, msg)
			e.message = ""
			e.showMessage = false
			cmds = append(cmds, cmd)

		case "enter":
			validate(&e.start)
			validate(&e.end)

			if !e.hasError() {
				now := time.Now()

				// FIXME: well...
				if e.start.Value() != "" {
					start, _ := time.Parse(time.TimeOnly, e.start.Value())
					start = time.Date(now.Year(), now.Month(), now.Day(), start.Hour(), start.Minute(), start.Second(), 0, now.Location())
					store.GetActiveEntry().Start = start
				}

				if e.end.Value() != "" {
					end, _ := time.Parse(time.TimeOnly, e.end.Value())
					end = time.Date(now.Year(), now.Month(), now.Day(), end.Hour(), end.Minute(), end.Second(), 0, now.Location())
					store.GetActiveEntry().End = end
				}

				// FIXME: persist in a nicer way
				cmds = append(cmds, store.SetEntries(store.GetEntries()))
			}

			e.showMessage = true
		}

	case messages.ViewChangedMsg:
		e.end.Blur()
		e.start.CursorEnd()
		cmds = append(cmds, e.start.Focus())
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

	e.start.Prompt = store.Strings().START + ": "
	e.end.Prompt = store.Strings().END + ": "

	if e.entry == nil {
		return box.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				caption.Render(store.Strings().EDIT_ENTRY+"\n"),
				store.Strings().NO_ENTRY_SELECTED,
			),
		)
	} else {
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

}

func (e *Edit) focusNext() {
	if e.start.Focused() {
		e.start.Blur()
		e.end.Focus()
	} else {
		e.end.Blur()
		e.start.Focus()
	}
}

func (e *Edit) autoFormatValue(input *textinput.Model, msg tea.KeyMsg) {
	v := input.Value()

	isRemovingOrNavigating := slices.Contains(
		[]string{"backspace", "delete", "left", "right"},
		msg.String(),
	)

	if isRemovingOrNavigating {
		return
	}

	// insert ":" automatically
	if len(v) == 2 || len(v) == 5 {
		input.SetValue(input.Value() + ":")
		input.CursorEnd()
	}
}

func validate(input *textinput.Model) {
	value := input.Value()

	if len(value) == 0 {
		input.Err = nil
	} else {
		input.Err = input.Validate(input.Value())
	}

}

func (e *Edit) getFocusedTextInput() *textinput.Model {
	if e.start.Focused() {
		return &e.start
	} else {
		return &e.end
	}
}

func (e *Edit) hasError() bool {
	return e.start.Err != nil || e.end.Err != nil
}

func (e *Edit) SetEntry(entry *types.Entry) {
	if entry == nil {
		e.entry = nil
		return
	}

	if entry.Start.IsZero() {
		e.start.SetValue("")
	} else {
		e.start.SetValue(entry.Start.Format(time.TimeOnly))
	}

	if entry.End.IsZero() {
		e.end.SetValue("")
	} else {
		e.end.SetValue(entry.End.Format(time.TimeOnly))
	}

	e.entry = entry
}
