package menu

import (
	"fmt"
	"garden-logger/internal/config"
	"garden-logger/internal/indexing"
	"os"
	"path/filepath"
)

func (menuState *MenuState) writeNote(name string) (string, error) {
	var targetDir string

	if menuState.DirState.RelativePath == "" {
		targetDir = filepath.Join(config.RootDir, config.InboxDir)
	} else {
		targetDir = menuState.DirState.AbsolutePath
	}
	log.Debug("Writing Note", "targetDir", targetDir, "relativePath", menuState.DirState.RelativePath)

	config, err := indexing.GetIndexConfig(targetDir)
	if err != nil {
		return "", err
	}

	var filename string
	if config.Strategy == indexing.IndexStrategyNumeric {
		nextIndex, err := indexing.FindNextIndex(targetDir)
		if err != nil {
			return "", err
		}
		filename = fmt.Sprintf("%d - %s.md", nextIndex, name)
	} else {
		filename = fmt.Sprintf("%s.md", name)
	}

	fullPath := filepath.Join(targetDir, filename)

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Failed to create file %s: %w", fullPath, err)
	}
	defer file.Close()

	frontmatter := fmt.Sprintf("# %s\n\n", name)

	if _, err := file.WriteString(frontmatter); err != nil {
		return "", fmt.Errorf("[ERROR] Failed to write frontmatter: %w", err)
	}

	return fullPath, nil
}

func launchDir(dirPath string) error {
	var sessionName string
	if dirPath == "" {
		sessionName = "The Garden Log"
	} else {
		sessionName = filepath.Base(dirPath)
	}

	fullPath := filepath.Join(config.RootDir, dirPath)

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

	cmd.Dir = config.RootDir

	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to launch kitty: %w", err)
	}

	os.Exit(0)
	return nil
}

func (menuState *MenuState) launchMenu() (string, error) {
	items, err := menuState.getMenuItems()
	if err != nil {
		return "", err
	}

	args := []string{"notes", "-dmenu", "-l", "10", "-i", "-p", menuState.getPrompt()}

	if menuState.DirState != nil {
		statusMsg := menuState.DirState.formatStatusMessage()
		args = append(args, "-mesg", statusMsg)
	}

	cmd := exec.Command("rofi-launcher", args...)
	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			switch exitError.ExitCode() {
			case 10: // kb-custom-1 (Ctrl+Alt+J) - move down
				log.Debug("MoveDownKeyBind")
				return "", nil
				// return ms.handleMoveDown(string(output))
			case 11: // kb-custom-2 (Ctrl+Alt+K) - move up
				// return ms.handleMoveUp(string(output))
				log.Debug("MoveDownKeyBind")
				return "", nil
			}
		}
		return "", fmt.Errorf("[ERROR] Error getting output when launching rofi: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func (menuState *MenuState) handleMoveDown(selectedItem string) (string, error) {
	if menuState.Mode != ModeBrowse {
		return "", fmt.Errorf("file movement only available in browse mode")
	}

	indexConfig, err := indexing.GetIndexConfig(menuState.DirState.AbsolutePath)
	if err != nil {
		return "", err
	}

	if indexConfig.Strategy != indexing.IndexStrategyNumeric {
		return "", fmt.Errorf("file movement only available with numeric indexing")
	}

	selectedItem = strings.TrimSpace(selectedItem)
	if selectedItem == "" {
		return "", fmt.Errorf("no file selected for movement")
	}

	err = indexing.MoveIndexedFileDown(menuState.DirState.AbsolutePath, selectedItem)
	if err != nil {
		log.Error("Move down failed", "error", err)
		return "", err
	}

	// Refresh directory info to reflect changes
	menuState.DirState, err = buildDirInfo(menuState.DirState.RelativePath)
	if err != nil {
		return "", err
	}

	return "", nil // Return to menu
}

func (menuState *MenuState) handleMoveUp(selectedItem string) (string, error) {
	if menuState.Mode != ModeBrowse {
		return "", fmt.Errorf("file movement only available in browse mode")
	}

	indexConfig, err := indexing.GetIndexConfig(menuState.DirState.AbsolutePath)
	if err != nil {
		return "", err
	}

	if indexConfig.Strategy != indexing.IndexStrategyNumeric {
		return "", fmt.Errorf("file movement only available with numeric indexing")
	}

	selectedItem = strings.TrimSpace(selectedItem)
	if selectedItem == "" {
		return "", fmt.Errorf("no file selected for movement")
	}

	err = indexing.MoveIndexedFileUp(menuState.DirState.AbsolutePath, selectedItem)
	if err != nil {
		log.Error("Move up failed", "error", err)
		return "", err
	}

	// Refresh directory info to reflect changes
	menuState.DirState, err = buildDirInfo(menuState.DirState.RelativePath)
	if err != nil {
		return "", err
	}

	return "", nil // Return to menu
}
