package main

import (
	"os"

	"garden-logger/internal/app"
)

func main() {
	if err := app.StartApp(); err != nil {
		os.Exit(1)
	}
}
