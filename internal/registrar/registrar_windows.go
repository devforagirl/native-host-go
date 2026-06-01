//go:build windows

package registrar

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

type windowsRegistrar struct{}

func newRegistrar() Registrar {
	return &windowsRegistrar{}
}

func (r *windowsRegistrar) Register(hostName string, hostPath string) error {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return fmt.Errorf("APPDATA environment variable not set")
	}

	manifestPath := filepath.Join(appData, "FlowMeter", "manifest.json")
	manifest := Manifest{
		Name:        hostName,
		Description: "FlowMeter Native Host",
		Path:        hostPath,
		Type:        "stdio",
	}

	if err := writeManifest(manifestPath, manifest); err != nil {
		return err
	}

	// Registry key for Chrome Native Messaging Hosts
	keyPath := fmt.Sprintf(`Software\Google\Chrome\NativeMessagingHosts\%s`, hostName)
	k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			k, _, err = registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
			if err != nil {
				return fmt.Errorf("failed to create registry key %s: %w", keyPath, err)
			}
		} else {
			return fmt.Errorf("failed to open registry key %s: %w", keyPath, err)
		}
	}
	defer k.Close()

	if err := k.SetStringValue("", manifestPath); err != nil {
		return fmt.Errorf("failed to set registry value for %s: %w", keyPath, err)
	}

	return nil
}
