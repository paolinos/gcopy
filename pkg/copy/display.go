package copy

import (
	"fmt"
	"time"

	"github.com/paolinos/gcopy/pkg/analyzer"
)

type fileProgress struct {
	Filename string
	Size     int64
}

func (f *fileProgress) SizeReadable() string {
	return analyzer.GetSizeReadable(f.Size)
}

func FormatSummaryLine(totalFiles int, sizeReadable string, elapsed time.Duration, errCount int) string {
	if errCount > 0 {
		return fmt.Sprintf("\nCopied %d files (%s) in %.2fs — %d file(s) failed",
			totalFiles, sizeReadable, elapsed.Seconds(), errCount)
	}
	return fmt.Sprintf("\nCopied %d files (%s) in %.2fs",
		totalFiles, sizeReadable, elapsed.Seconds())
}
