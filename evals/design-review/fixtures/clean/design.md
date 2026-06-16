# Orders List Pagination Design

## Summary

Add cursor-based pagination to the existing `GET /orders` endpoint so clients can page through large order sets without missing or duplicating rows when new orders arrive during paging.

## Goals

- Return at most 50 orders per page, ordered by `(created_at, id)`.
- Stable paging: no missed or duplicated orders when rows are inserted mid-pagination.
- p95 latency per page under 150 ms at the current 2M-row table size.

## Primary Users / Roles

- Internal dashboard clients listing a merchant's orders.
- Third-party API consumers exporting their own orders.

## Non-Goals

- Filtering or full-text search over orders (handled by a separate search endpoint).
- Changing the order schema or write path.

## Context

- `GET /orders` currently returns all rows for a merchant with no limit (`internal/api/orders.go`).
- The `orders` table already has a primary key `id` and a `created_at` column; an index on `(merchant_id, created_at, id)` exists.

## Decision Record

### Options Considered

- **Offset pagination** (`LIMIT/OFFSET`): trivial to implement, but skips or repeats rows when inserts happen during paging, and OFFSET scans degrade on deep pages.
- **Cursor pagination** on `(created_at, id)`: stable under concurrent inserts and uses the existing composite index; slightly more client-side complexity (opaque cursor token).

### Decision & Rationale

Cursor pagination. The stability goal rules out offset, and the existing `(merchant_id, created_at, id)` index makes the keyset query cheap. The cursor is an opaque base64 of `(created_at, id)`.

## Proposed Solution

High-level: keyset pagination using the existing composite index, with an opaque cursor token returned to the client.

### Architecture

`GET /orders?limit=&cursor=` → orders handler → repository keyset query → response with `data[]` and `next_cursor`.

### Components

- **Orders handler** (`internal/api/orders.go`): parses `limit` (default 50, max 50) and `cursor`; returns `{ data, next_cursor }`. `next_cursor` is null on the last page.
- **Cursor codec** (`internal/api/cursor.go`, new): encodes/decodes `(created_at, id)` to/from a base64 token; rejects malformed tokens.
- **Orders repository** (`internal/store/orders.go`): `ListAfter(merchantID, createdAt, id, limit)` runs `WHERE merchant_id=? AND (created_at, id) > (?, ?) ORDER BY created_at, id LIMIT ?`.

### Data Flow

1. Client calls `GET /orders` (no cursor) → repository runs the keyset query from the start → handler returns 50 rows + `next_cursor` from the last row.
2. Client calls `GET /orders?cursor=<token>` → codec decodes `(created_at, id)` → repository fetches the next 50 rows after that key.
3. When fewer than `limit` rows return, `next_cursor` is null.

## Error Handling

- Malformed or tampered `cursor` → `400` with `{ "error": "invalid_cursor" }`.
- `limit` out of range → clamped to `[1, 50]`, not an error.
- Repository/database error → `500`, logged with the merchant ID and decoded cursor key.

## Testing

- Unit: cursor codec round-trips `(created_at, id)`; rejects malformed tokens.
- Repository: keyset query returns correct rows and ordering; last page returns null cursor.
- Integration: inserting a row mid-pagination does not cause a missed or duplicated order across pages.

## Open Questions

None. Key decisions were confirmed during discovery.
