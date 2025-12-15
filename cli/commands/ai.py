"""
B.DEV CLI - AI Commands
Uses the new core.ai engine with memory support
"""
import typer
from pathlib import Path

from utils.ui import console, print_error, print_warning, print_success, print_info
from utils.context import get_project_context

app = typer.Typer(help="Assistant IA avec mémoire", no_args_is_help=True)

@app.command()
def check():
    """Vérifie la connexion à Ollama et l'état de la mémoire"""
    import subprocess
    
    console.print("[bold #FF6B35]Vérification du système IA...[/bold cyan]")
    
    try:
        subprocess.run(["ollama", "--version"], check=True, capture_output=True)
        print_success("Ollama est installé")
        
        # Check models
        result = subprocess.run(["ollama", "list"], capture_output=True, text=True, encoding='utf-8')
        lines = result.stdout.strip().split('\n')[1:]
        models = [line.split()[0] for line in lines if line]
        
        if models:
            print_success(f"Modèles: {', '.join(models)}")
        else:
            print_warning("Aucun modèle installé")
        
        # Check memory
        from core.ai.engine import get_ai_engine
        engine = get_ai_engine()
        console.print(f"[dim]{engine.get_memory_summary()}[/dim]")
        
    except (subprocess.CalledProcessError, FileNotFoundError):
        print_error("Ollama non accessible")

@app.command()
def chat(
    prompt: str = typer.Argument(..., help="Question pour l'IA"),
    context: bool = typer.Option(True, "--context/--no-context", help="Inclure contexte projet"),
    memory: bool = typer.Option(True, "--memory/--no-memory", help="Utiliser la mémoire de conversation")
):
    """Discuter avec l'IA (avec mémoire de conversation)"""
    from core.ai.engine import get_ai_engine
    
    engine = get_ai_engine()
    
    # Build context
    ctx = ""
    if context:
        console.print("[dim]Chargement du contexte...[/dim]")
        ctx = get_project_context(Path.cwd())
    
    console.print(f"[bold #FF6B35]B.AI ({engine.model}) répond...[/bold cyan]\n")
    
    # Stream response
    for chunk in engine.chat(prompt, context=ctx, stream=True):
        print(chunk, end="", flush=True)
    
    print()  # Newline after response

@app.command()
def forget():
    """Effacer la mémoire de conversation"""
    from core.ai.engine import get_ai_engine
    
    engine = get_ai_engine()
    engine.clear_memory()
    print_success("Mémoire effacée. L'IA a oublié notre conversation.")

@app.command()
def memory():
    """Afficher l'état de la mémoire"""
    from core.ai.engine import get_ai_engine
    from core.session import get_session
    
    session = get_session()
    messages = session.get_ai_context()
    
    if not messages:
        console.print("[dim]Aucune conversation en mémoire.[/dim]")
        return
    
    console.print(f"[bold #FF6B35]Mémoire de conversation ({len(messages)} messages):[/bold cyan]\n")
    
    for msg in messages[-10:]:
        role = "[green]Vous[/green]" if msg["role"] == "user" else "[#FF6B35]B.AI[/#FF6B35]"
        content = msg["content"][:100] + "..." if len(msg["content"]) > 100 else msg["content"]
        console.print(f"  {role}: {content}")

@app.command()
def generate(
    description: str = typer.Argument(..., help="Description de ce qu'il faut générer"),
    output: str = typer.Option(None, "--output", "-o", help="Fichier de sortie")
):
    """Générer du code à partir d'une description"""
    from core.ai.engine import get_ai_engine
    
    engine = get_ai_engine()
    
    gen_prompt = f"""Generate code based on this description:
{description}

Output ONLY the code, no explanations. Use proper formatting."""
    
    console.print("[bold #FF6B35]Génération en cours...[/bold cyan]\n")
    
    full_response = ""
    for chunk in engine.chat(gen_prompt, stream=True):
        print(chunk, end="", flush=True)
        full_response += chunk
    
    print()
    
    if output:
        Path(output).write_text(full_response)
        print_success(f"Code sauvegardé dans {output}")
