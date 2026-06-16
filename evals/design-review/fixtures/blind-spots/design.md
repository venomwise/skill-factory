# Nightly CRM Contact Sync Design

## Summary

Every night, pull all contacts from the third-party Acme CRM API and upsert them into our `contacts` table so sales reps see up-to-date contact data in our app.

## Goals

- Keep our `contacts` table in sync with Acme CRM, refreshed once per night.
- Reflect new and changed contacts from Acme within 24 hours.

## Primary Users / Roles

- Sales reps viewing contact records in our app.

## Non-Goals

- Real-time sync (nightly batch is acceptable).
- Writing changes back to Acme (read-only from our side).

## Context

- Acme exposes `GET /v1/contacts` returning all contacts as JSON pages.
- We have a `contacts` table keyed by our own `id`, plus an `acme_id` column.
- A cron runner already exists for nightly jobs.

## Decision Record

### Options Considered

- **Full reload each night**: simple; fetch everything and upsert by `acme_id`.
- **Incremental by `updated_since`**: fetch only changed records; smaller payloads but depends on Acme's change timestamps.

### Decision & Rationale

Full reload each night — simplest to implement and Acme's contact volume is small enough to fetch fully.

## Proposed Solution

A nightly cron job fetches all Acme contacts page by page and upserts each into `contacts` by `acme_id`.

### Architecture

Cron trigger → sync job → Acme client (paged GET) → upsert into `contacts`.

### Components

- **Sync job**: orchestrates paging and upsert; runs at 02:00 daily.
- **Acme client**: wraps `GET /v1/contacts?page=`; returns parsed contacts.
- **Contact upserter**: `INSERT ... ON CONFLICT (acme_id) DO UPDATE` for each record.

### Data Flow

1. Cron triggers the sync job at 02:00.
2. The job pages through `GET /v1/contacts` until no more pages.
3. Each contact is upserted into `contacts` by `acme_id`.
4. The job logs how many contacts were processed.

## Error Handling

- If a page request fails, the job logs the error and stops.

## Testing

- Unit: upserter maps Acme fields to our columns correctly.
- Integration: a full run upserts a sample set of contacts.

## Open Questions

None. Key decisions were confirmed during discovery.
