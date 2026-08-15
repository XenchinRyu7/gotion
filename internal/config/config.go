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

// SessionConfig holds persistent Notion authentication tokens.
type SessionConfig struct {
	TokenV2 string `json:"token_v2"`
}

func getSessionConfigPath() (string, error) {
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

	return filepath.Join(configDir, "session.json"), nil
}

// LoadSessionToken returns the saved session token or empty string.
func LoadSessionToken() string {
	path, err := getSessionConfigPath()
	if err != nil {
		return ""
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	var session SessionConfig
	if err := json.Unmarshal(data, &session); err != nil {
		return ""
	}

	return session.TokenV2
}

// SaveSessionToken persists the session token to disk.
func SaveSessionToken(token string) error {
	path, err := getSessionConfigPath()
	if err != nil {
		return err
	}

	session := SessionConfig{TokenV2: token}
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize session: %w", err)
	}

	return os.WriteFile(path, data, 0600)
}

// ClearSessionToken removes the saved session token.
func ClearSessionToken() error {
	path, err := getSessionConfigPath()
	if err != nil {
		return err
	}
	_ = os.Remove(path)
	return nil
}

