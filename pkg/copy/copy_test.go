package copy

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

const validFilePath = "./file"
const invalidFilePath = "unknown_file"

var ErrInvalidPath = errors.New("Invalid path")

type mockFile struct {
	os.File
}

type otherMock struct {
	io.Reader
}

func (m *mockFile) Close() error {
	return nil
}
func (m *otherMock) Close() error { return nil }

// Override the osOpen function v could be `error` or `os.File`. you can use os.CreateTemp("", "data") to create a temp file
func mockOpen(v interface{}) {
	osOpen = func(name string) (*os.File, error) {
		if err, ok := v.(error); ok {
			return &os.File{}, err
		}

		f := v.(*os.File)
		if f.Name() != name {
			return f, ErrInvalidPath
		}

		f, _ = os.Open(f.Name())
		return f, nil
	}
}

// Override the osCreate function
func mockCreate(v interface{}) {
	osCreate = func(name string) (*os.File, error) {
		if err, ok := v.(error); ok {
			return &os.File{}, err
		}

		f := v.(*os.File)
		if f.Name() != name {
			return f, ErrInvalidPath
		}
		return f, nil
	}
}

func TestCopy(t *testing.T) {

	sourceFile, _ := os.CreateTemp("", "source file")
	destinationFile, _ := os.CreateTemp("", "")

	defer func() {

		sourceFile.Close()
		destinationFile.Close()
		fmt.Print("Defer closing TestCopy")
	}()

	t.Run("Given a invalid source When trying to copy Then should return an error.", func(t *testing.T) {
		mockOpen(ErrSourcePath)

		res, err := CopyFromTo(CopyOptions{
			Source:      invalidFilePath,
			Destination: destinationFile.Name(),
			Override:    false,
		})

		expectError(t, err, ErrSourcePath)
		expectCopyResult(t, res, invalidFilePath, destinationFile.Name(), 0, 0)
	})

	t.Run("Given a invalid destination When trying to copy Then should return an error.", func(t *testing.T) {
		mockOpen(sourceFile)
		mockCreate(ErrDestinationPath)

		res, err := CopyFromTo(CopyOptions{
			Source:      sourceFile.Name(),
			Destination: invalidFilePath,
			Override:    false,
		})

		expectError(t, err, ErrDestinationPath)
		expectCopyResult(t, res, sourceFile.Name(), invalidFilePath, 0, 0)
	})

	t.Run("Given a valid source and destination When copy the file Then should return the result.", func(t *testing.T) {
		mockOpen(sourceFile)
		mockCreate(destinationFile)

		res, err := CopyFromTo(CopyOptions{
			Source:      sourceFile.Name(),
			Destination: destinationFile.Name(),
			Override:    false,
		})

		if err != nil {
			t.Errorf("Found an unexpected error %s", err)
		}
		expectCopyResult(t, res, sourceFile.Name(), destinationFile.Name(), 1, 0)
	})
}

func expectError(t *testing.T, err error, expectedError error) {
	if err == nil {
		t.Error("Expected to have an error, but received nil")
	}

	if expectedError != err {
		t.Error(fmt.Sprintf("Expected to have an '%s' but received:", expectedError), err)
	}
}

func expectCopyResult(t *testing.T, res CopyResult, source string, destination string, files int, folders int) {
	if res.GetSource() != source {
		t.Error("Invalid source.", "expected", source, "received", res.GetSource())
	}
	if res.GetDestination() != destination {
		t.Error("Invalid destination.", "expected", destination, "received", res.GetDestination())
	}

	if res.CopiedFiles() != files {
		t.Error("Invalid copied files.", "expected", files, "received", res.CopiedFiles())
	}

	if res.CopiedFolders() != folders {
		t.Error("Invalid copied folders.", "expected", folders, "received", res.CopiedFolders())
	}
}
