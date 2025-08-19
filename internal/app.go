// Gonna just use one entry struct, and every layer uses the filesystem as
// state directly, effectively we'll cd in and out of shit and re-read
// everything, instead of modeling the whole filesystem for now.
package internal

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

var rootDir = os.Getenv("GARDEN_LOG_DIR")
var inboxDir = os.Getenv("1 Inbox")
var Logger *slog.Logger

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
