package openai

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v4"
	openai "github.com/sashabaranov/go-openai"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"
)

// Client wraps the OpenAI API client with resilience and observability.
type Client struct {
	api         *openai.Client
	logger      *slog.Logger
	rateLimiter *rate.Limiter
	tracer      trace.Tracer
}

// Option allows customizing the Client.
type Option func(*Client)

// WithLogger sets the logger for the client.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithRateLimit sets the rate limit (requests per second) and burst size.
func WithRateLimit(r rate.Limit, b int) Option {
	return func(c *Client) {
		c.rateLimiter = rate.NewLimiter(r, b)
	}
}

// NewClient creates a new OpenAI client with default resilience patterns.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		api: openai.NewClient(apiKey),
		// Default no-op logger if not provided
		logger: slog.New(slog.Default().Handler()),
		// Default rate limit: 10 RPS, burst 20 (conservative default)
		rateLimiter: rate.NewLimiter(rate.Limit(10), 20),
		tracer:      otel.Tracer("openai-pdf-sdk/pkg/openai"),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// UploadFile uploads a file to OpenAI with retries and tracing.
func (c *Client) UploadFile(ctx context.Context, path string, purpose string) (openai.File, error) {
	ctx, span := c.tracer.Start(ctx, "UploadFile",
		trace.WithAttributes(
			attribute.String("file.path", path),
			attribute.String("file.purpose", purpose),
		),
	)
	defer span.End()

	if err := c.waitRateLimit(ctx); err != nil {
		return openai.File{}, err
	}

	var file openai.File
	operation := func() error {
		var err error
		c.logger.DebugContext(ctx, "Attempting file upload", "path", path)
		file, err = c.api.CreateFile(ctx, openai.FileRequest{
			FilePath: path,
			Purpose:  purpose,
		})
		if err != nil {
			// Check if error is permanent (e.g. invalid auth), wrap if needed
			// For simplicity assuming network/server errors are transient
			c.logger.WarnContext(ctx, "Upload failed, retrying", "error", err)
			return err
		}
		return nil
	}

	// Exponential backoff: max interval 10s, max time 1 min
	b := backoff.NewExponentialBackOff()
	b.MaxInterval = 10 * time.Second
	b.MaxElapsedTime = 1 * time.Minute

	if err := backoff.Retry(operation, backoff.WithContext(b, ctx)); err != nil {
		span.RecordError(err)
		c.logger.ErrorContext(ctx, "File upload failed permanently", "error", err)
		return openai.File{}, fmt.Errorf("failed to upload file after retries: %w", err)
	}

	c.logger.InfoContext(ctx, "File uploaded successfully", "file_id", file.ID)
	return file, nil
}

// SendText sends a prompt to OpenAI with retries, rate limiting, and observability.
func (c *Client) SendText(ctx context.Context, text string) (string, error) {
	ctx, span := c.tracer.Start(ctx, "SendText",
		trace.WithAttributes(attribute.Int("text.length", len(text))),
	)
	defer span.End()

	if err := c.waitRateLimit(ctx); err != nil {
		return "", err
	}

	var content string
	operation := func() error {
		c.logger.DebugContext(ctx, "Sending chat completion request")
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
			var apiErr *openai.APIError
			if errors.As(err, &apiErr) {
				// Don't retry on 400s (invalid request) or 401 (auth)
				if apiErr.HTTPStatusCode >= 400 && apiErr.HTTPStatusCode < 500 && apiErr.HTTPStatusCode != 429 {
					return backoff.Permanent(err)
				}
			}
			c.logger.WarnContext(ctx, "Chat completion failed, retrying", "error", err)
			return err
		}
		if len(resp.Choices) == 0 {
			return errors.New("empty choices in response")
		}
		content = resp.Choices[0].Message.Content
		return nil
	}

	b := backoff.NewExponentialBackOff()
	b.MaxInterval = 5 * time.Second
	b.MaxElapsedTime = 30 * time.Second

	if err := backoff.Retry(operation, backoff.WithContext(b, ctx)); err != nil {
		span.RecordError(err)
		c.logger.ErrorContext(ctx, "SendText failed permanently", "error", err)
		return "", fmt.Errorf("chat completion error: %w", err)
	}

	c.logger.InfoContext(ctx, "Chat completion successful")
	return content, nil
}

func (c *Client) waitRateLimit(ctx context.Context) error {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		c.logger.ErrorContext(ctx, "Rate limiter error", "error", err)
		return fmt.Errorf("rate limit exceeded: %w", err)
	}
	return nil
}
