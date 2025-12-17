# Codebase Analysis Report

## Code Quality

### Dead Code and Unused Variables
- **Detection**: Manual review and tool check (partial).
- **Findings**:
  - `internal/core/repl/completer.go`: `loadProjects` (unused locally, might be used by `RefreshProjects`).
  - `internal/core/git/mock.go`: Many methods of `MockGit` seem unused if tests don't cover them.
  - `internal/core/multi/multi.go`: `ExecuteShell` (unused?).
- **Recommendation**: Run `staticcheck` in CI to automatically detect these.

### Functions Exceeding 50 Lines
- **Findings**:
  - `internal/cmd/git/git.go`: `NewCommand` and subcommands likely contribute to high file line count (476 lines).
  - `internal/core/git/git.go`: `ParseStatus` (74 lines), `ParseLog` (might be long).
  - `internal/core/projects/projects.go`: `Detect` (60 lines), `Scan` (50+ lines).
  - `internal/ai/engine/client.go`: `streamChat` (60+ lines), `streamGenerate` (60+ lines).
  - `internal/core/workflow/workflow.go`: `Execute` (60+ lines), `executeStep` (60+ lines).
- **Recommendation**: Refactor long functions into smaller, testable helpers. For example, `ParseStatus` could be split into `parseBranchInfo` and `parseFileChanges`.

### Missing Error Handling
- **Findings**:
  - `internal/core/git/git_test.go`: `t.Errorf` is used, which allows execution to continue. In some cases (e.g. setup), `t.Fatalf` should be used to prevent panics.
  - `internal/core/workflow/workflow.go`: `e.executeStep` calls `cmd.CombinedOutput()` and sets `result.Error` but the caller `Execute` continues unless `!step.Continue`. This is by design, but ensure errors are logged if ignored.
- **Recommendation**: Use `t.Fatalf` for test setup failures. Review `workflow.go` to ensure errors aren't silently swallowed when `continue-on-error` is false.

### Hardcoded Values
- **Findings**:
  - Permissions: `0o600`, `0o644`, `0700`, `0755` scattered across `internal/core/config`, `internal/core/vault`, `internal/core/workflow`, `internal/core/session`.
  - `internal/ai/engine/client.go`: Defaults `http://localhost:11434`, `llama3.2`.
  - `pkg/ui/theme.go`: Hex colors.
- **Recommendation**: define constants for permissions in a shared package (e.g., `pkg/consts` or `internal/core/fs`). Move AI defaults to configuration.

### Global Mutable State
- **Findings**:
  - `pkg/ui/theme.go`: `var CurrentTheme = AnthropicTheme()`.
  - `internal/cmd/secrets/secrets.go`: `var v *vault.Vault`.
  - `pkg/ui/components.go`: `ActiveGlyphs` (global var).
- **Recommendation**: Avoid global state where possible. Pass `Theme` or `Vault` as dependencies to functions or structs.

## Security

### Unsafe File Permissions
- **Findings**:
  - `internal/core/workflow/workflow.go`: `os.WriteFile(path, data, 0644)`. Workflows might contain secrets or sensitive logic. `0600` is safer if they are user-specific.
  - `internal/core/git/git_test.go`: `0o644` (Test files, acceptable).
  - `internal/core/vault/vault.go`: `0600` (Safe).
- **Recommendation**: Change `internal/core/workflow/workflow.go` to use `0600` for workflow files if they are not meant to be shared globally.

### Missing Input Validation
- **Findings**:
  - `internal/cmd/secrets/secrets.go`: Basic password length check exists.
  - `internal/core/projects/projects.go`: `Analyze` trusts the file system structure.
  - Command arguments in `internal/cmd` should be validated more rigorously (e.g., checking for valid project names/paths).

### Secrets/Credentials in Code
- **Findings**:
  - None explicitly found in source code.
  - `internal/core/workflow/workflow.go`: Implements secret expansion `${{ secrets.KEY }}`, which is good practice (externalizing secrets).

### Unsafe HTTP Calls
- **Findings**:
  - `internal/ai/engine/client.go`: Uses `http.NewRequestWithContext`. **Safe**.
  - `internal/ai/engine/client.go`: `Client` struct has a configured `HTTPClient` with timeout. **Safe**.

## Testing

