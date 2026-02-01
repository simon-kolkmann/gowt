package time_progress

import (
	"gowt/types"
	"gowt/util"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	progress progress.Model
	label    string
	Elapsed  time.Duration
	Target   time.Duration
	Color    string
}

type Option func(*Model)

func WithLabel(label string) Option {
	return func(m *Model) {
		m.setLabel(label)
	}
}

func NewTimeProgress(opts ...Option) Model {
	m := Model{
		progress: progress.New(
			progress.WithSolidFill(types.Theme.Success),
			progress.WithWidth(50),
			progress.WithoutPercentage(),
		),
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

func (m *Model) setLabel(label string) {
	m.label = label
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
	row := lipgloss.NewStyle().Margin(0, 0, 1, 0).Render

	elapsed := m.Elapsed.String()
	target := m.Target.String()
	remaining := m.Remaining().String()
	percent, percentAsString := util.ElapsedInPercent(m.Elapsed, m.Target)

	m.progress.FullColor = m.Color

	components := []string{}

	if m.label != "" {
		components = append(components, row(m.label))
	}

	components = append(components,
		row(m.progress.ViewAs(percent/100)),
		row(elapsed+" / "+target+" ("+remaining+", "+percentAsString+")"))

	return lipgloss.JoinVertical(
		lipgloss.Center,
		components...,
	)
}

func (m *Model) Remaining() time.Duration {
	return m.Target - m.Elapsed
}
