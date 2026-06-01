package messenger

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"testing"
)

func TestReadMessage(t *testing.T) {
	t.Run("NormalMessage", func(t *testing.T) {
		msg := []byte(`{"status":"ok"}`)
		length := uint32(len(msg))

		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, length)
		buf.Write(msg)

		got, err := ReadMessage(buf)
		if err != nil {
			t.Fatalf("ReadMessage failed: %v", err)
		}

		if !reflect.DeepEqual(got, msg) {
			t.Errorf("ReadMessage = %v, want %v", got, msg)
		}
	})

	t.Run("EmptyMessage", func(t *testing.T) {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, uint32(0))

		got, err := ReadMessage(buf)
		if err != nil {
			t.Fatalf("ReadMessage failed: %v", err)
		}

		if len(got) != 0 {
			t.Errorf("ReadMessage = %v, want empty slice", got)
		}
	})

	t.Run("PartialRead", func(t *testing.T) {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, uint32(10))
		buf.Write([]byte("too short")) // only 9 bytes

		_, err := ReadMessage(buf)
		if err != io.ErrUnexpectedEOF {
			t.Errorf("ReadMessage error = %v, want %v", err, io.ErrUnexpectedEOF)
		}
	})

	t.Run("OversizedMessage", func(t *testing.T) {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, uint32(MaxMessageSize+1))

		_, err := ReadMessage(buf)
		if err == nil {
			t.Fatal("ReadMessage should have failed for oversized message")
		}
		expectedErr := "message too large"
		if !contains(err.Error(), expectedErr) {
			t.Errorf("ReadMessage error = %v, want it to contain %q", err, expectedErr)
		}
	})
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

func TestSendMessage(t *testing.T) {
	msg := []byte(`{"status":"ok"}`)
	buf := new(bytes.Buffer)

	err := SendMessage(buf, msg)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	// Verify first 4 bytes are the little-endian length
	var length uint32
	if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
		t.Fatalf("Failed to read length: %v", err)
	}

	if length != uint32(len(msg)) {
		t.Errorf("SendMessage length = %v, want %v", length, len(msg))
	}

	// Verify the message content
	got := buf.Bytes()
	if !reflect.DeepEqual(got, msg) {
		t.Errorf("SendMessage content = %v, want %v", got, msg)
	}
}
