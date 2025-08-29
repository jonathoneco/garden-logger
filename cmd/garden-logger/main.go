package main

import (
	"errors"
	"garden-logger/internal"
	"log/slog"
	"os"
)

func main() {
	slog.Info("Garden Logger main entry point")

	if err := internal.StartApp(); err != nil {
		var launchErr internal.LaunchSuccessError
		if errors.As(err, &launchErr) {
			os.Exit(0) // Success - program launched editor
		}
		slog.Error("Application Error", "error", err)
		os.Exit(1)
	}

	slog.Info("Garden Logger completed successfully")
}
