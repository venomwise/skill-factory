# Bulk Order Export Design

## Summary

Add a `POST /orders/export` endpoint that generates a CSV of a merchant's orders and emails a download link when ready.

## Goals

- Let merchants export up to 1M orders to CSV.
- Deliver the export within 5 minutes of the request.
- Notify the merchant by email and in-app when the export is ready.

## Primary Users / Roles

- Merchants exporting their order history for accounting.

## Non-Goals

- Scheduled/recurring exports.

## Context

- Orders live in the `orders` table; there is an existing transactional email service.

## Decision Record

### Options Considered

- TBD.

### Decision & Rationale

To be decided during implementation.

## Proposed Solution

The endpoint kicks off a background job that queries orders, writes a CSV to object storage, and sends the link.

### Architecture

Handler → background worker → object storage → email.

### Components

- **Export handler**: accepts the request and returns `202`.
- **Export worker**: builds the CSV.

### Data Flow

1. Merchant calls `POST /orders/export`.
2. A job is enqueued.
3. The worker writes the CSV and emails the link.
