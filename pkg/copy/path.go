package copy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/paolinos/gcopy/pkg/analyzer"
)

var filepathDir = filepath.Dir

// Copy all files from folder
func copyFolderData(source string, destination string, files []analyzer.FileInfo) error {

	err := checkOrCreateFolder(destination)
	if err != nil {
		return err
	}

	for _, file := range files {
		filename := strings.Replace(file.Path, source, "", 1)
		dstPath := filepath.Join(destination, filename)
		err := copyChunksFromSource(file.Path, dstPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(source string, destination string) error {

	destinationFolder := filepathDir(destination)
	err := checkOrCreateFolder(destinationFolder)
	if err != nil {
		return err
	}

	err = copyChunksFromSource(source, destination)
	if err != nil {
		return err
	}
	return nil
}

// Copy path using the analysis
func CopyPath(data analyzer.AnalyzeResult) {

	if len(data.Folders) == 0 && data.TotalFiles == 1 {
		// Copy file from to
		err := copyFile(data.Source, data.Destination)
		if err != nil {
			fmt.Printf("Error trying to copy file from:%s; to:%s; with error:%s", data.Source, data.Destination, err)
		}
		return
	}

	for _, folder := range data.Folders {
		folderName := filepath.Base(folder.Path)

		d := strings.Replace(folder.Path, data.Source, data.Destination, 1)
		// TODO: remove fmt.Printf
		fmt.Printf(" - Folder -> Full path:%s; Folder name:%s; Destination path: %s\n", folder.Path, folderName, d)
		err := copyFolderData(folder.Path, d, folder.Files)
		if err != nil {
			fmt.Printf("Error trying to copy from:%s; to:%s; with error:%s", folder.Path, d, err)
		}
	}
}
