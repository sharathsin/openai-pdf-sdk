package openai

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	api *openai.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		api: openai.NewClient(apiKey),
	}
}

// UploadFile uploads a file to OpenAI for a specific purpose (e.g., "assistants", "fine-tune").
func (c *Client) UploadFile(ctx context.Context, path string, purpose string) (openai.File, error) {
	req := openai.FileRequest{
		FilePath: path,
		Purpose:  purpose,
	}

	file, err := c.api.CreateFile(ctx, req)
	if err != nil {
		return openai.File{}, fmt.Errorf("failed to upload file: %w", err)
	}

	return file, nil
}

// SendText sends a prompt to the OpenAI Chat Completion API.
// This can be used to send the text extracted from the PDF.
func (c *Client) SendText(ctx context.Context, text string) (string, error) {
	resp, err := c.api.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("chat completion error: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
