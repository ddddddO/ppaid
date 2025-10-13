package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type selectCoverageFilesView struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func newSelectCoverageFilesView() (*selectCoverageFilesView, error) {
	paths, err := getPHPCodeFilePaths()
	if err != nil {
		return nil, err
	}

	// TODO: 多分ここで、選択されたカバレッジ取りたいファイルパスをマージして、最大公約数的なパスを算出して以下に設定できればよさそうな気がする？
	// - pcov.directory=
	// - phpunit.xmlの <coverage> or <source> or <filter> ...?

	return &selectCoverageFilesView{
		choices:  paths,
		selected: make(map[int]struct{}),
	}, nil
}

func (t *selectCoverageFilesView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		// case "ctrl+c", "q":
		// 	return t, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if t.cursor > 0 {
				t.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if t.cursor < len(t.choices)-1 {
				t.cursor++
			}

		// spaceで選択・選択解除
		case " ":
			_, ok := t.selected[t.cursor]
			if ok {
				delete(t.selected, t.cursor)
			} else {
				t.selected[t.cursor] = struct{}{}
			}

		// TODO: Enterで次に行きたい
		case "enter":
			// choiced := []string{}
			// for i := range t.selected {
			// 	choiced = append(choiced, t.choices[i])
			// }

			// fmt.Printf("Result2: %+v\n", choiced)

			m.quitting = true
			return m, tea.Quit
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (t *selectCoverageFilesView) view() string {
	// The header
	s := "Select target files you want to coverage (press Space)\n\n"

	// Iterate over our choices
	for i, choice := range t.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if t.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := t.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
