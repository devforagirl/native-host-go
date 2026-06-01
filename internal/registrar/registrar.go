package registrar

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Manifest defines the structure of the browser's native messaging host manifest file.
type Manifest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Type        string `json:"type"`
}

// Registrar defines the interface for registering the native messaging host on different operating systems.
type Registrar interface {
	Register(hostName string, hostPath string) error
}

// NewRegistrar returns the OS-specific implementation of the Registrar.
func NewRegistrar() Registrar {
	return newRegistrar()
}

// writeManifest writes the manifest JSON to the specified path.
func writeManifest(path string, manifest Manifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create manifest directory %s: %w", dir, err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file %s: %w", path, err)
	}

	return nil
}
