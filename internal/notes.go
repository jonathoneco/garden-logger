package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Note struct {
	Name string
	Path string
}

func createNewNote(dirPath string) error {
	nameCmd := exec.Command("dmenu", "-c", "-sb", "#a3be8c", "-nf", "#d8dee9", "-p", "Enter a name: ")
	nameCmd.Stdin = strings.NewReader("")

	output, err := nameCmd.Output()
	name := strings.TrimSpace(string(output))

	if err != nil || name == "" {
		name = time.Now().Format("2001-01-01")
	}

	var notePath string
	if dirPath == "" {
		notePath = "1 Inbox/"
	} else {
		notePath = dirPath
	}

	note := Note{
		Name: name,
	}

	filePath, err := writeNote(notePath, note)
	if err != nil {
		return err
	}

	return launchNote(filePath)
}

func writeNote(dirPath string, note Note) (string, error) {
	targetDir := filepath.Join(rootDir, dirPath)

	// if err := os.MkdirAll(targetDir, 0755); err != nil {
	// 	return "", err
	// }

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

	if _, err := file.WriteString(fmt.Sprintf("# %s\n\n", note.Name)); err != nil {
		return "", fmt.Errorf("failed to write title: %w", err)
	}

	return fullPath, nil
}
