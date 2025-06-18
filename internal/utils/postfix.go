package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func PostfixHostname() (string, error) {
	cmd := exec.Command("postconf", "-h", "myhostname")

	stderr := strings.Builder{}
	stdout := strings.Builder{}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("unable to get postfix hostname: %w", err)
	}
	if stderr.String() != "" {
		return "", fmt.Errorf("unexpected output: %s", strings.TrimSpace(stderr.String()))

	}
	return strings.TrimSpace(stdout.String()), nil
}
