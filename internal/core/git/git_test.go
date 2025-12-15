package git

import (
	"reflect"
	"testing"
)

func TestParseStatus(t *testing.T) {
	// Sample output from 'git status --porcelain=v2 --branch'
	output := `# branch.oid 0000000000000000000000000000000000000000
# branch.head master
# branch.upstream origin/master
# branch.ab +1 -2
1 .M N... 100644 100644 100644 9daeafb9864cf43055ae93beb0afd6c7d144bfa4 9daeafb9864cf43055ae93beb0afd6c7d144bfa4 modified_file.go
1 M. N... 100644 100644 100644 9daeafb9864cf43055ae93beb0afd6c7d144bfa4 9daeafb9864cf43055ae93beb0afd6c7d144bfa4 staged_file.go
? untracked_file.txt
`

	want := &Status{
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
		Deleted: make([]FileChange, 0),
	}

	got := ParseStatus(output)

	if got.Branch != want.Branch {
		t.Errorf("Branch = %q, want %q", got.Branch, want.Branch)
	}
	if got.Remote != want.Remote {
		t.Errorf("Remote = %q, want %q", got.Remote, want.Remote)
	}
	if got.Ahead != want.Ahead {
		t.Errorf("Ahead = %d, want %d", got.Ahead, want.Ahead)
	}
	if got.Behind != want.Behind {
		t.Errorf("Behind = %d, want %d", got.Behind, want.Behind)
	}
	if !reflect.DeepEqual(got.Modified, want.Modified) {
		t.Errorf("Modified = %v, want %v", got.Modified, want.Modified)
	}
	if !reflect.DeepEqual(got.Staged, want.Staged) {
		t.Errorf("Staged = %v, want %v", got.Staged, want.Staged)
	}
	if !reflect.DeepEqual(got.Untracked, want.Untracked) {
		t.Errorf("Untracked = %v, want %v", got.Untracked, want.Untracked)
	}
	if got.IsClean != want.IsClean {
		t.Errorf("IsClean = %v, want %v", got.IsClean, want.IsClean)
	}
	if !reflect.DeepEqual(got.Deleted, want.Deleted) {
		t.Errorf("Deleted = %v, want %v", got.Deleted, want.Deleted)
	}
}

func TestParseLog(t *testing.T) {
	output := `hash1|shorthash1|Author Name|author@example.com|2 hours ago|Commit message 1
hash2|shorthash2|Author Name|author@example.com|yesterday|Commit message 2`

	want := []Commit{
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
	}

	got := ParseLog(output)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseLog() = %v, want %v", got, want)
	}
}

func TestParseBranches(t *testing.T) {
	output := `* master                  abcdef1 Commit message
  feature/new-ui          1234567 Another commit
  remotes/origin/master   abcdef1 Commit message
  remotes/origin/feature  1234567 Another commit`

	want := []BranchInfo{
		{Name: "master", IsCurrent: true},
		{Name: "feature/new-ui", IsCurrent: false},
		{Name: "origin/master", IsCurrent: false, IsRemote: true},
		{Name: "origin/feature", IsCurrent: false, IsRemote: true},
	}

	got := ParseBranches(output)

	if len(got) != len(want) {
		t.Fatalf("ParseBranches() len = %d, want %d", len(got), len(want))
	}

	for i, w := range want {
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
}
