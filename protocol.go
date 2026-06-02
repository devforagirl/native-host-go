package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ReadMessage reads a Chrome Native Messaging message from the provided reader.
// The protocol consists of a 4-byte little-endian length prefix followed by the message.
func ReadMessage(r io.Reader) ([]byte, error) {
	var length uint32
	err := binary.Read(r, binary.LittleEndian, &length)
	if err != nil {
		return nil, fmt.Errorf("failed to read message length: %w", err)
	}

	message := make([]byte, length)
	_, err = io.ReadFull(r, message)
	if err != nil {
		return nil, fmt.Errorf("failed to read message body: %w", err)
	}

	return message, nil
}

// WriteMessage writes a Chrome Native Messaging message to the provided writer.
// The protocol consists of a 4-byte little-endian length prefix followed by the message.
func WriteMessage(w io.Writer, message []byte) error {
	length := uint32(len(message))
	err := binary.Write(w, binary.LittleEndian, length)
	if err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}

	_, err = w.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message body: %w", err)
	}

	return nil
}
