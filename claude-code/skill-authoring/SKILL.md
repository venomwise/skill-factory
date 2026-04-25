---
name: skill-authoring
description: >
  Guide for creating and optimizing AI agent skills. Use when the user wants to create a new skill,
  improve an existing skill, review skill quality, or understand skill authoring best practices.
  Covers skill structure, content guidelines, common patterns, and anti-patterns.
---

# Skill Authoring Best Practices

This skill provides universal guidelines for creating effective AI agent skills that work across different platforms and models.

## Core Principles

### 1. Conciseness is Key

The context window is a shared resource. Your skill competes with:
- System prompts
- Conversation history
- Other skills' metadata
- The actual user request

**Default assumption:** The AI agent is already intelligent.

Only add context the agent doesn't already have. Challenge each piece of information:
- "Does the agent really need this explanation?"
- "Can I assume the agent knows this?"
- "Does this paragraph justify its token cost?"

**Good example (concise):**
```markdown
## Extract PDF text

Use pdfplumber for text extraction:

```python
import pdfplumber
with pdfplumber.open("file.pdf") as pdf:
    text = pdf.pages[0].extract_text()
```
```

**Bad example (verbose):**
```markdown
## Extract PDF text

PDF (Portable Document Format) files are a common file format that contains
text, images, and other content. To extract text from a PDF, you'll need to
use a library. There are many libraries available...
```

### 2. Set Appropriate Degrees of Freedom

Match the level of specificity to the task's fragility and variability.

**High freedom** (text-based instructions):
- Multiple approaches are valid
- Decisions depend on context
- Heuristics guide the approach

Example:
```markdown
## Code review process

1. Analyze the code structure and organization
2. Check for potential bugs or edge cases
3. Suggest improvements for readability
4. Verify adherence to project conventions
```

**Medium freedom** (pseudocode or scripts with parameters):
- A preferred pattern exists
- Some variation is acceptable
- Configuration affects behavior

**Low freedom** (specific scripts, few parameters):
- Operations are fragile and error-prone
- Consistency is critical
- A specific sequence must be followed

**Analogy:** Think of the agent as exploring a path:
- **Narrow bridge with cliffs:** Provide specific guardrails (low freedom)
- **Open field:** Give general direction and trust the agent (high freedom)

### 3. Progressive Disclosure

Structure content in layers that load on-demand:

**Level 1: Metadata** (always loaded)
- Skill name and description
- Enables discovery without context penalty

**Level 2: Main instructions** (loaded when triggered)
- Core workflows and guidance
- Keep under 500 lines

**Level 3: Detailed resources** (loaded as needed)
- Reference documentation
- Detailed examples
- Utility scripts

### 4. Test-Driven Development

**Create evaluations BEFORE writing extensive documentation.**

1. Identify gaps: Run the agent on tasks without the skill
2. Create evaluations: Build 3+ test scenarios
3. Establish baseline: Measure performance without the skill
4. Write minimal instructions: Address only the identified gaps
5. Iterate: Execute evaluations and refine

## Skill Structure

### Required Components

Every skill needs a main file (e.g., `SKILL.md`) with metadata:

```yaml
---
name: your-skill-name
description: Brief description of what this skill does and when to use it
---

# Your Skill Name

## Instructions
[Clear, step-by-step guidance]

## Examples
[Concrete examples]
```

**Required fields:**
- `name`: Lowercase letters, numbers, hyphens only (max 64 chars)
- `description`: What the skill does AND when to use it (max 1024 chars)

### Naming Conventions

Use consistent naming patterns. Consider **gerund form** (verb + -ing):

**Good examples:**
- `processing-pdfs`
- `analyzing-spreadsheets`
- `managing-databases`
- `testing-code`

**Avoid:**
- Vague names: `helper`, `utils`, `tools`
- Overly generic: `documents`, `data`, `files`

### Writing Effective Descriptions

**Always write in third person** (description goes into system prompt):

- ✓ Good: "Processes Excel files and generates reports"
- ✗ Avoid: "I can help you process Excel files"
- ✗ Avoid: "You can use this to process Excel files"

**Be specific and include key terms:**

```yaml
description: Extract text and tables from PDF files, fill forms, merge documents. Use when working with PDF files or when the user mentions PDFs, forms, or document extraction.
```

Avoid vague descriptions:
```yaml
description: Helps with documents  # Too vague
```

## Progressive Disclosure Patterns

### Pattern 1: High-level guide with references

```markdown
# PDF Processing

## Quick start
[Basic usage here]

## Advanced features
**Form filling**: See [FORMS.md](FORMS.md)
**API reference**: See [REFERENCE.md](REFERENCE.md)
**Examples**: See [EXAMPLES.md](EXAMPLES.md)
```

