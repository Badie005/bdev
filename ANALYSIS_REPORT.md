# Codebase Analysis Report

## 1. Code Quality

### Dead Code & Unused Variables
- **Finding**: `IsClean` field in `Status` struct (`internal/core/git/git.go`) is assigned but not effectively used within the package logic other than for assignment.
- **Finding**: `Scan` function in `internal/core/projects/projects.go` has complex logic that might have unreachable paths if not carefully tested.

### Functions Exceeding 50 Lines
- **Finding**: `ParseStatus` in `internal/core/git/git.go` is ~60 lines long.
- **Finding**: `executeStep` in `internal/core/workflow/workflow.go` is ~55 lines long.
- **Finding**: `Execute` in `internal/core/workflow/workflow.go` is ~50 lines.

### Missing Error Handling
- **Critical**: `internal/core/git/git.go` ignores errors from `strconv.Atoi` when parsing ahead/behind status:
  ```go
  status.Ahead, _ = strconv.Atoi(strings.TrimPrefix(parts[0], "+"))
  status.Behind, _ = strconv.Atoi(strings.TrimPrefix(parts[1], "-"))
  ```
  If parsing fails, these default to 0 silently.

### Hardcoded Values
- **Finding**: `internal/cmd/secrets/secrets.go` contains hardcoded string literals for CLI prompts and messages (e.g., "Enter master password: ", "Vault created successfully").
- **Finding**: `internal/core/workflow/workflow.go` has hardcoded "powershell" and "sh" commands.

### Global Mutable State
- **Finding**: `pkg/ui/theme.go` exports `var CurrentTheme` which is mutable. This could lead to race conditions if modified concurrently.
- **Finding**: `pkg/ui/components.go` uses `var ActiveGlyphs` which is also global and mutable.

## 2. Security

### Unsafe File Permissions
- **Finding**: `internal/core/workflow/workflow.go` writes workflow files with `0644` (`os.WriteFile(path, data, 0644)`). If these files contain secrets or sensitive env vars, `0600` would be safer.
- **Finding**: Test files (`git_test.go`, `projects_test.go`) use `0644` or `0755` which is acceptable for tests but should be watched.
- **Positive**: `internal/core/vault` and `internal/ai/memory` correctly use `0600` for sensitive files.

### Secrets & Credentials
- **Positive**: `internal/cmd/secrets` uses `term.ReadPassword` for secure input.
- **Risk**: Environment variable expansion in `internal/core/workflow` (`expandEnv`) does simple string replacement and might leak secrets if not carefully handled (though it attempts to mask them).

### Unsafe HTTP Calls
- **Positive**: `internal/ai/engine/client.go` uses `http.NewRequestWithContext`, preventing orphaned requests.

## 3. Testing

### Low Coverage (< 50%)
- **Critical**: `internal/ai/agents` (0.0%)
- **Critical**: `internal/ai/engine` (0.0%)
- **Critical**: `internal/cmd/...` (0.0%) - CLI commands are untested.
- **Critical**: `internal/core/workflow` (0.0%)
- **Critical**: `pkg/ui` (0.0%) - UI logic untested.
- **Warning**: `internal/core/git` (42.1%)

### Missing Tests for Public Functions
- `pkg/ui` has 44 exported functions but 0% coverage.
- `internal/core/workflow` has complex logic in `Execute` and `expandEnv` but no tests.

## 4. Go Best Practices

### Missing Godoc
- **Finding**: `pkg/ui/theme.go` helper functions (`Primary`, `Secondary`, etc.) lack individual comments.
- **Finding**: `internal/cmd` packages generally lack documentation for their `NewCommand` functions.

### Context Usage
- **Finding**: `internal/core/git/git.go` methods (`Status`, `Log`, `Commit`) execute `exec.Command` without `context.Context`. This prevents cancellation and timeouts for long-running git operations.
- **Finding**: `internal/core/workflow/workflow.go` `Execute` does not accept a `context.Context`, making workflow execution uncancellable.

## 5. Performance

### Allocations in Loops
- **Finding**: `internal/core/git/git.go` in `ParseStatus` appends to `status.Staged`, `status.Modified` etc. inside a loop without preallocating.
  ```go
  status.Staged = append(status.Staged, change)
  ```
- **Finding**: `internal/core/workflow/workflow.go` appends to `result.Steps` in a loop.

## 6. CI/CD

### Workflows
- **Verified**: `.github/workflows/ci.yml` correctly disables caching for Windows (`cache: ${{ matrix.os != 'windows-latest' }}`), adhering to known issues.
- **Configuration**: `gosec` is configured to exclude `G306`, which allows `0644` permissions. This explains why the workflow file permission issue wasn't caught by CI.

## Recommendations

1.  **Fix Error Handling**: Handle `strconv.Atoi` errors in `git.go`.
2.  **Add Context**: Refactor `internal/core/git` to use `exec.CommandContext`.
3.  **Improve Testing**: Add tests for `internal/core/workflow` and `pkg/ui`.
4.  **Secure Permissions**: Change workflow file permissions to `0600`.
5.  **Refactor Globals**: Make `CurrentTheme` thread-safe or immutable.
