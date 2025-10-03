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

const default_chunks_bytes = 1024 * 1024 * 100 // 100MB
const prompt_help = `
GCopy (v %s) tool is used to copy files/folders
%s

Usage:

	gcopy [options] [source] [destination]

Options:
	-chunks				%s
	-o

`

var chunkOptions = fmt.Sprintf("Chunk size to copy files expressed in bytes. Default value:%d(%s)", default_chunks_bytes, analyzer.GetSizeReadable(default_chunks_bytes))
var ErrMissingArguments = errors.New("invalid arguments to run the program")

func GetHelper(version string, description string) string {
	return fmt.Sprintf(prompt_help, version, description, chunkOptions)
}

func GetAppOptions() (AppOptions, error) {
	chunks := flag.Int64("chunks", default_chunks_bytes, chunkOptions)
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
