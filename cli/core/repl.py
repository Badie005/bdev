"""
B.DEV CLI Core - Enhanced Interactive REPL Mode
Professional shell with history recall, shortcuts, and project context
"""
import sys
import os
import shlex
from typing import List, Optional, Callable
from pathlib import Path

try:
    import readline
    READLINE_AVAILABLE = True
except ImportError:
    try:
        import pyreadline3 as readline
        READLINE_AVAILABLE = True
    except ImportError:
        READLINE_AVAILABLE = False

from rich.console import Console
from rich.panel import Panel
from rich.text import Text
from rich.table import Table

from core.session import get_session
from core.config import get_config

console = Console()

class REPL:
    """Enhanced Interactive Read-Eval-Print Loop"""
    
    def __init__(self, app_callback: Callable):
        self.app_callback = app_callback
        self.session = get_session()
        self.config = get_config()
        self.running = True
        self.current_project: Optional[Path] = None
        
        # Command shortcuts - merge defaults with user config
        self.aliases = {
            "ls": "list",
            "ll": "list",
            "cd": "open",
            "s": "system status",
            "h": "system health",
            "q": "exit",
            "?": "help",
        }
        # Load user aliases from config
        try:
            from commands.config import load_user_config
            user_config = load_user_config()
            user_aliases = user_config.get("aliases", {})
            self.aliases.update(user_aliases)
        except:
            pass
        
        # Available commands for completion
        self.commands = [
            "list", "new", "open", "dashboard",
            "projects", "system", "dev", "ai",
            "help", "exit", "quit", "clear", "history",
            "ls", "ll", "cd", "s", "h", "q", "?",
            "ai chat", "ai memory", "ai forget", "ai generate",
            "system status", "system health", "system info",
            "dev audit", "dev git", "dev docker",
        ]
        
        self._setup_readline()
    
    def _setup_readline(self):
        if not READLINE_AVAILABLE:
            return
            
        readline.set_completer(self._completer)
        readline.parse_and_bind("tab: complete")
        
        history_file = Path.home() / "Dev" / ".bdev" / "cache" / ".repl_history"
        history_file.parent.mkdir(parents=True, exist_ok=True)
        try:
            readline.read_history_file(str(history_file))
        except:
            pass
        
        import atexit
        atexit.register(lambda: readline.write_history_file(str(history_file)))
    
    def _completer(self, text: str, state: int) -> Optional[str]:
        """Tab completion with project names"""
        line = readline.get_line_buffer() if READLINE_AVAILABLE else ""
        
        # If typing after 'open' or 'cd', complete project names
        if line.startswith(("open ", "cd ", "projects describe ", "projects run ")):
            from utils.config import DEV_PATH
            if DEV_PATH.exists():
                projects = [p.name for p in DEV_PATH.iterdir() if p.is_dir() and not p.name.startswith('.')]
                matches = [p for p in projects if p.lower().startswith(text.lower())]
                if state < len(matches):
                    return matches[state]
            return None
        
        # Default command completion
        options = [cmd for cmd in self.commands if cmd.startswith(text)]
        if state < len(options):
            return options[state]
        return None
    
    def _print_banner(self):
        """Print Claude-style welcome banner"""
        from utils.branding import get_repl_banner
        from utils.theme import get_theme_manager
        
        c = get_theme_manager().current_theme.colors
        
        # Clear and print banner
        console.print()
        console.print(get_repl_banner())
        console.print(f"[dim {c.border}]{'─' * 50}[/]")
        console.print()
    
    def _print_help(self):
        table = Table(title="Commandes REPL", show_header=False, box=None, padding=(0, 2))
        table.add_column("Cmd", style="bold green")
        table.add_column("Description")
        
        table.add_row("help, ?", "Afficher cette aide")
        table.add_row("clear", "Effacer l'écran")
        table.add_row("history", "Afficher l'historique")
        table.add_row("!n", "Ré-exécuter la commande n de l'historique")
        table.add_row("exit, quit, q", "Quitter le mode interactif")
        
        console.print(table)
        console.print()
        
        table2 = Table(title="Raccourcis", show_header=False, box=None, padding=(0, 2))
        table2.add_column("Alias", style="bold #FF6B35")
        table2.add_column("Commande")
        
        table2.add_row("ls, ll", "list")
        table2.add_row("cd <projet>", "open <projet>")
        table2.add_row("s", "system status")
        table2.add_row("h", "system health")
        
        console.print(table2)
        console.print()
        console.print("[dim]Toutes les commandes bdev fonctionnent sans préfixe[/dim]")
    
    def _handle_builtin(self, line: str) -> bool:
        """Handle REPL-specific commands. Returns True if handled."""
        cmd = line.split()[0].lower() if line else ""
        
        if cmd in ("exit", "quit", "q"):
            self.running = False
            console.print("[dim]Au revoir![/dim]")
            return True
        
        if cmd == "clear":
            os.system('cls' if os.name == 'nt' else 'clear')
            return True
        
        if cmd in ("help", "?"):
            self._print_help()
            return True
        
        if cmd == "history":
            for i, h in enumerate(self.session.history[-20:], 1):
                console.print(f"  [#FF6B35]{i:2}[/#FF6B35]. {h}")
            return True
        
        # History recall: !n
        if line.startswith("!"):
            try:
                idx = int(line[1:]) - 1
                history = self.session.history[-20:]
                if 0 <= idx < len(history):
                    recalled = history[idx]
                    console.print(f"[dim]→ {recalled}[/dim]")
                    self._execute(recalled)
                else:
                    console.print("[red]Index d'historique invalide[/red]")
            except ValueError:
                console.print("[red]Usage: !n (n = numéro de l'historique)[/red]")
            return True
        
        return False
    
    def _resolve_alias(self, line: str) -> str:
        """Resolve command aliases"""
        parts = line.split(maxsplit=1)
        if parts and parts[0] in self.aliases:
            resolved = self.aliases[parts[0]]
            if len(parts) > 1:
                resolved += " " + parts[1]
            return resolved
        return line
    
    def _execute(self, line: str):
        """Execute a command line with NLP fallback"""
        # Resolve aliases
        original_line = line
        line = self._resolve_alias(line)
        
        # Parse
        try:
            parts = shlex.split(line)
        except ValueError:
            parts = line.split()
        
        if not parts:
            return
        
        # Try to execute via typer
        try:
            sys.argv = ["bdev"] + parts
            self.app_callback()
        except SystemExit as e:
            # Check if it was a "no such command" error
            if e.code != 0:
                # Try natural language parsing
                from core.nlp import parse_natural_language
                command, confidence = parse_natural_language(original_line)
                
                if confidence >= 0.7 and command != original_line:
                    console.print(f"[dim]→ {command}[/dim]")
                    try:
                        parts = shlex.split(command)
                        sys.argv = ["bdev"] + parts
                        self.app_callback()
                    except SystemExit:
                        pass
        except Exception as e:
            console.print(f"[red]Erreur: {e}[/red]")
    
    def run(self):
        """Main REPL loop"""
        self._print_banner()
        
        while self.running:
            try:
                # Dynamic prompt with project context
                prompt = "[bold #FF6B35]bdev>[/bold cyan] "
                if self.current_project:
                    prompt = f"[bold #FF6B35]bdev[/bold cyan]:[magenta]{self.current_project.name}[/magenta]> "
                
                # Get input (use simple input, console.input causes issues)
                try:
                    line = input("bdev> ")
                except EOFError:
                    break
                
                line = line.strip()
                if not line:
                    continue
                
                # Save to history
                self.session.add_to_history(line)
                
                # Check for built-in REPL commands
                if self._handle_builtin(line):
                    continue
                
                # Execute command
                self._execute(line)
                console.print()  # Spacer
                
            except KeyboardInterrupt:
                console.print("\n[dim]Ctrl+C - Utilisez 'exit' pour quitter[/dim]")
                continue

def start_repl(app_callback: Callable):
    """Entry point for REPL mode"""
    repl = REPL(app_callback)
    repl.run()
