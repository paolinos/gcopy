package copy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/paolinos/gcopy/pkg/analyzer"
)

type fileProcess struct {
	Source            string
	FolderDestination string
	Destination       string
	Progress          *fileProgress
}

type CopyResult struct {
	Err             error
	Info            string
	CheckPassed     bool
	SourcePath      string
	DestinationPath string
}

type CopyManager struct {
	result           []CopyResult
	progress         []*fileProgress
	chFileProcessing chan fileProcess
	wg               sync.WaitGroup
	check            bool
	skip             bool
	force            bool
	activeCount      atomic.Int32
	doneCount        atomic.Int32
	totalFiles       int
}

func NewCopyManager(check, skip, force bool) *CopyManager {
	return &CopyManager{
		chFileProcessing: make(chan fileProcess),
		result:           []CopyResult{},
		progress:         []*fileProgress{},
		check:            check,
		skip:             skip,
		force:            force,
	}
}

func (c *CopyManager) copyChunksFromSource(srcFilePath string, dstFilePath string, chunkSize int64, progress *fileProgress) (int64, error) {
	source, err := os.Open(srcFilePath)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	folderPath := filepath.Dir(dstFilePath)
	err = checkOrCreateFolder(folderPath)
	if err != nil {
		return 0, err
	}

	destination, err := os.Create(dstFilePath)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	buf := make([]byte, chunkSize)
	copied := int64(0)

	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return copied, err
		}
		if n == 0 {
			break
		}

		_, err = destination.Write(buf[:n])
		if err != nil {
			return copied, err
		}
		copied += int64(n)

		progress.Copied = copied
		progress.End = time.Now()
	}

	return copied, nil
}

// displayProgress continuously redraws live progress for all active files.
func (c *CopyManager) displayProgress(stop <-chan struct{}) {
	if !IsTerminal() {
		<-stop
		return
	}

	fmt.Print("\033[?25l") // hide cursor
	defer fmt.Print("\033[?25h") // show cursor

	prevLineCount := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			// Final redraw
			c.redrawProgress(&prevLineCount)
			return
		case <-ticker.C:
			c.redrawProgress(&prevLineCount)
		}
	}
}

func (c *CopyManager) redrawProgress(prevLineCount *int) {
	active := make([]*fileProgress, 0)
	for _, fp := range c.progress {
		if fp.End.IsZero() || fp.Copied < fp.Size {
			active = append(active, fp)
		}
	}

	// Move cursor up to the start of the progress area and clear everything.
	if *prevLineCount > 0 {
		fmt.Printf("\033[%dA", *prevLineCount)
		for i := 0; i < *prevLineCount; i++ {
			fmt.Print("\033[2K\033[1B")
		}
		fmt.Printf("\033[%dA", *prevLineCount)
	}

	for _, fp := range active {
		fmt.Printf("\033[2K%s\033[1B", FormatFileLine(fp))
	}
	// Move cursor back to start of progress area.
	if len(active) > 0 {
		fmt.Printf("\033[%dA", len(active))
	}

	*prevLineCount = len(active)
}

func (c *CopyManager) Run(data analyzer.AnalyzeResult, chunkSize int64) {
	c.totalFiles = data.TotalFiles
	workerLimit := 10

	// Start live progress display.
	stopDisplay := make(chan struct{})
	go c.displayProgress(stopDisplay)

	for i := 0; i < workerLimit; i++ {
		go func(id int) {
			for file := range c.chFileProcessing {
				c.activeCount.Add(1)
				result := CopyResult{
					CheckPassed:     true,
					SourcePath:      file.Source,
					DestinationPath: file.Destination,
				}

				// Only skip if --skip is explicitly set and dest exists.
				if c.skip {
					if _, dstErr := os.Stat(file.Destination); dstErr == nil {
						result.Err = fmt.Errorf("skipped: %s already exists", filepath.Base(file.Destination))
						result.CheckPassed = false
						if file.Progress != nil {
							file.Progress.End = file.Progress.Start
						}
						c.result = append(c.result, result)
						c.activeCount.Add(-1)
						c.doneCount.Add(1)
						c.wg.Done()
						continue
					}
				}

				_, err := c.copyChunksFromSource(file.Source, file.Destination, chunkSize, file.Progress)
				if err != nil {
					result.Err = err
					result.CheckPassed = false
				}

				c.result = append(c.result, result)
				c.activeCount.Add(-1)
				c.doneCount.Add(1)
				c.wg.Done()
			}
		}(i)
	}

	for _, folder := range data.Folders {
		d, _ := filepath.Rel(data.Source, folder.Path)
		dstDir := filepath.Join(data.Destination, d)
		err := checkOrCreateFolder(dstDir)
		if err != nil {
			fmt.Printf("Error creating folder: %v\n", err)
			continue
		}

		for _, file := range folder.Files {
			filename := filepath.Base(file.Path)
			dstPath := filepath.Join(dstDir, filename)
			fp := &fileProgress{
				Filename: file.Path,
				Size:     file.Size,
				Start:    time.Now(),
			}
			c.progress = append(c.progress, fp)

			c.wg.Add(1)
			c.chFileProcessing <- fileProcess{Source: file.Path, FolderDestination: d, Destination: dstPath, Progress: fp}
		}
	}

	if len(data.Folders) == 0 && data.TotalFiles == 1 {
		fp := &fileProgress{
			Filename: data.Source,
			Size:     data.TotalSize,
			Start:    time.Now(),
		}
		c.progress = append(c.progress, fp)

		c.wg.Add(1)
		c.chFileProcessing <- fileProcess{Source: data.Source, Destination: data.Destination, Progress: fp}
	}

	c.wg.Wait()
	close(stopDisplay)

	// Final clear of progress area.
	if IsTerminal() && len(c.progress) > 0 {
		fmt.Printf("\033[%dA", len(c.progress))
		for i := 0; i < len(c.progress); i++ {
			fmt.Print("\033[2K\033[1B")
		}
		fmt.Printf("\033[%dA", len(c.progress))
	}
}

func (c *CopyManager) Result() []CopyResult {
	return c.result
}

func (c *CopyManager) Progress() []*fileProgress {
	return c.progress
}
