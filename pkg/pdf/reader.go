package pdf

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/ledongthuc/pdf"
)

// bufferPool reuses bytes.Buffers to reduce GC pressure during text extraction.
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// ExtractText reads a PDF file from the given path and returns its text content.
// It accepts a context for cancellation (though the underlying library might not support it per-page,
// we can check it before processing).
func ExtractText(ctx context.Context, path string) (string, error) {
	// Check context early
	if err := ctx.Err(); err != nil {
		return "", err
	}

	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open pdf: %w", err)
	}
	defer f.Close()

	// Get a buffer from the pool
	buf := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufferPool.Put(buf)
	}()

	// Since ledongthuc/pdf doesn't take context in GetPlainText, we verify context here again
	if err := ctx.Err(); err != nil {
		return "", err
	}

	b, err := r.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("failed to get plain text: %w", err)
	}

	// This ReadFrom is where the buffer is used.
	// We optimize by not creating a new buffer every time, but reusing the one from pool.
	_, err = buf.ReadFrom(b)
	if err != nil {
		return "", fmt.Errorf("failed to read text to buffer: %w", err)
	}

	return buf.String(), nil
}
