package modes

func getNewNotePrompt() string {
	return ""
}

func getNewNoteMenuItems() ([]string, error) {
	return []string{}, nil
}

func handleNewNoteChoice(choice string) error {
	name := choice
	if name == "" {
		name = time.Now().Format("2006-01-02")
	}

	filePath, err := ms.writeNote(name)
	if err != nil {
		return err
	}

	return launchNote(filePath)
}
