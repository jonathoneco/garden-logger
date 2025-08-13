#!/bin/bash

# Garden Logger Demo Script
# This shows how to set up and use garden-logger

set -e

echo "🌱 Garden Logger Setup Demo"
echo "============================"

# Check if dmenu is available
if ! command -v dmenu &> /dev/null; then
    echo "❌ dmenu not found. Please install dmenu first:"
    echo "   Arch Linux: sudo pacman -S dmenu"
    echo "   Ubuntu/Debian: sudo apt install dmenu"
    exit 1
fi

# Check if kitty is available
if ! command -v kitty &> /dev/null; then
    echo "⚠️  kitty not found. You can still use garden-logger but will need to modify the editor setup."
    echo "   Install kitty: sudo pacman -S kitty (Arch) or sudo apt install kitty (Ubuntu/Debian)"
fi

# Check if nvim is available  
if ! command -v nvim &> /dev/null; then
    echo "⚠️  nvim not found. You can still use garden-logger but will need to modify the editor setup."
    echo "   Install neovim: sudo pacman -S neovim (Arch) or sudo apt install neovim (Ubuntu/Debian)"
fi

echo
echo "📁 Setting up demo environment..."

# Create a demo notes directory
DEMO_DIR="$HOME/demo-garden-log"
export GARDEN_LOG_DIR="$DEMO_DIR"

# Clean up any existing demo
if [ -d "$DEMO_DIR" ]; then
    echo "🧹 Cleaning up existing demo directory..."
    rm -rf "$DEMO_DIR"
fi

# Create demo structure
mkdir -p "$DEMO_DIR"
mkdir -p "$DEMO_DIR/1 Inbox"
mkdir -p "$DEMO_DIR/2 Projects" 
mkdir -p "$DEMO_DIR/3 Areas"
mkdir -p "$DEMO_DIR/4 Archive"
mkdir -p "$DEMO_DIR/5 Resources/1 Templates"

# Create some demo files
cat > "$DEMO_DIR/1 Inbox/1 Welcome to Garden Logger.md" << 'EOF'
# Welcome to Garden Logger

This is a demo note created by the setup script.

Garden Logger is your dmenu-based note-taking companion that helps you:

- 📝 Create and organize notes with automatic indexing
- 📂 Navigate through directories with ease
- ⚡ Quick access via dmenu interface
- 🎯 Focus on writing, not organization

## Quick Tips

- Press **Enter** to select an item in dmenu
- Press **Escape** to cancel/go back
- Directories are shown with trailing `/`
- Type to filter items in dmenu

Happy note-taking! 🌱
EOF

cat > "$DEMO_DIR/2 Projects/1 Garden Logger Development.md" << 'EOF'
# Garden Logger Development

## Current Status
- ✅ Basic dmenu integration  
- ✅ Directory navigation
- ✅ Automatic indexing
- ✅ Note creation with date fallback
- ✅ Kitty + nvim integration

## Next Features
- 🔄 Template support
- 🔍 Global search
- ⚙️  Configuration files
- 📱 Better mobile sync story

## Notes
This project is a great way to learn Go while building something useful!
EOF

cat > "$DEMO_DIR/5 Resources/1 Templates/daily-standup.md" << 'EOF'
# Daily Standup - {{DATE}}

## Yesterday
- 

## Today  
- 

## Blockers
- 

## Notes
- 
EOF

echo "✅ Demo environment created at $DEMO_DIR"
echo
echo "🔧 Building garden-logger..."
make build > /dev/null 2>&1

echo "✅ Build complete!"
echo
echo "🚀 To try garden-logger:"
echo "   1. Set environment: export GARDEN_LOG_DIR=\"$DEMO_DIR\""
echo "   2. Run: ./garden-logger"
echo
echo "📖 Quick test (will show dmenu with demo content):"
echo "   You should see: New, 1 Inbox/, 2 Projects/, 3 Areas/, etc."
echo
echo "🎯 When you're ready to use with your real notes:"
echo "   export GARDEN_LOG_DIR=\"$HOME/src/garden-log\""
echo "   ./garden-logger"
echo
echo "📚 For installation and full usage instructions, see INSTALL.md"
echo
echo "Demo directory: $DEMO_DIR"
echo "You can safely delete this when done: rm -rf $DEMO_DIR"