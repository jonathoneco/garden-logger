# Code Review: Garden Logger MVP

## Executive Summary

The Garden Logger MVP successfully delivers its core functionality with a solid architectural foundation. The code is well-organized, builds cleanly, and provides real value as a dmenu_notes replacement. However, there are several areas that need attention before moving beyond MVP status, particularly around testing, error handling, and Go idioms.

## High Level Assessment

### Strengths
- **Clear Architecture**: Clean separation between main entry point and internal logic
- **Good Abstractions**: Well-defined domain models (`Entry`, `Directory`, `MenuState`)
- **Pragmatic Technology Choices**: Smart use of rofi for UI while keeping core logic in Go
- **Focused Scope**: MVP delivers core functionality without feature creep
- **Integration-First**: Works with existing tools rather than reinventing

### Areas for Improvement
- Some Go idiom violations that impact maintainability
- Error handling could be more robust
- Hard dependencies on external tools not well documented

## Resolved Issues ✅

The following critical issues have been addressed:
- ✅ **Index validation bug fixed** - Proper error messages and validation logic
- ✅ **Config struct implemented** - Eliminates global variables
- ✅ **Unreachable code removed** - Clean control flow in move operations

## Remaining Go Idiom Violations

### 1. Library Calling os.Exit()
**Risk Level: HIGH** | **File:** `internal/cmd.go:68,86`

Libraries should never call `os.Exit()` - only `main()` should control program termination.

**Current Issue:**
```go
func (menu *MenuState) LaunchDir(dirPath string) error {
    // ... launch logic ...
    os.Exit(0) // ❌ Library shouldn't exit
    return nil
}
```

**Recommended Fix:**
```go
// Create a special error type to signal successful launch
type LaunchSuccessError struct {
    Message string
}
func (e LaunchSuccessError) Error() string { return e.Message }

func (menu *MenuState) LaunchDir(dirPath string) error {
    // ... existing launch logic ...
    err := cmd.Start()
    if err != nil {
        return fmt.Errorf("failed to launch kitty: %w", err)
    }
    return LaunchSuccessError{Message: "directory editor launched successfully"}
}

// In main.go:
if err := internal.StartApp(); err != nil {
    var launchErr internal.LaunchSuccessError
    if errors.As(err, &launchErr) {
        os.Exit(0) // Success - program launched editor
    }
    slog.Error("Application Error", "error", err)
    os.Exit(1)
}
```

### 2. Inconsistent Receiver Names
**Risk Level: LOW** | **Files:** Throughout `internal/`

Go convention is to use short, consistent receiver names (1-2 characters).

**Current Issues:**
```go
func (dir *Directory) LoadEntry(...)     // ❌ Should be 'd'
func (menu *MenuState) formatStatusMessage() // ❌ Should be 'm'
func (entry *Entry) String()             // ❌ Should be 'e'
```

**Fix:**
```go
func (d *Directory) LoadEntry(...)
func (m *MenuState) formatStatusMessage()
func (e *Entry) String()
```

### 3. Magic Numbers and Constants
**Risk Level: LOW** | **File:** `internal/modes.go`

**Current:**
```go
switch exitError.ExitCode() {
case 10: // kb-custom-1 (Ctrl+Alt+J) - move down
case 11: // kb-custom-2 (Ctrl+Alt+K) - move up
```

**Fix:**
```go
const (
    RofiExitCodeMoveDown = 10 // Ctrl+Alt+J
    RofiExitCodeMoveUp   = 11 // Ctrl+Alt+K
)
```

## Architecture Improvement: Fix "God Object" Pattern

### Current Problem
`MenuState` has become a "god object" handling too many responsibilities:
- UI state management
- File operations (`writeNote`, `LoadDirectory`)
- Command execution (`LaunchDir`, `launchNote`) 
- Menu rendering

### Recommended Solution: Simple Service Pattern

**Create a focused service for note operations:**
```go
// internal/notes.go
type NotesService struct {
    config *Config
}

func NewNotesService(config *Config) *NotesService {
    return &NotesService{config: config}
}

func (s *NotesService) LoadDirectory(dirPath string) (*Directory, error) {
    absPath := filepath.Join(s.config.RootDir, dirPath)
    // ... move current LoadDirectory logic here
}

func (s *NotesService) CreateNote(dir *Directory, entry *Entry) (string, error) {
    var targetDir string
    if dir.Path == "" {
        targetDir = filepath.Join(s.config.RootDir, s.config.InboxDir)
    } else {
        targetDir = dir.AbsPath
    }
    // ... move current writeNote logic here
}

func (s *NotesService) LaunchNoteEditor(filePath string) error {
    // ... move current launchNote logic here
}

func (s *NotesService) LaunchDirectoryEditor(dirPath string) error {
    // ... move current LaunchDir logic here
}
```

