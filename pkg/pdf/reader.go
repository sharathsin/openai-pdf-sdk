package pdf

import (
	"bytes"
	"fmt"

	"github.com/ledongthuc/pdf"
)

// ExtractText reads a PDF file from the given path and returns its text content.
func ExtractText(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open pdf: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("failed to get plain text: %w", err)
	}
	buf.ReadFrom(b)

	return buf.String(), nil
}
