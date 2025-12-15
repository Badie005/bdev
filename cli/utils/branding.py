"""
B.DEV CLI - ASCII Art Mascot & Branding
Claude-style pixel art for headers
"""

# B.DEV Mascot (inspired by Claude's style but unique)
MASCOT_SMALL = """
[#d97757]  ██████  [/]
[#d97757] ████████ [/]
[#d97757]██[white]▓▓[/]██[white]▓▓[/]██[/]
[#d97757]██████████[/]
[#d97757]██[black]▄[/]████[black]▄[/]██[/]
[#d97757] ████████ [/]
[#d97757]  ██  ██  [/]
"""

MASCOT_LARGE = """
[#d97757]        ████████████        [/]
[#d97757]      ████████████████      [/]
[#d97757]    ████████████████████    [/]
[#d97757]   ██████[white]████[/]████[white]████[/]██████   [/]
[#d97757]  ████████████████████████  [/]
[#d97757]  ████████████████████████  [/]
[#d97757]  ████[black]██[/]████████████[black]██[/]████  [/]
[#d97757]   ██████████████████████   [/]
[#d97757]    ████████████████████    [/]
[#d97757]      ████████████████      [/]
[#d97757]        ████    ████        [/]
"""

# Simple B logo
LOGO_B = """
[#d97757]██████  [/]
[#d97757]██   ██ [/]
[#d97757]██████  [/]
[#d97757]██   ██ [/]
[#d97757]██████  [/]
"""

# Welcome message
WELCOME_ART = """
[#d97757]  ████  [/]     [bold #ebdbb2]B.DEV CLI[/]
[#d97757] ██  ██ [/]     [dim]Enterprise Developer Workstation[/]
[#d97757] ██████ [/]     [dim #928374]v1.0.0[/]
[#d97757] ██  ██ [/]
[#d97757] ██████ [/]
"""

# Stars decoration (like Claude's)
STARS = "[dim]*[/]"

def get_decorated_header():
    """Get a decorated header with stars"""
    import random
    stars = ""
    for _ in range(50):
        if random.random() > 0.92:
            stars += " * "
        else:
            stars += "   "
    return f"[dim]{stars}[/]"

def get_welcome_screen():
    """Full welcome screen like Claude Code"""
    return f"""
{get_decorated_header()}

{WELCOME_ART}

[dim]─────────────────────────────────────────────[/]

[bold #d97757]Let's get started.[/]

[dim]Type a command or use the interactive mode[/]
[dim]Run [bold #d97757]/help[/bold #d97757] for help, [bold #d97757]/theme[/bold #d97757] to customize[/]

{get_decorated_header()}
"""

def get_repl_banner():
    """REPL mode banner"""
    return f"""
[#d97757]  ████  [/]  [bold #ebdbb2]B.DEV Interactive Mode[/]
[#d97757] ██████ [/]  [dim]Tab to complete • !n for history • exit to quit[/]
[#d97757] ██  ██ [/]
"""
