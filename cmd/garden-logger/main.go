package main

import (
	"garden-logger/internal"
	"log/slog"
	"os"
)

func main() {
	slog.Info("Garden Logger main entry point")

	if err := internal.StartApp(); err != nil {
		slog.Error("Application Error", "error", err)
		os.Exit(1)
	}

	slog.Info("Garden Logger completed successfully")
}
