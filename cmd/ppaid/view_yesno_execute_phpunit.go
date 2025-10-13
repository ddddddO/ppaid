package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type yesnoView struct {
	cursor   int
	choices  []string
	selected string

	cmdPHPUnit *cmdPHPUnit
}

func newYesNoView() (*yesnoView, error) {
	return &yesnoView{
		choices:    []string{"Yes", "No"},
		cmdPHPUnit: &cmdPHPUnit{},
	}, nil
}

func (v *yesnoView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			v.selected = v.choices[v.cursor]

			if v.selected == "Yes" {
				// TODO: 一旦ここに置くだけ
				if err := os.RemoveAll(outputCoverageDir); err != nil {
					panic(err)
				}
				if err := os.Remove("tmp_phpunit.xml"); err != nil {
					panic(err)
				}

				targetTests := []string{}
				for s := range m.selectTestFilesView.selected {
					targetTests = append(targetTests, s)
				}

				if err := generatePHPUnitXML(targetTests, m.selectCoverageFilesView.longestMatchDirPath()); err != nil {
					panic(err)
				}

				return m, v.cmdPHPUnit.command()
			}

			m.quitting = true
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

func (v *yesnoView) view(cfv *selectCoverageFilesView) string {
	s := strings.Builder{}
	s.WriteString("Execute PHPUnit?\n\n")

	// v.cmdPHPUnit.build("./src", "PPAID", "tmp_phpunit.xml")
	v.cmdPHPUnit.build(cfv.longestMatchDirPath(), "PPAID", "tmp_phpunit.xml")

	s.WriteString(fmt.Sprintf("%s\n\n", v.cmdPHPUnit.rawCmd()))

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
