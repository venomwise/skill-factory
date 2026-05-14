# Clarification Benchmark - Iteration 1

| Eval | Config | Passed | Total | Pass rate | Time |
|---|---:|---:|---:|---:|---:|
| 1 dry-run-small-change | with_skill | 4 | 4 | 100% | 34.4s |
| 1 dry-run-small-change | without_skill | 4 | 4 | 100% | 62.5s |
| 2 rbac-hidden-large-change | with_skill | 3 | 4 | 75% | 45.1s |
| 2 rbac-hidden-large-change | without_skill | 3 | 4 | 75% | 31.5s |
| 3 ambiguous-priority-sort | with_skill | 4 | 4 | 100% | 51.2s |
| 3 ambiguous-priority-sort | without_skill | 4 | 4 | 100% | 69.5s |

## Aggregate

| Config | Passed | Total | Pass rate | Total time |
|---|---:|---:|---:|---:|
| with_skill | 11 | 12 | 92% | 130.7s |
| without_skill | 11 | 12 | 92% | 163.5s |

## Notes

- Both with-skill and baseline failed the explicit brainstorming handoff assertion for RBAC.
- Eval 3 with-skill met behavioral clarification assertions but produced a less ideal tie-break implementation than baseline; future evals should add code correctness assertions.
