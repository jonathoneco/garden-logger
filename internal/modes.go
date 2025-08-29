package internal

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// Browse Mode

func (menu *MenuState) getBrowseMenuItems() ([]string, error) {
	items := []string{MenuNew, MenuSettings}

	items = append(items, menu.Dir.ListEntries()...)
	if menu.Dir.Path != "" {
		items = append(items, MenuBack)
	}
	items = append(items, MenuOpenCurrentFolder)
	return items, nil
}

func (menu *MenuState) handleBrowseChoice(choice string) error {
	switch choice {
	case MenuNew:
		menu.Mode = ModeNew
	case MenuSettings:
		menu.Mode = ModeSettings
	case MenuBack:
		err := menu.navigateToParent()
		if err != nil {
			return err
		}
	case MenuOpenCurrentFolder:
		return launchDir(menu.Dir.Path)
	default:
		if strings.HasSuffix(choice, "/") {
			newDirPath := filepath.Join(menu.Dir.Path, choice)
			err := menu.navigateTo(newDirPath)
			if err != nil {
				return err
			}
			return nil
		}

		if strings.HasSuffix(strings.ToLower(choice), ".md") {
			fullFilePath := filepath.Join(menu.Dir.Path, choice)
			return launchNote(fullFilePath)
		}

		return fmt.Errorf("[ERROR] unexpected choice: %s", choice)
	}

	return nil
}

// New Mode

func getNewMenuItems() ([]string, error) {
	return []string{MenuNewNote, MenuNewDirectory, MenuNewNoteFromTemplate, MenuBack}, nil
}

func (menu *MenuState) handleNewChoice(choice string) error {
	switch choice {
	case MenuNewNote:
		menu.Mode = ModeNewNote
	case MenuNewDirectory:
	case MenuNewNoteFromTemplate:
	case MenuBack:
		menu.Mode = ModeBrowse
	}
	return nil
}

// New Note Mode

func getNewNoteMenuItems() ([]string, error) {
	return nil, nil
}

func (menu *MenuState) handleNewNoteChoice(choice string) error {
	name := choice
	if name == "" {
		name = time.Now().Format("2006-01-02")
	}

	entry := &Entry{
		Name:      name,
		NoteIndex: menu.Dir.NewFileIndex(),
		Ext:       ".md",
		IsDir:     false,
	}

	filePath, err := menu.writeNote(entry)
	if err != nil {
		return err
	}

	return launchNote(filePath)
}

// Settings Mode

func formatSelectedOption(text string, selected bool) string {
	if selected {
		return text + " ✓"
	}
	return text
}

func (menu *MenuState) getSettingsMenuItems() ([]string, error) {
	menuItems := []string{
		formatSelectedOption(MenuIndexSetting, menu.Dir.IsIndexed),
		MenuBack,
	}

	return menuItems, nil
}

func (menu *MenuState) handleSettingsChoice(choice string) error {

	switch choice {
	case MenuIndexSetting:
		menu.Dir.ApplyNumericIndexing()
	case formatSelectedOption(MenuIndexSetting, true):
		menu.Dir.RemoveIndexing()
	}

	menu.Mode = ModeBrowse
	return nil
}
