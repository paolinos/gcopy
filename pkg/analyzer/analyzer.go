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

func GetSizeReadable(size int64) string {
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

	sourcePath, err := os.Stat(res.Source)
	if err != nil {
		return res, err
	}

	// File
	if !sourcePath.IsDir() {
		res.TotalFiles = 1
		res.TotalSize = sourcePath.Size()
		sizeReadable := GetSizeReadable(res.TotalSize)
		res.SizeReadable = sizeReadable
		return res, nil
	}

	// Folders
	folders := make(map[string]FolderInfo)
	getFolderInfo := func(path string) *FolderInfo {
		f, ok := folders[path]
		if !ok {
			f = FolderInfo{
				Path:  path,
				Files: []FileInfo{},
			}
			folders[path] = f
		}
		return &f
	}

	err = filepathWalk(res.Source, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			folder := filepath.Dir(p)
			currentFolder := getFolderInfo(folder)

			bytes := info.Size()
			currentFolder.Size += bytes
			currentFolder.Files = append(currentFolder.Files, FileInfo{
				Path: p,
				Size: bytes,
			})

			folders[currentFolder.Path] = *currentFolder
		}
		return nil
	})

	for _, f := range folders {
		res.Folders = append(res.Folders, f)
		res.TotalSize += f.Size
		res.TotalFiles += len(f.Files)
	}
	res.SizeReadable = GetSizeReadable(res.TotalSize)

	return res, err
}