The agent loads additional files only when needed.

### Pattern 2: Domain-specific organization

For skills with multiple domains, organize by domain:

```
bigquery-skill/
├── SKILL.md (overview)
└── reference/
    ├── finance.md
    ├── sales.md
    └── product.md
```

```markdown
# BigQuery Analysis

## Available datasets
**Finance**: Revenue, ARR → See [reference/finance.md](reference/finance.md)
**Sales**: Pipeline, accounts → See [reference/sales.md](reference/sales.md)
**Product**: API usage → See [reference/product.md](reference/product.md)
```

### Pattern 3: Conditional details

```markdown
# Document Processing

## Creating documents
Use docx-js for new documents. See [DOCX-JS.md](DOCX-JS.md).

## Editing documents
For simple edits, modify directly.

**For tracked changes**: See [REDLINING.md](REDLINING.md)
**For advanced features**: See [ADVANCED.md](ADVANCED.md)
```

### Important: Avoid deeply nested references

Keep references **one level deep** from the main file. The agent may only partially read deeply nested files.

**Bad (too deep):**
```markdown
# SKILL.md → advanced.md → details.md → actual info
```

**Good (one level):**
```markdown
# SKILL.md → advanced.md (complete info)
# SKILL.md → reference.md (complete info)
```

### Structure longer reference files

For files over 100 lines, include a table of contents:

```markdown
# API Reference

## Contents
- Authentication and setup
- Core methods
- Advanced features
- Error handling
- Code examples

## Authentication and setup
...
```

## Content Guidelines

### Avoid Time-Sensitive Information

Don't include information that will become outdated:

**Bad (time-sensitive):**
```markdown
If you're doing this before August 2025, use the old API.
```

**Good (use "old patterns" section):**
```markdown
## Current method
Use the v2 API endpoint: `api.example.com/v2/messages`

## Old patterns
<details>
<summary>Legacy v1 API (deprecated 2025-08)</summary>
The v1 API used: `api.example.com/v1/messages`
This endpoint is no longer supported.
</details>
```

### Use Consistent Terminology

Choose one term and use it throughout:

**Good (consistent):**
- Always "API endpoint"
- Always "field"
- Always "extract"

**Bad (inconsistent):**
- Mix "API endpoint", "URL", "API route", "path"
- Mix "field", "box", "element", "control"

## Common Patterns

### Template Pattern

Provide templates for output format.

**For strict requirements:**
```markdown
## Report structure

ALWAYS use this exact template:

```markdown
# [Analysis Title]

## Executive summary
[One-paragraph overview]

## Key findings
- Finding 1 with data
- Finding 2 with data

## Recommendations
1. Specific recommendation
2. Specific recommendation
```
```

**For flexible guidance:**
```markdown
## Report structure

Here is a sensible default, but adapt as needed:

[template here]

Adjust sections based on the specific analysis type.
```

### Examples Pattern

Provide input/output pairs:

```markdown
## Commit message format

**Example 1:**
Input: Added user authentication with JWT tokens
Output:
```
feat(auth): implement JWT-based authentication

Add login endpoint and token validation middleware
```

**Example 2:**
Input: Fixed bug where dates displayed incorrectly
Output:
```
fix(reports): correct date formatting

Use UTC timestamps consistently
```
```

### Conditional Workflow Pattern

Guide the agent through decision points:

```markdown
## Document modification workflow

1. Determine the modification type:
   **Creating new content?** → Follow "Creation workflow"
   **Editing existing?** → Follow "Editing workflow"

2. Creation workflow:
   - Use docx-js library
   - Build from scratch
   - Export to .docx

3. Editing workflow:
   - Unpack existing document
   - Modify XML directly
   - Validate changes
   - Repack when complete
```

### Workflow Checklist Pattern

For complex multi-step tasks, provide a checklist:

```markdown
## Research synthesis workflow

Copy this checklist and track progress:

```
Research Progress:
- [ ] Step 1: Read all source documents
- [ ] Step 2: Identify key themes
- [ ] Step 3: Cross-reference claims
- [ ] Step 4: Create structured summary
- [ ] Step 5: Verify citations
```

**Step 1: Read all source documents**
Review each document in `sources/`. Note main arguments.

**Step 2: Identify key themes**
Look for patterns. Where do sources agree or disagree?

[Continue with detailed steps...]
```

### Feedback Loop Pattern

Implement validation cycles:

