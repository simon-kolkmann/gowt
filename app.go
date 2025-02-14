package main

import (
	"gowt/bubbles/help"
	"gowt/i18n"
	"gowt/types"
	"gowt/util"
	"gowt/views"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

type app struct {
	dump       io.Writer
	clock      views.Clock
	settings   views.Settings
	help       help.Model
	activeView types.View
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
		clock:      views.NewClock(),
		settings:   views.NewSettings(),
		help:       help.NewHelp(),
		activeView: types.ViewClock,
		dump:       dump,
	}
}

func (a app) Init() tea.Cmd {
	return setLanguage(i18n.LANG_GERMAN)
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if a.dump != nil {
		_, isTimeTickMsg := msg.(util.TimeTickMsg)
		_, isBlinkMsg := msg.(cursor.BlinkMsg)

		if !isTimeTickMsg && !isBlinkMsg {
			spew.Fdump(a.dump, msg)
		}
	}

	var cmd tea.Cmd
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

		case "ctrl+left":
			if a.activeView > types.ViewSettings {
				a.activeView--
			}

			cmds = append(cmds, sendViewChangedMsg(a.activeView))

		case "ctrl+right":
			if a.activeView < types.ViewClock {
				a.activeView++
			}

			cmds = append(cmds, sendViewChangedMsg(a.activeView))

		}

	case types.StoreChangedMsg:
		_, cmd = a.clock.Update(msg)
		cmds = append(cmds, cmd)
	}

	if a.activeView == types.ViewClock {
		_, cmd = a.clock.Update(msg)
		cmds = append(cmds, cmd)
	}

	if a.activeView == types.ViewSettings {
		_, cmd = a.settings.Update(msg)
		cmds = append(cmds, cmd)
	}

	// always visible
	_, cmd = a.help.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a app) View() string {
	var activeView string

	switch a.activeView {
	case types.ViewClock:
		activeView = a.clock.View()

	case types.ViewSettings:
		activeView = a.settings.View()

	default:
		activeView = "no active view"
	}

	box := lipgloss.
		NewStyle().
		Padding(1, 2, 0, 2).
		Margin(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(types.Theme.Primary))

	return box.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			activeView,
			a.help.View(),
		),
	)
}

func setLanguage(l types.Language) tea.Cmd {
	return func() tea.Msg {
		i18n.Selected = l
		return types.LanguageChangedMsg(l)
	}
}

func sendViewChangedMsg(v types.View) tea.Cmd {
	return func() tea.Msg {
		return types.ViewChangedMsg(v)
	}
}
