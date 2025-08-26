package menu

import (
	"fmt"
	"log/slog"
)

var log = slog.Default().With("package", "menu")

func (dirState *DirState) formatStatusMessage() string {
	path := dirState.RelativePath
	return fmt.Sprintf("Path: %s \nIndexing: %s", path, dirState.IndexingStrategy)
}

func formatOption(text string, selected bool) string {
	if selected {
		return text + " ✓"
	}
	return text
}

// func Browse() error {
// 	dirState, err := state.InitDirState()
// 	if err != nil {
// 		return err
// 	}
//
// 	state := &MenuState{Mode: ModeBrowse, DirState: dirState}
//
// 	for {
// 		log.Debug("Browsing", "Current Path", dirState.RelativePath, "Mode", state.Mode)
// 		choice, err := state.launchMenu()
// 		log.Debug("Selection", "Choice", choice)
// 		if err != nil {
// 			return err
// 		}
//
// 		// Skip handling if choice is empty (file movement operations return empty string)
// 		if choice == "" {
// 			continue
// 		}
//
// 		err = state.handleChoice(choice)
// 		if err != nil {
// 			return err
// 		}
// 	}
// }
