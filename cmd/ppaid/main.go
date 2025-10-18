package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddddddO/ppaid/internal/model"
)

func main() {
	m, err := model.New()
	if err != nil {
		fmt.Printf("XXxx: %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("XXxx: %v", err)
		os.Exit(1)
	}

}
