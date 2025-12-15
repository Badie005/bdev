from pathlib import Path
import subprocess
import json
from datetime import datetime
from .ui import print_error
from .config import SCRIPTS_PATH

def detect_project_type(project_path: Path) -> str:
    """Detecte le type de projet"""
    if (project_path / "package.json").exists():
        try:
            pkg = json.loads((project_path / "package.json").read_text())
            deps = pkg.get("dependencies", {})
            if "next" in deps: return "Next.js"
            if "@angular/core" in deps: return "Angular"
            if "react" in deps: return "React"
            if "vue" in deps: return "Vue.js"
        except:
            pass
        return "Node.js"
    if (project_path / "composer.json").exists():
        try:
            composer = json.loads((project_path / "composer.json").read_text())
            if "laravel/framework" in composer.get("require", {}):
                return "Laravel"
        except:
            pass
        return "PHP"
    if (project_path / "requirements.txt").exists(): return "Python"
    if (project_path / "Cargo.toml").exists(): return "Rust"
    if (project_path / "go.mod").exists(): return "Go"
    return "Other"

def get_git_status(project_path: Path) -> str:
    """Recupere le statut Git"""
    if not (project_path / ".git").exists():
        return "No Git"
    try:
        result = subprocess.run(
            ["git", "status", "--porcelain"],
            cwd=project_path,
            capture_output=True,
            text=True,
            encoding='utf-8', 
            errors='replace',
            timeout=5
        )
        if result.stdout.strip():
            lines = len(result.stdout.strip().split('\n'))
            return f"{lines} changes"
        return "Clean"
    except:
        return "?"

def format_time_ago(dt: datetime) -> str:
    """Formate une date en temps relatif"""
    delta = datetime.now() - dt
    if delta.days > 30:
        return f"{delta.days // 30}mo"
    if delta.days > 0:
        return f"{delta.days}d"
    if delta.seconds > 3600:
        return f"{delta.seconds // 3600}h"
    return "<1h"

def run_script(script_name: str, args: list = None):
    """Execute un script PowerShell"""
    script_path = SCRIPTS_PATH / script_name
    cmd = ["powershell", "-NoProfile", "-File", str(script_path)]
    if args:
        cmd.extend(args)
    try:
        subprocess.run(cmd, check=True, encoding='utf-8', errors='replace')
    except subprocess.CalledProcessError as e:
        print_error(f"Erreur lors de l'execution du script {script_name}: {e}")
