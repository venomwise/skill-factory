# Anti-Patterns: What to Avoid

Common mistakes and anti-patterns to avoid when creating AI agent skills.

---

## 1. Verbosity Anti-Patterns

### Over-Explaining Basic Concepts

**Bad:**
```markdown
Python is a high-level, interpreted programming language created by Guido van Rossum
in 1991. It emphasizes code readability and uses significant indentation. Python
supports multiple programming paradigms including procedural, object-oriented, and
functional programming...
```

**Good:**
```markdown
Install dependencies:
```bash
pip install requests
```
```

### Redundant Instructions

**Bad:**
```markdown
First, you need to open the file. To open the file, use the open() function.
The open() function takes a filename as a parameter. After opening, you can read it.
To read it, use the read() method...
```

**Good:**
```markdown
Read the file:
```python
with open("file.txt") as f:
    content = f.read()
```
```

---

## 2. Structure Anti-Patterns

### Deeply Nested References

**Bad:**
```
SKILL.md → "See advanced.md"
  advanced.md → "See details.md"
    details.md → "See examples.md"
      examples.md → actual content
```

**Good:**
```
SKILL.md → advanced.md (complete info)
SKILL.md → examples.md (complete info)
```

### Monolithic Files

**Bad:**
```markdown
# My Skill (2000 lines)

## Section 1 (500 lines of details)
## Section 2 (500 lines of details)
## Section 3 (500 lines of details)
## Section 4 (500 lines of details)
```

**Good:**
```markdown
# My Skill (300 lines)

## Section 1
[Brief overview]
For details, see [section1.md](references/section1.md)

## Section 2
[Brief overview]
For details, see [section2.md](references/section2.md)
```

---

## 3. Instruction Anti-Patterns

### Too Many Options Without Guidance

**Bad:**
```markdown
You can use pypdf, pdfplumber, PyMuPDF, pdf2image, camelot, tabula-py,
pdfminer, or PyPDF2. Each has different features. PyPDF2 is older but
stable. pdfplumber is newer and has better text extraction. PyMuPDF is
faster but has a different API...
```

**Good:**
```markdown
Use pdfplumber for text extraction:
```python
import pdfplumber
```

For scanned PDFs requiring OCR, use pdf2image with pytesseract.
For tables, use camelot-py.
```

### Vague Instructions

**Bad:**
```markdown
Process the data appropriately.
Handle errors as needed.
Format the output nicely.
```

**Good:**
```markdown
Process the data:
1. Parse CSV with pandas
2. Filter rows where status == "active"
3. Sort by date descending
4. Export to JSON with indent=2
```

### Heavy-Handed Language

**Bad:**
```markdown
You MUST ALWAYS use this exact format. NEVER deviate from these steps.
This is ABSOLUTELY CRITICAL. FAILURE TO COMPLY will result in errors.
```

**Good:**
```markdown
Use this format to ensure consistency across reports:

[format here]

If your specific context requires adaptation, maintain these core principles:
- Clear section headers
- Data-driven findings
- Actionable recommendations
```

---

## 4. Description Anti-Patterns

### Vague Descriptions

**Bad:**
```yaml
description: Helps with documents
```

**Bad:**
```yaml
description: A useful tool for data processing
```

**Good:**
```yaml
description: Extract text and tables from PDF files, fill forms, merge documents. Use when working with PDF files or when the user mentions PDFs, forms, or document extraction.
```

### First/Second Person

**Bad:**
```yaml
description: I can help you process Excel files and create charts
```

**Bad:**
```yaml
description: You can use this to analyze spreadsheets
```

**Good:**
```yaml
description: Processes Excel files, creates charts, and generates reports. Use when the user mentions Excel, spreadsheets, .xlsx files, or data analysis.
```

### Missing Trigger Context

**Bad:**
```yaml
description: Generates commit messages
```

**Good:**
```yaml
description: Generates conventional commit messages following the format type(scope): subject. Use when the user wants to commit code, mentions git commit, or asks for help writing commit messages.
```

---

## 5. Example Anti-Patterns

### Abstract Examples

**Bad:**
```markdown
**Example:**
Input: Your data
Output: Processed result
```

**Good:**
```markdown
**Example:**
Input: sales_data.csv with columns: date, product, revenue, region
Output: monthly_summary.json with aggregated revenue by product and region
```

### No Examples

**Bad:**
```markdown
## Usage

Use the format described above to create your output.
```

**Good:**
```markdown
## Usage

**Example 1:**
Input: "Add user authentication"
Output:
```
feat(auth): implement user authentication

Add login endpoint and session management
```

**Example 2:**
Input: "Fix crash on startup"
Output:
```
fix(app): prevent crash on startup

Initialize config before loading modules
```
```

---

## 6. Dependency Anti-Patterns

### Assuming Dependencies Exist

**Bad:**
```markdown
Use the requests library to fetch data.
```

**Good:**
```markdown
Install requests:
```bash
pip install requests
```

Then fetch data:
```python
import requests
response = requests.get(url)
```
```

### Platform-Specific Paths

**Bad:**
```markdown
Run the script:
```bash
python scripts\helper.py
```
```

**Good:**
```markdown
Run the script:
```bash
python scripts/helper.py
```
```

---

## 7. Time-Sensitivity Anti-Patterns

### Hardcoded Dates

**Bad:**
```markdown
If you're doing this before August 2025, use the old API.
After August 2025, use the new API.
```

