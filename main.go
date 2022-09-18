package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type duration struct {
	min int
	sec int
}
type model struct {
	duration            duration
	durations           map[byte]duration
	pattern             string
	currentCountdownIdx int
}

func initialModel(pattern string, w, s, l int) *model {
	durations := map[byte]duration{'w': {w, 0}, 's': {s, 0}, 'l': {l, 0}}
	start := pattern[0]
	return &model{
		currentCountdownIdx: 0,
		duration:            durations[start],
		durations:           durations,
		pattern:             pattern,
	}
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	// start ticking
	return doTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case TickMsg:
		if m.duration.min == 0 && m.duration.sec == 0 {
			// move to next counter in pattern
			idx := m.currentCountdownIdx + 1
			if idx > len(m.pattern)-1 {
				idx = 0
			}
			m.currentCountdownIdx = idx

			//set duration for current counter
			m.duration = m.durations[m.pattern[m.currentCountdownIdx]]
			return m, doTick()
		}

		min, sec := m.duration.min, m.duration.sec-1
		if sec < 0 {
			min, sec = min-1, 59
		}

		m.duration.min, m.duration.sec = min, sec
		return m, doTick()

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

func main() {
	var p = flag.String("p", "wswswl", "Pattern of a working session")
	var w = flag.Int("w", 25, "Duration of a working session")
	var s = flag.Int("s", 5, "Duration of a short break")
	var l = flag.Int("l", 10, "Duration of a long break")

	flag.Parse()
	// validate flags
	if *p == "" {
		fmt.Printf("Invalid pattern ''%s', should not be empty\n", *p)
		os.Exit(2)
	}
	for _, c := range *p {
		if c != 'w' && c != 'l' && c != 's' {
			fmt.Printf("Invalid pattern ''%s', should contain only w,s, or l\n", *p)
			os.Exit(2)
		}
	}

	prog := tea.NewProgram(initialModel(*p, *w, *s, *l))
	if err := prog.Start(); err != nil {
		log.Fatalf("There's been an error: %v", err)
	}
}
