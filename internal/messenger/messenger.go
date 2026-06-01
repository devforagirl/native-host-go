package messenger

import (
	"encoding/binary"
	"fmt"
	"io"
)

// MaxMessageSize is the maximum allowed message size to prevent DoS attacks.
const MaxMessageSize = 1 * 1024 * 1024 // 1MB

// ReadMessage reads a message from the given reader according to the Chrome Native Messaging protocol.
func ReadMessage(r io.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	if length > MaxMessageSize {
		return nil, fmt.Errorf("message too large: %d", length)
	}

	msg := make([]byte, length)
	if _, err := io.ReadFull(r, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

// SendMessage writes a message to the given writer according to the Chrome Native Messaging protocol.
func SendMessage(w io.Writer, msg []byte) error {
	if err := binary.Write(w, binary.LittleEndian, uint32(len(msg))); err != nil {
		return err
	}
	_, err := w.Write(msg)
	return err
}
