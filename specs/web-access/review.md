# Web Access Design Review

> Review of `specs/web-access/design.md` against the 7-dimension rubric. Structural headings are in English; finding content is in the design's current language.

## Verdict

- **Overall**: Revise
- **spec-plan readiness**: Conditional
- **Findings**: 0 Blocker · 2 Major · 2 Minor
- **Summary**: The design is structurally complete and fits the repository's Go-skill model, but it needs tighter command contracts and explicit old-skill routing changes before planning can be fully deterministic.

## Scope Reviewed

- **Reviewed file**: `specs/web-access/design.md`
- **Project context consulted**: `AGENTS.md`, `exa-search/SKILL.md`, `grok-search/SKILL.md`, `exa-search-go/cmd/*`, `exa-search-go/internal/config/loader.go`, `grok-search-go/cmd/*`, `grok-search-go/internal/config/*`, `.github/workflows/*search*.yml`, recent git history
- **Rubric**: 7 dimensions (D1 Completeness · D2 Usability · D3 Document Conformance · D4 Project Fit · D5 Blind Spots · D6 Over-Engineering · D7 Optimization)

## Findings

### Major

#### [M1] Command contracts are not concrete enough for the renamed and overlapping modes — D2 Usability / Actionability

- **Location**: `design.md §Proposed Solution / CLI Layer` lines 136-177; `design.md §Migration Documentation` lines 305-318; `exa-search-go/cmd/research.go`; `grok-search-go/cmd/research.go`
- **Issue**: The design lists command names and global flags, but it does not define per-command inputs, flags, defaults, and output modes for the merged CLI. This matters because `exa-search research` is mapped to `web-access extract`, while `web-access research` is assigned to Grok. Existing code shows these are materially different contracts: Exa `research` has `--num`, `--type`, `--text`, `--highlights`, date/domain/category/autoprompt flags and defaults text extraction on; Grok `research` only takes a research query plus global Grok options. `spec-plan` would have to infer the exact `extract` contract and how to avoid user confusion around the reused `research` name.
- **Recommendation**: Add a command contract table for every command with provider, required argument/flag, command-specific flags, default values, output payload family, and migration alias/rename notes. Call out explicitly that `extract` preserves Exa `research` behavior, including default text extraction.

#### [M2] Deprecation does not specify routing metadata changes for existing skills — D4 Project Fit

- **Location**: `design.md §Migration Documentation` line 320; `exa-search/SKILL.md` frontmatter; `grok-search/SKILL.md` frontmatter; `AGENTS.md §Skill File Format`
- **Issue**: The design says to update `exa-search/SKILL.md` and `grok-search/SKILL.md` so new work prefers `web-access`, but it does not state whether their YAML `description` fields must change. In this repository, `AGENTS.md` says the `description` field is critical for routing/matching, and both existing skill descriptions still strongly route source-first and real-time web work to the old skills. Updating only body text would leave the agent-facing routing signal ambiguous and undermine the goal of one reliable web access entry point.
- **Recommendation**: Specify the exact deprecation surface: update old skill frontmatter descriptions to point new work at `web-access`, keep compatibility usage in the body, and ensure the new `web-access/SKILL.md` description covers both source-first and live research routing.

### Minor

#### [m1] Release trigger language is slightly inconsistent with the repository workflow rule — D3 Document Conformance

- **Location**: `design.md §Goals` line 27; `design.md §Release Automation` lines 296-303; `AGENTS.md §GitHub Actions Workflow Requirements`
- **Issue**: The goal says "tag-only release builds", but the workflow section allows `web-access-release.yml` to run on `workflow_dispatch` as well. `AGENTS.md` says release workflows should be tag-only unless explicitly approved otherwise. Existing workflows are mixed (`db-explorer-release.yml` is tag-only; `grok-search-release.yml` has manual dispatch), so this is not a blocker, but the design should make the intended exception explicit if manual dispatch is required.
- **Recommendation**: Either remove `workflow_dispatch` from the release workflow requirement or add a short rationale that manual dispatch is an approved exception for `web-access`, while preserving no normal branch push trigger.

#### [m2] Cooldown state path semantics are underspecified — D5 Blind Spots

- **Location**: `design.md §Config Layer` lines 223-229; `grok-search-go/internal/config/template.go`
- **Issue**: The proposed config sets `state_file = "runtime/web-access-grok-cooldowns.json"`, but does not define what relative paths are resolved against. Existing Grok config also uses a relative cooldown path, so implementation can likely reuse that behavior, but the unified config is a good place to make the rule visible because users may run the binary from different working directories.
- **Recommendation**: State whether relative cooldown paths resolve against the config file directory, current working directory, or the user's home/config directory, and add that behavior to config tests.

## Dimension Summary

Status: ✓ sound · △ has findings to address · ✗ blocking issue.

| Dimension | Status | Finding refs |
|-----------|--------|--------------|
| D1 Completeness (完整性) | ✓ | — |
| D2 Usability (可用性) | △ | M1 |
| D3 Document Conformance (规范性) | △ | m1 |
| D4 Project Fit (符合项目规范) | △ | M2 |
| D5 Blind Spots (盲点) | △ | m2 |
| D6 Over-Engineering (过度设计) | ✓ | — |
| D7 Optimization (优化点) | ✓ | — |

## Recommended Next Step

Address the Major findings before `spec-plan`: add the per-command contract table and define the old-skill frontmatter deprecation behavior. After those are fixed, the design should be ready for `spec-plan`; the Minor items can be folded into the same revision or captured as planning details.
