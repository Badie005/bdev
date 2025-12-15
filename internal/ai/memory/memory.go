package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/badie/bdev/internal/ai/engine"
)

// Memory stores conversation history
type Memory struct {
	SessionID   string          `json:"session_id"`
	Messages    []StoredMessage `json:"messages"`
	MaxMessages int             `json:"-"`
	MaxTokens   int             `json:"-"`
	FilePath    string          `json:"-"`
}

// StoredMessage extends Message with metadata
type StoredMessage struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Tokens    int       `json:"tokens"`
	Timestamp time.Time `json:"timestamp"`
}

// New creates a new memory instance
func New(sessionID string, filePath string) *Memory {
	return &Memory{
		SessionID:   sessionID,
		Messages:    make([]StoredMessage, 0),
		MaxMessages: 20,
		MaxTokens:   4096,
		FilePath:    filePath,
	}
}

// Add adds a message to memory
func (m *Memory) Add(role, content string) {
	msg := StoredMessage{
		Role:      role,
		Content:   content,
		Tokens:    EstimateTokens(content),
		Timestamp: time.Now(),
	}

	m.Messages = append(m.Messages, msg)
	m.prune()
}

// GetContext returns messages for API call
func (m *Memory) GetContext() []engine.Message {
	messages := make([]engine.Message, len(m.Messages))
	for i, msg := range m.Messages {
		messages[i] = engine.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return messages
}

// GetContextWithSystem returns messages with a system prompt prepended
func (m *Memory) GetContextWithSystem(systemPrompt string) []engine.Message {
	messages := make([]engine.Message, 0, len(m.Messages)+1)

	// Add system prompt first
	messages = append(messages, engine.Message{
		Role:    "system",
		Content: systemPrompt,
	})

	// Add conversation history
	for _, msg := range m.Messages {
		messages = append(messages, engine.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return messages
}

// Clear clears all messages
func (m *Memory) Clear() {
	m.Messages = make([]StoredMessage, 0)
}

// Size returns the number of messages
func (m *Memory) Size() int {
	return len(m.Messages)
}

// TotalTokens returns estimated total tokens
func (m *Memory) TotalTokens() int {
	total := 0
	for _, msg := range m.Messages {
		total += msg.Tokens
	}
	return total
}

// prune removes old messages to stay within limits
func (m *Memory) prune() {
	// Prune by message count
	if len(m.Messages) > m.MaxMessages {
		m.Messages = m.Messages[len(m.Messages)-m.MaxMessages:]
	}

	// Prune by token count
	for m.TotalTokens() > m.MaxTokens && len(m.Messages) > 1 {
		m.Messages = m.Messages[1:]
	}
}

// Save persists memory to disk
func (m *Memory) Save() error {
	if m.FilePath == "" {
		return nil
	}

	// Ensure directory exists
	dir := filepath.Dir(m.FilePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.FilePath, data, 0o600)
}

// Load loads memory from disk
func (m *Memory) Load() error {
	if m.FilePath == "" {
		return nil
	}

	data, err := os.ReadFile(m.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file yet, that's OK
		}
		return err
	}

	return json.Unmarshal(data, m)
}

// LastMessage returns the last message if any
func (m *Memory) LastMessage() *StoredMessage {
	if len(m.Messages) == 0 {
		return nil
	}
	return &m.Messages[len(m.Messages)-1]
}

// EstimateTokens roughly estimates token count
// Uses ~4 chars per token approximation
func EstimateTokens(text string) int {
	return len(text) / 4
}

// Global memory instance
var globalMemory *Memory

// GetGlobal returns the global memory instance
func GetGlobal(filePath string) *Memory {
	if globalMemory == nil {
		globalMemory = New("default", filePath)
		_ = globalMemory.Load() // Ignore error: file may not exist yet
	}
	return globalMemory
}
