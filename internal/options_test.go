package internal

import (
	"fmt"
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

func TestRootCmd(t *testing.T) {
	source := "./source"
	destination := "./dest"

	t.Run("Given no arguments, When executing, Then should return error", func(t *testing.T) {
		cmd, _ := NewRootCmd("test", "test")
		cmd.SetArgs([]string{})

		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error but got nil")
		}
	})

	t.Run("Given only source, When executing, Then should return error", func(t *testing.T) {
		cmd, _ := NewRootCmd("test", "test")
		cmd.SetArgs([]string{source})

		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error but got nil")
		}
	})

	t.Run("Given too many arguments, When executing, Then should return error", func(t *testing.T) {
		cmd, _ := NewRootCmd("test", "test")
		cmd.SetArgs([]string{source, destination, "extra"})

		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error but got nil")
		}
	})

	t.Run("Given a source and destination, When executing, Then should parse correctly", func(t *testing.T) {
		cmd, opts := NewRootCmd("test", "test")
		cmd.SetArgs([]string{source, destination})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if opts.Source != source {
			t.Errorf("Expected source to be %s, got: %s", source, opts.Source)
		}
		if opts.Destination != destination {
			t.Errorf("Expected destination to be %s, got: %s", destination, opts.Destination)
		}
		if opts.Chunks != 32*1024 {
			t.Errorf("Expected default chunks to be 32768, got: %d", opts.Chunks)
		}
	})

	t.Run("Given a chunk flag, When executing, Then should parse chunk", func(t *testing.T) {
		chunk := int64(100)
		cmd, opts := NewRootCmd("test", "test")
		cmd.SetArgs([]string{fmt.Sprintf("--chunks=%d", chunk), source, destination})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if opts.Chunks != chunk {
			t.Errorf("Expected chunks to be %d, got: %d", chunk, opts.Chunks)
		}
	})

	t.Run("Given skip and force flags, When executing, Then should parse correctly", func(t *testing.T) {
		cmd, opts := NewRootCmd("test", "test")
		cmd.SetArgs([]string{source, destination, "--skip", "--force"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !opts.Skip {
			t.Error("Expected Skip to be true")
		}
		if !opts.Force {
			t.Error("Expected Force to be true")
		}
	})

	t.Run("Given show and parallel flags, When executing, Then should parse correctly", func(t *testing.T) {
		cmd, opts := NewRootCmd("test", "test")
		cmd.SetArgs([]string{source, destination, "--show=10", "--parallel=8"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if opts.Show != 10 {
			t.Errorf("Expected Show to be 10, got %d", opts.Show)
		}
		if opts.Parallel != 8 {
			t.Errorf("Expected Parallel to be 8, got %d", opts.Parallel)
		}
	})

	t.Run("Given flags before positional args, When executing, Then should parse correctly", func(t *testing.T) {
		cmd, opts := NewRootCmd("test", "test")
		cmd.SetArgs([]string{"--skip", "--chunks=512", source, destination})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !opts.Skip {
			t.Error("Expected Skip to be true")
		}
		if opts.Chunks != 512 {
			t.Errorf("Expected Chunks to be 512, got %d", opts.Chunks)
		}
		if opts.Source != source {
			t.Errorf("Expected source to be %s, got: %s", source, opts.Source)
		}
		if opts.Destination != destination {
			t.Errorf("Expected destination to be %s, got: %s", destination, opts.Destination)
		}
	})
}
