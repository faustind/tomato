package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/faustind/tomato/internal/model"
)

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

	prog := tea.NewProgram(model.Initial(*p, *w, *s, *l))
	if err := prog.Start(); err != nil {
		log.Fatalf("There's been an error: %v", err)
	}
}
