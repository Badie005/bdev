package gitcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/badie/bdev/internal/core/git"
	"github.com/badie/bdev/pkg/ui"
)

// NewCommand creates the git command group
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git",
		Short: "Git operations with enhanced output",
		Long:  "Native Git integration with colored output and simplified commands",
	}

	cmd.AddCommand(statusCmd())
	cmd.AddCommand(commitCmd())
	cmd.AddCommand(pushCmd())
	cmd.AddCommand(pullCmd())
	cmd.AddCommand(logCmd())
	cmd.AddCommand(branchCmd())
	cmd.AddCommand(stashCmd())
	cmd.AddCommand(diffCmd())
	cmd.AddCommand(addCmd())
	cmd.AddCommand(resetCmd())

	return cmd
}

// ============================================================
// STATUS
// ============================================================

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show working tree status",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.OpenCurrent()
			if err != nil {
				return err
			}

			status, err := repo.Status()
			if err != nil {
				return err
			}

			// Branch info
			branchInfo := status.Branch
			if status.Remote != "" {
				branchInfo += ui.Muted(" → " + status.Remote)
			}
			if status.Ahead > 0 {
				branchInfo += ui.Success(fmt.Sprintf(" ↑%d", status.Ahead))
			}
			if status.Behind > 0 {
				branchInfo += ui.Warning(fmt.Sprintf(" ↓%d", status.Behind))
			}
			fmt.Println(ui.Bold("Branch: ") + ui.Primary(branchInfo))

			if status.IsClean {
				fmt.Println(ui.Success("Working tree clean"))
				return nil
			}

			// Staged changes
			if len(status.Staged) > 0 {
				fmt.Println(ui.Bold("\nStaged for commit:"))
				for _, f := range status.Staged {
					fmt.Printf("  %s %s\n", ui.Success("+"), f.Path)
				}
			}

			// Modified
			if len(status.Modified) > 0 {
				fmt.Println(ui.Bold("\nModified:"))
				for _, f := range status.Modified {
					fmt.Printf("  %s %s\n", ui.Warning("M"), f.Path)
				}
			}

			// Untracked
			if len(status.Untracked) > 0 {
				fmt.Println(ui.Bold("\nUntracked:"))
				for _, f := range status.Untracked {
					fmt.Printf("  %s %s\n", ui.Muted("?"), f.Path)
				}
			}

			// Conflicts
			if status.HasConflicts {
				fmt.Println(ui.Error("Merge conflicts detected!"))
			}

			return nil
		},
	}
}

// ============================================================
// COMMIT
// ============================================================

func commitCmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "commit [message]",
		Short: "Create a new commit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.OpenCurrent()
			if err != nil {
				return err
			}

			message := args[0]
			if err := repo.Commit(message, all); err != nil {
				return err
			}

			fmt.Println(ui.Success("Commit created: ") + ui.Muted(message))
			return nil
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Stage all modified files")
	return cmd
}

// ============================================================
// PUSH
// ============================================================

func pushCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push commits to remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.OpenCurrent()
			if err != nil {
				return err
			}

			fmt.Println(ui.Primary("Pushing..."))
			if err := repo.Push(force); err != nil {
				return err
			}

			fmt.Println(ui.Success("Push successful"))
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force push (with lease)")
	return cmd
}

// ============================================================
// PULL
// ============================================================

func pullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Pull from remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.OpenCurrent()
			if err != nil {
				return err
			}

			fmt.Println(ui.Primary("Pulling..."))
			if err := repo.Pull(); err != nil {
				return err
			}

			fmt.Println(ui.Success("Pull successful"))
			return nil
		},
	}
}

// ============================================================
// LOG
// ============================================================

