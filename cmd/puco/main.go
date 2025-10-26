package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddddddO/puco/internal"
	"github.com/ddddddO/puco/internal/model"
)

func main() {
	var shouldRestoreLatestExecutedData bool
	flag.BoolVar(&shouldRestoreLatestExecutedData, "repeat", false, "This flag starts with data selected by the most recently executed puco.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "puco")
		fmt.Fprintln(os.Stderr, "\nOptions:")

		flag.PrintDefaults()

		fmt.Fprintln(os.Stderr, "\nExample:")
		fmt.Fprintf(os.Stderr, "  %s          # normal launch\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --repeat # launch using the most recent data\n", os.Args[0])

		fmt.Fprintln(os.Stderr, "\nProcessing description:")
		fmt.Fprintln(os.Stderr, "  1. You can select multiple test files to run (fuzzy search available).")
		fmt.Fprintln(os.Stderr, "  2. You can select multiple PHP files for which you want to calculate coverage (fuzzy search available).")
		fmt.Fprintln(os.Stderr, "  3. Calculate the longest matching directory path from multiple selected PHP file paths in step 2")
		fmt.Fprintln(os.Stderr, "    - ※ Note that only the PHP file paths selected in step 2 are not the target for coverage calculation. Instead, the directory path under the longest match calculated becomes the target for coverage calculation. If there are numerous PHP files under the calculated directory path, the coverage calculation process may become slow.")
		fmt.Fprintln(os.Stderr, "  4. If Steps 1 and 3 and an existing phpunit.xml are present, generate phpunitxml_generated_by_puco.xml based on them.")
		fmt.Fprintln(os.Stderr, "  5. Assemble and execute the php command.")
		fmt.Fprintln(os.Stderr, "  6. Coverage reports are generated under the coverage-puco directory.")

		fmt.Fprintln(os.Stderr, "\nWARNING:")
		fmt.Fprintln(os.Stderr, `  When puco is run for the first time, a configuration file named ~/.config/puco.toml is created. This configuration file contains a key: CommandToSpecifyBeforePHPCommand. It specifies that the PHP command should be executed via the Docker command. If you wish to execute the PHP command directly, please set the value of this key to "" or delete this entire line.`)

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
