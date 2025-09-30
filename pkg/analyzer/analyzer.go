package analyzer

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const (
	Copy     = "copy"
	Override = "override"
	None     = "none"
)

type FileInfo struct {
	Size int64
	Path string
}

type FolderInfo struct {
	Path         string
	Files        []FileInfo
	Size         int64
	SizeReadable string
}

type AnalyzeResult struct {
	TotalFiles   int
	TotalSize    int64
	SizeReadable string
	Folders      []FolderInfo
	Source       string
	Destination  string
}

var filepathWalk = filepath.Walk

func getSizeReadable(size int64) string {
	n := float64(size)
	var power = float64(1 << 10)
	for _, unit := range []string{"B", "KB", "MB", "GB", "TB", "PB"} {
		if n >= power {
			n = float64(n) / float64(power)
		} else {
			format := "%.1f%s"
			_, fracPart := math.Modf(n)
			d := fmt.Sprintf("%.1f", fracPart)
			if strings.HasPrefix(d, "0.0") {
				format = "%.0f%s"
			}
			return fmt.Sprintf(format, n, unit)
		}

	}
	return fmt.Sprintf("%d", size)
}

// Analyze path and return a list of files from the path
func AnalyzePath(source string, destination string) (AnalyzeResult, error) {

	res := AnalyzeResult{
		TotalFiles:  0,
		TotalSize:   0,
		Folders:     []FolderInfo{},
		Source:      filepath.Join(source),
		Destination: filepath.Join(destination),
	}

	var currentFolder = FolderInfo{}
	err := filepathWalk(res.Source, func(p string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			if currentFolder.Path != "" {
				currentFolder.SizeReadable = getSizeReadable(currentFolder.Size)
				res.Folders = append(res.Folders, currentFolder)

				res.TotalSize += currentFolder.Size
				res.TotalFiles += len(currentFolder.Files)
			}
			currentFolder = FolderInfo{
				Path:  p,
				Files: []FileInfo{},
			}
		} else {
			bytes := info.Size()
			currentFolder.Size += bytes
			currentFolder.Files = append(currentFolder.Files, FileInfo{
				Path: p,
				Size: bytes,
			})
		}
		return nil
	})

	res.SizeReadable = getSizeReadable(res.TotalSize)

	return res, err
}
