"""
B.DEV CLI - Quick Action Shortcuts
start, test, build, deploy - Universal commands that work on any project
"""
import typer
import subprocess
import json
from pathlib import Path
from typing import Optional

from utils.ui import console, print_header, print_error, print_success, print_warning
from utils.helpers import detect_project_type

app = typer.Typer(help="Actions rapides", no_args_is_help=False)

def get_project_info(path: Path) -> dict:
    """Detect project type and available scripts"""
    info = {
        "type": detect_project_type(path),
        "scripts": {},
        "has_tests": False,
        "has_build": False
    }
    
    # Node.js
    if (path / "package.json").exists():
        try:
            pkg = json.loads((path / "package.json").read_text())
            info["scripts"] = pkg.get("scripts", {})
            info["has_tests"] = "test" in info["scripts"]
            info["has_build"] = "build" in info["scripts"]
        except:
            pass
    
    # Python
    if (path / "manage.py").exists():
        info["type"] = "Django"
        info["has_tests"] = True
    elif (path / "pytest.ini").exists() or (path / "tests").exists():
        info["has_tests"] = True
    
    # Laravel
    if (path / "artisan").exists():
        info["type"] = "Laravel"
        info["has_tests"] = True
    
    return info

@app.command()
def start(
    project: Optional[str] = typer.Argument(None, help="Projet à démarrer (optionnel)")
):
    """Démarrer le serveur de développement"""
    from utils.config import DEV_PATH
    
    if project:
        matches = [p for p in DEV_PATH.iterdir() if p.is_dir() and project.lower() in p.name.lower()]
        if not matches:
            print_error(f"Projet '{project}' non trouvé")
            raise typer.Exit(1)
        path = matches[0]
    else:
        path = Path.cwd()
    
    info = get_project_info(path)
    
    print_header("Start Dev Server", f"{path.name} ({info['type']})")
    
    # Determine command based on project type
    if "dev" in info["scripts"]:
        cmd = "npm run dev"
    elif "start" in info["scripts"]:
        cmd = "npm start"
    elif "serve" in info["scripts"]:
        cmd = "npm run serve"
    elif info["type"] == "Django":
        cmd = "python manage.py runserver"
    elif info["type"] == "Flask":
        cmd = "flask run"
    elif info["type"] == "Laravel":
        cmd = "php artisan serve"
    elif (path / "package.json").exists():
        cmd = "npm start"
    else:
        print_error(f"Impossible de détecter comment démarrer ce projet ({info['type']})")
        console.print("[dim]Naviguez manuellement et lancez le serveur[/dim]")
        raise typer.Exit(1)
    
    console.print(f"[green]Exécution:[/green] {cmd}")
    console.print("[dim]Ctrl+C pour arrêter[/dim]\n")
    
    subprocess.run(cmd, cwd=path, shell=True)

@app.command()
def test(
    project: Optional[str] = typer.Argument(None, help="Projet à tester"),
    watch: bool = typer.Option(False, "--watch", "-w", help="Mode watch")
):
    """Lancer les tests"""
    from utils.config import DEV_PATH
    
    if project:
        matches = [p for p in DEV_PATH.iterdir() if p.is_dir() and project.lower() in p.name.lower()]
        if not matches:
            print_error(f"Projet '{project}' non trouvé")
            raise typer.Exit(1)
        path = matches[0]
    else:
        path = Path.cwd()
    
    info = get_project_info(path)
    
    print_header("Tests", path.name)
    
    # Determine test command
    if "test" in info["scripts"]:
        cmd = "npm test" + (" -- --watch" if watch else "")
    elif info["type"] == "Django":
        cmd = "python manage.py test"
    elif info["type"] == "Laravel":
        cmd = "php artisan test"
    elif (path / "pytest.ini").exists() or (path / "tests").is_dir():
        cmd = "pytest" + (" -f" if watch else "")
    elif (path / "package.json").exists():
        cmd = "npm test"
    else:
        print_error("Impossible de détecter comment tester ce projet")
        raise typer.Exit(1)
    
    console.print(f"[green]Exécution:[/green] {cmd}\n")
    subprocess.run(cmd, cwd=path, shell=True)

