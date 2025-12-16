package engine

import "context"

// MockClient is a test mock implementation of the Engine interface
type MockClient struct {
	// Configurable return values
	IsAvailableFunc     func() bool
	HealthCheckFunc     func() bool
	ChatFunc            func(ctx context.Context, messages []Message) (<-chan string, <-chan error)
	ChatSyncFunc        func(ctx context.Context, messages []Message) (string, error)
	GenerateFunc        func(ctx context.Context, prompt, system string) (<-chan string, <-chan error)
	ListModelsFunc      func() ([]ModelInfo, error)
	HasModelFunc        func(name string) bool
	GetCurrentModelFunc func() string

	// Call tracking
	Calls []string

	// Default responses
	DefaultResponse  string
	DefaultModel     string
	DefaultAvailable bool
}

// Ensure MockClient implements Engine interface
var _ Engine = (*MockClient)(nil)

// NewMockClient creates a new mock client with sensible defaults
func NewMockClient() *MockClient {
	return &MockClient{
		DefaultResponse:  "Mock AI response",
		DefaultModel:     "mock-model",
		DefaultAvailable: true,
		Calls:            make([]string, 0),
	}
}

func (m *MockClient) IsAvailable() bool {
	m.Calls = append(m.Calls, "IsAvailable")
	if m.IsAvailableFunc != nil {
		return m.IsAvailableFunc()
	}
	return m.DefaultAvailable
}

func (m *MockClient) HealthCheck() bool {
	m.Calls = append(m.Calls, "HealthCheck")
	if m.HealthCheckFunc != nil {
		return m.HealthCheckFunc()
	}
	return m.DefaultAvailable
}

func (m *MockClient) Chat(ctx context.Context, messages []Message) (<-chan string, <-chan error) {
	m.Calls = append(m.Calls, "Chat")
	if m.ChatFunc != nil {
		return m.ChatFunc(ctx, messages)
	}

	// Default: return response in channel
	respChan := make(chan string, 1)
	errChan := make(chan error, 1)
	go func() {
		respChan <- m.DefaultResponse
		close(respChan)
		close(errChan)
	}()
	return respChan, errChan
}

func (m *MockClient) ChatSync(ctx context.Context, messages []Message) (string, error) {
	m.Calls = append(m.Calls, "ChatSync")
	if m.ChatSyncFunc != nil {
		return m.ChatSyncFunc(ctx, messages)
	}
	return m.DefaultResponse, nil
}

func (m *MockClient) Generate(ctx context.Context, prompt, system string) (<-chan string, <-chan error) {
	m.Calls = append(m.Calls, "Generate")
	if m.GenerateFunc != nil {
		return m.GenerateFunc(ctx, prompt, system)
	}

	// Default: return response in channel
	respChan := make(chan string, 1)
	errChan := make(chan error, 1)
	go func() {
		respChan <- m.DefaultResponse
		close(respChan)
		close(errChan)
	}()
	return respChan, errChan
}

func (m *MockClient) ListModels() ([]ModelInfo, error) {
	m.Calls = append(m.Calls, "ListModels")
	if m.ListModelsFunc != nil {
		return m.ListModelsFunc()
	}
	return []ModelInfo{
		{Name: m.DefaultModel},
	}, nil
}

func (m *MockClient) HasModel(name string) bool {
	m.Calls = append(m.Calls, "HasModel")
	if m.HasModelFunc != nil {
		return m.HasModelFunc(name)
	}
	return name == m.DefaultModel
}

func (m *MockClient) GetCurrentModel() string {
	m.Calls = append(m.Calls, "GetCurrentModel")
	if m.GetCurrentModelFunc != nil {
		return m.GetCurrentModelFunc()
	}
	return m.DefaultModel
}
