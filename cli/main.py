#!/usr/bin/env python3
"""
B.DEV CLI - Enterprise Developer Workstation
Claude-style design with warm orange palette
"""

import sys
import os

# FORCE Claude colors by overriding Rich's default styles BEFORE any import
os.environ["FORCE_COLOR"] = "1"

from rich.console import Console
from rich.theme import Theme as RichTheme
from rich import style as rich_style

# ═══════════════════════════════════════════════════════════════
# CLAUDE CODE - EXACT PALETTE
# ═══════════════════════════════════════════════════════════════
CLAUDE_ORANGE = "#FF6B35"    # Primary accent
CLAUDE_CORAL = "#FF8F6B"     # Secondary/hover
CLAUDE_WHITE = "#FFFFFF"     # Text
CLAUDE_GRAY = "#8E8E93"      # Dim text
CLAUDE_MUTED = "#48484A"     # Disabled
CLAUDE_BORDER = "#2C2C2E"    # Borders
CLAUDE_SUCCESS = "#34C759"   # Success
CLAUDE_WARNING = "#FF9500"   # Warning
CLAUDE_ERROR = "#FF3B30"     # Error
CLAUDE_INFO = "#007AFF"      # Info

# Create Claude theme
CLAUDE_THEME = RichTheme({
    "info": CLAUDE_ORANGE,
    "warning": CLAUDE_WARNING,
    "danger": CLAUDE_ERROR,
    "error": CLAUDE_ERROR,
    "success": CLAUDE_SUCCESS,
    "cyan": CLAUDE_ORANGE,
    "bright_cyan": CLAUDE_ORANGE,
    "green": CLAUDE_SUCCESS,
    "bright_green": CLAUDE_SUCCESS,
    "yellow": CLAUDE_WARNING,
    "magenta": CLAUDE_CORAL,
    "blue": CLAUDE_INFO,
    "repr.number": CLAUDE_ORANGE,
    "repr.str": CLAUDE_WHITE,
    "repr.path": CLAUDE_ORANGE,
    "rule.line": CLAUDE_MUTED,
    "prompt": CLAUDE_ORANGE,
    "prompt.choices": CLAUDE_GRAY,
})

# Monkey-patch Rich's default console
_claude_console = Console(theme=CLAUDE_THEME, force_terminal=True)

# Override the default styles in rich
import rich.default_styles
from rich.style import Style
rich.default_styles.DEFAULT_STYLES["info"] = Style(color=CLAUDE_ORANGE)
rich.default_styles.DEFAULT_STYLES["cyan"] = Style(color=CLAUDE_ORANGE)
rich.default_styles.DEFAULT_STYLES["bright_cyan"] = Style(color=CLAUDE_ORANGE)
rich.default_styles.DEFAULT_STYLES["green"] = Style(color=CLAUDE_SUCCESS)

import typer
from typer import rich_utils

# Override ALL Typer's colors with Claude palette
# Command styling
rich_utils.STYLE_OPTION = CLAUDE_ORANGE
rich_utils.STYLE_SWITCH = CLAUDE_ORANGE
rich_utils.STYLE_ARGUMENT = CLAUDE_ORANGE  # Command arguments
rich_utils.STYLE_METAVAR = CLAUDE_WHITE
rich_utils.STYLE_METAVAR_SEP = CLAUDE_GRAY
rich_utils.STYLE_USAGE = CLAUDE_ORANGE
rich_utils.STYLE_USAGE_COMMAND = f"bold {CLAUDE_ORANGE}"

# Panel borders
rich_utils.STYLE_COMMANDS_PANEL_BORDER = CLAUDE_MUTED
rich_utils.STYLE_OPTIONS_PANEL_BORDER = CLAUDE_MUTED
rich_utils.STYLE_ARGUMENTS_PANEL_BORDER = CLAUDE_MUTED
rich_utils.STYLE_ERRORS_PANEL_BORDER = CLAUDE_ERROR

# Help text
rich_utils.STYLE_HELPTEXT = ""
rich_utils.STYLE_HELPTEXT_FIRST_LINE = ""
rich_utils.STYLE_OPTION_HELP = ""
rich_utils.STYLE_OPTION_DEFAULT = CLAUDE_GRAY