@app.command()
def build(
    project: Optional[str] = typer.Argument(None, help="Projet à build")
):
    """Build pour production"""
    from utils.config import DEV_PATH
    
    if project:
        matches = [p for p in DEV_PATH.iterdir() if p.is_dir() and project.lower() in p.name.lower()]
        if not matches:
            print_error(f"Projet '{project}' non trouvé")
            raise typer.Exit(1)
        path = matches[0]
    else:
        path = Path.cwd()
    
    info = get_project_info(path)
    
    print_header("Build Production", path.name)
    
    # Determine build command
    if "build" in info["scripts"]:
        cmd = "npm run build"
    elif info["type"] == "Laravel":
        cmd = "npm run build && php artisan optimize"
    elif (path / "setup.py").exists():
        cmd = "python setup.py build"
    elif (path / "pyproject.toml").exists():
        cmd = "python -m build"
    else:
        print_error("Impossible de détecter comment builder ce projet")
        raise typer.Exit(1)
    
    console.print(f"[green]Exécution:[/green] {cmd}\n")
    result = subprocess.run(cmd, cwd=path, shell=True)
    
    if result.returncode == 0:
        print_success("Build réussi!")
    else:
        print_error("Build échoué")

@app.command()
def fix(
    project: Optional[str] = typer.Argument(None, help="Projet à fixer")
):
    """Auto-fix (lint, format)"""
    from utils.config import DEV_PATH
    
    if project:
        matches = [p for p in DEV_PATH.iterdir() if p.is_dir() and project.lower() in p.name.lower()]
        if not matches:
            print_error(f"Projet '{project}' non trouvé")
            raise typer.Exit(1)
        path = matches[0]
    else:
        path = Path.cwd()
    
    info = get_project_info(path)
    
    print_header("Auto-Fix", path.name)
    
    # Try various fix commands
    commands_run = []
    
    if "lint:fix" in info["scripts"]:
        console.print("[#FF6B35]Running lint:fix...[/#FF6B35]")
        subprocess.run("npm run lint:fix", cwd=path, shell=True)
        commands_run.append("lint:fix")
    elif "lint" in info["scripts"]:
        console.print("[#FF6B35]Running lint --fix...[/#FF6B35]")
        subprocess.run("npm run lint -- --fix", cwd=path, shell=True)
        commands_run.append("lint --fix")
    
    if "format" in info["scripts"]:
        console.print("[#FF6B35]Running format...[/#FF6B35]")
        subprocess.run("npm run format", cwd=path, shell=True)
        commands_run.append("format")
    
    if (path / "pyproject.toml").exists() or info["type"] == "Python":
        console.print("[#FF6B35]Running black + isort...[/#FF6B35]")
        subprocess.run("black . && isort .", cwd=path, shell=True)
        commands_run.append("black + isort")
    
    if commands_run:
        print_success(f"Exécuté: {', '.join(commands_run)}")
    else:
        print_warning("Aucun outil de fix trouvé")

@app.command()
def deploy(
    env: str = typer.Option("staging", "--env", "-e", help="Environnement: staging, production")
):
    """Déployer le projet"""
    path = Path.cwd()
    info = get_project_info(path)
    
    print_header("Deploy", f"{path.name} → {env}")
    
    # Check for common deployment setups
    if (path / "vercel.json").exists() or "vercel" in str(info.get("scripts", {})):
        cmd = f"vercel {'--prod' if env == 'production' else ''}"
    elif (path / "netlify.toml").exists():
        cmd = f"netlify deploy {'--prod' if env == 'production' else ''}"
    elif "deploy" in info["scripts"]:
        cmd = "npm run deploy"
    elif (path / "Dockerfile").exists():
        print_warning("Docker detected. Manual deployment required.")
        raise typer.Exit(0)
    else:
        print_error("Aucune configuration de déploiement détectée")
        console.print("[dim]Ajoutez vercel.json, netlify.toml, ou un script 'deploy'[/dim]")
        raise typer.Exit(1)
    
    console.print(f"[green]Exécution:[/green] {cmd}\n")
    subprocess.run(cmd, cwd=path, shell=True)
