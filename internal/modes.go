package internal

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Browse Mode

func (m *MenuState) getBrowseMenuItems() ([]string, error) {
	items := []string{MenuNew, MenuSettings}

	items = append(items, m.Dir.ListEntries()...)
	if m.Dir.Path != "" {
		items = append(items, MenuBack)
	}
	items = append(items, MenuOpenCurrentFolder)
	return items, nil
}

func (m *MenuState) handleBrowseChoice(choice string) error {
	switch choice {
	case MenuNew:
		m.Mode = ModeNew
	case MenuSettings:
		m.Mode = ModeSettings
	case MenuBack:
		err := m.navigateToParent()
		if err != nil {
			return err
		}
	case MenuOpenCurrentFolder:
		return m.notes.LaunchDirectoryEditor(m.Dir.Path)
	default:
		if strings.HasSuffix(choice, "/") {
			newDirPath := filepath.Join(m.Dir.Path, choice)
			err := m.navigateTo(newDirPath)
			if err != nil {
				return err
			}
			return nil
		}

		if strings.HasSuffix(strings.ToLower(choice), ".md") {
			fullFilePath := filepath.Join(m.Dir.Path, choice)
			return m.notes.LaunchNoteEditor(fullFilePath)
		}

		return fmt.Errorf("unexpected menu choice %q in %s mode", choice, m.Mode)
	}

	return nil
}

// New Mode

func getNewMenuItems() ([]string, error) {
	return []string{MenuNewNote, MenuNewDirectory, MenuNewNoteFromTemplate, MenuBack}, nil
}

func (m *MenuState) handleNewChoice(choice string) error {
	switch choice {
	case MenuNewNote:
		m.Mode = ModeNewNote
	case MenuNewDirectory:
	case MenuNewNoteFromTemplate:
	case MenuBack:
		m.Mode = ModeBrowse
	}
	return nil
}

// New Note Mode

func getNewNoteMenuItems() ([]string, error) {
	return nil, nil
}

func (m *MenuState) handleNewNoteChoice(choice string) error {
	filePath, err := m.notes.CreateNoteFromUserInput(m.Dir, choice)
	if err != nil {
		return err
	}

	return m.notes.LaunchNoteEditor(filePath)
}

// Settings Mode

func formatSelectedOption(text string, selected bool) string {
	if selected {
		return text + " ✓"
	}
	return text
}

func (m *MenuState) getSettingsMenuItems() ([]string, error) {
	menuItems := []string{
		formatSelectedOption(MenuIndexSetting, m.Dir.IsIndexed),
		MenuBack,
	}

	return menuItems, nil
}

func (m *MenuState) handleSettingsChoice(choice string) error {

	switch choice {
	case MenuIndexSetting:
		m.Dir.ApplyNumericIndexing()
	case formatSelectedOption(MenuIndexSetting, true):
		m.Dir.RemoveIndexing()
	}

	dir, err := m.notes.LoadDirectory(m.Dir.Path)
	if err != nil {
		return err
	}
	m.Dir = dir

	m.Mode = ModeBrowse
	return nil
}
