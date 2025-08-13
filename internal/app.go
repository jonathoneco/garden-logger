// Gonna just use one entry struct, and every layer uses the filesystem as
// state directly, effectively we'll cd in and out of shit and re-read
// everything, instead of modeling the whole filesystem for now.
package internal

import (
	"fmt"
	"os"
)

var rootDir = os.Getenv("GARDEN_LOG_DIR")

func StartApp() error {
	if rootDir == "" {
		return fmt.Errorf("GARDEN_LOG_DIR environment variable is not set")
	}

	return browse("")
}
