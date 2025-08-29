package internal

import (
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// Entries

type Kind int

const (
	KindFile Kind = iota
	KindDir
)

type Entry struct {
	Index int
	Name  string
	Ext   string
	IsDir bool
}

func LoadEntry(dirEntry os.DirEntry) (*Entry, error) {
	slog.Debug("Processing directory entry", "name", dirEntry.Name(), "isDir", dirEntry.IsDir())

	index, name, err := parseEntryName(dirEntry.Name())
	if err != nil {
		return nil, err
	}

	isDir := dirEntry.IsDir()
	ext := "/"
	if !isDir {
		ext = filepath.Ext(dirEntry.Name())
	}

	entry := &Entry{index, name, ext, isDir}
	slog.Debug("Entry loaded successfully", "index", index, "name", name, "ext", ext, "isDir", isDir)
	return entry, nil
}

// Parses entry name and returns Index, CleanedName
func parseEntryName(name string) (int, string, error) {
	slog.Debug("Parsing entry name", "rawName", name)

	re := regexp.MustCompile(`^(?:(\d{2})\.\s+)?([^.]+)(?:\.(.+))?$`)
	matches := re.FindStringSubmatch(name)

	if len(matches) <= 0 {
		slog.Debug("Entry name parsing fallback", "originalName", name)
		return -1, name, nil
	}

	indexString := matches[1]
	cleanName := matches[2]
	slog.Debug("Entry name parsed successfully", "index", indexString, "cleanName", cleanName)

	index, err := strconv.Atoi(indexString)
	if err != nil {
		return 0, "", err
	}

	return index, cleanName, nil

}

func (e *Entry) String() string {
	return fmt.Sprintf("%02d. %s%s", e.Index, e.Name, e.Ext)
}

// Directories

type Directory struct {
	Path      string
	AbsPath   string
	IsIndexed bool
	Entries   []*Entry
}

func LoadDirectory(dirPath string) (*Directory, error) {
	slog.Info("Loading directory state", "dirPath", dirPath)
	absPath := filepath.Join(RootDir, dirPath)
	slog.Debug("Directory paths resolved", "dirPath", dirPath, "absPath", absPath)

	isIndexed := false

	entries, err := LoadEntries(dirPath)
	if err != nil {
		return nil, err
	}

	slog.Info("Directory state loaded successfully", "dirPath", dirPath, "entryCount", len(entries))

	return &Directory{
		Path:      dirPath,
		AbsPath:   absPath,
		IsIndexed: isIndexed,
		Entries:   entries,
	}, err
}

func LoadEntries(dirPath string) ([]*Entry, error) {
	absPath := filepath.Join(RootDir, dirPath)
	slog.Debug("Reading directory entries", "absPath", absPath)

	dirEntries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list entries for %s: %w", dirPath, err)
	}

	slog.Debug("Directory read successfully", "rawEntryCount", len(dirEntries))

	var entries []*Entry

	for _, dirEntry := range dirEntries {
		if strings.HasPrefix(dirEntry.Name(), ".") {
			continue
		}

		entry, err := LoadEntry(dirEntry)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	slices.SortStableFunc(entries, func(i, j *Entry) int { return cmp.Compare(i.Index, j.Index) })

	return entries, nil
}

func (dir *Directory) ListEntries() []string {
	slog.Debug("Listing entries for directory", "dirPath", dir.Path, "entryCount", len(dir.Entries))

	var entries []string
	for _, e := range dir.Entries {
		entries = append(entries, e.String())
	}
	sort.Strings(entries)

	slog.Debug("Entries listed and sorted", "sortedCount", len(entries))
	return entries
}

// Indexing

func LoadIsIndexed(absPath string) bool {
	indexFilePath := filepath.Join(absPath, ".index")

	if _, err := os.Stat(indexFilePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func (dir *Directory) UpdateIsIndex(isIndexed bool) {
	indexFilePath := filepath.Join(dir.AbsPath, ".index")

	if isIndexed {
		os.WriteFile(indexFilePath, nil, 0644)
	}

	os.Remove(indexFilePath)
}

func (dir *Directory) NewDirIndex() int {
	if !dir.IsIndexed {
		return -1

	}

	maxDirIndex := 0
	for index, entry := range dir.Entries {
		if !entry.IsDir {
			break
		}
		maxDirIndex = index
	}
	return maxDirIndex + 1
}

func (dir *Directory) NewFileIndex() int {
	if !dir.IsIndexed {
		return -1

	}

	entries := dir.Entries
	return entries[len(entries)-1].Index + 1
}

func (dir *Directory) ApplyNumericIndexing() error

func (dir *Directory) RemoveIndexing() error

func (dir *Directory) ValidateIndexing() error {
	if dir.IsIndexed {
		for index, entry := range dir.Entries {
			if entry.Index != index+1 {
				return fmt.Errorf("Index Validation Error")
			}
		}
	}

	return nil
}