# Status colors
rich_utils.STYLE_REQUIRED_SHORT = CLAUDE_ERROR
rich_utils.STYLE_REQUIRED_LONG = CLAUDE_ERROR
rich_utils.STYLE_NEGATIVE_OPTION = CLAUDE_ERROR
rich_utils.STYLE_DEPRECATED = CLAUDE_GRAY
rich_utils.STYLE_ABORTED = CLAUDE_ERROR

# AGGRESSIVE: Override Rich Table column styling for commands
from rich.table import Table
_orig_add_column = Table.add_column

def _claude_add_column(self, header="", *args, style=None, **kwargs):
    # Replace any cyan-ish styles with Claude orange
    if style and "cyan" in str(style).lower():
        style = CLAUDE_ORANGE
    # Also force command columns to orange
    if header.lower() in ["commands", "command", "name", "nom"]:
        style = f"bold {CLAUDE_ORANGE}"
    return _orig_add_column(self, header, *args, style=style, **kwargs)

Table.add_column = _claude_add_column

# AGGRESSIVE: Patch typer's help formatter directly
try:
    from typer.core import TyperGroup
    _get_help_text = getattr(TyperGroup, 'format_help_text', None)
except:
    pass

from commands import projects, system, devtools, ai, agents, analytics, multi, git, quick, config, theme, workflow, secrets
from utils.ui import console, print_header, print_info
from rich.traceback import install
install(show_locals=True, theme="monokai")

app = typer.Typer(
    name="bdev",
    help="B.DEV CLI - Enterprise Developer Workstation",
    add_completion=True,
    no_args_is_help=False,
    rich_markup_mode="rich"
)

# Integration des modules
app.add_typer(projects.app, name="projects", help="Gestion des projets")
app.add_typer(system.app, name="system", help="Maintenance système")
app.add_typer(devtools.app, name="dev", help="Outils de développement")
app.add_typer(ai.app, name="ai", help="Assistant IA")
app.add_typer(agents.app, name="agent", help="Agents IA autonomes")
app.add_typer(analytics.app, name="analytics", help="Analytiques")
app.add_typer(multi.app, name="multi", help="Multi-projets")
app.add_typer(git.app, name="git", help="Git intégré")
app.add_typer(config.app, name="config", help="Configuration")
app.add_typer(theme.app, name="theme", help="Thèmes visuels")
app.add_typer(workflow.app, name="workflow", help="Workflows automatisés")
app.add_typer(secrets.app, name="secrets", help="Vault secrets chiffré")

# Raccourcis root-level
@app.command("list", help="Liste les projets")
def list_shortcut(
    sort: str = typer.Option("name", "--sort", "-s"),
    filter_type: str = typer.Option(None, "--type", "-t")
):
    projects.list_projects(sort, filter_type)

@app.command("open", help="Ouvrir un projet")
def open_shortcut(name: str):
    projects.open(name)

@app.command("new", help="Nouveau projet")
def new_shortcut(template: str = None, name: str = None):
    projects.new(template, name)

# Quick actions at root level
@app.command("start", help="Démarrer le serveur dev")
def start_shortcut(project: str = None):
    quick.start(project)

@app.command("test", help="Lancer les tests")
def test_shortcut(project: str = None, watch: bool = False):
    quick.test(project, watch)

@app.command("build", help="Build production")
def build_shortcut(project: str = None):
    quick.build(project)

@app.command("fix", help="Auto-fix lint/format")
def fix_shortcut(project: str = None):
    quick.fix(project)

@app.command("deploy", help="Déployer")
def deploy_shortcut(env: str = "staging"):
    quick.deploy(env)

