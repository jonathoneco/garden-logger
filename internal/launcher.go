package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func launchMenu(items []string, prompt string, colors ...string) (string, error) {
	args := []string{"-c", "-l", "10", "-i", "-p", prompt}

	if len(colors) >= 2 {
		args = append(args, "-sb", colors[0], "-nf", colors[1])
	}

	cmd := exec.Command("dmenu", args...)
	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func launchNote(filePath string) error {
	cmd := exec.Command("kitty", "-e", "nvim", filePath)

	cmd.Dir = rootDir

	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to launch kitty: %w", err)
	}

	return nil
}
