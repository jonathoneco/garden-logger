package main

import (
	"errors"
	"flag"
	"fmt"
	"garden-logger/internal"
	"log/slog"
	"os"
)

func main() {
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.Parse()

	internal.InitLogger(verbose)

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	if err := handleCommand(args); err != nil {
		var launchErr internal.LaunchSuccessError
		if errors.As(err, &launchErr) {
			os.Exit(0)
		}
		slog.Error("CLI Error", "error", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: garden-logger-cli <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  new          Create a new note in inbox with current date")
	fmt.Println("  open <path>  Open note at specified path")
}

func handleCommand(args []string) error {
	command := args[0]

	switch command {
	case "new":
		return handleNewCommand()
	case "open":
		if len(args) < 2 {
			return fmt.Errorf("open command requires a path argument")
		}
		return handleOpenCommand(args[1])
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func handleNewCommand() error {
	config, err := internal.LoadConfig()
	if err != nil {
		return err
	}

	notes := internal.NewNotesService(config)
	nav := internal.NewNavigator(notes)

	err = nav.NavigateTo("")
	if err != nil {
		return err
	}

	err = nav.NavigateTo(config.InboxDir)
	if err != nil {
		return err
	}

	filePath, err := notes.CreateEntryFromUserInput(nav.CurrentDirectory(), "", false)
	if err != nil {
		return err
	}

	return notes.LaunchNoteEditor(filePath)
}

func handleOpenCommand(path string) error {
	config, err := internal.LoadConfig()
	if err != nil {
		return err
	}

	notes := internal.NewNotesService(config)
	return notes.LaunchNoteEditor(path)
}