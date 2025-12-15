"""
B.DEV CLI - UI Components
Rich UI components using the centralized Theme Engine
"""
from rich.console import Console
from rich.table import Table
from rich.panel import Panel
from rich.text import Text
from rich.prompt import Prompt, Confirm
from rich import box
from rich.layout import Layout
from rich.align import Align
from rich.padding import Padding
from rich.progress import Progress, SpinnerColumn, TextColumn, BarColumn, TimeRemainingColumn

from utils.theme import get_theme_manager

# Global Console instance managed by theme
console = get_theme_manager().get_console()

def reload_console():
    """Reload console with latest theme"""
    global console
    console = get_theme_manager().get_console()

def print_header(title: str, subtitle: str = ""):
    """Affiche un en-tête stylisé avec le thème actif"""
    theme = get_theme_manager().current_theme
    c = theme.colors
    
    reload_console()
    console.print()
    
    # Header grid
    grid = Table.grid(expand=True)
    grid.add_column(justify="left", ratio=1)
    grid.add_column(justify="right")
    
    # Title
    title_text = Text(title.upper(), style=f"bold {c.primary}")
    if subtitle:
        title_text.append(f" • {subtitle}", style=f"dim {c.secondary}")
    
    # Status badges
    badges = Text("B.DEV CLI", style=f"bold {c.accent}")
    
    grid.add_row(title_text, badges)
    
    # Render with custom border color
    console.print(
        Panel(
            Padding(grid, (0, 1)),
            border_style=c.border,
            box=box.HEAVY_EDGE if theme.ascii_art_style == "block" else box.ROUNDED
        )
    )
    console.print()

def print_success(message: str):
    c = get_theme_manager().current_theme.colors
    console.print(f"[{c.success}]✔ {message}[/{c.success}]")

def print_error(message: str):
    c = get_theme_manager().current_theme.colors
    console.print(f"[{c.error}]✖ {message}[/{c.error}]")

def print_warning(message: str):
    c = get_theme_manager().current_theme.colors
    console.print(f"[{c.warning}]⚠ {message}[/{c.warning}]")

def print_info(message: str):
    c = get_theme_manager().current_theme.colors
    console.print(f"[{c.info}]ℹ {message}[/{c.info}]")

def print_panel(content, title: str = None, style: str = None):
    """Show a content panel"""
    c = get_theme_manager().current_theme.colors
    border = style or c.primary
    
    console.print(
        Panel(
            content,
            title=f"[bold {c.secondary}]{title}[/]" if title else None,
            border_style=border,
            box=box.ROUNDED,
            padding=(1, 2)
        )
    )

def create_table(columns: list[str], title: str = None) -> Table:
    """Crée une table stylisée"""
    c = get_theme_manager().current_theme.colors
    
    table = Table(
        show_header=True, 
        header_style=f"bold {c.secondary}", 
        title=f"[bold {c.primary}]{title}[/]" if title else None,
        border_style=c.border,
        box=box.SIMPLE,
        padding=(0, 2),
        collapse_padding=True
    )
    for col in columns:
        table.add_column(col)
    return table

def interactive_list(options: list[str], title: str = "Select option") -> str:
    """
    Select an option from a list
    (Simple implementation for now, could be enhanced with prompt_toolkit)
    """
    c = get_theme_manager().current_theme.colors
    
    console.print(f"\n[bold {c.primary}]{title}:[/]")
    for i, opt in enumerate(options, 1):
        console.print(f"  [{c.accent}]{i}.[/] {opt}")
    
    choice = Prompt.ask(f"[{c.secondary}]Choice[/]", choices=[str(i) for i in range(1, len(options)+1)])
    return options[int(choice)-1]

def create_progress():
    """Create a styled progress bar"""
    c = get_theme_manager().current_theme.colors
    return Progress(
        SpinnerColumn(style=c.accent),
        TextColumn("[progress.description]{task.description}"),
        BarColumn(complete_style=c.primary, finished_style=c.success),
        TimeRemainingColumn(),
        console=console
    )
