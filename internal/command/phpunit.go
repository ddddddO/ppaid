package command

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type PHPUitFinishedMsg struct {
	err error
}

func (p PHPUitFinishedMsg) Err() error {
	return p.err
}

type CmdPHPUnit struct {
	CommandToSpecifyBeforePHPCommand string
	cmd                              *exec.Cmd
}

const OutputCoverageDir = "coverage-ppaid"

func (c *CmdPHPUnit) Build(targetCoverageDir string, testSuiteName string, configFile string) {
	if len(c.CommandToSpecifyBeforePHPCommand) == 0 {
		c.cmd = exec.Command("php", []string{
			"-d",
			fmt.Sprintf("pcov.directory=%s", targetCoverageDir), // ここカンマ区切りで複数指定可能のようだけど、どうもうまくいってない。なので、最大公約数的なパスを一旦指定しておく
			"vendor/bin/phpunit",
			"--testsuite",
			testSuiteName,
			"--configuration",
			configFile,
			"--coverage-html",
			OutputCoverageDir,
		}...)

		return
	}

	parsedCmd := strings.Split(c.CommandToSpecifyBeforePHPCommand, " ")
	if len(parsedCmd) == 1 {
		c.cmd = exec.Command(parsedCmd[0], []string{
			"php",
			"-d",
			fmt.Sprintf("pcov.directory=%s", targetCoverageDir),
			"vendor/bin/phpunit",
			"--testsuite",
			testSuiteName,
			"--configuration",
			configFile,
			"--coverage-html",
			OutputCoverageDir,
		}...)

		return
	}

	args := append(parsedCmd[1:], []string{
		"php",
		"-d",
		fmt.Sprintf("pcov.directory=%s", targetCoverageDir),
		"vendor/bin/phpunit",
		"--testsuite",
		testSuiteName,
		"--configuration",
		configFile,
		"--coverage-html",
		OutputCoverageDir,
	}...)
	c.cmd = exec.Command(parsedCmd[0], args...)
}

func (c *CmdPHPUnit) RawCmd() string {
	return c.cmd.String()
}

func (c *CmdPHPUnit) Command() tea.Cmd {
	return tea.ExecProcess(c.cmd, func(err error) tea.Msg {
		return PHPUitFinishedMsg{err}
	})
}
