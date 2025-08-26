package modes

import "garden-logger/internal/state"

func getNewPrompt() string {
	return ""
}

func getNewMenuItems() ([]string, error) {
	return []string{state.MenuNewNote, state.MenuNewDirectory, state.MenuNewNoteFromTemplate, state.MenuBack}, nil
}

func handleNewChoice(choice string) error {
	switch choice {
	case state.MenuNewNote:
		ms.Mode = ModeNewNote
	case state.MenuNewDirectory:
	case state.MenuNewNoteFromTemplate:
	case state.MenuBack:
		ms.Mode = ModeBrowse
	}
	return nil
}
