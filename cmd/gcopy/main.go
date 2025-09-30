package main

import (
	"fmt"
	"os"
	"time"

	"github.com/paolinos/gcopy/pkg/analyzer"
	"github.com/paolinos/gcopy/pkg/copy"
)

var (
	Version     string = "0.0.1"
	Description string
	BuildTime   string
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Bold   = "\033[1m"
)

const PROMPT_HELP = `
GCopy (v %s) tool is used to copy files/folders
%s

Usage:

	gcopy [source] [destination]
`

func main() {

	// TODO: move this validation of the cmd
	args := os.Args[1:]
	if len(args) != 2 {

		fmt.Printf(PROMPT_HELP, Version, Description)
		os.Exit(0)
	}

	start := time.Now()

	res, err := analyzer.AnalyzePath(args[0], args[1])
	if err != nil {
		fmt.Printf("%s%sError:%s\n %s \n -----\n%s", Red, Bold, Bold, err, Reset)
		os.Exit(0)
	}

	deltaT := time.Since(start)

	fmt.Printf(`%sPrepare to copy:
	- from:%s; 
	- to:%s
Total Files: %d
Total Size: %s
Time: %f
%s`, Green, res.Source, res.Destination, res.TotalFiles, res.SizeReadable, deltaT.Seconds(), Reset)

	for _, v := range res.Folders {
		fmt.Printf("Path %s, Size: %d, SizeReadable: %s, Files to copy: %d,\n", v.Path, v.Size, v.SizeReadable, len(v.Files))
		for _, f := range v.Files {
			fmt.Printf("		- Filepath:%s; Size:%d;\n", f.Path, f.Size)
		}
	}

	start = time.Now()
	copy.CopyPath(res)
	deltaT = time.Since(start)

	// TODO: improve result message
	fmt.Printf("Total files to copied: %d, size:%s, Time:%f seconds, error:%s", res.TotalFiles, res.SizeReadable, deltaT.Seconds(), err)
}
