package engine

import "context"

// Engine defines the interface for AI engine operations.
// This interface enables dependency injection and mocking in tests.
type Engine interface {
	// IsAvailable checks if the AI service is reachable
	IsAvailable() bool

	// HealthCheck is an alias for IsAvailable
	HealthCheck() bool

	// Chat sends a streaming chat request
	// Returns channels for response chunks and errors
	Chat(ctx context.Context, messages []Message) (<-chan string, <-chan error)

	// ChatSync sends a non-streaming chat request
	// Returns the complete response
	ChatSync(ctx context.Context, messages []Message) (string, error)

	// Generate sends a streaming generate request
	Generate(ctx context.Context, prompt, system string) (<-chan string, <-chan error)

	// ListModels returns available models
	ListModels() ([]ModelInfo, error)

	// HasModel checks if a specific model is available
	HasModel(name string) bool

	// GetCurrentModel returns the active model name
	GetCurrentModel() string
}

// Ensure Client implements Engine interface
var _ Engine = (*Client)(nil)

// StreamHandler is a callback type for processing streaming responses
type StreamHandler func(chunk string) error

// EngineWithStreaming extends Engine with callback-style streaming
type EngineWithStreaming interface {
	Engine
	ChatWithHandler(ctx context.Context, messages []Message, handler StreamHandler) error
	GenerateWithHandler(ctx context.Context, prompt, system string, handler StreamHandler) error
}
