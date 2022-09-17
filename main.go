package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	hour int
	min  int
	sec  int
}

func initialModel() *model {
	now := time.Now()
	hour, min, sec := now.Clock()
	return &model{
		hour,
		min,
		sec,
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
		t := time.Time(msg)
		m.hour, m.min, m.sec = t.Clock()
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
	s += fmt.Sprintf("%2d:%2d:%2d", m.hour, m.min, m.sec)
	s += "\npress q to quit.\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		log.Fatalf("There's been an error: %v", err)
	}
}
