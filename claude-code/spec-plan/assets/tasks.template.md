# Implementation Plan: <Project Name>

## Overview

This plan implements the requirements defined in [requirements.md](requirements.md).

<!-- Summarize:
     - Number of phases and overall approach
     - Execution order rationale (what comes first and why)
     - Key technical decisions (language, frameworks, dependencies)
     Keep it brief — the detail lives in the tasks below. -->

<Summarize the implementation strategy here.>

## Tasks

<!-- FORMAT RULES:

1. Phase heading — use checkbox list item, not markdown heading:
   CORRECT:  - [ ] Phase 1: Title
   WRONG:    ### Phase 1: Title

2. Task details — indented bullets under the task title:
   CORRECT:
     - [ ] 1.1 Task title
       - Create `path/to/file.py` with `function_name` implementing <behavior>
       - _Requirements: 1.1, 1.2_

3. Optional tasks — append [optional] to the task title:
   CORRECT:  - [ ] 1.3 Write unit tests for feature [optional]
   Use for: test tasks, verification tasks, nice-to-have features

4. Requirement references — use "N.M" format:
   CORRECT:  _Requirements: 1.1, 2.3_
   WRONG:    _Requirements: R1, R2_

5. Checkpoints — insert between major phases:
   - [ ] Checkpoint: Verify <what to validate before proceeding>

6. Task size — each task should be completable in one focused session.
   If it needs sub-phases, split it into multiple tasks. -->

- [ ] Phase 1: <Title>
  - [ ] 1.1 <Task title>
    - <Concrete step: create/modify file, implement function, configure setting>
    - <Additional steps as needed>
    - _Requirements: 1.1, 1.2_
  - [ ] 1.2 <Task title>
    - <Concrete steps>
    - _Requirements: 1.3_
  - [ ] 1.3 Write tests for <feature> [optional]
    - Test <normal scenario>
    - Test <error scenario>
    - Test <boundary condition>
    - _Requirements: 1.1, 1.2, 1.3_

- [ ] Checkpoint: Verify <scope of Phase 1>

- [ ] Phase 2: <Title>
  - [ ] 2.1 <Task title>
    - <Concrete steps>
    - _Requirements: 2.1_
  - [ ] 2.2 Write tests for <feature> [optional]
    - <Test scenarios>
    - _Requirements: 2.1_

## Notes

<!-- Include as applicable:
     - Architecture decisions and rationale
     - Dependencies between phases
     - Testing strategy (unit, integration, e2e)
     - Known risks or open questions -->

- Tasks marked `[optional]` can be skipped for an MVP.
- Every task references requirement IDs for traceability.
