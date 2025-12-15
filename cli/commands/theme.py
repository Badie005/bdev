"""
B.DEV CLI - Theme Switcher
Change the visual style of the CLI
"""
import typer
from utils.ui import console, print_header, print_success, print_error
from utils.theme import get_theme_manager, PRESETS

app = typer.Typer(help="Gestion des thèmes", no_args_is_help=True)

@app.command("list")
def list_themes():
    """Lister les thèmes disponibles"""
    current = get_theme_manager().current_theme.name
    
    print_header("Thèmes Visuels", "B.DEV Experience")
    
    from rich.table import Table
    table = Table(show_header=True, box=None)
    table.add_column("Nom", style="bold #FF6B35")
    table.add_column("Description")
    table.add_column("Status")
    
    for name, theme in PRESETS.items():
        if name == "default": continue
        status = "✔ Actif" if name == current else ""
        table.add_row(name, theme.description, f"[green]{status}[/]")
    
    console.print(table)

@app.command("set")
def set_theme(
    name: str = typer.Argument(..., help="Nom du thème")
):
    """Changer le thème actif"""
    if name not in PRESETS:
        print_error(f"Thème '{name}' inconnu")
        console.print("Utilisez 'bdev theme list' pour voir les options.")
        raise typer.Exit(1)
    
    # Update config
    from commands.config import load_user_config, save_user_config
    
    config = load_user_config()
    if "display" not in config:
        config["display"] = {}
    
    config["display"]["theme"] = name
    save_user_config(config)
    
    print_success(f"Thème '{name}' activé!")
    console.print("\n[dim]Redémarrez le CLI pour voir tous les changements[/dim]")

@app.command("preview")
def preview_theme(name: str):
    """Prévisualiser un thème"""
    if name not in PRESETS:
        print_error(f"Thème '{name}' inconnu")
        return

    theme = PRESETS[name]
    c = theme.colors
    
    console.print(f"\n[bold {c.primary}]=== Theme Preview: {name.upper()} ===[/]")
    console.print(f"Primary: [{c.primary}]██████[/] {c.primary}")
    console.print(f"Secondary: [{c.secondary}]██████[/] {c.secondary}")
    console.print(f"Accent: [{c.accent}]██████[/] {c.accent}")
    console.print(f"Success: [{c.success}]✔ Success[/]")
    console.print(f"Error: [{c.error}]✖ Error[/]")
    console.print(f"Warning: [{c.warning}]⚠ Warning[/]")
    console.print(f"Border: [{c.border}]──────[/]")
