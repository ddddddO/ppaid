package main

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

const EXEC_PHPUNIT = `php -d pcov.directory=%s vendor/bin/phpunit --testsuite "%s" --configuration %s --coverage-html coverage`

type phpnuitFinishedMsg struct {
	err error
}

type cmdPHPUnit struct {
	cmd *exec.Cmd
}

const outputCoverageDir = "coverage-ppaid"

func (c *cmdPHPUnit) build(targetCoverageDir string, testSuiteName string, configFile string) {
	c.cmd = exec.Command("php", []string{
		"-d",
		fmt.Sprintf("pcov.directory=%s", targetCoverageDir), // ここカンマ区切りで複数指定可能のようだけど、どうもうまくいってない。なので、最大公約数的なパスを一旦指定しておく
		"vendor/bin/phpunit",
		"--testsuite",
		testSuiteName,
		"--configuration",
		configFile,
		"--coverage-html",
		outputCoverageDir,
	}...)
}

func (c *cmdPHPUnit) rawCmd() string {
	return c.cmd.String()
}

func (c *cmdPHPUnit) command() tea.Cmd {
	return tea.ExecProcess(c.cmd, func(err error) tea.Msg {
		return phpnuitFinishedMsg{err}
	})
}
