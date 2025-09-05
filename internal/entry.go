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

type Entry struct {
	EntryIndex int
	Name       string
	Ext        string
	IsDir      bool
	ParentPath string
}

func (e *Entry) IsAnchor() bool {
	return e.EntryIndex == 0
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

	entry := &Entry{index, name, ext, isDir, d.AbsPath}
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
	if e.EntryIndex == -1 {
		return fmt.Sprintf("%s%s", e.Name, e.Ext)
	}
	return fmt.Sprintf("%02d. %s%s", e.EntryIndex, e.Name, e.Ext)
}

func (e *Entry) FilePath() string {
	return filepath.Join(e.ParentPath, e.String())
}

func (e *Entry) Move(newIndex int) error {
	if e.IsAnchor() {
		slog.Debug("Cannot move anchor entry")
		return nil
	}

	oldPath := e.FilePath()
	e.EntryIndex = newIndex
	newPath := e.FilePath()

	slog.Debug("Calling move entry on ", "entry", e.String(), "oldPath", oldPath, "newPath", newPath)
	return os.Rename(oldPath, newPath)
}

func (e *Entry) Remove() error {
	slog.Debug("Calling remove entry on ", "entry", e.String(), "path", e.ParentPath)
	return os.Remove(e.FilePath())
}

// Directories

type Directory struct {
	Path      string
	AbsPath   string
	IsIndexed bool
	Entries   []*Entry
}

func (d *Directory) GetEntryByIndex(index int) *Entry {
	for _, entry := range d.Entries {
		if index == entry.EntryIndex {
			return entry
		}
	}
	return nil
}

