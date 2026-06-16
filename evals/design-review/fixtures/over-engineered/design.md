# Theme Preference Design

## Summary

Persist a single user setting — preferred UI theme, either `light` or `dark` — so the app renders the user's last choice on next login.

## Goals

- Store one value per user: `theme ∈ {light, dark}`.
- Read it on login to render the correct theme.
- Default to `light` when unset.

## Primary Users / Roles

- End users toggling light/dark mode in settings.

## Non-Goals

- Any preference other than theme (no language, density, notifications, etc.).
- Per-device or per-session themes — one value per user account.
- Custom or user-defined themes beyond light/dark.

## Context

- Users table exists; the app has a settings page with a theme toggle that currently only affects the current session and is lost on logout.

## Decision Record

### Options Considered

- **Single column `users.theme`**: one column, read on login. Minimal.
- **Generic preference framework**: a pluggable engine supporting arbitrary typed preferences with multiple storage backends.

### Decision & Rationale

Chose the generic preference framework so future preferences can be added without schema changes.

## Proposed Solution

A `PreferenceEngine` with a plugin registry, pluggable storage backends, a typed schema system, and an event bus that broadcasts preference changes.

### Architecture

`PreferenceEngine` → `PreferenceProvider` plugins (one per preference type) → `StorageBackend` interface (SQL, Redis, in-memory implementations) → `PreferenceChangedEvent` published on an internal event bus → subscribers (cache invalidator, audit logger, future webhooks).

### Components

- **PreferenceEngine**: orchestrates provider lookup, validation, read/write, and event publication.
- **ThemePreferenceProvider**: the single concrete provider, validating `theme ∈ {light, dark}`.
- **StorageBackend (interface)** with `SqlBackend`, `RedisBackend`, `InMemoryBackend` implementations; backend selected via config.
- **PreferenceCache**: write-through cache layer over the backend with a 5-minute TTL.
- **EventBus**: in-process pub/sub for `PreferenceChangedEvent`; an audit subscriber records every change.

### Data Flow

1. User toggles theme → `PreferenceEngine.set(userId, "theme", value)`.
2. Engine resolves `ThemePreferenceProvider`, validates, writes through `PreferenceCache` to the configured `StorageBackend`.
3. Engine publishes `PreferenceChangedEvent`; the audit subscriber logs it.
4. On login, `PreferenceEngine.get(userId, "theme")` reads cache, falling back to backend, defaulting to `light`.

## Error Handling

- Invalid theme value → validation error from the provider.
- Backend unavailable → fall back to the next configured backend, then to default.
- Event bus delivery failure → retried by the bus with backoff.

## Testing

- Unit: provider validation; cache hit/miss; backend selection; event publication.
- Integration: set then get returns the stored theme across each backend implementation.

## Open Questions

None. Key decisions were confirmed during discovery.