**MenuState becomes pure UI state:**
```go
type MenuState struct {
    Dir       *Directory
    Mode      Mode
    Selection string
    notes     *NotesService  // Single clean dependency
}

func (m *MenuState) navigateTo(dirPath string) error {
    dir, err := m.notes.LoadDirectory(dirPath)  // Clean delegation
    if err != nil {
        return err
    }
    m.Dir = dir
    m.Selection = ""
    return nil
}
```

**Benefits:**
- Clear separation of concerns without over-engineering
- No repetitive config passing
- MenuState methods become thin coordination layers
- Easy to understand what each struct does

## Simple Go Idiom Fixes

The following are straightforward improvements that follow Go conventions:

## Comprehensive Naming Convention Fixes

### 1. Receiver Name Consistency
**Current Violations:**
```go
func (dir *Directory) LoadEntry(...)         // ❌ 'dir' too long
func (menu *MenuState) formatStatusMessage() // ❌ 'menu' too long
func (entry *Entry) String()                 // ❌ 'entry' too long
```

**Go Standard Fix:**
```go
func (d *Directory) LoadEntry(...)           // ✅ Single char
func (m *MenuState) formatStatusMessage()    // ✅ Single char  
func (e *Entry) String()                     // ✅ Single char
```

**Apply Across All Files:**
- `files.go`: Replace all `dir` → `d`, `entry` → `e`
- `menu.go`: Replace all `menu` → `m`
- `modes.go`: Replace all `menu` → `m`

### 2. Method Name Inconsistencies

**Current Issues:**
```go
func (menu *MenuState) formatStatusMessage() string  // ❌ private method, inconsistent casing
func (menu *MenuState) LaunchDir(dirPath string)    // ✅ public method
func (menu *MenuState) launchNote(filePath string)  // ❌ private method, inconsistent with LaunchDir
```

**Decision Needed:** Choose consistent approach:

**Option A - Keep Both Private:**
```go
func (m *MenuState) formatStatusMessage() string
func (m *MenuState) launchDir(dirPath string) error    // Make private
func (m *MenuState) launchNote(filePath string) error  // Already private
```

**Option B - Make Launch Methods Public:**
```go
func (m *MenuState) FormatStatusMessage() string       // Make public
func (m *MenuState) LaunchDir(dirPath string) error    // Already public
func (m *MenuState) LaunchNote(filePath string) error  // Make public  
```

### 3. Variable Naming Improvements

**Current Unclear Names:**
```go
dirTrip := false  // ❌ What does this mean?
```

**Descriptive Alternatives:**
```go
foundFirstFile := false     // ✅ Clear intent
hasSeenFile := false        // ✅ Clear state
filesStarted := false       // ✅ Clear boundary
```

### 4. Function Naming Patterns

**Current Generic Names:**
```go
func getNewMenuItems() ([]string, error)
func (menu *MenuState) getMenuItems() ([]string, error)
func (menu *MenuState) handleChoice(choice string) error
```

**More Descriptive:**
```go
func buildNewModeMenuItems() ([]string, error)         // ✅ Describes action
func (m *MenuState) buildCurrentModeMenuItems() ([]string, error)  // ✅ Specific
func (m *MenuState) processMenuSelection(choice string) error       // ✅ Clear purpose
```

## Simple Error Handling Improvements

Just improve error messages with more context - no complex error types needed:

**Current Generic Errors:**
```go
return fmt.Errorf("Index Validation Error Directory after File")
return fmt.Errorf("[ERROR] unexpected choice: %s", choice)
```

**Simple Improvements:**
```go
return fmt.Errorf("validation failed: found directory after file in %s", d.Path)
return fmt.Errorf("unexpected menu choice %q in %s mode", choice, m.Mode)
```

**Pattern to Follow:**
- Add context (which directory, which mode, which operation)
- Use `%w` for error wrapping when chaining calls
- Keep error messages user-readable

## Package Structure - Keep It Simple

**Current structure is fine for this project size.** Only consider splitting if files get unwieldy (>500 lines).

**If you do split later, keep it minimal:**
```
internal/
├── config.go     # Config, constants, environment  
├── models.go     # Entry, Directory, MenuState types
├── indexing.go   # All indexing operations
├── menu.go       # UI and navigation
└── app.go        # Application startup
```

Don't create deep hierarchies or multiple interface files for a personal tool.

## Action Items Summary

### High Priority (Architecture)
1. **Extract NotesService**: Move file/command operations out of MenuState
2. **Remove os.Exit() calls**: Implement `LaunchSuccessError` pattern

### Medium Priority (Go Idioms)  
3. **Receiver naming**: Global find/replace `(dir *Directory)` → `(d *Directory)` etc.
4. **Constants extraction**: Move magic numbers to `const` blocks
5. **Error message enhancement**: Add context to existing generic errors
6. **Variable naming**: Replace `dirTrip` with `foundFirstFile`

### Low Priority (Polish)
7. **Method visibility consistency**: Align `launchNote` with `LaunchDir`
8. **Add godoc comments**: Document exported functions only

The recent fixes have significantly improved code quality. These remaining items will complete the transformation into idiomatic, maintainable Go code ready for future expansion.
