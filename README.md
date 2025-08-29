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
    - Currently, if I change a directories indexing setting, or reorder indexed files, it breaks links to those files
    - I can launch a headless neovim instance and use `:ObsidianRename` to rename files while preserving links to them
    - I don't use links very often at the moment so this isn't a super high priority
- Improved Index Validation / Repair
    - Right now if index validation fails it just errors
    - I want to add some sort of indexing repair funcionality so I can delete files without too much concern
- Deletion
    - Maybe support deletion
    - This is meant to be a quick-touch tool so I don't want to accidentally delete files or directories while moving around
    - Some sort of "type the note or directories name to confirm deletion" interface
- Templates
    - I still need to add the template functionality
    - Support for frontmatter with configurable / auto-populated fields
    - Might be nice to have template defaults for some directories
- Sync Surface
    - Haven't yet decided how I'm actually syncing my notes across surfaces, I want something responsive that handles offline edits well
- Across Layer Note Movement
    - Seems like a weird edge case, but might be a nice to have
- Dependencies
    - If this gets any attention at all for some reason, I may try to abstract the dependencies a little better so it works for other people, with some sort of config file, but not a priority for me at all right now


# Immediate TODOs
- Fix directory creation
- Template Functionality
- Some QOL logging
- Update docs below this point

## Rofi
I'm running into some drawbacks with rofi that are starting to get cumbersome
For one, I am creating entity objects to make indexing, files, and dir operations simpler, but rofi slection only works with text so I have to do a lot of reverse engineering whereas being able to have selection return references entity objects directly would be ideal
Also, I'd like some conditional rendering for things like directories vs files
What I like is that I can use common configs for things like my color scheme and font from rofi so that this ui updates with my note changer, and largely I don't have to implement a ui from scratch I can focus on data manipulation
Is there any happy middleground where I can get richer behavior while still getting the benefits from rofi, without having to implement my own ui from scratch


- Create a map of string to entities, use choice as key for how to handle
## Data Modeling
I'm modeling entries more explicitly and separately from current directory state as it makes data operations easier, but the thing I suspect I'll run into is that getting the choice as an entry is going to be a bit of a pain. For now I'll just do so explicitly but if there were a way to get the menu / rofi return to directly associate entries that'd be ideal

The naive approach is just parsing the choice filename and looking for a matching entry in the directory (maybe I don't even need to find the original, I can just create a new one at the point of need?)

I'm running up against the limitations of using rofi as the UI, ideally I can just get richer behavior with object linking but if not I might have to switch to a native ui solution

I'd really rather not though because rofi gives me good fzf behavior out of the box and is stylized with the rest of my system
Is there a way to get rofi to pass around objects and let me handle selection?

## Indexing
- I should swap out the way I look for an index
- I've simplified indexing a lot, instead of json index configs with for one type of indexing, I just touch or delete an empty .index file to indicate whether or not to numerically index my folders
- All the indexing tools should interact with is a list of indices
    - Functions to get the indices in a directory
    - I'll need a way to reverse look up the document so I'll need some sort of ''


## Configuration
- note root folder
- default indexing option,can also set per-directory index strategy with a .indexing file dictating the strategy

## Indexing
- Numeric
- Datetime
- None

## Operation Syncing
- Things I've Considered
    - use an sync file that I append file operations to, then a small obsidian plugin that reads that sync file and performs the vault operations internally
        - links would be broken until obsidian opens
        - once the files are moved the vault operations will fail
    - use obsidian's http rest api for manipulation with link updates
        - obsidian doesn't run headlessly, will require obsidian to already be running or launching it, that feels clunky
    - obsidian-cli or some sync plugins
        - all of them have some poorly handled edge case or other, and I don't like the idea of running into a dependency issue if obsidian upgrades or something
- Indexing poses a problem here, since re-indexing is a rename, that means links need to be updated as frequently / deeply as things are re-indexed, not to mention moves or directory renames
- I think for now, since I don't really use links heavily, I'll skip this and deal with it once I need to
- Eureka! I can use nvim and ObsidianRename for a headless rename, this is huge
- I'm thinking in the process of moving things around there will be a lot of moving (i.e. if I'm moving one document up it'll trigger a bunch of updates for that document and the ones around it), I'll do the naive approach first but I think if this hits performance issues I'll need to swap to a model where I cache updates, execute on ui navigation, and don't perform no-ops



### Requirements

TBD

## Usage

TBD

## Contributing

Not accepting PRs

## License

MIT © Jonathon Corrales de Oliveira
