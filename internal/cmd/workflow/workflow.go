package workflowcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/badie/bdev/internal/core/config"
	"github.com/badie/bdev/internal/core/vault"
	"github.com/badie/bdev/internal/core/workflow"
	"github.com/badie/bdev/pkg/ui"
)

var engine *workflow.Engine

func getEngine() *workflow.Engine {
	if engine == nil {
		cfg := config.Get()
		engine = workflow.New(filepath.Join(cfg.Paths.Bdev, "workflows"))
		// Inject Vault instance (locked initially)
		engine.Vault = vault.New(cfg.VaultFile())
	}
	return engine
}

// NewCommand creates the workflow command group
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "workflow",
		Aliases: []string{"wf", "flow"},
		Short:   "Run automation workflows",
		Long:    "Execute YAML-based automation workflows",
	}

	cmd.AddCommand(listCmd())
	cmd.AddCommand(runCmd())
	cmd.AddCommand(showCmd())

	return cmd
}

// ============================================================
// LIST - List available workflows
// ============================================================

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available workflows",
		RunE: func(cmd *cobra.Command, args []string) error {
			eng := getEngine()

			workflows, err := eng.List()
			if err != nil {
				return err
			}

			if len(workflows) == 0 {
				fmt.Println(ui.Muted("No workflows found"))
				fmt.Println(ui.Muted("Create workflows in: " + eng.WorkflowDir))
				return nil
			}

			fmt.Println(ui.Bold(fmt.Sprintf("Workflows (%d)", len(workflows))))
			for _, name := range workflows {
				wf, err := eng.Load(name)
				if err != nil {
					fmt.Printf("  %s %s\n", ui.Primary(name), ui.Error("(invalid)"))
					continue
				}
				desc := wf.Description
				if desc == "" {
					desc = fmt.Sprintf("%d steps", len(wf.Steps))
				}
				fmt.Printf("  %s  %s\n", ui.Primary(name), ui.Muted(desc))
			}

			return nil
		},
	}
}

func runCmd() *cobra.Command {
	var verbose bool
	var withSecrets bool

	cmd := &cobra.Command{
		Use:   "run <name>",
		Short: "Execute a workflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			eng := getEngine()
			eng.Verbose = verbose

			name := args[0]
			wf, err := eng.Load(name)
			if err != nil {
				return err
			}

			// Check if workflow needs secrets
			needsSecrets := false
			// Naive check in steps
			for _, s := range wf.Steps {
				if strings.Contains(s.Run, "secrets.") || strings.Contains(s.Cwd, "secrets.") {
					needsSecrets = true
					break
				}
				for _, v := range s.Env {
					if strings.Contains(v, "secrets.") {
						needsSecrets = true
						break
					}
				}
			}

			// Also check workflow env
			if !needsSecrets {
				for _, v := range wf.Env {
					if strings.Contains(v, "secrets.") {
						needsSecrets = true
						break
					}
				}
			}

			// Handle Vault Unlock
			if (needsSecrets || withSecrets) && eng.Vault != nil {
				// We need to match the Vault interface method Unlock, but Engine defines Vault as interface
				// We know concrete type is *vault.Vault so we assert or use the concrete type method
				if v, ok := eng.Vault.(*vault.Vault); ok {
					if v.Exists() {
						fmt.Print(ui.Bold("Vault password required: "))
						bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
						fmt.Println()
						if err != nil {
							return fmt.Errorf("failed to read password: %w", err)
						}

						if err := v.Unlock(string(bytePassword)); err != nil {
							return fmt.Errorf("failed to unlock vault: %w", err)
						}
						fmt.Println(ui.Success("Vault unlocked"))
						fmt.Println()
					} else {
						fmt.Println(ui.Warning("Vault not initialized. Secrets will not be expanded."))
					}
				}
			}

			fmt.Println(ui.Bold("Running: " + wf.Name))
			if wf.Description != "" {
				fmt.Println(ui.Muted(wf.Description))
			}
			fmt.Println()

			// ... (Callback setup remains same)

			stepNum := 0
			eng.OnStep = func(step workflow.Step, result *workflow.StepResult) {
				stepNum++
				status := ui.Success(ui.ActiveGlyphs.Check)
				if !result.Success {
					status = ui.Error(ui.ActiveGlyphs.Cross)
				}
				fmt.Printf("%s %d. %s %s\n", status, stepNum, step.Name, ui.Muted(fmt.Sprintf("(%s)", result.Duration.Round(100*1e6))))

				if verbose && result.Output != "" {
					fmt.Println(ui.Muted(result.Output))
				}
			}

			result := eng.Execute(wf)

			fmt.Println()
			if result.Success {
				fmt.Println(ui.Success("Workflow completed successfully"))
			} else {
				fmt.Println(ui.Error("Workflow failed"))
			}
			fmt.Printf("Total time: %s\n", ui.Muted(result.Duration.Round(100*1e6).String()))

			if !result.Success {
				return fmt.Errorf("workflow failed")
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show step output")
	cmd.Flags().BoolVar(&withSecrets, "secrets", false, "Force unlock vault")
	return cmd
}

// ============================================================
// SHOW - Show workflow details
// ============================================================

func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show workflow details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			eng := getEngine()

			name := args[0]
			wf, err := eng.Load(name)
			if err != nil {
				return err
			}

			fmt.Println(ui.Bold(wf.Name))
			if wf.Description != "" {
				fmt.Println(ui.Muted(wf.Description))
			}
			fmt.Println()

			fmt.Println(ui.Bold("Steps:"))
			for i, step := range wf.Steps {
				fmt.Printf("  %d. %s\n", i+1, ui.Primary(step.Name))
				fmt.Printf("     %s\n", ui.Muted(step.Run))
			}

			if len(wf.OnSuccess) > 0 {
				fmt.Println()
				fmt.Println(ui.Bold("On Success:"))
				for _, step := range wf.OnSuccess {
					fmt.Printf("  - %s\n", step.Name)
				}
			}

			if len(wf.OnFailure) > 0 {
				fmt.Println()
				fmt.Println(ui.Bold("On Failure:"))
				for _, step := range wf.OnFailure {
					fmt.Printf("  - %s\n", step.Name)
				}
			}

			return nil
		},
	}
}
