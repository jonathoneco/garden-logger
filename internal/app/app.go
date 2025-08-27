package app

import (
	"flag"
	"fmt"
	"garden-logger/internal/state"
	"log/slog"
)

var log = slog.Default().With("package", "app")

func StartApp() error {
	var verbose bool

	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.Parse()

	state.InitLogger(verbose)

	if state.RootDir == "" {
		err := fmt.Errorf("GARDEN_LOG_DIR environment variable is not set")
		log.Error("Startup Error", "error", err)
		return err
	}

	log.Debug("Garden Logger Started", "rootDir", state.RootDir)

	// err := menu.Browse()
	// if err != nil {
	// 	log.Error("Operational Error", "error", err)
	// }
	// return err
	state.LoadDirState("")
	state.LoadDirState("01. Inbox")

	return nil
}
