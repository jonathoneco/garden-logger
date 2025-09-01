package internal

import (
	"fmt"
	"path/filepath"
)

// General

func (m *MenuState) handleFileSelection(choice string, onFileSelect func(string) error) error {
	switch choice {
	case MenuBack:
		return m.nav.NavigateToParent()
	default:
		entry := m.nav.CurrentDirectory().FindEntryFromFilename(choice)
		if entry == nil {
			return fmt.Errorf("entry not found: %q", choice)
		}

		fullPath := filepath.Join(m.nav.CurrentDirectory().Path, choice)
		if entry.IsDir {
			return m.nav.NavigateTo(fullPath)
		} else {
			return onFileSelect(fullPath)
		}
	}
}

// Browse Mode

func (m *MenuState) getBrowseMenuItems() ([]string, error) {
	items := []string{MenuNew, MenuSettings}

	items = append(items, m.getNavigationMenuItems()...)
	items = append(items, MenuOpenCurrentFolder)
	return items, nil
}

func (m *MenuState) handleBrowseChoice(choice string) error {
	switch choice {
	case MenuNew:
		m.Mode = ModeNew
		return nil
	case MenuSettings:
		m.Mode = ModeSettings
		return nil
	case MenuOpenCurrentFolder:
		return m.notes.LaunchDirectoryEditor(m.nav.CurrentDirectory().Path)
	}

	return m.handleFileSelection(choice, m.notes.LaunchNoteEditor)
}

// New Mode

func getNewMenuItems() ([]string, error) {
	return []string{MenuNewNote, MenuNewDirectory, MenuNewNoteFromTemplate, MenuBack}, nil
}

func (m *MenuState) handleNewChoice(choice string) error {
	switch choice {
	case MenuNewNote:
		if m.nav.CurrentDirectory().Path == "" {
			err := m.nav.NavigateTo(m.config.InboxDir)
			if err != nil {
				return err
			}
		}
		m.Mode = ModeNewNote
	case MenuNewDirectory:
		m.Mode = ModeNewDirectory
	case MenuNewNoteFromTemplate:
		if m.nav.CurrentDirectory().Path == "" {
			err := m.nav.NavigateTo(m.config.InboxDir)
			if err != nil {
				return err
			}
		}
		m.nav.Save()
		m.nav.NavigateTo(m.config.TemplateDir)
		m.Mode = ModePickTemplate
	case MenuBack:
		m.Mode = ModeBrowse
	}
	return nil
}

// New Note Mode

func (m *MenuState) handleNewEntry(choice string, isDir bool) error {
	var filePath string
	var err error

	templatePath, templateErr := m.nav.RestoreTemplate()
	if templateErr == nil {
		filePath, err = m.notes.CreateEntryFromTemplate(m.nav.CurrentDirectory(), choice, templatePath)
	} else {
		filePath, err = m.notes.CreateEntryFromUserInput(m.nav.CurrentDirectory(), choice, isDir)
	}

	if err != nil {
		return err
	}

	m.Mode = ModeBrowse

	if isDir {
		return m.nav.NavigateTo(m.nav.CurrentDirectory().Path)
	}

	return m.notes.LaunchNoteEditor(filePath)
}

// Settings Mode

func formatSelectedOption(text string, selected bool) string {
	if selected {
		return text + "   ✓"
	}
	return text
}

func (m *MenuState) getSettingsMenuItems() ([]string, error) {
	menuItems := []string{
		formatSelectedOption(MenuIndexSetting, m.nav.CurrentDirectory().IsIndexed),
		MenuBack,
	}

	return menuItems, nil
}

func (m *MenuState) handleSettingsChoice(choice string) error {
	currentDir := m.nav.CurrentDirectory()

	switch choice {
	case MenuIndexSetting:
		currentDir.ApplyNumericIndexing()
	case formatSelectedOption(MenuIndexSetting, true):
		currentDir.RemoveIndexing()
	}

	err := m.nav.NavigateTo(currentDir.Path)
	if err != nil {
		return err
	}

	m.Mode = ModeBrowse
	return nil
}

// Template Mode

func (m *MenuState) handleTemplateChoice(choice string) error {
	return m.handleFileSelection(choice, func(templatePath string) error {
		originalDir, err := m.nav.Restore()
		if err != nil {
			return err
		}

		err = m.nav.NavigateTo(originalDir.Path)
		if err != nil {
			return err
		}

		m.nav.SaveTemplate(templatePath)
		m.Mode = ModeNewNote
		return nil
	})
}
