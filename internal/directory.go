package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DirInfo

type DirInfo struct {
	AbsolutePath     string
	IndexingStrategy string
	RelativePath     string
}

func buildDirInfo(dirPath string) (*DirInfo, error) {
	absPath := filepath.Join(rootDir, dirPath)

	config, err := getIndexConfig(absPath)
	if err != nil {
		return nil, err
	}

	relativePath := dirPath

	indexingStrategy := config.Strategy.String()

	return &DirInfo{
		AbsolutePath:     absPath,
		IndexingStrategy: indexingStrategy,
		RelativePath:     relativePath,
	}, nil
}

func (dirInfo *DirInfo) listEntries() ([]string, error) {
	fullPath := dirInfo.AbsolutePath
	dirEntries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error reading entries for path %s: %w", fullPath, err)

	}

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
	return entries, nil
}

func (dirInfo *DirInfo) formatStatusMessage() string {
	path := dirInfo.RelativePath
	return fmt.Sprintf("Path: %s \nIndexing: %s", path, dirInfo.IndexingStrategy)
}
