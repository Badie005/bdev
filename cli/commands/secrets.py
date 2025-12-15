"""
B.DEV CLI - Secrets Commands
Encrypted secrets management
"""
import typer
from getpass import getpass

from utils.ui import console, print_header, print_success, print_error, print_warning

app = typer.Typer(help="Gestion des secrets", no_args_is_help=True)

def ensure_unlocked():
    """Ensure vault is unlocked"""
    from core.vault import get_vault
    
    vault = get_vault()
    
    if not vault.is_initialized():
        console.print("[yellow]Vault non initialisé[/yellow]")
        if typer.confirm("Créer un nouveau vault?"):
            password = getpass("Mot de passe: ")
            confirm = getpass("Confirmer: ")
            if password != confirm:
                print_error("Mots de passe différents")
                raise typer.Exit(1)
            vault.init(password)
            print_success("Vault créé!")
        else:
            raise typer.Exit(1)
    
    if vault._key is None:
        password = getpass("Mot de passe vault: ")
        if not vault.unlock(password):
            print_error("Mot de passe incorrect")
            raise typer.Exit(1)
    
    return vault

@app.command("init")
def init_vault():
    """Initialiser le vault"""
    from core.vault import get_vault
    
    vault = get_vault()
    
    if vault.is_initialized():
        print_warning("Vault déjà initialisé")
        if not typer.confirm("Réinitialiser? (données perdues)"):
            return
    
    password = getpass("Nouveau mot de passe: ")
    confirm = getpass("Confirmer: ")
    
    if password != confirm:
        print_error("Mots de passe différents")
        raise typer.Exit(1)
    
    vault.init(password)
    print_success("Vault initialisé!")

@app.command("set")
def set_secret(
    key: str = typer.Argument(..., help="Nom du secret"),
    value: str = typer.Argument(None, help="Valeur (ou input sécurisé)")
):
    """Définir un secret"""
    vault = ensure_unlocked()
    
    if value is None:
        value = getpass(f"Valeur pour '{key}': ")
    
    vault.set(key, value)
    print_success(f"Secret '{key}' enregistré")

@app.command("get")
def get_secret(
    key: str = typer.Argument(..., help="Nom du secret"),
    show: bool = typer.Option(False, "--show", "-s", help="Afficher la valeur")
):
    """Récupérer un secret"""
    vault = ensure_unlocked()
    
    value = vault.get(key)
    
    if value is None:
        print_error(f"Secret '{key}' non trouvé")
        raise typer.Exit(1)
    
    if show:
        console.print(f"[#FF6B35]{key}[/#FF6B35] = {value}")
    else:
        console.print(f"[#FF6B35]{key}[/#FF6B35] = [dim]********[/dim]")
        console.print("[dim]Utilisez --show pour afficher[/dim]")

@app.command("list")
def list_secrets():
    """Lister les secrets"""
    vault = ensure_unlocked()
    
    print_header("Secrets Vault", "Chiffré")
    
    keys = vault.list_keys()
    
    if not keys:
        console.print("[dim]Aucun secret stocké[/dim]")
        return
    
    for key in keys:
        console.print(f"  • [#FF6B35]{key}[/#FF6B35]")
    
    console.print(f"\n[dim]Total: {len(keys)} secrets[/dim]")

@app.command("delete")
def delete_secret(
    key: str = typer.Argument(..., help="Nom du secret")
):
    """Supprimer un secret"""
    vault = ensure_unlocked()
    
    if typer.confirm(f"Supprimer '{key}'?"):
        if vault.delete(key):
            print_success(f"Secret '{key}' supprimé")
        else:
            print_error(f"Secret '{key}' non trouvé")

@app.command("export")
def export_secrets():
    """Exporter les secrets comme variables d'environnement"""
    vault = ensure_unlocked()
    
    env = vault.export_env()
    
    print_header("Export Env", "Variables")
    
    for key, value in env.items():
        console.print(f"export {key}=\"{value}\"")
