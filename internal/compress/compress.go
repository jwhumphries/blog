package compress

import (
	"bytes"

	"github.com/andybalholm/brotli"
)

// Brotli compresses data using Brotli compression at the default quality level.
func Brotli(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := brotli.NewWriterLevel(&buf, brotli.DefaultCompression)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
