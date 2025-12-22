package ui

import (
	"fmt"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════════════════════════════
// B.DEV OMEGA DESIGN SYSTEM v2.0 - COMPONENTS
// Architecture: Claude Code Anthropic (Terminal-First)
// ═══════════════════════════════════════════════════════════════════════════════

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 1. PROGRESS BAR (9 Phases per User Spec)                                    │
// └─────────────────────────────────────────────────────────────────────────────┘

// ProgressFrame defines a single frame of a progress bar animation
type ProgressFrame struct {
	Bar         string
	Percentage  int
	Description string
}

// ProgressBuild defines the standard build progress animation frames
var ProgressBuild = []ProgressFrame{
	{Bar: "[░░░░░░░░░░░░░░░░░░░░]", Percentage: 0, Description: "Initializing"},
	{Bar: "[██░░░░░░░░░░░░░░░░░░]", Percentage: 12, Description: "Parsing modules"},
	{Bar: "[████░░░░░░░░░░░░░░░░]", Percentage: 24, Description: "Type checking"},
	{Bar: "[███████░░░░░░░░░░░░░]", Percentage: 38, Description: "Compiling packages"},
	{Bar: "[██████████░░░░░░░░░░]", Percentage: 51, Description: "Linking dependencies"},
	{Bar: "[████████████░░░░░░░░]", Percentage: 63, Description: "Optimizing output"},
	{Bar: "[███████████████░░░░░]", Percentage: 77, Description: "Stripping symbols"},
	{Bar: "[█████████████████░░░]", Percentage: 89, Description: "Finalizing build"},
	{Bar: "[████████████████████]", Percentage: 100, Description: "Complete " + ActiveGlyphs.Check},
}

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 2. STATUS INDICATORS (Spinners per User Spec)                               │
// └─────────────────────────────────────────────────────────────────────────────┘

// StatusFrames defines a sequence of strings for status animation
type StatusFrames []string

var (
	// StatusLoading is the default spinner animation
	StatusLoading    = StatusFrames(ActiveGlyphs.SpinnerDefault)
	// StatusProcessing is the build spinner animation
	StatusProcessing = StatusFrames(ActiveGlyphs.SpinnerBuild)
	// StatusNetwork is the network spinner animation
	StatusNetwork    = StatusFrames(ActiveGlyphs.SpinnerNetwork)
)

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 3. SEPARATORS (User Spec)                                                   │
// └─────────────────────────────────────────────────────────────────────────────┘

const (
	// SeparatorLight is a light separator line
	SeparatorLight  = `─────────────────────────────────────────────────────────────`
	// SeparatorMedium is a medium (double) separator line
	SeparatorMedium = `═════════════════════════════════════════════════════════════`
	// SeparatorDotted is a dotted separator line
	SeparatorDotted = `·························································`
)

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 4. BOX COMPONENT (User Spec - 3 Styles)                                     │
// └─────────────────────────────────────────────────────────────────────────────┘

// BoxStyle defines the visual style of a box component
type BoxStyle string

const (
	// BoxStyleLight uses light box drawing characters
	BoxStyleLight   BoxStyle = "light"
	// BoxStyleRounded uses rounded corner box drawing characters
	BoxStyleRounded BoxStyle = "rounded"
	// BoxStyleDouble uses double line box drawing characters
	BoxStyleDouble  BoxStyle = "double"
)

// Box represents a boxed content area
type Box struct {
	Width   int
	Title   string
	Content string
	Style   BoxStyle
}

// Render returns the string representation of the box
func (b *Box) Render() string {
	var tl, tr, bl, br, h, v string

	switch b.Style {
	case BoxStyleRounded:
		tl, tr = ActiveGlyphs.RoundTopLeft, ActiveGlyphs.RoundTopRight
		bl, br = ActiveGlyphs.RoundBottomLeft, ActiveGlyphs.RoundBottomRight
		h, v = ActiveGlyphs.BoxHorizontal, ActiveGlyphs.BoxVertical
	case BoxStyleDouble:
		tl, tr = ActiveGlyphs.DoubleTopLeft, ActiveGlyphs.DoubleTopRight
		bl, br = ActiveGlyphs.DoubleBottomLeft, ActiveGlyphs.DoubleBottomRight
		h, v = ActiveGlyphs.DoubleHorizontal, ActiveGlyphs.DoubleVertical
	default: // light
		tl, tr = ActiveGlyphs.BoxTopLeft, ActiveGlyphs.BoxTopRight
		bl, br = ActiveGlyphs.BoxBottomLeft, ActiveGlyphs.BoxBottomRight
		h, v = ActiveGlyphs.BoxHorizontal, ActiveGlyphs.BoxVertical
	}

	// Header
	titleLen := len(b.Title)
	headerWidth := b.Width - titleLen - 4
	if headerWidth < 0 {
		headerWidth = 0
	}
	header := fmt.Sprintf("%s%s %s %s%s", tl, h, b.Title, strings.Repeat(h, headerWidth), tr)

	// Content
	contentLines := strings.Split(b.Content, "\n")
	body := ""
	for _, line := range contentLines {
		padding := b.Width - len(line) - 4
		if padding < 0 {
			padding = 0
		}
		body += fmt.Sprintf("%s  %s%s  %s\n", v, line, strings.Repeat(" ", padding), v)
	}

	// Footer
	footer := fmt.Sprintf("%s%s%s", bl, strings.Repeat(h, b.Width-2), br)

	return header + "\n" + body + footer
}

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 5. PROMPT STYLES (User Spec)                                                │
// └─────────────────────────────────────────────────────────────────────────────┘

var (
	// PromptSpine is the vertical line used in prompts
	PromptSpine   = "┃"
	// PromptPointer is the arrow used in prompts
	PromptPointer = ActiveGlyphs.Pointer
)

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 6. STATUS MESSAGES (User Spec)                                              │
// └─────────────────────────────────────────────────────────────────────────────┘

// MessageSuccess formats a success message with an icon
func MessageSuccess(msg string) string {
	return fmt.Sprintf("%s  %s", Success(ActiveGlyphs.Check), msg)
}

// MessageError formats an error message with an icon
func MessageError(msg string) string {
	return fmt.Sprintf("%s  %s", Error(ActiveGlyphs.Cross), msg)
}

// MessageWarning formats a warning message with an icon
func MessageWarning(msg string) string {
	return fmt.Sprintf("%s  %s", Warning(ActiveGlyphs.Warning), msg)
}

// MessageInfo formats an info message with an icon
func MessageInfo(msg string) string {
	return fmt.Sprintf("%s  %s", Info(ActiveGlyphs.Info), msg)
}

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 7. WELCOME HEADER (User Spec - Double Box)                                  │
// └─────────────────────────────────────────────────────────────────────────────┘

// PrintWelcome prints the standard welcome banner
func PrintWelcome() {
	width := 58

	tl, tr := ActiveGlyphs.DoubleTopLeft, ActiveGlyphs.DoubleTopRight
	bl, br := ActiveGlyphs.DoubleBottomLeft, ActiveGlyphs.DoubleBottomRight
	h, v := ActiveGlyphs.DoubleHorizontal, ActiveGlyphs.DoubleVertical

	topBorder := tl + strings.Repeat(h, width) + tr
	emptyLine := v + strings.Repeat(" ", width) + v
	botBorder := bl + strings.Repeat(h, width) + br

	title := "B·DEV OMEGA"
	subtitle := "Elite Development Environment"
	version := "Version 2.5.0"
	status := fmt.Sprintf("System Ready %s %s", ActiveGlyphs.Online, time.Now().Format("15:04:05"))

	centerPad := func(s string) string {
		pad := (width - len(s)) / 2
		return strings.Repeat(" ", pad) + s + strings.Repeat(" ", width-len(s)-pad)
	}

	fmt.Println()
	fmt.Println(Primary(topBorder))
	fmt.Println(Primary(emptyLine))
	fmt.Println(Primary(v) + Bold(centerPad(title)) + Primary(v))
	fmt.Println(Primary(v) + Muted(centerPad(subtitle)) + Primary(v))
	fmt.Println(Primary(emptyLine))
	fmt.Println(Primary(v) + Muted(centerPad(version)) + Primary(v))
	fmt.Println(Primary(v) + Success(centerPad(status)) + Primary(v))
	fmt.Println(Primary(emptyLine))
	fmt.Println(Primary(botBorder))
	fmt.Println()
}

// PrintWelcomeWithVersion prints the welcome banner (version arguments are deprecated/unused in current implementation)
func PrintWelcomeWithVersion(version, goVersion string) {
	PrintWelcome()
}

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 8. UTILITY FUNCTIONS                                                        │
// └─────────────────────────────────────────────────────────────────────────────┘

// PrintHeader prints a section header with a separator
func PrintHeader(header string) {
	fmt.Println()
	fmt.Println(Bold(header))
	fmt.Println(Muted(SeparatorMedium))
}

// PrintProgressBar prints a single frame of the progress bar
func PrintProgressBar(frame ProgressFrame) {
	fmt.Printf("\r%s %3d%%  %s",
		Primary(frame.Bar),
		frame.Percentage,
		Muted(frame.Description))
}

// AnimateProgress runs the progress bar animation
func AnimateProgress(frames []ProgressFrame, delayMs int) {
	for _, frame := range frames {
		PrintProgressBar(frame)
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
	fmt.Println()
}

// AnimateStatus runs a status spinner animation for a fixed duration
func AnimateStatus(frames StatusFrames, description string, duration time.Duration) {
	start := time.Now()
	frameIndex := 0

	for time.Since(start) < duration {
		fmt.Printf("\r%s %s",
			Primary(frames[frameIndex%len(frames)]),
			Muted(description))
		frameIndex++
		time.Sleep(80 * time.Millisecond) // 80ms per user spec
	}
	fmt.Println()
}

// ClearLine clears the current terminal line
func ClearLine() {
	fmt.Print("\r\033[K")
}

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 9. SPINNER COMPONENT (User Spec - 80ms)                                     │
// └─────────────────────────────────────────────────────────────────────────────┘

// Spinner represents a loading spinner
type Spinner struct {
	Frames      StatusFrames
	Description string
	Delay       time.Duration
	stopCh      chan bool
	active      bool
}

// NewSpinner creates a new default spinner
func NewSpinner(description string) *Spinner {
	return &Spinner{
		Frames:      StatusLoading,
		Description: description,
		Delay:       80 * time.Millisecond, // User spec: 80ms
		stopCh:      make(chan bool),
	}
}

// NewSpinnerWithStyle creates a new spinner with a specific style (build, network, or default)
func NewSpinnerWithStyle(description, style string) *Spinner {
	var frames StatusFrames
	switch style {
	case "build":
		frames = StatusProcessing
	case "network":
		frames = StatusNetwork
	default:
		frames = StatusLoading
	}
	return &Spinner{
		Frames:      frames,
		Description: description,
		Delay:       80 * time.Millisecond,
		stopCh:      make(chan bool),
	}
}

// Start begins the spinner animation in a goroutine
func (s *Spinner) Start() {
	if s.active {
		return
	}
	s.active = true
	go func() {
		frameIndex := 0
		for {
			select {
			case <-s.stopCh:
				return
			default:
				fmt.Printf("\r%s %s",
					Primary(s.Frames[frameIndex%len(s.Frames)]),
					Muted(s.Description))
				frameIndex++
				time.Sleep(s.Delay)
			}
		}
	}()
}

// Stop halts the spinner animation
func (s *Spinner) Stop() {
	if !s.active {
		return
	}
	s.active = false
	s.stopCh <- true
	ClearLine()
}

// Success stops the spinner and prints a success message
func (s *Spinner) Success(msg string) {
	s.Stop()
	fmt.Println(MessageSuccess(msg))
}

// Error stops the spinner and prints an error message
func (s *Spinner) Error(msg string) {
	s.Stop()
	fmt.Println(MessageError(msg))
}

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 10. BOOT SEQUENCE (User Spec - 600ms)                                       │
// └─────────────────────────────────────────────────────────────────────────────┘

// AnimateBoot runs the system boot animation
func AnimateBoot() {
	AnimateStatus(StatusProcessing, "Initializing system core...", 600*time.Millisecond)
	ClearLine()
	fmt.Printf("\r%s %s\n", Success(ActiveGlyphs.Check), Primary("System Online"))
	time.Sleep(100 * time.Millisecond)
}

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 11. DEMO (Full Component Showcase)                                          │
// └─────────────────────────────────────────────────────────────────────────────┘

// Demo runs a demonstration of all UI components
func Demo() {
	PrintHeader("OMEGA DESIGN SYSTEM v2.0 :: DEMO")

	fmt.Println(Bold("1. Welcome Header"))
	PrintWelcome()

	fmt.Println(Bold("2. Progress Bar (9 phases)"))
	AnimateProgress(ProgressBuild, 150)
	fmt.Println()

	fmt.Println(Bold("3. Spinners"))
	s := NewSpinner("Processing request...")
	s.Start()
	time.Sleep(800 * time.Millisecond)
	s.Success("Request completed")
	fmt.Println()

	fmt.Println(Bold("4. Status Messages"))
	fmt.Println(MessageSuccess("Operation completed successfully"))
	fmt.Println(MessageWarning("Disk space low"))
	fmt.Println(MessageError("Connection failed"))
	fmt.Println(MessageInfo("Update available"))
	fmt.Println()

	fmt.Println(Bold("5. Box Components"))
	box := Box{
		Width:   40,
		Title:   "System Status",
		Content: fmt.Sprintf("%s Services      Running\n%s Database      Connected\n%s Cache         Operational", ActiveGlyphs.Online, ActiveGlyphs.Pointer, ActiveGlyphs.Pointer),
		Style:   BoxStyleLight,
	}
	fmt.Println(box.Render())
	fmt.Println()

	roundedBox := Box{
		Width:   40,
		Title:   "Fatal Error",
		Content: fmt.Sprintf("%s Database Connection Failed\n\nHost: localhost:5432\nError: ECONNREFUSED", ActiveGlyphs.Cross),
		Style:   BoxStyleRounded,
	}
	fmt.Println(Error(roundedBox.Render()))
	fmt.Println()
}
