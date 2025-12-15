"""
B.DEV CLI - User Preferences & Customization
Allow users to configure the CLI to their liking
"""
import typer
import json
from pathlib import Path
from typing import Optional

from utils.ui import console, print_header, print_success, print_error, print_warning

app = typer.Typer(help="Configuration utilisateur", no_args_is_help=True)

CONFIG_FILE = Path.home() / "Dev" / ".bdev" / "config.json"

def load_user_config() -> dict:
    """Load user configuration"""
    if CONFIG_FILE.exists():
        try:
            return json.loads(CONFIG_FILE.read_text())
        except:
            pass
    return {}

def save_user_config(config: dict):
    """Save user configuration"""
    CONFIG_FILE.parent.mkdir(parents=True, exist_ok=True)
    CONFIG_FILE.write_text(json.dumps(config, indent=2))

@app.command("show")
def show_config():
    """Afficher la configuration actuelle"""
    config = load_user_config()
    
    print_header("Configuration", str(CONFIG_FILE))
    
    if not config:
        console.print("[dim]Aucune configuration personnalisée[/dim]")
        console.print("\n[#FF6B35]Utilisez 'bdev config set <clé> <valeur>' pour configurer[/#FF6B35]")
        return
    
    from rich.tree import Tree
    tree = Tree("[bold #FF6B35]Config[/bold cyan]")
    
    def add_to_tree(node, data, prefix=""):
        for key, value in data.items():
            if isinstance(value, dict):
                branch = node.add(f"[#FF6B35]{key}[/#FF6B35]")
                add_to_tree(branch, value, prefix + "  ")
            else:
                node.add(f"[#FF6B35]{key}[/#FF6B35] = [green]{value}[/green]")
    
    add_to_tree(tree, config)
    console.print(tree)

@app.command("set")
def set_config(
    key: str = typer.Argument(..., help="Clé (ex: ai.model, aliases.gs)"),
    value: str = typer.Argument(..., help="Valeur")
):
    """Définir une valeur de configuration"""
    config = load_user_config()
    
    # Support dot notation
    keys = key.split('.')
    target = config
    for k in keys[:-1]:
        target = target.setdefault(k, {})
    
    # Try to parse as JSON for complex values
    try:
        parsed_value = json.loads(value)
    except:
        parsed_value = value
    
    target[keys[-1]] = parsed_value
    save_user_config(config)
    
    print_success(f"Configuré: {key} = {parsed_value}")

@app.command("get")
def get_config(key: str = typer.Argument(..., help="Clé à récupérer")):
    """Obtenir une valeur de configuration"""
    config = load_user_config()
    
    keys = key.split('.')
    value = config
    for k in keys:
        if isinstance(value, dict) and k in value:
            value = value[k]
        else:
            print_error(f"Clé '{key}' non trouvée")
            raise typer.Exit(1)
    
    console.print(f"[#FF6B35]{key}[/#FF6B35] = [green]{value}[/green]")

@app.command("reset")
def reset_config(
    confirm: bool = typer.Option(False, "--yes", "-y", help="Confirmer")
):
    """Réinitialiser la configuration"""
    if not confirm:
        if not typer.confirm("Réinitialiser toute la configuration?"):
            raise typer.Exit(0)
    
    if CONFIG_FILE.exists():
        CONFIG_FILE.unlink()
    
    print_success("Configuration réinitialisée")

@app.command("alias")
def manage_alias(
    name: str = typer.Argument(None, help="Nom de l'alias"),
    command: str = typer.Argument(None, help="Commande associée"),
    delete: bool = typer.Option(False, "--delete", "-d", help="Supprimer l'alias")
):
    """Gérer les alias personnalisés"""
    config = load_user_config()
    aliases = config.setdefault("aliases", {})
    
    if not name:
        # List aliases
        print_header("Aliases", "Raccourcis personnalisés")
        
        if not aliases:
            console.print("[dim]Aucun alias défini[/dim]")
            console.print("\n[#FF6B35]bdev config alias <nom> <commande>[/#FF6B35]")
            return
        
        from rich.table import Table
        table = Table(show_header=True)
        table.add_column("Alias", style="bold cyan")
        table.add_column("Commande", style="green")
        
        for alias, cmd in aliases.items():
            table.add_row(alias, cmd)
        
        console.print(table)
        return
    
    if delete:
        if name in aliases:
            del aliases[name]
            save_user_config(config)
            print_success(f"Alias '{name}' supprimé")
        else:
            print_error(f"Alias '{name}' non trouvé")
        return
    
    if not command:
        print_error("Spécifiez la commande pour cet alias")
        raise typer.Exit(1)
    
    aliases[name] = command
    save_user_config(config)
    print_success(f"Alias créé: {name} → {command}")

@app.command("init")
def init_config():
    """Initialiser la configuration avec les valeurs par défaut"""
    default_config = {
        "user": {
            "name": "",
            "editor": "code"
        },
        "ai": {
            "model": "codellama:7b",
            "fallback_model": "phi3:mini",
            "memory_enabled": True
        },
        "repl": {
            "prompt": "bdev> ",
            "history_size": 500
        },
        "aliases": {
            "gs": "git status",
            "gp": "git push",
            "gl": "git pull",
            "gc": "git commit",
            "t": "test",
            "b": "build",
            "d": "deploy"
        },
        "projects": {
            "default_template": "basic",
            "auto_open": True
        }
    }
    
    if CONFIG_FILE.exists():
        if not typer.confirm("Configuration existante. Écraser?"):
            raise typer.Exit(0)
    
    save_user_config(default_config)
    print_success("Configuration initialisée!")
    show_config()

@app.command("edit")
def edit_config():
    """Ouvrir le fichier de configuration dans l'éditeur"""
    import subprocess
    
    config = load_user_config()
    editor = config.get("user", {}).get("editor", "code")
    
    # Ensure config file exists
    if not CONFIG_FILE.exists():
        save_user_config({})
    
    console.print(f"[#FF6B35]Ouverture avec {editor}...[/#FF6B35]")
    subprocess.run([editor, str(CONFIG_FILE)], shell=True)