func (d *Directory) LoadEntries() error {
	d.Entries = nil

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

	slices.SortStableFunc(d.Entries, func(i, j *Entry) int { return cmp.Compare(i.EntryIndex, j.EntryIndex) })

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

func (d *Directory) GetEntryByFilename(filename string) *Entry {
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

// | 0-index | 1-index |
// | 0 | 1 |
// | 1 | 2 |
// | 2 | 3 |

func (d *Directory) MoveEntryUp(entry *Entry) error {
	if entry.IsAnchor() {
		slog.Debug("Cannot move anchor entry")
		return nil
	}
	if entry.IsDir && entry.EntryIndex <= 1 {
		slog.Debug("Cannot move directory above position 1")
		return nil
	}
	if !entry.IsDir && entry.EntryIndex <= d.NewDirIndex()-1 {
		slog.Debug("Cannot move entry: would conflict with directory ordering")
		return nil
	}

	swapEntry := d.GetEntryByIndex(entry.EntryIndex - 1)

	err := entry.Move(entry.EntryIndex - 1)
	if err != nil {
		return err
	}

	err = swapEntry.Move(entry.EntryIndex)
	if err != nil {
		return err
	}

	return nil
}

func (d *Directory) MoveEntryDown(entry *Entry) error {
	if entry.IsAnchor() {
		slog.Debug("Cannot move anchor entry")
		return nil
	}
	if entry.IsDir && entry.EntryIndex >= d.NewDirIndex()-1 {
		slog.Debug("Cannot move entry: would conflict with directory ordering")
		return nil
	}
	if !entry.IsDir && entry.EntryIndex >= d.NewFileIndex()-1 {
		slog.Debug("Cannot move directory below last position")
		return nil
	}

	swapEntry := d.GetEntryByIndex(entry.EntryIndex - 1)

	err := entry.Move(entry.EntryIndex + 1)
	if err != nil {
		return err
	}

	err = swapEntry.Move(entry.EntryIndex)
	if err != nil {
		return err
	}

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
	for _, entry := range d.Entries {
		if !entry.IsDir {
			continue
		}
		if entry.IsAnchor() {
			continue
		}
		if entry.EntryIndex > maxDirIndex {
			maxDirIndex = entry.EntryIndex
		}
	}
	return maxDirIndex + 1
}

func (d *Directory) NewFileIndex() int {
	if !d.IsIndexed {
		return -1
	}

	maxFileIndex := 0
	for _, entry := range d.Entries {
		if entry.IsAnchor() {
			continue
		}
		if entry.EntryIndex > maxFileIndex {
			maxFileIndex = entry.EntryIndex
		}
	}
	return maxFileIndex + 1
}

func (d *Directory) InsertEntry(e *Entry) error {
	if e.EntryIndex == -1 {
		// Non-indexed entry, just append to the list
		d.Entries = append(d.Entries, e)
		return nil
	}

	d.Entries = append(d.Entries, e)
	for _, entry := range d.Entries {
		if entry.IsAnchor() {
			continue
		}
		if entry.EntryIndex >= e.EntryIndex {
			entry.Move(entry.EntryIndex + 1)
		}
	}
	return nil
}

func (d *Directory) DeleteEntry(e *Entry) error {
	if e.IsAnchor() {
		slog.Debug("Cannot delete anchor entry")
		return nil
	}

	if d.IsIndexed {
		for i := e.EntryIndex - 1; i < len(d.Entries); i++ {
			entry := d.Entries[i]
			if e == entry || entry.IsAnchor() {
				continue
			}
			err := entry.Move(entry.EntryIndex - 1)
			if err != nil {
				return err
			}
		}
	}

	e.Remove()

	return nil
}

func (d *Directory) ApplyNumericIndexing() error {
	slog.Debug("Starting ApplyNumericIndexing", "path", d.Path, "currentEntries", len(d.Entries))
	d.UpdateIsIndex(true)

	// Move Dirs first
	dirIndex := 1
	for _, entry := range d.Entries {
		if entry.IsDir && !entry.IsAnchor() {
			slog.Debug("Moving directory", "entry", entry.String(), "fromIndex", entry.EntryIndex, "toIndex", dirIndex)
			err := entry.Move(dirIndex)
			if err != nil {
				return err
			}
			dirIndex++
		}
	}

	// Move Files
	fileIndex := dirIndex
	for _, entry := range d.Entries {
		if !entry.IsDir && !entry.IsAnchor() {
			slog.Debug("Moving file", "entry", entry.String(), "fromIndex", entry.EntryIndex, "toIndex", fileIndex)
			err := entry.Move(fileIndex)
			if err != nil {
				return err
			}
			fileIndex++
		}
	}

	slog.Debug("Completed ApplyNumericIndexing", "path", d.Path)
	return nil
}

func (d *Directory) RemoveIndexing() error {
	slog.Debug("Starting RemoveIndexing", "path", d.Path, "currentEntries", len(d.Entries))
	d.UpdateIsIndex(false)
	for _, entry := range d.Entries {
		if !entry.IsAnchor() {
			slog.Debug("Removing index from entry", "entry", entry.String(), "fromIndex", entry.EntryIndex)
			err := entry.Move(-1)
			if err != nil {
				return err
			}
		}
	}
	slog.Debug("Completed RemoveIndexing", "path", d.Path)
	return nil
}

func (d *Directory) ValidateIndexing() error {
	if d.IsIndexed {
		foundFirstFile := false
		nonAnchorIndex := 1
		for _, entry := range d.Entries {
			// Skip anchor entries in validation
			if entry.IsAnchor() {
				continue
			}

			// Errors if we run into a directory after flipping foundFirstFile at the first file
			if !entry.IsDir && !foundFirstFile {
				foundFirstFile = true
			} else if foundFirstFile {
				return fmt.Errorf("validation failed: found directory after file in %s", d.Path)
			}

			if entry.EntryIndex != nonAnchorIndex {
				return fmt.Errorf("index validation failed: entry %q has index %d, expected %d",
					entry.Name, entry.EntryIndex, nonAnchorIndex)
			}
			nonAnchorIndex++
		}
	}

	return nil
}
