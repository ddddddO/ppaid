package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddddddO/ppaid/internal"
	"github.com/ddddddO/ppaid/internal/model"
)

func main() {
	cfg, err := internal.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config. %v", err)
		os.Exit(1)
	}

	m, err := model.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize model. %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to program run. %v", err)
		os.Exit(1)
	}

}
