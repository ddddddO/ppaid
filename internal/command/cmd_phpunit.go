package command

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type PHPUitFinishedMsg struct {
	err error
}

func (p PHPUitFinishedMsg) Err() error {
	return p.err
}

type CmdPHPUnit struct {
	cmd *exec.Cmd
}

const OutputCoverageDir = "coverage-ppaid"

func (c *CmdPHPUnit) Build(targetCoverageDir string, testSuiteName string, configFile string) {
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
}

func (c *CmdPHPUnit) RawCmd() string {
	return c.cmd.String()
}

func (c *CmdPHPUnit) Command() tea.Cmd {
	return tea.ExecProcess(c.cmd, func(err error) tea.Msg {
		return PHPUitFinishedMsg{err}
	})
}
