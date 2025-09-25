package copy

import (
	"errors"
	"fmt"
	"os"
)

type CopyOptions struct {
	Source      string
	Destination string
	Override    bool
}

// TODO: Change this to have dynamic chunk sizes
const chunkSize = 1024 * 1024 // 1MB

var ErrSourcePath = errors.New("the source path is invalid")
var ErrDestinationPath = errors.New("the destination path is invalid")
var ErrUnexpected = errors.New("Unexpected error")

var osOpen = os.Open
var osCreate = os.Create

func CopyFromTo(options CopyOptions) (CopyResult, error) {

	result := copyResult{
		Source:        options.Source,
		Destination:   options.Destination,
		copiedFiles:   0,
		copiedFolders: 0,
	}

	srcFile, err := osOpen(options.Source)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ERROR %s %s\n", ErrUnexpected.Error(), r)
		}
		srcFile.Close()
	}()

	if err != nil {
		return result, ErrSourcePath
	}
	dstFile, err := osCreate(options.Destination)
	if err != nil {
		return result, ErrDestinationPath
	}

	err = copyChunks(srcFile, dstFile)
	if err != nil {
		return result, err
	}

	result.copiedFiles++

	return result, nil
}
