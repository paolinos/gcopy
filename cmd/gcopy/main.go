package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/paolinos/gcopy/internal"
	"github.com/paolinos/gcopy/pkg/analyzer"
	"github.com/paolinos/gcopy/pkg/copy"
)

var (
	Version     = "0.0.1"
	Description = "Cross-platform file copy utility"
	BuildTime   string

	reset = "\033[0m"
	red   = "\033[31m"
	green = "\033[32m"
	bold  = "\033[1m"
)

func getFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func main() {
	cmd, options := internal.NewRootCmd(Version, Description)

	if err := cmd.Execute(); err != nil {
		os.Exit(2)
	}
	if options.Source == "" {
		return
	}

	start := time.Now()

	res, err := analyzer.AnalyzePath(options.Source, options.Destination)
	if err != nil {
		fmt.Printf("%s%sError:%s %s\n -----\n", red, bold, reset, err)
		cmd.Usage()
		os.Exit(2)
	}

	fmt.Printf("%sPreparing to copy:%s\n  from: %s\n  to:   %s\n  files: %d  size: %s\n\n",
		green, reset, res.Source, res.Destination, res.TotalFiles, res.SizeReadable)

	start = time.Now()
	cm := copy.NewCopyManager(options.Check, options.Skip, options.Force)
	cm.SetHeaderInfo(res.Source, res.Destination, res.TotalFiles, res.SizeReadable)
	cm.Run(res, options.Chunks)
	elapsed := time.Since(start)

	// Count failures.
	errCount := 0
	for _, r := range cm.Result() {
		if r.Err != nil {
			errCount++
		}
	}

	// Display results.
	// copy.DisplayResults(res, cm.Result(), cm.Progress(), options.Show)
	fmt.Println(copy.FormatSummaryLine(res.TotalFiles, res.SizeReadable, elapsed, errCount))

	// If --check is enabled, validate each file's source and destination sizes.
	if options.Check {
		var failures []string
		for _, result := range cm.Result() {
			base := filepath.Base(result.SourcePath)
			if !result.CheckPassed || result.Err != nil {
				failures = append(failures, fmt.Sprintf("%s: %v", base, result.Err))
				fmt.Printf("%s%sError:%s %s — %v\n", red, bold, reset, base, result.Err)
				continue
			}
			srcSize, err := getFileSize(result.SourcePath)
			if err != nil {
				failures = append(failures, fmt.Sprintf("%s: source stat failed: %v", base, err))
				fmt.Printf("%s%sError:%s could not stat source %s: %v\n", red, bold, reset, base, err)
				continue
			}
			dstSize, err := getFileSize(result.DestinationPath)
			if err != nil {
				failures = append(failures, fmt.Sprintf("%s: dest stat failed: %v", base, err))
				fmt.Printf("%s%sError:%s could not stat destination %s: %v\n", red, bold, reset, base, err)
				continue
			}
			if srcSize != dstSize {
				failures = append(failures, fmt.Sprintf("%s: size mismatch src=%d dst=%d", base, srcSize, dstSize))
				fmt.Printf("%s%sError:%s size mismatch %s (src=%d, dst=%d)\n", red, bold, reset, base, srcSize, dstSize)
			} else {
				fmt.Printf("  ✓ %s (%d bytes)\n", base, srcSize)
			}
		}

		if len(failures) > 0 {
			fmt.Printf("\n%s%s%d file(s) failed validation%s\n", red, bold, len(failures), reset)
			os.Exit(1)
		}
	}
}