### Packages with <50% Test Coverage
- **Findings**:
  - `internal/core/git` (Low coverage, mostly mock).
  - `internal/core/multi` (0%).
  - `internal/core/repl` (0%).
  - `internal/core/session` (0%).
  - `internal/core/workflow` (0%).
  - `internal/cmd/*` (Mostly 0%).
  - `pkg/ui` (0%).
  - `internal/ai/agents` (0%).
  - `internal/ai/memory` (0%).
- **Recommendation**: Prioritize adding tests for `pkg/ui` (components), `internal/core/workflow`, and `internal/ai`.

### Missing Tests for Public Functions
- **Findings**:
  - `pkg/ui/components.go`: `Render`, `MessageSuccess`, etc.
  - `internal/core/workflow/workflow.go`: `Execute`, `Load`, `Save`.
  - `internal/ai/engine/client.go`: `Chat`, `Generate`.
- **Recommendation**: Add unit tests for these critical functions.

### Test Files without Table-Driven Tests
- **Findings**:
  - `internal/core/runner/runner_test.go`: Partially table-driven, but uses `t.Run("echo_command", ...)` individually.
  - `internal/core/git/git_test.go`: Uses `t.Run` for distinct scenarios but could be more table-driven for `ParseStatus` (it is table-driven now, good).
- **Recommendation**: Refactor `runner_test.go` to fully utilize table-driven tests.

## Go Best Practices

### Missing Godoc Comments
- **Findings**:
  - `internal/cmd/ai/ai.go`: Missing comments for command constructors.
  - `internal/ai/agents/agents.go`: Missing comments for `Agent` interface methods.
  - `pkg/ui/components.go`: Some comments exist but could be more comprehensive.
- **Recommendation**: Add godoc comments for all exported functions and types.

### Incorrect Error Message Capitalization
- **Findings**:
  - `internal/cmd/ai/ai.go`: `fmt.Errorf("Ollama is not running...")` -> Capitalized.
  - `internal/cmd/secrets/secrets.go`: `fmt.Errorf("Vault created successfully")` (This is a message, not an error return, but uses fmt.Errorf? No, wait, it was `fmt.Println` in my read). `fmt.Errorf("passwords do not match")` is Lowercase (Correct).
- **Recommendation**: Ensure all error strings returned by `fmt.Errorf` or `errors.New` start with a lowercase letter and do not end with punctuation (style guide).

### Missing Context.Context
- **Findings**:
  - `internal/core/multi/multi.go`: `Execute` takes `context.Context`. **Good**.
  - `internal/ai/engine/client.go`: methods take `context.Context`. **Good**.
  - `internal/core/git/git.go`: `run` uses `exec.Command`, which does not take context. It should use `exec.CommandContext` to allow cancellation/timeouts.
- **Recommendation**: Update `internal/core/git/git.go` to use `exec.CommandContext`.

## Performance

### Unnecessary Allocations
- **Findings**:
  - `internal/core/git/git.go`: `status.Staged = append(status.Staged, change)`. Repeated appending without preallocation.
  - `internal/core/repl/completer.go`: `c.projects = make([]string, 0)`.
- **Recommendation**: Use `make([]T, 0, estimatedSize)` where possible. For git status, we don't know the size upfront easily, but for `completer.go` we might.

## CI/CD

### Failing GitHub Actions Workflows
- **Findings**:
  - `.github/workflows/ci.yml`:
    - `golangci-lint` action is used.
    - Tests run on Ubuntu, Windows, MacOS.
    - `cache: true` is used for `setup-go` on Windows. **Issue detected**: The prompt memory says "Windows CI workflows ... must have caching disabled ... to avoid tar/cache errors".
- **Recommendation**: Disable caching for Windows in `ci.yml`.

### Missing or Outdated Dependencies
- **Findings**:
  - `golang.org/x/crypto v0.16.0` (Somewhat old).
  - `spf13/cobra v1.8.0` (Recent).
- **Recommendation**: Run `go get -u ./...` to update dependencies.

## Summary of Critical Issues
1.  **CI/CD**: Fix Windows caching issue in `ci.yml`.
2.  **Security**: Change permissions in `internal/core/workflow/workflow.go` to `0600`.
3.  **Testing**: Significantly increase coverage, especially in `internal/core/workflow` and `pkg/ui`.
4.  **Best Practices**: Add context support to `internal/core/git` for timeout control.
