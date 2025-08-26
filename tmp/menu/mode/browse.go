package modes

func getBrowseMenuItems() ([]string, error) {
	entries, err := dirState.ListEntries()

	if err != nil {
		return entries, err
	}

	items := []string{config.MenuNew, config.MenuSettings}
	items = append(items, entries...)
	if dirState.RelativePath != "" {
		items = append(items, config.MenuBack)
	}
	items = append(items, config.MenuOpenCurrentFolder)
	return items, nil
}

func handleBrowseChoice(choice string) error {
	switch choice {
	case config.MenuNew:
		ms.Mode = ModeNew
	case config.MenuSettings:
		ms.Mode = ModeIndexing
	case config.MenuBack:
		err := ms.navigateToParent()
		if err != nil {
			return err
		}
	case config.MenuOpenCurrentFolder:
		return launchDir(ms.DirState.RelativePath)
	default:
		if strings.HasSuffix(choice, "/") {
			newDirPath := filepath.Join(ms.DirState.RelativePath, choice)
			ms.navigateTo(newDirPath)
			return nil
		}

		if strings.HasSuffix(strings.ToLower(choice), ".md") {
			fullFilePath := filepath.Join(ms.DirState.RelativePath, choice)
			return launchNote(fullFilePath)
		}

		return fmt.Errorf("[ERROR] unexpected choice: %s", choice)
	}

	return nil
}
