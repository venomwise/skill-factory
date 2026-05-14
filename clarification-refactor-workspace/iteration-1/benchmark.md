# Clarification Refactor Benchmark - Iteration 1

| Eval | Config | Passed | Total | Pass rate | Time |
|---|---:|---:|---:|---:|---:|
| 1 dry-run-small-change | with_skill | 3 | 4 | 75% | 35.2s |
| 1 dry-run-small-change | without_skill | 3 | 4 | 75% | 31.0s |
| 2 rbac-hidden-large-change | with_skill | 2 | 4 | 50% | 25.1s |
| 2 rbac-hidden-large-change | without_skill | 2 | 4 | 50% | 25.9s |
| 3 ambiguous-priority-sort | with_skill | 1 | 4 | 25% | 27.3s |
| 3 ambiguous-priority-sort | without_skill | 1 | 4 | 25% | 31.7s |

## Aggregate

| Config | Passed | Total | Pass rate | Total time |
|---|---:|---:|---:|---:|
| with_skill | 6 | 12 | 50% | 87.7s |
| without_skill | 6 | 12 | 50% | 88.5s |

## Notes

- Refactored skill did not materially change behavior versus baseline in this harness.
- The model did not perform the required `brainstorming` handoff for cross-cutting RBAC.
- The model still directly implemented user-visible ordering behavior without asking blind-spot questions.
- This suggests either the skill is not being selected/applied strongly enough by metadata-only triggering, or the skill needs an even shorter, more front-loaded operational contract.
