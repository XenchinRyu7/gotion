package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// WindowState stores window dimensions, position, and maximized state.
type WindowState struct {
	X           int  `json:"x"`
	Y           int  `json:"y"`
	Width       int  `json:"width"`
	Height      int  `json:"height"`
	IsMaximized bool `json:"is_maximized"`
}

// DefaultWindowState returns sensible initial window parameters.
func DefaultWindowState() WindowState {
	return WindowState{
		X:           -1, // -1 indicates window should be centered
		Y:           -1,
		Width:       1280,
		Height:      800,
		IsMaximized: false,
	}
}

func getConfigPath() (string, error) {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to locate config directory: %w", err)
		}
		localAppData = configDir
	}

	configDir := filepath.Join(localAppData, "Gotion", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	return filepath.Join(configDir, "window.json"), nil
}

// LoadWindowState loads saved window geometry or returns defaults.
func LoadWindowState() WindowState {
	def := DefaultWindowState()

	path, err := getConfigPath()
	if err != nil {
		return def
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return def
	}

	var state WindowState
	if err := json.Unmarshal(data, &state); err != nil {
		return def
	}

	// Validate min bounds
	if state.Width < 600 {
		state.Width = def.Width
	}
	if state.Height < 400 {
		state.Height = def.Height
	}

	return state
}

// SaveWindowState writes window state to disk.
func SaveWindowState(state WindowState) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize window state: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
