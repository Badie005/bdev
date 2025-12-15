"""
B.DEV CLI - Workflow Engine
YAML-based automation like GitHub Actions
"""
import yaml
import subprocess
from pathlib import Path
from typing import Dict, List, Optional, Any
from dataclasses import dataclass, field
from datetime import datetime

from utils.ui import console, print_header, print_success, print_error, print_warning

WORKFLOWS_DIR = Path.home() / "Dev" / ".bdev" / "workflows"

@dataclass
class WorkflowStep:
    name: str
    run: str
    continue_on_error: bool = False
    env: Dict[str, str] = field(default_factory=dict)
    condition: Optional[str] = None

@dataclass
class Workflow:
    name: str
    description: str = ""
    steps: List[WorkflowStep] = field(default_factory=list)
    env: Dict[str, str] = field(default_factory=dict)
    on_success: Optional[str] = None
    on_failure: Optional[str] = None

class WorkflowEngine:
    """Execute YAML workflows"""
    
    def __init__(self):
        WORKFLOWS_DIR.mkdir(parents=True, exist_ok=True)
    
    def list_workflows(self) -> List[str]:
        """List available workflows"""
        return [f.stem for f in WORKFLOWS_DIR.glob("*.yml")]
    
    def load(self, name: str) -> Optional[Workflow]:
        """Load a workflow from YAML"""
        path = WORKFLOWS_DIR / f"{name}.yml"
        if not path.exists():
            return None
        
        try:
            data = yaml.safe_load(path.read_text(encoding='utf-8'))
            
            steps = []
            for step_data in data.get("steps", []):
                if isinstance(step_data, str):
                    steps.append(WorkflowStep(name=step_data, run=step_data))
                elif isinstance(step_data, dict):
                    steps.append(WorkflowStep(
                        name=step_data.get("name", step_data.get("run", "step")),
                        run=step_data.get("run", ""),
                        continue_on_error=step_data.get("continue_on_error", False),
                        env=step_data.get("env", {}),
                        condition=step_data.get("if")
                    ))
            
            return Workflow(
                name=data.get("name", name),
                description=data.get("description", ""),
                steps=steps,
                env=data.get("env", {}),
                on_success=data.get("on_success"),
                on_failure=data.get("on_failure")
            )
        except Exception as e:
            print_error(f"Erreur parsing workflow: {e}")
            return None
    
    def run(self, name: str, cwd: Path = None) -> bool:
        """Execute a workflow"""
        from utils.theme import get_theme_manager
        c = get_theme_manager().current_theme.colors
        
        workflow = self.load(name)
        if not workflow:
            print_error(f"Workflow '{name}' non trouvé")
            return False
        
        print_header(f"Workflow: {workflow.name}", workflow.description)
        
        start_time = datetime.now()
        success_count = 0
        fail_count = 0
        
        for i, step in enumerate(workflow.steps, 1):
            console.print(f"\n[{c.accent}]Step {i}/{len(workflow.steps)}[/] {step.name}")
            console.print(f"[dim]$ {step.run}[/dim]")
            
            # Merge env
            env = {**workflow.env, **step.env}
            
            try:
                result = subprocess.run(
                    step.run,
                    shell=True,
                    cwd=cwd or Path.cwd(),
                    env={**dict(__import__('os').environ), **env},
                    capture_output=False
                )
                
                if result.returncode == 0:
                    console.print(f"[{c.success}]✔ Passed[/]")
                    success_count += 1
                else:
                    console.print(f"[{c.error}]✖ Failed (code {result.returncode})[/]")
                    fail_count += 1
                    if not step.continue_on_error:
                        if workflow.on_failure:
                            subprocess.run(workflow.on_failure, shell=True)
                        return False
                        
            except Exception as e:
                console.print(f"[{c.error}]✖ Error: {e}[/]")
                fail_count += 1
                if not step.continue_on_error:
                    return False
        
        duration = (datetime.now() - start_time).total_seconds()
        console.print(f"\n[{c.border}]{'─' * 40}[/]")
        console.print(f"[{c.success}]✔ {success_count}[/] passed  [{c.error}]✖ {fail_count}[/] failed  [{c.accent}]⏱ {duration:.1f}s[/]")
        
        if fail_count == 0 and workflow.on_success:
            subprocess.run(workflow.on_success, shell=True)
        
        return fail_count == 0
    
    def create_template(self, name: str) -> Path:
        """Create a workflow template"""
        template = """name: {name}
description: My workflow

# Environment variables
env:
  NODE_ENV: development

# Steps to execute
steps:
  - name: Install dependencies
    run: npm install
    
  - name: Run tests
    run: npm test
    continue_on_error: false
    
  - name: Build
    run: npm run build

# Hooks
on_success: echo "Workflow completed!"
on_failure: echo "Workflow failed!"
"""
        path = WORKFLOWS_DIR / f"{name}.yml"
        path.write_text(template.format(name=name))
        return path

def get_workflow_engine() -> WorkflowEngine:
    return WorkflowEngine()
