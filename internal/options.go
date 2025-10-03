package internal

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/paolinos/gcopy/pkg/analyzer"
)

type AppOptions struct {
	Source      string
	Destination string
	Chunks      int64
}

const DEFAULT_CHUNKS_SIZE = 1024 * 1024 * 100 // 100MB
const PROMPT_HELP = `
GCopy (v %s) tool is used to copy files/folders
%s

Usage:

	gcopy [options] [source] [destination]

Options:
	-chunks				%s			

`

var chunkOptions = fmt.Sprintf("Chunk size to copy files expressed in bytes. Default value:%d(%s)", DEFAULT_CHUNKS_SIZE, analyzer.GetSizeReadable(DEFAULT_CHUNKS_SIZE))
var ErrMissingArguments = errors.New("invalid arguments to run the program")

func GetHelper(version string, description string) string {
	return fmt.Sprintf(PROMPT_HELP, version, description, chunkOptions)
}

func GetAppOptions() (AppOptions, error) {
	chunks := flag.Int64("chunks", DEFAULT_CHUNKS_SIZE, fmt.Sprintf("chunk size to copy files expressed in bytes. default value %d (%s)", DEFAULT_CHUNKS_SIZE, analyzer.GetSizeReadable(DEFAULT_CHUNKS_SIZE)))
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		return AppOptions{}, ErrMissingArguments
	}
	// NOTE: args only should be the source and destination
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			return AppOptions{}, ErrMissingArguments
		}
	}

	res := AppOptions{
		Source:      args[0],
		Destination: args[1],
		Chunks:      *chunks,
	}

	return res, nil
}
