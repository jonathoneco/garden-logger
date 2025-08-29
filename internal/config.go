package internal

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	RootDir  string
	InboxDir string
}

func LoadConfig() (*Config, error) {
	rootDir := os.Getenv("GARDEN_LOG_DIR")
	if rootDir == "" {
		return nil, fmt.Errorf("GARDEN_LOG_DIR environment variable is not set")
	}
	return &Config{
		RootDir:  rootDir,
		InboxDir: "01. Inbox",
	}, nil
}

const (
	MenuIndexSetting        = "   Numeric Indexing"
	MenuIndexDatetime       = "󰃭   Datetime"
	MenuIndexNone           = "󰟢   None"
	MenuNew                 = "   New"
	MenuNewNote             = "   New Note"
	MenuNewDirectory        = "   New Directory"
	MenuNewNoteFromTemplate = "   New Note from Template"
	MenuBack                = "←   Back"
	MenuSettings            = "   Settings"
	MenuOpenCurrentFolder   = "   Open Current Folder"
)

func InitLogger(verbose bool) {
	level := slog.LevelInfo

	if verbose {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)
}
