#!/usr/bin/env python3
"""
Python CLI Template - Point d'entrée
=====================================
"""

import typer
from rich.console import Console

app = typer.Typer(
    name="mycli",
    help="🐍 Mon CLI Python personnalisé",
    add_completion=False
)
console = Console()


@app.command()
def hello(name: str = typer.Argument("World", help="Nom à saluer")):
    """👋 Saluer quelqu'un"""
    console.print(f"[bold green]Hello, {name}![/bold green]")


@app.command()
def info():
    """ℹ️ Afficher les informations"""
    console.print("[cyan]CLI créé avec B.DEV Template[/cyan]")
    console.print("[dim]Version: 0.1.0[/dim]")


@app.callback(invoke_without_command=True)
def main(ctx: typer.Context):
    """🐍 Mon CLI Python personnalisé"""
    if ctx.invoked_subcommand is None:
        console.print("[bold]Bienvenue dans votre CLI ![/bold]")
        console.print("[dim]Utilisez --help pour voir les commandes[/dim]")


if __name__ == "__main__":
    app()
