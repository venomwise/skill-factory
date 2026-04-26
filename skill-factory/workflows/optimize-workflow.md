# Optimize Workflow: Improving an Existing Skill

This workflow guides you through systematically improving an existing skill with baseline comparison and validation.

**Estimated time:** 1-3 hours depending on complexity

---

## Step 1: Analyze Current State

**Goal:** Understand the existing skill and identify what needs improvement.

### Actions

1. **Locate and read the skill:**
   - Find the skill directory
   - Read SKILL.md thoroughly
   - Check for reference files
   - Review any existing scripts

2. **Understand current behavior:**
   - What does it do now?
   - When does it trigger?
   - What are known issues?

3. **Identify problems:**
   - User-reported issues
   - Triggering problems (over/under-triggering)
   - Quality issues in outputs
   - Performance problems
   - Missing functionality

4. **Check for existing tests:**
   - Look for `evals/evals.json`
   - Review test coverage
   - Note if tests are outdated

### Checkpoint 1: Issues Identified

**Before proceeding, verify:**
- [ ] Skill location confirmed
- [ ] Current functionality understood
- [ ] Specific problems identified and documented
- [ ] Existing tests reviewed (if any)
- [ ] User has confirmed these are the right issues to fix

**If issues are unclear, ask user for more details.**

---

## Step 2: Define Optimization Goal

**Goal:** Set clear, measurable objectives for the optimization.

### Actions

1. **Categorize the optimization type:**
   
   **Quality Improvement:**
   - Outputs don't meet requirements
   - Instructions are unclear
   - Missing edge case handling
   - Examples are inadequate
   
   **Triggering Optimization:**
   - Skill doesn't trigger when it should
   - Skill triggers when it shouldn't
   - Description needs refinement
   
   **Performance Optimization:**
   - Takes too long to execute
   - Uses too many tokens
   - Inefficient workflow
   
   **Feature Addition:**
   - Add new capabilities
   - Support new use cases
   - Extend existing functionality

2. **Define success criteria:**
   - What does "better" look like?
   - How will you measure improvement?
   - What metrics matter? (pass rate, speed, quality)

3. **Set scope:**
   - What will change?
   - What must stay the same?
   - Are there constraints?

### Checkpoint 2: Goal Confirmed

**Before proceeding, verify:**
- [ ] Optimization type identified
- [ ] Success criteria defined
- [ ] Metrics for comparison chosen
- [ ] Scope is clear
- [ ] User agrees with the goal

**If goal is vague, refine it with user.**

---

## Step 3: Create/Update Tests

**Goal:** Ensure adequate test coverage for the optimization.

### Actions

1. **Review existing tests:**
   - Are current tests adequate?
   - Do they cover the problem areas?
   - Are they still relevant?

2. **Add new tests if needed:**
   - Create tests that expose the current problems
   - Cover new functionality being added
   - Include edge cases
   - Use realistic user prompts

3. **Update `evals/evals.json`:**
   ```json
   {
     "skill_name": "example-skill",
     "evals": [
       {
         "id": 1,
         "prompt": "Realistic user prompt",
         "expected_output": "What should happen",
         "files": []
       }
     ]
   }
   ```

4. **Aim for 3-5 test cases minimum:**
   - At least one test that currently fails
   - Tests that should continue to pass
   - Edge cases

### Checkpoint 3: Tests Ready

**Before proceeding, verify:**
- [ ] Test coverage is adequate
- [ ] Tests expose current problems
- [ ] New functionality is covered
- [ ] `evals/evals.json` updated
- [ ] User has approved test cases

**If test coverage is inadequate, add more tests.**

---

## Step 4: Snapshot Baseline

**Goal:** Capture current performance for comparison.

### Actions

1. **Create workspace:**
   ```
   <skill-name>-workspace/
   └── iteration-1/
       ├── skill-snapshot/  (copy of original skill)
       └── eval-0/, eval-1/, ...
   ```

