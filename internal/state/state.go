package state

import (
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var Log *slog.Logger

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
	DirPath string
	AbsPath string
	// Strategy *indexing.IndexConfig
	Entries []*Entry
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

	// config, err := indexing.GetIndexConfig(absPath)
	// if err != nil {
	// 	return nil, err
	// }

	Log.Debug("Test")

	entries, err := ListEntries(dirPath)
	if err != nil {
		return nil, err
	}

	Log.Debug("Loaded Directory", "dirPath", dirPath)

	return &DirState{
		DirPath: dirPath,
		AbsPath: absPath,
		// Strategy: config,
		Entries: entries,
	}, err
}

func ListEntries(dirPath string) ([]*Entry, error) {
	absPath := filepath.Join(RootDir, dirPath)

	dirEntries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, &StateError{Op: "ListEntries", Path: dirPath, Err: err}
	}

	var entries []*Entry

	for _, dirEntry := range dirEntries {
		if strings.HasPrefix(dirEntry.Name(), ".") {
			continue
		}

		entry := LoadEntry(dirEntry)

		entries = append(entries, entry)
	}

	return entries, nil
}

func LoadEntry(dirEntry os.DirEntry) *Entry {
	index, name := parseEntryName(dirEntry.Name())
	isDir := dirEntry.IsDir()
	ext := "/"
	if !isDir {
		ext = filepath.Ext(dirEntry.Name())
	}

	entry := &Entry{index, name, ext, isDir}
	Log.Debug("Loaded Entry", "entry", entry)
	return entry
}

// Parses entry name and returns Index, CleanedName
func parseEntryName(name string) (string, string) {
	// This seems like a poor use of branching, is there a better way to handle the optional file extension with regex
	// re := regexp.MustCompile(`^(\d{2})\.\s*(.+)\.(.+)$|^(\d{2})\.\s*(.+)$`)
	re := regexp.MustCompile(`^(\d{2})\.\s*(.+)(\.?)(.*)$`)
	matches := re.FindStringSubmatch(name)
	Log.Debug("Regex Matches", "matches", matches, "len", len(matches))
	if len(matches) > 0 {
		return matches[1], matches[2]
	}

	return "", name
}
