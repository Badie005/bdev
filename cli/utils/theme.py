"""
B.DEV CLI - Theme Engine
Claude Code EXACT Design System
"""
from dataclasses import dataclass
from typing import Dict, Optional
from rich.style import Style
from rich.theme import Theme as RichTheme
from rich.console import Console

# ═══════════════════════════════════════════════════════════════
# CLAUDE CODE EXACT PALETTE
# ═══════════════════════════════════════════════════════════════

@dataclass
class ThemeColors:
    # Primary
    primary: str       # Main accent (actions)
    secondary: str     # Hover/alt accent
    # Background
    background: str    # Main bg
    background_alt: str # Secondary bg
    # Text
    text: str          # Primary text
    text_dim: str      # Secondary text
    muted: str         # Disabled
    # Semantic
    success: str
    warning: str
    error: str
    info: str
    # UI
    border: str

@dataclass
class ThemeDefinition:
    name: str
    description: str
    colors: ThemeColors

# ═══════════════════════════════════════════════════════════════
# CLAUDE CODE THEME (EXACT)
# ═══════════════════════════════════════════════════════════════

THEME_CLAUDE = ThemeDefinition(
    name="claude",
    description="Claude Code - Official Palette",
    colors=ThemeColors(
        # Primary - Claude Orange
        primary="#FF6B35",
        secondary="#FF8F6B",
        # Background - Pure Dark
        background="#1A1A1A",
        background_alt="#0F0F0F",
        # Text
        text="#FFFFFF",
        text_dim="#8E8E93",
        muted="#48484A",
        # Semantic
        success="#34C759",
        warning="#FF9500",
        error="#FF3B30",
        info="#007AFF",
        # UI
        border="#2C2C2E"
    )
)

# Alternative themes
THEME_GEMINI = ThemeDefinition(
    name="gemini",
    description="Gemini CLI Style",
    colors=ThemeColors(
        primary="#4285F4",
        secondary="#669DF6",
        background="#1A1A1A",
        background_alt="#0F0F0F",
        text="#FFFFFF",
        text_dim="#9AA0A6",
        muted="#5F6368",
        success="#34A853",
        warning="#FBBC04",
        error="#EA4335",
        info="#4285F4",
        border="#3C4043"
    )
)

THEME_MATRIX = ThemeDefinition(
    name="matrix",
    description="Cyberpunk Green",
    colors=ThemeColors(
        primary="#00FF00",
        secondary="#00CC00",
        background="#000000",
        background_alt="#001100",
        text="#00FF00",
        text_dim="#008800",
        muted="#004400",
        success="#00FF00",
        warning="#88FF00",
        error="#FF0000",
        info="#00FF00",
        border="#003300"
    )
)

PRESETS = {
    "claude": THEME_CLAUDE,
    "gemini": THEME_GEMINI,
    "matrix": THEME_MATRIX,
    "default": THEME_CLAUDE
}

class ThemeManager:
    _instance = None
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._init()
        return cls._instance
    
    def _init(self):
        self.current_theme = THEME_CLAUDE
        self._load_config()
    
    def _load_config(self):
        try:
            from commands.config import load_user_config
            config = load_user_config()
            theme_name = config.get("display", {}).get("theme", "claude")
            if theme_name in PRESETS:
                self.current_theme = PRESETS[theme_name]
        except:
            pass
    
    def get_console(self) -> Console:
        """Get a configured rich Console instance"""
        c = self.current_theme.colors
        
        rich_theme = RichTheme({
            # Core colors
            "primary": c.primary,
            "secondary": c.secondary,
            "success": c.success,
            "warning": c.warning,
            "error": c.error,
            "info": c.info,
            "dim": c.text_dim,
            "muted": c.muted,
            "border": c.border,
            # Override standard colors to Claude palette
            "cyan": c.primary,
            "bright_cyan": c.primary,
            "green": c.success,
            "bright_green": c.success,
            "red": c.error,
            "bright_red": c.error,
            "yellow": c.warning,
            "bright_yellow": c.warning,
            "blue": c.info,
            "bright_blue": c.info,
            "magenta": c.secondary,
            "bright_magenta": c.secondary,
            # Semantic
            "header": f"bold {c.primary}",
            "link": f"underline {c.primary}",
            "code": f"{c.text}",
        })
        
        return Console(theme=rich_theme)

def get_theme_manager() -> ThemeManager:
    return ThemeManager()
