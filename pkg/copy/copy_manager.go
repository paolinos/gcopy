package copy

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"

	"github.com/paolinos/gcopy/pkg/analyzer"
)

type fileProcess struct {
	Source            string
	FolderDestination string
	Destination       string
}

type CopyResult struct {
	Err  error
	Info string
}

type CopyManager struct {
	result           []CopyResult
	chFileProcessing chan fileProcess
	wg               sync.WaitGroup
}

func NewCopyManager() *CopyManager {
	return &CopyManager{
		chFileProcessing: make(chan fileProcess),
		result:           []CopyResult{},
	}
}

// Copy chunks from source to destination
func (c *CopyManager) copyChunksFromSource(srcFilePath string, dstFilePath string, chunkSize int64) error {

	source, err := osOpen(srcFilePath)
	if err != nil {
		return err
	}
	defer source.Close()

	folderPath := filepath.Dir(dstFilePath)
	err = checkOrCreateFolder(folderPath)
	if err != nil {
		return err
	}

	destination, err := osCreate(dstFilePath)
	if err != nil {
		return err
	}
	defer destination.Close()

	buf := make([]byte, chunkSize)

	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			// TODO: Error reading. we should retry
			return err
		}
		if n == 0 {
			return nil
		}

		_, err = destination.Write(buf[:n])
		if err != nil {
			// TODO: Error writing. We should retry
			return err
		}
	}
}

func (c *CopyManager) Run(data analyzer.AnalyzeResult, chunkSize int64) {

	fmt.Printf("Copy total files:%d", data.TotalFiles)

	workerLimit := 10

	for i := 0; i < workerLimit; i++ {
		go func(id int) {
			for file := range c.chFileProcessing {
				err := c.copyChunksFromSource(file.Source, file.Destination, chunkSize)
				result := CopyResult{
					Err:  err,
					Info: "",
				}
				if err != nil {
					result.Info = fmt.Sprintf("File copy from:%s, to:%s", file.Source, file.Destination)
				}
				c.result = append(c.result, result)
				c.wg.Done()
			}
		}(i)
	}

	//c.files = append(c.files, fileProcess{Source: source, Destination: destination})
	for _, folder := range data.Folders {
		d := strings.Replace(folder.Path, data.Source, data.Destination, 1)
		/*err := checkOrCreateFolder(data.Destination)
		if err != nil {
			// TODO: need to move to other level
			fmt.Printf("Error to create folder")
			continue
		}*/

		for _, file := range folder.Files {
			filename := filepath.Base(file.Path)
			//filename := strings.Replace(file.Path, data.Source, "", 1)
			dstPath := filepath.Join(d, filename)
			//c.files = append(c.files, fileProcess{Source: file.Path, Destination: dstPath})
			c.wg.Add(1)
			c.chFileProcessing <- fileProcess{Source: file.Path, FolderDestination: d, Destination: dstPath}
		}
	}

	c.wg.Wait()
	close(c.chFileProcessing)

	fmt.Printf("Copy ends %v", c.result)
}
