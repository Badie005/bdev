from pathlib import Path
import os

def get_project_context(path: Path = None) -> str:
    """Récupère le contexte du projet courant pour l'IA"""
    if path is None:
        path = Path.cwd()
        
    context = []
    
    # 1. Info de base
    context.append(f"Project Path: {path}")
    
    # NEW: Si on est à la racine ou dans Dev, on liste les projets dispos
    # Cela aide l'IA à répondre à "quels sont mes projets?"
    dev_projects = Path.home() / "Dev" / "Projects"
    if path == Path.home() or path == Path.home() / "Dev":
        context.append("\nAvailable Projects:")
        if dev_projects.exists():
            try:
                # On liste les dossiers dans Projects
                projs = [p.name for p in dev_projects.iterdir() if p.is_dir() and not p.name.startswith('.')]
                context.append(", ".join(projs))
            except Exception:
                context.append("(Unable to list projects)")

    # 2. Structure des fichiers (limité au niveau 1 pour l'instant pour la rapidité)
    context.append("\nFile Structure (Current Dir):")
    try:
        items = [p.name for p in path.iterdir() if not p.name.startswith('.')]
        context.append(", ".join(items))
    except Exception as e:
        context.append(f"Error reading directory: {e}")
        
    # 3. Contenu des fichiers clés
    key_files = ["package.json", "composer.json", "requirements.txt", "Cargo.toml", "go.mod", "README.md"]
    
    context.append("\nKey Configuration Files:")
    for filename in key_files:
        p = path / filename
        if p.exists() and p.is_file():
            try:
                # On lit seulement les 50 premières lignes pour ne pas surcharger le prompt
                content = p.read_text(encoding='utf-8').splitlines()[:50]
                context.append(f"\n--- {filename} ---")
                context.append("\n".join(content))
            except Exception:
                pass
                
    return "\n".join(context)
