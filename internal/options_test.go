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

		expectErrAndEmptyOptions(t, o, err)
	})

	t.Run("Given options after the source and destination, When parsing the options, Then should return error", func(t *testing.T) {
		os.Args = []string{"cmd", source, destination, fmt.Sprintf("-chunks=%d", chunk)}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		o, err := GetAppOptions()

		expectErrAndEmptyOptions(t, o, err)
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
		if o.Chunks != 0 {
			t.Errorf("Expected chunk to be:0, received: %d", o.Chunks)
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

	t.Run("Given skip and force flags, When parsing the options, Then should return correct values", func(t *testing.T) {
		os.Args = []string{"cmd", source, destination, "--skip", "--force"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		o, err := GetAppOptions()
		if err != nil {
			t.Errorf("Unexpected error: %e", err)
		}
		if !o.Skip {
			t.Errorf("Expected Skip to be true")
		}
		if !o.Force {
			t.Errorf("Expected Force to be true")
		}
	})

	t.Run("Given show and parallel flags, When parsing the options, Then should return correct values", func(t *testing.T) {
		os.Args = []string{"cmd", source, destination, "--show=10", "--parallel=8"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		o, err := GetAppOptions()
		if err != nil {
			t.Errorf("Unexpected error: %e", err)
		}
		if o.Show != 10 {
			t.Errorf("Expected Show to be 10, got %d", o.Show)
		}
		if o.Parallel != 8 {
			t.Errorf("Expected Parallel to be 8, got %d", o.Parallel)
		}
	})

}

func expectErrAndEmptyOptions(t *testing.T, o *AppOptions, err error) {
	if err == nil {
		t.Error("Expected error but got nil")
	}
	if o != nil {
		if len(o.Source) != 0 {
			t.Errorf("Source should be empty, got: %s", o.Source)
		}
		if len(o.Destination) != 0 {
			t.Errorf("Destination should be empty, got: %s", o.Destination)
		}
		if o.Chunks != 0 {
			t.Errorf("Chunks should be 0, got: %d", o.Chunks)
		}
		if o.Skip {
			t.Errorf("Skip should be false")
		}
		if o.Force {
			t.Errorf("Force should be false")
		}
		if o.Show != 5 {
			t.Errorf("Show should be 5, got: %d", o.Show)
		}
		if o.Parallel != 5 {
			t.Errorf("Parallel should be 5, got: %d", o.Parallel)
		}
	}
}
