# Garden Logger

My personal notes I call my Garden Log. I was keeping my notes in notion for a while but found navigation to and interaction with my notes often randomizes me. I am swapping to keeping my notes in raw markdown that I edit with neovim where I can, and obsidian where I can't. I'm writing this tool for two reasons:
1) I want to reduce the friction of getting to where I need to be as much as possible
2) I'd like to learn Golang and this seems like a good first project for it.

My system is based on the PARA model, currently I've got the following folder structure

$HOME/src/garden-log
├── 1 Inbox
├── 2 Projects
├── 3 Areas
├── 4 Archive
└── 5 Resources
    └─ 1 Templates

Inspired by https://github.com/BreadOnPenguins/scripts/blob/master/dmenu_notes
I originally modified it a bit to work with my personal system and subdirectories, and eventually just wanted a lot more functionality than it made sense to implement in bash scripts

## Spec / Features
### Configuration
- Root Directory: Directory where my notes are stored
- Inbox Directory: Where to put "quick notes"
- Template Directory: Template Directory


### Note Management
- Directory Navigation within my root Notes Directory
- Create a Directory, Note, or Note from a Template
- Open a selected note in Neovm
- Open a selected directory in a tmux session
- Unnamed notes are titled with the current date (for more easily logging things like daily logs or saxophone practice)

### Indexing
- Setting toggle for wheter or not to index the current directory
- Support for reordering indexed entries
- Dir-priority indexing, directories are sorted to the top

### CLI Entry Point
- CLI Entry point to enable use of the indexing and quality of life functionality from scripts or keyboard shortcuts

## Dependencies
This is intended to be integrated to my personal work machine, I wanted a quick and dirty tool to make my life easier so I rely heavily on some other quality of life tools specific to my environment. For more information about how I've configured these, take a look at my dotfiles repo
- Obsidian
    - This isn't strictly necessary for this app but it's what I use to interact with my synced notes outside of my workstation
- Neovim
    - I use this for text editing at my workstation, gives me:
        - syntax highlighting
        - buffer and directory scoped fuzzy finding
        - obsidian tooling integration with obsidian.nvim
- Tmux
    - I use this for launching directory specific sessions if I want to do more than just take notes on a particular subject
- Rofi
    - Provides a simple, extensible UI
    - I use different config files for different rofi use-cases so I have a 'rofi-launcher' script in my dotfiles, that's what this project actually calls when launching rofi
    - Part of my notes rofi configs includes custom keybinds for ctrl-j / ctrl-k for entry navigation and ctrl-alt-j / ctrl-alt-k for indexed entry re-ordering

## Follow Ups

- Improved Renaming
    - Add support for renaming in the UI
    - Currently, if I change a directories indexing setting, or reorder indexed files, it breaks links to those files
    - I can launch a headless neovim instance and use `:ObsidianRename` to rename files while preserving links to them
        - I should see what other obsidian operations I may want to support that I can take this approach for
    - I don't use links very often at the moment so this isn't a super high priority
- Improved Index Validation / Repair
    - Right now if index validation fails it just errors
    - I want to add some sort of indexing repair funcionality so I can delete files without too much concern
- Deletion
    - Maybe support deletion
    - This is meant to be a quick-touch tool so I don't want to accidentally delete files or directories while moving around
    - Some sort of "type the note or directories name to confirm deletion" interface
- Templates
    - Support for frontmatter with configurable / auto-populated fields
    - Might be nice to have template defaults for some directories
- Sync Surface
    - Haven't yet decided how I'm actually syncing my notes across surfaces, I want something responsive that handles offline edits well
- Across Layer Note Movement
    - Seems like a weird edge case, but might be a nice to have
- Dependencies
    - If this gets any attention at all for some reason, I may try to abstract the dependencies a little better so it works for other people, with some sort of config file, but not a priority for me at all right now


# Immediate TODOs
- Some QOL logging

### Requirements

TBD

## Usage

TBD

## Contributing

Not accepting PRs

## License

MIT © Jonathon Corrales de Oliveira
