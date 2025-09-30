package copy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type dirInfo struct {
	files   []string
	folders []string
}

type pathCopyResult struct {
	files   int
	folders int
	paths   []string
	errors  []string
}

func (p pathCopyResult) HasError() bool {
	return len(p.errors) > 0
}

// Read and copy, is going to check the path if is a directory or a file to copy and start to copy
func readAndCopy(source *os.File, destination *os.File) pathCopyResult {
	res := pathCopyResult{
		files:   0,
		folders: 0,
		paths:   []string{},
		errors:  []string{},
	}
	fileinfo, err := source.Stat()
	if err != nil {
		res.errors = append(res.errors, err.Error())
		return res
	}

	if !fileinfo.IsDir() {
		// Copy the file
		err = copyChunksFromSource(source, destination)
		if err != nil {
			res.errors = append(res.errors, err.Error())
		} else {
			res.files++
		}
		return res
	}

	// TODO: review
	copyFolder(source.Name(), destination.Name())

	return res
}

// Copy files from source to destination folder. return list of errors
func copyFolderFiles(sourceFolder string, destinationFolder string, files []string) []string {
	var errors []string

	err := checkOrCreateFolder(destinationFolder)
	if err != nil {
		errors = append(errors, err.Error())
		return errors
	}

	// Copy all the files from the folder
	for _, file := range files {
		srcPath := filepath.Join(sourceFolder, file)
		dstPath := filepath.Join(destinationFolder, file)
		err := copyChunksFromPath(srcPath, dstPath)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	return errors
}

func copyChunksFromPath(srcFilePath string, dstFilePath string) error {
	source, err := osOpen(srcFilePath)
	if err != nil {
		return err
	}
	destination, err := osCreate(dstFilePath)
	if err != nil {
		return err
	}

	return copyChunksFromSource(source, destination)
}

// Copy chunks from source to destination
func copyChunksFromSource(source *os.File, destination *os.File) error {
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
