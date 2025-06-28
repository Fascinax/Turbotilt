package runtime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanupTempFiles(t *testing.T) {
	// Create temporary files for the test
	tempDir := t.TempDir()
	tempFiles := []string{
		filepath.Join(tempDir, "file1.tmp"),
		filepath.Join(tempDir, "file2.tmp"),
	}

	// Create the files
	for _, file := range tempFiles {
		if err := os.WriteFile(file, []byte("test content"), 0644); err != nil {
			t.Fatalf("Unable to create temporary file: %v", err)
		}
	}

	// Verify the files exist
	for _, file := range tempFiles {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("File %s should exist: %v", file, err)
		}
	}

	// Execute the cleanup function
	CleanupTempFiles(tempFiles)

	// Verify the files have been deleted
	for _, file := range tempFiles {
		if _, err := os.Stat(file); err == nil {
			t.Errorf("File %s should have been deleted", file)
		} else if !os.IsNotExist(err) {
			t.Errorf("Unexpected error when checking file %s: %v", file, err)
		}
	}
}
