package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func sortEntriesDirFirst(dirEntries []os.DirEntry) []string {
	var dirs, files []string

	for _, entry := range dirEntries {
		name := entry.Name()

		if strings.HasPrefix(name, ".") {
			continue
		}

		if entry.IsDir() {
			dirs = append(dirs, name+"/")
		} else if filepath.Ext(entry.Name()) == ".md" {
			files = append(files, name)
		}
	}

	sort.Strings(dirs)
	sort.Strings(files)

	result := append(dirs, files...)
	return result
}

func sortEntries(dirEntries []os.DirEntry) []string {
	var entries []string

	for _, entry := range dirEntries {
		name := entry.Name()

		if strings.HasPrefix(name, ".") {
			continue
		}

		if entry.IsDir() {
			entries = append(entries, name+"/")
		} else if filepath.Ext(entry.Name()) == ".md" {
			entries = append(entries, name)
		}
	}

	sort.Strings(entries)
	return entries
}

func listEntries(dirPath string) ([]string, error) {
	fullPath := filepath.Join(rootDir, dirPath)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	// TODO: Currently sorts with directory priority,
	// eventually have this as a component of Indexing strategy
	result := sortEntries(entries)
	return result, nil
}

func browse(dirPath string) error {
	entries, err := listEntries(dirPath)
	if err != nil {
		return err
	}

	var menuItems []string
	menuItems = append(menuItems, "New")

	menuItems = append(menuItems, entries...)

	if dirPath != "" {
		menuItems = append(menuItems, "..")
	}

	choice, err := launchMenu(menuItems, "Choose note or create new: ")
	if err != nil {
		return err
	}

	switch choice {
	case "New":
		return createNewNote(dirPath)
	case "..":
		parentFolder := filepath.Dir(strings.TrimSuffix(dirPath, "/"))
		if parentFolder == "." {
			parentFolder = ""
		} else {
			parentFolder += "/"
		}
		return browse(parentFolder)
	default:
		return handleFileOrDirectory(choice, dirPath)
	}
}

func handleFileOrDirectory(choice, dirPath string) error {
	if strings.HasSuffix(choice, "/") {
		newDirPath := dirPath + choice
		return browse(newDirPath)
	}

	if strings.HasSuffix(strings.ToLower(choice), ".md") {
		fullFilePath := filepath.Join(dirPath, choice)
		return launchNote(fullFilePath)
	}

	return fmt.Errorf("unexpected choice: %s", choice)
}

func findNextIndex(targetDir string) (int, error) {
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 1, nil
		}
		return 0, err
	}

	var maxIndex int
	indexPattern := regexp.MustCompile(`^(\d+)\s`)

	for _, entry := range entries {
		if matches := indexPattern.FindStringSubmatch(entry.Name()); matches !=
			nil {
			if index, err := strconv.Atoi(matches[1]); err == nil && index >
				maxIndex {
				maxIndex = index
			}
		}
	}

	return maxIndex + 1, nil
}
