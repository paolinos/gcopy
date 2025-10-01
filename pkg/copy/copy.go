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
const chunkSize = 1024 * 1024 * 100 // 100MB

var ErrSourcePath = errors.New("the source path is invalid")
var ErrDestinationPath = errors.New("the destination path is invalid")
var ErrUnexpected = errors.New("Unexpected error")

var osOpen = os.Open
var osCreate = os.Create
var osStat = os.Stat

// Copy chunks from source to destination
func copyChunksFromSource(srcFilePath string, dstFilePath string) error {

	source, err := osOpen(srcFilePath)
	if err != nil {
		return err
	}
	destination, err := osCreate(dstFilePath)
	if err != nil {
		return err
	}

	buf := make([]byte, chunkSize)

	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			// TODO: Error reading. we should retry
			fmt.Println(err)
			return err
		}
		if n == 0 {
			return nil
		}

		_, err = destination.Write(buf[:n])
		if err != nil {
			// TODO: Error writing. We should retry
			fmt.Println(err)
			return err
		}
	}
}
