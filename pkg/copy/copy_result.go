package copy

type copyResult struct {
	Source        string
	Destination   string
	copiedFiles   int
	copiedFolders int
}

type CopyStatus struct {
	TotalFiles    int
	CopiedFiles   int
	TotalFolders  int
	CopiedFolders int
	TotalErrors   int
}

type CopyResult interface {
	GetSource() string
	GetDestination() string
	CopiedFiles() int
	CopiedFolders() int
}

func (c copyResult) GetSource() string {
	return c.Source
}

func (c copyResult) GetDestination() string {
	return c.Destination
}

func (c copyResult) CopiedFiles() int {
	return c.copiedFiles
}

func (c copyResult) CopiedFolders() int {
	return c.copiedFolders
}
