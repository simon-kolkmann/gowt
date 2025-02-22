package views

import (
	"gowt/types"
	"gowt/util"
	"slices"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Edit struct {
	begin   textinput.Model
	end     textinput.Model
	message string
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
		begin: begin,
		end:   end,
	}
}

func (e *Edit) Init() tea.Cmd {
	return nil
}

func (e *Edit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

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
			input := e.getFocusedTextInput()
			*input, cmd = input.Update(msg)
			e.autoFormatValue(input, msg)
			e.message = ""

		case "enter":
			validate(&e.begin)
			validate(&e.end)

			if e.hasError() {
				e.message = "❌ Mindestens eine Eingabe ist \nfehlerhaft und kann nicht \ngespeichert werden."
			} else {
				e.message = "Die Eingaben wurden gespeichert."
			}
		}
	}

	return e, cmd
}

func (e *Edit) View() string {
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

	e.begin.Prompt = "Beginn: "
	e.end.Prompt = "Ende: "

	return box.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			caption.Render("Eintrag bearbeiten\n"),
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				e.begin.View(),
				"   ",
				e.end.View(),
			),
			"",
			message.Render(e.message),
		),
	)
}

func (e *Edit) focusNext() {
	if e.begin.Focused() {
		e.begin.Blur()
		e.end.Focus()
	} else {
		e.end.Blur()
		e.begin.Focus()
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
	input.Err = input.Validate(input.Value())
}

func (e *Edit) getFocusedTextInput() *textinput.Model {
	if e.begin.Focused() {
		return &e.begin
	} else {
		return &e.end
	}
}

func (e *Edit) hasError() bool {
	return e.begin.Err != nil || e.end.Err != nil
}
