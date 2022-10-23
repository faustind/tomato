// Package internal contains the bubbletea model and its helper
package internal

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type duration struct {
	min int
	sec int
}

type window struct {
	height, width int
}

type tickMsg time.Time

// Tick returns a tea.Cmd that returns a tickMsg every 1 second
func Tick() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type model struct {
	window              window // current size of terminal
	duration            duration
	durations           map[byte]duration
	pattern             string
	currentCountdownIdx int
}

// NewModel returns an initial model according to given params
func NewModel(pattern string, w, s, l int) *model {
	durations := map[byte]duration{
		'w': {min: w, sec: 0}, // duration of a working session
		's': {min: s, sec: 0}, // duration of a short break
		'l': {min: l, sec: 0}, // duration of a long break
	}
	start := pattern[0]

	mod := &model{
		currentCountdownIdx: 0,
		duration:            durations[start],
		durations:           durations,
		pattern:             pattern,
	}

	return mod
}

func (m model) Init() tea.Cmd {
	// start ticking
	return Tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.window.width, m.window.height = msg.Width, msg.Height
		return m, nil

	case tickMsg:
		var durationIdx byte //:= m.pattern[m.currentCountdownIdx]
		if m.duration.min == 0 && m.duration.sec == 0 {
			// move to next counter in pattern
			idx := m.currentCountdownIdx + 1
			if idx > len(m.pattern)-1 {
				idx = 0
			}
			m.currentCountdownIdx = idx
			durationIdx = m.pattern[m.currentCountdownIdx]

			//set duration for current counter
			m.duration = m.durations[durationIdx]
			return m, Tick()
		}

		min, sec := m.duration.min, m.duration.sec-1
		if sec < 0 {
			min, sec = min-1, 59
		}

		m.duration.min, m.duration.sec = min, sec

		return m, Tick()

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

	countDownStyle := lipgloss.NewStyle().Align(lipgloss.Center)

	timerString := fmt.Sprintf("%02d:%02d", m.duration.min, m.duration.sec)
	timer := ""
	for _, char := range timerString {
		timer = lipgloss.JoinHorizontal(lipgloss.Center, timer, drawChar(char))
	}

	countDown := countDownStyle.Render(timer)

	statusMsg := lipgloss.NewStyle().Align(lipgloss.Center).Render(
		strings.ToUpper(msgs[m.pattern[m.currentCountdownIdx]]),
	)

	ui := lipgloss.Place(
		m.window.width, m.window.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, countDown, statusMsg),
	)

	return ui
}
