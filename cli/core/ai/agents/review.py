"""
B.DEV CLI - Review Agent
AI-powered code review
"""
from pathlib import Path
from typing import Generator

from core.ai.agents.base import BaseAgent, AgentTask, AgentResult, AgentStatus, register_agent

@register_agent
class ReviewAgent(BaseAgent):
    """Agent for code review"""
    
    name = "review"
    description = "Analyse et review de code avec IA"
    
    def execute(self, task: AgentTask) -> Generator[str, None, AgentResult]:
        """Execute code review on target file or project"""
        self.status = AgentStatus.THINKING
        
        target = task.target
        if not target:
            yield "[ERROR] Aucune cible spécifiée"
            return AgentResult(success=False, message="No target specified")
        
        if not target.exists():
            yield f"[ERROR] {target} n'existe pas"
            return AgentResult(success=False, message=f"{target} does not exist")
        
        # Collect code to review
        yield f"📂 Analyse de {target.name}...\n"
        
        code_content = ""
        files_reviewed = []
        
        if target.is_file():
            # Single file
            code_content = target.read_text(encoding='utf-8', errors='replace')
            files_reviewed.append(target.name)
        else:
            # Directory - get key files
            extensions = {'.py', '.js', '.ts', '.jsx', '.tsx', '.vue', '.php'}
            for f in target.rglob('*'):
                if f.suffix in extensions and 'node_modules' not in str(f) and '.git' not in str(f):
                    try:
                        content = f.read_text(encoding='utf-8', errors='replace')
                        if len(code_content) + len(content) < 10000:  # Limit size
                            code_content += f"\n\n=== {f.relative_to(target)} ===\n{content}"
                            files_reviewed.append(str(f.relative_to(target)))
                    except:
                        pass
        
        if not code_content:
            yield "[WARNING] Aucun fichier de code trouvé"
            return AgentResult(success=False, message="No code files found")
        
        yield f"📝 {len(files_reviewed)} fichier(s) à analyser\n"
        yield "🤖 Review en cours...\n\n"
        
        # Build review prompt
        strict = task.options.get('strict', False)
        focus = task.options.get('focus', 'general')
        
        prompt = f"""You are an expert code reviewer. Review the following code thoroughly.

Focus areas: {focus}
Strictness: {'Very strict, catch everything' if strict else 'Normal, focus on important issues'}

For each issue found, provide:
1. **Severity**: 🔴 Critical, 🟠 Warning, 🟡 Suggestion
2. **Location**: File and approximate location
3. **Issue**: What's wrong
4. **Fix**: How to fix it

Also provide:
- Overall code quality score (1-10)
- Best practices followed
- Improvement recommendations

CODE TO REVIEW:
{code_content[:8000]}
"""
        
        self.status = AgentStatus.EXECUTING
        
        # Stream response
        full_response = ""
        for chunk in self.engine.chat(prompt, stream=True):
            full_response += chunk
            yield chunk
        
        self.status = AgentStatus.COMPLETED
        
        yield "\n\n✅ Review terminée!"
        
        return AgentResult(
            success=True,
            message="Review completed",
            data={"files": files_reviewed, "response": full_response}
        )
