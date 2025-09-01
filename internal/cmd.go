package internal

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
	RofiExitCodeMoveDown = 10 // Ctrl+Alt+J
	RofiExitCodeMoveUp   = 11 // Ctrl+Alt+K
	RofiExitCodeDelete   = 12 // Ctrl+Alt+K
)

func (m *MenuState) launchMenu() (string, error) {
	items, err := m.getMenuItems()
	if err != nil {
		return "", err
	}

	args := []string{"notes", "-dmenu", "-l", "10", "-i", "-p", m.getPrompt()}

	if m.Selection != "" {
		for i, item := range items {
			if item == m.Selection {
				args = append(args, "-selected-row", strconv.Itoa(i))
				break
			}
		}
	}

	// if m.nav.CurrentDirectory() != nil {
	// 	statusMsg := m.formatStatusMessage()
	// 	args = append(args, "-mesg", statusMsg)
	// }

	cmd := exec.Command("rofi-launcher", args...)
	menuInput := strings.Join(items, "\n")
	cmd.Stdin = strings.NewReader(menuInput)

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			selection := strings.TrimSpace(string(output))
			entry := m.nav.CurrentDirectory().GetEntryByFilename(selection)
			switch exitError.ExitCode() {
			case RofiExitCodeMoveDown:
				m.nav.CurrentDirectory().MoveEntryDown(entry)
				m.Selection = entry.String()
				return "", nil
			case RofiExitCodeMoveUp:
				m.nav.CurrentDirectory().MoveEntryUp(entry)
				m.Selection = entry.String()
				return "", nil
			case RofiExitCodeDelete:
				m.nav.CurrentDirectory().DeleteEntry(entry)
				return "", nil
			}
		}
		return "", fmt.Errorf("rofi command failed in %s mode: %w", m.Mode, err)
	}
	selection := strings.TrimSpace(string(output))

	return selection, nil
}
