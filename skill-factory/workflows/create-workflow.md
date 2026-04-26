# Create Workflow: Building a New Skill

This workflow guides you through creating a new skill from scratch with systematic testing and validation.

**Estimated time:** 30-60 minutes for a simple skill, 2-4 hours for complex skills

---

## Step 1: Capture Intent

**Goal:** Understand what the user wants the skill to do.

### Actions

1. **Check conversation history** - The user may have already demonstrated a workflow. Look for:
   - Tools they used
   - Sequence of steps
   - Corrections they made
   - Input/output formats

2. **Ask clarifying questions:**
   - What should this skill enable Claude to do?
   - When should this skill trigger? (what phrases/contexts)
   - What's the expected output format?
   - Are there specific edge cases to handle?

3. **Determine if testing is needed:**
   - **Objective outputs** (file transforms, data extraction, code generation) → Need test cases
   - **Subjective outputs** (writing style, creative work) → May skip formal testing
   - Suggest the appropriate default, but let user decide

### Checkpoint 1: Intent Confirmed

**Before proceeding, verify:**
- [ ] You can clearly explain what the skill does in 1-2 sentences
- [ ] You know when it should trigger (specific user phrases or contexts)
- [ ] You understand the expected output format
- [ ] User has confirmed the intent is correct
- [ ] Decision made on whether to create test cases

**If any item is unclear, go back and ask more questions.**

---

## Step 2: Interview & Research

**Goal:** Gather all necessary context to write a complete skill.

### Actions

1. **Proactive questioning:**
   - What are common edge cases?
   - Are there example files I should see?
   - What does success look like?
   - Are there dependencies or prerequisites?
   - What should happen when things go wrong?

2. **Research (if needed):**
   - Check available MCPs for relevant documentation
   - Look for similar existing skills
   - Research best practices for this domain
   - Use subagents for parallel research if available

3. **Document findings:**
   - Note key requirements
   - List dependencies
   - Identify potential challenges
   - Collect example inputs/outputs

### Checkpoint 2: Requirements Clear

**Before proceeding, verify:**
- [ ] All edge cases identified
- [ ] Input/output formats documented
- [ ] Dependencies and prerequisites listed
- [ ] Success criteria defined
- [ ] User has reviewed and confirmed requirements
- [ ] You have enough context to write the skill

**If research is incomplete, continue gathering information.**

---

## Step 3: Write SKILL.md

**Goal:** Create the skill definition following best practices.

### Actions

1. **Create frontmatter:**
   ```yaml
   ---
   name: skill-name
   description: >
     What it does AND when to use it. Include specific contexts and keywords
     that should trigger this skill. Be slightly "pushy" to combat undertriggering.
   ---
   ```

2. **Write main content sections:**
   - **Overview**: Brief introduction
   - **When to use**: Specific scenarios (if not fully covered in description)
   - **Instructions**: Step-by-step guidance
   - **Examples**: Concrete input/output pairs
   - **Edge cases**: How to handle special situations
   - **References**: Links to additional resources (if needed)

3. **Apply best practices:**
   - Keep main file < 500 lines
   - Use imperative form for instructions
   - Explain the "why" behind requirements
   - Provide concrete examples
   - Move heavy content to reference files
   - Use progressive disclosure

4. **Create supporting files (if needed):**
   - `scripts/` - Reusable helper scripts
   - `references/` - Detailed documentation
   - `assets/` - Templates or resources

### Checkpoint 3: Draft Reviewed

**Before proceeding, verify:**
- [ ] Frontmatter has `name` and `description`
- [ ] Description includes what it does AND when to use it
- [ ] Main file is under 500 lines (or has clear reference structure)
- [ ] Instructions are clear and actionable
- [ ] Examples are concrete and realistic
- [ ] User has reviewed and approved the draft
- [ ] No obvious gaps or ambiguities

**If draft needs work, revise before continuing.**

---

## Step 4: Create Test Cases

**Goal:** Define realistic test scenarios to validate the skill.

**Note:** Skip this step if user decided testing isn't needed (subjective outputs).

### Actions

1. **Generate 2-3 realistic test prompts:**
   - Use language a real user would actually say
   - Include specific details (file paths, names, values)
   - Cover different aspects of the skill
   - Mix simple and complex cases

2. **Share with user:**
   "Here are test cases I'd like to try. Do these look right, or should we add more?"

3. **Create `evals/evals.json`:**
   ```json
   {
     "skill_name": "example-skill",
     "evals": [
       {
         "id": 1,
         "prompt": "User's task prompt with specific details",
         "expected_output": "Description of expected result",
         "files": []
       }
     ]
   }
   ```

4. **Don't write assertions yet** - You'll draft them while tests run.

### Checkpoint 4: Tests Approved

