package registrar

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteManifest(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "manifest.json")
	manifest := Manifest{
		Name:        "test.host",
		Description: "Test Host",
		Path:        "/tmp/test-host",
		Type:        "stdio",
	}

	err := writeManifest(manifestPath, manifest)
	if err != nil {
		t.Fatalf("writeManifest failed: %v", err)
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("failed to read manifest file: %v", err)
	}

	var decoded Manifest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal manifest: %v", err)
	}

	if decoded.Name != manifest.Name {
		t.Errorf("expected name %s, got %s", manifest.Name, decoded.Name)
	}
	if decoded.Description != manifest.Description {
		t.Errorf("expected description %s, got %s", manifest.Description, decoded.Description)
	}
	if decoded.Path != manifest.Path {
		t.Errorf("expected path %s, got %s", manifest.Path, decoded.Path)
	}
	if decoded.Type != manifest.Type {
		t.Errorf("expected type %s, got %s", manifest.Type, decoded.Type)
	}
}

func TestNewRegistrar(t *testing.T) {
	reg := NewRegistrar()
	if reg == nil {
		t.Fatal("NewRegistrar returned nil")
	}
}
