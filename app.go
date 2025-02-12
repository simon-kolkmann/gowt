package main

import (
	"gowt/views"

	tea "github.com/charmbracelet/bubbletea"
)

type app struct {
	clock views.Clock
}

func NewApp() app {
	return app{
		clock: views.NewClock(),
	}
}

func (a app) Init() tea.Cmd {
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return a, tea.Quit

		}

	}

	var cmd tea.Cmd

	_, cmd = a.clock.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a app) View() string {
	return a.clock.View()
}
