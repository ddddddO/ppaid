package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddddddO/gtree"
)

type coverageListView struct {
}

func newCoverageListView() *coverageListView {
	return &coverageListView{}
}

func (c *coverageListView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
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

func (c *coverageListView) view() string {
	s := &strings.Builder{}
	lvl := 2
	s.WriteString(fmt.Sprintf("\n\n--- Cverage list (Max depth: %d) ---\n\n", lvl+1))

	coverages, err := getCoverageList(lvl)
	if err != nil {
		panic(err)
	}

	root := gtree.NewRoot(outputCoverageDir)
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
		// s.WriteString(fmt.Sprintln(iter.Path()))
		s.WriteString(fmt.Sprintln(iter.Row()))
	}
	s.WriteString("\n")

	// TODO: なんとかしたい
	s.WriteString("If it is not finished, press any key to finish it...")

	return s.String()
}

func getCoverageList(level int) ([]string, error) {
	ignore := []string{"_css", "_icons", "_js"}

	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	target := filepath.Join(root, outputCoverageDir)

	paths := []string{}
	err = filepath.WalkDir(target, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == root {
			return nil
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		lvl := strings.Count(relativePath, string(filepath.Separator))
		if lvl > level {
			return nil
		}

		if d.IsDir() && slices.Contains(ignore, d.Name()) {
			return fs.SkipDir
		}

		paths = append(paths, relativePath)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}
