package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// ==============================================================
// Parser Tests - Table-Driven
// ==============================================================

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *Status
	}{
		{
			name: "full_status_with_all_fields",
			input: `# branch.oid 0000000000000000000000000000000000000000
# branch.head master
# branch.upstream origin/master
# branch.ab +1 -2
1 .M N... 100644 100644 100644 9daeafb9864cf43055ae93beb0afd6c7d144bfa4 9daeafb9864cf43055ae93beb0afd6c7d144bfa4 modified_file.go
1 M. N... 100644 100644 100644 9daeafb9864cf43055ae93beb0afd6c7d144bfa4 9daeafb9864cf43055ae93beb0afd6c7d144bfa4 staged_file.go
? untracked_file.txt
`,
			want: &Status{
				Branch:  "master",
				Remote:  "origin/master",
				Ahead:   1,
				Behind:  2,
				IsClean: false,
				Modified: []FileChange{
					{Path: "modified_file.go", Status: ".M"},
				},
				Staged: []FileChange{
					{Path: "staged_file.go", Status: "M."},
				},
				Untracked: []FileChange{
					{Path: "untracked_file.txt", Status: "?"},
				},
				Deleted: []FileChange{},
			},
		},
		{
			name: "clean_repo",
			input: `# branch.oid abc123
# branch.head main
# branch.upstream origin/main
# branch.ab +0 -0
`,
			want: &Status{
				Branch:    "main",
				Remote:    "origin/main",
				Ahead:     0,
				Behind:    0,
				IsClean:   true,
				Modified:  []FileChange{},
				Staged:    []FileChange{},
				Untracked: []FileChange{},
				Deleted:   []FileChange{},
			},
		},
		{
			name:  "empty_output",
			input: "",
			want: &Status{
				Branch:    "",
				Remote:    "",
				Ahead:     0,
				Behind:    0,
				IsClean:   true,
				Modified:  []FileChange{},
				Staged:    []FileChange{},
				Untracked: []FileChange{},
				Deleted:   []FileChange{},
			},
		},
		{
			name: "branch_only_no_remote",
			input: `# branch.oid def456
# branch.head feature/test
`,
			want: &Status{
				Branch:    "feature/test",
				Remote:    "",
				Ahead:     0,
				Behind:    0,
				IsClean:   true,
				Modified:  []FileChange{},
				Staged:    []FileChange{},
				Untracked: []FileChange{},
				Deleted:   []FileChange{},
			},
		},
		{
			name: "conflict_marker",
			input: `# branch.head main
u UU N... 100644 100644 100644 100644 hash1 hash2 hash3 conflict_file.go
`,
			want: &Status{
				Branch:       "main",
				IsClean:      false,
				HasConflicts: true,
				Modified:     []FileChange{},
				Staged:       []FileChange{},
				Untracked:    []FileChange{},
				Deleted:      []FileChange{},
			},
		},
		{
			name: "multiple_untracked",
			input: `# branch.head develop
? file1.txt
? file2.txt
? dir/file3.txt
`,
			want: &Status{
				Branch:   "develop",
				IsClean:  false,
				Modified: []FileChange{},
				Staged:   []FileChange{},
				Untracked: []FileChange{
					{Path: "file1.txt", Status: "?"},
					{Path: "file2.txt", Status: "?"},
					{Path: "dir/file3.txt", Status: "?"},
				},
				Deleted: []FileChange{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseStatus(tt.input)

			if got.Branch != tt.want.Branch {
				t.Errorf("Branch = %q, want %q", got.Branch, tt.want.Branch)
			}
			if got.Remote != tt.want.Remote {
				t.Errorf("Remote = %q, want %q", got.Remote, tt.want.Remote)
			}
			if got.Ahead != tt.want.Ahead {
				t.Errorf("Ahead = %d, want %d", got.Ahead, tt.want.Ahead)
			}
			if got.Behind != tt.want.Behind {
				t.Errorf("Behind = %d, want %d", got.Behind, tt.want.Behind)
			}
			if got.IsClean != tt.want.IsClean {
				t.Errorf("IsClean = %v, want %v", got.IsClean, tt.want.IsClean)
			}
			if got.HasConflicts != tt.want.HasConflicts {
				t.Errorf("HasConflicts = %v, want %v", got.HasConflicts, tt.want.HasConflicts)
			}
			if !reflect.DeepEqual(got.Modified, tt.want.Modified) {
				t.Errorf("Modified = %v, want %v", got.Modified, tt.want.Modified)
			}
			if !reflect.DeepEqual(got.Staged, tt.want.Staged) {
				t.Errorf("Staged = %v, want %v", got.Staged, tt.want.Staged)
			}
			if !reflect.DeepEqual(got.Untracked, tt.want.Untracked) {
				t.Errorf("Untracked = %v, want %v", got.Untracked, tt.want.Untracked)
			}
			if !reflect.DeepEqual(got.Deleted, tt.want.Deleted) {
				t.Errorf("Deleted = %v, want %v", got.Deleted, tt.want.Deleted)
			}
		})
	}
}

