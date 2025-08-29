package internal

import (
	"fmt"
	"log/slog"
	"path/filepath"
)

type Mode int

const (
	ModeBrowse Mode = iota
	ModeNew
	ModeNewNote
	ModeNewTemplatedNote
	ModeSettings
)

type MenuState struct {
	Dir       *Directory
	Mode      Mode
	Selection string
}

func (menu *MenuState) formatStatusMessage() string {
	path := menu.Dir.Path
	return fmt.Sprintf("Path: %s \nIndexing: %s", path, "TEMP")
}

func InitMenuState() (*MenuState, error) {
	slog.Info("Initializing menu state")

	dir, err := LoadDirectory("")
	if err != nil {
		return nil, err
	}

	menuState := &MenuState{dir, ModeBrowse, ""}
	slog.Info("Menu state initialized successfully", "initialMode", ModeBrowse, "rootEntries", len(dir.Entries))
	return menuState, nil
}

func (menu *MenuState) navigateTo(dirPath string) error {
	slog.Info("Navigating to directory", "fromPath", menu.Dir.Path, "toPath", dirPath)

	dir, err := LoadDirectory(dirPath)
	if err != nil {
		return err
	}

	menu.Dir = dir
	menu.Selection = ""
	slog.Info("Navigation completed successfully", "newPath", dirPath, "entryCount", len(dir.Entries))
	return nil
}

func (menu *MenuState) navigateToParent() error {
	slog.Debug("Attempting to navigate to parent directory", "currentPath", menu.Dir.Path)

	if menu.Dir.Path == "" {
		return fmt.Errorf("already at root directory")
	}

	parentPath := filepath.Dir(menu.Dir.Path)
	if parentPath == "." {
		parentPath = ""
	}
	slog.Debug("Parent path resolved", "parentPath", parentPath)

	return menu.navigateTo(parentPath)
}

func (menu *MenuState) getPrompt() string {
	switch menu.Mode {
	case ModeNew:
		return "New: "
	case ModeNewNote:
		return "Enter a name: "
	case ModeSettings:
		return "Indexing: "
	default:
		return "Browse: "
	}
}

func (menu *MenuState) handleChoice(choice string) error {
	var err error = nil
	switch menu.Mode {
	case ModeBrowse:
		err = menu.handleBrowseChoice(choice)
	case ModeNew:
		err = menu.handleNewChoice(choice)
	// case ModeSettings:
	// 	err = menu.handleSettingsChoice(choice)
	case ModeNewNote:
		err = menu.handleNewNoteChoice(choice)
	}

	return err
}

func (menu *MenuState) getMenuItems() ([]string, error) {
	switch menu.Mode {
	case ModeNew:
		return getNewMenuItems()
	case ModeSettings:
		return menu.getSettingsMenuItems()
	case ModeBrowse: // browse
		return menu.getBrowseMenuItems()
	default:
		return nil, nil
	}
}

func Browse() error {
	menu, err := InitMenuState()
	if err != nil {
		return err
	}

	for {
		slog.Debug("Browsing", "Current Path", menu.Dir.Path, "Mode", menu.Mode)

		choice, err := menu.launchMenu()
		slog.Debug("Selection", "Choice", choice)
		if err != nil {
			return err
		}

		// Skip handling if choice is empty (file movement operations return empty string)
		if choice == "" {
			continue
		}

		err = menu.handleChoice(choice)
		if err != nil {
			return err
		}
	}
}
