"""
B.DEV CLI - Workflow Commands
YAML-based automation
"""
import typer
from pathlib import Path

from utils.ui import console, print_header, print_success, print_error

app = typer.Typer(help="Workflows automatisés", no_args_is_help=True)

@app.command("list")
def list_workflows():
    """Lister les workflows disponibles"""
    from core.workflow import get_workflow_engine
    
    engine = get_workflow_engine()
    workflows = engine.list_workflows()
    
    print_header("Workflows", "Automatisation")
    
    if not workflows:
        console.print("[dim]Aucun workflow. Créez-en un avec 'bdev workflow create'[/dim]")
        return
    
    from rich.table import Table
    table = Table(show_header=True)
    table.add_column("Nom", style="bold #FF6B35")
    table.add_column("Description")
    
    for name in workflows:
        wf = engine.load(name)
        desc = wf.description if wf else "-"
        table.add_row(name, desc)
    
    console.print(table)

@app.command("run")
def run_workflow(
    name: str = typer.Argument(..., help="Nom du workflow"),
    cwd: str = typer.Option(None, "--cwd", "-C", help="Répertoire de travail")
):
    """Exécuter un workflow"""
    from core.workflow import get_workflow_engine
    
    engine = get_workflow_engine()
    work_dir = Path(cwd) if cwd else Path.cwd()
    
    success = engine.run(name, cwd=work_dir)
    
    if success:
        print_success("Workflow terminé avec succès!")
    else:
        print_error("Workflow échoué")
        raise typer.Exit(1)

@app.command("create")
def create_workflow(
    name: str = typer.Argument(..., help="Nom du workflow")
):
    """Créer un nouveau workflow"""
    from core.workflow import get_workflow_engine
    
    engine = get_workflow_engine()
    path = engine.create_template(name)
    
    print_success(f"Workflow créé: {path}")
    console.print("[dim]Éditez le fichier YAML pour personnaliser[/dim]")
    
    # Open in editor
    if typer.confirm("Ouvrir dans l'éditeur?"):
        import subprocess
        import os
        # Windows: use start command
        os.startfile(str(path))

@app.command("edit")
def edit_workflow(name: str = typer.Argument(..., help="Nom du workflow")):
    """Éditer un workflow"""
    from core.workflow import WORKFLOWS_DIR
    import subprocess
    
    path = WORKFLOWS_DIR / f"{name}.yml"
    if not path.exists():
        print_error(f"Workflow '{name}' non trouvé")
        raise typer.Exit(1)
    
    import os
    os.startfile(str(path))

@app.command("show")
def show_workflow(name: str = typer.Argument(..., help="Nom du workflow")):
    """Afficher le contenu d'un workflow"""
    from core.workflow import get_workflow_engine
    from utils.theme import get_theme_manager
    
    c = get_theme_manager().current_theme.colors
    engine = get_workflow_engine()
    wf = engine.load(name)
    
    if not wf:
        print_error(f"Workflow '{name}' non trouvé")
        raise typer.Exit(1)
    
    print_header(wf.name, wf.description)
    
    console.print(f"[bold {c.secondary}]Steps:[/]")
    for i, step in enumerate(wf.steps, 1):
        console.print(f"  [{c.accent}]{i}.[/] {step.name}")
        console.print(f"     [dim]$ {step.run}[/dim]")
