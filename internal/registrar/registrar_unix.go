//go:build !windows

package registrar

import (
	"fmt"
	"os"
	"path/filepath"
)

type unixRegistrar struct{}

func newRegistrar() Registrar {
	return &unixRegistrar{}
}

func (r *unixRegistrar) Register(hostName string, hostPath string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Path for Chrome Native Messaging Hosts on Unix
	manifestPath := filepath.Join(home, ".config", "google-chrome", "NativeMessagingHosts", fmt.Sprintf("%s.json", hostName))

	manifest := Manifest{
		Name:        hostName,
		Description: "FlowMeter Native Host",
		Path:        hostPath,
		Type:        "stdio",
	}

	if err := writeManifest(manifestPath, manifest); err != nil {
		return err
	}

	return nil
}
