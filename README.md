# Garden Logger

I have a personal note taking system I call my Garden Log, this tool is one I'm writing to help manage and use my note taking system. I keep my notes in raw markdown that I edit with neovim where I can, and obsidian where I can't.

My system is based on the PARA model, currently I've got the following folder structure

$HOME/src/garden-log
├── 1 Inbox
├── 2 Projects
├── 3 Areas
├── 4 Archive
└── 5 Resources
    └─ 1 Templates

I'm also using this as an excuse to learn Go.

Inspired by https://github.com/BreadOnPenguins/scripts/blob/master/dmenu_notes
I originally modified it a bit to work with my personal system and subdirectories, and eventually just wanted a lot more functionality than it made sense to implement in bash scripts

## Spec / Features
This is the intended feature spec for this tool

### Note Creation / Management

Minimal ui, inspired by dmenu
By default starts with root in my notes directory
Should be able to Ctrl-C or Escape at any step with no effect

Entries:
- A "New" entry at the top
- A "New from Template" entry
- Each item in the current directory
- A `..` entry at the bottom

Selection:
- New: Swaps to ui for creating a new note or directory (determined by whether the entry ends in a /)
- New From Template: Show's available templates, upon selection of one goes to the New Entry UI (with some flavor text about the template being created)
- A note: Opens the note in neovim in a new terminal instance, with my notes directory as the root
- A directory: Navigates to that directory in the UI
- `..`: Navigates to the parent directory in the UI

Indexing:
- Everything inside a directory is 1-indexed
- Auto indexes when creating a note or directory
- Directories and files share indexes, i.e. if a directory contains a directory and a file, it would look like
1 Example Directory
├── 1 Example Subdirectory
└── 2 Example Note.md

Field:
- FZF for the entries in the directory

"New Entry" UI:
- If entering clean note / directory name, the entry get's auto-populated to the bottom
- If an index is provided, the entry gets inserted and shuffles everything else down
- If after selecting a template, can only enter a filename for the note, not a directory

Updates:
- Support for reordering with automatic index updates
- Moving/Renaming a note updates it's header and any links to it
- Moving/Renaming a folder updates the semantic id's of all it's subnotes and any links to those semantic ids

Templates:
- Support for some auto-populated properties i.e. created-date, grow with need

### Questions
- Integration
    - Obsidian
        - I'm only really doing the header thing because I thought I'd need to handle link update propogation manually outside of obsidian
        - If I can make use of obsidian's linking and update propagation I don't need them on every note and can just add properties in templates where relevant, i.e. created time for saxophone practice logs
    - Git
        - For now using git to save my notes remotely, want to figure out a better solution to track across my machines and my phone
- There are some cases where this indexing doesn't make sense, i.e. for saxophone practice I want my more recent sessions up top so I'd want to index by date-time, I should have some sort of way to enable this
- Looks like I'll need some sort of config to declare
    - where my root notes are
    - directories where I don't want auto indexing (maybe I actually use some sort of hidden file for this? like a .noindex that way it persists across moves)
    - preferred editor
- Not sure how to handle moving files / directories across layers, some weird edge cases:
    - moving file/directory from unindexed directory to an indexed one, and vice versa
    - moving file/directory in between indexed directories, if there's overlap new file take's priority, if not (old index larger than max index) it gets added to the bottom
    - moving file/directory in between unindexed directories
- Concurrency
    - Should I create some sort of lock, not sure if I need to worry about concurrency
    - Only case I really think this applies is if I try to manipulate notes open in some instance of neovim, not sure what to do about this yet
- File watching
    - The tool should not be constantly running, it should respond to the state of my notes on start
- Error Recovery
    - Mostly worried about corrupt indexing
- Other filetypes
    - I'll want to support things like images and other stuff I do in my notes, not sure what the best way to do that with all this indexing is. Maybe just don't index those?

# Follow Ups
- Add support for `.` to open the folder in it's own tmux session
- Turn off obsidian.nvim frontmatter
- Add support for vim marks
- Add support for global search
    - fzf file path
- Recent Notes support

- Is it worth using linked lists for the indexed directories?

- Folder priority as bool in indexing strategy

# Configuration
- note root folder
- default indexing option,can also set per-directory index strategy with a .indexing file dictating the strategy

# Indexing
- Numeric
- Datetime
- None



## Install

### Requirements

TBD

## Usage

TBD

## Contributing

Not accepting PRs

## License

MIT © Jonathon Corrales de Oliveira
