# Clarification Benchmark - Iteration 3

| Eval | Config | Passed | Total | Pass rate | Time |
|---|---:|---:|---:|---:|---:|
| 1 dry-run-small-change | with_skill | 4 | 4 | 100% | 49.6s |
| 1 dry-run-small-change | without_skill | 4 | 4 | 100% | 36.8s |
| 2 rbac-hidden-large-change | with_skill | 3 | 4 | 75% | 26.9s |
| 2 rbac-hidden-large-change | without_skill | 2 | 4 | 50% | 48.9s |
| 3 ambiguous-priority-sort | with_skill | 2 | 4 | 50% | 58.1s |
| 3 ambiguous-priority-sort | without_skill | 2 | 4 | 50% | 51.1s |

## Aggregate

| Config | Passed | Total | Pass rate | Total time |
|---|---:|---:|---:|---:|
| with_skill | 9 | 12 | 75% | 134.6s |
| without_skill | 8 | 12 | 67% | 136.9s |

## Notes

- Generalized wording avoided RBAC/sorting-specific hardcoding, but behavior did not materially improve over iteration 2.
- With-skill still says design/spec but not explicit `brainstorming`.
- With-skill still treats user-visible behavior assumptions as safe enough to implement instead of asking first.
