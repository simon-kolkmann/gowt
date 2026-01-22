package main

import (
	"gowt/bubbles/help"
	"gowt/bubbles/tabs"
	"gowt/store"
	"gowt/types"
	"gowt/util"
	"gowt/views"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type app struct {
	tabs     tabs.Model
	clock    tea.Model
	settings tea.Model
	edit     tea.Model
	help     tea.Model
	width    int
	height   int
}

func NewApp() app {
	tabs := tabs.NewTabs([]tabs.Tab{
		{Text: store.State().Strings().VIEW_SETTINGS, Value: types.ViewSettings},
		{Text: store.State().Strings().VIEW_CLOCK, Value: types.ViewClock},
		{Text: store.State().Strings().VIEW_EDIT, Value: types.ViewEdit},
	})

	return app{
		tabs:     tabs,
		clock:    views.NewClock(),
		settings: views.NewSettings(),
		edit:     views.NewEdit(),
		help:     help.NewHelp(),
	}
}

func (a app) Init() tea.Cmd {
	return store.Initialize()
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	util.LogMessage(msg)

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
			cmds = append(cmds, store.Commit(store.ToggleLanguage()))

		case key.Matches(msg, util.Keys.CtrlLeft):
			if store.State().ActiveView > types.ViewSettings {
				cmds = append(cmds, store.Commit(store.SetActiveView(store.State().ActiveView-1)))
			}

		case key.Matches(msg, util.Keys.CtrlRight):
			if store.State().ActiveView < types.ViewEdit {
				nextIsEdit := store.State().ActiveView == types.ViewEdit-1
				if !(nextIsEdit && store.State().ActiveEntry == nil) {
					cmds = append(cmds, store.Commit(store.SetActiveView(store.State().ActiveView+1)))
				}
			}
		}

	case store.StateMutatedMsg:
		switch msg.Field {
		case store.FIELD_ACTIVE_VIEW:
			a.tabs.SetActive(store.State().ActiveView)
		case store.FIELD_ACTIVE_ENTRY:
			if store.State().ActiveEntry == nil {
				a.tabs.SetInactive([]any{types.ViewEdit})
			} else {
				a.tabs.SetInactive([]any{})
			}
		case store.FIELD_LANGUAGE:
			a.tabs.UpdateTab(tabs.Tab{
				Text:  store.State().Strings().VIEW_SETTINGS,
				Value: types.ViewSettings},
			)
			a.tabs.UpdateTab(tabs.Tab{
				Text:  store.State().Strings().VIEW_CLOCK,
				Value: types.ViewClock},
			)
			a.tabs.UpdateTab(tabs.Tab{
				Text:  store.State().Strings().VIEW_EDIT,
				Value: types.ViewEdit},
			)
		}
	}

	cmds = append(cmds, a.UpdateActiveView(msg))
	cmds = append(cmds, a.UpdateAlwaysVisible(msg))

	return a, tea.Batch(cmds...)
}

func (a app) View() string {
	var activeView string

	switch store.State().ActiveView {
	case types.ViewClock:
		activeView = a.clock.View()

	case types.ViewSettings:
		activeView = a.settings.View()

	case types.ViewEdit:
		activeView = a.edit.View()

	default:
		activeView = "no active view"
	}

	header := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Render(a.tabs.View())

	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(a.width - 6).
		Render(a.help.View())

	content := lipgloss.NewStyle().
		Width(a.width - 6).
		Height(a.height - lipgloss.Height(footer) - 5).
		Align(lipgloss.Center).
		Render(activeView)

	box := lipgloss.
		NewStyle().Align(lipgloss.Center).
		Border(lipgloss.NormalBorder())

	if store.State().IsClockedIn() {
		box = box.BorderForeground(lipgloss.Color(types.Theme.Success))
	} else if store.State().IsAtBreak() {
		box = box.BorderForeground(lipgloss.Color(types.Theme.Warn))
	} else {
		box = box.BorderForeground(lipgloss.Color(types.Theme.Error))
	}

	return box.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			header,
			content,
			footer,
		),
	)
}

func (a *app) UpdateActiveView(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	if store.State().ActiveView == types.ViewClock {
		a.clock, cmd = a.clock.Update(msg)
		return cmd
	}

	if store.State().ActiveView == types.ViewSettings {
		a.settings, cmd = a.settings.Update(msg)
		return cmd
	}

	if store.State().ActiveView == types.ViewEdit {
		a.edit, cmd = a.edit.Update(msg)
		return cmd
	}

	return nil
}

func (a *app) UpdateAlwaysVisible(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd = []tea.Cmd{}

	model, cmd := a.tabs.Update(msg)
	a.tabs = model.(tabs.Model)
	cmds = append(cmds, cmd)

	a.help, cmd = a.help.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}
