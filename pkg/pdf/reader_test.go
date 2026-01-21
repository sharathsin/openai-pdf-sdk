package pdf_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jung-kurt/gofpdf"
	"github.com/sharathjillela/markets-sdk/openai-pdf-sdk/pkg/pdf"
)

func TestExtractText(t *testing.T) {
	// 1. Generate a temporary PDF file
	tempDir := t.TempDir()
	pdfPath := filepath.Join(tempDir, "test.pdf")

	pdfGen := gofpdf.New("P", "mm", "A4", "")
	pdfGen.AddPage()
	pdfGen.SetFont("Arial", "B", 16)
	pdfGen.Cell(40, 10, "Hello World")
	err := pdfGen.OutputFileAndClose(pdfPath)
	if err != nil {
		t.Fatalf("Failed to generate test PDF: %v", err)
	}

	// 2. Test extraction
	text, err := pdf.ExtractText(pdfPath)
	if err != nil {
		t.Fatalf("ExtractText returned error: %v", err)
	}

	// 3. Verify content
	// Note: PDF extraction might contain extra whitespace/newlines
	if !strings.Contains(text, "Hello World") {
		t.Errorf("Expected text to contain 'Hello World', got: %q", text)
	}
}

func TestExtractText_FileNotFound(t *testing.T) {
	_, err := pdf.ExtractText("non-existent-file.pdf")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
