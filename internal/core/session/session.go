package session

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Session represents the current CLI session state
type Session struct {
	StartTime    time.Time         `json:"start_time"`
	CurrentDir   string            `json:"current_dir"`
	CommandCount int               `json:"command_count"`
	AIContext    []Message         `json:"ai_context"`
	LastCommand  string            `json:"last_command"`
	Metadata     map[string]string `json:"metadata"`
}

// Message represents an AI conversation message
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Global session instance
var globalSession *Session

// New creates a new session
func New() *Session {
	cwd, _ := os.Getwd()
	return &Session{
		StartTime:    time.Now(),
		CurrentDir:   cwd,
		CommandCount: 0,
		AIContext:    make([]Message, 0),
		LastCommand:  "",
		Metadata:     make(map[string]string),
	}
}

// Get returns the global session instance
func Get() *Session {
	if globalSession == nil {
		globalSession = New()
	}
	return globalSession
}

// ProjectName returns the current project name from cwd
func (s *Session) ProjectName() string {
	return filepath.Base(s.CurrentDir)
}

// GitBranch returns the current git branch (cached for performance)
func (s *Session) GitBranch() string {
	// Check cache first
	if branch, ok := s.Metadata["git_branch"]; ok {
		// Cache for 5 seconds
		if ts, ok := s.Metadata["git_branch_ts"]; ok {
			if cached, err := time.Parse(time.RFC3339, ts); err == nil {
				if time.Since(cached) < 5*time.Second {
					return branch
				}
			}
		}
	}

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = s.CurrentDir
	out, err := cmd.Output()
	if err != nil {
		s.Metadata["git_branch"] = ""
		return ""
	}

	branch := strings.TrimSpace(string(out))
	s.Metadata["git_branch"] = branch
	s.Metadata["git_branch_ts"] = time.Now().Format(time.RFC3339)
	return branch
}

// IsGitRepo checks if current directory is a git repository
func (s *Session) IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = s.CurrentDir
	return cmd.Run() == nil
}

// IncrementCommandCount increments the command counter
func (s *Session) IncrementCommandCount() {
	s.CommandCount++
}

// SetLastCommand records the last executed command
func (s *Session) SetLastCommand(cmd string) {
	s.LastCommand = cmd
}

// UpdateCurrentDir updates the current working directory
func (s *Session) UpdateCurrentDir() {
	cwd, err := os.Getwd()
	if err == nil {
		s.CurrentDir = cwd
		// Invalidate git cache
		delete(s.Metadata, "git_branch")
		delete(s.Metadata, "git_branch_ts")
	}
}

// AddAIMessage adds a message to AI context
func (s *Session) AddAIMessage(role, content string) {
	s.AIContext = append(s.AIContext, Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	})

	// Keep only last 20 messages for context window management
	if len(s.AIContext) > 20 {
		s.AIContext = s.AIContext[len(s.AIContext)-20:]
	}
}

// GetAIContext returns the AI conversation context
func (s *Session) GetAIContext() []Message {
	return s.AIContext
}

// ClearAIContext clears the AI conversation history
func (s *Session) ClearAIContext() {
	s.AIContext = make([]Message, 0)
}

// AIContextSize returns the number of messages in AI context
func (s *Session) AIContextSize() int {
	return len(s.AIContext)
}

// Uptime returns session duration
func (s *Session) Uptime() time.Duration {
	return time.Since(s.StartTime)
}

// UptimeFormatted returns a human-readable uptime
func (s *Session) UptimeFormatted() string {
	d := s.Uptime()

	if d < time.Minute {
		return d.Round(time.Second).String()
	}
	if d < time.Hour {
		return d.Round(time.Minute).String()
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 && minutes > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dm", minutes)
}

// Save persists session to disk
func (s *Session) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Load restores session from disk
func Load(path string) (*Session, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return New(), nil // File doesn't exist, return new session
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return New(), nil // Invalid JSON, return new session
	}

	// Update current directory
	s.UpdateCurrentDir()

	return &s, nil
}

// Stats returns session statistics
func (s *Session) Stats() map[string]interface{} {
	return map[string]interface{}{
		"uptime":      s.UptimeFormatted(),
		"commands":    s.CommandCount,
		"ai_messages": len(s.AIContext),
		"project":     s.ProjectName(),
		"git_branch":  s.GitBranch(),
		"is_git_repo": s.IsGitRepo(),
	}
}