func TestParseLog(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Commit
	}{
		{
			name: "multiple_commits",
			input: `hash1|shorthash1|Author Name|author@example.com|2 hours ago|Commit message 1
hash2|shorthash2|Author Name|author@example.com|yesterday|Commit message 2`,
			want: []Commit{
				{
					Hash:      "hash1",
					ShortHash: "shorthash1",
					Author:    "Author Name",
					Email:     "author@example.com",
					Date:      "2 hours ago",
					Message:   "Commit message 1",
				},
				{
					Hash:      "hash2",
					ShortHash: "shorthash2",
					Author:    "Author Name",
					Email:     "author@example.com",
					Date:      "yesterday",
					Message:   "Commit message 2",
				},
			},
		},
		{
			name:  "empty_log",
			input: "",
			want:  []Commit{},
		},
		{
			name:  "single_commit",
			input: "abc123|abc|John Doe|john@example.com|just now|Initial commit",
			want: []Commit{
				{
					Hash:      "abc123",
					ShortHash: "abc",
					Author:    "John Doe",
					Email:     "john@example.com",
					Date:      "just now",
					Message:   "Initial commit",
				},
			},
		},
		{
			name:  "malformed_line_skipped",
			input: "not a valid line\nhash|short|Author|email@x.com|now|message",
			want: []Commit{
				{
					Hash:      "hash",
					ShortHash: "short",
					Author:    "Author",
					Email:     "email@x.com",
					Date:      "now",
					Message:   "message",
				},
			},
		},
		{
			name:  "message_with_pipes",
			input: "h|s|A|e@x|now|feat: add X | Y | Z support",
			want: []Commit{
				{
					Hash:      "h",
					ShortHash: "s",
					Author:    "A",
					Email:     "e@x",
					Date:      "now",
					Message:   "feat: add X | Y | Z support",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLog(tt.input)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBranches(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []BranchInfo
	}{
		{
			name: "local_and_remote",
			input: `* master                  abcdef1 Commit message
  feature/new-ui          1234567 Another commit
  remotes/origin/master   abcdef1 Commit message
  remotes/origin/feature  1234567 Another commit`,
			want: []BranchInfo{
				{Name: "master", IsCurrent: true},
				{Name: "feature/new-ui", IsCurrent: false},
				{Name: "origin/master", IsCurrent: false, IsRemote: true},
				{Name: "origin/feature", IsCurrent: false, IsRemote: true},
			},
		},
		{
			name:  "empty_output",
			input: "",
			want:  []BranchInfo{},
		},
		{
			name: "only_local",
			input: `  dev
* main
  feature/test`,
			want: []BranchInfo{
				{Name: "dev", IsCurrent: false},
				{Name: "main", IsCurrent: true},
				{Name: "feature/test", IsCurrent: false},
			},
		},
		{
			name: "only_remotes",
			input: `  remotes/origin/main
  remotes/origin/dev`,
			want: []BranchInfo{
				{Name: "origin/main", IsCurrent: false, IsRemote: true},
				{Name: "origin/dev", IsCurrent: false, IsRemote: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseBranches(tt.input)

			if len(got) != len(tt.want) {
				t.Fatalf("ParseBranches() len = %d, want %d", len(got), len(tt.want))
			}

			for i, w := range tt.want {
				if got[i].Name != w.Name {
					t.Errorf("Branch[%d].Name = %q, want %q", i, got[i].Name, w.Name)
				}
				if got[i].IsCurrent != w.IsCurrent {
					t.Errorf("Branch[%d].IsCurrent = %v, want %v", i, got[i].IsCurrent, w.IsCurrent)
				}
				if got[i].IsRemote != w.IsRemote {
					t.Errorf("Branch[%d].IsRemote = %v, want %v", i, got[i].IsRemote, w.IsRemote)
				}
			}
		})
	}
}

// ==============================================================
// Integration Tests - Real Git Repository
// ==============================================================

// createTestRepo creates a temporary git repository for testing
func createTestRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	// Initialize git repo
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")

	return dir
}

// runGit runs a git command in the specified directory
func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, output)
	}
	return strings.TrimSpace(string(output))
}

// createFile creates a file with content
func createFile(t *testing.T, dir, name, content string) {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}

