package copy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/paolinos/gcopy/pkg/analyzer"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

type fileProcess struct {
	Source            string
	FolderDestination string
	Destination       string
	Progress          *fileProgress
	Bar               *mpb.Bar
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
	resultMu         sync.Mutex
	progress         []*fileProgress
	chFileProcessing chan fileProcess
	wg               sync.WaitGroup
	check            bool
	skip             bool
	force            bool
	sourcePath       string
	destinationPath  string
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

func copyFileChunks(srcFilePath string, dstFilePath string, chunkSize int64, progress *fileProgress, bar ...*mpb.Bar) (int64, error) {
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

		if len(bar) > 0 && bar[0] != nil {
			bar[0].IncrBy(n)
		}
	}

	return copied, nil
}

func (c *CopyManager) Run(data analyzer.AnalyzeResult, chunkSize int64) {
	c.totalFiles = data.TotalFiles
	workerLimit := 5

	// Create mpb progress container.
	p := mpb.New(
		mpb.WithWidth(getTerminalWidth()),
		mpb.WithRefreshRate(100*time.Millisecond),
	)

	for i := 0; i < workerLimit; i++ {
		go func(id int) {
			for file := range c.chFileProcessing {
				result := CopyResult{
					CheckPassed:     true,
					SourcePath:      file.Source,
					DestinationPath: file.Destination,
				}

				if c.skip {
					if _, dstErr := os.Stat(file.Destination); dstErr == nil {
						result.Err = fmt.Errorf("skipped: %s already exists", filepath.Base(file.Destination))
						result.CheckPassed = false
						file.Bar.SetCurrent(file.Progress.Size)
						c.addResult(result)
						c.wg.Done()
						continue
					}
				}

				_, err := copyFileChunks(file.Source, file.Destination, chunkSize, file.Progress, file.Bar)
				if err != nil {
					result.Err = err
					result.CheckPassed = false
				}

				c.addResult(result)
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
			fp := &fileProgress{
				Filename: file.Path,
				Size:     file.Size,
			}

			name := filepath.Base(file.Path)
			bar := p.AddBar(file.Size,
				mpb.BarWidth(getTerminalWidth()-50),
				mpb.PrependDecorators(
					decor.Name(name, decor.WC{W: 20}),
				),
				mpb.AppendDecorators(
					decor.Counters(decor.UnitKiB, " % .1f / % .1f"),
					decor.Percentage(decor.WC{W: 5}),
				),
			)

			c.progress = append(c.progress, fp)
			c.wg.Add(1)
			c.chFileProcessing <- fileProcess{
				Source:      file.Path,
				Destination: filepath.Join(dstDir, filepath.Base(file.Path)),
				Progress:    fp,
				Bar:         bar,
			}
		}
	}

	if len(data.Folders) == 0 && data.TotalFiles == 1 {
		fp := &fileProgress{
			Filename: data.Source,
			Size:     data.TotalSize,
		}

		name := filepath.Base(data.Source)
		bar := p.AddBar(data.TotalSize,
			mpb.BarWidth(getTerminalWidth()-50),
			mpb.PrependDecorators(
				decor.Name(name, decor.WC{W: 20}),
			),
			mpb.AppendDecorators(
				decor.Counters(decor.UnitKiB, " % .1f / % .1f"),
				decor.Percentage(decor.WC{W: 5}),
			),
		)

		c.progress = append(c.progress, fp)
		c.wg.Add(1)
		c.chFileProcessing <- fileProcess{
			Source:      data.Source,
			Destination: data.Destination,
			Progress:    fp,
			Bar:         bar,
		}
	}

	c.wg.Wait()
	p.Wait()
}

func (c *CopyManager) addResult(result CopyResult) {
	c.resultMu.Lock()
	defer c.resultMu.Unlock()
	c.result = append(c.result, result)
}

func (c *CopyManager) Result() []CopyResult {
	c.resultMu.Lock()
	defer c.resultMu.Unlock()
	out := make([]CopyResult, len(c.result))
	copy(out, c.result)
	return out
}

func (c *CopyManager) Progress() []*fileProgress {
	return c.progress
}

func (c *CopyManager) SetHeaderInfo(sourcePath, destinationPath string, totalFiles int, sizeReadable string) {
	c.sourcePath = sourcePath
	c.destinationPath = destinationPath
	c.totalFiles = totalFiles
}
