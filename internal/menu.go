package internal

type Mode int

const (
	ModeBrowse Mode = iota
	ModeNew
	ModeNewNote
	ModeNewDirectory
	ModePickTemplate
	ModeSettings
)

func (mode Mode) String() string {
	switch mode {
	case ModeBrowse:
		return "ModeBrowse"
	case ModeNew:
		return "ModeNew"
	case ModeNewNote:
		return "ModeNewNote"
	case ModeNewDirectory:
		return "ModeNewDirectory"
	case ModePickTemplate:
		return "ModePickTemplate"
	case ModeSettings:
		return "ModeSettings"
	default:
		return ""
	}
}

type MenuState struct {
	Mode      Mode
	Selection string
	config    *Config
	nav       *Navigator
	notes     *EntryService
}

// func (m *MenuState) formatStatusMessage() string {
// 	path := m.nav.CurrentDirectory().Path
// 	return fmt.Sprintf("Path: %s \nIndexing: %s", path, "TEMP")
// }

func InitMenuState() (*MenuState, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	notes := NewNotesService(config)
	nav := NewNavigator(notes)

	err = nav.NavigateTo("")
	if err != nil {
		return nil, err
	}

	menu := &MenuState{ModeBrowse, "", config, nav, notes}
	return menu, nil
}

func (m *MenuState) getNavigationMenuItems() []string {
	entries := m.nav.ListEntries()
	if m.nav.CurrentDirectory().Path != "" {
		entries = append(entries, MenuBack)
	}
	return entries
}

func (m *MenuState) getPrompt() string {
	switch m.Mode {
	case ModeNew:
		return "New: "
	case ModeNewNote:
		return "Enter a file name: "
	case ModeNewDirectory:
		return "Enter a folder name: "
	case ModePickTemplate:
		return "Pick a template: "
	case ModeSettings:
		return "Indexing: "
	default:
		return "Browse: "
	}
}

func (m *MenuState) handleChoice(choice string) error {
	var err error = nil
	switch m.Mode {
	case ModeBrowse:
		err = m.handleBrowseChoice(choice)
	case ModeNew:
		err = m.handleNewChoice(choice)
	case ModePickTemplate:
		err = m.handleTemplateChoice(choice)
	case ModeSettings:
		err = m.handleSettingsChoice(choice)
	case ModeNewNote:
		err = m.handleNewEntry(choice, false)
	case ModeNewDirectory:
		err = m.handleNewEntry(choice, true)
	}

	return err
}

func (m *MenuState) getMenuItems() ([]string, error) {
	switch m.Mode {
	case ModeNew:
		return getNewMenuItems()
	case ModeSettings:
		return m.getSettingsMenuItems()
	case ModeBrowse: // browse
		return m.getBrowseMenuItems()
	case ModePickTemplate:
		return m.getNavigationMenuItems(), nil
	default:
		return nil, nil
	}
}

func Browse() error {
	menu, err := InitMenuState()
	if err != nil {
		return err
	}

	for {
		choice, err := menu.launchMenu()
		if err != nil {
			return err
		}

		// Skip handling if choice is empty (file movement operations return empty string)
		if choice == "" && menu.Mode != ModeNewNote {
			continue
		}

		err = menu.handleChoice(choice)
		if err != nil {
			return err
		}

		menu.nav.Reload()
	}
}