@app.command("do", help="Exécute en langage naturel (Gemini-style)")
def do_command(
    instruction: str = typer.Argument(..., help="Ce que vous voulez faire")
):
    """
    Exécute une instruction en langage naturel.
    Exemples:
      bdev do "ouvre mon projet"
      bdev do "commit mes changements"
      bdev do "lance les tests"
    """
    from core.nlp import parse_natural_language
    
    command, confidence = parse_natural_language(instruction)
    
    if confidence >= 0.7:
        console.print(f"[dim]→ {command}[/dim]")
        # Execute the parsed command
        import shlex
        try:
            parts = shlex.split(command)
            sys.argv = ["bdev"] + parts
            app()
        except SystemExit:
            pass
    else:
        # Low confidence - use AI to interpret
        console.print(f"[yellow]Je ne suis pas sûr... Je demande à l'IA[/yellow]")
        from core.ai.engine import get_ai_engine
        engine = get_ai_engine()
        
        prompt = f"""The user wants to: "{instruction}"
        
Available bdev commands:
- list: list projects
- open <name>: open project
- start: start dev server
- test: run tests
- build: build for production
- git status/commit/push/pull: git operations
- agent review/explain/document: AI code analysis

What command should I run? Reply with just the command, nothing else."""
        
        response = ""
        for chunk in engine.chat(prompt, stream=False):
            response += chunk
        
        suggested = response.strip().split('\n')[0]
        console.print(f"[#FF6B35]Suggestion:[/#FF6B35] {suggested}")
        
        if typer.confirm("Exécuter?"):
            try:
                parts = shlex.split(suggested)
                sys.argv = ["bdev"] + parts
                app()
            except:
                pass

@app.command("dashboard", help="Affiche le menu principal")
def dashboard():
    """Affiche le dashboard B.DEV - Claude Style"""
    from rich.panel import Panel
    from rich.layout import Layout
    from rich.table import Table
    from rich.align import Align
    from rich.columns import Columns
    from rich import box
    from datetime import datetime
    from utils.branding import WELCOME_ART, get_decorated_header
    from utils.theme import get_theme_manager
    
    c = get_theme_manager().current_theme.colors
    current_time = datetime.now().strftime("%A %d %B %H:%M")
    
    # Welcome screen with mascot
    console.print()
    console.print(get_decorated_header())
    console.print(WELCOME_ART)
    console.print(f"[dim {c.border}]{'─' * 50}[/]")
    console.print()
    
    # Quick commands in Claude style
    console.print(f"[bold {c.primary}]Quick Commands[/]")
    console.print()
    
    cmds = Table.grid(padding=(0, 4))
    cmds.add_column(style=f"bold {c.accent}")
    cmds.add_column(style="dim")
    
    cmds.add_row("bdev", "Interactive mode (REPL)")
    cmds.add_row("bdev start", "Start dev server")
    cmds.add_row("bdev ai chat", "AI assistant")
    cmds.add_row("bdev list", "List projects")
    cmds.add_row("bdev theme set <name>", "Change theme")
    
    console.print(cmds)
    console.print()
    
    # Recent activity
    console.print(f"[bold {c.secondary}]Recent Projects[/]")
    
    from pathlib import Path
    DEV_PATH = Path.home() / "Dev" / "Projects"
    
    if DEV_PATH.exists():
        proj_list = [p for p in DEV_PATH.iterdir() if p.is_dir() and not p.name.startswith('.')]
        proj_list.sort(key=lambda x: x.stat().st_mtime, reverse=True)
        
        for p in proj_list[:3]:
            time_str = datetime.fromtimestamp(p.stat().st_mtime).strftime("%H:%M")
            console.print(f"  [{c.accent}]•[/{c.accent}] {p.name} [dim]({time_str})[/dim]")
    
    console.print()
    console.print(f"[dim {c.border}]{'─' * 50}[/]")
    console.print(f"[dim]Type [bold {c.primary}]bdev[/] for interactive mode[/dim]")

@app.callback(invoke_without_command=True)
def main(ctx: typer.Context):
    """B.DEV CLI - Enterprise Developer Workstation"""
    if ctx.invoked_subcommand is None:
        # No command provided -> Launch REPL
        try:
            from core.repl import start_repl
            start_repl(app)
        except ImportError as e:
            # Fallback to dashboard if REPL not available
            console.print(f"[yellow]REPL non disponible: {e}[/yellow]")
            dashboard()

if __name__ == "__main__":
    app()