**Before proceeding, verify:**
- [ ] 2-3 realistic test cases created
- [ ] Test prompts use natural language
- [ ] Tests cover key functionality
- [ ] User has approved the test cases
- [ ] `evals/evals.json` file created
- [ ] Workspace directory structure planned: `<skill-name>-workspace/iteration-1/`

**If tests are inadequate, create better ones.**

---

## Step 5: Run Evaluations

**Goal:** Execute all test cases with and without the skill.

**IMPORTANT:** This is one continuous sequence. Do NOT stop partway through.

### Actions

1. **Create workspace structure:**
   ```
   <skill-name>-workspace/
   └── iteration-1/
       ├── eval-0/
       │   ├── with_skill/outputs/
       │   └── without_skill/outputs/
       ├── eval-1/
       │   ├── with_skill/outputs/
       │   └── without_skill/outputs/
       └── ...
   ```

2. **Spawn all runs in the SAME turn:**
   
   For each test case, spawn TWO subagents simultaneously:
   
   **With-skill run:**
   ```
   Execute this task:
   - Skill path: <path-to-skill>
   - Task: <eval prompt>
   - Input files: <eval files if any, or "none">
   - Save outputs to: <workspace>/iteration-1/eval-<ID>/with_skill/outputs/
   - Outputs to save: <what user cares about>
   ```
   
   **Baseline run (without skill):**
   ```
   Execute this task:
   - NO skill provided
   - Task: <same eval prompt>
   - Input files: <same files>
   - Save outputs to: <workspace>/iteration-1/eval-<ID>/without_skill/outputs/
   - Outputs to save: <same outputs>
   ```

3. **Create `eval_metadata.json` for each test:**
   ```json
   {
     "eval_id": 0,
     "eval_name": "descriptive-name-here",
     "prompt": "The user's task prompt",
     "assertions": []
   }
   ```

4. **While runs are in progress, draft assertions:**
   - Don't wait idle - use this time productively
   - Create objectively verifiable checks
   - Use descriptive names
   - Explain to user what each assertion checks
   - Update `eval_metadata.json` and `evals/evals.json`

5. **As runs complete, capture timing data:**
   
   When each subagent completes, save `timing.json` immediately:
   ```json
   {
     "total_tokens": 84852,
     "duration_ms": 23332,
     "total_duration_seconds": 23.3
   }
   ```
   
   **This is the only chance to capture this data!**

### Checkpoint 5: All Runs Complete

**Before proceeding, verify:**
- [ ] All with_skill runs completed successfully
- [ ] All without_skill runs completed successfully
- [ ] Outputs saved to correct directories
- [ ] Timing data captured for all runs
- [ ] Assertions drafted and documented
- [ ] `eval_metadata.json` created for each test
- [ ] No runs failed or timed out

**If any runs failed, investigate and rerun.**

---

## Step 6: Review Results

**Goal:** Present results to user and collect feedback.

### Actions

1. **Grade each run:**
   - Spawn grader subagent (or grade inline)
   - Read `agents/grader.md` for instructions
   - Evaluate each assertion against outputs
   - Save to `grading.json` in each run directory
   - Use exact field names: `text`, `passed`, `evidence`

2. **Aggregate benchmark:**
   ```bash
   python -m scripts.aggregate_benchmark <workspace>/iteration-1 --skill-name <name>
   ```
   This produces `benchmark.json` and `benchmark.md`

3. **Analyst pass:**
   - Read benchmark data
   - Look for patterns (see `agents/analyzer.md`)
   - Note non-discriminating assertions
   - Identify high-variance evals
   - Check time/token tradeoffs

4. **Launch eval viewer:**
   ```bash
   nohup python <skill-factory-path>/eval-viewer/generate_review.py \
     <workspace>/iteration-1 \
     --skill-name "my-skill" \
     --benchmark <workspace>/iteration-1/benchmark.json \
     > /dev/null 2>&1 &
   VIEWER_PID=$!
   ```
   
   **For Cowork/headless:** Use `--static <output_path>` instead

5. **Tell user:**
   "I've opened the results in your browser. Two tabs:
   - 'Outputs': Click through each test case, leave feedback
   - 'Benchmark': Quantitative comparison
   
   When done, come back and let me know."

6. **Wait for user to finish reviewing**

7. **Read feedback:**
   ```bash
   cat <workspace>/iteration-1/feedback.json
   ```
   
   Empty feedback = user thought it was fine
   Focus on test cases with specific complaints

8. **Kill viewer:**
   ```bash
   kill $VIEWER_PID 2>/dev/null
   ```

### Checkpoint 6: Feedback Collected

**Before proceeding, verify:**
- [ ] All runs graded with assertions
- [ ] Benchmark aggregated successfully
- [ ] Analyst observations documented
- [ ] Eval viewer launched and user reviewed results
- [ ] Feedback collected from user
- [ ] Viewer process killed
- [ ] You understand what needs improvement

