"""
B.DEV CLI Core - Configuration System
Hierarchical configuration with defaults
"""
import json
from pathlib import Path
from typing import Any, Dict, Optional

# Paths
HOME = Path.home()
BDEV_PATH = HOME / "Dev" / ".bdev"
GLOBAL_CONFIG = BDEV_PATH / "config.json"
SESSION_FILE = BDEV_PATH / "cache" / "session.json"

DEFAULT_CONFIG = {
    "ai": {
        "model": "codellama:7b",
        "fallback_model": "phi3:mini",
        "memory_enabled": True,
        "max_context_tokens": 4000
    },
    "projects": {
        "path": str(HOME / "Dev" / "Projects"),
        "templates_path": str(BDEV_PATH / "templates")
    },
    "ui": {
        "color": True,
        "verbose": False
    },
    "repl": {
        "history_size": 500,
        "prompt": "bdev> "
    }
}

class Config:
    """Hierarchical configuration manager"""
    
    _instance: Optional['Config'] = None
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._load()
        return cls._instance
    
    def _load(self):
        self._config = DEFAULT_CONFIG.copy()
        
        # Load global config
        if GLOBAL_CONFIG.exists():
            try:
                user_config = json.loads(GLOBAL_CONFIG.read_text())
                self._deep_merge(self._config, user_config)
            except:
                pass
    
    def _deep_merge(self, base: dict, override: dict):
        """Deep merge override into base"""
        for key, value in override.items():
            if key in base and isinstance(base[key], dict) and isinstance(value, dict):
                self._deep_merge(base[key], value)
            else:
                base[key] = value
    
    def get(self, key: str, default: Any = None) -> Any:
        """Get config value using dot notation (e.g., 'ai.model')"""
        keys = key.split('.')
        value = self._config
        for k in keys:
            if isinstance(value, dict) and k in value:
                value = value[k]
            else:
                return default
        return value
    
    def set(self, key: str, value: Any):
        """Set config value using dot notation"""
        keys = key.split('.')
        target = self._config
        for k in keys[:-1]:
            target = target.setdefault(k, {})
        target[keys[-1]] = value
    
    def save(self):
        """Save current config to disk"""
        GLOBAL_CONFIG.write_text(json.dumps(self._config, indent=2))
    
    @property
    def all(self) -> Dict:
        return self._config

# Singleton accessor
def get_config() -> Config:
    return Config()
