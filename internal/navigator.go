package internal

import (
	"fmt"
	"path/filepath"
)

type Navigator struct {
	currentDir    *Directory
	savedDir      *Directory
	savedTemplate string
	notes         *EntryService
}

func NewNavigator(notes *EntryService) *Navigator {
	return &Navigator{
		notes: notes,
	}
}

func (n *Navigator) NavigateTo(dirPath string) error {
	dir, err := n.notes.LoadDirectory(dirPath)
	if err != nil {
		return err
	}
	n.currentDir = dir
	return nil
}

func (n *Navigator) NavigateToParent() error {
	if n.currentDir.Path == "" {
		return fmt.Errorf("already at root directory")
	}

	parentPath := filepath.Dir(n.currentDir.Path)
	if parentPath == "." {
		parentPath = ""
	}

	return n.NavigateTo(parentPath)
}

func (n *Navigator) Reload() error {
	return n.currentDir.LoadEntries()
}

func (n *Navigator) Save() {
	n.savedDir = n.currentDir
}

func (n *Navigator) Restore() (*Directory, error) {
	if n.savedDir == nil {
		return nil, fmt.Errorf("no saved directory to restore")
	}

	restored := n.savedDir
	n.savedDir = nil
	return restored, nil
}

func (n *Navigator) CurrentDirectory() *Directory {
	return n.currentDir
}

func (n *Navigator) ListEntries() []string {
	return n.currentDir.ListEntries()
}

func (n *Navigator) SaveTemplate(templatePath string) {
	n.savedTemplate = templatePath
}

func (n *Navigator) RestoreTemplate() (string, error) {
	if n.savedTemplate == "" {
		return "", fmt.Errorf("no saved template to restore")
	}

	template := n.savedTemplate
	n.savedTemplate = ""
	return template, nil
}
