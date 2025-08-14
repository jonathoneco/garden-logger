package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Note struct {
	Name string
	Path string
}

func createNewNote(dirPath string) error {
	name, err := launchMenu([]string{}, "Enter a name: ")

	if err != nil {
		return err
	}

	if name == "" {
		name = time.Now().Format("2006-01-02")
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

func generateSemanticID(dirPath, noteName string) string {
	// Remove rootDir prefix and clean the path
	cleanPath := strings.TrimPrefix(dirPath, rootDir)
	cleanPath = strings.Trim(cleanPath, "/")

	// Split path and remove indexes from each segment
	var segments []string
	if cleanPath != "" {
		pathParts := strings.Split(cleanPath, "/")
		for _, part := range pathParts {
			// Remove leading index (e.g., "5 Resources" -> "Resources")
			cleaned := regexp.MustCompile(`^\d+\s+`).ReplaceAllString(part, "")
			if cleaned != "" {
				segments = append(segments, cleaned)
			}
		}
	}

	// Add note name (also strip index)
	cleanedNoteName := regexp.MustCompile(`^\d+\s+`).ReplaceAllString(noteName, "")
	segments = append(segments, cleanedNoteName)

	// Convert to lowercase, replace spaces with hyphens, join with underscores
	var cleanedSegments []string
	for _, segment := range segments {
		cleaned := strings.ToLower(segment)
		cleaned = strings.ReplaceAll(cleaned, " ", "-")
		cleanedSegments = append(cleanedSegments, cleaned)
	}

	return strings.Join(cleanedSegments, "_")
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

	// semanticID := generateSemanticID(dirPath, note.Name)

	// frontmatter := fmt.Sprintf("---\naliases: [%s]\n---\n# %s\n\n", semanticID, note.Name)
	frontmatter := fmt.Sprintf("# %s\n\n", note.Name)

	if _, err := file.WriteString(frontmatter); err != nil {
		return "", fmt.Errorf("failed to write frontmatter: %w", err)
	}

	return fullPath, nil
}
