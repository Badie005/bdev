import typer
from utils.helpers import run_script

app = typer.Typer(help="Outils de développement", no_args_is_help=True)

@app.command()
def git(
    action: str = typer.Argument("status-all", help="Action: status-all, pull-all, find-dirty, clean-branches")
):
    """Utilitaires Git multi-projets"""
    run_script("git-helper.ps1", [action])

@app.command()
def docker(
    action: str = typer.Argument("status", help="Action: status, start, stop, clean, logs"),
    container: str = typer.Option(None, "--container", "-c", help="Nom du conteneur")
):
    """Gestion Docker"""
    args = [action]
    if container:
        args.extend(["-Container", container])
    run_script("docker-helper.ps1", args)

@app.command()
def ports(
    action: str = typer.Argument("list", help="Action: list, find, kill"),
    port: int = typer.Option(None, "--port", "-p", help="Numero de port")
):
    """Gestion des ports"""
    args = [action]
    if port:
        args.extend(["-Port", str(port)])
    run_script("port-manager.ps1", args)

@app.command()
def npm(
    action: str = typer.Argument("outdated-all", help="Action: outdated-all, update-all, clean-all, audit-all")
):
    """Utilitaires NPM multi-projets"""
    run_script("npm-helper.ps1", [action])

@app.command()
def db(
    action: str = typer.Argument("list", help="Action: list, create, drop, backup, restore"),
    name: str = typer.Option(None, "--name", "-n", help="Nom de la base"),
    file: str = typer.Option(None, "--file", "-f", help="Fichier SQL")
):
    """Gestion des bases de donnees"""
    args = [action]
    if name:
        args.extend(["-Name", name])
    if file:
        args.extend(["-File", file])
    run_script("db-helper.ps1", args)

@app.command()
def dotfiles(
    push: bool = typer.Option(False, "--push", help="Push local -> repo"),
    pull: bool = typer.Option(False, "--pull", help="Pull repo -> local")
):
    """Synchronisation des configurations"""
    args = []
    if push: args.append("-Push")
    elif pull: args.append("-Pull")
    else: args.append("-Status")
    run_script("sync-dotfiles.ps1", args)

@app.command()
def audit(
    project: str = typer.Argument(None, help="Nom du projet (optionnel, sinon dossier courant)")
):
    """Scan de sécurité des dépendances"""
    import subprocess
    from pathlib import Path
    from utils.config import DEV_PATH
    from utils.ui import console, print_error, print_success, print_warning
    
    # Determine project path
    if project:
        matches = [p for p in DEV_PATH.iterdir() if p.is_dir() and project.lower() in p.name.lower()]
        if not matches:
            print_error(f"Projet '{project}' non trouvé")
            raise typer.Exit(1)
        proj_path = matches[0]
    else:
        proj_path = Path.cwd()
    
    console.print(f"[bold #FF6B35]Audit de sécurité: {proj_path.name}[/bold cyan]")
    
    # Detect and run appropriate audit
    if (proj_path / "package.json").exists():
        console.print("[dim]Détecté: Node.js - Exécution de npm audit...[/dim]")
        result = subprocess.run(["npm", "audit", "--json"], cwd=proj_path, capture_output=True, text=True, encoding='utf-8', errors='replace')
        
        try:
            import json
            data = json.loads(result.stdout)
            vulns = data.get("metadata", {}).get("vulnerabilities", {})
            total = sum(vulns.values())
            
            if total == 0:
                print_success("Aucune vulnérabilité trouvée!")
            else:
                print_warning(f"Vulnérabilités trouvées: {total}")
                console.print(f"  Critical: [red]{vulns.get('critical', 0)}[/red]")
                console.print(f"  High: [yellow]{vulns.get('high', 0)}[/yellow]")
                console.print(f"  Moderate: {vulns.get('moderate', 0)}")
                console.print(f"  Low: {vulns.get('low', 0)}")
                console.print("\n[dim]Exécutez 'npm audit fix' pour corriger automatiquement.[/dim]")
        except:
            # Fallback: just show raw output summary
            if "found 0 vulnerabilities" in result.stdout.lower() or result.returncode == 0:
                print_success("Aucune vulnérabilité détectée.")
            else:
                print_warning("Vulnérabilités potentielles détectées. Lancez 'npm audit' pour détails.")
                
    elif (proj_path / "requirements.txt").exists():
        console.print("[dim]Détecté: Python - Vérification basique...[/dim]")
        # pip-audit n'est pas toujours installé, on fait simple
        console.print("[yellow]Pour un audit complet, installez et lancez: pip-audit[/yellow]")
    else:
        print_error("Type de projet non supporté pour l'audit.")
