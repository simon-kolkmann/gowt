package main

import (
	"gowt/i18n"
	"gowt/types"
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
	return setLanguage(i18n.LANG_GERMAN)
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return a, tea.Quit

		case "ctrl+l":
			if i18n.Selected == i18n.LANG_GERMAN {
				cmds = append(cmds, setLanguage(i18n.LANG_ENGLISH))
			}

			if i18n.Selected == i18n.LANG_ENGLISH {
				cmds = append(cmds, setLanguage(i18n.LANG_GERMAN))
			}

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

func setLanguage(l i18n.Language) tea.Cmd {
	return func() tea.Msg {
		i18n.Selected = l
		return types.LanguageChangedMsg(l)
	}
}
