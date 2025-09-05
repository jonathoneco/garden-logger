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

// EntryService handles all note and directory operations
type EntryService struct {
	config *Config
}

func NewNotesService(config *Config) *EntryService {
	return &EntryService{config: config}
}

func (s *EntryService) LoadDirectory(dirPath string) (*Directory, error) {
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

func (s *EntryService) CreateEntry(d *Directory, entry *Entry) (string, error) {
	slog.Debug("Creating entry", "name", entry.Name, "index", entry.EntryIndex, "isDir", entry.IsDir, "parentPath", entry.ParentPath)

	var targetDir string

	if d.Path == "" && !entry.IsDir {
		targetDir = filepath.Join(s.config.RootDir, s.config.InboxDir)
	} else {
		targetDir = d.AbsPath
	}

	fullPath := filepath.Join(targetDir, entry.String())
	slog.Debug("Creating at path", "fullPath", fullPath, "targetDir", targetDir)

	if entry.IsDir {
		err := os.Mkdir(fullPath, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create note directory %s: %w", fullPath, err)
		}
	} else {
		file, err := os.Create(fullPath)
		if err != nil {
			return "", fmt.Errorf("failed to create note file %s: %w", fullPath, err)
		}
		defer file.Close()
		frontmatter := fmt.Sprintf("# %s\n\n", entry.Name)

		if _, err := file.WriteString(frontmatter); err != nil {
			return "", fmt.Errorf("failed to write frontmatter to note: %w", err)
		}
	}

	if err := d.InsertEntry(entry); err != nil {
		return "", err
	}

	return entry.FilePath(), nil
}

func (s *EntryService) LaunchNoteEditor(filePath string) error {
	fullPath := filepath.Join(s.config.RootDir, filePath)

	cmd := exec.Command("kitty", "--title", "The Garden Log", "-e", "nvim", fullPath)

	cmd.Dir = s.config.RootDir
	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to launch note editor: %w", err)
	}

	return LaunchSuccessError{Message: "note editor launched successfully"}
}

func (s *EntryService) LaunchDirectoryEditor(dirPath string) error {
	fullPath := filepath.Join(s.config.RootDir, dirPath)

	cmd := exec.Command("kitty", "-e", "tmux-sessionizer", fullPath)
	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to launch directory editor: %w", err)
	}

	return LaunchSuccessError{Message: "directory editor launched successfully"}
}

func (s *EntryService) CreateEntryFromUserInput(d *Directory, name string, isDir bool) (string, error) {
	slog.Debug("Creating entry from user input", "name", name, "isDir", isDir, "dirPath", d.Path)
	if name == "" {
		name = time.Now().Format("2006-01-02")
	}

	ext := ""
	if !isDir {
		ext = ".md"
	}

	index := d.NewFileIndex()
	if isDir {
		index = d.NewDirIndex()
	}

	entry := &Entry{
		Name:       name,
		EntryIndex: index,
		Ext:        ext,
		IsDir:      isDir,
		ParentPath: d.Path,
	}

	return s.CreateEntry(d, entry)
}

func (s *EntryService) CreateEntryFromTemplate(d *Directory, name string, templatePath string) (string, error) {
	slog.Debug("Creating entry from template", "name", name, "templatePath", templatePath, "dirPath", d.Path)
	if name == "" {
		name = time.Now().Format("2006-01-02")
	}

	entry := &Entry{
		Name:       name,
		EntryIndex: d.NewFileIndex(),
		Ext:        ".md",
		IsDir:      false,
		ParentPath: d.AbsPath,
	}

	filePath, err := s.CreateEntry(d, entry)
	if err != nil {
		return "", err
	}

	absTemplatePath := filepath.Join(s.config.RootDir, templatePath)
	templateContent, err := os.ReadFile(absTemplatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", absTemplatePath, err)
	}

	frontmatter := fmt.Sprintf("# %s\n\n", entry.Name)
	finalContent := frontmatter + string(templateContent)

	err = os.WriteFile(filePath, []byte(finalContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write template content to %s: %w", filePath, err)
	}

	return filePath, nil
}
