package copy

import (
	"os"

	"golang.org/x/term"
)

const (
	defaultTermWidth = 80
	filenameWidth    = 20
	sizeColumnWidth  = 10
	speedColumnWidth = 10
	timeColumnWidth  = 6
	minBarWidth      = 10

	// Color codes (used for error messages).
	reset = "\033[0m"
	red   = "\033[31m"
	green = "\033[32m"
	bold  = "\033[1m"
)

func getTerminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return defaultTermWidth
	}
	return w
}

// IsTerminal reports whether stdout is connected to a terminal.
func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
