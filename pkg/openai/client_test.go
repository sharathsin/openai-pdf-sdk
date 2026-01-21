package openai_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sharathjillela/markets-sdk/openai-pdf-sdk/pkg/openai"
)

func TestNewClient(t *testing.T) {
	client := openai.NewClient("test-key")
	if client == nil {
		t.Error("NewClient returned nil")
	}
}

func TestIntegration_UploadAndChat(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: OPENAI_API_KEY not set")
	}

	client := openai.NewClient(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Test SendText (Chat Completion)
	t.Run("SendText", func(t *testing.T) {
		response, err := client.SendText(ctx, "Say 'test' and nothing else.")
		if err != nil {
			t.Fatalf("SendText failed: %v", err)
		}
		if len(response) == 0 {
			t.Error("Received empty response from OpenAI")
		}
		t.Logf("OpenAI Response: %s", response)
	})

	// 2. Test UploadFile
	t.Run("UploadFile", func(t *testing.T) {
		// Create a temporary file to upload
		tempDir := t.TempDir()
		tempFile := filepath.Join(tempDir, "test_upload.jsonl")
		// Creating a valid JSONL file for fine-tuning or just a text file for assistants
		content := []byte(`{"messages": [{"role": "system", "content": "Marv is a factual chatbot that is also sarcastic."}, {"role": "user", "content": "What's the capital of France?"}, {"role": "assistant", "content": "Paris, as if everyone doesn't know that already."}]}`)
		if err := os.WriteFile(tempFile, content, 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Use "fine-tune" as purpose since it's a standard one, or "assistants"
		file, err := client.UploadFile(ctx, tempFile, "fine-tune")
		if err != nil {
			t.Fatalf("UploadFile failed: %v", err)
		}

		if file.ID == "" {
			t.Error("Returned file ID is empty")
		}
		t.Logf("Uploaded File ID: %s", file.ID)
	})
}
