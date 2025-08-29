package internal

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// LaunchSuccessError signals that the application should exit after successful launch
type LaunchSuccessError struct {
	Message string
}

func (e LaunchSuccessError) Error() string { return e.Message }

// NotesService handles all note and directory operations
type NotesService struct {
	config *Config
}

func NewNotesService(config *Config) *NotesService {
	return &NotesService{config: config}
}

func (s *NotesService) LoadDirectory(dirPath string) (*Directory, error) {
	absPath := filepath.Join(s.config.RootDir, dirPath)

	isIndexed := LoadIsIndexed(absPath)
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

	slog.Debug("Loaded Directory", "directory", dirPath, "entries", dir.Entries)

	return dir, nil
}

func (s *NotesService) CreateNote(dir *Directory, entry *Entry) (string, error) {

	var targetDir string

	if dir.Path == "" {
		targetDir = filepath.Join(s.config.RootDir, s.config.InboxDir)
	} else {
		targetDir = dir.AbsPath
	}

	fullPath := filepath.Join(targetDir, entry.String())

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create note file %s: %w", fullPath, err)
	}
	defer file.Close()
	frontmatter := fmt.Sprintf("# %s\n\n", entry.Name)

	if _, err := file.WriteString(frontmatter); err != nil {
		return "", fmt.Errorf("failed to write frontmatter to note: %w", err)
	}

	return fullPath, nil
}

func (s *NotesService) LaunchNoteEditor(filePath string) error {

	cmd := exec.Command("kitty", "-e", "nvim", filePath)

	cmd.Dir = s.config.RootDir
	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to launch note editor: %w", err)
	}

	return LaunchSuccessError{Message: "note editor launched successfully"}
}

func (s *NotesService) LaunchDirectoryEditor(dirPath string) error {

	var sessionName string
	if dirPath == "" {
		sessionName = "The Garden Log"
	} else {
		sessionName = filepath.Base(dirPath)
	}

	fullPath := filepath.Join(s.config.RootDir, dirPath)

	cmd := exec.Command("kitty", "-e", "tmux", "new-session", "-s", sessionName, "-c", fullPath, "nvim .")

	cmd.Dir = fullPath
	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to launch directory editor: %w", err)
	}

	return LaunchSuccessError{Message: "directory editor launched successfully"}
}

func (s *NotesService) CreateNoteFromUserInput(dir *Directory, name string) (string, error) {
	if name == "" {
		name = time.Now().Format("2006-01-02")
	}

	entry := &Entry{
		Name:      name,
		NoteIndex: dir.NewFileIndex(),
		Ext:       ".md",
		IsDir:     false,
	}

	return s.CreateNote(dir, entry)
}

