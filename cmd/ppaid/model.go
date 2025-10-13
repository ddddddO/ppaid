package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ViewOfSelectTestFiles = iota
	ViewOfSelectCoverageFiles
)

type model struct {
	currentView int
	quitting    bool

	selectTestFilesView     *selectTestFilesView
	selectCoverageFilesView *selectCoverageFilesView
}

func initialModel() (model, error) {
	tfv, err := newSelectTestFilesView()
	if err != nil {
		return model{}, err
	}

	cfv, err := newSelectCoverageFilesView()
	if err != nil {
		return model{}, err
	}

	return model{
		currentView: ViewOfSelectTestFiles,

		selectTestFilesView:     tfv,
		selectCoverageFilesView: cfv,
	}, nil
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	switch m.currentView {
	case ViewOfSelectTestFiles:
		return m.selectTestFilesView.update(msg, m)
	case ViewOfSelectCoverageFiles:
		return m.selectCoverageFilesView.update(msg, m)
	default:
		return m, nil
	}
}

func (m model) View() string {
	// 最終結果出力
	if m.quitting {
		var sb strings.Builder
		sb.WriteString("Result:\n\n")

		sb.WriteString("Selected test files:\n")
		if len(m.selectTestFilesView.selected) == 0 {
			sb.WriteString("  (no selected))\n")
		} else {
			for choice := range m.selectTestFilesView.selected {
				sb.WriteString(fmt.Sprintf("  - %s\n", choice))
			}
		}

		sb.WriteString("\nSelected coverage target:\n")
		if len(m.selectCoverageFilesView.selected) == 0 {
			sb.WriteString("  (no selected)\n")
		} else {
			for choice := range m.selectCoverageFilesView.selected {
				sb.WriteString(fmt.Sprintf("  - %s\n", choice))
			}
		}

		return sb.String()
	}

	switch m.currentView {
	case ViewOfSelectTestFiles:
		return m.selectTestFilesView.view()
	case ViewOfSelectCoverageFiles:
		return m.selectCoverageFilesView.view()
	default:
		return "unknown view"
	}
}
