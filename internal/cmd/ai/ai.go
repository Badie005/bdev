package aicmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/cobra"

	"github.com/badie/bdev/internal/ai/engine"
	"github.com/badie/bdev/internal/ai/memory"
	"github.com/badie/bdev/internal/core/config"
	"github.com/badie/bdev/pkg/ui"
)

var client *engine.Client

func getClient() *engine.Client {
	if client == nil {
		cfg := config.Get()
		client = engine.New(engine.Config{
			BaseURL:     cfg.AI.BaseURL,
			Model:       cfg.AI.Model,
			Fallback:    cfg.AI.FallbackModel,
			Timeout:     120 * 1e9, // 120 seconds
			MaxRetries:  3,
			Temperature: 0.7,
		})
	}
	return client
}

func getMemory() *memory.Memory {
	cfg := config.Get()
	return memory.GetGlobal(cfg.AIMemoryFile())
}

// NewCommand creates the AI command group
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI assistant powered by Ollama",
		Long:  "Chat with AI, check status, and manage conversation memory",
	}

	cmd.AddCommand(checkCmd())
	cmd.AddCommand(chatCmd())
	cmd.AddCommand(forgetCmd())
	cmd.AddCommand(modelsCmd())
	cmd.AddCommand(memoryCmd())

	return cmd
}

// ============================================================
// CHECK - Check Ollama status
// ============================================================

func checkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check if Ollama is running",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()

			if !c.IsAvailable() {
				fmt.Println(ui.Error("Ollama is not running"))
				fmt.Println(ui.Muted("Start it with: ollama serve"))
				return nil
			}

			fmt.Println(ui.Success("Ollama is running"))
			fmt.Printf("  %s %s\n", ui.Primary("Model:"), c.GetCurrentModel())

			models, err := c.ListModels()
			if err == nil && len(models) > 0 {
				fmt.Printf("  %s %d available\n", ui.Primary("Models:"), len(models))
			}

			return nil
		},
	}
}

// ============================================================
// CHAT - Interactive chat
// ============================================================

func chatCmd() *cobra.Command {
	var noMemory bool

	cmd := &cobra.Command{
		Use:   "chat [prompt]",
		Short: "Chat with the AI assistant",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			m := getMemory()

			if !c.IsAvailable() {
				return fmt.Errorf("ollama is not running. Start it with: ollama serve")
			}

			prompt := strings.Join(args, " ")

			// Add user message to memory
			if !noMemory {
				m.Add("user", prompt)
			}

			// Prepare messages
			var messages []engine.Message
			if noMemory {
				messages = []engine.Message{{Role: "user", Content: prompt}}
			} else {
				messages = m.GetContextWithSystem(getSystemPrompt())
			}

			// Create context with cancellation
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Handle Ctrl+C
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)
			go func() {
				<-sigChan
				cancel()
			}()

			// Stream response
			fmt.Println()
			output, errChan := c.Chat(ctx, messages)

			var fullResponse strings.Builder
			for chunk := range output {
				fmt.Print(chunk)
				fullResponse.WriteString(chunk)
			}
			fmt.Println()
			fmt.Println()

			// Check for errors
			select {
			case err := <-errChan:
				if err != nil {
					return err
				}
			default:
			}

			// Save response to memory
			if !noMemory && fullResponse.Len() > 0 {
				m.Add("assistant", fullResponse.String())
				_ = m.Save() // Non-critical: best effort persistence
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&noMemory, "no-memory", false, "Don't use conversation memory")
	return cmd
}

// ============================================================
// FORGET - Clear memory
// ============================================================

func forgetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "forget",
		Short: "Clear conversation memory",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := getMemory()
			count := m.Size()
			m.Clear()
			_ = m.Save() // Non-critical: best effort persistence

			fmt.Println(ui.Success(fmt.Sprintf("Cleared %d messages from memory", count)))
			return nil
		},
	}
}

// ============================================================
// MODELS - List available models
// ============================================================

func modelsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "models",
		Short: "List available Ollama models",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()

			if !c.IsAvailable() {
				return fmt.Errorf("ollama is not running")
			}

			models, err := c.ListModels()
			if err != nil {
				return err
			}

			if len(models) == 0 {
				fmt.Println(ui.Muted("No models installed"))
				fmt.Println(ui.Muted("Pull a model with: ollama pull llama3.2"))
				return nil
			}

			fmt.Println(ui.Bold("Available Models:"))
			for _, m := range models {
				size := float64(m.Size) / 1024 / 1024 / 1024
				current := ""
				if m.Name == c.Model || m.Name == c.GetCurrentModel() {
					current = ui.Success(" (active)")
				}
				fmt.Printf("  %s%s %s\n",
					ui.Primary(m.Name),
					current,
					ui.Muted(fmt.Sprintf("(%.1f GB)", size)))
			}

			return nil
		},
	}
}

// ============================================================
// MEMORY - Show memory status
// ============================================================

func memoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "memory",
		Short: "Show conversation memory status",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := getMemory()

			fmt.Println(ui.Bold("Memory Status:"))
			fmt.Printf("  %s %d\n", ui.Primary("Messages:"), m.Size())
			fmt.Printf("  %s ~%d\n", ui.Primary("Tokens:"), m.TotalTokens())

			if m.Size() > 0 {
				fmt.Println()
				fmt.Println(ui.Bold("Recent conversation:"))

				// Show last 5 messages
				start := 0
				if m.Size() > 5 {
					start = m.Size() - 5
				}

				for i := start; i < m.Size(); i++ {
					msg := m.Messages[i]
					var role string
					switch msg.Role {
					case "user":
						role = ui.Primary("you:")
					case "assistant":
						role = ui.Success("ai:")
					default:
						role = ui.Muted(msg.Role + ":")
					}

					content := msg.Content
					if len(content) > 80 {
						content = content[:80] + "..."
					}
					fmt.Printf("  %s %s\n", role, content)
				}
			}

			return nil
		},
	}
}

// ============================================================
// HELPERS
// ============================================================

func getSystemPrompt() string {
	return `You are B.DEV AI, a helpful coding assistant integrated into a CLI tool.
You help developers with:
- Code review and suggestions
- Debugging and problem solving
- Architecture and design decisions
- Documentation and explanations
- Best practices and patterns

Be concise and practical. Provide code examples when helpful.
Format responses for terminal display (no markdown images or complex tables).`
}
