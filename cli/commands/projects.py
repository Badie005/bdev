import typer
from pathlib import Path
from datetime import datetime
import subprocess
import json
from utils.config import DEV_PATH, TEMPLATES_PATH
from utils.helpers import detect_project_type, get_git_status, format_time_ago, run_script
from utils.ui import console, create_table, print_header, print_error, print_info, print_success

app = typer.Typer(help="Gestion des projets", no_args_is_help=True)

@app.command("list")
def list_projects(
    sort: str = typer.Option("name", "--sort", "-s", help="Trier par: name, type, date"),
    filter_type: str = typer.Option(None, "--type", "-t", help="Filtrer par type")
):
    """Liste tous les projets"""
    projects = [p for p in DEV_PATH.iterdir() if p.is_dir() and not p.name.startswith('.')]
    
    project_data = []
    with console.status("[bold green]Analyse des projets...[/bold green]"):
        for proj in projects:
            proj_type = detect_project_type(proj)
            if filter_type and proj_type.lower() != filter_type.lower():
                continue
            project_data.append({
                "name": proj.name,
                "type": proj_type,
                "mtime": datetime.fromtimestamp(proj.stat().st_mtime),
                "git": get_git_status(proj)
            })
    
    if sort == "date":
        project_data.sort(key=lambda x: x["mtime"], reverse=True)
    elif sort == "type":
        project_data.sort(key=lambda x: x["type"])
    else:
        project_data.sort(key=lambda x: x["name"].lower())
    
    table = create_table(["Nom", "Type", "Modifié", "Git"], title="Projets B.DEV")
    
    for proj in project_data:
        git_style = "green" if proj["git"] == "Clean" else "yellow"
        if proj["git"] == "No Git": git_style = "dim white"
        
        table.add_row(
            f"[bold #FF6B35]{proj['name']}[/bold cyan]",
            proj['type'],
            format_time_ago(proj['mtime']),
            f"[{git_style}]{proj['git']}[/{git_style}]"
        )
    
    console.print(table)
    console.print(f"Total: [bold #FF6B35]{len(project_data)}[/bold cyan] projets")

@app.command()
def new(
    template: str = typer.Argument(None, help="Nom du template"),
    name: str = typer.Argument(None, help="Nom du nouveau projet")
):
    """Creer un nouveau projet (Wizard Interactif)"""
    
    # Mode Wizard si arguments manquants
    if not template or not name:
        print_header("Nouveau Projet", "Mode Interactif")
        
        if not name:
            name = typer.prompt("Nom du projet")
            
        if not template:
            # On liste les templates dispos pour le choix
            templates = [t.name for t in TEMPLATES_PATH.iterdir() if t.is_dir()]
            if not templates:
                print_error("Aucun template trouvé dans .bdev/templates")
                raise typer.Exit(1)
                
            console.print("\n[bold]Templates disponibles:[/bold]")
            for i, t in enumerate(templates):
                console.print(f"  [#FF6B35]{i+1}. {t}[/#FF6B35]")
                
            choice = typer.prompt("\nChoisissez un template (numéro)", type=int)
            if 0 < choice <= len(templates):
                template = templates[choice - 1]
            else:
                print_error("Choix invalide")
                raise typer.Exit(1)

    print_header("Création en cours", f"Template: {template} | Nom: {name}")
    run_script("new-project.ps1", [template, name])

@app.command()
def open(name: str = typer.Argument(..., help="Nom du projet")):
    """Ouvrir un projet dans VSCode"""
    matches = [p for p in DEV_PATH.iterdir() if p.is_dir() and name.lower() in p.name.lower()]
    
    if not matches:
        print_error(f"Projet '{name}' non trouve")
        raise typer.Exit(1)
    
    if len(matches) == 1:
        project = matches[0]
    else:
        print_info("Plusieurs projets correspondent:")
        for i, p in enumerate(matches):
            console.print(f"  {i+1}. {p.name}")
        choice = typer.prompt("Choisissez un numero", type=int)
        project = matches[choice - 1]
    
    print_success(f"Ouverture de {project.name}...")
    subprocess.run(["code", str(project)], shell=True)

@app.command()
def templates():
    """Liste les templates disponibles"""
    table = create_table(["Template", "Description"], title="Templates Disponibles")
    
    for template in TEMPLATES_PATH.iterdir():
        if template.is_dir():
            readme = template / "README.md"
            desc = ""
            if readme.exists():
                first_line = readme.read_text().split('\n')[0]
                desc = first_line.replace('#', '').strip()
            
            table.add_row(f"[yellow]{template.name}[/yellow]", desc)
            
    console.print(table)
    console.print("[dim]Utilisez 'bdev new <template> <nom>' pour créer un projet[/dim]")

