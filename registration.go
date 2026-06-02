package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/windows/registry"
)

type Manifest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Path            string   `json:"path"`
	Type            string   `json:"type"`
	AllowedOrigins  []string `json:"allowed_origins"`
}

func HandleRegistration() {
	fmt.Println(`==================================================`)
	fmt.Println(`          FlowMeter Native Host Registration`)
	fmt.Println(`==================================================`)
	fmt.Println(`\nConfiguring environment...`)

	err := performRegistration()
	if err != nil {
		fmt.Printf(`❌ Registration failed: %v\n`, err)
		fmt.Println(`\nPlease check if you have the necessary permissions.`)
	} else {
		fmt.Println(`✅ Registration successful!`)
		fmt.Println(`\nNext steps:`)
		fmt.Println(`1. Open Chrome Browser`)
		fmt.Println(`2. Click the FlowMeter extension icon`)
		fmt.Println(`3. Click "Connect" to start monitoring`)
	}

	fmt.Println(`\n==================================================`)
	fmt.Println(`Press any key to exit...`)
	bufio.NewReader(os.Stdin).ReadByte()
}

func performRegistration() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	manifest := Manifest{
		Name:           "com.flowmeter.host",
		Description:    "FlowMeter Native Host",
		Path:           exePath,
		Type:           "stdio",
		AllowedOrigins: []string{"chrome-extension://omhgobopmdmnbcanhbcpfcdaphllgbkk/"},
	}

	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if runtime.GOOS == "windows" {
		return registerWindows(manifestBytes)
	} else if runtime.GOOS == "darwin" {
		return registerMacOS(manifestBytes)
	}

	return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
}

func registerWindows(manifestBytes []byte) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %w", err)
	}

	configDir := filepath.Join(home, "AppData", "Local", "FlowMeter")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	manifestPath := filepath.Join(configDir, "manifest.json")
	if err := os.WriteFile(manifestPath, manifestBytes, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	keyPath := `Software\Google\Chrome\NativeMessagingHosts\com.flowmeter.host`
	k, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to create registry key: %w", err)
	}
	defer k.Close()

	if err := k.SetStringValue("", manifestPath); err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	return nil
}

func registerMacOS(manifestBytes []byte) error {
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
