package config

import (
	"testing"
)

func TestDefaultWindowState(t *testing.T) {
	def := DefaultWindowState()
	if def.Width != 1280 || def.Height != 800 {
		t.Errorf("Unexpected default window dimensions: %dx%d", def.Width, def.Height)
	}
	if def.IsMaximized {
		t.Errorf("Default window should not be maximized")
	}
}

func TestSaveAndLoadWindowState(t *testing.T) {
	testState := WindowState{
		X:           150,
		Y:           120,
		Width:       1400,
		Height:      900,
		IsMaximized: true,
	}

	err := SaveWindowState(testState)
	if err != nil {
		t.Fatalf("SaveWindowState failed: %v", err)
	}

	loaded := LoadWindowState()
	if loaded.X != testState.X || loaded.Y != testState.Y || loaded.Width != testState.Width || loaded.Height != testState.Height || loaded.IsMaximized != testState.IsMaximized {
		t.Errorf("Loaded state %+v did not match saved state %+v", loaded, testState)
	}
}
