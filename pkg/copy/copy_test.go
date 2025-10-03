package copy

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

const invalidFilePath = "unknown_file"
const DEFAULT_SIZE = 100

var ErrUnexpected = errors.New("unexpected error")
var ErrInvalidPath = errors.New("invalid path")

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
	}()

	t.Run("Given a invalid source When trying to copy chunks Then should return an error.", func(t *testing.T) {
		var ErrSourcePath = errors.New("the source path is invalid")
		mockOpen(ErrSourcePath)

		err := copyChunksFromSource(invalidFilePath, destinationFile.Name(), DEFAULT_SIZE)

		expectError(t, err, ErrSourcePath)
	})

	t.Run("Given a invalid destination When trying to copy chunks Then should return an error.", func(t *testing.T) {
		var ErrDestinationPath = errors.New("the destination path is invalid")
		mockOpen(sourceFile)
		mockCreate(ErrDestinationPath)

		err := copyChunksFromSource(sourceFile.Name(), invalidFilePath, DEFAULT_SIZE)

		expectError(t, err, ErrDestinationPath)
	})

	t.Run("Given a valid source and destination When copy the file Then should return the result.", func(t *testing.T) {
		mockOpen(sourceFile)
		mockCreate(destinationFile)

		err := copyChunksFromSource(sourceFile.Name(), destinationFile.Name(), DEFAULT_SIZE)
		if err != nil {
			t.Errorf("Found an unexpected error %s", err)
		}
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
