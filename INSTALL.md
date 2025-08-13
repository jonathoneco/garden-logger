# Installation and Usage Guide

## Quick Start

Your Garden Logger is ready to use! Here's how to get it running as a dmenu_notes replacement.

## Installation

### Prerequisites

- dmenu
- kitty terminal
- nvim
- Go 1.19+ (for building)

```bash
# Install prerequisites (Arch Linux)
sudo pacman -S dmenu kitty neovim

# Install prerequisites (Ubuntu/Debian)
sudo apt install dmenu kitty neovim
```

### Build and Install

```bash
# You're already in the project directory
cd /home/jonco/src/garden-logger

# Build (already done!)
go build ./cmd/garden-logger

# Install for easy access
sudo cp garden-logger /usr/local/bin/
# Or just for your user
mkdir -p ~/.local/bin
cp garden-logger ~/.local/bin/
```

## Configuration

Set your environment variables (add to ~/.bashrc or ~/.zshrc):

```bash
# Your notes directory (update path as needed)
export GARDEN_LOG_DIR="$HOME/src/garden-log"

# Optional: Terminal preference  
export TERMINAL="kitty"
```

Then reload your shell:
```bash
source ~/.bashrc  # or source ~/.zshrc
```

## Usage

### Running the Application

```bash
# If installed system-wide
garden-logger

# Or run directly from project
./garden-logger

# Or create an alias
alias notes='garden-logger'
```

### How It Works (Current Implementation)

1. **Start**: Opens dmenu with current directory contents
2. **Navigation**:
   - Directories shown with trailing `/`
   - `.md` files listed after directories
   - Choose "New" to create a note
   - Choose ".." to go up a directory
3. **Creating Notes**:
   - Enter a name when prompted
   - Press Escape for automatic date-based naming (YYYY-MM-DD)
   - Notes automatically get indexed (1 Note.md, 2 Another.md, etc.)
   - New notes go to "1 Inbox/" if you're in the root

### Key Features Working Now

✅ **dmenu integration** - Browse with dmenu interface  
✅ **Directory navigation** - Navigate folders, directories listed first  
✅ **Automatic indexing** - Notes get numbered automatically  
✅ **Kitty + nvim integration** - Opens notes in Kitty terminal with nvim  
✅ **Date-based naming** - Defaults to current date format  
✅ **Markdown support** - Filters for .md files  
✅ **Inbox system** - New notes go to "1 Inbox/" by default  

### Example Workflow

```bash
# Set up your notes directory
export GARDEN_LOG_DIR="$HOME/src/garden-log"
mkdir -p "$GARDEN_LOG_DIR"

# Run the application
garden-logger

# In dmenu:
# 1. Choose "New" 
# 2. Type "Daily Standup" (or press Escape for date)
# 3. Note gets created as "1 Inbox/1 Daily Standup.md"
# 4. Opens in Kitty with nvim
```

### Integration with Your System

Since you mentioned using this as a dmenu_notes replacement, you can:

1. **Create a keyboard shortcut** in your window manager to run `garden-logger`
2. **Replace dmenu_notes** by aliasing or symlinking
3. **Use with existing workflow** - it respects your PARA structure

### Directory Structure (Matches Your Spec)

```
$HOME/src/garden-log/
├── 1 Inbox/
│   ├── 1 Daily Standup.md
│   ├── 2 2024-08-13.md
│   └── 3 Quick Ideas.md
├── 2 Projects/
│   ├── 1 Garden Logger.md
│   └── 2 Website Redesign.md
├── 3 Areas/
├── 4 Archive/
└── 5 Resources/
    └── 1 Templates/
```

## Customization

### Change Date Format

Edit `internal/app.go`, line 303:
```go
name = time.Now().Format("2006-01-02")  // Change this
```

Examples:
- `"2006-01-02"` → "2024-08-13"
- `"Jan-02-2006"` → "Aug-13-2024"
- `"2006-01-02_15-04"` → "2024-08-13_14-30"

### Change Colors

Edit the dmenu color args in `internal/app.go`:
```go
args := []string{"-c", "-l", "10", "-i", "-p", prompt, "-sb", "#your-color"}
```

### Change Editor

Modify `openInKitty()` function in `internal/app.go` to use different editor/terminal.

## Troubleshooting

### Notes directory not found
```bash
# Make sure the directory exists and environment variable is set
echo $GARDEN_LOG_DIR
mkdir -p "$GARDEN_LOG_DIR"
```

### dmenu not working
```bash
# Test dmenu directly
echo -e "test1\ntest2" | dmenu
```

### Kitty not opening
```bash
# Test kitty directly
kitty --version
```

## Next Steps

Your current implementation covers the core dmenu_notes functionality. Based on your spec, you might want to add:

- Templates support (`createFromTemplate()` function)
- Better error handling for edge cases
- Configuration file support
- More indexing strategies (datetime, none)
- Global search functionality

The foundation is solid for building out your full vision!