import typer
from utils.helpers import run_script

app = typer.Typer(help="Maintenance système", no_args_is_help=True)

@app.command()
def status():
    """Résumé rapide du système (une ligne)"""
    import psutil
    from utils.ui import console
    
    cpu = psutil.cpu_percent(interval=0.5)
    ram = psutil.virtual_memory()
    disk = psutil.disk_usage('C:/')
    
    ram_used = ram.used / (1024**3)
    ram_total = ram.total / (1024**3)
    disk_free = disk.free / (1024**3)
    
    status_line = f"CPU: {cpu}% | RAM: {ram_used:.1f}/{ram_total:.1f}GB | Disque: {disk_free:.0f}GB libres"
    
    # Colorize based on thresholds
    cpu_color = "green" if cpu < 70 else ("yellow" if cpu < 90 else "red")
    ram_color = "green" if ram.percent < 70 else ("yellow" if ram.percent < 90 else "red")
    
    console.print(f"[{cpu_color}]CPU: {cpu}%[/{cpu_color}] | [{ram_color}]RAM: {ram_used:.1f}/{ram_total:.1f}GB[/{ram_color}] | Disque: {disk_free:.0f}GB libres")

@app.command()
def health():
    """Diagnostic système complet"""
    run_script("health-check.ps1")

@app.command()
def info():
    """Versions des outils installés"""
    import subprocess
    from utils.ui import console, create_table
    
    tools = [
        ("Node.js", ["node", "-v"]),
        ("npm", ["npm", "-v"]),
        ("Git", ["git", "--version"]),
        ("Python", ["python", "--version"]),
        ("PHP", ["php", "-v"]),
        ("Docker", ["docker", "--version"]),
    ]
    
    table = create_table(["Outil", "Version"], title="Outils Installés")
    
    for name, cmd in tools:
        try:
            result = subprocess.run(cmd, capture_output=True, text=True, encoding='utf-8', errors='replace', timeout=5)
            version = result.stdout.strip().split('\n')[0] if result.returncode == 0 else "[red]Non installé[/red]"
        except:
            version = "[red]Non installé[/red]"
        table.add_row(name, version)
    
    console.print(table)

@app.command()
def watch(interval: int = typer.Option(5, "--interval", "-i", help="Intervalle en secondes")):
    """Surveillance continue du systeme"""
    run_script("watch.ps1", ["-Interval", str(interval)])

@app.command()
def backup(
    full: bool = typer.Option(False, "--full", "-f", help="Backup complet"),
    dry_run: bool = typer.Option(False, "--dry-run", "-n", help="Simulation")
):
    """Backup des configurations"""
    args = []
    if full: args.append("-Full")
    if dry_run: args.append("-DryRun")
    run_script("backup.ps1", args if args else None)

@app.command()
def clean(dry_run: bool = typer.Option(False, "--dry-run", "-n", help="Simulation")):
    """Nettoyage du systeme"""
    args = ["-DryRun"] if dry_run else None
    run_script("clean-system.ps1", args)

@app.command()
def heavy(
    threshold_mb: int = typer.Option(100, help="Seuil en MB pour afficher les dossiers")
):
    """Scan les dossiers les plus lourds (node_modules, venv...)"""
    from utils.config import DEV_PATH
    from rich.progress import track
    from rich.table import Table
    from utils.ui import console
    
    console.print(f"[bold #FF6B35]Recherche des dossiers lourds (> {threshold_mb}MB)...[/bold cyan]")
    
    heavy_folders = []
    
    # On scanne les projets
    targets = ["node_modules", "venv", "target", "vendor", ".git", "dist", "build"]
    
    projects = [p for p in DEV_PATH.iterdir() if p.is_dir()]
    
    for proj in track(projects, description="Scanning..."):
        for target in targets:
            folder = proj / target
            if folder.exists() and folder.is_dir():
                # Calcul taille
                size = 0
                for f in folder.rglob('*'):
                     if f.is_file(): size += f.stat().st_size
                
                size_mb = size / (1024 * 1024)
                if size_mb > threshold_mb:
                    heavy_folders.append((folder, size_mb))
    
    # Sort
    heavy_folders.sort(key=lambda x: x[1], reverse=True)
    
    table = Table(title=f"Dossiers > {threshold_mb}MB")
    table.add_column("Dossier", style="#FF6B35")
    table.add_column("Taille", style="red bold")
    
    for folder, size in heavy_folders:
        table.add_row(str(folder), f"{size:.0f} MB")
        
    console.print(table)

@app.command()
def prune(
    confirm: bool = typer.Option(False, "--yes", "-y", help="Confirmer suppression")
):
    """Nettoyer interactivement les dossiers lourds"""
    # Simple wrapper pour l'instant
    console.print("[yellow]Scan pour nettoyage...[/yellow]")
    heavy(threshold_mb=200)
    console.print("\n[bold]Pour supprimer, supprimez manuellement ou utilisez 'bdev clean'[/bold]")
