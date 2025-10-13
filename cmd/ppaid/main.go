package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m, err := initialModel()
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

func getPHPTestFilePaths() ([]string, error) {
	return getPHPFilePaths("./tests", []string{"vendor", ".git"})
}

func getPHPCodeFilePaths() ([]string, error) {
	paths := []string{}
	if ps, err := getPHPFilePaths("./src", []string{"vendor", ".git"}); err == nil {
		paths = append(paths, ps...)
	}
	if ps, err := getPHPFilePaths("./app", []string{"vendor", ".git"}); err == nil {
		paths = append(paths, ps...)
	}

	return paths, nil
}

func getPHPFilePaths(targetDir string, ignore []string) ([]string, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	target := filepath.Join(root, targetDir)

	paths := []string{}
	err = filepath.WalkDir(target, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == root {
			return nil
		}

		relativePath, err := filepath.Rel(target, path)
		if err != nil {
			return err
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
