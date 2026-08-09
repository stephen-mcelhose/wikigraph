package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// applySed pipes html through sed with the given expressions (each passed as a
// separate -e argument). Requires sed to be available on PATH.
func applySed(html string, exprs []string) (string, error) {
	args := make([]string, 0, len(exprs)*2)
	for _, e := range exprs {
		args = append(args, "-e", e)
	}
	cmd := exec.Command("sed", args...)
	cmd.Stdin = strings.NewReader(html)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("sed: %w", err)
	}
	return string(out), nil
}
