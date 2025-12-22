package ui_test

import (
	"strings"
	"testing"
	"time"

	"github.com/badie/bdev/pkg/ui"
)

func TestStatusMessages(t *testing.T) {
	msg := "test message"

	if !strings.Contains(ui.MessageSuccess(msg), msg) {
		t.Error("MessageSuccess should contain the message")
	}
	if !strings.Contains(ui.MessageError(msg), msg) {
		t.Error("MessageError should contain the message")
	}
	if !strings.Contains(ui.MessageWarning(msg), msg) {
		t.Error("MessageWarning should contain the message")
	}
	if !strings.Contains(ui.MessageInfo(msg), msg) {
		t.Error("MessageInfo should contain the message")
	}
}

func TestBoxRender(t *testing.T) {
	box := ui.Box{
		Width:   20,
		Title:   "Test",
		Content: "Hello\nWorld",
		Style:   ui.BoxStyleRounded,
	}

	rendered := box.Render()
	if !strings.Contains(rendered, "Test") {
		t.Error("Box should contain title")
	}
	if !strings.Contains(rendered, "Hello") {
		t.Error("Box should contain content")
	}
	if !strings.Contains(rendered, ui.ActiveGlyphs.RoundTopLeft) {
		t.Error("Box should use requested style glyphs")
	}
}

func TestSpinner(t *testing.T) {
	s := ui.NewSpinner("Test")
	if s == nil {
		t.Error("NewSpinner should return instance")
	}

	// Start and stop quickly
	s.Start()
	time.Sleep(10 * time.Millisecond)
	s.Stop()

	// Test helpers
	s.Success("Done")
	s.Error("Failed")
}
