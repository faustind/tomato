package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	min int
	sec int
}

func initialModel() *model {
	return &model{
		min: 1,
		sec: 0,
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
		if m.min == 0 && m.sec == 0 {
			return m, nil
		}

		min, sec := m.min, m.sec-1
		if sec < 0 {
			min, sec = min-1, 59
		}

		m.min, m.sec = min, sec
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
	s := fmt.Sprintf("tomato\n")
	s += fmt.Sprintf("%2d:%2d", m.min, m.sec)
	s += "\npress q to quit.\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		log.Fatalf("There's been an error: %v", err)
	}
}
