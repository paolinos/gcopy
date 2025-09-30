package copy

import (
	"errors"
	"os"
)

type CopyOptions struct {
	Source      string
	Destination string
	Override    bool
}

// TODO: Change this to have dynamic chunk sizes
const chunkSize = 1024 * 1024 * 100 // 100MB

var ErrSourcePath = errors.New("the source path is invalid")
var ErrDestinationPath = errors.New("the destination path is invalid")
var ErrUnexpected = errors.New("Unexpected error")

var osOpen = os.Open
var osCreate = os.Create
var osStat = os.Stat
