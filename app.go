package main

import (
	"gowt/bubbles/help"
	"gowt/i18n"
	"gowt/store"
	"gowt/types"
	"gowt/util"
	"gowt/views"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type app struct {
	clock    tea.Model
	settings tea.Model
	edit     tea.Model
	help     tea.Model
	width    int
	height   int
}

func NewApp() app {
	return app{
		clock:    views.NewClock(),
		settings: views.NewSettings(),
		edit:     views.NewEdit(),
		help:     help.NewHelp(),
	}
}

func (a app) Init() tea.Cmd {
	return tea.Batch(
		store.Init(),
	)
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	util.LogMessage(msg)

	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.height = msg.Height
		a.width = msg.Width

	case tea.KeyMsg:
		switch {

		// These keys should exit the program.
		case key.Matches(msg, util.Keys.Quit):
			return a, tea.Quit

		case key.Matches(msg, util.Keys.CtrlL):
			if store.GetLanguage() == i18n.LANG_ENGLISH {
				cmds = append(cmds, store.SetLanguage(i18n.LANG_GERMAN))
			} else {
				cmds = append(cmds, store.SetLanguage(i18n.LANG_ENGLISH))
			}

		case key.Matches(msg, util.Keys.CtrlLeft):
			if store.GetActiveView() > types.ViewSettings {
				cmds = append(cmds, store.SetActiveView(store.GetActiveView()-1))
			}

		case key.Matches(msg, util.Keys.CtrlRight):
			if store.GetActiveView() < types.ViewEdit {
				cmds = append(cmds, store.SetActiveView(store.GetActiveView()+1))
			}

		}
	}

	if store.GetActiveView() == types.ViewClock {
		a.clock, cmd = a.clock.Update(msg)
		cmds = append(cmds, cmd)
	}

	if store.GetActiveView() == types.ViewSettings {
		a.settings, cmd = a.settings.Update(msg)
		cmds = append(cmds, cmd)
	}

	if store.GetActiveView() == types.ViewEdit {
		a.edit, cmd = a.edit.Update(msg)
		cmds = append(cmds, cmd)
	}

	// always visible
	a.help, cmd = a.help.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a app) View() string {
	var activeView string

	switch store.GetActiveView() {
	case types.ViewClock:
		activeView = a.clock.View()

	case types.ViewSettings:
		activeView = a.settings.View()

	case types.ViewEdit:
		activeView = a.edit.View()

	default:
		activeView = "no active view"
	}

	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(a.width - 6).
		Render(a.help.View())

	content := lipgloss.NewStyle().
		Width(a.width - 6).
		Height(a.height - lipgloss.Height(footer) - 3).
		Align(lipgloss.Center).
		Render(activeView)

	box := lipgloss.
		NewStyle().Align(lipgloss.Center).
		Padding(1, 2, 0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(types.Theme.Primary))

	return box.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			content,
			footer,
		),
	)
}
