package internal

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	// "log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

func (menu *MenuState) writeNote(entry *Entry) (string, error) {
	slog.Info("Starting note creation", "noteName", entry.Name, "currentDir", menu.Dir.Path)

	var targetDir string

	if menu.Dir.Path == "" {
		targetDir = filepath.Join(RootDir, InboxDir)
		slog.Debug("Using inbox directory for note", "targetDir", targetDir)
	} else {
		targetDir = menu.Dir.AbsPath
		slog.Debug("Using current directory for note", "targetDir", targetDir)
	}

	fullPath := filepath.Join(targetDir, entry.String())

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Failed to create file %s: %w", fullPath, err)
	}
	defer file.Close()
	frontmatter := fmt.Sprintf("# %s\n\n", entry.Name)

	if _, err := file.WriteString(frontmatter); err != nil {
		return "", fmt.Errorf("[ERROR] Failed to write frontmatter: %w", err)
	}

	slog.Info("Note created successfully", "fullPath", fullPath)
	return fullPath, nil
}

func launchDir(dirPath string) error {
	slog.Info("Launching directory editor", "dirPath", dirPath)

	var sessionName string
	if dirPath == "" {
		sessionName = "The Garden Log"
	} else {
		sessionName = filepath.Base(dirPath)
	}

	fullPath := filepath.Join(RootDir, dirPath)

	cmd := exec.Command("kitty", "-e", "tmux", "new-session", "-s", sessionName, "-c", fullPath, "nvim .")

	cmd.Dir = fullPath
	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to launch kitty: %w", err)
	}

	slog.Info("Directory editor launched successfully")
	os.Exit(0)
	return nil
}

func launchNote(filePath string) error {
	slog.Info("Launching note editor", "filePath", filePath)

	cmd := exec.Command("kitty", "-e", "nvim", filePath)

	cmd.Dir = RootDir
	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to launch kitty: %w", err)
	}

	slog.Info("Note editor launched successfully")
	os.Exit(0)
	return nil
}

func (menu *MenuState) launchMenu() (string, error) {
	slog.Debug("Launching menu interface", "mode", menu.Mode, "currentDir", menu.Dir.Path)

	items, err := menu.getMenuItems()
	if err != nil {
		return "", err
	}

	args := []string{"notes", "-dmenu", "-l", "10", "-i", "-p", menu.getPrompt()}

	if menu.Selection != "" {
		for i, item := range items {
			if item == menu.Selection {
				args = append(args, "-selected-row", strconv.Itoa(i))
				break
			}
		}
	}

	if menu.Dir != nil {
		statusMsg := menu.formatStatusMessage()
		args = append(args, "-mesg", statusMsg)
	}

	cmd := exec.Command("rofi-launcher", args...)
	menuInput := strings.Join(items, "\n")
	cmd.Stdin = strings.NewReader(menuInput)

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			selection := strings.TrimSpace(string(output))
			entry := menu.Dir.FindEntryFromFilename(selection)
			switch exitError.ExitCode() {
			case 10: // kb-custom-1 (Ctrl+Alt+J) - move down
				slog.Info("Move down keybind triggered", "exitCode", 10, "output", selection)
				menu.Dir.MoveEntryDown(entry)
				menu.Selection = entry.String()
				return "", nil
			case 11: // kb-custom-2 (Ctrl+Alt+K) - move up
				slog.Info("Move up keybind triggered", "exitCode", 11, "output", selection)
				menu.Dir.MoveEntryUp(entry)
				menu.Selection = entry.String()
				return "", nil
			}
		}
		return "", fmt.Errorf("[ERROR] Error getting output when launching rofi: %w", err)
	}
	selection := strings.TrimSpace(string(output))

	slog.Info("Menu selection made", "selection", selection, "mode", menu.Mode)
	return selection, nil
}
