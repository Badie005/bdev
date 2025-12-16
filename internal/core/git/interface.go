package git

import "context"

// Repo defines the interface for git repository operations.
// This interface enables dependency injection and mocking in tests.
type Repo interface {
	// Status returns the current repository status
	Status() (*Status, error)

	// CurrentBranch returns the name of the current branch
	CurrentBranch() string

	// Log returns the commit history
	Log(count int) ([]Commit, error)

	// LogGraph returns commit history with ASCII graph
	LogGraph(count int) (string, error)

	// ListBranches returns all local and remote branches
	ListBranches() ([]BranchInfo, error)

	// Diff returns the diff output
	Diff(staged bool) (string, error)

	// DiffStat returns diff statistics
	DiffStat(staged bool) (string, error)

	// Remote returns the origin remote URL
	Remote() string

	// LastCommit returns the most recent commit
	LastCommit() (*Commit, error)

	// Add stages files for commit
	Add(files ...string) error

	// Commit creates a new commit with the given message
	Commit(message string, all bool) error

	// Push pushes commits to remote
	Push(force bool) error

	// Pull pulls changes from remote
	Pull() error

	// Fetch fetches from all remotes
	Fetch() error

	// CreateBranch creates a new branch
	CreateBranch(name string) error

	// SwitchBranch switches to an existing branch
	SwitchBranch(name string) error

	// CreateAndSwitch creates a new branch and switches to it
	CreateAndSwitch(name string) error

	// DeleteBranch deletes a branch
	DeleteBranch(name string, force bool) error

	// Stash stashes current changes
	Stash(message string) error

	// StashPop pops the latest stash
	StashPop() error

	// StashList lists all stashes
	StashList() ([]string, error)

	// StashDrop drops a stash by index
	StashDrop(index int) error

	// Reset resets changes
	Reset(hard bool, ref string) error

	// ResetFile unstages a specific file
	ResetFile(file string) error
}

// RepoOpener is a factory interface for opening repositories
type RepoOpener interface {
	Open(path string) (Repo, error)
	OpenCurrent() (Repo, error)
	IsRepo(path string) bool
}

// Ensure Repository implements Repo interface
var _ Repo = (*Repository)(nil)

// ContextualRepo extends Repo with context-aware operations
// This is for future use with cancellation and timeouts
type ContextualRepo interface {
	Repo
	StatusWithContext(ctx context.Context) (*Status, error)
	CommitWithContext(ctx context.Context, message string, all bool) error
}
