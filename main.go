package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/faustind/tomato/internal"
)

func main() {
	var (
		p = flag.String("p", "wswswl", "Pattern of a working session")
		w = flag.Int("w", 25, "Duration of a working session")
		s = flag.Int("s", 5, "Duration of a short break")
		l = flag.Int("l", 10, "Duration of a long break")
	)

	flag.Parse()
	// validate flags
	if *p == "" {
		fmt.Printf("Invalid pattern ''%s', should not be empty\n", *p)
		os.Exit(2)
	}

	for _, c := range *p {
		if c != 'w' && c != 'l' && c != 's' {
			fmt.Printf("Invalid pattern '%s', should contain only w,s, or l\n", *p)
			os.Exit(2)
		}
	}

	opts := []tea.ProgramOption{
		tea.WithAltScreen(), // full screen
	}

	initialModel := internal.NewModel(*p, *w, *s, *l)

	prog := tea.NewProgram(initialModel, opts...)
	if err := prog.Start(); err != nil {
		log.Fatalf("There's been an error: %v", err)
	}
}
