# Design du CLI Claude Code d'Anthropic

**Date de vérification:** 2025-12-15  
**Version documentée:** Claude Code 1.0+ (versions jusqu'à 2.0.52+)

---

## Table des matières

1. [Vue d'ensemble](#vue-densemble)
2. [Installation et Authentification](#installation-et-authentification)
3. [Commandes principales](#commandes-principales)
4. [Commandes slash (Slash Commands)](#commandes-slash-slash-commands)
5. [Composants UI du CLI](#composants-ui-du-cli)
6. [Design Visuel et Couleurs](#design-visuel-et-couleurs)
7. [Principes UX/Design](#principes-uxdesign)
8. [Architecture Technique](#architecture-technique)
9. [Éléments non publiquement vérifiables](#éléments-non-publiquement-vérifiables)
10. [Sources et Références](#sources-et-références)

---

## Vue d'ensemble

### Qu'est-ce que Claude Code ?

Claude Code est un outil de développement "agentic" créé par Anthropic qui fonctionne directement dans le terminal. Il permet aux développeurs d'interagir avec l'IA Claude pour automatiser des tâches de développement, comprendre des bases de code et gérer des workflows Git, le tout via des commandes en langage naturel.

**Caractéristiques principales:**
- Interface en ligne de commande (CLI) interactive
- Compréhension contextuelle du code et de l'architecture du projet
- Exécution de commandes shell, édition de fichiers, et gestion Git
- Support pour plusieurs modèles Claude (Opus, Sonnet, Haiku)
- Extensibilité via plugins, commandes personnalisées et MCP (Model Context Protocol)
- Intégration avec VS Code, JetBrains et autres IDE

**Référentiel officiel:** [github.com/anthropics/claude-code](https://github.com/anthropics/claude-code)  
**Documentation officielle:** [code.claude.com/docs](https://code.claude.com/docs/en/overview)

---

## Installation et Authentification

### Méthodes d'installation

**macOS/Linux:**
```bash
curl -fsSL https://claude.ai/install.sh | bash
```

**macOS (Homebrew):**
```bash
brew install --cask claude-code
```

**Windows:**
```powershell
irm https://claude.ai/install.ps1 | iex
```

**NPM (multi-plateforme):**
```bash
npm install -g @anthropic-ai/claude-code
```
*Note: Nécessite Node.js 18+*

### Authentification

**Deux méthodes principales:**

1. **Authentification par navigateur** (recommandée pour les abonnés Pro/Max):
   - Lancez `claude` dans votre projet
   - Une page d'authentification s'ouvre automatiquement dans le navigateur
   - Connexion via compte claude.ai

2. **Clé API Anthropic** (pour usage pay-as-you-go):
   ```bash
   export ANTHROPIC_API_KEY="votre_clé_api"
   ```
   - Obtenir une clé depuis [console.anthropic.com](https://console.anthropic.com)
   - À ajouter dans `.bashrc`, `.zshrc`, ou fichier de configuration shell

**Première utilisation:**
```bash
cd votre-projet
claude
```

---

## Commandes principales

### Syntaxe de base

```bash
claude [options] ["requête en langage naturel"]
```

### Commandes du shell

| Commande | Description | Exemple |
|----------|-------------|---------|
| `claude` | Démarre une session interactive | `claude` |
| `claude "query"` | Exécute une requête directe | `claude "Fixe ce bug"` |
| `claude -p "query"` | Mode headless/print (non-interactif) | `claude -p "Analyse ce fichier"` |
| `claude -c` | Continue la session la plus récente | `claude -c` |
| `claude --resume <id>` ou `-r` | Reprend une session spécifique par ID | `claude -r "abc123" "Continue ce PR"` |
| `claude update` | Met à jour vers la dernière version | `claude update` |
| `claude mcp` | Configure les serveurs MCP | `claude mcp` |

### Flags CLI principaux

| Flag | Description | Exemple |
|------|-------------|---------|
| `--add-dir` | Ajoute des répertoires de travail supplémentaires | `claude --add-dir ../apps ../lib` |
| `--allowedTools` | Liste d'outils autorisés sans demande de permission | `claude --allowedTools "Bash(git log:*)"` |
| `--system-prompt` | Remplace complètement le prompt système par défaut | `claude --system-prompt "You are..."` |
| `--append-system-prompt` | Ajoute au prompt système (recommandé) | `claude --append-system-prompt "Focus on security"` |
| `--system-prompt-file` | Charge un prompt système depuis un fichier | `claude --system-prompt-file prompt.txt` |
| `--model` ou `-m` | Spécifie le modèle Claude à utiliser | `claude -m sonnet` |
| `--output-format` | Format de sortie (json, stream-json) | `claude -p --output-format json "query"` |
| `--debug` | Active le mode debug détaillé | `claude --debug` |

### Mode print (headless/automatisation)

Le mode print (`-p`) est conçu pour l'automatisation, les scripts CI/CD et les bots SRE:

```bash
# Sortie simple
claude -p "Analyse ces erreurs"

# Sortie JSON avec métadonnées
claude -p --output-format json "Analyse ce code"

# Streaming JSON (messages en temps réel)
claude -p --output-format stream-json "Génère un rapport"

# Avec customisation du system prompt
claude -p "Analyse ces erreurs" \
  --append-system-prompt "You are an SRE expert" \
  --output-format json \
  --allowedTools "Bash,Read,mcp__datadog"
```

---

## Commandes slash (Slash Commands)

Les commandes slash sont des raccourcis internes utilisables durant une session interactive Claude Code.

### Commandes slash natives (built-in)

| Commande | Description |
|----------|-------------|
| `/help` | Affiche la liste complète des commandes disponibles (natives + custom + MCP) |
| `/config` | Ouvre l'interface de configuration interactive |
| `/clear` | Efface l'historique de conversation et démarre une nouvelle session |
| `/compact` | Résume et compacte l'historique pour libérer de l'espace dans la fenêtre de contexte |
| `/context` | Affiche l'utilisation actuelle des tokens et le contexte |
| `/model` | Change le modèle Claude en cours (Opus/Sonnet/Haiku) |
| `/exit` ou `/quit` | Quitte la session interactive |
| `/status` | Affiche un panneau de statut complet (usage, session, configuration) |
| `/export` | Exporte la session actuelle (markdown, JSONL) |
| `/init` | Initialise le contexte de projet avec structure et documentation |
| `/hooks` | Ouvre l'interface de configuration des hooks |
| `/vim` | Active/désactive le mode Vim keybindings |
| `/terminal-setup` | Configure automatiquement le terminal (ex: Shift+Enter pour nouvelles lignes) |
| `/mcp` | Configure les serveurs Model Context Protocol |
| `/install-github-app` | Installe l'app GitHub pour revue automatique de PRs |
| `/agent` | Ouvre l'interface de création/gestion de subagents |
| `/bug` | Rapporte un bug directement à Anthropic |
| `/upgrade` | Interface pour upgrader l'abonnement Claude |
| `/plugin` | Gère les plugins Claude Code |

### Commandes slash personnalisées

Les utilisateurs peuvent créer leurs propres commandes slash en plaçant des fichiers Markdown dans:
- **Projet:** `.claude/commands/` (partagées avec l'équipe)
- **Personnel:** `~/.claude/commands/` (disponibles partout)

**Structure d'une commande personnalisée:**

```markdown
---
description: Brève description de la commande
allowed-tools: Read, Grep, Glob, Bash(git:*)
model: claude-sonnet-4
argument-hint: [filename]
---

# Instructions détaillées pour Claude

Analysez le fichier spécifié: $ARGUMENTS

Effectuez les étapes suivantes:
1. Lecture du fichier
2. Analyse de la structure
3. Suggestions d'amélioration
```

**Exemple d'utilisation:**
```bash
# Dans Claude Code:
/optimize src/auth.js
```

**Caractéristiques:**
- Utilisation de `$ARGUMENTS` pour arguments dynamiques
- Arguments positionnels: `$1`, `$2`, `$3`, etc.
- Exécution bash: `` !`command` ``
- Références de fichiers: `@filename`
- Organisation par sous-répertoires pour namespacing

---

## Composants UI du CLI

### Interface interactive

Claude Code utilise une interface textuelle riche (TUI - Terminal User Interface) avec plusieurs composants:

#### 1. **Zone de chat principale**
- Affichage du dialogue entre l'utilisateur et Claude
- Messages différenciés par rôle (user/assistant)
- Support du markdown dans les réponses

#### 2. **Indicateurs de statut**

**Progress indicators:**
- Barres de progression pour tâches longues
- Spinners d'attente durant l'exécution
- Indicateurs de complétion (`[████████████░░░░] 73.9%`)

**Session status:**
- Informations sur la session actuelle (ID, durée)
- Usage des tokens en temps réel
- Modèle actuellement utilisé

#### 3. **Messages d'erreur et avertissements**

Les messages d'erreur utilisent des codes de couleur pour la visibilité:
- Erreurs critiques en rouge
- Avertissements en jaune/orange
- Confirmations en vert

#### 4. **Prompts de permission**

Claude Code demande l'autorisation avant d'exécuter certaines actions:
- Édition de fichiers
- Exécution de commandes shell
- Création/suppression de fichiers

Format typique:
```
⚠️  Claude wants to edit: src/auth.js
   Allow? [y/n/always]
```

#### 5. **Indicateurs visuels d'activité**

**Exemples reconstitués:**
```
⏺ Bash(cd /path/project && git commit -m "message")
✓ File written successfully
⏺ Analyzing codebase...
```

**Symboles couramment utilisés:**
- `⏺` : Exécution en cours
- `✓` : Action réussie
- `✗` : Échec
- `⚠️` : Avertissement
- `💡` : Suggestion/conseil
- `🔍` : Recherche/analyse

#### 6. **Status line (barre de statut)**

Claude Code 2.0+ supporte des status lines personnalisables affichant:
- Modèle actuel
- Répertoire de travail
- Branche Git active
- Usage des tokens
- Timer de session (5h blocks)
- Métriques personnalisées

**Configuration via SDK ou outils tiers comme `ccstatusline`:**
- Style Powerline avec flèches et séparateurs
- Thèmes personnalisables (dark/light)
- Support multi-lignes (jusqu'à 3 lignes)
- Widgets personnalisés

#### 7. **Navigation dans l'historique**

- **Flèche haut:** Parcourt les messages précédents
- **Escape (2x):** Affiche liste de tous les messages précédents pour navigation rapide
- **Tab:** Autocomplétion des noms de fichiers et chemins

#### 8. **Indicateurs de contexte**

Affichage de fichiers et outils référencés:
```
📁 Referenced files:
  - src/auth.js
  - tests/auth.test.js

🔧 Tools used:
  - Read
  - Bash
```

---

## Design Visuel et Couleurs

### Palette de couleurs (basée sur chalk)

Claude Code utilise la bibliothèque `chalk` pour la coloration ANSI dans le terminal.

**Couleurs primaires identifiées:**

| Élément | Couleur Chalk | Code ANSI | Utilisation |
|---------|---------------|-----------|-------------|
| Erreur | `.red` / `.redBright` | `\x1b[31m` / `\x1b[91m` | Messages d'erreur critiques |
| Avertissement | `.yellow` / `.bgRed.white` | `\x1b[33m` / `\x1b[41m\x1b[37m` | Warnings, tips startup |
| Succès | `.green` / `.greenBright` | `\x1b[32m` / `\x1b[92m` | Confirmations, actions réussies |
| Info | `.blue` / `.blueBright` | `\x1b[34m` / `\x1b[94m` | Informations générales |
| Cyan | `.cyan` | `\x1b[36m` | Liens, références |
| Gris/Dim | `.gray` / `.dim` | `\x1b[90m` / `\x1b[2m` | Métadonnées, contexte |

**Styles de texte:**
- `.bold` : Texte en gras (`\x1b[1m`)
- `.italic` : Italique (`\x1b[3m`)
- `.underline` : Souligné (`\x1b[4m`)
- `.reset` : Réinitialisation (`\x1b[0m`)

**Note importante:** Un bug identifié (#1341) montre que les couleurs de background (`bgRed`, `bgWhite`) peuvent "saigner" dans la sortie suivante si non terminées avec `.reset()`.

### Support de couleurs

**Niveaux de couleur:**
- **ANSI basic (16 couleurs):** Mode par défaut
- **256 couleurs:** Supporté via détection automatique
- **Truecolor (24-bit RGB):** Activé avec `COLORTERM=truecolor`

**Exemple d'activation truecolor:**
```bash
export COLORTERM=truecolor
```

Vérification du support:
```bash
tput colors  # Affiche le nombre de couleurs supportées (>= 256 recommandé)
```

### Thèmes du CLI

Claude Code propose **6 thèmes prédéfinis** (depuis v1.0.3):

1. **Dark mode** (défaut)
2. **Light mode**
3. **Dark mode (colorblind-friendly)**
4. **Light mode (colorblind-friendly)**
5. **Dark mode (ANSI colors only)**
6. **Light mode (ANSI colors only)**

**Configuration:**
```
/config
# Puis naviguer vers Appearance > Theme
```

**Limitations connues:**
- Pas de support pour thèmes complètement personnalisés (demandé dans issue #1302)
- Les thèmes peuvent ne pas respecter les couleurs du terminal configurées par l'utilisateur
- Requête communautaire pour support de formats type base16, iTerm2

### Hiérarchie visuelle

**Principe de design:**
- Messages utilisateur: Style normal, couleur neutre
- Messages Claude: Légère indentation ou préfixe visuel
- Code: Blocs délimités avec backticks markdown, syntax highlighting si supporté par terminal
- Commandes exécutées: Préfixe `⏺` ou `$`, couleur cyan/blue
- Résultats: Indentés, couleur gris/dim pour différenciation
- Erreurs: Rouge vif, préfixe `✗` ou `ERROR:`

---

## Principes UX/Design

### 1. **Lisibilité et clarté**

- **Messages concis:** Claude évite la verbosité, présente l'essentiel
- **Formatage structuré:** Utilisation de listes, tableaux Markdown dans les réponses
- **Séparation visuelle:** Espaces et lignes pour aérer le contenu
- **Troncature intelligente:** Les longues sorties de commandes peuvent être résumées

### 2. **Feedback utilisateur continu**

- **Indicateurs de progression:** Spinners et barres pour opérations longues
- **Confirmations explicites:** "File written successfully", "Tests passed"
- **Messages d'erreur informatifs:** Explication de l'erreur + suggestions de résolution
- **Temps réel:** Streaming des réponses pour feedback instantané

### 3. **Gestion du contexte**

**Optimisation de la fenêtre de contexte:**
- Commande `/compact`: Résumé automatique de l'historique
- Commande `/clear`: Reset complet pour nouveau départ
- Affichage usage tokens via `/context` ou `/status`
- Stratégies de conservation: Extraction du contexte essentiel, suppression du superflu

**Best practice:** Compacter à chaque checkpoint naturel (feature complétée, bug fixé, commit effectué).

### 4. **Permissions et sécurité**

**Contrôle utilisateur:**
- Demandes de permission avant actions sensibles (édition, shell, suppression)
- Options: `y` (oui), `n` (non), `always` (toujours autoriser cet outil/action)
- Configuration via `allowedTools` pour pré-approuver certains outils
- Fichiers `.claude/settings.json` pour définir permissions par projet

**Exemple de configuration:**
```json
{
  "allowedTools": [
    "Bash(git log:*)",
    "Bash(npm test:*)",
    "Read",
    "Grep"
  ]
}
```

### 5. **Accessibilité**

**Support clavier:**
- Raccourcis intuitifs (Escape pour arrêter, Up/Down pour navigation)
- Mode Vim disponible (`/vim`) pour utilisateurs avancés
- Pas de dépendance à la souris

**Thèmes colorblind-friendly:**
- Palettes adaptées pour daltonisme
- Reliance sur symboles en plus des couleurs (✓, ✗, ⚠️)

**Support écrans:**
- Adaptation à différentes tailles de terminal
- Wrapping automatique du texte
- Pas de hard-coded widths (responsive)

### 6. **Interactivité et aide contextuelle**

- **`/help`:** Accès immédiat à toutes les commandes disponibles
- **Tips au démarrage:** Messages utiles affichés lors du lancement
- **Error recovery:** Suggestions automatiques en cas d'échec
- **Documentation inline:** Descriptions claires dans `/help`

### 7. **Cohérence**

- **Langage uniforme:** Terminologie cohérente (session, compact, context, etc.)
- **Patterns répétés:** Structure similaire pour toutes les commandes slash
- **Prévisibilité:** Comportements attendus et documentés

---

## Architecture Technique

### Stack technique identifié

**Langages:**
- **TypeScript** (34.0%): Langage principal
- **Python** (25.2%): Composants et outils
- **Shell** (22.5%): Scripts d'installation et automatisation
- **PowerShell** (12.4%): Support Windows
- **Dockerfile** (5.9%): Conteneurisation

**Bibliothèques et dépendances:**
- **Chalk:** Coloration terminal ANSI
- **Anthropic TS SDK:** Communication avec l'API Claude (`beta.messages.create`)
- **Node.js 18+:** Runtime
- **Possible:** Ink (React pour terminal) pour composants UI complexes (non confirmé)

### Architecture client-serveur

**Communication:**
1. **CLI (client):** Interface utilisateur dans le terminal
2. **Anthropic API (serveur):** Modèle Claude hébergé par Anthropic

**Flux:**
```
[Utilisateur] → [Claude Code CLI] → [Anthropic API] → [Modèle Claude]
                     ↓                       ↓
              [Filesystem]             [Réponse]
              [Git, Shell]                 ↓
                     ← ← ← ← ← ← ← ← ← ← ← ←
```

**Caractéristiques:**
- Requêtes API via `beta.messages.create` (Anthropic TS SDK)
- Pas de stockage local persistant des conversations (sauf `.claude/projects/`)
- Streaming des réponses pour affichage en temps réel

### Authentification et tokens

**Méthodes:**
1. **Browser-based OAuth:** Authentification via navigateur (abonnés Pro/Max)
2. **API Key:** Variable d'environnement `ANTHROPIC_API_KEY`

**Token refresh:**
- Sessions de 5 heures
- Refresh automatique pour les abonnés (browser auth)
- Gestion manuelle pour clés API

### Système de fichiers

**Répertoires principaux:**
- `~/.claude/`: Configuration globale utilisateur
  - `~/.claude/commands/`: Commandes slash personnelles
  - `~/.claude/CLAUDE.md`: Instructions globales pour Claude
  - `~/.claude/projects/`: Sessions et historique de projets
  
- `.claude/` (dans projet): Configuration spécifique au projet
  - `.claude/commands/`: Commandes slash du projet
  - `.claude/CLAUDE.md`: Instructions spécifiques au projet
  - `.claude/agents/`: Définitions de subagents
  - `.claude/hooks/`: Scripts de hooks (pre/post édition)
  - `.claude/settings.json`: Configuration projet

**Format de stockage:**
- Conversations: JSONL (JSON Lines)
- Configuration: JSON
- Commandes/agents: Markdown avec frontmatter YAML

### Hooks système

Les hooks permettent d'exécuter du code avant/après certaines actions:

**Types de hooks:**
- `pre_edit`: Avant édition de fichier
- `post_edit`: Après édition de fichier
- `pre_bash`: Avant exécution commande shell
- `post_bash`: Après exécution commande shell

**Exemple de hook (TypeScript/JavaScript):**
```javascript
// ~/.claude/hooks/format-on-save.js
export default async function postEdit({ filepath, content }) {
  // Exécuter Prettier sur le fichier
  const formatted = await runPrettier(filepath);
  return { filepath, content: formatted };
}
```

**Configuration:**
```
/hooks
# Interface interactive pour activer/configurer les hooks
```

### Model Context Protocol (MCP)

MCP permet d'étendre Claude Code avec outils et intégrations externes:

**Serveurs MCP courants:**
- **GitHub:** Accès repos, issues, PRs
- **Databases:** PostgreSQL, MySQL, SQLite
- **Browser automation:** Puppeteer, Playwright
- **File systems:** Accès étendu aux fichiers
- **APIs:** Intégrations tierces (Jira, Slack, etc.)

**Ajout d'un serveur MCP:**
```bash
claude mcp add <name> <command> [args...]
claude mcp add --transport stdio github --env GITHUB_TOKEN=xxx -- npx github-mcp
```

**Transports supportés:**
- `stdio`: Standard input/output
- `sse`: Server-Sent Events
- `http`: HTTP/REST

### Subagents (Multi-Agent)

Claude Code peut déléguer des tâches à des "subagents" spécialisés:

**Architecture:**
- **Main Agent:** Contexte principal de conversation
- **Sub Agent:** Agent spécialisé avec contexte isolé pour tâche spécifique
- **Tool "Task":** Permet d'invoquer un subagent depuis le Main Agent

**Avantages:**
- Optimisation du contexte principal (moins de tokens gaspillés)
- Spécialisation des tâches (ex: agent de debug, agent de review)
- Retour du résultat final uniquement au contexte principal

**Définition d'un subagent:**
```markdown
---
name: code-reviewer
description: Reviews code for best practices and bugs
model: claude-sonnet-4
allowed-tools: Read, Grep, Glob
---

You are a senior code reviewer specializing in security and performance.
Review the provided code and identify issues.
```

**Emplacement:**
- `.claude/agents/` (projet)
- `~/.claude/agents/` (personnel)

**Invocation:**
```
/agent
# ou automatiquement si configuré
@code-reviewer analyze src/auth.js
```

### SDK Claude Code

Anthropic propose un SDK TypeScript/JavaScript pour intégrer Claude Code programmatiquement:

```typescript
import { query } from "@anthropic-ai/claude-code";

for await (const message of query({ 
  prompt: "Analyze this file", 
  options: { maxTurns: 5 } 
})) {
  if (message.type === "assistant") {
    console.log(message.message);
  }
}
```

**Use cases:**
- Automatisation CI/CD
- Bots de développement
- Intégrations personnalisées

---

## Éléments non publiquement vérifiables

Cette section liste les éléments du design de Claude Code qui n'ont pas pu être trouvés dans les sources publiques disponibles.

### 1. Code source complet

**Status:** ❌ Non disponible

- Le dépôt GitHub [anthropics/claude-code](https://github.com/anthropics/claude-code) est public mais ne contient pas le code source de l'application CLI elle-même
- Contient uniquement: README, CHANGELOG, LICENSE, exemples, plugins
- Le code est distribué sous forme de binaire compilé/uglified

**Tentatives de reverse engineering:**
- Certains projets communautaires ont tenté d'analyser le code uglify (ex: `Yuyz0112/claude-code-reverse`)
- Ces tentatives ont été découragées ou retirées par demande d'Anthropic
- Approche v2 basée sur monitoring des requêtes API plutôt que décompilation

### 2. Codes couleurs HEX exacts

**Status:** ⚠️ Partiellement vérifiable

**Connu:**
- Utilisation de `chalk` pour ANSI colors
- Couleurs nommées (red, green, blue, yellow, cyan, gray)

**Non connu:**
- Valeurs HEX/RGB exactes pour thèmes personnalisés
- Palette complète des 6 thèmes prédéfinis
- Mapping précis ANSI → RGB pour truecolor

### 3. Typographie et polices

**Status:** ❌ Non spécifié

- Claude Code utilise la police du terminal configurée par l'utilisateur
- Pas de police imposée ou recommandée officiellement
- Pas de typo custom pour branding

### 4. Composants UI internes

**Status:** ⚠️ Reconstitués par observation

- Spinners, progress bars: Identifiés par issues et communauté
- Implémentation exacte (library utilisée): Non confirmée
- Possibilité: `ora`, `cli-progress`, ou custom

### 5. Limites techniques précises

**Status:** ⚠️ Partiellement documentées

**Non confirmés publiquement:**
- Taille maximale exacte de contexte (tokens) par session
- Limites de rate limiting API exactes
- Timeout par défaut des commandes shell
- Nombre maximum de hooks/subagents configurables

### 6. Telemetry et analytics

**Status:** ⚠️ Mentions générales seulement

**Documenté:**
- Collection de feedback (acceptation/rejet de code)
- Données de conversation associées
- Feedback via `/bug`

**Non détaillé:**
- Formats de données exactes collectées
- Fréquence d'envoi des metrics
- Endpoints de telemetry

### 7. Algorithmes internes

**Status:** ❌ Non divulgués

- Algorithme de compaction (`/compact`)
- Stratégies de sélection d'outils
- Logique de détection de contexte et résumé
- Prompt engineering interne exact

### 8. Performance et benchmarks

**Status:** ❌ Pas de métriques officielles

- Temps de réponse moyens par type de requête
- Latence réseau vs compute
- Benchmarks de vitesse Opus vs Sonnet vs Haiku
- Resource usage (CPU, RAM) du CLI

---

## Sources et Références

### Sources officielles Anthropic

1. **Documentation Claude Code (principale)**  
   URL: [https://code.claude.com/docs/en/overview](https://code.claude.com/docs/en/overview)  
   Contexte: Documentation officielle complète avec guides d'installation, référence CLI, commandes slash, MCP, SDK

2. **Documentation CLI Reference**  
   URL: [https://code.claude.com/docs/en/cli-reference](https://code.claude.com/docs/en/cli-reference)  
   Contexte: Référence complète des flags CLI, options, et commandes principales

3. **Documentation Slash Commands**  
   URL: [https://code.claude.com/docs/en/slash-commands](https://code.claude.com/docs/en/slash-commands)  
   Contexte: Documentation exhaustive des commandes slash natives et personnalisées

4. **Dépôt GitHub officiel**  
   URL: [https://github.com/anthropics/claude-code](https://github.com/anthropics/claude-code)  
   Contexte: README, changelog, examples, plugins - 42.4k stars, 2.8k forks

5. **Package NPM officiel**  
   URL: [https://www.npmjs.com/package/@anthropic-ai/claude-code](https://www.npmjs.com/package/@anthropic-ai/claude-code)  
   Contexte: Distribution NPM avec instructions d'installation

6. **Page produit Claude Code**  
   URL: [https://www.claude.com/product/claude-code](https://www.claude.com/product/claude-code)  
   Contexte: Présentation marketing et fonctionnalités principales

### Issues GitHub officielles

7. **Issue #1341: Background color bleed (chalk)**  
   URL: [https://github.com/anthropics/claude-code/issues/1341](https://github.com/anthropics/claude-code/issues/1341)  
   Contexte: Bug détaillant l'utilisation de chalk pour couleurs et problèmes ANSI

8. **Issue #1302: Custom terminal themes**  
   URL: [https://github.com/anthropics/claude-code/issues/1302](https://github.com/anthropics/claude-code/issues/1302)  
   Contexte: Liste complète des 6 thèmes prédéfinis, requête pour customisation

9. **Issue #2686: Terminal progress bars/spinners**  
   URL: [https://github.com/anthropics/claude-code/issues/2686](https://github.com/anthropics/claude-code/issues/2686)  
   Contexte: Traitement des spinners et barres de progression VT100

10. **Issue #12405: Progress bar garish appearance**  
    URL: [https://github.com/anthropics/claude-code/issues/12405](https://github.com/anthropics/claude-code/issues/12405)  
    Contexte: Design de barre de progression dans version 2.0.52+

### Guides et tutoriels communautaires

11. **Shipyard Claude Code Cheatsheet**  
    URL: [https://shipyard.build/blog/claude-code-cheat-sheet/](https://shipyard.build/blog/claude-code-cheat-sheet/)  
    Contexte: Guide complet avec commandes, configuration, workflows, best practices

12. **First Principles: Complete Slash Commands Reference**  
    URL: [https://firstprinciplescg.com/resources/claude-code-slash-commands-the-complete-reference-guide/](https://firstprinciplescg.com/resources/claude-code-slash-commands-the-complete-reference-guide/)  
    Contexte: Liste exhaustive incluant commandes non documentées officiellement

13. **Builder.io: How I use Claude Code**  
    URL: [https://www.builder.io/blog/claude-code](https://www.builder.io/blog/claude-code)  
    Contexte: Best practices, tips, workflows réels, configuration hooks et commandes

14. **ClaudeLog (ressources communautaires)**  
    URL: [https://claudelog.com/](https://claudelog.com/)  
    Contexte: Collection de MCP servers, plugins, status line formatters (ccstatusline)

15. **Awesome Claude Code (GitHub)**  
    URL: [https://github.com/hesreallyhim/awesome-claude-code](https://github.com/hesreallyhim/awesome-claude-code)  
    Contexte: Liste curée de commandes, workflows, plugins, hooks de la communauté

### Outils et extensions tiers

16. **ccstatusline (GitHub)**  
    URL: [https://github.com/sirmalloc/ccstatusline](https://github.com/sirmalloc/ccstatusline)  
    Contexte: Status line customizable avec powerline, thèmes, widgets - insights sur design visuel

17. **Claude Code UI (siteboon)**  
    URL: [https://github.com/siteboon/claudecodeui](https://github.com/siteboon/claudecodeui)  
    Contexte: Interface web/mobile pour Claude Code - détails sur composants UI

18. **Claudia GUI**  
    URL: [https://claudia.so/](https://claudia.so/)  
    Contexte: Interface graphique pour Claude Code (Tauri, React) - design patterns UI

### Analyses techniques

19. **Medium: Fixing Claude Code's Remote Colors**  
    URL: Article par Martin Thorsen Ranang (Juin 2025)  
    Contexte: Détails techniques sur COLORTERM=truecolor et support couleur 24-bit

20. **AI Engineer Guide: Claude Code Prompts & Tools**  
    URL: [https://aiengineerguide.com/blog/claude-code-prompt/](https://aiengineerguide.com/blog/claude-code-prompt/)  
    Contexte: Analyse des prompts système et définitions d'outils

21. **Claude Code Reverse Engineering (Yuyz0112)**  
    URL: [https://github.com/Yuyz0112/claude-code-reverse](https://github.com/Yuyz0112/claude-code-reverse)  
    Contexte: Tentative d'analyse du CLI - insights sur architecture Sub Agent et API

---

## Notes finales

Ce document a été compilé à partir de sources publiquement disponibles le 2025-12-15. Claude Code étant en développement actif avec mises à jour fréquentes, certaines informations peuvent évoluer rapidement.

**Recommandations:**
- Toujours consulter la documentation officielle pour informations à jour
- Utiliser `/help` dans Claude Code pour liste complète des commandes disponibles
- Suivre le changelog officiel pour nouvelles fonctionnalités

**Limitations de cette documentation:**
- Pas d'accès au code source interne
- Certains détails techniques (couleurs exactes HEX, algorithmes internes) non disponibles publiquement
- Exemples de sortie CLI reconstitués à partir de descriptions et captures d'écran communautaires

**Pour aller plus loin:**
- Expérimenter directement avec Claude Code
- Rejoindre le Discord Claude Developers
- Contribuer aux projets open-source communautaires

---

*Document généré dans le cadre d'une recherche technique exhaustive sur le design du CLI Claude Code d'Anthropic.*