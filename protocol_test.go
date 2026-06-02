package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestReadWriteMessage(t *testing.T) {
	tests := []struct {
		name    string
		message []byte
	}{
		{
			name:    "Simple message",
			message: []byte(`{"status":"ok"}`),
		},
		{
			name:    "Empty message",
			message: []byte(""),
		},
		{
			name:    "Large message",
			message: bytes.Repeat([]byte("a"), 1024),
		},
		{
			name:    "Binary data",
			message: []byte{0x00, 0x01, 0x02, 0x03, 0xff, 0xfe},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			// Test WriteMessage
			err := WriteMessage(&buf, tt.message)
			if err != nil {
				t.Fatalf("WriteMessage failed: %v", err)
			}

			// Test ReadMessage
			got, err := ReadMessage(&buf)
			if err != nil {
				t.Fatalf("ReadMessage failed: %v", err)
			}

			if !reflect.DeepEqual(got, tt.message) {
				t.Errorf("ReadMessage() = %v, want %v", got, tt.message)
			}
		})
	}
}

func TestReadMessageError(t *testing.T) {
	t.Run("Short length prefix", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{1, 0}) // Only 2 bytes, needs 4
		_, err := ReadMessage(buf)
		if err == nil {
			t.Error("Expected error for short length prefix, got nil")
		}
	})

	t.Run("Short body", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{10, 0, 0, 0, 1, 2, 3}) // Length 10, but only 3 bytes
		_, err := ReadMessage(buf)
		if err == nil {
			t.Error("Expected error for short body, got nil")
		}
	})
}
