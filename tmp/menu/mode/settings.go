package modes

func getSettingsPrompt() string {
	return ""
}

func getSettingsMenuItems() ([]string, error) {
	current := dirState.IndexingStrategy
	return []string{
		formatOption(config.MenuIndexNumeric, current == indexing.IndexStrategyNumeric),
		formatOption(config.MenuIndexNone, current == indexing.IndexStrategyNone),
		config.MenuBack,
	}, nil
}

func handleSettingsChoice(choice string) error {
	switch choice {
	case config.MenuIndexNumeric:
		if ms.DirState.IndexingStrategy != indexing.IndexStrategyNumeric {
			config := &indexing.IndexConfig{
				Strategy: indexing.IndexStrategyNumeric,
				NumericConfig: &indexing.NumericConfig{
					DirPriority: true, // Default to directory priority
				},
			}
			if err := indexing.WriteIndexConfig(ms.DirState.AbsolutePath, config); err != nil {
				return fmt.Errorf("failed to set numeric indexing: %w", err)
			}
			if err := indexing.ApplyNumericIndexing(ms.DirState.AbsolutePath, true); err != nil {
				return fmt.Errorf("failed to apply numeric indexing: %w", err)
			}
			ms.DirState.IndexingStrategy = indexing.IndexStrategyNumeric
		}
	case config.MenuIndexNone:
		if ms.DirState.IndexingStrategy != indexing.IndexStrategyNone {
			if err := indexing.RemoveIndexing(ms.DirState.AbsolutePath); err != nil {
				return fmt.Errorf("failed to remove indexing: %w", err)
			}
			config := &indexing.IndexConfig{Strategy: indexing.IndexStrategyNone}
			if err := indexing.WriteIndexConfig(ms.DirState.AbsolutePath, config); err != nil {
				return fmt.Errorf("failed to clear indexing config: %w", err)
			}
			ms.DirState.IndexingStrategy = indexing.IndexStrategyNone
		}
	}
	ms.Mode = ModeBrowse
	return nil
}
