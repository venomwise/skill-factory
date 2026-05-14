# Clarification Benchmark - Iteration 2

| Eval | Config | Passed | Total | Pass rate | Time |
|---|---:|---:|---:|---:|---:|
| 1 dry-run-small-change | with_skill | 4 | 4 | 100% | 47.8s |
| 1 dry-run-small-change | without_skill | 4 | 4 | 100% | 50.7s |
| 2 rbac-hidden-large-change | with_skill | 3 | 4 | 75% | 33.6s |
| 2 rbac-hidden-large-change | without_skill | 1 | 4 | 25% | 117.8s |
| 3 ambiguous-priority-sort | with_skill | 2 | 4 | 50% | 83.4s |
| 3 ambiguous-priority-sort | without_skill | 2 | 4 | 50% | 54.0s |

## Aggregate

| Config | Passed | Total | Pass rate | Total time |
|---|---:|---:|---:|---:|
| with_skill | 9 | 12 | 75% | 164.8s |
| without_skill | 7 | 12 | 58% | 222.5s |

## Notes

- With-skill improved over baseline on RBAC by refusing to code, but still failed the explicit `brainstorming` wording assertion.
- Eval 3 assertion was tightened: both configurations still implemented after assumptions instead of asking sorting blind-spot questions first.
- Next iteration should make `brainstorming` wording and user-visible behavior-change questions more forceful/earlier in the workflow.
