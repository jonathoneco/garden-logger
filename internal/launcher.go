package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func browse(dirPath string) error {
	dirInfo, err := buildDirInfo(dirPath)
	if err != nil {
		return err
	}

	Logger.Debug("Built Status Info",
		"Absolute Path", dirInfo.AbsolutePath,
		"Indexing Strategy", dirInfo.IndexingStrategy,
		"Relative Path", dirInfo.RelativePath,
	)

	state := &MenuState{Mode: ModeBrowse, DirInfo: dirInfo}

	for {
		Logger.Debug("Browsing", "Current Path", dirInfo.RelativePath, "Mode", state.Mode)
		choice, err := state.launchMenu()
		Logger.Debug("Selection", "Choice", choice)
		if err != nil {
			return err
		}

		err = state.handleChoice(choice)
		if err != nil {
			return err
		}
	}
}

func launchDir(dirPath string) error {
	var sessionName string
	if dirPath == "" {
		sessionName = "The Garden Log"
	} else {
		sessionName = filepath.Base(dirPath)
	}

	fullPath := filepath.Join(rootDir, dirPath)

	cmd := exec.Command("kitty", "-e", "tmux", "new-session", "-s", sessionName, "-c", fullPath, "nvim .")

	cmd.Dir = fullPath

	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to launch kitty: %w", err)
	}

	os.Exit(0)
	return nil
}

func launchNote(filePath string) error {
	cmd := exec.Command("kitty", "-e", "nvim", filePath)

	cmd.Dir = rootDir

	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to launch kitty: %w", err)
	}

	os.Exit(0)
	return nil
}

func (ms *MenuState) launchMenu() (string, error) {
	items, err := ms.getMenuItems()
	if err != nil {
		return "", err
	}

	args := []string{"notes", "-dmenu", "-l", "10", "-i", "-p", ms.getPrompt()}

	if ms.DirInfo != nil {
		statusMsg := ms.DirInfo.formatStatusMessage()
		args = append(args, "-mesg", statusMsg)
	}

	cmd := exec.Command("rofi-launcher", args...)
	// Logger.Debug("Launching Menu", "cmd", cmd.String())
	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("[ERROR] Error getting output when launching rofi: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