func logCmd() *cobra.Command {
	var count int
	var graph bool

	cmd := &cobra.Command{
		Use:   "log",
		Short: "Show commit history",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.OpenCurrent()
			if err != nil {
				return err
			}

			if graph {
				output, err := repo.LogGraph(count)
				if err != nil {
					return err
				}
				fmt.Println(output)
				return nil
			}

			commits, err := repo.Log(count)
			if err != nil {
				return err
			}

			for _, c := range commits {
				fmt.Printf("%s %s %s\n",
					ui.Warning(c.ShortHash),
					c.Message,
					ui.Muted("("+c.Date+" by "+c.Author+")"))
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&count, "count", "n", 10, "Number of commits to show")
	cmd.Flags().BoolVarP(&graph, "graph", "g", false, "Show graph view")
	return cmd
}

// ============================================================
// BRANCH
// ============================================================

func branchCmd() *cobra.Command {
	var delete bool
	var force bool

	cmd := &cobra.Command{
		Use:   "branch [name]",
		Short: "List, create, or switch branches",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.OpenCurrent()
			if err != nil {
				return err
			}

			// No args = list branches
			if len(args) == 0 {
				branches, err := repo.ListBranches()
				if err != nil {
					return err
				}

				for _, b := range branches {
					prefix := "  "
					if b.IsCurrent {
						prefix = ui.Success("* ")
					}
					name := b.Name
					if b.IsRemote {
						name = ui.Muted(name)
					}
					fmt.Println(prefix + name)
				}
				return nil
			}

			name := args[0]

			// Delete branch
			if delete {
				if err := repo.DeleteBranch(name, force); err != nil {
					return err
				}
				fmt.Println(ui.Success("Branch deleted: ") + name)
				return nil
			}

			// Switch or create
			if err := repo.SwitchBranch(name); err != nil {
				// Try creating new branch
				if err := repo.CreateAndSwitch(name); err != nil {
					return err
				}
				fmt.Println(ui.Success("Created and switched to: ") + name)
				return nil
			}

			fmt.Println(ui.Success("Switched to: ") + name)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&delete, "delete", "d", false, "Delete branch")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force delete")
	return cmd
}

// ============================================================
// STASH
// ============================================================

func stashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stash [push|pop|list|drop]",
		Short: "Stash changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.OpenCurrent()
			if err != nil {
				return err
			}

			action := "push"
			if len(args) > 0 {
				action = args[0]
			}

			switch action {
			case "push":
				if err := repo.Stash(""); err != nil {
					return err
				}
				fmt.Println(ui.Success("Changes stashed"))

			case "pop":
				if err := repo.StashPop(); err != nil {
					return err
				}
				fmt.Println(ui.Success("Stash popped"))

			case "list":
				stashes, err := repo.StashList()
				if err != nil {
					return err
				}
				if len(stashes) == 0 {
					fmt.Println(ui.Muted("No stashes"))
					return nil
				}
				for _, s := range stashes {
					fmt.Println(s)
				}

			case "drop":
				if err := repo.StashDrop(0); err != nil {
					return err
				}
				fmt.Println(ui.Success("Stash dropped"))

			default:
				return fmt.Errorf("unknown stash action: %s", action)
			}

			return nil
		},
	}

	return cmd
}

// ============================================================
// DIFF
// ============================================================

func diffCmd() *cobra.Command {
	var staged bool
	var stat bool

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Show changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.OpenCurrent()
			if err != nil {
				return err
			}

			var output string
			if stat {
				output, err = repo.DiffStat(staged)
			} else {
				output, err = repo.Diff(staged)
			}

			if err != nil {
				return err
			}

			if output == "" {
				fmt.Println(ui.Muted("No changes"))
				return nil
			}

			fmt.Println(output)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&staged, "staged", "s", false, "Show staged changes")
	cmd.Flags().BoolVar(&stat, "stat", false, "Show statistics only")
	return cmd
}

// ============================================================
// ADD
// ============================================================

func addCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [files...]",
		Short: "Stage files for commit",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.OpenCurrent()
			if err != nil {
				return err
			}

			if err := repo.Add(args...); err != nil {
				return err
			}

			if len(args) == 0 {
				fmt.Println(ui.Success("All files staged"))
			} else {
				fmt.Println(ui.Success("Staged: ") + ui.Muted(fmt.Sprintf("%v", args)))
			}
			return nil
		},
	}
}

// ============================================================
// RESET
// ============================================================

func resetCmd() *cobra.Command {
	var hard bool

	cmd := &cobra.Command{
		Use:   "reset [ref]",
		Short: "Reset changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.OpenCurrent()
			if err != nil {
				return err
			}

			ref := ""
			if len(args) > 0 {
				ref = args[0]
			}

			if hard {
				fmt.Println(ui.Warning("Hard reset will discard all changes!"))
			}

			if err := repo.Reset(hard, ref); err != nil {
				return err
			}

			fmt.Println(ui.Success("Reset complete"))
			return nil
		},
	}

	cmd.Flags().BoolVar(&hard, "hard", false, "Hard reset (discard all changes)")
	return cmd
}
