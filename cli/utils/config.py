from pathlib import Path
import json
from typing import Dict, Any

# Chemins
HOME = Path.home()
DEV_PATH = HOME / "Dev" / "Projects"
BDEV_PATH = HOME / "Dev" / ".bdev"
CLI_PATH = BDEV_PATH / "cli"
TEMPLATES_PATH = BDEV_PATH / "templates"
SCRIPTS_PATH = BDEV_PATH / "scripts"
DATA_PATH = BDEV_PATH / "data"
CONFIG_PATH = BDEV_PATH / "config.json"

def load_config() -> Dict[str, Any]:
    """Charge la configuration"""
    if CONFIG_PATH.exists():
        try:
            return json.loads(CONFIG_PATH.read_text())
        except json.JSONDecodeError:
            return {}
    return {}

def save_config(config: Dict[str, Any]):
    """Sauvegarde la configuration"""
    CONFIG_PATH.write_text(json.dumps(config, indent=2))
