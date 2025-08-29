package internal

import (
	"fmt"
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
	config    *Config
	notes     *NotesService
}

func (m *MenuState) formatStatusMessage() string {
	path := m.Dir.Path
	return fmt.Sprintf("Path: %s \nIndexing: %s", path, "TEMP")
}

func InitMenuState() (*MenuState, error) {
	config, err := LoadConfig()

	if err != nil {
		return nil, err
	}

	notes := NewNotesService(config)
	menu := &MenuState{nil, ModeBrowse, "", config, notes}
	dir, err := notes.LoadDirectory("")
	if err != nil {
		return nil, err
	}
	menu.Dir = dir

	return menu, nil
}

func (m *MenuState) navigateTo(dirPath string) error {
	dir, err := m.notes.LoadDirectory(dirPath)
	if err != nil {
		return err
	}
	m.Dir = dir
	m.Selection = ""
	return nil
}

func (m *MenuState) navigateToParent() error {

	if m.Dir.Path == "" {
		return fmt.Errorf("already at root directory")
	}

	parentPath := filepath.Dir(m.Dir.Path)
	if parentPath == "." {
		parentPath = ""
	}

	return m.navigateTo(parentPath)
}

func (m *MenuState) getPrompt() string {
	switch m.Mode {
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

func (m *MenuState) handleChoice(choice string) error {
	var err error = nil
	switch m.Mode {
	case ModeBrowse:
		err = m.handleBrowseChoice(choice)
	case ModeNew:
		err = m.handleNewChoice(choice)
	case ModeSettings:
		err = m.handleSettingsChoice(choice)
	case ModeNewNote:
		err = m.handleNewNoteChoice(choice)
	}

	return err
}

func (m *MenuState) getMenuItems() ([]string, error) {
	switch m.Mode {
	case ModeNew:
		return getNewMenuItems()
	case ModeSettings:
		return m.getSettingsMenuItems()
	case ModeBrowse: // browse
		return m.getBrowseMenuItems()
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

		choice, err := menu.launchMenu()
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
