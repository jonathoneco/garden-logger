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

func (d *Directory) LoadEntry(dirEntry os.DirEntry) (*Entry, error) {
	index, name, err := parseEntryName(dirEntry.Name())
	if err != nil {
		return nil, err
	}

	isDir := dirEntry.IsDir()
	ext := ""
	if !isDir {
		ext = filepath.Ext(dirEntry.Name())
	}

	entry := &Entry{index, name, ext, isDir}
	return entry, nil
}

// Parses entry name and returns Index, CleanedName
func parseEntryName(name string) (int, string, error) {
	re := regexp.MustCompile(`^(?:(\d{2})\.\s+)?([^.]+)(?:\.(.+))?$`)
	matches := re.FindStringSubmatch(name)

	if len(matches) <= 0 {
		return -1, name, nil
	}

	cleanName := matches[2]

	if matches[1] == "" {
		return -1, cleanName, nil
	}

	index, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, "", err
	}

	return index, cleanName, nil

}

func (e *Entry) String() string {
	if e.NoteIndex == -1 {
		return fmt.Sprintf("%s%s", e.Name, e.Ext)
	}
	return fmt.Sprintf("%02d. %s%s", e.NoteIndex, e.Name, e.Ext)
}

func (d *Directory) FilePath(entry *Entry) string {
	return filepath.Join(d.AbsPath, entry.String())
}

// Directories

type Directory struct {
	Path      string
	AbsPath   string
	IsIndexed bool
	Entries   []*Entry
}

func (d *Directory) LoadEntries() error {

	dirEntries, err := os.ReadDir(d.AbsPath)
	if err != nil {
		return fmt.Errorf("failed to list entries for %s: %w", d.Path, err)
	}

	for _, dirEntry := range dirEntries {
		if strings.HasPrefix(dirEntry.Name(), ".") {
			continue
		}

		entry, err := d.LoadEntry(dirEntry)
		if err != nil {
			return err
		}
		d.Entries = append(d.Entries, entry)
	}

	slices.SortStableFunc(d.Entries, func(i, j *Entry) int { return cmp.Compare(i.NoteIndex, j.NoteIndex) })

	return nil
}

func (d *Directory) ListEntries() []string {

	var entries []string
	for _, entry := range d.Entries {
		entries = append(entries, entry.String())
	}
	sort.Strings(entries)

	return entries
}

func (d *Directory) FindEntryFromFilename(filename string) *Entry {
	for _, entry := range d.Entries {
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

func (d *Directory) MoveEntry(entry *Entry, newIndex int) error {

	oldPath := d.FilePath(entry)
	entry.NoteIndex = newIndex
	newPath := d.FilePath(entry)

	slog.Debug("Calling move entry on ", "entry", entry.String(), "oldPath", oldPath, "newPath", newPath)
	return os.Rename(oldPath, newPath)
}

// | 0-index | 1-index |
// | 0 | 1 |
// | 1 | 2 |
// | 2 | 3 |

func (d *Directory) MoveEntryUp(entry *Entry) error {
	if entry.IsDir && entry.NoteIndex <= 1 {
		slog.Debug("Cannot move directory above position 1")
		return nil
	}
	if !entry.IsDir && entry.NoteIndex <= d.NewDirIndex()-1 {
		slog.Debug("Cannot move entry: would conflict with directory ordering")
		return nil
	}

	entries := d.Entries

	index := entry.NoteIndex - 1
	swapIndex := index - 1
	err := d.MoveEntry(entries[swapIndex], entry.NoteIndex)
	if err != nil {
		return err
	}

	err = d.MoveEntry(entry, entry.NoteIndex-1)
	if err != nil {
		return err
	}

	entries[index], entries[swapIndex] = entries[swapIndex], entries[index]

	return nil
}

func (d *Directory) MoveEntryDown(entry *Entry) error {
	if entry.IsDir && entry.NoteIndex >= d.NewDirIndex()-1 {
		slog.Debug("Cannot move entry: would conflict with directory ordering")
		return nil
	}
	if !entry.IsDir && entry.NoteIndex >= d.NewFileIndex()-1 {
		slog.Debug("Cannot move directory above position 1")
		return nil
	}

	entries := d.Entries

	index := entry.NoteIndex - 1
	swapIndex := index + 1
	err := d.MoveEntry(entries[swapIndex], entry.NoteIndex)
	if err != nil {
		return err
	}

	err = d.MoveEntry(entry, entry.NoteIndex+1)
	if err != nil {
		return err
	}

	entries[index], entries[swapIndex] = entries[swapIndex], entries[index]
	return nil
}

func (d *Directory) UpdateIsIndex(isIndexed bool) error {
	indexFilePath := filepath.Join(d.AbsPath, ".index")

	if isIndexed {
		return os.WriteFile(indexFilePath, nil, 0644)
	}

	return os.Remove(indexFilePath)
}

func (d *Directory) NewDirIndex() int {
	if !d.IsIndexed {
		return -1

	}

	maxDirIndex := 0
	for index, entry := range d.Entries {
		if !entry.IsDir {
			break
		}
		maxDirIndex = index
	}
	return maxDirIndex + 1
}

func (d *Directory) NewFileIndex() int {
	if !d.IsIndexed {
		return -1
	}

	entries := d.Entries
	return entries[len(entries)-1].NoteIndex + 1
}

func (d *Directory) ApplyNumericIndexing() error {
	d.UpdateIsIndex(true)

	//Move Dirs
	dirIndex := 1
	for _, entry := range d.Entries {
		if entry.IsDir {
			err := d.MoveEntry(entry, dirIndex)
			if err != nil {
				return err
			}
			dirIndex++
		}
	}

	// Move Files
	fileIndex := dirIndex + 1
	for _, entry := range d.Entries {
		if !entry.IsDir {
			err := d.MoveEntry(entry, fileIndex)
			if err != nil {
				return err
			}
			fileIndex++
		}
	}
	return nil
}

func (d *Directory) RemoveIndexing() error {
	d.UpdateIsIndex(false)
	for _, entry := range d.Entries {
		err := d.MoveEntry(entry, -1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Directory) ValidateIndexing() error {
	if d.IsIndexed {
		foundFirstFile := false
		for index, entry := range d.Entries {
			// Errors if we run into a directory after flipping foundFirstFile at the first file
			if !entry.IsDir && !foundFirstFile {
				foundFirstFile = true
			} else if foundFirstFile {
				return fmt.Errorf("validation failed: found directory after file in %s", d.Path)
			}

			expectedIndex := index + 1
			if entry.NoteIndex != expectedIndex {
				return fmt.Errorf("index validation failed: entry %q has index %d, expected %d",
					entry.Name, entry.NoteIndex, expectedIndex)
			}
		}
	}

	return nil
}
