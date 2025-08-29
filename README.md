# Garden Logger

I have a personal note taking system I call my Garden Log, this tool is one I'm writing to help manage and use my note taking system. I keep my notes in raw markdown that I edit with neovim where I can, and obsidian where I can't. The reason I'm writing this tool is that I often find navigation to and interaction with my notes often randomizes me, especially with Notion which is what I am currently using. I want this tool to reduce the friction of getting to where I need to be as much as possible

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
- Fix directory creation
- Add support for `.` to open the folder in it's own tmux session
- Turn off obsidian.nvim frontmatter
- Add support for vim marks
- Add support for global search
    - fzf file path
- Recent Notes support

- Is it worth using linked lists for the indexed directories?

- Folder priority as bool in indexing strategy

- Add logging

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

## Docs
- Update to rofi info
    - Provide info about rofi-launcher script

## Templates
- Template mode for menu
- Default title, with option to change

## Logging
- Fix debug logging

## Install

### Requirements

TBD

## Usage

TBD

## Contributing

Not accepting PRs

## License

MIT © Jonathon Corrales de Oliveira
