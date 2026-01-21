package pdf_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/jung-kurt/gofpdf"
	"github.com/sharathjillela/markets-sdk/openai-pdf-sdk/pkg/pdf"
)

// BenchmarkExtractText benchmarks the text extraction to verify optimizations.
func BenchmarkExtractText(b *testing.B) {
	// Setup: Generate a PDF once
	tempDir := b.TempDir()
	pdfPath := filepath.Join(tempDir, "bench.pdf")
	pdfGen := gofpdf.New("P", "mm", "A4", "")
	pdfGen.AddPage()
	pdfGen.SetFont("Arial", "B", 12)
	// Add substantial text to make allocation differences visible
	for i := 0; i < 100; i++ {
		pdfGen.Cell(40, 10, "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.")
		pdfGen.Ln(5)
	}
	if err := pdfGen.OutputFileAndClose(pdfPath); err != nil {
		b.Fatalf("Failed to generate PDF: %v", err)
	}

	b.ResetTimer()
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		_, err := pdf.ExtractText(ctx, pdfPath)
		if err != nil {
			b.Fatalf("ExtractText failed: %v", err)
		}
	}
}
