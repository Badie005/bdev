package root

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/badie/bdev/internal/ai/agents"
	aicmd "github.com/badie/bdev/internal/cmd/ai"
	gitcmd "github.com/badie/bdev/internal/cmd/git"
	multicmd "github.com/badie/bdev/internal/cmd/multi"
	projectcmd "github.com/badie/bdev/internal/cmd/projects"
	secretscmd "github.com/badie/bdev/internal/cmd/secrets"
	"github.com/badie/bdev/internal/cmd/version"
	workflowcmd "github.com/badie/bdev/internal/cmd/workflow"
	"github.com/badie/bdev/internal/core/repl"
	"github.com/badie/bdev/internal/core/runner"
	"github.com/badie/bdev/pkg/ui"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:                        "bdev",
	Short:                      "B.DEV CLI - Enterprise Developer Agent",
	Long:                       ui.Muted("  Enterprise Developer Workstation\n  Powered by Go & Ollama AI"),
	SilenceUsage:               true,
	SilenceErrors:              true,
	SuggestionsMinimumDistance: 1,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		noColor, _ := cmd.Flags().GetBool("no-color")
		if noColor {
			color.NoColor = true
		}
	},
	Run: func(cmd *cobra.Command, _ []string) {
		// No args = start interactive REPL
		repl.Start(cmd)
	},
}

// Execute runs the root command with unified error handling
func Execute() {
	// Trap panics for robust crash handling
	defer func() {
		if r := recover(); r != nil {
			fmt.Println()
			fmt.Println(ui.Error("CRITICAL SYSTEM FAILURE"))
			fmt.Println(ui.Muted(fmt.Sprintf("Panic: %v", r)))
			os.Exit(1)
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		// Handle specific known errors with style
		msg := err.Error()

		if strings.Contains(msg, "unknown command") {
			fmt.Println(ui.Error("Unknown command. Try 'help' to see available tools."))
		} else if strings.Contains(msg, "connection refused") {
			fmt.Println(ui.Error("Connection check failed. Is the service running?"))
		} else {
			fmt.Println(ui.Error(msg))
		}

		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.bdev/config.json)")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose logging")

	// Register subcommands
	rootCmd.AddCommand(version.NewCommand())
	rootCmd.AddCommand(gitcmd.NewCommand())
	rootCmd.AddCommand(projectcmd.NewCommand())
	rootCmd.AddCommand(aicmd.NewCommand())
	rootCmd.AddCommand(agents.NewCommand())
	rootCmd.AddCommand(secretscmd.NewCommand())
	rootCmd.AddCommand(workflowcmd.NewCommand())
	rootCmd.AddCommand(multicmd.NewCommand())

	// Register quick action shortcuts
	registerShortcuts()

	// book II: THE SPECTRUM - Omega Help Renderer
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		ui.PrintWelcome()

		fmt.Println(ui.Bold("USAGE"))
		fmt.Printf("  %s %s\n\n", ui.Primary("bdev"), ui.Muted("[command] [flags]"))

		fmt.Println(ui.Bold("CORE COMMANDS"))
		printCommandGroup(cmd, []string{"projects", "git", "ai", "agents", "workflow"})

		fmt.Println(ui.Bold("UTILITIES"))
		printCommandGroup(cmd, []string{"secrets", "multi", "config", "version", "help"})

		fmt.Println(ui.Bold("FLAGS"))
		cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
			fmt.Printf("  %s %s\n",
				ui.Info(fmt.Sprintf("--%s", f.Name)),
				ui.Muted(f.Usage))
		})
		fmt.Println()

		fmt.Println(ui.Muted("Use 'bdev [command] --help' for details."))
		fmt.Println()
	})
}

func printCommandGroup(parent *cobra.Command, names []string) {
	for _, name := range names {
		for _, cmd := range parent.Commands() {
			if cmd.Name() == name {
				fmt.Printf("  %s  %s\n",
					ui.Cyan(fmt.Sprintf("%-12s", cmd.Name())),
					ui.Muted(cmd.Short))
			}
		}
	}
	fmt.Println()
}

func registerShortcuts() {
	// bdev list - shortcut for projects list
	rootCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all projects (shortcut for 'projects list')",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Delegate to projects list
			projectsCmd, _, _ := rootCmd.Find([]string{"projects", "list"})
			if projectsCmd != nil {
				return projectsCmd.RunE(projectsCmd, args)
			}
			return nil
		},
	})

	// bdev start - start dev server in current project
	rootCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start development server (current project)",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := runner.NewFromCwd()
			if err != nil {
				return err
			}
			return r.Start()
		},
	})

	// bdev test - run tests in current project
	rootCmd.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "Run tests (current project)",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := runner.NewFromCwd()
			if err != nil {
				return err
			}
			return r.Test(false)
		},
	})

	// bdev build - build current project
	rootCmd.AddCommand(&cobra.Command{
		Use:   "build",
		Short: "Build for production (current project)",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := runner.NewFromCwd()
			if err != nil {
				return err
			}
			return r.Build()
		},
	})

	// bdev install - install dependencies
	rootCmd.AddCommand(&cobra.Command{
		Use:   "install",
		Short: "Install dependencies (current project)",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := runner.NewFromCwd()
			if err != nil {
				return err
			}
			return r.Install()
		},
	})

	// bdev lint - run linter
	var fix bool
	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Run linter (current project)",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := runner.NewFromCwd()
			if err != nil {
				return err
			}
			return r.Lint(fix)
		},
	}
	lintCmd.Flags().BoolVar(&fix, "fix", false, "Auto-fix issues")
	rootCmd.AddCommand(lintCmd)

	// bdev demo - show visual components
	rootCmd.AddCommand(&cobra.Command{
		Use:     "demo",
		Aliases: []string{"ui"},
		Short:   "Showcase Omega Design System components",
		Run: func(cmd *cobra.Command, args []string) {
			ui.Demo()
		},
	})
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return
		}

		viper.AddConfigPath(home + "/Dev/.bdev")
		viper.AddConfigPath(home + "/.bdev")
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	_ = viper.ReadInConfig()
}
