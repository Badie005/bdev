"""
B.DEV CLI - Analytics Commands
Developer productivity tracking
"""
import typer
from datetime import datetime

from utils.ui import console, print_header, create_table

app = typer.Typer(help="Analytiques développeur", no_args_is_help=True)

@app.command()
def today():
    """Statistiques d'aujourd'hui (Visuel)"""
    from core.analytics import get_tracker
    from utils.theme import get_theme_manager
    
    c = get_theme_manager().current_theme.colors
    tracker = get_tracker()
    stats = tracker.get_today_stats()
    
    print_header("Activité du Jour", datetime.now().strftime("%A %d %B"))
    
    from rich.panel import Panel
    from rich.table import Table
    from rich import box
    from rich.columns import Columns
    
    # KPIs
    kpi_grid = Table.grid(expand=True, padding=(0, 2))
    kpi_grid.add_column(justify="center", ratio=1)
    kpi_grid.add_column(justify="center", ratio=1)
    kpi_grid.add_column(justify="center", ratio=1)
    
    kpi_grid.add_row(
        f"[bold {c.primary}]{stats['commands']}[/]\n[dim]Commandes[/]",
        f"[bold {c.accent}]{stats['projects']}[/]\n[dim]Projets[/]",
        f"[bold {c.success}]{stats['ai_queries']}[/]\n[dim]IA Utils[/]"
    )
    
    console.print(Panel(kpi_grid, border_style=c.border, box=box.ROUNDED, padding=(1, 2)))
    console.print()
    
    # Command breakdown visualization
    breakdown = tracker.get_command_breakdown()
    if breakdown:
        console.print(f"[bold {c.secondary}]Top Commandes[/]")
        max_val = max(breakdown.values()) if breakdown else 1
        
        # Block characters for charts
        blocks = [" ", "1", "2", "3", "4", "5", "6", "7", "8"] # not used directly, using full blocks
        
        table = Table(box=None, show_header=False, padding=(0, 1))
        table.add_column("Cmd", style=c.primary, width=15)
        table.add_column("Bar", ratio=1)
        table.add_column("Count", justify="right", style=c.accent)
        
        for cmd, count in sorted(breakdown.items(), key=lambda x: -x[1])[:5]:
            width = int((count / max_val) * 40)
            bar = f"[{c.primary}]" + ("█" * width) + f"[/{c.primary}]"
            table.add_row(cmd, bar, str(count))
            
        console.print(table)

@app.command()
def week():
    """Statistiques de la semaine"""
    from core.analytics import get_tracker
    from utils.theme import get_theme_manager
    
    c = get_theme_manager().current_theme.colors
    tracker = get_tracker()
    stats = tracker.get_week_stats()
    
    print_header("Activité Semaine", "7 derniers jours")
    
    from rich.panel import Panel
    from rich.table import Table
    
    grid = Table.grid(expand=True, padding=(1, 4))
    grid.add_column(justify="center")
    grid.add_column(justify="center")
    
    # Total commands big number
    grid.add_row(
        f"[bold {c.primary} size=3]{stats.get('total_commands', 0)}[/]",
        f"[bold {c.accent} size=3]{stats.get('active_days', 0)}/7[/]"
    )
    grid.add_row(f"[dim]Total Commandes[/]", f"[dim]Jours Actifs[/]")
    
    console.print(Panel(grid, border_style=c.border, padding=(2, 2)))

@app.command()
def summary():
    """Résumé rapide"""
    from core.analytics import get_tracker
    from utils.theme import get_theme_manager
    c = get_theme_manager().current_theme.colors
    
    tracker = get_tracker()
    stats = tracker.get_today_stats()
    week = tracker.get_week_stats()
    
    console.print(f"[{c.primary}]Today:[/{c.primary}] {stats['commands']} cmds | [{c.accent}]Week:[/{c.accent}] {week.get('total_commands', 0)}")
