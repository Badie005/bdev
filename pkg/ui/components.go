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

type ProgressFrame struct {
	Bar         string
	Percentage  int
	Description string
}

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

type StatusFrames []string

var (
	StatusLoading    = StatusFrames(ActiveGlyphs.SpinnerDefault)
	StatusProcessing = StatusFrames(ActiveGlyphs.SpinnerBuild)
	StatusNetwork    = StatusFrames(ActiveGlyphs.SpinnerNetwork)
)

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 3. SEPARATORS (User Spec)                                                   │
// └─────────────────────────────────────────────────────────────────────────────┘

const (
	SeparatorLight  = `─────────────────────────────────────────────────────────────`
	SeparatorMedium = `═════════════════════════════════════════════════════════════`
	SeparatorDotted = `·························································`
)

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 4. BOX COMPONENT (User Spec - 3 Styles)                                     │
// └─────────────────────────────────────────────────────────────────────────────┘

type BoxStyle string

const (
	BoxStyleLight   BoxStyle = "light"
	BoxStyleRounded BoxStyle = "rounded"
	BoxStyleDouble  BoxStyle = "double"
)

type Box struct {
	Width   int
	Title   string
	Content string
	Style   BoxStyle
}

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
	PromptSpine   = "┃"
	PromptPointer = ActiveGlyphs.Pointer
)

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 6. STATUS MESSAGES (User Spec)                                              │
// └─────────────────────────────────────────────────────────────────────────────┘

func MessageSuccess(msg string) string {
	return fmt.Sprintf("%s  %s", Success(ActiveGlyphs.Check), msg)
}

func MessageError(msg string) string {
	return fmt.Sprintf("%s  %s", Error(ActiveGlyphs.Cross), msg)
}

func MessageWarning(msg string) string {
	return fmt.Sprintf("%s  %s", Warning(ActiveGlyphs.Warning), msg)
}

func MessageInfo(msg string) string {
	return fmt.Sprintf("%s  %s", Info(ActiveGlyphs.Info), msg)
}

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 7. WELCOME HEADER (User Spec - Double Box)                                  │
// └─────────────────────────────────────────────────────────────────────────────┘

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

func PrintWelcomeWithVersion(version, goVersion string) {
	PrintWelcome()
}

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 8. UTILITY FUNCTIONS                                                        │
// └─────────────────────────────────────────────────────────────────────────────┘

func PrintHeader(header string) {
	fmt.Println()
	fmt.Println(Bold(header))
	fmt.Println(Muted(SeparatorMedium))
}

func PrintProgressBar(frame ProgressFrame) {
	fmt.Printf("\r%s %3d%%  %s",
		Primary(frame.Bar),
		frame.Percentage,
		Muted(frame.Description))
}

func AnimateProgress(frames []ProgressFrame, delayMs int) {
	for _, frame := range frames {
		PrintProgressBar(frame)
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
	fmt.Println()
}

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

func ClearLine() {
	fmt.Print("\r\033[K")
}

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 9. SPINNER COMPONENT (User Spec - 80ms)                                     │
// └─────────────────────────────────────────────────────────────────────────────┘

type Spinner struct {
	Frames      StatusFrames
	Description string
	Delay       time.Duration
	stopCh      chan bool
	active      bool
}

func NewSpinner(description string) *Spinner {
	return &Spinner{
		Frames:      StatusLoading,
		Description: description,
		Delay:       80 * time.Millisecond, // User spec: 80ms
		stopCh:      make(chan bool),
	}
}

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

func (s *Spinner) Stop() {
	if !s.active {
		return
	}
	s.active = false
	s.stopCh <- true
	ClearLine()
}

func (s *Spinner) Success(msg string) {
	s.Stop()
	fmt.Println(MessageSuccess(msg))
}

func (s *Spinner) Error(msg string) {
	s.Stop()
	fmt.Println(MessageError(msg))
}

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 10. BOOT SEQUENCE (User Spec - 600ms)                                       │
// └─────────────────────────────────────────────────────────────────────────────┘

func AnimateBoot() {
	AnimateStatus(StatusProcessing, "Initializing system core...", 600*time.Millisecond)
	ClearLine()
	fmt.Printf("\r%s %s\n", Success(ActiveGlyphs.Check), Primary("System Online"))
	time.Sleep(100 * time.Millisecond)
}

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 11. DEMO (Full Component Showcase)                                          │
// └─────────────────────────────────────────────────────────────────────────────┘

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
