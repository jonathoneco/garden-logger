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

		// Skip handling if choice is empty (file movement operations return empty string)
		if choice == "" {
			continue
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
	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			switch exitError.ExitCode() {
			case 10: // kb-custom-1 (Ctrl+Alt+J) - move down
				return ms.handleMoveDown(string(output))
			case 11: // kb-custom-2 (Ctrl+Alt+K) - move up  
				return ms.handleMoveUp(string(output))
			}
		}
		return "", fmt.Errorf("[ERROR] Error getting output when launching rofi: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func (ms *MenuState) handleMoveDown(selectedItem string) (string, error) {
	if ms.Mode != ModeBrowse {
		return "", fmt.Errorf("file movement only available in browse mode")
	}
	
	config, err := getIndexConfig(ms.DirInfo.AbsolutePath)
	if err != nil {
		return "", err
	}
	
	if config.Strategy != IndexStrategyNumeric {
		return "", fmt.Errorf("file movement only available with numeric indexing")
	}

	selectedItem = strings.TrimSpace(selectedItem)
	if selectedItem == "" {
		return "", fmt.Errorf("no file selected for movement")
	}

	err = moveIndexedFileDown(ms.DirInfo.AbsolutePath, selectedItem)
	if err != nil {
		Logger.Error("Move down failed", "error", err)
		return "", err
	}
	
	// Refresh directory info to reflect changes
	ms.DirInfo, err = buildDirInfo(ms.DirInfo.RelativePath)
	if err != nil {
		return "", err
	}
	
	return "", nil // Return to menu
}

func (ms *MenuState) handleMoveUp(selectedItem string) (string, error) {
	if ms.Mode != ModeBrowse {
		return "", fmt.Errorf("file movement only available in browse mode")
	}
	
	config, err := getIndexConfig(ms.DirInfo.AbsolutePath)
	if err != nil {
		return "", err
	}
	
	if config.Strategy != IndexStrategyNumeric {
		return "", fmt.Errorf("file movement only available with numeric indexing")
	}

	selectedItem = strings.TrimSpace(selectedItem)
	if selectedItem == "" {
		return "", fmt.Errorf("no file selected for movement")
	}

	err = moveIndexedFileUp(ms.DirInfo.AbsolutePath, selectedItem)
	if err != nil {
		Logger.Error("Move up failed", "error", err)
		return "", err
	}
	
	// Refresh directory info to reflect changes
	ms.DirInfo, err = buildDirInfo(ms.DirInfo.RelativePath)
	if err != nil {
		return "", err
	}
	
	return "", nil // Return to menu
}
