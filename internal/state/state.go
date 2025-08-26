package state

import (
	"fmt"
	"garden-logger/internal/indexing"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Types

type Kind int

const (
	KindFile Kind = iota
	KindDir
)

type Entry struct {
	Index string
	Name  string
	Ext   string
	IsDir bool
}

type DirState struct {
	DirPath  string
	AbsPath  string
	Strategy *indexing.IndexConfig
	Entries  []*Entry
}

type Mode int

const (
	ModeBrowse Kind = iota
	ModeNew
	ModeNewNote
	ModeSettings
)

type MenuState struct {
	DirState string
	Mode     Mode
}

// Methods

func LoadDirState(dirPath string) (*DirState, error) {
	absPath := filepath.Join(RootDir, dirPath)

	config, err := indexing.GetIndexConfig(absPath)
	if err != nil {
		return nil, err
	}

	entries, err := ListEntries(dirPath)
	if err != nil {
		return nil, err
	}

	return &DirState{
		DirPath:  dirPath,
		AbsPath:  absPath,
		Strategy: config,
		Entries:  entries,
	}, err
}

func ListEntries(dirPath string) ([]*Entry, error) {
	absPath := filepath.Join(RootDir, dirPath)

	dirEntries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, &StateError{Op: "ListEntries", Path: dirPath, Err: err}
	}

	var entries []*Entry

	for _, entry := range dirEntries {
		name := entry.Name()

		if strings.HasPrefix(name, ".") {
			continue
		}

		entry, err := LoadEntry(name, entry.IsDir())

		if err != nil {
			return entries, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// func LoadEntry(name string, isDir bool) (*Entry, error) {
func LoadEntry(name string, isDir bool) {
	entry := &Entry{}

	re := regexp.MustCompile(`^(\d{2})\.\s*(.+)$`)
	matches := re.FindStringSubmatch(name)
	if len(matches) == 3 {
		return &Entry{matches[1], matches[2], true}

	}
	return 0, name, false

	return &Entry{}, nil
}
