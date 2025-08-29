package internal

import (
	"flag"
	"log/slog"
)

func StartApp() error {
	var verbose bool

	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.Parse()

	InitLogger(verbose)

	slog.Info("Application startup initiated", "verbose", verbose)

	err := Browse()
	if err != nil {
		return err
	}

	slog.Info("Application startup completed successfully")
	return nil
}
