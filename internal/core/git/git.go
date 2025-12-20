package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Repository represents a git repository
type Repository struct {
	Path   string
	branch string // cached
}

// Status represents parsed git status
type Status struct {
	Branch       string
	Remote       string
	Ahead        int
	Behind       int
	Staged       []FileChange
	Modified     []FileChange
	Untracked    []FileChange
	Deleted      []FileChange
	HasConflicts bool
	IsClean      bool
}

// FileChange represents a changed file
type FileChange struct {
	Path   string
	Status string // M, A, D, R, C, U, ?
}

// Commit represents a git commit
type Commit struct {
	Hash      string
	ShortHash string
	Author    string
	Email     string
	Date      string
	Message   string
	Graph     string // For log --graph
}

// BranchInfo represents branch information
type BranchInfo struct {
	Name      string
	IsCurrent bool
	IsRemote  bool
	Upstream  string
}

// Open opens a git repository at path
func Open(path string) (*Repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if !IsRepo(absPath) {
		return nil, fmt.Errorf("not a git repository: %s", absPath)
	}

	return &Repository{Path: absPath}, nil
}

// OpenCurrent opens the git repository in current directory
func OpenCurrent() (*Repository, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return Open(cwd)
}

// IsRepo checks if path is a git repository
func IsRepo(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	return cmd.Run() == nil
}

// run executes a git command and returns output
func (r *Repository) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("%s", errMsg)
	}

	return strings.TrimSpace(stdout.String()), nil
}

// Status returns parsed git status
func (r *Repository) Status() (*Status, error) {
	output, err := r.run("status", "--porcelain=v2", "--branch")
	if err != nil {
		return nil, err
	}

	return ParseStatus(output), nil
}

// CurrentBranch returns the current branch name
func (r *Repository) CurrentBranch() string {
	if r.branch != "" {
		return r.branch
	}

	branch, err := r.run("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return ""
	}

	r.branch = branch
	return branch
}

// Commit creates a new commit
func (r *Repository) Commit(message string, all bool) error {
	var args []string
	if all {
		args = []string{"commit", "-a", "-m", message}
	} else {
		args = []string{"commit", "-m", message}
	}

	_, err := r.run(args...)
	return err
}

// Add stages files
func (r *Repository) Add(files ...string) error {
	if len(files) == 0 {
		files = []string{"."}
	}
	args := append([]string{"add"}, files...)
	_, err := r.run(args...)
	return err
}

// Push pushes to remote
func (r *Repository) Push(force bool) error {
	args := []string{"push"}
	if force {
		args = append(args, "--force-with-lease")
	}
	_, err := r.run(args...)
	return err
}

// Pull pulls from remote
func (r *Repository) Pull() error {
	_, err := r.run("pull", "--ff-only")
	return err
}

// Fetch fetches from remote
func (r *Repository) Fetch() error {
	_, err := r.run("fetch", "--all", "--prune")
	return err
}

// Log returns commit history
func (r *Repository) Log(count int) ([]Commit, error) {
	format := "--format=%H|%h|%an|%ae|%ar|%s"
	output, err := r.run("log", fmt.Sprintf("-%d", count), format, "--no-merges")
	if err != nil {
		return nil, err
	}

	return ParseLog(output), nil
}

// LogGraph returns commit history with graph
func (r *Repository) LogGraph(count int) (string, error) {
	return r.run("log", fmt.Sprintf("-%d", count), "--oneline", "--graph", "--decorate", "--all")
}

// ListBranches returns all branches
func (r *Repository) ListBranches() ([]BranchInfo, error) {
	output, err := r.run("branch", "-a", "-v")
	if err != nil {
		return nil, err
	}

	return ParseBranches(output), nil
}

// CreateBranch creates a new branch
func (r *Repository) CreateBranch(name string) error {
	_, err := r.run("branch", name)
	return err
}

// SwitchBranch switches to a branch
func (r *Repository) SwitchBranch(name string) error {
	_, err := r.run("checkout", name)
	return err
}

// CreateAndSwitch creates and switches to a new branch
func (r *Repository) CreateAndSwitch(name string) error {
	_, err := r.run("checkout", "-b", name)
	return err
}

// DeleteBranch deletes a branch
func (r *Repository) DeleteBranch(name string, force bool) error {
	flag := "-d"
	if force {
		flag = "-D"
	}
	_, err := r.run("branch", flag, name)
	return err
}

// Stash stashes current changes
func (r *Repository) Stash(message string) error {
	args := []string{"stash", "push"}
	if message != "" {
		args = append(args, "-m", message)
	}
	_, err := r.run(args...)
	return err
}

// StashPop pops the latest stash
func (r *Repository) StashPop() error {
	_, err := r.run("stash", "pop")
	return err
}

// StashList lists all stashes
func (r *Repository) StashList() ([]string, error) {
	output, err := r.run("stash", "list")
	if err != nil {
		return nil, err
	}

	if output == "" {
		return []string{}, nil
	}

	return strings.Split(output, "\n"), nil
}