2. **Snapshot the current skill:**
   ```bash
   cp -r <skill-path> <workspace>/skill-snapshot/
   ```

3. **Run baseline tests:**
   
   For each test case, spawn subagent:
   ```
   Execute this task:
   - Skill path: <workspace>/skill-snapshot/
   - Task: <eval prompt>
   - Input files: <eval files if any>
   - Save outputs to: <workspace>/iteration-1/eval-<ID>/old_skill/outputs/
   - Outputs to save: <what user cares about>
   ```

4. **Capture timing data** as runs complete

5. **Create `eval_metadata.json`** for each test

### Checkpoint 4: Baseline Captured

**Before proceeding, verify:**
- [ ] Original skill snapshotted
- [ ] All baseline tests completed
- [ ] Outputs saved correctly
- [ ] Timing data captured
- [ ] `eval_metadata.json` created for each test
- [ ] No runs failed

**If baseline runs failed, investigate and rerun.**

---

## Step 5: Apply Improvements

**Goal:** Make targeted improvements to the skill.

### Actions

1. **Analyze the problems:**
   - Review baseline outputs
   - Read transcripts if available
   - Identify root causes
   - Look for patterns

2. **Plan improvements:**
   
   **For quality issues:**
   - Clarify ambiguous instructions
   - Add missing examples
   - Explain the "why" behind requirements
   - Remove unhelpful content
   
   **For triggering issues:**
   - Enhance description with keywords
   - Add specific use cases
   - Make description more "pushy"
   - Include context clues
   
   **For performance issues:**
   - Bundle repeated work into scripts
   - Remove unnecessary steps
   - Optimize workflow
   - Use progressive disclosure
   
   **For feature additions:**
   - Add new sections to SKILL.md
   - Create reference files if needed
   - Add helper scripts
   - Update examples

3. **Apply changes:**
   - Edit SKILL.md in original location
   - Update reference files
   - Add/modify scripts
   - Keep changes focused

4. **Document changes:**
   - Note what was changed and why
   - Keep track for comparison later

### Checkpoint 5: Changes Reviewed

**Before proceeding, verify:**
- [ ] Root causes identified
- [ ] Improvements planned
- [ ] Changes applied to skill
- [ ] Changes are focused and targeted
- [ ] No unintended modifications
- [ ] User has reviewed the changes

**If changes are too broad or risky, refine them.**

---

## Step 6: Run Evaluations

**Goal:** Test the improved skill against baseline.

**IMPORTANT:** This is one continuous sequence. Do NOT stop partway through.

### Actions

1. **Spawn all test runs in the SAME turn:**
   
   For each test case, spawn subagent with improved skill:
   ```
   Execute this task:
   - Skill path: <improved-skill-path>
   - Task: <eval prompt>
   - Input files: <eval files if any>
   - Save outputs to: <workspace>/iteration-1/eval-<ID>/with_skill/outputs/
   - Outputs to save: <what user cares about>
   ```

2. **While runs are in progress:**
   - Draft or update assertions
   - Review baseline outputs
   - Prepare comparison criteria

3. **Capture timing data** as runs complete

4. **Grade all runs:**
   - Grade improved skill runs
   - Grade baseline runs (if not done yet)
   - Use same assertions for both
   - Save to `grading.json`

5. **Aggregate benchmark:**
   ```bash
   python -m scripts.aggregate_benchmark <workspace>/iteration-1 --skill-name <name>
   ```

### Checkpoint 6: All Runs Complete

**Before proceeding, verify:**
- [ ] All improved skill runs completed
- [ ] All baseline runs completed (from Step 4)
- [ ] Timing data captured for all runs
- [ ] All runs graded with assertions
- [ ] Benchmark aggregated
- [ ] No runs failed

**If any runs failed, investigate and rerun.**

---

## Step 7: Compare Results

**Goal:** Analyze improvements and present comparison to user.

### Actions

