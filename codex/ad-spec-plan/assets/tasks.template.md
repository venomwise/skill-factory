# Implementation Plan: <Project Name>

## Overview

This implementation plan is driven by the requirements in [requirements.md](requirements.md).

<!-- Describe:
     - How many phases and the overall approach
     - The execution order rationale (what comes first and why)
     - Key technical decisions (language, frameworks, dependencies)
-->

<Summarize the implementation strategy: number of phases, execution order rationale, and key technical decisions.>

## Tasks

<!-- FORMAT RULES — follow these exactly:

1. Phase heading: use checkbox list item, NOT markdown heading
   CORRECT:  - [ ] 1. Phase 1: Title
   WRONG:    ### Phase 1: Title

2. Task descriptions: use indented bullet list under the task title
   CORRECT:
     - [ ] 1.1 Task title
       - Create `src/path/file.py` with `function_name` implementing behavior
       - _Requirements: 1.1, 1.2_
   WRONG:
     - [ ] 1.1 Task title
     Create src/path/file.py with function_name...
     _Requirements: R1_

3. Optional tasks: use `*` after the checkbox brackets
   Phase level:     - [ ]* 4. Optional Phase: Title
   Sub-task level:  - [ ]* 1.3 Write unit tests for feature
   WRONG:           - [ ] 1.3 可选：Write unit tests
   Mark as optional: test tasks, verification tasks, and nice-to-have features
   Sub-tasks under an optional Phase inherit that optional status and do not need their own `*`.

4. Requirements references: use "RequirementNumber.CriterionNumber" format
   CORRECT:  _Requirements: 1.1, 1.2, 2.3_
   WRONG:    _Requirements: R1, R2_

5. Each phase should contain implementation tasks + verification/test task(s)
-->

- [ ] 1. Phase 1: Example Phase (replace)
  - [ ] 1.1 Example task title
    - Create `src/module/file.py` with `function_name` implementing <behavior>
    - Integrate with existing `src/other/module.py` by importing `ExistingClass`
    - _Requirements: 1.1, 1.2_
  - [ ] 1.2 Another implementation task
    - Modify `src/module/config.py` to add `<new_setting>` with default value
    - Update `src/module/__init__.py` to export the new component
    - _Requirements: 1.3, 1.4_
  - [ ]* 1.3 Write unit tests for <feature>
    - Test <normal scenario>
    - Test <error scenario>
    - Test <boundary condition>
    - _Requirements: 1.1, 1.2, 1.3_

- [ ] 2. Phase 2: <Phase Title> (replace)
  - [ ] 2.1 <Task Title>
    - Create/modify `src/<module>/<file>` with `<function/class>` implementing <behavior>
    - <Additional concrete steps or deliverables>
    - _Requirements: 2.1_
  - [ ]* 2.2 Write unit tests for <Phase 2 feature>
    - Test <normal scenario>
    - Test <error scenario>
    - Test <boundary condition>
    - _Requirements: 2.1_

- [ ] 3. Checkpoint - Verify <scope>
  - Ensure all tests pass, ask the user if questions arise.

- [ ]* 4. Optional Phase: <Phase Title>
  - [ ] 4.1 <Task Title>
    - <Concrete steps or deliverables>
    - _Requirements: 4.1_

## Notes

- Tasks marked with `*` are optional and can be skipped for an MVP.
- Each task should reference one or more requirement IDs for traceability.
- Keep task numbering stable so requirement references stay valid.

<!-- Also include as applicable:
     - Implementation language and key frameworks
     - Architecture decisions and rationale
     - Module extraction / reuse strategy
     - Testing strategy (unit, integration, property-based)
-->
