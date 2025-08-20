package internal

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

var rootDir = os.Getenv("GARDEN_LOG_DIR")
var inboxDir = "1 Inbox"
var Logger *slog.Logger

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

func initLogger(verbose bool) {
	level := slog.LevelInfo

	if verbose {
		level = slog.LevelDebug
	}

	Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
}

func StartApp() error {

	var verbose bool

	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.Parse()

	initLogger(verbose)

	if rootDir == "" {
		err := fmt.Errorf("GARDEN_LOG_DIR environment variable is not set")
		Logger.Error("Startup Error", "error", err)
		return err
	}

	Logger.Debug("Garden Logger Started", "rootDir", rootDir)

	err := browse("")
	if err != nil {
		Logger.Error("Operational Error", "error", err)
	}
	return err
}
