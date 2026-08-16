package copy

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"time"

	"github.com/paolinos/gcopy/pkg/analyzer"
)

// fileProgress tracks per-file transfer statistics.
type fileProgress struct {
	Filename string
	Size     int64
	Copied   int64
	Start    time.Time
	End      time.Time
}

func (f *fileProgress) AverageCopyMiB() float64 {
	if f.End.IsZero() || f.Start.Equal(f.End) {
		return 0
	}
	dur := f.End.Sub(f.Start).Seconds()
	if dur == 0 {
		return 0
	}
	bytesPerSec := float64(f.Copied) / dur
	return bytesPerSec / (1024 * 1024)
}

func (f *fileProgress) ElapsedTime() string {
	if f.End.IsZero() {
		return "00:00"
	}
	dur := f.End.Sub(f.Start).Round(time.Second)
	if dur < 0 {
		return "00:00"
	}
	mins := int(dur.Minutes())
	secs := int(dur.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}

func (f *fileProgress) SizeReadable() string {
	return analyzer.GetSizeReadable(f.Size)
}

// FormatProgressBar renders a progress bar of the given width without brackets.
func FormatProgressBar(pct float64, width int) string {
	if width <= 0 {
		return ""
	}
	filled := int(math.Round(float64(width) * pct / 100))
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	empty := width - filled
	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}

// FormatFileLine renders a single file line with dynamic column sizing.
func FormatFileLine(fp *fileProgress) string {
	pct := float64(0)
	if fp.Size > 0 {
		pct = float64(fp.Copied) / float64(fp.Size) * 100
	}

	termWidth := getTerminalWidth()
	// Layout: name + space + size + 2 spaces + speed + space + bar + space + pct + space + time
	overhead := 1 + sizeColumnWidth + 2 + speedColumnWidth + 1 + pctColumnWidth + 1 + timeColumnWidth
	barWidth := termWidth - overhead - filenameMinWidth
	if barWidth < 10 {
		barWidth = 10
	}
	nameWidth := termWidth - overhead - barWidth
	if nameWidth < filenameMinWidth {
		nameWidth = filenameMinWidth
	}

	name := fp.Filename
	if len(name) > nameWidth-1 {
		name = name[:nameWidth-1] + "…"
	} else {
		name = strings.Repeat(" ", nameWidth-len(name)) + name
	}

	bar := FormatProgressBar(pct, barWidth)
	return fmt.Sprintf("%s %s  %7.2f MiB/s %s %3d%% %s",
		name,
		fp.SizeReadable(),
		fp.AverageCopyMiB(),
		bar,
		int(math.Round(pct)),
		fp.ElapsedTime(),
	)
}

// DisplayResults prints the final copy summary.
func DisplayResults(data analyzer.AnalyzeResult, results []CopyResult, progress []*fileProgress, showLimit int) {
	fmt.Println("Copy:")
	fmt.Printf("  from: %s\n", data.Source)
	fmt.Printf("  to:   %s\n", data.Destination)
	fmt.Printf("\nTotal: %d files to copy\n", data.TotalFiles)
	fmt.Println(strings.Repeat("-", 40))

	sorted := make([]*fileProgress, len(progress))
	copy(sorted, progress)
	for i := range sorted {
		if sorted[i] != nil {
			sorted[i].Filename = filepath.Base(sorted[i].Filename)
		}
	}

	n := len(sorted)
	if n > showLimit {
		sorted = sorted[n-showLimit:]
	}

	for _, fp := range sorted {
		fmt.Println(FormatFileLine(fp))
	}
}

// FormatSummaryLine formats the final stats line.
func FormatSummaryLine(totalFiles int, sizeReadable string, elapsed time.Duration, errCount int) string {
	if errCount > 0 {
		return fmt.Sprintf("\nCopied %d files (%s) in %.2fs — %d file(s) failed",
			totalFiles, sizeReadable, elapsed.Seconds(), errCount)
	}
	return fmt.Sprintf("\nCopied %d files (%s) in %.2fs",
		totalFiles, sizeReadable, elapsed.Seconds())
}
