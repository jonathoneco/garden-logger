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
	NoteIndex int
	Name      string
	Ext       string
	IsDir     bool
}

func (dir *Directory) LoadEntry(dirEntry os.DirEntry) (*Entry, error) {
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

	slog.Debug("Regex matches while parsing entry name", "matches", matches)

	cleanName := matches[2]

	if matches[1] == "" {
		return -1, cleanName, nil
	}

	index, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, "", err
	}

	slog.Debug("Entry name parsed successfully", "index", index, "cleanName", cleanName)

	return index, cleanName, nil

}

func (entry *Entry) String() string {
	if entry.NoteIndex == -1 {
		return fmt.Sprintf("%s%s", entry.Name, entry.Ext)
	}
	return fmt.Sprintf("%02d. %s%s", entry.NoteIndex, entry.Name, entry.Ext)
}

func (dir *Directory) FilePath(entry *Entry) string {
	return filepath.Join(dir.AbsPath, entry.String())
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
	dir := &Directory{
		Path:      dirPath,
		AbsPath:   absPath,
		IsIndexed: isIndexed,
		Entries:   nil,
	}

	err := dir.LoadEntries()
	if err != nil {
		return nil, err
	}

	slog.Info("Directory state loaded successfully", "dirPath", dirPath, "entryCount", len(dir.Entries))

	return dir, err
}

func (dir *Directory) LoadEntries() error {
	slog.Debug("Reading directory entries", "absPath", dir.AbsPath)

	dirEntries, err := os.ReadDir(dir.AbsPath)
	if err != nil {
		return fmt.Errorf("failed to list entries for %s: %w", dir.Path, err)
	}

	slog.Debug("Directory read successfully", "rawEntryCount", len(dirEntries))

	for _, dirEntry := range dirEntries {
		if strings.HasPrefix(dirEntry.Name(), ".") {
			continue
		}

		entry, err := dir.LoadEntry(dirEntry)
		if err != nil {
			return err
		}
		dir.Entries = append(dir.Entries, entry)
	}

	slices.SortStableFunc(dir.Entries, func(i, j *Entry) int { return cmp.Compare(i.NoteIndex, j.NoteIndex) })

	return nil
}

func (dir *Directory) ListEntries() []string {
	slog.Debug("Listing entries for directory", "dirPath", dir.Path, "entryCount", len(dir.Entries))

	var entries []string
	for _, entry := range dir.Entries {
		entries = append(entries, entry.String())
	}
	sort.Strings(entries)

	slog.Debug("Entries listed and sorted", "sortedCount", len(entries))
	return entries
}

func (dir *Directory) FindEntryFromFilename(filename string) *Entry {
	for _, entry := range dir.Entries {
		if entry.String() == filename {
			return entry
		}
	}
	return nil
}

// Indexing
// A lot of this indexing relies on the entries array maintaining sorting

func LoadIsIndexed(absPath string) bool {
	indexFilePath := filepath.Join(absPath, ".index")

	if _, err := os.Stat(indexFilePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func (dir *Directory) MoveEntry(entry *Entry, newIndex int) error {

	oldPath := dir.FilePath(entry)
	entry.NoteIndex = newIndex
	newPath := dir.FilePath(entry)

	slog.Debug("Calling move entry on ", "entry", entry.String(), "oldPath", oldPath, "newPath", newPath)
	return os.Rename(oldPath, newPath)
}

// | 0-index | 1-index |
// | 0 | 1 |
// | 1 | 2 |
// | 2 | 3 |

func (dir *Directory) MoveEntryUp(entry *Entry) error {
	if entry.IsDir && entry.NoteIndex <= 1 {
		slog.Debug("Attempted to move directory out of bounds")
		return nil
	} else if entry.NoteIndex <= dir.NewDirIndex()-1 {
		slog.Debug("Attempted to move directory out of bounds")
	}

	entries := dir.Entries

	index := entry.NoteIndex - 1
	swapIndex := index - 1
	err := dir.MoveEntry(entries[swapIndex], entry.NoteIndex)
	if err != nil {
		return err
	}

	err = dir.MoveEntry(entry, entry.NoteIndex-1)
	if err != nil {
		return err
	}

	entries[index], entries[swapIndex] = entries[swapIndex], entries[index]

	return nil
}

func (dir *Directory) MoveEntryDown(entry *Entry) error {
	if entry.IsDir && entry.NoteIndex >= dir.NewDirIndex()-1 {
		slog.Debug("Attempted to move directory out of bounds")
	} else if entry.NoteIndex >= dir.NewFileIndex()-1 {
		slog.Debug("Attempted to move directory out of bounds")
	}

	entries := dir.Entries

	index := entry.NoteIndex - 1
	swapIndex := index + 1
	err := dir.MoveEntry(entries[swapIndex], entry.NoteIndex)
	if err != nil {
		return err
	}

	err = dir.MoveEntry(entry, entry.NoteIndex+1)
	if err != nil {
		return err
	}

	entries[index], entries[swapIndex] = entries[swapIndex], entries[index]
	return nil
}

func (dir *Directory) UpdateIsIndex(isIndexed bool) error {
	indexFilePath := filepath.Join(dir.AbsPath, ".index")

	if isIndexed {
		return os.WriteFile(indexFilePath, nil, 0644)
	}

	return os.Remove(indexFilePath)
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
	return entries[len(entries)-1].NoteIndex + 1
}

func (dir *Directory) ApplyNumericIndexing() error {
	dir.UpdateIsIndex(true)

	//Move Dirs
	dirIndex := 1
	for _, entry := range dir.Entries {
		if entry.IsDir {
			err := dir.MoveEntry(entry, dirIndex)
			if err != nil {
				return err
			}
			dirIndex++
		}
	}

	// Move Files
	fileIndex := dirIndex + 1
	for _, entry := range dir.Entries {
		if !entry.IsDir {
			err := dir.MoveEntry(entry, fileIndex)
			if err != nil {
				return err
			}
			fileIndex++
		}
	}
	return nil
}

func (dir *Directory) RemoveIndexing() error {
	dir.UpdateIsIndex(false)
	for _, entry := range dir.Entries {
		err := dir.MoveEntry(entry, -1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dir *Directory) ValidateIndexing() error {
	if dir.IsIndexed {
		dirTrip := false
		for index, entry := range dir.Entries {
			// Errors if we run into a directory after flipping dirTrip at the first file
			if !entry.IsDir && !dirTrip {
				dirTrip = true
			} else if dirTrip {
				return fmt.Errorf("Index Validation Error Directory after File")
			}

			if entry.NoteIndex != index+1 {
				return fmt.Errorf("Index Validation Error")
			}
		}
	}

	return nil
}
