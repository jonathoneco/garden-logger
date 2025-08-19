package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Status

type StatusInfo struct {
	AbsolutePath     string
	IndexingStrategy string
	RelativePath     string
}

func formatStatusMessage(info *StatusInfo) string {
	path := info.RelativePath
	if path == "" {
		path = "/"
	}
	return fmt.Sprintf("Path: %s | Indexing: %s", path, info.IndexingStrategy)
}

func buildStatusInfo(dirPath string) (*StatusInfo, error) {
	absPath := filepath.Join(rootDir, dirPath)

	config, err := getIndexConfig(absPath)
	if err != nil {
		return nil, err
	}

	relativePath := dirPath
	if relativePath == "" {
		relativePath = "/"
	}

	indexingStrategy := config.Strategy.String()

	return &StatusInfo{
		AbsolutePath:     absPath,
		IndexingStrategy: indexingStrategy,
		RelativePath:     relativePath,
	}, nil
}

func listEntries(statusInfo *StatusInfo) ([]string, error) {
	fullPath := statusInfo.AbsolutePath
	dirEntries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading entries for path %s: %w", fullPath, err)

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

func browse(dirPath string) error {

	// Indexing
	// TODO: Parse .index file and pass indexing strategy alongside dirPath

	Logger.Debug("Browsing", "Current Path", dirPath)

	statusInfo, err := buildStatusInfo(dirPath)
	if err != nil {
		return err
	}

	Logger.Debug("Built Status Info",
		"Absolute Path", statusInfo.AbsolutePath,
		"Indexing Strategy", statusInfo.IndexingStrategy,
		"Relative Path", statusInfo.RelativePath,
	)

	entries, err := listEntries(statusInfo)
	if err != nil {
		return err
	}

	var menuItems []string

	menuItems = append(menuItems, "New")
	menuItems = append(menuItems, entries...)

	if dirPath != "" {
		menuItems = append(menuItems, "..")
	}
	menuItems = append(menuItems, "/")

	choice, err := launchMenu(menuItems, "The Garden Log: ", statusInfo)
	if err != nil {
		return err
	}

	switch choice {
	case "New":
		Logger.Debug("Starting new note creation")
		return createNewNote(statusInfo)
	case "..":
		Logger.Debug("Navigating to parent folder")
		parentFolder := filepath.Dir(strings.TrimSuffix(dirPath, "/"))
		if parentFolder == "." {
			parentFolder = ""
		} else {
			parentFolder += "/"
		}
		return browse(parentFolder)
	case "/":
		Logger.Debug("Launching directory session")
		return launchDir(dirPath)
	default:
		Logger.Debug("Navigating to File or Directory")
		return handleFileOrDirectory(choice, dirPath)
	}
}

func handleFileOrDirectory(choice, dirPath string) error {
	if strings.HasSuffix(choice, "/") {
		newDirPath := dirPath + choice
		return browse(newDirPath)
	}

	if strings.HasSuffix(strings.ToLower(choice), ".md") {
		fullFilePath := filepath.Join(dirPath, choice)
		return launchNote(fullFilePath)
	}

	return fmt.Errorf("[ERROR] unexpected choice: %s", choice)
}

func launchNote(filePath string) error {
	cmd := exec.Command("kitty", "-e", "nvim", filePath)

	cmd.Dir = rootDir

	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to launch kitty: %w", err)
	}

	return nil
}

func launchDir(dirPath string) error {
	var sessionName string
	if dirPath == "" {
		sessionName = "The Garden Log"
	} else {
		sessionName = filepath.Base(dirPath)
	}

	fullPath := filepath.Join(rootDir, dirPath)

	cmd := exec.Command("kitty", "-e", "tmux", "new-session", "-s", sessionName, "-c", fullPath, "nvim .")

	cmd.Dir = fullPath

	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to launch kitty: %w", err)
	}

	return nil
}

func launchMenu(items []string, prompt string, statusInfo *StatusInfo) (string, error) {
	args := []string{"notes", "-dmenu", "-l", "10", "-i", "-p", prompt}

	if statusInfo != nil {
		statusMsg := formatStatusMessage(statusInfo)
		args = append(args, "-mesg", statusMsg)
	}

	cmd := exec.Command("rofi-launcher", args...)
	Logger.Debug("Launching Menu", "cmd", cmd.String())
	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting output when launching rofi: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