func TestOpen(t *testing.T) {
	t.Run("valid_repo", func(t *testing.T) {
		dir := createTestRepo(t)

		repo, err := Open(dir)
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
		if repo == nil {
			t.Fatal("Open() returned nil")
		}
		if repo.Path != dir {
			t.Errorf("repo.Path = %q, want %q", repo.Path, dir)
		}
	})

	t.Run("not_a_repo", func(t *testing.T) {
		dir := t.TempDir() // Not initialized

		_, err := Open(dir)
		if err == nil {
			t.Fatal("Open() expected error for non-repo")
		}
		if !strings.Contains(err.Error(), "not a git repository") {
			t.Errorf("error = %q, want 'not a git repository'", err)
		}
	})

	t.Run("nonexistent_path", func(t *testing.T) {
		_, err := Open("/nonexistent/path/12345")
		if err == nil {
			t.Fatal("Open() expected error for nonexistent path")
		}
	})
}

func TestIsRepo(t *testing.T) {
	t.Run("is_repo", func(t *testing.T) {
		dir := createTestRepo(t)
		if !IsRepo(dir) {
			t.Error("IsRepo() = false, want true")
		}
	})

	t.Run("not_repo", func(t *testing.T) {
		dir := t.TempDir()
		if IsRepo(dir) {
			t.Error("IsRepo() = true, want false")
		}
	})
}

func TestRepository_Status_Integration(t *testing.T) {
	dir := createTestRepo(t)

	// Create initial commit
	createFile(t, dir, "initial.txt", "initial content")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "Initial commit")

	repo, err := Open(dir)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	t.Run("clean_repo", func(t *testing.T) {
		status, err := repo.Status()
		if err != nil {
			t.Fatalf("Status() error = %v", err)
		}
		if !status.IsClean {
			t.Error("IsClean = false, want true")
		}
		if len(status.Modified) != 0 {
			t.Errorf("Modified count = %d, want 0", len(status.Modified))
		}
	})

	t.Run("modified_file", func(t *testing.T) {
		createFile(t, dir, "initial.txt", "modified content")

		status, err := repo.Status()
		if err != nil {
			t.Fatalf("Status() error = %v", err)
		}
		if status.IsClean {
			t.Error("IsClean = true, want false")
		}
		if len(status.Modified) != 1 {
			t.Errorf("Modified count = %d, want 1", len(status.Modified))
		}

		// Reset for next test
		runGit(t, dir, "checkout", "--", "initial.txt")
	})

	t.Run("untracked_file", func(t *testing.T) {
		createFile(t, dir, "untracked.txt", "new file")

		status, err := repo.Status()
		if err != nil {
			t.Fatalf("Status() error = %v", err)
		}
		if status.IsClean {
			t.Error("IsClean = true, want false")
		}
		if len(status.Untracked) != 1 {
			t.Errorf("Untracked count = %d, want 1", len(status.Untracked))
		}

		// Cleanup
		os.Remove(filepath.Join(dir, "untracked.txt"))
	})

	t.Run("staged_file", func(t *testing.T) {
		createFile(t, dir, "staged.txt", "staged content")
		runGit(t, dir, "add", "staged.txt")

		status, err := repo.Status()
		if err != nil {
			t.Fatalf("Status() error = %v", err)
		}
		if status.IsClean {
			t.Error("IsClean = true, want false")
		}
		if len(status.Staged) != 1 {
			t.Errorf("Staged count = %d, want 1", len(status.Staged))
		}

		// Cleanup
		runGit(t, dir, "reset", "HEAD", "staged.txt")
		os.Remove(filepath.Join(dir, "staged.txt"))
	})
}

func TestRepository_CurrentBranch_Integration(t *testing.T) {
	dir := createTestRepo(t)

	// Need at least one commit to have a branch
	createFile(t, dir, "test.txt", "content")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "Initial")

	repo, err := Open(dir)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	branch := repo.CurrentBranch()
	// Git 2.28+ defaults to "main", older versions use "master"
	if branch != "main" && branch != "master" {
		t.Errorf("CurrentBranch() = %q, want 'main' or 'master'", branch)
	}

	// Test caching
	firstCall := repo.CurrentBranch()
	secondCall := repo.CurrentBranch()
	if firstCall != secondCall {
		t.Errorf("CurrentBranch() not cached: %q != %q", firstCall, secondCall)
	}
}

