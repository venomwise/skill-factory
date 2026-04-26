# Cross-Agent Skill Architecture

This document defines the standard architecture for skills that need to be shared across multiple AI agents (pi, codex, Claude Code, etc.).

## Problem Statement

When the same skill is used by multiple AI agents, traditional approaches lead to:
- **Duplication**: Each agent has its own copy of the skill
- **Inconsistency**: Updates must be synchronized across multiple locations
- **Wasted resources**: Multiple virtual environments for the same dependencies
- **Configuration sprawl**: API keys and settings duplicated in multiple places

## Solution: Unified Skill Directory

Skills should be installed in a **shared, agent-agnostic location** that all AI agents can reference.

### Directory Structure

```
~/.AI-Skills/                          # Shared skill directory
├── <skill-name>/
│   ├── .venv/                         # Python virtual environment (if needed)
│   ├── scripts/                       # Executable scripts
│   ├── references/                    # Documentation
│   ├── assets/                        # Templates, data files
│   ├── config.json                    # Default configuration (can be committed)
│   ├── config.local.json              # User overrides (gitignored)
│   ├── .gitignore                     # Protect sensitive files
│   └── SKILL.md                       # Skill definition
└── README.md                          # Architecture documentation

~/.config/ai-skills/                   # XDG-compliant shared configuration
├── <skill-name>.json                  # Cross-agent configuration (API keys, etc.)
└── ...
```

### Agent Integration

Each AI agent creates a symbolic link to the shared skill:

```bash
# Pi agent
ln -s ~/.AI-Skills/<skill-name> ~/.pi/agent/skills/<skill-name>

# Codex
ln -s ~/.AI-Skills/<skill-name> ~/.codex/skills/<skill-name>

# Claude Code
ln -s ~/.AI-Skills/<skill-name> ~/.claude/skills/<skill-name>
```

This allows agents to discover skills without duplication.

## Configuration Architecture

### Configuration Priority (Highest to Lowest)

1. **Command-line arguments**: `--api-key`, `--config`, etc.
2. **Environment variables**: `SKILL_API_KEY`, `SKILL_CONFIG_PATH`, etc.
3. **Skill-specific local override**: `~/.AI-Skills/<skill-name>/config.local.json`
4. **Skill-specific defaults**: `~/.AI-Skills/<skill-name>/config.json`
5. **Cross-agent shared config**: `~/.config/ai-skills/<skill-name>.json` ⭐ Recommended for API keys
6. **Legacy locations**: Agent-specific config directories (for backward compatibility)

### Configuration File Purposes

**`~/.config/ai-skills/<skill-name>.json`** (Recommended for API keys)
- Cross-agent shared configuration
- API keys, credentials, tokens
- User-level settings that apply to all agents
- Follows XDG Base Directory Specification

**`~/.AI-Skills/<skill-name>/config.json`** (Default settings)
- Default configuration that can be committed to version control
- Base URLs, timeout values, default parameters
- No sensitive information
- Distributed with the skill

**`~/.AI-Skills/<skill-name>/config.local.json`** (User overrides)
- User-specific overrides for this skill
- Gitignored
- Takes precedence over config.json
- Use when you need skill-specific settings that differ from cross-agent defaults

### Example Configuration Setup

**Cross-agent shared (API keys):**
```bash
# ~/.config/ai-skills/exa-search.json
{
  "profiles": [
    { "id": "main", "api_key": "your-api-key-here" },
    { "id": "backup", "api_key": "backup-key-here" }
  ]
}
```

**Skill defaults (no sensitive data):**
```bash
# ~/.AI-Skills/exa-search/config.json
{
  "base_url": "https://api.exa.ai",
  "timeout_seconds": 30,
  "default_num_results": 5
}
```

**User overrides (optional):**
```bash
# ~/.AI-Skills/exa-search/config.local.json
{
  "timeout_seconds": 60,
  "default_num_results": 10
}
```

## Implementation Guidelines

### For Skills with Python Scripts

**1. Virtual Environment Location**

