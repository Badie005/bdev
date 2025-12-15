"""
B.DEV CLI - Agent Commands
AI Agents for autonomous tasks
"""
import typer
from pathlib import Path
from typing import Optional

from utils.ui import console, print_header, print_error, print_success

app = typer.Typer(help="Agents IA autonomes", no_args_is_help=True)

@app.command("list")
def list_agents():
    """Liste les agents disponibles"""
    from core.ai.agents.base import AgentRegistry
    # Import agents to register them
    from core.ai.agents import review
    
    print_header("Agents IA", "Assistants autonomes")
    
    agents = AgentRegistry.list_agents()
    
    if not agents:
        console.print("[dim]Aucun agent disponible[/dim]")
        return
    
    from rich.table import Table
    table = Table(show_header=True)
    table.add_column("Agent", style="bold #FF6B35")
    table.add_column("Description")
    
    for name, desc in agents.items():
        table.add_row(name, desc)
    
    console.print(table)

@app.command()
def review(
    target: str = typer.Argument(".", help="Fichier ou dossier à analyser"),
    strict: bool = typer.Option(False, "--strict", "-s", help="Mode strict"),
    focus: str = typer.Option("general", "--focus", "-f", help="Focus: security, performance, style, general")
):
    """Review de code avec IA"""
    from core.ai.agents.base import AgentRegistry, AgentTask
    from core.ai.agents import review as review_module  # Force import
    
    agent = AgentRegistry.get("review")
    if not agent:
        print_error("Agent 'review' non disponible")
        raise typer.Exit(1)
    
    target_path = Path(target).resolve()
    
    print_header("Code Review", target_path.name)
    
    task = AgentTask(
        name="code-review",
        description=f"Review code in {target_path}",
        target=target_path,
        options={"strict": strict, "focus": focus}
    )
    
    # Execute and stream output
    result = None
    for output in agent.execute(task):
        if isinstance(output, str):
            print(output, end="", flush=True)
        else:
            result = output
    
    print()

@app.command()
def document(
    target: str = typer.Argument(".", help="Fichier ou dossier à documenter"),
    style: str = typer.Option("jsdoc", "--style", "-s", help="Style: jsdoc, docstring, markdown")
):
    """Générer de la documentation avec IA"""
    from core.ai.engine import get_ai_engine
    
    target_path = Path(target).resolve()
    print_header("Documentation Generator", target_path.name)
    
    if not target_path.exists():
        print_error(f"{target_path} n'existe pas")
        raise typer.Exit(1)
    
    # Read code
    if target_path.is_file():
        code = target_path.read_text(encoding='utf-8', errors='replace')
    else:
        # Get main files
        code = ""
        for ext in ['.py', '.js', '.ts']:
            for f in target_path.glob(f'*{ext}'):
                code += f"\n\n=== {f.name} ===\n" + f.read_text(encoding='utf-8', errors='replace')
                if len(code) > 5000:
                    break
    
    if not code:
        print_error("Aucun code trouvé")
        raise typer.Exit(1)
    
    prompt = f"""Generate {style} documentation for this code.
Include:
- Function/class descriptions
- Parameter types and descriptions
- Return values
- Usage examples

CODE:
{code[:6000]}
"""
    
    console.print("[#FF6B35]Génération en cours...[/#FF6B35]\n")
    
    engine = get_ai_engine()
    for chunk in engine.chat(prompt, stream=True):
        print(chunk, end="", flush=True)
    
    print()
    print_success("Documentation générée!")

@app.command()
def explain(
    target: str = typer.Argument(..., help="Fichier à expliquer")
):
    """Expliquer du code avec IA"""
    from core.ai.engine import get_ai_engine
    
    target_path = Path(target).resolve()
    
    if not target_path.exists() or not target_path.is_file():
        print_error(f"{target_path} n'est pas un fichier valide")
        raise typer.Exit(1)
    
    print_header("Code Explainer", target_path.name)
    
    code = target_path.read_text(encoding='utf-8', errors='replace')
    
    prompt = f"""Explain this code in detail. Be thorough but clear.
Include:
1. What the code does overall
2. Key components and their roles
3. Important patterns or techniques used
4. Potential issues or improvements

CODE:
{code[:8000]}
"""
    
    console.print("[#FF6B35]Analyse en cours...[/#FF6B35]\n")
    
    engine = get_ai_engine()
    for chunk in engine.chat(prompt, stream=True):
        print(chunk, end="", flush=True)
    
    print()