1. **Analyst pass:**
   - Read `benchmark.json`
   - Compare pass rates (improved vs baseline)
   - Compare timing and token usage
   - Look for regressions
   - Identify unexpected changes
   - See `agents/analyzer.md` for guidance

2. **Launch eval viewer:**
   ```bash
   nohup python <skill-factory-path>/eval-viewer/generate_review.py \
     <workspace>/iteration-1 \
     --skill-name "my-skill" \
     --benchmark <workspace>/iteration-1/benchmark.json \
     > /dev/null 2>&1 &
   VIEWER_PID=$!
   ```
   
   **For Cowork/headless:** Use `--static <output_path>`

3. **Tell user:**
   "I've opened the comparison results. You can see:
   - Side-by-side outputs (improved vs baseline)
   - Quantitative metrics (pass rate, time, tokens)
   - Per-test breakdown
   
   Review and let me know what you think."

4. **Wait for user feedback**

5. **Read feedback:**
   ```bash
   cat <workspace>/iteration-1/feedback.json
   ```

6. **Kill viewer:**
   ```bash
   kill $VIEWER_PID 2>/dev/null
   ```

7. **Summarize improvements:**
   - What got better?
   - What stayed the same?
   - Any regressions?
   - Metrics comparison

### Checkpoint 7: Comparison Done

**Before proceeding, verify:**
- [ ] Benchmark analysis completed
- [ ] Eval viewer launched
- [ ] User reviewed comparison
- [ ] Feedback collected
- [ ] Improvements quantified
- [ ] Regressions identified (if any)
- [ ] User satisfied with direction

**If results are unsatisfactory, proceed to Step 8 for iteration.**

---

## Step 8: Iterate (If Needed)

**Goal:** Further refine the skill based on comparison results.

**Decision point:**
- If improvements are satisfactory → Skip to Step 9
- If more work needed → Continue iteration

### Actions

1. **Analyze what didn't work:**
   - Which tests still fail?
   - What feedback did user give?
   - Are there regressions?
   - What needs more work?

2. **Plan next iteration:**
   - Focus on remaining issues
   - Don't undo what worked
   - Consider different approaches
   - Keep changes incremental

3. **Apply refinements:**
   - Edit skill based on learnings
   - Try alternative approaches
   - Add missing pieces

4. **Run new iteration:**
   - Create `iteration-2/` directory
   - Run all tests again
   - Compare to iteration-1 (use `--previous-workspace`)
   - Collect feedback

5. **Repeat until:**
   - User is satisfied
   - Improvements plateau
   - Success criteria met

### Checkpoint 8: Improvements Done

**Before proceeding, verify:**
- [ ] All critical issues resolved
- [ ] Success criteria met
- [ ] No significant regressions
- [ ] User satisfied with improvements
- [ ] Metrics show improvement

**If issues remain, continue iterating.**

---

## Step 9: Final Validation

**Goal:** Ensure optimized skill meets quality standards.

### Actions

1. **Run Essential Quality checklist:**
   - [ ] Description complete (what + when)
   - [ ] Main file < 500 lines
   - [ ] At least 3 test cases
   - [ ] All tests pass
   - [ ] Terminology consistent

2. **Run Recommended Quality checklist:**
   - [ ] No time-sensitive info
   - [ ] Examples concrete
   - [ ] References organized
   - [ ] Progressive disclosure used
   - [ ] Workflows clear

3. **Run Optimization-Specific checklist:**
   - [ ] Baseline comparison done
   - [ ] Metrics improved (pass rate, speed, or quality)
   - [ ] No regressions in existing functionality
   - [ ] All user feedback addressed

4. **Final review with user:**
   - Show before/after comparison
   - Demonstrate improvements
   - Confirm it meets their needs

5. **Optional: Description optimization**
   
   If optimization goal included triggering improvements:
   ```bash
   python -m scripts.run_loop \
     --eval-set <trigger-evals.json> \
     --skill-path <skill-path> \
     --model <model-id> \
     --max-iterations 5
   ```

