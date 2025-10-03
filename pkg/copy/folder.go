package copy

import (
	"os"
	"path/filepath"
)

// check if folder exist if not create it
func checkOrCreateFolder(path string) error {
	p := filepath.Join(path)
	if _, err := osStat(p); osIsNotExist(err) {

		err := osMkdir(p, os.ModePerm)
		return err
	}
	return nil
}
