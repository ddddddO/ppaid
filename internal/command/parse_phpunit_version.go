package command

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// メジャーバージョンを返す
func ParsePHPUnitVersion() (int, error) {
	cmd := exec.Command("php", []string{
		"vendor/bin/phpunit",
		"--version",
	}...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	return validatePHPUnitVersion(string(stdoutStderr))
}

func validatePHPUnitVersion(s string) (int, error) {
	parsed := strings.Split(s, " ")
	fmt.Println(parsed)
	if len(parsed) < 2 {
		return 0, errors.New("failed to parse PHPUnit version (1)")
	}
	if parsed[0] != "PHPUnit" {
		return 0, errors.New("failed to parse PHPUnit version (2)")
	}
	rawVersion := parsed[1]
	splitedVersion := strings.Split(rawVersion, ".")
	if len(splitedVersion) < 1 {
		return 0, errors.New("failed to parse PHPUnit version (3)")
	}
	majorVersion, err := strconv.Atoi(splitedVersion[0])
	if err != nil {
		return 0, errors.New("failed to parse PHPUnit version (4)")
	}
	return majorVersion, nil
}
