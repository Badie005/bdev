package workflow

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Workflow represents a workflow definition
type Workflow struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
	Steps       []Step            `yaml:"steps"`
	OnSuccess   []Step            `yaml:"on_success,omitempty"`
	OnFailure   []Step            `yaml:"on_failure,omitempty"`
}

// Step represents a workflow step
type Step struct {
	Name     string            `yaml:"name"`
	Run      string            `yaml:"run"`
	Cwd      string            `yaml:"cwd,omitempty"`
	Env      map[string]string `yaml:"env,omitempty"`
	If       string            `yaml:"if,omitempty"`
	Continue bool              `yaml:"continue_on_error,omitempty"`
	Timeout  string            `yaml:"timeout,omitempty"`
}

// StepResult holds the result of a step execution
type StepResult struct {
	Step     Step
	Success  bool
	Output   string
	Error    error
	Duration time.Duration
}

// WorkflowResult holds the result of a workflow execution
type WorkflowResult struct {
	Workflow  *Workflow
	Steps     []StepResult
	Success   bool
	Duration  time.Duration
	StartTime time.Time
}

// Engine executes workflows
type Engine struct {
	WorkflowDir string
	Env         map[string]string
	Verbose     bool
	Vault       interface { // Interface to avoid strict dependency on specific Vault implementation details if unused
		Get(key string) (string, error)
		IsUnlocked() bool
	}
	OnStep func(step Step, result *StepResult)
}

// New creates a new workflow engine
func New(workflowDir string) *Engine {
	return &Engine{
		WorkflowDir: workflowDir,
		Env:         make(map[string]string),
	}
}

// List returns all available workflows
func (e *Engine) List() ([]string, error) {
	entries, err := os.ReadDir(e.WorkflowDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	workflows := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			workflows = append(workflows, strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml"))
		}
	}

	return workflows, nil
}

// Load loads a workflow by name
func (e *Engine) Load(name string) (*Workflow, error) {
	// Try .yaml first, then .yml
	path := filepath.Join(e.WorkflowDir, name+".yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = filepath.Join(e.WorkflowDir, name+".yml")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("workflow not found: %s", name)
	}

	var wf Workflow
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("invalid workflow: %w", err)
	}

	if wf.Name == "" {
		wf.Name = name
	}

	return &wf, nil
}

// Execute runs a workflow
func (e *Engine) Execute(wf *Workflow) *WorkflowResult {
	result := &WorkflowResult{
		Workflow:  wf,
		Steps:     make([]StepResult, 0),
		Success:   true,
		StartTime: time.Now(),
	}

	// Merge workflow env with engine env
	env := make(map[string]string)
	for k, v := range e.Env {
		env[k] = v
	}
	for k, v := range wf.Env {
		env[k] = v
	}

	// Execute steps
	for _, step := range wf.Steps {
		stepResult := e.executeStep(step, env)
		result.Steps = append(result.Steps, stepResult)

		if e.OnStep != nil {
			e.OnStep(step, &stepResult)
		}

		if !stepResult.Success && !step.Continue {
			result.Success = false
			break
		}
	}

	// Execute on_success or on_failure hooks
	if result.Success && len(wf.OnSuccess) > 0 {
		for _, step := range wf.OnSuccess {
			stepResult := e.executeStep(step, env)
			result.Steps = append(result.Steps, stepResult)
		}
	} else if !result.Success && len(wf.OnFailure) > 0 {
		for _, step := range wf.OnFailure {
			stepResult := e.executeStep(step, env)
			result.Steps = append(result.Steps, stepResult)
		}
	}

	result.Duration = time.Since(result.StartTime)
	return result
}

// executeStep runs a single step
func (e *Engine) executeStep(step Step, env map[string]string) StepResult {
	start := time.Now()
	result := StepResult{Step: step}

	// Merge step env
	stepEnv := make(map[string]string)
	for k, v := range env {
		stepEnv[k] = v
	}
	for k, v := range step.Env {
		stepEnv[k] = v
	}

	// Expand environment variables in command
	cmdLine := e.expandEnv(step.Run, stepEnv)

	// Create command (use shell)
	var cmd *exec.Cmd
	if isWindows() {
		cmd = exec.Command("powershell", "-Command", cmdLine)
	} else {
		cmd = exec.Command("sh", "-c", cmdLine)
	}

	// Set working directory
	if step.Cwd != "" {
		cmd.Dir = e.expandEnv(step.Cwd, stepEnv)
	}

	// Set environment
	cmd.Env = os.Environ()
	for k, v := range stepEnv {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Run and capture output
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	result.Duration = time.Since(start)

	if err != nil {
		result.Error = err
		result.Success = false
	} else {
		result.Success = true
	}

	return result
}

// expandEnv expands environment variables and secrets in a string
func (e *Engine) expandEnv(s string, env map[string]string) string {
	result := s

	// 1. Context variables ${{ env.VAR }} or ${{ secrets.KEY }}
	// Simple replacement for now, regex ideal but simple substitution works for MVP

	// Helper to replace secrets
	if strings.Contains(result, "secrets.") && e.Vault != nil && e.Vault.IsUnlocked() {
		// This is a naive implementation. For robustness, regex should be used to find all ${{ secrets.KEY }}
		// But let's try to find keys that are likely requested.
		// Actually, we can't iterate secrets if we don't know keys.
		// Valid Workflow syntax: ${{ secrets.MY_KEY }}
		// We can scan the string for this pattern.

		start := 0
		for {
			idx := strings.Index(result[start:], "${{ secrets.")
			if idx == -1 {
				break
			}
			idx += start
			end := strings.Index(result[idx:], " }}")
			if end == -1 {
				break
			}
			end += idx

			// Extract key: ${{ secrets.KEY }} -> KEY
			keyRaw := result[idx+len("${{ secrets.") : end]
			key := strings.TrimSpace(keyRaw)

			secretVal, err := e.Vault.Get(key)
			if err == nil {
				// Replace substring
				result = result[:idx] + secretVal + result[end+3:]
				// Adjust start to avoid infinite loop if replacement contains pattern (unlikely for secrets)
				// But we should restart search or move past? moving past is safer.
				start = idx + len(secretVal)
			} else {
				// Warn? For now keep placeholder or empty?
				// Common behavior matches empty string if missing?
				// Let's leave it to show error or just skip.
				start = end + 3
			}
		}
	}

	// 2. Standard Env vars
	for k, v := range env {
		result = strings.ReplaceAll(result, "$"+k, v)
		result = strings.ReplaceAll(result, "${"+k+"}", v)
		result = strings.ReplaceAll(result, "${{ "+k+" }}", v)
		result = strings.ReplaceAll(result, "${{ env."+k+" }}", v)
	}

	return os.ExpandEnv(result)
}

// isWindows checks if running on Windows
func isWindows() bool {
	return os.PathSeparator == '\\'
}

// Save saves a workflow to disk
func (e *Engine) Save(wf *Workflow) error {
	data, err := yaml.Marshal(wf)
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(e.WorkflowDir, 0755); err != nil {
		return err
	}

	path := filepath.Join(e.WorkflowDir, wf.Name+".yaml")
	return os.WriteFile(path, data, 0644)
}

// Delete removes a workflow
func (e *Engine) Delete(name string) error {
	path := filepath.Join(e.WorkflowDir, name+".yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = filepath.Join(e.WorkflowDir, name+".yml")
	}
	return os.Remove(path)
}
