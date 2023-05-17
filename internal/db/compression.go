package db

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"strings"
)

// compress compresses a string using zlib
func compress(text string) (string, error) {
	buf := new(bytes.Buffer)

	writer := zlib.NewWriter(buf)
	_, err := writer.Write([]byte(text))
	if err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// decompress decompresses a string using zlib
func decompress(zlibString string) (string, error) {
	reader, err := zlib.NewReader(strings.NewReader(zlibString))
	if err != nil {
		return "", fmt.Errorf("opening zlib reader: %w", err)
	}

	out, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("reading from zlib reader: %w", err)
	}

	if err := reader.Close(); err != nil {
		return "", err
	}

	return string(out), nil
}
