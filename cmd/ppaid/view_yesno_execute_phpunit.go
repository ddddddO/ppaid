package main

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type yesnoView struct {
	cursor   int
	choices  []string
	selected string
}

func newYesNoView() (*yesnoView, error) {
	return &yesnoView{
		choices: []string{"Yes", "No"},
	}, nil
}

func (v *yesnoView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			m.quitting = true
			v.selected = v.choices[v.cursor]

			if v.selected == "Yes" {
				// TODO: 一旦ここに置くだけ
				targetTests := []string{}
				for s := range m.selectTestFilesView.selected {
					targetTests = append(targetTests, filepath.Join("./tests", s))
				}

				if err := generatePHPUnitXML(targetTests); err != nil {
					panic(err)
				}
			}

			return m, tea.Quit

		case "down", "j":
			v.cursor++
			if v.cursor >= len(v.choices) {
				v.cursor = 0
			}

		case "up", "k":
			v.cursor--
			if v.cursor < 0 {
				v.cursor = len(v.choices) - 1
			}
		}
	}

	return m, nil
}

const EXEC_PHPUNIT = `php -d pcov.directory=%s vendor/bin/phpunit --testsuite "%s" --configuration %s --coverage-html coverage`

func (v *yesnoView) view() string {
	s := strings.Builder{}
	s.WriteString("Execute PHPUnit?\n\n")

	execCmd := fmt.Sprintf(EXEC_PHPUNIT, "./src", "PPAID", "tmp_phpunit.xml")

	s.WriteString(fmt.Sprintf("%s\n\n", execCmd))

	for i := 0; i < len(v.choices); i++ {
		if v.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(v.choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}
