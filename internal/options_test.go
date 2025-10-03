package internal

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGetHelper(t *testing.T) {

	version := "version-test"
	description := "description test"
	helper := GetHelper(version, description)

	if !strings.Contains(helper, version) {
		t.Errorf("No version found in the helper. Helper: %s", helper)
	}
	if !strings.Contains(helper, description) {
		t.Errorf("No description found in the helper. Helper: %s", helper)
	}

}

func TestGetAppOptions(t *testing.T) {

	source := "./source"
	destination := "./source"
	chunk := int64(100)

	t.Run("Given no arguments, When parsing the options, Then should return error", func(t *testing.T) {
		os.Args = []string{"cmd"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		o, err := GetAppOptions()

		expectErrAndEmptyOptions(t, ErrMissingArguments, o, err)
	})

	t.Run("Given options after the source and destination, When parsing the options, Then should return error", func(t *testing.T) {
		os.Args = []string{"cmd", source, destination, fmt.Sprintf("-chunks=%d", chunk)}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		o, err := GetAppOptions()

		expectErrAndEmptyOptions(t, ErrMissingArguments, o, err)
	})

	t.Run("Given a source and destination, When parsing the options, Then should return options and default chunk", func(t *testing.T) {
		os.Args = []string{"cmd", source, destination}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		o, err := GetAppOptions()
		if err != nil {
			t.Errorf("Unexpected error: %e", err)
		}
		if o.Source != source {
			t.Errorf("Expected source to be:%s, received: %s", source, o.Source)
		}
		if o.Destination != destination {
			t.Errorf("Expected source to be:%s, received: %s", destination, o.Destination)
		}
		if o.Chunks != default_chunks_bytes {
			t.Errorf("Expected chunk to be:%d, received: %d", default_chunks_bytes, o.Chunks)
		}
	})

	t.Run("Given a chunk, When parsing the options, Then should return passed chunk", func(t *testing.T) {
		os.Args = []string{"cmd", fmt.Sprintf("-chunks=%d", chunk), source, destination}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		o, err := GetAppOptions()
		if err != nil {
			t.Errorf("Unexpected error: %e", err)
		}
		if o.Chunks != chunk {
			t.Errorf("Expected chunk to be:%d, received: %d", chunk, o.Chunks)
		}
	})
}

func expectErrAndEmptyOptions(t *testing.T, expectedError error, o AppOptions, err error) {
	if err != expectedError {
		t.Error("Found options after the source and destination ")
	}
	if len(o.Source) != 0 {
		t.Errorf("Source should be empty")
	}
	if len(o.Destination) != 0 {
		t.Errorf("Destination should be empty")
	}
	if o.Chunks != 0 {
		t.Errorf("Chunks should be empty")
	}
}
