"""
B.DEV CLI - Native Git Commands
Complete Git integration without needing to remember options
"""
import typer
import subprocess
from pathlib import Path
from typing import Optional

from utils.ui import console, print_header, print_error, print_success, print_warning

app = typer.Typer(help="Git intégré", no_args_is_help=True)

def run_git(args: list, cwd: Path = None) -> subprocess.CompletedProcess:
    """Execute git command"""
    return subprocess.run(
        ["git"] + args,
        cwd=cwd or Path.cwd(),
        capture_output=True,
        text=True,
        encoding='utf-8',
        errors='replace'
    )

@app.command()
def status():
    """Status du repository (amélioré)"""
    result = run_git(["status", "--short", "--branch"])
    
    if result.returncode != 0:
        print_error("Pas un repository Git")
        raise typer.Exit(1)
    
    lines = result.stdout.strip().split('\n')
    
    # Parse branch info
    branch_line = lines[0] if lines else ""
    console.print(f"[bold #FF6B35]Branch:[/bold cyan] {branch_line.replace('## ', '')}")
    
    if len(lines) > 1:
        console.print("\n[bold]Changements:[/bold]")
        for line in lines[1:]:
            if line.startswith('M'):
                console.print(f"  [yellow]M[/yellow] {line[2:].strip()}")
            elif line.startswith('A'):
                console.print(f"  [green]A[/green] {line[2:].strip()}")
            elif line.startswith('D'):
                console.print(f"  [red]D[/red] {line[2:].strip()}")
            elif line.startswith('??'):
                console.print(f"  [dim]?[/dim] {line[2:].strip()}")
            else:
                console.print(f"  {line}")
    else:
        print_success("Working tree clean")

@app.command()
def commit(
    message: str = typer.Argument(..., help="Message de commit"),
    all: bool = typer.Option(False, "--all", "-a", help="Ajouter tous les fichiers modifiés")
):
    """Commit les changements"""
    if all:
        result = run_git(["add", "-A"])
        if result.returncode != 0:
            print_error(f"Erreur git add: {result.stderr}")
            raise typer.Exit(1)
    
    result = run_git(["commit", "-m", message])
    
    if result.returncode == 0:
        print_success(f"Commit créé: {message}")
    else:
        if "nothing to commit" in result.stdout.lower():
            print_warning("Rien à committer")
        else:
            print_error(result.stderr or result.stdout)

@app.command()
def push(
    force: bool = typer.Option(False, "--force", "-f", help="Force push")
):
    """Push vers le remote"""
    args = ["push"]
    if force:
        args.append("--force-with-lease")
    
    console.print("[#FF6B35]Pushing...[/#FF6B35]")
    result = run_git(args)
    
    if result.returncode == 0:
        print_success("Push réussi")
    else:
        print_error(result.stderr)

@app.command()
def pull():
    """Pull depuis le remote"""
    console.print("[#FF6B35]Pulling...[/#FF6B35]")
    result = run_git(["pull", "--ff-only"])
    
    if result.returncode == 0:
        if "Already up to date" in result.stdout:
            console.print("[dim]Déjà à jour[/dim]")
        else:
            print_success("Pull réussi")
            console.print(result.stdout)
    else:
        print_error(result.stderr)

@app.command()
def log(
    count: int = typer.Option(10, "--count", "-n", help="Nombre de commits")
):
    """Historique des commits (formaté)"""
    result = run_git(["log", f"-{count}", "--oneline", "--graph", "--decorate"])
    
    if result.returncode == 0:
        console.print(result.stdout)
    else:
        print_error(result.stderr)

@app.command()
def branch(
    name: Optional[str] = typer.Argument(None, help="Nom de la branche à créer/switcher"),
    delete: bool = typer.Option(False, "--delete", "-d", help="Supprimer la branche")
):
    """Gestion des branches"""
    if not name:
        # List branches
        result = run_git(["branch", "-v"])
        console.print(result.stdout)
    elif delete:
        result = run_git(["branch", "-d", name])
        if result.returncode == 0:
            print_success(f"Branche '{name}' supprimée")
        else:
            print_error(result.stderr)
    else:
        # Switch or create
        result = run_git(["checkout", name])
        if result.returncode != 0:
            # Try creating
            result = run_git(["checkout", "-b", name])
        
        if result.returncode == 0:
            print_success(f"Switched to '{name}'")
        else:
            print_error(result.stderr)

@app.command()
def stash(
    action: str = typer.Argument("push", help="Action: push, pop, list, drop")
):
    """Gestion du stash"""
    if action == "push":
        result = run_git(["stash", "push"])
    elif action == "pop":
        result = run_git(["stash", "pop"])
    elif action == "list":
        result = run_git(["stash", "list"])
    elif action == "drop":
        result = run_git(["stash", "drop"])
    else:
        print_error(f"Action inconnue: {action}")
        raise typer.Exit(1)
    
    if result.returncode == 0:
        console.print(result.stdout or "[dim]OK[/dim]")
    else:
        print_error(result.stderr)

@app.command()
def diff(
    staged: bool = typer.Option(False, "--staged", "-s", help="Voir les changements staged")
):
    """Voir les différences"""
    args = ["diff"]
    if staged:
        args.append("--staged")
    
    result = run_git(args)
    
    if result.stdout:
        console.print(result.stdout)
    else:
        console.print("[dim]Aucune différence[/dim]")

@app.command()
def add(
    files: Optional[str] = typer.Argument(None, help="Fichiers à ajouter (ou '.' pour tout)")
):
    """Ajouter des fichiers au staging"""
    if not files:
        files = "."
    
    result = run_git(["add", files])
    
    if result.returncode == 0:
        print_success(f"Ajouté: {files}")
    else:
        print_error(result.stderr)

@app.command()
def reset(
    hard: bool = typer.Option(False, "--hard", help="Reset hard")
):
    """Reset les changements"""
    args = ["reset"]
    if hard:
        args.append("--hard")
    
    result = run_git(args)
    
    if result.returncode == 0:
        print_success("Reset effectué")
    else:
        print_error(result.stderr)
