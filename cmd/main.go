package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sharathjillela/markets-sdk/openai-pdf-sdk/pkg/openai"
	"github.com/sharathjillela/markets-sdk/openai-pdf-sdk/pkg/pdf"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}

	filePath := flag.String("file", "", "Path to the PDF file")
	purpose := flag.String("purpose", "assistants", "Purpose of the file upload (e.g., assistants, fine-tune)")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("Please provide a file path using -file flag")
	}

	fmt.Printf("Processing file: %s\n", *filePath)

	// 1. Extract Text
	start := time.Now()
	text, err := pdf.ExtractText(*filePath)
	if err != nil {
		log.Fatalf("Error extracting text: %v", err)
	}
	fmt.Printf("Successfully extracted %d characters in %v\n", len(text), time.Since(start))
	if len(text) > 200 {
		fmt.Printf("Preview: %s...\n", text[:200])
	} else {
		fmt.Printf("Preview: %s\n", text)
	}

	// 2. Upload File
	client := openai.NewClient(apiKey)
	fmt.Println("\nUploading file to OpenAI...")
	start = time.Now()
	file, err := client.UploadFile(context.Background(), *filePath, *purpose)
	if err != nil {
		log.Fatalf("Error uploading file: %v", err)
	}
	fmt.Printf("Successfully uploaded file in %v\n", time.Since(start))
	fmt.Printf("File ID: %s\n", file.ID)
	fmt.Printf("Filename: %s\n", file.FileName)
	fmt.Printf("Status: %s\n", file.Status)

	// 3. Send Summary Request (Optional demo)
	fmt.Println("\nSending text snippet to OpenAI for summary...")
	snippet := text
	if len(snippet) > 2000 {
		snippet = snippet[:2000] // Truncate for demo to avoid token limits
	}
	prompt := fmt.Sprintf("Please summarize this text from a PDF: %s", snippet)
	summary, err := client.SendText(context.Background(), prompt)
	if err != nil {
		log.Printf("Error getting summary: %v", err)
	} else {
		fmt.Printf("\nSummary from OpenAI:\n%s\n", summary)
	}
}
