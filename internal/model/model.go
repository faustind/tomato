package model

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/faustind/tomato/internal/cmd"
)

type duration struct {
	min int
	sec int
}

type window struct {
	height, width int
}

type model struct {
	window              window
	duration            duration
	durations           map[byte]duration
	pattern             string
	currentCountdownIdx int
}

func Initial(pattern string, w, s, l int) *model {
	durations := map[byte]duration{'w': {w, 0}, 's': {s, 0}, 'l': {l, 0}}
	start := pattern[0]
	return &model{
		currentCountdownIdx: 0,
		duration:            durations[start],
		durations:           durations,
		pattern:             pattern,
	}
}

func (m model) Init() tea.Cmd {
	// start ticking
	return cmd.Tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.window.width, m.window.height = msg.Width, msg.Height
		return m, nil

	case cmd.TickMsg:
		if m.duration.min == 0 && m.duration.sec == 0 {
			// move to next counter in pattern
			idx := m.currentCountdownIdx + 1
			if idx > len(m.pattern)-1 {
				idx = 0
			}
			m.currentCountdownIdx = idx

			//set duration for current counter
			m.duration = m.durations[m.pattern[m.currentCountdownIdx]]
			return m, cmd.Tick()
		}

		min, sec := m.duration.min, m.duration.sec-1
		if sec < 0 {
			min, sec = min-1, 59
		}

		m.duration.min, m.duration.sec = min, sec
		return m, cmd.Tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	msgs := map[byte]string{
		'w': "working session",
		's': "short break",
		'l': "long break",
	}

	timer := lipgloss.NewStyle().Align(lipgloss.Center).Render(
		fmt.Sprintf("%02d:%02d", m.duration.min, m.duration.sec),
	)

	statusMsg := lipgloss.NewStyle().Align(lipgloss.Center).Render(
		strings.Title(msgs[m.pattern[m.currentCountdownIdx]]),
	)

	ui := lipgloss.Place(
		m.window.width, m.window.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, timer, statusMsg),
	)

	return ui
}
