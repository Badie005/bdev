# Codebase Analysis Report

## 1. Code Quality

### Functions Exceeding 50 Lines
The following functions exceed the recommended 50-line limit, making them harder to maintain and test:

*   **`internal/cmd/workflow/workflow.go`**: `runCmd` (110 lines) - *Refactoring recommended*
*   **`internal/cmd/root/root.go`**: `registerShortcuts` (92 lines)
*   **`internal/ai/agents/agents.go`**: `runAgentOnFile` (79 lines)
*   **`internal/core/git/git.go`**: `ParseStatus` (74 lines)
*   **`internal/cmd/ai/ai.go`**: `chatCmd` (75 lines)
*   **`internal/cmd/git/git.go`**: `statusCmd` (65 lines), `branchCmd` (63 lines), `stashCmd` (56 lines)
*   **`internal/cmd/multi/multi.go`**: `execCmd` (58 lines), `gitCmd` (65 lines)
*   **`internal/cmd/projects/projects.go`**: `runCmd` (56 lines)

### Global Mutable State
Several global variables were found which can lead to race conditions and testing issues:

*   **`pkg/ui/theme.go`**: `ActiveGlyphs`, `FallbackGlyphs` (Exported mutable structs)
*   **`pkg/ui/components.go`**: `ProgressBuild` (Exported mutable slice)
*   **`internal/cmd/root/root.go`**: `rootCmd` (Package-level variable)
*   **`internal/ai/agents/agents.go`**: `ReviewerAgent`, `ExplainerAgent`, etc. are global variables.

### Missing Error Handling & Conventions
*   **Error Message Capitalization**: several error messages start with a capital letter, violating Go style guide (ST1005).
    *   `internal/ai/agents/agents.go`: "Ollama is not running..."
    *   `internal/cmd/secrets/secrets.go`: "Vault already exists..."
    *   `internal/core/vault/vault.go`: "Wrong password or..."

### Hardcoded Values
*   **`internal/ai/engine/client.go`**: Hardcoded timeouts (`2*time.Second`, `10*time.Second`) should be configurable or constants.
*   **`internal/cmd/secrets/secrets.go`**: Hardcoded password length requirement (8).

## 2. Security

### Unsafe File Permissions
The following files are written with permissions that may be too open (not `0o600` or `0o755`):

*   **`internal/core/workflow/workflow.go`**: `os.WriteFile(path, data, 0644)`
    *   Workflows may contain sensitive environment variables. Permission `0644` (readable by all users) is risky. Recommended: `0o600`.
*   **`internal/core/git/git_test.go`**: `0o644` (Test file)
*   **`internal/core/projects/projects_test.go`**: `0o644` (Test file)

### Secrets in Code
*   No hardcoded secrets (API keys, passwords) were found in the source code.
*   Secrets management in `internal/core/vault` uses AES-256-GCM, which is secure.

### Input Validation
*   `internal/cmd/secrets/secrets.go` validates password length, but complexity checks are missing.

## 3. Testing

### Low Coverage Packages (< 50%)
The following packages have critical gaps in test coverage:

*   **`internal/ai/agents`**: 0.0%
*   **`internal/ai/engine`**: 0.0%
*   **`internal/ai/memory`**: 0.0%
*   **`internal/cmd/...`**: 0.0% (All CLI commands)
*   **`internal/core/multi`**: 0.0%
*   **`internal/core/repl`**: 0.0%
*   **`internal/core/session`**: 0.0%
*   **`internal/core/workflow`**: 0.0%
*   **`pkg/ui`**: 0.0%
*   **`internal/core/git`**: 42.1%

### Missing Tests
*   Public functions in `pkg/ui` are untested.
*   `internal/core/workflow` engine logic is completely untested.

## 4. Performance

### Allocations
*   **`internal/core/workflow/workflow.go`**: In `List()`, `workflows` slice is appended to in a loop without pre-allocation.
    *   Current: `workflows := make([]string, 0)`
    *   Recommended: `workflows := make([]string, 0, len(entries))`

## 5. CI/CD

### Dependencies
*   **`gopkg.in/yaml.v3`**: Used correctly, but caused minor linter confusion in some environments.
*   **Build**: Project builds successfully with Go 1.24.

## Recommendations
1.  **Restrict Permissions**: Change `os.WriteFile` in `workflow.go` to use `0o600`.
2.  **Increase Coverage**: Add unit tests for `internal/core/workflow` and `pkg/ui`.
3.  **Refactor**: Break down `runCmd` in `workflow.go` and `registerShortcuts` in `root.go`.
4.  **Fix Constants**: Move timeouts and magic numbers to `const` blocks.
