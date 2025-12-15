package engine

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"` // system, user, assistant
	Content string `json:"content"`
}

// ChatRequest represents an Ollama chat request
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Options  *Options  `json:"options,omitempty"`
}

// Options represents model options
type Options struct {
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
}

// ChatResponse represents a streaming chat response chunk
type ChatResponse struct {
	Model     string  `json:"model"`
	Message   Message `json:"message"`
	Done      bool    `json:"done"`
	CreatedAt string  `json:"created_at,omitempty"`
	// Final response fields
	TotalDuration      int64 `json:"total_duration,omitempty"`
	LoadDuration       int64 `json:"load_duration,omitempty"`
	PromptEvalCount    int   `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64 `json:"prompt_eval_duration,omitempty"`
	EvalCount          int   `json:"eval_count,omitempty"`
	EvalDuration       int64 `json:"eval_duration,omitempty"`
}

// ModelInfo represents model information
type ModelInfo struct {
	Name       string    `json:"name"`
	ModifiedAt time.Time `json:"modified_at"`
	Size       int64     `json:"size"`
}

// ModelsResponse represents the list models response
type ModelsResponse struct {
	Models []ModelInfo `json:"models"`
}

// GenerateRequest represents a generate request
type GenerateRequest struct {
	Model   string   `json:"model"`
	Prompt  string   `json:"prompt"`
	Stream  bool     `json:"stream"`
	System  string   `json:"system,omitempty"`
	Options *Options `json:"options,omitempty"`
}

// GenerateResponse represents a generate response chunk
type GenerateResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at,omitempty"`
}

// Client is the Ollama API client
type Client struct {
	BaseURL     string
	Model       string
	Fallback    string
	HTTPClient  *http.Client
	Timeout     time.Duration
	MaxRetries  int
	Temperature float64
}

// Config holds client configuration
type Config struct {
	BaseURL     string
	Model       string
	Fallback    string
	Timeout     time.Duration
	MaxRetries  int
	Temperature float64
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		BaseURL:     "http://localhost:11434",
		Model:       "llama3.2",
		Fallback:    "phi3:mini",
		Timeout:     120 * time.Second,
		MaxRetries:  3,
		Temperature: 0.7,
	}
}

// New creates a new Ollama client
func New(cfg Config) *Client {
	return &Client{
		BaseURL:     cfg.BaseURL,
		Model:       cfg.Model,
		Fallback:    cfg.Fallback,
		Timeout:     cfg.Timeout,
		MaxRetries:  cfg.MaxRetries,
		Temperature: cfg.Temperature,
		HTTPClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// IsAvailable checks if Ollama is running
func (c *Client) IsAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/api/tags", nil)
	if err != nil {
		return false
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// HealthCheck is an alias for IsAvailable to satisfy new interface requirements
func (c *Client) HealthCheck() bool {
	return c.IsAvailable()
}

// ListModels returns available models
func (c *Client) ListModels() ([]ModelInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/api/tags", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	var result ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse models: %w", err)
	}

	return result.Models, nil
}

// HasModel checks if a specific model is available
func (c *Client) HasModel(name string) bool {
	models, err := c.ListModels()
	if err != nil {
		return false
	}

	for _, m := range models {
		if m.Name == name {
			return true
		}
	}
	return false
}

// Chat sends a chat request with streaming
func (c *Client) Chat(ctx context.Context, messages []Message) (<-chan string, <-chan error) {
	output := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(output)
		defer close(errChan)

		model := c.Model
		if !c.HasModel(model) && c.Fallback != "" {
			model = c.Fallback
		}

		req := ChatRequest{
			Model:    model,
			Messages: messages,
			Stream:   true,
			Options: &Options{
				Temperature: c.Temperature,
			},
		}

		err := c.streamChat(ctx, req, output)
		if err != nil {
			errChan <- err
		}
	}()

	return output, errChan
}

// ChatSync sends a non-streaming chat request
func (c *Client) ChatSync(ctx context.Context, messages []Message) (string, error) {
	model := c.Model
	if !c.HasModel(model) && c.Fallback != "" {
		model = c.Fallback
	}

	req := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
		Options: &Options{
			Temperature: c.Temperature,
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Message.Content, nil
}

// streamChat handles streaming chat response
func (c *Client) streamChat(ctx context.Context, req ChatRequest, output chan<- string) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Ollama error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		var chunk ChatResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			continue // Skip malformed lines
		}

		if chunk.Message.Content != "" {
			output <- chunk.Message.Content
		}

		if chunk.Done {
			return nil
		}
	}
}

// Generate sends a generate request with streaming
func (c *Client) Generate(ctx context.Context, prompt string, system string) (<-chan string, <-chan error) {
	output := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(output)
		defer close(errChan)

		model := c.Model
		if !c.HasModel(model) && c.Fallback != "" {
			model = c.Fallback
		}

		req := GenerateRequest{
			Model:  model,
			Prompt: prompt,
			Stream: true,
			System: system,
			Options: &Options{
				Temperature: c.Temperature,
			},
		}

		err := c.streamGenerate(ctx, req, output)
		if err != nil {
			errChan <- err
		}
	}()

	return output, errChan
}

// streamGenerate handles streaming generate response
func (c *Client) streamGenerate(ctx context.Context, req GenerateRequest, output chan<- string) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		var chunk GenerateResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			continue
		}

		if chunk.Response != "" {
			output <- chunk.Response
		}

		if chunk.Done {
			return nil
		}
	}
}

// GetCurrentModel returns the active model name
func (c *Client) GetCurrentModel() string {
	if c.HasModel(c.Model) {
		return c.Model
	}
	if c.HasModel(c.Fallback) {
		return c.Fallback
	}
	return "none"
}
