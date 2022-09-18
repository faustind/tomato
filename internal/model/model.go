package model

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/faustind/tomato/internal/cmd"
)

type model struct {
	duration            duration
	durations           map[byte]duration
	pattern             string
	currentCountdownIdx int
}

type duration struct {
	min int
	sec int
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

	s := fmt.Sprintf("tomato: %s\n", msgs[m.pattern[m.currentCountdownIdx]])
	s += fmt.Sprintf("%02d:%02d", m.duration.min, m.duration.sec)
	s += "\npress q to quit.\n"
	return s
}
