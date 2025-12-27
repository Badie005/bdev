package agents_test

import (
	"testing"

	"github.com/badie/bdev/internal/ai/agents"
)

func TestGetAgent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
	}{
		{"Reviewer", "reviewer", "reviewer"},
		{"ReviewerCase", "Reviewer", "reviewer"},
		{"Explainer", "explainer", "explainer"},
		{"Documenter", "documenter", "documenter"},
		{"Debugger", "debugger", "debugger"},
		{"Architect", "architect", "architect"},
		{"Unknown", "unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := agents.GetAgent(tt.input)
			if tt.wantName == "" {
				if got != nil {
					t.Errorf("GetAgent(%q) = %v, want nil", tt.input, got)
				}
			} else {
				if got == nil {
					t.Fatalf("GetAgent(%q) returned nil", tt.input)
				}
				if got.Name() != tt.wantName {
					t.Errorf("GetAgent(%q).Name() = %q, want %q", tt.input, got.Name(), tt.wantName)
				}
			}
		})
	}
}

func TestAgents(t *testing.T) {
	// Verify all agents have content
	for _, agent := range agents.AllAgents {
		if agent.Name() == "" {
			t.Errorf("Agent has empty name")
		}
		if agent.Description() == "" {
			t.Errorf("Agent %s has empty description", agent.Name())
		}
		if agent.SystemPrompt() == "" {
			t.Errorf("Agent %s has empty prompt", agent.Name())
		}
	}
}
