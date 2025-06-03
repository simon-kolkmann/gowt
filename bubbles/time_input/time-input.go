package time_input

import (
	"gowt/util"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Time  time.Time
	Input textinput.Model
}

func New(prompt string) Model {
	input := textinput.New()
	input.Placeholder = "hh:mm:ss"
	input.CharLimit = 8
	input.Width = 8
	input.Validate = util.Validators.Time
	input.Prompt = prompt

	return Model{
		Input: input,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.autoFormatValue(msg)

		switch msg.String() {
		case "enter":
			value := m.Input.Value()
			m.validate()

			if m.Input.Err != nil || value == "" {
				break
			}

			now := time.Now()
			t, _ := time.Parse(time.TimeOnly, value)

			m.Time = time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				t.Hour(),
				t.Minute(),
				t.Second(),
				0,
				now.Location(),
			)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.Input.View()
}

func (m *Model) autoFormatValue(msg tea.KeyMsg) {
	v := m.Input.Value()

	isRemovingOrNavigating := slices.Contains(
		[]string{"backspace", "delete", "left", "right"},
		msg.String(),
	)

	if isRemovingOrNavigating {
		return
	}

	// insert ":" automatically
	if len(v) == 2 || len(v) == 5 {
		m.Input.SetValue(m.Input.Value() + ":")
		m.Input.CursorEnd()
	}
}

func (m *Model) validate() {
	input := m.Input
	value := m.Input.Value()

	if len(value) == 0 {
		input.Err = nil
	} else {
		input.Err = input.Validate(value)
	}
}