// StashDrop drops a stash
func (r *Repository) StashDrop(index int) error {
	_, err := r.run("stash", "drop", fmt.Sprintf("stash@{%d}", index))
	return err
}

// Diff returns diff output
func (r *Repository) Diff(staged bool) (string, error) {
	args := []string{"diff", "--color=always"}
	if staged {
		args = append(args, "--staged")
	}
	return r.run(args...)
}

// DiffStat returns diff statistics
func (r *Repository) DiffStat(staged bool) (string, error) {
	args := []string{"diff", "--stat"}
	if staged {
		args = append(args, "--staged")
	}
	return r.run(args...)
}

// Reset resets changes
func (r *Repository) Reset(hard bool, ref string) error {
	args := []string{"reset"}
	if hard {
		args = append(args, "--hard")
	}
	if ref != "" {
		args = append(args, ref)
	}
	_, err := r.run(args...)
	return err
}

// ResetFile unstages a file
func (r *Repository) ResetFile(file string) error {
	_, err := r.run("reset", "HEAD", "--", file)
	return err
}

// Remote returns the remote URL
func (r *Repository) Remote() string {
	url, _ := r.run("remote", "get-url", "origin")
	return url
}

// LastCommit returns the last commit info
func (r *Repository) LastCommit() (*Commit, error) {
	commits, err := r.Log(1)
	if err != nil {
		return nil, err
	}
	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits")
	}
	return &commits[0], nil
}

// ============================================================
// Parser Functions
// ============================================================

// ParseStatus parses git status --porcelain=v2 --branch output
func ParseStatus(output string) *Status {
	status := &Status{
		Staged:    make([]FileChange, 0),
		Modified:  make([]FileChange, 0),
		Untracked: make([]FileChange, 0),
		Deleted:   make([]FileChange, 0),
		IsClean:   true,
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		// Branch info
		if strings.HasPrefix(line, "# branch.head ") {
			status.Branch = strings.TrimPrefix(line, "# branch.head ")
			continue
		}
		if strings.HasPrefix(line, "# branch.upstream ") {
			status.Remote = strings.TrimPrefix(line, "# branch.upstream ")
			continue
		}
		if strings.HasPrefix(line, "# branch.ab ") {
			ab := strings.TrimPrefix(line, "# branch.ab ")
			parts := strings.Fields(ab)
			if len(parts) >= 2 {
				ahead, err := strconv.Atoi(strings.TrimPrefix(parts[0], "+"))
				if err == nil {
					status.Ahead = ahead
				}
				behind, err := strconv.Atoi(strings.TrimPrefix(parts[1], "-"))
				if err == nil {
					status.Behind = behind
				}
			}
			continue
		}

		// File changes (porcelain v2)
		if strings.HasPrefix(line, "1 ") || strings.HasPrefix(line, "2 ") {
			// Changed file
			parts := strings.Fields(line)
			if len(parts) >= 9 {
				xy := parts[1] // XY status
				path := parts[len(parts)-1]

				change := FileChange{Path: path, Status: xy}

				// Index status (staged)
				if xy[0] != '.' {
					status.Staged = append(status.Staged, change)
					status.IsClean = false
				}
				// Worktree status (modified)
				if xy[1] != '.' {
					status.Modified = append(status.Modified, change)
					status.IsClean = false
				}
			}
			continue
		}

		// Untracked
		if strings.HasPrefix(line, "? ") {
			path := strings.TrimPrefix(line, "? ")
			status.Untracked = append(status.Untracked, FileChange{Path: path, Status: "?"})
			status.IsClean = false
			continue
		}

		// Conflicts
		if strings.HasPrefix(line, "u ") {
			status.HasConflicts = true
			status.IsClean = false
		}
	}

	return status
}

// ParseLog parses git log output
func ParseLog(output string) []Commit {
	commits := make([]Commit, 0)

	if output == "" {
		return commits
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "|", 6)
		if len(parts) == 6 {
			commits = append(commits, Commit{
				Hash:      parts[0],
				ShortHash: parts[1],
				Author:    parts[2],
				Email:     parts[3],
				Date:      parts[4],
				Message:   parts[5],
			})
		}
	}

	return commits
}

// ParseBranches parses git branch output
func ParseBranches(output string) []BranchInfo {
	branches := make([]BranchInfo, 0)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		branch := BranchInfo{}

		// Current branch
		if strings.HasPrefix(line, "* ") {
			branch.IsCurrent = true
			line = strings.TrimPrefix(line, "* ")
		} else {
			line = strings.TrimPrefix(line, "  ")
		}

		// Remote branch
		if strings.HasPrefix(line, "remotes/") {
			branch.IsRemote = true
			line = strings.TrimPrefix(line, "remotes/")
		}

		// Extract branch name (first field)
		parts := strings.Fields(line)
		if len(parts) > 0 {
			branch.Name = parts[0]
		}

		branches = append(branches, branch)
	}

	return branches
}
