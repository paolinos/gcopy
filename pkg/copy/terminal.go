package copy

import (
	"os"

	"golang.org/x/term"
)

const (
	defaultTermWidth  = 80
	progressBarWidth  = 50
	filenameMinWidth  = 20
	headerPadding    = 2
	sizeColumnWidth   = 10
	speedColumnWidth  = 12
	pctColumnWidth    = 5
	timeColumnWidth   = 7
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
