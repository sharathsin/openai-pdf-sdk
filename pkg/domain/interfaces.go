package domain

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// TextExtractor defines the behavior for extracting text from a source.
type TextExtractor interface {
	ExtractText(ctx context.Context, path string) (string, error)
}

// AIClient defines the behavior for interacting with an AI provider.
type AIClient interface {
	UploadFile(ctx context.Context, path string, purpose string) (openai.File, error)
	SendText(ctx context.Context, text string) (string, error)
}