func TestRepository_Commit_Integration(t *testing.T) {
	dir := createTestRepo(t)

	repo, err := Open(dir)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	t.Run("commit_staged", func(t *testing.T) {
		createFile(t, dir, "file1.txt", "content1")
		runGit(t, dir, "add", "file1.txt")

		err := repo.Commit("test commit", false)
		if err != nil {
			t.Fatalf("Commit() error = %v", err)
		}

		// Verify commit exists
		log := runGit(t, dir, "log", "--oneline", "-1")
		if !strings.Contains(log, "test commit") {
			t.Errorf("Commit message not found in log: %s", log)
		}
	})

	t.Run("commit_all", func(t *testing.T) {
		createFile(t, dir, "file1.txt", "modified content")

		err := repo.Commit("commit all", true)
		if err != nil {
			t.Fatalf("Commit(all=true) error = %v", err)
		}

		log := runGit(t, dir, "log", "--oneline", "-1")
		if !strings.Contains(log, "commit all") {
			t.Errorf("Commit message not found in log: %s", log)
		}
	})

	t.Run("commit_nothing_staged_fails", func(t *testing.T) {
		err := repo.Commit("should fail", false)
		if err == nil {
			t.Error("Commit() with nothing staged should fail")
		}
	})
}

func TestRepository_Add_Integration(t *testing.T) {
	dir := createTestRepo(t)

	repo, err := Open(dir)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	// Create initial commit
	createFile(t, dir, "initial.txt", "initial")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "Initial")

	t.Run("add_specific_file", func(t *testing.T) {
		createFile(t, dir, "new.txt", "new content")

		err := repo.Add("new.txt")
		if err != nil {
			t.Fatalf("Add() error = %v", err)
		}

		status, _ := repo.Status()
		if len(status.Staged) != 1 {
			t.Errorf("Staged count = %d, want 1", len(status.Staged))
		}

		runGit(t, dir, "reset", "HEAD", "new.txt")
		os.Remove(filepath.Join(dir, "new.txt"))
	})

	t.Run("add_all", func(t *testing.T) {
		createFile(t, dir, "file1.txt", "content1")
		createFile(t, dir, "file2.txt", "content2")

		err := repo.Add()
		if err != nil {
			t.Fatalf("Add() error = %v", err)
		}

		status, _ := repo.Status()
		if len(status.Staged) < 2 {
			t.Errorf("Staged count = %d, want >= 2", len(status.Staged))
		}
	})
}

func TestRepository_ListBranches_Integration(t *testing.T) {
	dir := createTestRepo(t)

	// Create initial commit
	createFile(t, dir, "test.txt", "content")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "Initial")

	// Create additional branches
	runGit(t, dir, "branch", "feature/test")
	runGit(t, dir, "branch", "develop")

	repo, err := Open(dir)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	branches, err := repo.ListBranches()
	if err != nil {
		t.Fatalf("ListBranches() error = %v", err)
	}

	if len(branches) < 3 {
		t.Errorf("Expected at least 3 branches, got %d", len(branches))
	}

	// Check current branch is marked
	hasCurrent := false
	for _, b := range branches {
		if b.IsCurrent {
			hasCurrent = true
			break
		}
	}
	if !hasCurrent {
		t.Error("No current branch marked")
	}
}

func TestRepository_CreateAndSwitchBranch_Integration(t *testing.T) {
	dir := createTestRepo(t)

	// Create initial commit
	createFile(t, dir, "test.txt", "content")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "Initial")

	repo, err := Open(dir)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	err = repo.CreateAndSwitch("feature/new-branch")
	if err != nil {
		t.Fatalf("CreateAndSwitch() error = %v", err)
	}

	// Clear cached branch
	repo.branch = ""

	branch := repo.CurrentBranch()
	if branch != "feature/new-branch" {
		t.Errorf("CurrentBranch() = %q, want 'feature/new-branch'", branch)
	}
}

func TestRepository_Log_Integration(t *testing.T) {
	dir := createTestRepo(t)

	// Create multiple commits
	for i := 1; i <= 3; i++ {
		createFile(t, dir, "test.txt", strings.Repeat("x", i))
		runGit(t, dir, "add", ".")
		runGit(t, dir, "commit", "-m", "Commit "+string(rune('0'+i)))
	}

	repo, err := Open(dir)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	commits, err := repo.Log(2)
	if err != nil {
		t.Fatalf("Log() error = %v", err)
	}

	if len(commits) != 2 {
		t.Errorf("Log(2) returned %d commits, want 2", len(commits))
	}

	// Most recent first
	if !strings.Contains(commits[0].Message, "3") {
		t.Errorf("First commit should be most recent, got: %s", commits[0].Message)
	}
}

func TestRepository_Remote_Integration(t *testing.T) {
	dir := createTestRepo(t)

	repo, err := Open(dir)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	// No remote initially
	remote := repo.Remote()
	if remote != "" {
		t.Errorf("Remote() = %q, want empty", remote)
	}

	// Add a remote
	runGit(t, dir, "remote", "add", "origin", "https://github.com/test/repo.git")

	remote = repo.Remote()
	if remote != "https://github.com/test/repo.git" {
		t.Errorf("Remote() = %q, want 'https://github.com/test/repo.git'", remote)
	}
}
