package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
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

	return doPlatformRegistration(manifestBytes)
}
