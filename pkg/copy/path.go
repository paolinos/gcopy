package copy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/paolinos/gcopy/pkg/analyzer"
)

var filepathDir = filepath.Dir

func copyFolderData(source string, destination string, files []analyzer.FileInfo, chunkSize int64) error {
	err := checkOrCreateFolder(destination)
	if err != nil {
		return err
	}

	for _, file := range files {
		filename, _ := filepath.Rel(source, file.Path)
		dstPath := filepath.Join(destination, filename)
		fp := &fileProgress{}
		err := copyChunksFromSource(file.Path, dstPath, chunkSize, fp)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(source string, destination string, chunkSize int64) error {
	destinationFolder := filepathDir(destination)
	err := checkOrCreateFolder(destinationFolder)
	if err != nil {
		return err
	}

	fp := &fileProgress{}
	err = copyChunksFromSource(source, destination, chunkSize, fp)
	if err != nil {
		return err
	}
	return nil
}

func CopyPath(data analyzer.AnalyzeResult, chunkSize int64) {
	if len(data.Folders) == 0 && data.TotalFiles == 1 {
		err := copyFile(data.Source, data.Destination, chunkSize)
		if err != nil {
			fmt.Printf("Error trying to copy file from:%s; to:%s; with error:%s", data.Source, data.Destination, err)
		}
		return
	}

	for _, folder := range data.Folders {
		folderName := filepath.Base(folder.Path)
		d := strings.Replace(folder.Path, data.Source, data.Destination, 1)
		fmt.Printf(" - Folder -> Full path:%s; Folder name:%s; Destination path: %s\n", folder.Path, folderName, d)
		err := copyFolderData(folder.Path, d, folder.Files, chunkSize)
		if err != nil {
			fmt.Printf("Error trying to copy from:%s; to:%s; with error:%s", folder.Path, d, err)
		}
	}
}
