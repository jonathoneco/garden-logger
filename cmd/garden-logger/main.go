package main

import (
	"fmt"
	"os"

	"garden-logger/internal"
)

func main() {
	// Start the application
	if err := internal.StartApp(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
