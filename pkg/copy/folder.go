package copy

import (
	"fmt"
	"os"
	"path/filepath"
)

// Get information from directory, returning all the files and sub-folders
func getInfoFromDir(source *os.File) (dirInfo, error) {
	res := dirInfo{
		files:   []string{},
		folders: []string{},
	}

	paths, err := source.ReadDir(0)
	if err != nil {
		return res, nil
	}

	for _, e := range paths {
		if e.IsDir() {
			res.folders = append(res.folders, e.Name())
		} else {
			res.files = append(res.files, e.Name())
		}
	}
	return res, nil
}

func copyFolder(sourceFolder string, destinationFolder string) {

	source, err := osOpen(sourceFolder)
	if err != nil {
		// TODO:
		fmt.Println("Error opening source folder", err)
	}

	di, err := getInfoFromDir(source)
	if err != nil {
		// TODO:
		// res.errors = append(res.errors, err.Error())
		fmt.Println("Error get information from directory", err)
		return
	}
	fmt.Println("Prepare to copy files:", di.files, "folders:", di.folders)

	errs := copyFolderFiles(sourceFolder, destinationFolder, di.files)
	if len(errs) > 0 {
		//res.errors = append(res.errors, errs...)
		fmt.Println("Error trying to copy files", errs)
	}

	for _, folder := range di.folders {
		copyFolder(filepath.Join(sourceFolder, folder), filepath.Join(destinationFolder, folder))
	}
}

// check if folder exist if not create it
func checkOrCreateFolder(path string) error {
	p := filepath.Join(path)
	if _, err := os.Stat(p); os.IsNotExist(err) {

		err := os.Mkdir(p, os.ModePerm)
		return err
	}
	return nil
}
