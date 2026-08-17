package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

// AppOptions holds all parsed application options.
type AppOptions struct {
	Source      string
	Destination string
	Chunks      int64
	Check       bool
	Skip        bool
	Force       bool
	Show        int
	Parallel    int
}

// GetHelper returns a formatted helper message with version info.
func GetHelper(version, description string) string {
	return fmt.Sprintf("gcopy %s - %s", version, description)
}

// NewRootCmd creates the root cobra command with all flags.
func NewRootCmd(version, description string) (*cobra.Command, *AppOptions) {
	opts := &AppOptions{}

	cmd := &cobra.Command{
		Use:   "gcopy [flags] <source> <destination>",
		Short: description,
		Long:  GetHelper(version, description),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Source = args[0]
			opts.Destination = args[1]
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64Var(&opts.Chunks, "chunks", 32*1024, "Number of chunks for parallel copying.")
	flags.BoolVar(&opts.Check, "check", false, "Validate each file's source and destination sizes after copy.")
	flags.BoolVar(&opts.Skip, "skip", false, "Skip existing files at the destination without overwriting.")
	flags.BoolVar(&opts.Force, "force", false, "Overwrite all existing files at the destination.")
	flags.IntVar(&opts.Show, "show", 5, "Number of files to display in results summary.")
	flags.IntVar(&opts.Parallel, "parallel", 5, "Max concurrent file transfers.")

	return cmd, opts
}
