package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/faustind/tomato/internal/cmd"
)

const (
	PROGRESS = "PROGRESS"
	DIGITAL  = "DIGITAL"
)

type duration struct {
	min int
	sec int
}

type window struct {
	height, width int
}

type progressBar struct {
	bar               progress.Model
	padding, maxWidth int
}

type model struct {
	ui                  string // how to show the timer
	window              window // current size of terminal
	progress            progressBar
	duration            duration
	durations           map[byte]duration
	pattern             string
	currentCountdownIdx int
}

// Initial returns an initial model according to given params
func Initial(pattern string, w, s, l int, withProgress bool) *model {
	durations := map[byte]duration{'w': {min: w, sec: 0}, 's': {min: s, sec: 0}, 'l': {min: l, sec: 0}}
	start := pattern[0]

	ui := DIGITAL

	mod := &model{
		currentCountdownIdx: 0,
		duration:            durations[start],
		durations:           durations,
		pattern:             pattern,
		ui:                  ui,
	}

	if withProgress {
		ui = PROGRESS
		prog := &progressBar{
			padding:  2,
			maxWidth: 80,
			bar: progress.New(
				progress.WithSolidFill("#dbbe88"),
				progress.WithoutPercentage(),
			),
		}
		prog.bar.SetPercent(1.0)
		mod.progress = *prog
	}

	return mod
}

func (m model) Init() tea.Cmd {
	// start ticking
	return cmd.Tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.window.width, m.window.height = msg.Width, msg.Height

		if m.ui == PROGRESS {
			m.progress.bar.Width = msg.Width - m.progress.padding*2
			if m.progress.bar.Width > m.progress.maxWidth {
				m.progress.bar.Width = m.progress.maxWidth
			}
		}
		return m, nil

	case cmd.TickMsg:
		durationIdx := m.pattern[m.currentCountdownIdx]
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
			return m, cmd.Tick()
		}

		min, sec := m.duration.min, m.duration.sec-1
		if sec < 0 {
			min, sec = min-1, 59
		}

		m.duration.min, m.duration.sec = min, sec

		perc := (m.duration.min*60 + m.duration.sec) * 100 / (m.durations[durationIdx].min*60 + m.durations[durationIdx].sec)

		progressCmd := m.progress.bar.SetPercent(float64(perc) / 100.0)
		return m, tea.Batch(cmd.Tick(), progressCmd)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case progress.FrameMsg:
		progressModel, cmd := m.progress.bar.Update(msg)
		m.progress.bar = progressModel.(progress.Model)
		return m, cmd
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

	countDown := countDownStyle.Render(fmt.Sprintf("%02d:%02d", m.duration.min, m.duration.sec))

	if m.ui == PROGRESS {
		countDown = m.progress.bar.View() + "\n"
	}

	statusMsg := lipgloss.NewStyle().Align(lipgloss.Center).Render(
		strings.Title(msgs[m.pattern[m.currentCountdownIdx]]),
	)

	ui := lipgloss.Place(
		m.window.width, m.window.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, countDown, statusMsg),
	)

	return ui
}
