# Skill Factory

A comprehensive skill for creating and optimizing AI agent skills with systematic workflows, automated testing, and quality assurance.

## Overview

This skill integrates the best practices from both `skill-creator` (Anthropic) and `skill-authoring` (internal) to provide:

- **Two distinct workflows**: Create new skills or optimize existing ones
- **Mandatory checkpoints**: Prevent skipping steps and ensure quality
- **Automated testing**: Evaluation framework with baseline comparison
- **Quality validation**: Essential and recommended quality checklists
- **Tool support**: Scripts for testing, benchmarking, and packaging

## Quick Start

### Identify Your Scenario

**Creating a new skill?**
→ Follow [Create Workflow](workflows/create-workflow.md)

**Improving an existing skill?**
→ Follow [Optimize Workflow](workflows/optimize-workflow.md)

## Structure

```
skill-factory/
├── SKILL.md                    # Main entry point with scenario routing
├── workflows/
│   ├── create-workflow.md      # 9-step creation process with checkpoints
│   └── optimize-workflow.md    # 10-step optimization process with checkpoints
├── guides/
│   ├── principles.md           # Core design principles
│   ├── patterns.md             # Common reusable patterns
│   └── anti-patterns.md        # What to avoid
├── scripts/                    # Automation tools
│   ├── package_skill.py        # Package .skill files
│   ├── aggregate_benchmark.py  # Aggregate test results
│   └── run_loop.py             # Description optimization
├── eval-viewer/                # Interactive result viewer
├── agents/                     # Subagent instructions
│   ├── grader.md              # Assertion evaluation
│   ├── comparator.md          # Blind A/B comparison
│   └── analyzer.md            # Benchmark analysis
├── references/
│   ├── schemas.md             # JSON schemas for evals
│   └── examples.md            # Real-world examples
└── assets/                     # Templates and resources
```

## Key Features

### 1. Clear Workflow Separation

**Create Workflow** (9 steps):
1. Capture Intent → Checkpoint: Intent confirmed
2. Interview & Research → Checkpoint: Requirements clear
3. Write SKILL.md → Checkpoint: Draft reviewed
4. Create Test Cases → Checkpoint: Tests approved
5. Run Evaluations → Checkpoint: All runs complete
6. Review Results → Checkpoint: Feedback collected
7. Iterate (if needed) → Checkpoint: Improvements done
8. Final Validation → Checkpoint: Quality checklist
9. Package & Deliver → Checkpoint: Skill packaged

**Optimize Workflow** (10 steps):
1. Analyze Current State → Checkpoint: Issues identified
2. Define Optimization Goal → Checkpoint: Goal confirmed
3. Create/Update Tests → Checkpoint: Tests ready
4. Snapshot Baseline → Checkpoint: Baseline captured
5. Apply Improvements → Checkpoint: Changes reviewed
6. Run Evaluations → Checkpoint: All runs complete
7. Compare Results → Checkpoint: Comparison done
8. Iterate (if needed) → Checkpoint: Improvements done
9. Final Validation → Checkpoint: Quality checklist
10. Update & Deliver → Checkpoint: Skill updated

### 2. Mandatory Checkpoints

Each step has a checkpoint with clear verification criteria. The AI cannot proceed until all checkpoint items are satisfied.

### 3. Quality Checklist

**Essential Quality (Must Pass):**
- Description complete (what + when)
- Main file < 500 lines
- At least 3 test cases
- All tests pass
- Terminology consistent

**Recommended Quality (Should Pass):**
- No time-sensitive info
- Examples concrete
- References organized
- Progressive disclosure
- Workflows clear

**Optimization-Specific:**
- Baseline comparison done
- Metrics improved
- No regressions
- User feedback addressed

### 4. Automated Tools

```bash
# Package a skill
python -m scripts.package_skill <path/to/skill>

# Aggregate benchmark results
python -m scripts.aggregate_benchmark <workspace>/iteration-N --skill-name <name>

# Optimize skill description
python -m scripts.run_loop --eval-set <path> --skill-path <path> --model <model-id>

# Launch evaluation viewer
python eval-viewer/generate_review.py <workspace>/iteration-N \
  --skill-name "my-skill" \
  --benchmark <workspace>/iteration-N/benchmark.json
```

## Design Principles

### 1. Conciseness
Only include what the AI doesn't already know. Challenge every sentence.

### 2. Progressive Disclosure
- Level 1: Metadata (always loaded)
- Level 2: Main SKILL.md (loaded when triggered)
- Level 3: References (loaded as needed)

### 3. Test-Driven Development
Create evaluations BEFORE writing extensive documentation.

### 4. Explain the Why
Modern LLMs are smart. Explain reasoning rather than rigid rules.

**For detailed principles, see:** [guides/principles.md](guides/principles.md)

## Platform Support

### Claude Code (Full Features)
✅ Subagents for parallel testing
✅ Browser-based eval viewer
✅ Description optimization
✅ All automation scripts

### Claude.ai (Limited)
❌ No subagents → Sequential testing
❌ No browser → Inline results
❌ No description optimization
✅ Packaging works

### Cowork (Adapted)
✅ Subagents available
⚠️ No browser → Use `--static` flag
⚠️ Feedback downloads as file
✅ Description optimization works

## Usage Example

```markdown
User: "I want to create a skill for analyzing SQL query performance"

AI: "We're creating a new skill. Let me follow the Create Workflow.

**Step 1: Capture Intent**

Let me understand what you need:
1. What should this skill enable Claude to do?
   - Analyze SQL queries for performance issues?
   - Suggest optimizations?
   - Explain execution plans?

2. When should it trigger?
   - When user mentions "SQL performance"?
   - When user pastes a slow query?
   - When user asks about query optimization?

3. Expected output format?
   - Analysis report?
   - Optimized query?
   - List of recommendations?

4. Should we create test cases?
   - Yes, this has objective outputs (performance metrics, suggestions)

Please confirm or clarify..."
```

## Differences from Original Skills

### vs skill-creator (Anthropic)
- ✅ Added clear workflow separation (create vs optimize)
- ✅ Added mandatory checkpoints at each step
- ✅ Added comprehensive quality checklist
- ✅ Reorganized into guides/ and workflows/
- ✅ Removed platform-specific complexity from main file

### vs skill-authoring (Internal)
- ✅ Added complete tooling support (scripts, eval-viewer)
- ✅ Added systematic workflows with checkpoints
- ✅ Added automated testing framework
- ✅ Removed cross-agent architecture (as requested)
- ✅ Simplified quality checklist (was 13 items, now tiered)

## Contributing

When improving this skill:

1. Test changes with real skill creation/optimization tasks
2. Ensure checkpoints remain clear and actionable
3. Keep workflows focused and step-by-step
4. Update quality checklist if adding new requirements
5. Maintain separation between create and optimize workflows

## License

Incorporates components from Anthropic's skill-creator (see LICENSE.txt in skill-creator/).
Internal components are proprietary.

## Additional Resources

- [Design Principles](guides/principles.md)
- [Common Patterns](guides/patterns.md)
- [Anti-Patterns](guides/anti-patterns.md)
- [JSON Schemas](references/schemas.md)
- [Examples](references/examples.md)
