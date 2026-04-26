---
name: skill-factory
description: >
  Create new AI agent skills from scratch or optimize existing skills through systematic workflows.
  Use when the user wants to create a new skill, improve an existing skill's quality or performance,
  run evaluations to test a skill, or optimize a skill's triggering accuracy. Provides step-by-step
  workflows with checkpoints to ensure nothing is missed.
---

# Skill Factory

A comprehensive skill for creating and optimizing AI agent skills with systematic workflows, automated testing, and quality assurance.

## Quick Start: Identify Your Scenario

Before starting, determine which workflow you need:

### Scenario 1: Creating a New Skill
**Use when:**
- User says "create a skill for X"
- User wants to capture a workflow into a reusable skill
- User describes a capability they want Claude to have
- No existing skill exists for this purpose

**Go to:** [Create Workflow](workflows/create-workflow.md)

### Scenario 2: Optimizing an Existing Skill
**Use when:**
- User says "improve this skill" or "optimize X skill"
- A skill exists but doesn't work well
- User wants better triggering accuracy
- User wants to add test cases or improve quality
- User reports issues with an existing skill

**Go to:** [Optimize Workflow](workflows/optimize-workflow.md)

**If unclear:** Ask the user: "Are we creating a new skill from scratch, or improving an existing one?"

---

## Core Principles

Before diving into workflows, understand these fundamental principles that guide all skill development:

### 1. Conciseness Over Completeness
The context window is shared. Only include what the AI doesn't already know. Challenge every sentence: "Does this justify its token cost?"

**Good (concise):**
```python
Use pdfplumber for text extraction:
import pdfplumber
with pdfplumber.open("file.pdf") as pdf:
    text = pdf.pages[0].extract_text()
```

**Bad (verbose):**
```
PDF files are a common format that contains text and images. 
To extract text, you need a library. There are many options...
```

### 2. Progressive Disclosure
Structure content in three layers:
- **Level 1: Metadata** (name + description) - Always loaded (~100 words)
- **Level 2: Main SKILL.md** - Loaded when triggered (<500 lines ideal)
- **Level 3: References** - Loaded as needed (unlimited)

### 3. Test-Driven Development
Create evaluations BEFORE writing extensive documentation. Let real usage guide what to include.

### 4. Explain the Why
Modern LLMs are smart. Explain reasoning rather than rigid rules. Avoid heavy-handed MUSTs when possible.

**For detailed principles, see:** [guides/principles.md](guides/principles.md)

---

## Workflow Overview

Both workflows follow a similar structure with mandatory checkpoints:

```
┌─────────────────────────────────────────────────────────────┐
│ CREATE WORKFLOW                                             │
├─────────────────────────────────────────────────────────────┤
│ 1. Capture Intent          → Checkpoint: Intent confirmed   │
│ 2. Interview & Research    → Checkpoint: Requirements clear │
│ 3. Write SKILL.md          → Checkpoint: Draft reviewed     │
│ 4. Create Test Cases       → Checkpoint: Tests approved     │
│ 5. Run Evaluations         → Checkpoint: All runs complete  │
│ 6. Review Results          → Checkpoint: Feedback collected │
│ 7. Iterate (if needed)     → Checkpoint: Improvements done  │
│ 8. Final Validation        → Checkpoint: Quality checklist  │
│ 9. Package & Deliver       → Checkpoint: Skill packaged     │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ OPTIMIZE WORKFLOW                                           │
├─────────────────────────────────────────────────────────────┤
│ 1. Analyze Current State   → Checkpoint: Issues identified  │
│ 2. Define Optimization Goal→ Checkpoint: Goal confirmed     │
│ 3. Create/Update Tests     → Checkpoint: Tests ready        │
│ 4. Snapshot Baseline       → Checkpoint: Baseline captured  │
│ 5. Apply Improvements      → Checkpoint: Changes reviewed   │
│ 6. Run Evaluations         → Checkpoint: All runs complete  │
│ 7. Compare Results         → Checkpoint: Comparison done    │
│ 8. Iterate (if needed)     → Checkpoint: Improvements done  │
│ 9. Final Validation        → Checkpoint: Quality checklist  │
│ 10. Update & Deliver       → Checkpoint: Skill updated      │
└─────────────────────────────────────────────────────────────┘
```

**Key Features:**
- ✅ Mandatory checkpoints prevent skipping steps
- ✅ Clear success criteria for each step
- ✅ Automated tools for testing and evaluation
- ✅ Quality validation before delivery

---

## Communication Guidelines

Users range from coding experts to beginners. Adapt your language:

**Default approach:**
- "evaluation" and "benchmark" are OK
- For "JSON" and "assertion", check if user understands before using
- Briefly explain technical terms when in doubt