@app.command()
def describe(name: str = typer.Argument(..., help="Nom du projet")):
    """Analyse détaillée d'un projet"""
    matches = [p for p in DEV_PATH.iterdir() if p.is_dir() and name.lower() in p.name.lower()]
    if not matches:
        print_error(f"Projet '{name}' non trouve")
        raise typer.Exit(1)
        
    project = matches[0]
    print_header("Analyse", project.name)
    
    # 1. Taille sur le disque (approximative, sans node_modules/venv pour aller vite ou avec ?)
    # On va faire simple :
    total_size = 0
    file_count = 0
    for p in project.rglob('*'):
        if p.is_file() and not "node_modules" in p.parts and not ".git" in p.parts and not "venv" in p.parts:
            total_size += p.stat().st_size
            file_count += 1
            
    size_mb = total_size / (1024 * 1024)
    
    # 2. Dépendances
    deps = []
    if (project / "package.json").exists():
        try:
            pkg = json.loads((project / "package.json").read_text())
            deps = list(pkg.get("dependencies", {}).keys())
        except: pass
    elif (project / "requirements.txt").exists():
         deps = (project / "requirements.txt").read_text().splitlines()
         
    # 3. Last Commit
    last_commit = "Inconnu"
    if (project / ".git").exists():
        try:
             res = subprocess.run(["git", "log", "-1", "--format=%cd (%cr)"], cwd=project, capture_output=True, text=True, encoding='utf-8')
             if res.returncode == 0:
                 last_commit = res.stdout.strip()
        except: pass

    # Affichage
    from rich.columns import Columns
    from rich.panel import Panel

    console.print(f"[bold]Chemin:[/bold] {project}")
    console.print(f"[bold]Taille Code:[/bold] {size_mb:.2f} MB ({file_count} fichiers sources)")
    console.print(f"[bold]Dernier Commit:[/bold] {last_commit}")
    console.print()
    
    if deps:
        console.print(Panel("\n".join(deps[:15]) + ("\n..." if len(deps) > 15 else ""), title=f"Dépendances ({len(deps)})", border_style="green"))
    else:
        console.print("[dim]Aucune dépendance détectée.[/dim]")

@app.command()
def find(query: str = typer.Argument(..., help="Terme de recherche")):
    """Recherche rapide de projets par nom"""
    matches = [p for p in DEV_PATH.iterdir() if p.is_dir() and query.lower() in p.name.lower()]
    
    if not matches:
        print_error(f"Aucun projet correspondant à '{query}'")
        raise typer.Exit(1)
    
    console.print(f"[bold #FF6B35]Résultats pour '{query}':[/bold cyan]")
    for p in matches:
        proj_type = detect_project_type(p)
        console.print(f"  {p.name} [dim]({proj_type})[/dim]")

@app.command()
def run(name: str = typer.Argument(..., help="Nom du projet")):
    """Lance le serveur de développement du projet"""
    matches = [p for p in DEV_PATH.iterdir() if p.is_dir() and name.lower() in p.name.lower()]
    
    if not matches:
        print_error(f"Projet '{name}' non trouvé")
        raise typer.Exit(1)
    
    project = matches[0]
    proj_type = detect_project_type(project)
    
    print_header("Dev Server", f"{project.name} ({proj_type})")
    
    # Auto-detect and run
    if (project / "package.json").exists():
        # Check for common scripts
        try:
            pkg = json.loads((project / "package.json").read_text())
            scripts = pkg.get("scripts", {})
            if "dev" in scripts:
                cmd = "npm run dev"
            elif "start" in scripts:
                cmd = "npm start"
            else:
                cmd = "npm start"
        except:
            cmd = "npm start"
        
        console.print(f"[green]Exécution: {cmd}[/green]")
        subprocess.run(cmd, cwd=project, shell=True)
        
    elif (project / "manage.py").exists():
        # Django
        console.print("[green]Exécution: python manage.py runserver[/green]")
        subprocess.run(["python", "manage.py", "runserver"], cwd=project)
        
    elif (project / "artisan").exists():
        # Laravel
        console.print("[green]Exécution: php artisan serve[/green]")
        subprocess.run(["php", "artisan", "serve"], cwd=project)
        
    else:
        print_error(f"Type de projet '{proj_type}' non supporté pour le démarrage automatique.")
        console.print("[dim]Naviguez manuellement vers le projet et lancez le serveur.[/dim]")
