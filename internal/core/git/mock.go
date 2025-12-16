package git

// MockRepository is a test mock implementation of the Repo interface
type MockRepository struct {
	// Configurable return values
	StatusFunc          func() (*Status, error)
	CurrentBranchFunc   func() string
	LogFunc             func(count int) ([]Commit, error)
	LogGraphFunc        func(count int) (string, error)
	ListBranchesFunc    func() ([]BranchInfo, error)
	DiffFunc            func(staged bool) (string, error)
	DiffStatFunc        func(staged bool) (string, error)
	RemoteFunc          func() string
	LastCommitFunc      func() (*Commit, error)
	AddFunc             func(files ...string) error
	CommitFunc          func(message string, all bool) error
	PushFunc            func(force bool) error
	PullFunc            func() error
	FetchFunc           func() error
	CreateBranchFunc    func(name string) error
	SwitchBranchFunc    func(name string) error
	CreateAndSwitchFunc func(name string) error
	DeleteBranchFunc    func(name string, force bool) error
	StashFunc           func(message string) error
	StashPopFunc        func() error
	StashListFunc       func() ([]string, error)
	StashDropFunc       func(index int) error
	ResetFunc           func(hard bool, ref string) error
	ResetFileFunc       func(file string) error

	// Call tracking
	Calls []string
}

// Ensure MockRepository implements Repo interface
var _ Repo = (*MockRepository)(nil)

func (m *MockRepository) Status() (*Status, error) {
	m.Calls = append(m.Calls, "Status")
	if m.StatusFunc != nil {
		return m.StatusFunc()
	}
	return &Status{IsClean: true, Branch: "main"}, nil
}

func (m *MockRepository) CurrentBranch() string {
	m.Calls = append(m.Calls, "CurrentBranch")
	if m.CurrentBranchFunc != nil {
		return m.CurrentBranchFunc()
	}
	return "main"
}

func (m *MockRepository) Log(count int) ([]Commit, error) {
	m.Calls = append(m.Calls, "Log")
	if m.LogFunc != nil {
		return m.LogFunc(count)
	}
	return []Commit{}, nil
}

func (m *MockRepository) LogGraph(count int) (string, error) {
	m.Calls = append(m.Calls, "LogGraph")
	if m.LogGraphFunc != nil {
		return m.LogGraphFunc(count)
	}
	return "", nil
}

func (m *MockRepository) ListBranches() ([]BranchInfo, error) {
	m.Calls = append(m.Calls, "ListBranches")
	if m.ListBranchesFunc != nil {
		return m.ListBranchesFunc()
	}
	return []BranchInfo{{Name: "main", IsCurrent: true}}, nil
}

func (m *MockRepository) Diff(staged bool) (string, error) {
	m.Calls = append(m.Calls, "Diff")
	if m.DiffFunc != nil {
		return m.DiffFunc(staged)
	}
	return "", nil
}

func (m *MockRepository) DiffStat(staged bool) (string, error) {
	m.Calls = append(m.Calls, "DiffStat")
	if m.DiffStatFunc != nil {
		return m.DiffStatFunc(staged)
	}
	return "", nil
}

func (m *MockRepository) Remote() string {
	m.Calls = append(m.Calls, "Remote")
	if m.RemoteFunc != nil {
		return m.RemoteFunc()
	}
	return "https://github.com/test/repo.git"
}

func (m *MockRepository) LastCommit() (*Commit, error) {
	m.Calls = append(m.Calls, "LastCommit")
	if m.LastCommitFunc != nil {
		return m.LastCommitFunc()
	}
	return &Commit{Hash: "abc123", Message: "test commit"}, nil
}

func (m *MockRepository) Add(files ...string) error {
	m.Calls = append(m.Calls, "Add")
	if m.AddFunc != nil {
		return m.AddFunc(files...)
	}
	return nil
}

func (m *MockRepository) Commit(message string, all bool) error {
	m.Calls = append(m.Calls, "Commit")
	if m.CommitFunc != nil {
		return m.CommitFunc(message, all)
	}
	return nil
}

func (m *MockRepository) Push(force bool) error {
	m.Calls = append(m.Calls, "Push")
	if m.PushFunc != nil {
		return m.PushFunc(force)
	}
	return nil
}

func (m *MockRepository) Pull() error {
	m.Calls = append(m.Calls, "Pull")
	if m.PullFunc != nil {
		return m.PullFunc()
	}
	return nil
}

func (m *MockRepository) Fetch() error {
	m.Calls = append(m.Calls, "Fetch")
	if m.FetchFunc != nil {
		return m.FetchFunc()
	}
	return nil
}

func (m *MockRepository) CreateBranch(name string) error {
	m.Calls = append(m.Calls, "CreateBranch")
	if m.CreateBranchFunc != nil {
		return m.CreateBranchFunc(name)
	}
	return nil
}

func (m *MockRepository) SwitchBranch(name string) error {
	m.Calls = append(m.Calls, "SwitchBranch")
	if m.SwitchBranchFunc != nil {
		return m.SwitchBranchFunc(name)
	}
	return nil
}

func (m *MockRepository) CreateAndSwitch(name string) error {
	m.Calls = append(m.Calls, "CreateAndSwitch")
	if m.CreateAndSwitchFunc != nil {
		return m.CreateAndSwitchFunc(name)
	}
	return nil
}

func (m *MockRepository) DeleteBranch(name string, force bool) error {
	m.Calls = append(m.Calls, "DeleteBranch")
	if m.DeleteBranchFunc != nil {
		return m.DeleteBranchFunc(name, force)
	}
	return nil
}

func (m *MockRepository) Stash(message string) error {
	m.Calls = append(m.Calls, "Stash")
	if m.StashFunc != nil {
		return m.StashFunc(message)
	}
	return nil
}

func (m *MockRepository) StashPop() error {
	m.Calls = append(m.Calls, "StashPop")
	if m.StashPopFunc != nil {
		return m.StashPopFunc()
	}
	return nil
}

func (m *MockRepository) StashList() ([]string, error) {
	m.Calls = append(m.Calls, "StashList")
	if m.StashListFunc != nil {
		return m.StashListFunc()
	}
	return []string{}, nil
}

func (m *MockRepository) StashDrop(index int) error {
	m.Calls = append(m.Calls, "StashDrop")
	if m.StashDropFunc != nil {
		return m.StashDropFunc(index)
	}
	return nil
}

func (m *MockRepository) Reset(hard bool, ref string) error {
	m.Calls = append(m.Calls, "Reset")
	if m.ResetFunc != nil {
		return m.ResetFunc(hard, ref)
	}
	return nil
}

func (m *MockRepository) ResetFile(file string) error {
	m.Calls = append(m.Calls, "ResetFile")
	if m.ResetFileFunc != nil {
		return m.ResetFileFunc(file)
	}
	return nil
}
