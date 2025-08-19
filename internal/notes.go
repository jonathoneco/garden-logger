package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Note struct {
	Name string
	Path string
}

func createNewNote(statusInfo *StatusInfo) error {
	name, err := launchMenu([]string{}, "Enter a name: ", statusInfo)

	if err != nil {
		return err
	}

	if name == "" {
		name = time.Now().Format("2006-01-02")
	}

	note := Note{
		Name: name,
	}

	filePath, err := writeNote(statusInfo, note)
	if err != nil {
		return err
	}

	return launchNote(filePath)
}

func writeNote(statusInfo *StatusInfo, note Note) (string, error) {
	var targetDir string

	if statusInfo.RelativePath == "/" {
		targetDir = filepath.Join(rootDir, inboxDir)
	} else {
		targetDir = statusInfo.AbsolutePath
	}

	nextIndex, err := findNextIndex(targetDir)
	if err != nil {
		return "", err
	}

	indexedFilename := fmt.Sprintf("%d %s.md", nextIndex, note.Name)
	fullPath := filepath.Join(targetDir, indexedFilename)

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", fullPath, err)
	}
	defer file.Close()

	frontmatter := fmt.Sprintf("# %s\n\n", note.Name)

	if _, err := file.WriteString(frontmatter); err != nil {
		return "", fmt.Errorf("failed to write frontmatter: %w", err)
	}

	return fullPath, nil
}
