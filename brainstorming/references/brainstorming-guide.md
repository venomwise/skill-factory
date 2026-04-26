# Brainstorming Guide

## Context exploration priorities

1. README, CLAUDE.md, or equivalent project docs.
2. Project config: package.json, pyproject.toml, Cargo.toml, go.mod.
3. Entry points: main files, index files, route definitions.
4. Recent commits (up to 10) for current momentum and conventions.
5. Stop when you can state the project's purpose, tech stack, and directory structure.

## Decomposition heuristics

- Split when the request spans multiple domains or independent capabilities.
- Split when parts would require different data models or external integrations.
- Split when parts could be owned by different teams or deployed separately.
- When splitting, list sub-projects, dependencies, and a recommended order.

## Good scope split examples

- "Chat + file storage + billing" becomes three sub-projects with their own specs.
- "Redesign UI + add analytics" becomes a UI sub-project and a metrics sub-project.

## Question patterns

- Purpose: who is the user and what decision does this enable?
- Constraints: latency, scale, security, compliance, or operational limits.
- Success criteria: how will we know it worked?
- Data: inputs, outputs, and the source of truth.
- Risks: key failure modes or edge cases to handle.

## Design depth guidance

- Simple changes still use all sections, but keep each section to a few sentences.
- Moderate changes include a step-by-step data flow and key error cases.
- Complex changes include primary happy path plus 2-3 critical failure scenarios.

## Existing codebase tactics

- Identify seams to reuse and respect local patterns.
- Refactors are allowed only when they unblock the current goal.
- Prefer changes that reduce coupling without altering unrelated behavior.

## Review gate checklist

- The design doc covers architecture, components, data flow, error handling, and testing.
- Open questions are listed with an owner or next step.
- Scope matches a single sub-project with explicit non-goals.

## Fuzziness diagnosis

Assess the user's request against these levels before choosing a strategy:

| Level | Signal | Example | Strategy |
|-------|--------|---------|----------|
| Problem unclear | User describes symptoms, not goals; "it feels wrong" | "The system is hard to use" | Reframe the problem |
| Direction unclear | User has a goal but no sense of approach | "I want better user engagement" | Explore possibilities |
| Boundaries unclear | User knows what they want but not the edges | "I want user auth" | Scan for blind spots |
| Solution unclear | User knows what and scope, needs technical approach | "I want SSO with SAML" | Compare approaches (skip to step 4) |

## Assumption challenging

For non-trivial requests, identify potential assumptions and offer them to the user for confirmation. Frame them as "worth checking" rather than "you missed this".

1. Identify 1-3 assumptions the user may be implicitly making.
2. For each, ask: "Is this necessarily true? What if it were not?"
3. Common assumptions worth checking:
   - "Users will use this feature the way I imagine"
   - "The current architecture can support this without changes"
   - "This needs to be done in one release"
   - "Performance or scale will not be an issue"
   - "The existing data model is sufficient"

## Blind spot scanning

Pick the 1-2 checklists most relevant to the project's domain. Do not walk through every list; select by domain and skip items that clearly do not apply.

**Any user-facing feature:**
- Empty states, loading states, error states
- Permissions and access control
- Undo or rollback
- Offline behavior
- Accessibility
- Internationalization or localization
- Mobile or responsive behavior

**Any data feature:**
- Data migration from existing state
- Consistency and conflict resolution
- Retention and cleanup policies
- Privacy and compliance (GDPR, etc.)
- Backup and recovery

**Any integration:**
- Rate limits and quotas
- Authentication and credential rotation
- Failure modes and fallback behavior
- Versioning and backwards compatibility
- Monitoring and alerting

## Divergent exploration techniques

Use these during brainstorming to expand the problem space:

- **Constraint removal**: "If we had no technical, time, or budget limits, what would the ideal look like?" Then add constraints back one at a time.
- **Negative brainstorming**: "What would make this feature fail completely?" Flip each failure into a requirement.
- **Perspective switching**: Consider the same feature from different roles (end user, admin, ops engineer, new hire, power user).
- **Time horizon**: "How will this need to evolve in 6 months? A year?" Identify which decisions are hard to reverse.
- **Priority forcing**: "If you could ship only 3 things, which 3?" This forces the user to reveal what truly matters.
- **Analogy**: "How do similar products or domains solve this?" Borrow proven patterns.

## When to stop brainstorming

- The user can state what they want, why they want it, and what they do not want.
- Key constraints and success criteria are known or explicitly recorded as assumptions.
- For non-trivial requests, at least one assumption has been confirmed and one blind spot has been considered.
- The problem statement is stable and has not changed in the last 2 exchanges.
