package main

import (
	"fmt"
	"os"

	"github.com/paolinos/gcopy/pkg/copy"
)

var (
	Version     string = "0.0.0"
	Description string
	BuildTime   string
)

func main() {

	// TODO: move this validation of the cmd
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Printf(`
GCopy (v %s) tool is used to copy files/folders
%s

Usage:

	gcopy [source] [destination]
`, Version, Description)
		os.Exit(0)
	}

	res, err := copy.CopyFromTo(copy.CopyOptions{
		Source:      args[0],
		Destination: args[1],
	})

	if err != nil {
		fmt.Println(err)
	}

	// TODO: improve result message
	fmt.Println("Result", res)
}
