package menu

func (dirState *DirState) getPrompt() string {
	switch menuState.Mode {
	case ModeNew:
		return "New: "
	case ModeNewNote:
		return "Enter a name: "
	case ModeIndexing:
		return "Indexing: "
	default:
		return "Browse: "
	}
}

func (ms *MenuState) handleChoice(choice string) error {
	var err error = nil
	switch ms.Mode {
	case ModeBrowse:
		err = ms.handleBrowseChoice(choice)
	case ModeNew:
		err = ms.handleNewChoice(choice)
	case ModeIndexing:
		err = ms.handleIndexingChoice(choice)
	case ModeNewNote:
		err = ms.handleNewNote(choice)
	}

	return err
}

func (dirState *DirState) getMenuItems() ([]string, error) {
	switch menuState.Mode {
	case ModeNew:
	case ModeIndexing:
	case ModeBrowse: // browse
	default:
		return []string{}, nil
	}
}