**If feedback is unclear, ask user for clarification.**

---

## Step 7: Iterate (If Needed)

**Goal:** Improve the skill based on feedback.

**Decision point:** 
- If feedback is all positive/empty → Skip to Step 8
- If there are issues → Continue with iteration

### Actions

1. **Analyze feedback:**
   - Generalize from specific complaints
   - Look for patterns across test cases
   - Read transcripts, not just outputs
   - Identify root causes

2. **Plan improvements:**
   - Keep changes lean - remove what doesn't work
   - Explain the "why" behind requirements
   - Look for repeated work to bundle into scripts
   - Avoid overfitting to test cases

3. **Apply improvements:**
   - Edit SKILL.md
   - Update or add reference files
   - Create helper scripts if needed
   - Keep terminology consistent

4. **Rerun evaluations:**
   - Create `iteration-2/` directory
   - Spawn all test cases again (with_skill + without_skill)
   - Follow same process as Step 5
   - Launch viewer with `--previous-workspace <workspace>/iteration-1`

5. **Collect new feedback**

6. **Repeat until:**
   - User says they're happy
   - Feedback is all empty
   - Not making meaningful progress

### Checkpoint 7: Improvements Done

**Before proceeding, verify:**
- [ ] All feedback addressed
- [ ] Improvements applied to skill
- [ ] New iteration tested
- [ ] Results compared to previous iteration
- [ ] User satisfied with improvements
- [ ] No regressions introduced

**If issues remain, continue iterating.**

---

## Step 8: Final Validation

**Goal:** Ensure skill meets all quality standards.

### Actions

1. **Run through Essential Quality checklist:**
   - [ ] Description complete (what + when)
   - [ ] Main file < 500 lines
   - [ ] At least 3 test cases
   - [ ] All tests pass
   - [ ] Terminology consistent

2. **Run through Recommended Quality checklist:**
   - [ ] No time-sensitive info
   - [ ] Examples concrete
   - [ ] References organized
   - [ ] Progressive disclosure used
   - [ ] Workflows clear

3. **Final review with user:**
   - Show final skill content
   - Demonstrate test results
   - Confirm it meets their needs

4. **Optional: Description optimization**
   
   Ask user: "Would you like me to optimize the skill description for better triggering accuracy?"
   
   If yes, follow description optimization process (see main SKILL.md)

### Checkpoint 8: Quality Validated

**Before proceeding, verify:**
- [ ] All Essential Quality items pass
- [ ] Most Recommended Quality items pass
- [ ] User has done final review
- [ ] User confirms skill is ready
- [ ] Description optimization done (if requested)

**Do not proceed to packaging until all Essential items pass.**

---

## Step 9: Package & Deliver

**Goal:** Package the skill for installation and delivery.

### Actions

1. **Package the skill:**
   ```bash
   python -m scripts.package_skill <path/to/skill-folder>
   ```
   
   This creates `<skill-name>.skill` file

2. **Present to user:**
   - Show path to `.skill` file
   - Provide installation instructions
   - Explain how to use the skill

3. **Document the skill:**
   - Ensure README or documentation exists
   - Include usage examples
   - Note any dependencies

### Checkpoint 9: Skill Packaged

**Final verification:**
- [ ] Skill packaged successfully
- [ ] `.skill` file created
- [ ] User knows where to find it
- [ ] Installation instructions provided
- [ ] Documentation complete

**Congratulations! Skill creation complete.**

---

## Troubleshooting

### Common Issues

**Runs timing out:**
- Run tests sequentially instead of parallel
- Simplify test cases
- Check for infinite loops in skill instructions

**Assertions always pass/fail:**
- Assertions may be non-discriminating
- Revise to be more specific
- Consider removing unhelpful assertions

**User feedback is vague:**
- Ask specific questions
- Show side-by-side comparisons
- Request concrete examples of what's wrong

**Skill file too large:**
- Move content to reference files
- Remove redundant explanations
- Use progressive disclosure

---

## Platform-Specific Notes

### Claude.ai
- No subagents → Run tests yourself sequentially
- No browser → Present results inline
- Skip baseline comparisons
- Skip description optimization

### Cowork
- Use `--static` for eval viewer
- Feedback downloads as file
- May need to run tests in series if timeouts occur

---

## Summary

This workflow ensures systematic skill creation with quality validation:

1. ✅ Capture Intent → Understand requirements
2. ✅ Interview & Research → Gather context
3. ✅ Write SKILL.md → Create skill definition
4. ✅ Create Test Cases → Define validation
5. ✅ Run Evaluations → Test with/without skill
6. ✅ Review Results → Collect user feedback
7. ✅ Iterate → Improve based on feedback
8. ✅ Final Validation → Quality checklist
9. ✅ Package & Deliver → Ship the skill

**Each checkpoint prevents costly mistakes and ensures quality.**
