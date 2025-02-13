package main

import (
	"gowt/i18n"
	"gowt/types"
	"gowt/views"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

type app struct {
	dump  io.Writer
	clock views.Clock
}

func NewApp() app {
	var dump *os.File
	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}

	return app{
		clock: views.NewClock(),
		dump:  dump,
	}
}

func (a app) Init() tea.Cmd {
	return setLanguage(i18n.LANG_GERMAN)
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if a.dump != nil {
		spew.Fdump(a.dump, msg)
	}

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
