"""
B.DEV CLI Core - AI Engine
Ollama wrapper with streaming and memory
"""
import subprocess
from typing import Generator, Optional, List, Dict
from pathlib import Path

from core.config import get_config
from core.session import get_session

class AIEngine:
    """Unified AI interface with memory support"""
    
    def __init__(self):
        self.config = get_config()
        self.session = get_session()
        self.model = self.config.get("ai.model", "codellama:7b")
        self.fallback = self.config.get("ai.fallback_model", "phi3:mini")
    
    def _build_prompt_with_memory(self, user_message: str, context: str = "") -> str:
        """Build prompt including conversation history"""
        memory = self.session.get_ai_context()
        
        parts = []
        
        # System instruction
        parts.append("You are B.AI, a helpful developer assistant for the B.DEV CLI.")
        parts.append("Answer concisely and accurately. Use the conversation history for context.")
        parts.append("")
        
        # Add context if provided
        if context:
            parts.append("=== PROJECT CONTEXT ===")
            parts.append(context)
            parts.append("")
        
        # Add conversation history
        if memory:
            parts.append("=== CONVERSATION HISTORY ===")
            for msg in memory[-10:]:  # Last 10 messages
                role = "User" if msg["role"] == "user" else "Assistant"
                parts.append(f"{role}: {msg['content'][:500]}")  # Truncate for token limit
            parts.append("")
        
        # Current question
        parts.append("=== CURRENT QUESTION ===")
        parts.append(f"User: {user_message}")
        parts.append("")
        parts.append("Assistant:")
        
        return "\n".join(parts)
    
    def chat(self, message: str, context: str = "", stream: bool = True) -> Generator[str, None, None]:
        """Send message to AI and yield response chunks"""
        
        # Build prompt with memory
        full_prompt = self._build_prompt_with_memory(message, context)
        
        # Save user message to memory
        self.session.add_ai_message("user", message)
        
        response_text = ""
        
        try:
            process = subprocess.Popen(
                ["ollama", "run", self.model, full_prompt],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                encoding='utf-8',
                errors='replace',
                bufsize=1
            )
            
            # Stream output
            while True:
                char = process.stdout.read(1)
                if not char and process.poll() is not None:
                    break
                if char:
                    response_text += char
                    if stream:
                        yield char
            
            # Save assistant response to memory
            self.session.add_ai_message("assistant", response_text.strip())
            
            if not stream:
                yield response_text
                
        except FileNotFoundError:
            yield "[ERROR] Ollama n'est pas installé."
        except Exception as e:
            yield f"[ERROR] {str(e)}"
    
    def clear_memory(self):
        """Clear conversation history"""
        self.session.clear_ai_context()
    
    def get_memory_summary(self) -> str:
        """Get summary of conversation memory"""
        memory = self.session.get_ai_context()
        if not memory:
            return "Aucune conversation en mémoire."
        return f"{len(memory)} messages en mémoire (depuis {memory[0].get('timestamp', 'inconnu')})"

# Singleton
_engine: Optional[AIEngine] = None

def get_ai_engine() -> AIEngine:
    global _engine
    if _engine is None:
        _engine = AIEngine()
    return _engine
