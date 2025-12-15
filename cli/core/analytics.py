"""
B.DEV CLI - Analytics Tracker
Track developer activity and productivity
"""
import json
from pathlib import Path
from datetime import datetime, timedelta
from typing import Dict, List, Optional
from collections import defaultdict

ANALYTICS_DIR = Path.home() / "Dev" / ".bdev" / "analytics"
DAILY_FILE_TEMPLATE = "activity_{date}.json"

class ActivityTracker:
    """Tracks CLI activity for analytics"""
    
    _instance: Optional['ActivityTracker'] = None
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._init()
        return cls._instance
    
    def _init(self):
        ANALYTICS_DIR.mkdir(parents=True, exist_ok=True)
        self._today_file = ANALYTICS_DIR / DAILY_FILE_TEMPLATE.format(date=datetime.now().strftime("%Y-%m-%d"))
        self._load()
    
    def _load(self):
        if self._today_file.exists():
            try:
                self._data = json.loads(self._today_file.read_text())
            except:
                self._data = self._empty_day()
        else:
            self._data = self._empty_day()
    
    def _empty_day(self) -> Dict:
        return {
            "date": datetime.now().strftime("%Y-%m-%d"),
            "commands": [],
            "projects_opened": [],
            "ai_queries": 0,
            "time_spent_minutes": 0,
            "first_activity": None,
            "last_activity": None
        }
    
    def _save(self):
        self._today_file.write_text(json.dumps(self._data, indent=2))
    
    def log_command(self, command: str, project: Optional[str] = None):
        """Log a command execution"""
        now = datetime.now().isoformat()
        
        self._data["commands"].append({
            "command": command,
            "project": project,
            "timestamp": now
        })
        
        if not self._data["first_activity"]:
            self._data["first_activity"] = now
        self._data["last_activity"] = now
        
        if project and project not in self._data["projects_opened"]:
            self._data["projects_opened"].append(project)
        
        self._save()
    
    def log_ai_query(self):
        """Log an AI interaction"""
        self._data["ai_queries"] += 1
        self._save()
    
    def get_today_stats(self) -> Dict:
        """Get today's statistics"""
        return {
            "commands": len(self._data["commands"]),
            "projects": len(self._data["projects_opened"]),
            "ai_queries": self._data["ai_queries"],
            "first_activity": self._data["first_activity"],
            "last_activity": self._data["last_activity"]
        }
    
    def get_week_stats(self) -> Dict:
        """Get this week's statistics"""
        stats = defaultdict(int)
        today = datetime.now().date()
        
        for i in range(7):
            date = today - timedelta(days=i)
            file = ANALYTICS_DIR / DAILY_FILE_TEMPLATE.format(date=date.strftime("%Y-%m-%d"))
            if file.exists():
                try:
                    data = json.loads(file.read_text())
                    stats["total_commands"] += len(data.get("commands", []))
                    stats["total_projects"] += len(data.get("projects_opened", []))
                    stats["total_ai_queries"] += data.get("ai_queries", 0)
                    stats["active_days"] += 1
                except:
                    pass
        
        return dict(stats)
    
    def get_command_breakdown(self) -> Dict[str, int]:
        """Get command usage breakdown"""
        breakdown = defaultdict(int)
        for cmd in self._data["commands"]:
            # Extract base command
            base = cmd["command"].split()[0] if cmd["command"] else "unknown"
            breakdown[base] += 1
        return dict(breakdown)

def get_tracker() -> ActivityTracker:
    return ActivityTracker()
