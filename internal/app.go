package internal

import (
	"flag"
	"fmt"
	"log/slog"
)

func StartApp() error {
	var verbose bool

	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.Parse()

	InitLogger(verbose)

	slog.Info("Application startup initiated", "verbose", verbose)

	if RootDir == "" {
		err := fmt.Errorf("GARDEN_LOG_DIR environment variable is not set")
		slog.Error("Startup Error", "error", err)
		return err
	}

	Browse()

	slog.Info("Application startup completed successfully")
	return nil
}
