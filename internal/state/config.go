package state

import (
	"log/slog"
	"os"
)

var RootDir = os.Getenv("GARDEN_LOG_DIR")
var InboxDir = "01. Inbox"

const (
	MenuIndexNumeric        = "   Numeric"
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

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)
}
