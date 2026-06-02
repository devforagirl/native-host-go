//go:build darwin

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func doPlatformRegistration(manifestBytes []byte) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %w", err)
	}

	configDir := filepath.Join(home, "Library", "Application Support", "Google", "Chrome", "NativeMessagingHosts")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	manifestPath := filepath.Join(configDir, "com.flowmeter.host.json")
	if err := os.WriteFile(manifestPath, manifestBytes, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	return nil
}