```markdown
## Content review process

1. Draft content following guidelines
2. Review against checklist:
   - Check terminology consistency
   - Verify examples follow format
   - Confirm all sections present
3. If issues found:
   - Note each issue
   - Revise content
   - Review checklist again
4. Only proceed when all requirements met
5. Finalize document
```

## Anti-Patterns to Avoid

### 1. Don't Be Too Verbose

Assume the agent is intelligent. Don't over-explain.

### 2. Don't Offer Too Many Options

Provide a default recommendation:

**Bad:**
```markdown
You can use pypdf, or pdfplumber, or PyMuPDF, or pdf2image...
```

**Good:**
```markdown
Use pdfplumber for text extraction:
```python
import pdfplumber
```

For scanned PDFs requiring OCR, use pdf2image with pytesseract instead.
```

### 3. Avoid Platform-Specific Paths

Always use forward slashes:

- ✓ Good: `scripts/helper.py`, `reference/guide.md`
- ✗ Avoid: `scripts\helper.py`, `reference\guide.md`

Unix-style paths work across all platforms.

### 4. Don't Assume Dependencies Are Installed

Be explicit about requirements:

**Bad:**
```markdown
Use the pdf library to process the file.
```

**Good:**
```markdown
Install required package: `pip install pypdf`

Then use it:
```python
from pypdf import PdfReader
reader = PdfReader("file.pdf")
```
```

## Iterative Development

### Work with the Agent

The most effective development process involves the agent itself:

1. **Complete a task without a skill**: Work through a problem using normal prompting
2. **Identify the reusable pattern**: Notice what context you repeatedly provide
3. **Ask the agent to create a skill**: "Create a skill that captures this pattern"
4. **Review for conciseness**: Remove unnecessary explanations
5. **Test on similar tasks**: Use the skill with fresh agent instances
6. **Iterate based on observation**: Refine based on actual usage

### Observe Agent Behavior

Pay attention to how the agent uses your skill:

- **Unexpected exploration paths**: Structure may not be intuitive
- **Missed connections**: Links might need to be more explicit
- **Overreliance on certain sections**: Consider moving content to main file
- **Ignored content**: May be unnecessary or poorly signaled

### Gather Feedback

1. Share skills with teammates
2. Observe their usage
3. Ask: Does it activate when expected? Are instructions clear?
4. Incorporate feedback

## Quality Checklist

Before sharing a skill, verify:

### Core Quality
- [ ] Description is specific and includes key terms
- [ ] Description includes what it does AND when to use it
- [ ] Main file is under 500 lines
- [ ] Additional details in separate files (if needed)
- [ ] No time-sensitive information
- [ ] Consistent terminology throughout
- [ ] Examples are concrete
- [ ] File references are one level deep
- [ ] Progressive disclosure used appropriately
- [ ] Workflows have clear steps

### Testing
- [ ] At least 3 evaluations created
- [ ] Tested with real usage scenarios
- [ ] Team feedback incorporated (if applicable)

## Advanced: Skills with Executable Code

For skills that include scripts:

### Solve, Don't Punt

Handle error conditions explicitly:

**Good:**
```python
def process_file(path):
    try:
        with open(path) as f:
            return f.read()
    except FileNotFoundError:
        print(f"File {path} not found, creating default")
        with open(path, "w") as f:
            f.write("")
        return ""
```

**Bad:**
```python
def process_file(path):
    return open(path).read()  # Just fail
```

### Provide Utility Scripts

Pre-made scripts offer advantages:
- More reliable than generated code
- Save tokens
- Save time
- Ensure consistency

Example:
```markdown
## Utility scripts

**analyze_form.py**: Extract form fields from PDF

```bash
python scripts/analyze_form.py input.pdf > fields.json
```

**validate_boxes.py**: Check for overlapping boxes

```bash
python scripts/validate_boxes.py fields.json
```
```

### Create Verifiable Intermediate Outputs

Use "plan-validate-execute" pattern:

1. Agent creates a plan file (e.g., `changes.json`)
2. Validation script checks the plan
3. If valid, execute the plan
4. Verify the output

This catches errors before applying changes.

## Additional Resources

For more detailed guidance, see:
- [Common Patterns](references/patterns.md) - Detailed pattern examples
- [Anti-Patterns](references/anti-patterns.md) - What to avoid
- [Examples](references/examples.md) - Real-world skill examples

## Summary

Effective skills are:
- **Concise**: Only include what the agent doesn't know
- **Well-structured**: Use progressive disclosure
- **Tested**: Built with evaluations
- **Iterated**: Improved based on real usage
- **Consistent**: Use clear terminology and patterns

Focus on making the agent successful, not on comprehensive documentation.
