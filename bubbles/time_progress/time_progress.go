package time_progress

import (
	"gowt/types"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	progress progress.Model
	Label    string
	Elapsed  time.Duration
	Target   time.Duration
	Color    string
}

func NewTimeProgress() Model {
	return Model{
		progress: progress.New(
			progress.WithSolidFill(types.Theme.Success),
			progress.WithWidth(50),
			progress.WithoutPercentage(),
		),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)

	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Target == 0 {
		return ""
	}

	row := lipgloss.NewStyle().Margin(0, 0, 0, 0).Render

	elapsed := m.Elapsed.String()
	target := m.Target.String()
	remaining := m.Remaining().String()
	percent := m.RemainingInPercent()
	percentAsString := strconv.FormatFloat(percent, 'f', 2, 64) + "%"

	m.progress.FullColor = m.Color

	components := []string{}

	if m.Label != "" {
		components = append(components, row(m.Label))
	}

	components = append(components,
		row(m.progress.ViewAs(percent/100)),
		row(elapsed+" / "+target+" ("+remaining+", "+percentAsString+")"))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		components...,
	)
}

func (m *Model) Remaining() time.Duration {
	return m.Target - m.Elapsed
}

func (m *Model) RemainingInPercent() float64 {
	return m.Elapsed.Seconds() / (m.Target.Seconds() / 100)
}
