package tabs

import (
	"gowt/types"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Tab struct {
	Text  string
	Value any
}

type Model struct {
	tabs     []Tab
	Position lipgloss.Position
	active   any
	inactive []any
}

func NewTabs(tabs []Tab) Model {
	return Model{
		tabs:     tabs,
		Position: lipgloss.Left,
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return model, nil
}

func (model Model) View() string {
	styleTab := lipgloss.NewStyle().
		Padding(0, 2, 0, 2)

	switch model.Position {
	// case lipgloss.Center:
	// 	styleTab = styleTab.Border(lipgloss.NormalBorder(), false, true, false, true)
	case lipgloss.Left:
		styleTab = styleTab.Border(lipgloss.NormalBorder(), false, true, false, false)
	case lipgloss.Right:
		styleTab = styleTab.Border(lipgloss.NormalBorder(), false, false, false, true)
	}

	styleActiveTab := styleTab.Background(lipgloss.Color(types.Theme.Primary))
	styleInactiveTab := styleTab.Faint(true)

	renderedTabs := make([]string, len(model.tabs))

	for _, tab := range model.tabs {
		if model.active == tab.Value {
			renderedTabs = append(renderedTabs, styleActiveTab.Render(tab.Text))
		} else if slices.Contains(model.inactive, tab.Value) {
			renderedTabs = append(renderedTabs, styleInactiveTab.Render(tab.Text))
		} else {
			renderedTabs = append(renderedTabs, styleTab.Render(tab.Text))
		}
	}

	return strings.Join(renderedTabs, "")
}

func (model *Model) SetActive(value any) {
	model.active = value
}

func (model *Model) SetInactive(values []any) {
	model.inactive = values
}

func (model *Model) AddTab(tab Tab) {
	model.tabs = append(model.tabs, tab)
}

func (model *Model) UpdateTab(tab Tab) {
	idx := slices.IndexFunc(model.tabs, func(candidate Tab) bool {
		return candidate.Value == tab.Value
	})

	if idx == -1 {
		return
	}

	model.tabs[idx] = tab
}

func (model *Model) RemoveTab(tab Tab) {
	model.tabs = slices.DeleteFunc(model.tabs, func(candidate Tab) bool {
		return candidate.Value == tab.Value
	})
}