Always create the virtual environment in the skill directory:

```bash
python3 -m venv ~/.AI-Skills/<skill-name>/.venv
```

**2. Script Paths**

Use absolute paths to the shared location:

```bash
~/.AI-Skills/<skill-name>/.venv/bin/python ~/.AI-Skills/<skill-name>/scripts/script.py
```

**3. Configuration Loading**

Implement configuration priority in your scripts:

```python
def _default_config_paths() -> list[str]:
    root = _skill_root()
    home = os.path.expanduser("~")
    return [
        # Legacy fallback (agent-specific)
        os.path.join(home, ".codex", "config", "<skill-name>.json"),
        # XDG standard (cross-agent shared) ⭐
        os.path.join(home, ".config", "ai-skills", "<skill-name>.json"),
        # Skill defaults
        os.path.join(root, "config.json"),
        # User overrides
        os.path.join(root, "config.local.json"),
    ]
```

**4. .gitignore**

Always include a `.gitignore` in the skill directory:

```gitignore
.venv/
config.local.json
__pycache__/
*.pyc
*.pyo
*.pyd
.DS_Store
```

### For Skills with Configuration Files

**1. Separate Sensitive from Non-Sensitive**

- **config.json**: Default settings, can be committed
- **config.local.json**: User overrides, gitignored
- **~/.config/ai-skills/<skill>.json**: API keys, gitignored by default

**2. Document Configuration Locations**

In your SKILL.md, clearly document where configuration should go:

```markdown
## Setup

**1. Install dependencies:**
```bash
~/.AI-Skills/<skill-name>/.venv/bin/pip install -r ~/.AI-Skills/<skill-name>/requirements.txt
```

**2. Configure API key (recommended - shared across all agents):**

Create `~/.config/ai-skills/<skill-name>.json`:
```json
{
  "api_key": "YOUR_API_KEY"
}
```

Alternatively, create `~/.AI-Skills/<skill-name>/config.local.json` for skill-specific configuration.
```

**3. Provide Configuration Reference**

Create a `references/CONFIG.md` that documents:
- Configuration priority
- All available options
- Example configurations
- Troubleshooting

### For Skills Without External Dependencies

Skills that don't require Python packages or configuration files can still benefit from the shared architecture:

```
~/.AI-Skills/<skill-name>/
├── SKILL.md
├── references/
│   └── examples.md
└── assets/
    └── template.md
```

Agents link to this location, ensuring updates propagate to all agents.

## SKILL.md Updates

When documenting setup in SKILL.md, use the shared paths:

**Before (project-specific):**
```markdown
## Setup

```bash
python -m venv codex/<skill-name>/.venv
codex/<skill-name>/.venv/bin/pip install -r codex/<skill-name>/requirements.txt
```
```

**After (cross-agent shared):**
```markdown
## Setup

**1. The skill is installed at `~/.AI-Skills/<skill-name>/` (shared across all AI agents)**

**2. Create virtual environment:**
```bash
python3 -m venv ~/.AI-Skills/<skill-name>/.venv
```

**3. Install dependencies:**
```bash
~/.AI-Skills/<skill-name>/.venv/bin/pip install -r ~/.AI-Skills/<skill-name>/requirements.txt
```

**4. Configure API key (recommended - shared across all agents):**

Create `~/.config/ai-skills/<skill-name>.json`:
```json
{
  "api_key": "YOUR_API_KEY"
}
```
```

## Benefits

### For Users
- **Install once, use everywhere**: No need to set up the same skill multiple times
- **Unified configuration**: API keys in one place
- **Consistent behavior**: Same skill version across all agents
- **Easy updates**: Update once, all agents benefit

### For Developers
- **Single source of truth**: One codebase to maintain
- **Easier testing**: Test once, works for all agents
- **Clear separation**: Sensitive vs. non-sensitive configuration
- **Standard structure**: Predictable layout across skills

### For Teams
- **Shared skills**: Team members can share custom skills easily
- **Version control**: Track skill changes in one place
- **Onboarding**: New team members set up once

