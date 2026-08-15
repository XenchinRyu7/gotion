package profile

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetUserDataDir returns the persistent user data path for WebView2 on Windows.
// By default, it resolves to %LOCALAPPDATA%\Gotion\profile.
func GetUserDataDir() (string, error) {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		// Fallback to user config dir or APPDATA
		configDir, err := os.UserConfigDir()
		if err != nil {
			appData := os.Getenv("APPDATA")
			if appData != "" {
				localAppData = appData
			} else {
				return "", fmt.Errorf("failed to locate LOCALAPPDATA or UserConfigDir: %w", err)
			}
		} else {
			localAppData = configDir
		}
	}

	profileDir := filepath.Join(localAppData, "Gotion", "profile")

	// Ensure the directory exists
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create user data directory at %s: %w", profileDir, err)
	}

	return profileDir, nil
}
