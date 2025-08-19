package main

import (
	"os"

	"garden-logger/internal"
)

func main() {
	if err := internal.StartApp(); err != nil {
		os.Exit(1)
	}
}
