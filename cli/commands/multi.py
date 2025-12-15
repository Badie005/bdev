"""
B.DEV CLI - Multi-Project Operations
Execute commands across multiple projects
"""
import typer
import subprocess
from pathlib import Path
from typing import List, Optional

from utils.config import DEV_PATH
from utils.ui import console, print_header, print_error, print_success, print_warning
from utils.helpers import detect_project_type, get_git_status

app = typer.Typer(help="Opérations multi-projets", no_args_is_help=True)

def get_projects(filter_type: Optional[str] = None, filter_tag: Optional[str] = None) -> List[Path]:
    """Get filtered list of projects"""
    projects = [p for p in DEV_PATH.iterdir() if p.is_dir() and not p.name.startswith('.')]
    
    if filter_type:
        projects = [p for p in projects if detect_project_type(p).lower() == filter_type.lower()]
    
    return projects

@app.command("status")
def multi_status(
    filter_type: Optional[str] = typer.Option(None, "--type", "-t", help="Filtrer par type (node, python, etc.)")
):
    """Status Git de tous les projets"""
    print_header("Multi-Project Status", "Git Status")
    
    projects = get_projects(filter_type)
    
    from rich.table import Table
    table = Table(show_header=True)
    table.add_column("Projet", style="#FF6B35")
    table.add_column("Type", style="dim")
    table.add_column("Git Status")
    table.add_column("Branch")
    
    for proj in projects:
        proj_type = detect_project_type(proj)
        git_status = get_git_status(proj)
        
        # Get branch
        branch = "-"
        if (proj / ".git").exists():
            try:
                result = subprocess.run(
                    ["git", "branch", "--show-current"],
                    cwd=proj, capture_output=True, text=True, encoding='utf-8'
                )
                branch = result.stdout.strip() or "detached"
            except:
                pass
        
        status_style = "green" if git_status == "Clean" else ("yellow" if git_status != "No Git" else "dim")
        table.add_row(proj.name, proj_type, f"[{status_style}]{git_status}[/{status_style}]", branch)
    
    console.print(table)
    console.print(f"\nTotal: [bold]{len(projects)}[/bold] projets")

@app.command("pull")
def multi_pull(
    filter_type: Optional[str] = typer.Option(None, "--type", "-t", help="Filtrer par type")
):
    """Git pull sur tous les projets"""
    print_header("Multi-Project Pull", "Mise à jour")
    
    projects = get_projects(filter_type)
    success = 0
    failed = 0
    
    for proj in projects:
        if not (proj / ".git").exists():
            continue
        
        console.print(f"[#FF6B35]{proj.name}[/cyan]... ", end="")
        
        try:
            result = subprocess.run(
                ["git", "pull", "--ff-only"],
                cwd=proj, capture_output=True, text=True, encoding='utf-8', timeout=30
            )
            if result.returncode == 0:
                console.print("[green]OK[/green]")
                success += 1
            else:
                console.print(f"[yellow]SKIP[/yellow] ({result.stderr.strip()[:30]})")
                failed += 1
        except Exception as e:
            console.print(f"[red]FAIL[/red]")
            failed += 1
    
    console.print(f"\n[green]{success} réussis[/green], [yellow]{failed} échoués[/yellow]")

@app.command("audit")
def multi_audit(
    filter_type: Optional[str] = typer.Option(None, "--type", "-t", help="Filtrer par type")
):
    """Audit de sécurité sur tous les projets Node"""
    print_header("Multi-Project Audit", "Sécurité")
    
    projects = get_projects(filter_type or "node")
    
    from rich.table import Table
    table = Table(show_header=True)
    table.add_column("Projet", style="#FF6B35")
    table.add_column("Vulnérabilités")
    table.add_column("Critical", style="red")
    table.add_column("High", style="yellow")
    
    for proj in projects:
        if not (proj / "package.json").exists():
            continue
        
        console.print(f"Scanning {proj.name}...", end="\r")
        
        try:
            result = subprocess.run(
                ["npm", "audit", "--json"],
                cwd=proj, capture_output=True, text=True, encoding='utf-8', timeout=60
            )
            
            import json
            try:
                data = json.loads(result.stdout)
                vulns = data.get("metadata", {}).get("vulnerabilities", {})
                total = sum(vulns.values())
                critical = vulns.get("critical", 0)
                high = vulns.get("high", 0)
                
                status = "[green]Safe[/green]" if total == 0 else f"[yellow]{total}[/yellow]"
                table.add_row(proj.name, status, str(critical), str(high))
            except:
                table.add_row(proj.name, "[dim]N/A[/dim]", "-", "-")
        except:
            table.add_row(proj.name, "[red]Error[/red]", "-", "-")
    
    console.print(table)

@app.command("run")
def multi_run(
    command: str = typer.Argument(..., help="Commande à exécuter"),
    filter_type: Optional[str] = typer.Option(None, "--type", "-t", help="Filtrer par type")
):
    """Exécuter une commande dans tous les projets"""
    print_header("Multi-Project Run", command)
    
    projects = get_projects(filter_type)
    
    for proj in projects:
        console.print(f"\n[bold #FF6B35]═══ {proj.name} ═══[/bold cyan]")
        
        try:
            subprocess.run(command, cwd=proj, shell=True, timeout=120)
        except subprocess.TimeoutExpired:
            print_warning("Timeout")
        except Exception as e:
            print_error(str(e))

@app.command("update")
def multi_update(
    filter_type: Optional[str] = typer.Option(None, "--type", "-t", help="Filtrer par type"),
    dry_run: bool = typer.Option(False, "--dry-run", "-n", help="Simulation")
):
    """Mettre à jour les dépendances de tous les projets"""
    print_header("Multi-Project Update", "Dry run" if dry_run else "Live")
    
    projects = get_projects(filter_type)
    
    for proj in projects:
        console.print(f"\n[bold #FF6B35]═══ {proj.name} ═══[/bold cyan]")
        
        if (proj / "package.json").exists():
            cmd = "npm outdated" if dry_run else "npm update"
            subprocess.run(cmd, cwd=proj, shell=True)
        elif (proj / "requirements.txt").exists():
            if dry_run:
                console.print("[dim]Python project - check requirements.txt manually[/dim]")
            else:
                subprocess.run("pip install -r requirements.txt --upgrade", cwd=proj, shell=True)
