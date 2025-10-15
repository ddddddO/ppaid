package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sahilm/fuzzy"
)

type selectTestFilesView struct {
	height int

	choices         []string      // 選択肢のリスト
	filteredChoices []fuzzy.Match // 絞り込まれた選択肢のリスト

	cursor      int
	searchInput textinput.Model

	selected map[string]struct{}
}

func newSelectTestFilesView() (*selectTestFilesView, error) {
	paths, err := getPHPTestFilePaths()
	if err != nil {
		return nil, err
	}

	ti := textinput.New()
	ti.Placeholder = "Filter test files..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return &selectTestFilesView{
		height:          4,
		choices:         paths,
		filteredChoices: fuzzy.Find("", paths),
		selected:        make(map[string]struct{}),
		searchInput:     ti,
	}, nil
}

func (t *selectTestFilesView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if t.cursor > 0 {
				t.cursor--
			}

			t.searchInput, cmd = t.searchInput.Update(msg)
			return m, cmd

		case "down", "j":
			// 検索が空っぽの時
			if len(strings.TrimSpace(t.searchInput.Value())) == 0 {
				if t.cursor < len(t.choices)-1 {
					t.cursor++
				}
			}

			if t.cursor < len(t.filteredChoices)-1 {
				t.cursor++
			}

			t.searchInput, cmd = t.searchInput.Update(msg)
			return m, cmd

		// spaceで選択・選択解除
		case " ":
			// 検索が空っぽの時
			if len(strings.TrimSpace(t.searchInput.Value())) == 0 {
				choice := t.choices[t.cursor]
				_, ok := t.selected[choice]
				if ok {
					delete(t.selected, choice)
				} else {
					t.selected[choice] = struct{}{}
				}
			}

			if len(t.filteredChoices) > 0 {
				choice := t.filteredChoices[t.cursor].Str
				if _, ok := t.selected[choice]; ok {
					delete(t.selected, choice)
				} else {
					t.selected[choice] = struct{}{}
				}
			}

			// スペースで選択するため、検索窓への入力でスペースは許容しない
			// t.searchInput, cmd = t.searchInput.Update(msg)
			return m, nil

		case "enter":
			m.currentView = ViewOfSelectCoverageFiles
			return m, nil

		default:
			t.searchInput, cmd = t.searchInput.Update(msg)
			t.filteredChoices = fuzzy.Find(t.searchInput.Value(), t.choices)
			t.cursor = 0
			return m, cmd
		}
	}

	t.searchInput, cmd = t.searchInput.Update(msg)

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

func (t *selectTestFilesView) view() string {
	var sb strings.Builder
	sb.WriteString("Select the test files you want to run (press Space)\n\n")
	sb.WriteString(t.searchInput.View())
	sb.WriteString("\n\n")

	// t.filteredChoices は、fuzzy.Findの第一引数が空文字だとマッチしない。検索文字列が空なら全権表示するようにするため以下の処理
	matchs := []string{}
	if len(strings.TrimSpace(t.searchInput.Value())) == 0 {
		matchs = t.choices
	} else {
		for _, m := range t.filteredChoices {
			matchs = append(matchs, m.Str)
		}
	}

	// マイナスしてるのは、パス一覧を除いた高さを一旦決め打ちした数
	height := min(len(matchs), t.height-7)
	height = max(0, height) // 起動時、heightがマイナス値になることあってパニックになるから

	for i, match := range matchs[:height] {
		cursor := " " // no cursor
		if t.cursor == i {
			cursor = ">" // cursor!
		}

		checked := " " // not selected
		if _, ok := t.selected[match]; ok {
			checked = "x" // selected!
		}

		// Render the row
		sb.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, match))
	}
	if len(matchs) > height {
		sb.WriteString(fmt.Sprintf("  ... %d more", len(matchs)-height))
	}

	// The footer
	sb.WriteString("\n(↑↓: move, Space: select, Enter: next, q: quit)\n")

	// Send the UI for rendering
	return sb.String()
}
