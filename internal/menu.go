package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MenuMode int

const (
	ModeBrowse MenuMode = iota
	ModeNew
	ModeIndexing
	ModeNewNote
)

type MenuState struct {
	Mode    MenuMode // "browse", "new", "indexing"
	DirInfo *DirInfo
}

func formatOption(text string, selected bool) string {
	if selected {
		return text + " ✓"
	}
	return text
}

func (ms *MenuState) getMenuItems() ([]string, error) {
	switch ms.Mode {
	case ModeNew:
		return []string{MenuNewNote, MenuNewDirectory, MenuNewNoteFromTemplate, MenuBack}, nil
	case ModeIndexing:
		current := ms.DirInfo.IndexingStrategy
		return []string{
			formatOption(MenuIndexNumeric, current == IndexStrategyNumeric),
			formatOption(MenuIndexNone, current == IndexStrategyNone),
			MenuBack,
		}, nil
	case ModeBrowse: // browse
		entries, err := ms.DirInfo.listEntries()

		if err != nil {
			return entries, err
		}

		items := []string{MenuNew, MenuSettings}
		items = append(items, entries...)
		if ms.DirInfo.RelativePath != "" {
			items = append(items, MenuBack)
		}
		items = append(items, MenuOpenCurrentFolder)
		return items, nil
	default:
		return []string{}, nil
	}
}

func (ms *MenuState) getPrompt() string {
	switch ms.Mode {
	case ModeNew:
		return "New: "
	case ModeNewNote:
		return "Enter a name: "
	case ModeIndexing:
		return "Indexing: "
	default:
		return "Browse: "
	}
}

func (ms *MenuState) navigateTo(dirPath string) error {
	Logger.Debug("Navigating to", "dirPath", dirPath)
	parentDirInfo, err := buildDirInfo(dirPath)

	if err != nil {
		return err
	}

	ms.DirInfo = parentDirInfo
	return nil
}

func (ms *MenuState) navigateToParent() error {
	parentFolder := filepath.Dir(strings.TrimSuffix(ms.DirInfo.RelativePath, "/"))
	if parentFolder == "." {
		parentFolder = ""
	}
	return ms.navigateTo(parentFolder)
}

func (ms *MenuState) handleBrowseChoice(choice string) error {
	switch choice {
	case MenuNew:
		ms.Mode = ModeNew
	case MenuSettings:
		ms.Mode = ModeIndexing
	case MenuBack:
		err := ms.navigateToParent()
		if err != nil {
			return err
		}
	case MenuOpenCurrentFolder:
		return launchDir(ms.DirInfo.RelativePath)
	default:
		if strings.HasSuffix(choice, "/") {
			newDirPath := filepath.Join(ms.DirInfo.RelativePath, choice)
			ms.navigateTo(newDirPath)
			return nil
		}

		if strings.HasSuffix(strings.ToLower(choice), ".md") {
			fullFilePath := filepath.Join(ms.DirInfo.RelativePath, choice)
			return launchNote(fullFilePath)
		}

		return fmt.Errorf("[ERROR] unexpected choice: %s", choice)
	}

	return nil
}

func (ms *MenuState) handleNewChoice(choice string) error {
	switch choice {
	case MenuNewNote:
		ms.Mode = ModeNewNote
	case MenuNewDirectory:
	case MenuNewNoteFromTemplate:
	case MenuBack:
		ms.Mode = ModeBrowse
	}
	return nil
}

func (ms *MenuState) handleIndexingChoice(choice string) error {
	switch choice {
	case MenuIndexNumeric:
		if ms.DirInfo.IndexingStrategy != IndexStrategyNumeric {
			config := &IndexConfig{
				Strategy: IndexStrategyNumeric,
				NumericConfig: &NumericConfig{
					DirPriority: true, // Default to directory priority
				},
			}
			if err := writeIndexConfig(ms.DirInfo.AbsolutePath, config); err != nil {
				return fmt.Errorf("failed to set numeric indexing: %w", err)
			}
			if err := applyNumericIndexing(ms.DirInfo.AbsolutePath, true); err != nil {
				return fmt.Errorf("failed to apply numeric indexing: %w", err)
			}
			ms.DirInfo.IndexingStrategy = IndexStrategyNumeric
		}
	case MenuIndexNone:
		if ms.DirInfo.IndexingStrategy != IndexStrategyNone {
			if err := removeIndexing(ms.DirInfo.AbsolutePath); err != nil {
				return fmt.Errorf("failed to remove indexing: %w", err)
			}
			config := &IndexConfig{Strategy: IndexStrategyNone}
			if err := writeIndexConfig(ms.DirInfo.AbsolutePath, config); err != nil {
				return fmt.Errorf("failed to clear indexing config: %w", err)
			}
			ms.DirInfo.IndexingStrategy = IndexStrategyNone
		}
	}
	ms.Mode = ModeBrowse
	return nil
}

func (ms *MenuState) handleNewNote(choice string) error {
	name := choice
	if name == "" {
		name = time.Now().Format("2006-01-02")
	}

	filePath, err := ms.writeNote(name)
	if err != nil {
		return err
	}

	return launchNote(filePath)
}
func (ms *MenuState) writeNote(name string) (string, error) {
	var targetDir string

	if ms.DirInfo.RelativePath == "" {
		targetDir = filepath.Join(rootDir, inboxDir)
	} else {
		targetDir = ms.DirInfo.AbsolutePath
	}
	Logger.Debug("Writing Note", "targetDir", targetDir, "relativePath", ms.DirInfo.RelativePath)

	config, err := getIndexConfig(targetDir)
	if err != nil {
		return "", err
	}

	var filename string
	if config.Strategy == IndexStrategyNumeric {
		nextIndex, err := findNextIndex(targetDir)
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

func (ms *MenuState) handleChoice(choice string) error {
	var err error = nil
	switch ms.Mode {
	case ModeBrowse:
		err = ms.handleBrowseChoice(choice)
	case ModeNew:
		err = ms.handleNewChoice(choice)
	case ModeIndexing:
		err = ms.handleIndexingChoice(choice)
	case ModeNewNote:
		err = ms.handleNewNote(choice)
	}

	return err
}
