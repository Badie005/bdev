"""
B.DEV CLI - AI Agents Base
Foundation for autonomous AI agents
"""
from abc import ABC, abstractmethod
from typing import Dict, Any, Optional, Generator
from pathlib import Path
from dataclasses import dataclass
from enum import Enum

from core.ai.engine import get_ai_engine
from core.session import get_session

class AgentStatus(Enum):
    IDLE = "idle"
    THINKING = "thinking"
    EXECUTING = "executing"
    COMPLETED = "completed"
    FAILED = "failed"

@dataclass
class AgentTask:
    """Represents a task for an agent"""
    name: str
    description: str
    target: Optional[Path] = None
    options: Dict[str, Any] = None
    
    def __post_init__(self):
        if self.options is None:
            self.options = {}

@dataclass
class AgentResult:
    """Result from an agent execution"""
    success: bool
    message: str
    data: Optional[Dict[str, Any]] = None
    changes: Optional[list] = None

class BaseAgent(ABC):
    """Base class for all AI agents"""
    
    name: str = "base"
    description: str = "Base agent"
    
    def __init__(self):
        self.engine = get_ai_engine()
        self.session = get_session()
        self.status = AgentStatus.IDLE
    
    @abstractmethod
    def execute(self, task: AgentTask) -> Generator[str, None, AgentResult]:
        """Execute the agent's task. Yields progress messages."""
        pass
    
    def _ask_ai(self, prompt: str, context: str = "") -> str:
        """Helper to ask the AI engine"""
        response = ""
        for chunk in self.engine.chat(prompt, context=context, stream=False):
            response += chunk
        return response.strip()
    
    def _stream_ai(self, prompt: str, context: str = "") -> Generator[str, None, str]:
        """Helper to stream AI response"""
        full = ""
        for chunk in self.engine.chat(prompt, context=context, stream=True):
            full += chunk
            yield chunk
        return full

class AgentRegistry:
    """Registry for available agents"""
    
    _agents: Dict[str, type] = {}
    
    @classmethod
    def register(cls, agent_class: type):
        """Register an agent class"""
        cls._agents[agent_class.name] = agent_class
    
    @classmethod
    def get(cls, name: str) -> Optional[BaseAgent]:
        """Get an agent instance by name"""
        if name in cls._agents:
            return cls._agents[name]()
        return None
    
    @classmethod
    def list_agents(cls) -> Dict[str, str]:
        """List all registered agents"""
        return {name: agent.description for name, agent in cls._agents.items()}

def register_agent(cls):
    """Decorator to register an agent"""
    AgentRegistry.register(cls)
    return cls
