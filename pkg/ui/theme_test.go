package ui_test

import (
	"testing"

	"github.com/badie/bdev/pkg/ui"
)

func TestTheme(t *testing.T) {
	// Test default theme
	if ui.CurrentTheme == nil {
		t.Error("CurrentTheme should not be nil")
	}

	// Test colors
	if ui.Primary("test") == "" {
		t.Error("Primary color output should not be empty")
	}
	if ui.Secondary("test") == "" {
		t.Error("Secondary color output should not be empty")
	}
	if ui.Success("test") == "" {
		t.Error("Success color output should not be empty")
	}
	if ui.Warning("test") == "" {
		t.Error("Warning color output should not be empty")
	}
	if ui.Error("test") == "" {
		t.Error("Error color output should not be empty")
	}
	if ui.Muted("test") == "" {
		t.Error("Muted color output should not be empty")
	}
	if ui.Info("test") == "" {
		t.Error("Info color output should not be empty")
	}
	if ui.Cyan("test") == "" {
		t.Error("Cyan color output should not be empty")
	}
	if ui.Bold("test") == "" {
		t.Error("Bold color output should not be empty")
	}
}

func TestGlyphs(t *testing.T) {
	// Test default glyphs are populated
	g := ui.ActiveGlyphs

	if g.Check == "" {
		t.Error("Check glyph should not be empty")
	}
	if g.Cross == "" {
		t.Error("Cross glyph should not be empty")
	}

	// Test fallback glyphs
	g = ui.FallbackGlyphs
	if g.Check == "" {
		t.Error("Fallback Check glyph should not be empty")
	}
}

func TestAnthropicTheme(t *testing.T) {
	theme := ui.AnthropicTheme()
	if theme == nil {
		t.Error("AnthropicTheme should not return nil")
	}

	if theme.Rust == nil {
		t.Error("Theme functions should not be nil")
	}
}
