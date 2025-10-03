package analyzer

import (
	"fmt"
	"testing"
)

func TestGetSizeFormatted(t *testing.T) {

	t.Run("when size is lower than 1024 Then should return value in bytes", func(t *testing.T) {
		size := int64(900)
		res := GetSizeReadable(size)

		expected := fmt.Sprintf("%dB", size)
		if res != expected {
			t.Errorf("Expected ")
		}
	})

	t.Run("when size is 1024 Then should return 1kb", func(t *testing.T) {
		size := int64(1024)
		res := GetSizeReadable(size)

		expected := "1KB"
		if res != expected {
			t.Errorf("Expected ")
		}
	})

	t.Run("when size bigger than 1024 Then should return value in kb", func(t *testing.T) {
		size := int64(1500)
		res := GetSizeReadable(size)

		expected := "1.5KB"
		if res != expected {
			t.Errorf("Expected ")
		}
	})

}
