"""
B.DEV CLI - Natural Language Command Parser
Gemini CLI-style: "do something" gets parsed into actual commands
"""
import re
from typing import Optional, Tuple, List
from dataclasses import dataclass

@dataclass
class ParsedIntent:
    """Parsed user intent"""
    action: str
    target: Optional[str] = None
    args: List[str] = None
    confidence: float = 0.0
    
    def __post_init__(self):
        if self.args is None:
            self.args = []

class NaturalLanguageParser:
    """Parse natural language into CLI commands"""
    
    # Intent patterns
    PATTERNS = {
        # Project management
        r"(list|show|voir|affiche).*(projects?|projets?)": ("list", None),
        r"(open|ouvre|ouvrir)\s+(.+)": ("open", 1),
        r"(create|crée|créer|new|nouveau).*(project|projet)\s+(.+)": ("new", 2),
        
        # Git
        r"(commit|commiter)\s+['\"]?(.+?)['\"]?$": ("git commit", 1),
        r"(push|pousse|envoie)": ("git push", None),
        r"(pull|tire|récupère)": ("git pull", None),
        r"(status|statut|état)\s*(git)?": ("git status", None),
        
        # Quick actions
        r"(start|démarre|lance)\s*(le)?\s*(serveur|server|dev)?": ("start", None),
        r"(test|teste|tester)": ("test", None),
        r"(build|compile|construire)": ("build", None),
        r"(deploy|déploie|déployer)": ("deploy", None),
        r"(fix|corrige|corriger|lint)": ("fix", None),
        
        # System
        r"(health|santé|diagnostic)": ("system health", None),
        r"(clean|nettoie|nettoyer)": ("system clean", None),
        
        # AI
        r"(review|revue|analyse)\s*(code|le code)?": ("agent review", None),
        r"(explain|explique|expliquer)\s+(.+)": ("agent explain", 1),
        r"(document|documente|documenter)": ("agent document", None),
        
        # Multi
        r"(update|màj|mise à jour)\s*(all|tous|tout)?": ("multi update", None),
        r"(audit|sécurité|security)\s*(all|tous)?": ("multi audit", None),
    }
    
    def parse(self, text: str) -> Optional[ParsedIntent]:
        """Parse natural language into command intent"""
        text = text.lower().strip()
        
        for pattern, (action, capture_group) in self.PATTERNS.items():
            match = re.search(pattern, text, re.IGNORECASE)
            if match:
                target = None
                if capture_group is not None and len(match.groups()) > capture_group:
                    target = match.group(capture_group + 1)
                
                return ParsedIntent(
                    action=action,
                    target=target.strip() if target else None,
                    confidence=0.9
                )
        
        # Fallback: AI chat
        return ParsedIntent(
            action="ai chat",
            target=text,
            confidence=0.5
        )
    
    def to_command(self, intent: ParsedIntent) -> str:
        """Convert intent to CLI command"""
        cmd = intent.action
        if intent.target:
            # Quote if has spaces
            if ' ' in intent.target:
                cmd += f' "{intent.target}"'
            else:
                cmd += f' {intent.target}'
        return cmd

# Global instance
_parser: Optional[NaturalLanguageParser] = None

def get_parser() -> NaturalLanguageParser:
    global _parser
    if _parser is None:
        _parser = NaturalLanguageParser()
    return _parser

def parse_natural_language(text: str) -> Tuple[str, float]:
    """
    Parse natural language and return (command, confidence)
    """
    parser = get_parser()
    intent = parser.parse(text)
    return parser.to_command(intent), intent.confidence
