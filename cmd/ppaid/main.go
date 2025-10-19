package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddddddO/ppaid/internal"
	"github.com/ddddddO/ppaid/internal/model"
)

func main() {
	var shouldRestoreLatestExecutedData bool
	flag.BoolVar(&shouldRestoreLatestExecutedData, "repeat", false, "This flag starts with data selected by the most recently executed ppaid.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "ppaid")
		fmt.Fprintln(os.Stderr, "\nOptions:")

		flag.PrintDefaults()

		fmt.Fprintln(os.Stderr, "\nExample:")
		fmt.Fprintf(os.Stderr, "  %s          # normal launch\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --repeat # launch using the most recent data\n", os.Args[0])
	}
	flag.Parse()

	cfg, err := internal.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config. %v", err)
		os.Exit(1)
	}

	m, err := model.New(cfg, shouldRestoreLatestExecutedData)
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