## Migration Guide

### Migrating Existing Skills

**1. Create shared directory:**
```bash
mkdir -p ~/.AI-Skills/<skill-name>
mkdir -p ~/.config/ai-skills
```

**2. Move skill files:**
```bash
cp -r <agent-dir>/skills/<skill-name>/* ~/.AI-Skills/<skill-name>/
```

**3. Create symbolic link:**
```bash
rm -rf <agent-dir>/skills/<skill-name>
ln -s ~/.AI-Skills/<skill-name> <agent-dir>/skills/<skill-name>
```

**4. Move configuration:**
```bash
# Extract API keys to shared config
mv ~/.AI-Skills/<skill-name>/config.local.json ~/.config/ai-skills/<skill-name>.json

# Keep only non-sensitive defaults in skill directory
# Edit ~/.AI-Skills/<skill-name>/config.json to remove sensitive data
```

**5. Update script paths:**
- Update SKILL.md to use `~/.AI-Skills/<skill-name>/` paths
- Update configuration loading code to check `~/.config/ai-skills/`

**6. Add .gitignore:**
```bash
cat > ~/.AI-Skills/<skill-name>/.gitignore << 'EOF'
.venv/
config.local.json
__pycache__/
*.pyc
.DS_Store
EOF
```

**7. Test:**
```bash
# Verify skill works from shared location
<agent> "test the <skill-name> skill"
```

## Best Practices

### DO
- ✅ Use `~/.AI-Skills/` for skill installation
- ✅ Use `~/.config/ai-skills/` for API keys and credentials
- ✅ Document configuration priority clearly
- ✅ Provide migration instructions for existing users
- ✅ Use absolute paths in SKILL.md examples
- ✅ Include .gitignore to protect sensitive files
- ✅ Test with multiple agents to ensure compatibility

### DON'T
- ❌ Hardcode agent-specific paths (e.g., `~/.pi/`, `~/.codex/`)
- ❌ Put API keys in config.json (use config.local.json or ~/.config/ai-skills/)
- ❌ Assume a specific agent is being used
- ❌ Use relative paths that depend on current directory
- ❌ Commit .venv/ or config.local.json to version control

## Real-World Example: exa-search

The `exa-search` skill demonstrates this architecture:

**Directory structure:**
```
~/.AI-Skills/exa-search/
├── .venv/                    # Python dependencies
├── scripts/
│   └── exa_search.py         # Main script
├── references/
│   ├── CONFIG.md             # Configuration guide
│   └── query-recipes.md      # Usage examples
├── config.json               # Default settings (no API key)
├── .gitignore                # Protects .venv/ and config.local.json
└── SKILL.md                  # Skill definition

~/.config/ai-skills/
└── exa-search.json           # API key (cross-agent shared)
```

**Configuration priority in exa_search.py:**
```python
def _default_config_paths() -> list[str]:
    root = _skill_root()
    home = os.path.expanduser("~")
    return [
        os.path.join(home, ".codex", "config", "exa-search.json"),  # Legacy
        os.path.join(home, ".config", "ai-skills", "exa-search.json"),  # Shared ⭐
        os.path.join(root, "config.json"),  # Defaults
        os.path.join(root, "config.local.json"),  # User overrides
    ]
```

**Agent integration:**
```bash
# Pi
ln -s ~/.AI-Skills/exa-search ~/.pi/agent/skills/exa-search

# Codex
ln -s ~/.AI-Skills/exa-search ~/.codex/skills/exa-search
```

**Result**: One installation, one API key, works with all agents.

## Summary

For skills with Python scripts or configuration:

1. **Install to** `~/.AI-Skills/<skill-name>/`
2. **Store API keys in** `~/.config/ai-skills/<skill-name>.json`
3. **Link from agent directories** using symbolic links
4. **Document paths** using absolute paths to `~/.AI-Skills/`
5. **Implement configuration priority** with cross-agent shared config
6. **Protect sensitive files** with .gitignore

This architecture ensures skills are reusable, maintainable, and consistent across all AI agents.
