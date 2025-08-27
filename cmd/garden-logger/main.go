package main

import (
	"garden-logger/internal/app"
	"os"
)

func main() {
	if err := app.StartApp(); err != nil {
		os.Exit(1)
	}
}
