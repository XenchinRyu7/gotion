package profile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetUserDataDir(t *testing.T) {
	dir, err := GetUserDataDir()
	if err != nil {
		t.Fatalf("GetUserDataDir() returned error: %v", err)
	}

	if dir == "" {
		t.Fatalf("GetUserDataDir() returned empty path")
	}

	if !strings.Contains(strings.ToLower(dir), "gotion") {
		t.Errorf("GetUserDataDir() path %q should contain 'gotion'", dir)
	}

	// Verify directory was created
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Failed to stat user data dir %s: %v", dir, err)
	}

	if !info.IsDir() {
		t.Errorf("Path %s is not a directory", dir)
	}

	// Verify clean path
	cleaned := filepath.Clean(dir)
	if cleaned != dir {
		t.Errorf("Path is not clean: %s vs %s", dir, cleaned)
	}
}
