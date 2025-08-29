package internal

import (
	"log/slog"
	"os"
)

var RootDir = os.Getenv("GARDEN_LOG_DIR")
var InboxDir = "01. Inbox"

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
