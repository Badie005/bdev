"""
B.DEV CLI Core - Session Management
Persistent state across commands
"""
import json
from pathlib import Path
from datetime import datetime
from typing import Any, Dict, List, Optional

SESSION_DIR = Path.home() / "Dev" / ".bdev" / "cache"
SESSION_FILE = SESSION_DIR / "session.json"
HISTORY_FILE = SESSION_DIR / "history.json"

class Session:
    """Manages session state and command history"""
    
    _instance: Optional['Session'] = None
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._load()
        return cls._instance
    
    def _load(self):
        SESSION_DIR.mkdir(parents=True, exist_ok=True)
        
        self._state: Dict[str, Any] = {
            "current_project": None,
            "last_command": None,
            "started_at": datetime.now().isoformat(),
            "ai_context": []  # Conversation memory
        }
        
        self._history: List[str] = []
        
        # Load persisted session
        if SESSION_FILE.exists():
            try:
                self._state.update(json.loads(SESSION_FILE.read_text()))
            except:
                pass
        
        # Load history
        if HISTORY_FILE.exists():
            try:
                self._history = json.loads(HISTORY_FILE.read_text())
            except:
                pass
    
    def get(self, key: str, default: Any = None) -> Any:
        return self._state.get(key, default)
    
    def set(self, key: str, value: Any):
        self._state[key] = value
        self._save()
    
    def _save(self):
        SESSION_FILE.write_text(json.dumps(self._state, indent=2, default=str))
    
    # --- Command History ---
    
    def add_to_history(self, command: str):
        self._history.append(command)
        # Keep last 500
        self._history = self._history[-500:]
        HISTORY_FILE.write_text(json.dumps(self._history))
    
    @property
    def history(self) -> List[str]:
        return self._history
    
    # --- AI Conversation Memory ---
    
    def add_ai_message(self, role: str, content: str):
        """Add message to AI conversation memory"""
        self._state.setdefault("ai_context", []).append({
            "role": role,
            "content": content,
            "timestamp": datetime.now().isoformat()
        })
        # Keep last 20 messages
        self._state["ai_context"] = self._state["ai_context"][-20:]
        self._save()
    
    def get_ai_context(self) -> List[Dict]:
        return self._state.get("ai_context", [])
    
    def clear_ai_context(self):
        self._state["ai_context"] = []
        self._save()
    
    # --- Current Project ---
    
    @property
    def current_project(self) -> Optional[Path]:
        p = self._state.get("current_project")
        return Path(p) if p else None
    
    @current_project.setter
    def current_project(self, path: Optional[Path]):
        self._state["current_project"] = str(path) if path else None
        self._save()

# Singleton accessor
def get_session() -> Session:
    return Session()
