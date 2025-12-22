package ui

import (
	"github.com/fatih/color"
)

// ═══════════════════════════════════════════════════════════════════════════════
// B.DEV OMEGA DESIGN SYSTEM v2.0
// Architecture: Claude Code Anthropic (Terminal-First)
// ═══════════════════════════════════════════════════════════════════════════════

// Theme defines the contract for a visual theme
type Theme struct {
	// Primary Accent
	Rust      func(a ...interface{}) string
	RustDark  func(a ...interface{}) string
	RustLight func(a ...interface{}) string

	// Backgrounds
	Graphite      func(a ...interface{}) string
	GraphiteLight func(a ...interface{}) string

	// Text
	Snow       func(a ...interface{}) string
	Slate      func(a ...interface{}) string
	SlateLight func(a ...interface{}) string

	// Status
	Success func(a ...interface{}) string
	Warning func(a ...interface{}) string
	Error   func(a ...interface{}) string
	Info    func(a ...interface{}) string
}

// CurrentTheme is the global Theme Instance
var CurrentTheme = AnthropicTheme()

// AnthropicTheme returns the Anthropic Rust/Warm palette
// Rust #E07B39 approximated to HiYellow (closest ANSI match for orange)
func AnthropicTheme() *Theme {
	return &Theme{
		// Rust: #E07B39 - Anthropic Orange (approximated)
		Rust:      color.New(color.FgHiYellow, color.Bold).SprintFunc(),
		RustDark:  color.New(color.FgYellow).SprintFunc(),
		RustLight: color.New(color.FgHiYellow).SprintFunc(),

		// Graphite: Dark backgrounds
		Graphite:      color.New(color.BgBlack).SprintFunc(),
		GraphiteLight: color.New(color.BgHiBlack).SprintFunc(),

		// Snow: White text
		Snow:       color.New(color.FgWhite, color.Bold).SprintFunc(),
		Slate:      color.New(color.FgHiBlack).SprintFunc(),
		SlateLight: color.New(color.FgWhite).SprintFunc(),

		// Status colors
		Success: color.New(color.FgGreen).SprintFunc(),
		Warning: color.New(color.FgYellow).SprintFunc(),
		Error:   color.New(color.FgRed).SprintFunc(),
		Info:    color.New(color.FgBlue).SprintFunc(),
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// LEGACY PROXIES (Backward Compatibility)
// ═══════════════════════════════════════════════════════════════════════════════

// Primary returns the primary accent color
func Primary(a ...interface{}) string { return CurrentTheme.Rust(a...) }

// Secondary returns the secondary accent color
func Secondary(a ...interface{}) string { return CurrentTheme.Slate(a...) }

// Success returns the success color
func Success(a ...interface{}) string { return CurrentTheme.Success(a...) }

// Warning returns the warning color
func Warning(a ...interface{}) string { return CurrentTheme.Warning(a...) }

// Error returns the error color
func Error(a ...interface{}) string { return CurrentTheme.Error(a...) }

// Muted returns the muted color
func Muted(a ...interface{}) string { return CurrentTheme.Slate(a...) }

// Info returns the info color
func Info(a ...interface{}) string { return CurrentTheme.Info(a...) }

// Cyan returns the cyan color
func Cyan(a ...interface{}) string { return CurrentTheme.Rust(a...) }

// Bold returns the bold color
func Bold(a ...interface{}) string { return CurrentTheme.Snow(a...) }

// ═══════════════════════════════════════════════════════════════════════════════
// GLYPH REGISTRY (User Spec v2.0)
// ═══════════════════════════════════════════════════════════════════════════════

// Glyphs defines the contract for visual glyphs
type Glyphs struct {
	// Spinners (80ms interval)
	SpinnerDefault []string
	SpinnerBuild   []string
	SpinnerNetwork []string

	// Status Icons
	Check   string
	Cross   string
	Warning string
	Info    string
	Pointer string
	Bullet  string
	Online  string
	Offline string
	Folder  string
	File    string
	Branch  string
	Lock    string
	Sparkle string

	// Box Drawing (Light)
	BoxTopLeft     string
	BoxTopRight    string
	BoxBottomLeft  string
	BoxBottomRight string
	BoxHorizontal  string
	BoxVertical    string

	// Box Drawing (Rounded)
	RoundTopLeft     string
	RoundTopRight    string
	RoundBottomLeft  string
	RoundBottomRight string

	// Box Drawing (Double)
	DoubleTopLeft     string
	DoubleTopRight    string
	DoubleBottomLeft  string
	DoubleBottomRight string
	DoubleHorizontal  string
	DoubleVertical    string
}

// ActiveGlyphs - User Design System v2.0 spec
var ActiveGlyphs = Glyphs{
	// Spinners
	SpinnerDefault: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	SpinnerBuild:   []string{"▱▱▱", "▰▱▱", "▰▰▱", "▰▰▰", "▱▰▰", "▱▱▰"},
	SpinnerNetwork: []string{"→", "↗", "↑", "↖", "←", "↙", "↓", "↘"},

	// Status Icons (User Spec)
	Check:   "✓",
	Cross:   "✗",
	Warning: "!",
	Info:    "i",
	Pointer: "→",
	Bullet:  "•",
	Online:  "●",
	Offline: "○",
	Folder:  "▸",
	File:    "—",
	Branch:  "⎇",
	Lock:    "🔒",
	Sparkle: "✓",

	// Light Box
	BoxTopLeft:     "┌",
	BoxTopRight:    "┐",
	BoxBottomLeft:  "└",
	BoxBottomRight: "┘",
	BoxHorizontal:  "─",
	BoxVertical:    "│",

	// Rounded Box
	RoundTopLeft:     "╭",
	RoundTopRight:    "╮",
	RoundBottomLeft:  "╰",
	RoundBottomRight: "╯",

	// Double Box
	DoubleTopLeft:     "╔",
	DoubleTopRight:    "╗",
	DoubleBottomLeft:  "╚",
	DoubleBottomRight: "╝",
	DoubleHorizontal:  "═",
	DoubleVertical:    "║",
}

// FallbackGlyphs for basic ASCII terminals
var FallbackGlyphs = Glyphs{
	SpinnerDefault: []string{"|", "/", "-", "\\"},
	SpinnerBuild:   []string{"...", "..-", ".-.", "--.", "---"},
	SpinnerNetwork: []string{"->", ">>", "<-", "<<"},

	Check:   "[OK]",
	Cross:   "[X]",
	Warning: "[!]",
	Info:    "[i]",
	Pointer: "->",
	Bullet:  "-",
	Online:  "[*]",
	Offline: "[ ]",
	Folder:  ">",
	File:    "-",
	Branch:  "[b]",
	Lock:    "[L]",
	Sparkle: "*",

	BoxTopLeft:     "+",
	BoxTopRight:    "+",
	BoxBottomLeft:  "+",
	BoxBottomRight: "+",
	BoxHorizontal:  "-",
	BoxVertical:    "|",

	RoundTopLeft:     "+",
	RoundTopRight:    "+",
	RoundBottomLeft:  "+",
	RoundBottomRight: "+",

	DoubleTopLeft:     "+",
	DoubleTopRight:    "+",
	DoubleBottomLeft:  "+",
	DoubleBottomRight: "+",
	DoubleHorizontal:  "=",
	DoubleVertical:    "|",
}