**Good:**
```markdown
## Current method
Use the v2 API: `api.example.com/v2/`

## Legacy method
<details>
<summary>v1 API (deprecated 2025-08)</summary>
Use `api.example.com/v1/` - no longer supported
</details>
```

### Version-Specific Instructions

**Bad:**
```markdown
As of version 3.2 (current as of March 2024), use this approach...
```

**Good:**
```markdown
For version 3.2+, use this approach:
[instructions]

For older versions, see [legacy.md](references/legacy.md)
```

---

## 8. Testing Anti-Patterns

### No Test Cases

**Bad:**
Creating a skill without any test cases or validation.

**Good:**
Create at least 3 realistic test cases before finalizing the skill.

### Overfitting to Test Cases

**Bad:**
```markdown
If the input is exactly "process sales data", do X.
If the input is exactly "generate report", do Y.
```

**Good:**
```markdown
When the user wants to process data:
1. Identify the data source
2. Determine the processing steps needed
3. Apply transformations
4. Generate output in requested format
```

### Non-Discriminating Assertions

**Bad:**
```json
{
  "assertion": "Output exists",
  "passed": true
}
```

**Good:**
```json
{
  "assertion": "Output contains all required sections: summary, findings, recommendations",
  "passed": true,
  "evidence": "Found sections: Executive Summary, Key Findings (3 items), Recommendations (5 items)"
}
```

---

## 9. Workflow Anti-Patterns

### Missing Error Handling

**Bad:**
```markdown
1. Read the file
2. Process the data
3. Write the output
```

**Good:**
```markdown
1. Read the file
   - If file not found, check common locations
   - If unreadable, try different encodings
2. Process the data
   - Validate format before processing
   - Handle missing or malformed data
3. Write the output
   - Create directory if needed
   - Verify write succeeded
```

### Rigid Workflows

**Bad:**
```markdown
ALWAYS follow these steps in EXACTLY this order:
1. Step A
2. Step B
3. Step C
NEVER skip steps or change the order.
```

**Good:**
```markdown
Typical workflow:
1. Step A - Prepare data
2. Step B - Process data
3. Step C - Generate output

Adapt as needed:
- If data is pre-processed, skip Step A
- For simple cases, combine Steps B and C
```

---

## 10. Terminology Anti-Patterns

### Inconsistent Terms

**Bad:**
```markdown
First, open the file. Then read the document. Next, parse the data from the file.
Finally, process the contents of the document.
```

**Good:**
```markdown
First, open the file. Then read the file. Next, parse the file.
Finally, process the file.
```

### Jargon Without Context

**Bad:**
```markdown
Use the factory pattern to instantiate the singleton with dependency injection.
```

**Good:**
```markdown
Create a single shared instance (singleton) that other parts of the code can use:

```python
class DataProcessor:
    _instance = None
    
    @classmethod
    def get_instance(cls):
        if cls._instance is None:
            cls._instance = cls()
        return cls._instance
```
```

---

## 11. Scope Anti-Patterns

### Scope Creep

**Bad:**
A skill called "pdf-processor" that also handles Word docs, Excel files,
images, videos, and audio files.

**Good:**
A skill called "pdf-processor" that focuses on PDF operations.
Create separate skills for other file types.

### Too Narrow

**Bad:**
A skill called "extract-text-from-single-page-pdfs-in-portrait-orientation"

**Good:**
A skill called "pdf-text-extraction" that handles various PDF types.

---

## 12. Script Anti-Patterns

### Punting on Errors

**Bad:**
```python
def process_file(path):
    return open(path).read()  # Just crash if file doesn't exist
```

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

### No Validation

**Bad:**
```python
def apply_changes(changes):
    for change in changes:
        execute(change)  # Hope for the best
```

**Good:**
```python
def apply_changes(changes):
    # Validate first
    for change in changes:
        if not validate_change(change):
            raise ValueError(f"Invalid change: {change}")
    
    # Then execute
    for change in changes:
        execute(change)
```

---

## 13. Documentation Anti-Patterns

### No Documentation

**Bad:**
A skill with only code and no explanation of what it does or how to use it.

**Good:**
Clear SKILL.md with overview, usage instructions, and examples.

### Documentation Separate from Code

**Bad:**
Documentation in a separate wiki or Google Doc that gets out of sync.

**Good:**
Documentation in the skill directory, versioned with the code.

---

## Summary

Avoid these anti-patterns:

1. ❌ Over-explaining basic concepts
2. ❌ Deeply nested references
3. ❌ Too many options without guidance
4. ❌ Vague descriptions or instructions
5. ❌ Abstract examples
6. ❌ Assuming dependencies exist
7. ❌ Time-sensitive information
8. ❌ No test cases
9. ❌ Missing error handling
10. ❌ Inconsistent terminology
11. ❌ Scope creep or too narrow scope
12. ❌ Poor error handling in scripts
13. ❌ Missing or outdated documentation

Instead:

1. ✅ Be concise and assume intelligence
2. ✅ Keep references one level deep
3. ✅ Provide clear default recommendations
4. ✅ Give specific, actionable instructions
5. ✅ Use concrete examples
6. ✅ Explicitly state dependencies
7. ✅ Use version-agnostic patterns
8. ✅ Create realistic test cases
9. ✅ Handle errors gracefully
10. ✅ Use consistent terminology
11. ✅ Define appropriate scope
12. ✅ Validate before executing
13. ✅ Keep documentation current and co-located
