package copy

import (
	"errors"
	"fmt"
	"io"
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

func CopyFromTo(options CopyOptions) (CopyResult, error) {

	srcFile, err := os.Open(options.Source)
	if err != nil {
		return copyResult{
			Source:        options.Source,
			Destination:   options.Destination,
			copiedFiles:   0,
			copiedFolders: 0,
		}, ErrSourcePath
	}
	dstFile, err := os.Create(options.Destination)
	if err != nil {
		return copyResult{
			Source:        options.Source,
			Destination:   options.Destination,
			copiedFiles:   0,
			copiedFolders: 0,
		}, ErrDestinationPath
	}

	buf := make([]byte, chunkSize)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			// TODO: Error reading. we should retry
			fmt.Println(err)
			break
		}
		if n == 0 {
			break // completed
		}

		_, err = dstFile.Write(buf[:n])
		if err != nil {
			// TODO: Error writing. We should retry
			fmt.Println(err)
			break
		}
	}

	return copyResult{
		Source:        options.Source,
		Destination:   options.Destination,
		copiedFiles:   1,
		copiedFolders: 0,
	}, nil
}
