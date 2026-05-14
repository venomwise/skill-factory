# Clarification Benchmark - Iteration 5

Natural prompts; removed harness instructions that explicitly encouraged direct modification.

| Eval | Config | Passed | Total | Pass rate | Time |
|---|---:|---:|---:|---:|---:|
| 1 dry-run-small-change | with_skill | 4 | 4 | 100% | 43.7s |
| 1 dry-run-small-change | without_skill | 4 | 4 | 100% | 30.3s |
| 2 rbac-hidden-large-change | with_skill | 3 | 4 | 75% | 26.9s |
| 2 rbac-hidden-large-change | without_skill | 3 | 4 | 75% | 30.4s |
| 3 ambiguous-priority-sort | with_skill | 1 | 4 | 25% | 28.9s |
| 3 ambiguous-priority-sort | without_skill | 1 | 4 | 25% | 23.6s |

## Aggregate

| Config | Passed | Total | Pass rate | Total time |
|---|---:|---:|---:|---:|
| with_skill | 8 | 12 | 67% | 99.6s |
| without_skill | 8 | 12 | 67% | 84.2s |

## Notes

- Removing direct-implementation harness text did not improve the two target behaviors.
- With-skill still asks for framework/files rather than recommending `brainstorming` on cross-cutting RBAC.
- With-skill still directly edits user-visible ordering behavior and even regresses preserved created_at semantics by choosing ascending order.
- This suggests the skill is either not being triggered strongly enough or instructions need to be more front-loaded/operational.
