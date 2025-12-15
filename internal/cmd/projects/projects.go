package projectcmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/badie/bdev/internal/core/config"
	"github.com/badie/bdev/internal/core/projects"
	"github.com/badie/bdev/pkg/ui"
)

// NewCommand creates the projects command group
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "projects",
		Aliases: []string{"p", "proj"},
		Short:   "Manage development projects",
		Long:    "List, create, and manage your development projects",
	}

	cmd.AddCommand(listCmd())
	cmd.AddCommand(openCmd())
	cmd.AddCommand(runCmd())
	cmd.AddCommand(findCmd())
	cmd.AddCommand(describeCmd())

	return cmd
}

// ============================================================
// LIST
// ============================================================

func listCmd() *cobra.Command {
	var sortBy string
	var filterType string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Get()
			projectsDir := cfg.Paths.Projects

			projectList, err := projects.Scan(projectsDir)
			if err != nil {
				return err
			}

			if len(projectList) == 0 {
				fmt.Println(ui.Muted("No projects found in " + projectsDir))
				return nil
			}

			fmt.Println(ui.Bold(fmt.Sprintf("Projects (%d)", len(projectList))))
			fmt.Println(ui.Muted(ui.SeparatorLight))
			fmt.Println()

			for _, p := range projectList {
				icon := p.Type.Icon()
				typeName := p.Type.String()

				line := fmt.Sprintf("  %s %s", icon, ui.Bold(p.Name))
				line += ui.Muted(fmt.Sprintf(" (%s)", typeName))

				if p.GitBranch != "" {
					line += ui.Primary(fmt.Sprintf(" [%s]", p.GitBranch))
				}

				fmt.Println(line)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&sortBy, "sort", "s", "name", "Sort by: name, type, date")
	cmd.Flags().StringVarP(&filterType, "type", "t", "", "Filter by type")
	return cmd
}

// ============================================================
// OPEN
// ============================================================

func openCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "open [name]",
		Short: "Open project in VS Code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Get()
			projectName := args[0]

			// Find project
			projectPath := filepath.Join(cfg.Paths.Projects, projectName)
			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				return fmt.Errorf("project not found: %s", projectName)
			}

			// Open in VS Code
			vscode := exec.Command("code", projectPath)
			if err := vscode.Start(); err != nil {
				return fmt.Errorf("failed to open VS Code: %v", err)
			}

			fmt.Println(ui.Success("Opening ") + ui.Bold(projectName) + ui.Muted(" in VS Code"))
			return nil
		},
	}
}

// ============================================================
// RUN
// ============================================================

func runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run [name] [script]",
		Short: "Run a project script",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Get()
			projectName := args[0]
			script := "dev"
			if len(args) > 1 {
				script = args[1]
			}

			projectPath := filepath.Join(cfg.Paths.Projects, projectName)
			project := projects.Analyze(projectPath)
			if project == nil {
				return fmt.Errorf("project not found: %s", projectName)
			}

			// Get the appropriate command
			var runCmd *exec.Cmd
			switch script {
			case "dev", "start":
				bin, cmdArgs := project.GetStartCommand()
				if bin == "" {
					return fmt.Errorf("no start command for this project type")
				}
				runCmd = exec.Command(bin, cmdArgs...)
			case "test":
				bin, cmdArgs := project.GetTestCommand()
				if bin == "" {
					return fmt.Errorf("no test command for this project type")
				}
				runCmd = exec.Command(bin, cmdArgs...)
			case "build":
				bin, cmdArgs := project.GetBuildCommand()
				if bin == "" {
					return fmt.Errorf("no build command for this project type")
				}
				runCmd = exec.Command(bin, cmdArgs...)
			default:
				// Run npm script
				runCmd = exec.Command("npm", "run", script)
			}

			runCmd.Dir = projectPath
			runCmd.Stdout = os.Stdout
			runCmd.Stderr = os.Stderr
			runCmd.Stdin = os.Stdin

			fmt.Println(ui.Primary("Running ") + ui.Bold(script) + ui.Muted(" in ") + project.Name)
			fmt.Println()

			return runCmd.Run()
		},
	}
}

// ============================================================
// FIND
// ============================================================

func findCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "find [query]",
		Short: "Search for projects by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Get()
			query := args[0]

			projectList, err := projects.Scan(cfg.Paths.Projects)
			if err != nil {
				return err
			}

			matches := make([]projects.Project, 0)
			queryLower := strings.ToLower(query)

			for _, p := range projectList {
				if containsIgnoreCase(p.Name, queryLower) {
					matches = append(matches, p)
				}
			}

			if len(matches) == 0 {
				fmt.Println(ui.Muted("No projects matching: " + query))
				return nil
			}

			fmt.Println(ui.Bold(fmt.Sprintf("Found %d project(s)", len(matches))))
			for _, p := range matches {
				fmt.Printf("  %s %s %s\n", p.Type.Icon(), ui.Bold(p.Name), ui.Muted(p.Path))
			}

			return nil
		},
	}
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// ============================================================
// DESCRIBE
// ============================================================

func describeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "describe [name]",
		Short: "Show detailed project information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Get()
			projectName := args[0]

			projectPath := filepath.Join(cfg.Paths.Projects, projectName)
			project := projects.Analyze(projectPath)
			if project == nil {
				return fmt.Errorf("project not found: %s", projectName)
			}

			fmt.Println(ui.Bold(project.Name))
			fmt.Println(ui.Muted(ui.SeparatorLight))
			fmt.Printf("  %s %s\n", ui.Primary("Type:"), project.Type.String())
			fmt.Printf("  %s %s\n", ui.Primary("Path:"), project.Path)

			if project.Framework != "" {
				fmt.Printf("  %s %s\n", ui.Primary("Framework:"), project.Framework)
			}

			if project.HasGit {
				fmt.Printf("  %s %s\n", ui.Primary("Git:"), ui.Success("yes"))
				fmt.Printf("  %s %s\n", ui.Primary("Branch:"), project.GitBranch)
			}

			// Scripts
			if len(project.Scripts) > 0 {
				fmt.Println()
				fmt.Println(ui.Bold("Available Scripts:"))
				for name, script := range project.Scripts {
					if len(script) > 50 {
						script = script[:50] + "..."
					}
					fmt.Printf("  %s %s\n", ui.Primary(name+":"), ui.Muted(script))
				}
			}

			return nil
		},
	}
}
