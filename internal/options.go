package internal

import (
	"flag"
	"fmt"
	"strings"
)

// GetHelper returns a formatted helper message with version info.
func GetHelper(version, description string) string {
	return fmt.Sprintf("gcopy %s - %s", version, description)
}

// AppOptions holds all parsed application options.
type AppOptions struct {
	Source      string
	Destination string
	Chunks     int64
	Check       bool
	Skip        bool
	Force       bool
	Show        int
	Parallel    int
}

var errMissingArguments = fmt.Errorf("missing arguments")

// GetAppOptions parses command-line flags and returns the populated AppOptions.
func GetAppOptions() (*AppOptions, error) {
	o := &AppOptions{}

	flag.StringVar(&o.Source, "source", "", "Source path to copy from.")
	flag.StringVar(&o.Destination, "destination", "", "Destination path to copy to.")
	flag.Int64Var(&o.Chunks, "chunks", 32*1024, "Number of chunks for parallel copying.")
	flag.BoolVar(&o.Check, "check", false, "Validate each file's source and destination sizes after copy.")
	flag.BoolVar(&o.Skip, "skip", false, "Skip existing files at the destination without overwriting.")
	flag.BoolVar(&o.Force, "force", false, "Overwrite all existing files at the destination.")
	flag.IntVar(&o.Show, "show", 5, "Number of files to display in results summary (default: 5).")
	flag.IntVar(&o.Parallel, "parallel", 5, "Max concurrent file transfers (default: 5).")

	flag.Parse()

	// If no flags were set but positional args exist, use them as source/destination.
	if o.Source == "" && o.Destination == "" {
		args := flag.Args()
		if len(args) >= 2 {
			o.Source = args[0]
			o.Destination = args[1]
		} else if len(args) == 1 {
			o.Source = args[0]
		}
	}

	// Validate that source and destination are provided.
	if o.Source == "" || o.Destination == "" {
		return nil, fmt.Errorf("missing arguments")
	}

	// Check if options come after positional arguments (invalid usage).
	args := flag.Args()
	for i, arg := range args {
		if i >= 2 && !strings.HasPrefix(arg, "-") {
			return nil, fmt.Errorf("options after source and destination are not allowed")
		}
	}

	return o, nil
}
