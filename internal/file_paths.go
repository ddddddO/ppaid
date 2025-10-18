package internal

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ddddddO/ppaid/internal/command"
)

func GetPHPTestFilePaths() ([]string, error) {
	return GetPHPFilePaths("./tests", []string{"vendor", ".git"})
}

func GetPHPCodeFilePaths() ([]string, error) {
	paths := []string{}
	if ps, err := GetPHPFilePaths("./src", []string{"vendor", ".git"}); err == nil {
		paths = append(paths, ps...)
	}
	if ps, err := GetPHPFilePaths("./app", []string{"vendor", ".git"}); err == nil {
		paths = append(paths, ps...)
	}

	return paths, nil
}

func GetPHPFilePaths(targetDir string, ignore []string) ([]string, error) {
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

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if d.IsDir() && slices.Contains(ignore, d.Name()) {
			return fs.SkipDir
		}

		if !d.IsDir() && filepath.Ext(relativePath) != ".php" {
			return nil
		}

		paths = append(paths, relativePath)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func GetCoveragedFilePaths(level int) ([]string, error) {
	ignore := []string{"_css", "_icons", "_js"}

	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	target := filepath.Join(root, command.OutputCoverageDir)

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
