package copy

import (
	"fmt"
	"io"
	"os"
)

type CopyOptions struct {
	Source      string
	Destination string
	Override    bool
}

var osOpen = os.Open
var osCreate = os.Create
var osStat = os.Stat
var osMkdir = os.Mkdir
var osIsNotExist = os.IsNotExist

// Copy chunks from source to destination
func copyChunksFromSource(srcFilePath string, dstFilePath string, chunkSize int64) error {

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
