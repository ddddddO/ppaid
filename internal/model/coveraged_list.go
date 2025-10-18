package model

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddddddO/gtree"
	"github.com/ddddddO/ppaid/internal"
	"github.com/ddddddO/ppaid/internal/command"
)

type coveragedListView struct {
}

func newCoverageListView() *coveragedListView {
	return &coveragedListView{}
}

func (c *coveragedListView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			m.quitting = true
			return m, tea.Quit

			// case "down", "j":
			// 	v.cursor++
			// 	if v.cursor >= len(v.choices) {
			// 		v.cursor = 0
			// 	}

			// case "up", "k":
			// 	v.cursor--
			// 	if v.cursor < 0 {
			// 		v.cursor = len(v.choices) - 1
			// 	}
		}
	}

	return m, tea.Quit
}

func (c *coveragedListView) view() string {
	s := &strings.Builder{}
	lvl := 2
	s.WriteString(fmt.Sprintf("\n\n--- Cveraged file list (Max depth: %d) ---\n\n", lvl+1))

	coverages, err := internal.GetCoveragedFilePaths(lvl)
	if err != nil {
		panic(err)
	}

	root := gtree.NewRoot(command.OutputCoverageDir)
	var node *gtree.Node
	for i := range coverages {
		for i, name := range strings.Split(coverages[i], string(filepath.Separator)) {
			if i == 0 {
				node = root
				continue
			}

			node = node.Add(name)
		}
	}
	for iter, err := range gtree.WalkIterFromRoot(root) {
		if err != nil {
			panic(err)
		}

		if iter.Level() == 1 {
			s.WriteString(fmt.Sprintf("%s\n", iter.Row()))
			continue
		}
		s.WriteString(fmt.Sprintf("  %s\n", iter.Row()))
	}
	s.WriteString("\n")

	// TODO: なんとかしたい
	s.WriteString("If it is not finished, press any key to finish it...")

	return s.String()
}