### Checkpoint 9: Quality Validated

**Before proceeding, verify:**
- [ ] All Essential Quality items pass
- [ ] All Optimization-Specific items pass
- [ ] Most Recommended Quality items pass
- [ ] User has done final review
- [ ] User confirms optimization is complete
- [ ] Description optimization done (if needed)

**Do not proceed until all Essential and Optimization-Specific items pass.**

---

## Step 10: Update & Deliver

**Goal:** Package and deliver the optimized skill.

### Actions

1. **Preserve the original name:**
   - Keep same directory name
   - Keep same `name` in frontmatter
   - Don't create v2 variants

2. **Package the skill:**
   ```bash
   python -m scripts.package_skill <path/to/skill-folder>
   ```

3. **Present to user:**
   - Show path to updated `.skill` file
   - Summarize improvements made
   - Provide metrics comparison
   - Include installation instructions

4. **Document changes:**
   - Update README or changelog
   - Note what was improved
   - Include before/after metrics

### Checkpoint 10: Skill Updated

**Final verification:**
- [ ] Skill packaged successfully
- [ ] Original name preserved
- [ ] `.skill` file created
- [ ] User knows where to find it
- [ ] Improvements documented
- [ ] Installation instructions provided

**Congratulations! Skill optimization complete.**

---

## Troubleshooting

### Common Issues

**No improvement in metrics:**
- Changes may not address root cause
- Try different approach
- Review transcripts for insights
- Ask user for more context

**Regressions introduced:**
- Changes too broad
- Unintended side effects
- Revert problematic changes
- Make more targeted improvements

**Baseline comparison unclear:**
- Assertions may be non-discriminating
- Add more specific checks
- Use qualitative comparison
- Focus on user feedback

**Triggering still problematic:**
- Description may need optimization
- Run description optimization loop
- Add more specific keywords
- Make description more "pushy"

---

## Optimization Strategies

### Quality Improvements

**Clarify instructions:**
- Add concrete examples
- Explain the "why"
- Remove ambiguity
- Use imperative form

**Handle edge cases:**
- Add specific guidance
- Provide fallback behavior
- Include error handling

**Improve examples:**
- Make them more realistic
- Add input/output pairs
- Cover common scenarios

### Triggering Improvements

**Enhance description:**
- Add specific keywords
- Include use case contexts
- Be more "pushy"
- Cover edge cases

**Run optimization loop:**
- Generate trigger evals
- Test current description
- Iterate improvements
- Validate on held-out set

### Performance Improvements

**Bundle repeated work:**
- Create helper scripts
- Reuse common patterns
- Avoid redundant steps

**Optimize workflow:**
- Remove unnecessary steps
- Parallelize when possible
- Use progressive disclosure

**Reduce token usage:**
- Remove verbose explanations
- Use references for details
- Keep main file lean

---

## Platform-Specific Notes

### Claude.ai
- No subagents → Run tests yourself sequentially
- No browser → Present results inline
- Compare outputs manually
- Skip description optimization

### Cowork
- Use `--static` for eval viewer
- Feedback downloads as file
- May need series execution if timeouts occur
- Description optimization works

---

## Summary

This workflow ensures systematic skill optimization with baseline comparison:

1. ✅ Analyze Current State → Identify problems
2. ✅ Define Optimization Goal → Set objectives
3. ✅ Create/Update Tests → Ensure coverage
4. ✅ Snapshot Baseline → Capture current performance
5. ✅ Apply Improvements → Make targeted changes
6. ✅ Run Evaluations → Test improved skill
7. ✅ Compare Results → Analyze improvements
8. ✅ Iterate → Refine further if needed
9. ✅ Final Validation → Quality checklist
10. ✅ Update & Deliver → Ship improved skill

**Each checkpoint ensures improvements are validated and no regressions occur.**
