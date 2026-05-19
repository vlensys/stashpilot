package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vlensys/stashpilot/internal/git"
	"github.com/vlensys/stashpilot/internal/tui"
)

func main() {
	if !git.IsRepo() {
		fmt.Fprintln(os.Stderr, "stashpilot: not a git repository")
		os.Exit(1)
	}

	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