**Context cues:**
- If user mentions "test cases" or "assertions" → they're technical
- If user says "make it work better" → use simpler language
- When unsure, provide a short definition: "assertions (automated checks)"

---

## Tools and Scripts

This skill includes automated tools from the `scripts/` directory:

### Testing & Evaluation
```bash
# Aggregate benchmark results
python -m scripts.aggregate_benchmark <workspace>/iteration-N --skill-name <name>

# Run description optimization loop
python -m scripts.run_loop --eval-set <path> --skill-path <path> --model <model-id>

# Package skill for distribution
python -m scripts.package_skill <path/to/skill-folder>
```

### Evaluation Viewer
```bash
# Launch interactive review interface
python eval-viewer/generate_review.py <workspace>/iteration-N \
  --skill-name "my-skill" \
  --benchmark <workspace>/iteration-N/benchmark.json
```

**For complete tool documentation, see:** [references/schemas.md](references/schemas.md)

---

## Quality Checklist

Before delivering any skill (create or optimize), verify these requirements:

### Essential Quality (Must Pass)
- [ ] **Description complete**: Includes what it does AND when to use it
- [ ] **Main file size**: SKILL.md < 500 lines (or has clear reference structure)
- [ ] **Test coverage**: At least 3 realistic test cases exist
- [ ] **Tests pass**: All test cases execute successfully
- [ ] **Terminology consistent**: Same terms used throughout

### Recommended Quality (Should Pass)
- [ ] **No time-sensitive info**: No "before August 2025" type content
- [ ] **Examples concrete**: Real, specific examples (not abstract)
- [ ] **References organized**: One level deep, with clear pointers
- [ ] **Progressive disclosure**: Heavy content moved to reference files
- [ ] **Workflows clear**: Step-by-step instructions are unambiguous

### Optimization-Specific (For Optimize Workflow)
- [ ] **Baseline comparison**: Performance compared to previous version
- [ ] **Metrics improved**: Pass rate, speed, or quality measurably better
- [ ] **No regressions**: Existing functionality still works
- [ ] **User feedback addressed**: All reported issues resolved

**Checkpoint:** Do not proceed to packaging until all Essential Quality items pass.

---

## Platform-Specific Adaptations

### Claude Code (Full Features)
- ✅ Subagents available for parallel testing
- ✅ Browser-based eval viewer
- ✅ Description optimization with `claude -p`
- ✅ All automation scripts work

### Claude.ai (Limited)
- ❌ No subagents → Run tests sequentially yourself
- ❌ No browser → Present results inline in conversation
- ❌ No description optimization → Skip that step
- ✅ Packaging works

### Cowork (Adapted)
- ✅ Subagents available (may timeout, use series if needed)
- ⚠️ No browser → Use `--static <output>` for HTML file
- ⚠️ Feedback downloads as file → Read from Downloads folder
- ✅ Description optimization works

---

## Getting Started

1. **Identify your scenario** (create vs optimize)
2. **Read the appropriate workflow**:
   - [Create Workflow](workflows/create-workflow.md)
   - [Optimize Workflow](workflows/optimize-workflow.md)
3. **Follow each step with checkpoints**
4. **Use the quality checklist before delivery**

---

## Additional Resources

### Guides
- [Design Principles](guides/principles.md) - Core concepts and best practices
- [Common Patterns](guides/patterns.md) - Reusable skill patterns
- [Anti-Patterns](guides/anti-patterns.md) - What to avoid

### References
- [JSON Schemas](references/schemas.md) - Data structures for evals, grading, benchmarks
- [Examples](references/examples.md) - Real-world skill examples

### Agents
- [Grader](agents/grader.md) - How to evaluate assertions
- [Comparator](agents/comparator.md) - Blind A/B comparison
- [Analyzer](agents/analyzer.md) - Benchmark analysis

---

## Important Notes

### Principle of Lack of Surprise
Skills must not contain malware, exploits, or malicious content. Don't create misleading skills or facilitate unauthorized access.

### Updating Existing Skills
When updating an installed skill:
- **Preserve the original name** (don't create v2 variants)
- **Copy to writable location** before editing (installed paths may be read-only)
- **Stage in /tmp/** if packaging manually

### Workflow Discipline
**Do not skip checkpoints.** Each checkpoint ensures quality and prevents costly mistakes later. If you find yourself wanting to skip ahead, stop and complete the current checkpoint first.

---

## Summary

This skill provides two systematic workflows:
1. **Create**: Build new skills from scratch with testing and validation
2. **Optimize**: Improve existing skills with baseline comparison

Both workflows include:
- Clear step-by-step instructions
- Mandatory checkpoints
- Automated testing tools
- Quality validation
- Packaging and delivery

**Next step:** Identify your scenario and jump to the appropriate workflow.
