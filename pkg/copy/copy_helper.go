package copy

import (
	"fmt"
	"io"
	"os"
)

// Copy chunks from source to destination
func copyChunks(source *os.File, destination *os.File) error {
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
