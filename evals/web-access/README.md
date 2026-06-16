# Web Access Skill Evaluation

This directory contains evaluation cases for the `web-access` unified CLI skill.

## Purpose

These eval cases validate:
- Correct routing decisions between Exa (source-first) and Grok (live synthesis) providers
- Command flag and default behavior
- Output format variations (JSON, plain, URLs-only)
- Multi-profile failover behavior

## Structure

- `evals.json` - Structured evaluation cases with expected routing and example commands

## Usage

Eval cases focus on **routing correctness** and **command structure**, not network-dependent result accuracy.

Each case includes:
- Scenario description
- Expected provider routing and reasoning
- Example command demonstrating correct usage
- Acceptance criteria for validation

## Validation Approach

- **Routing decisions**: Verified through unit tests in `web-access-go/cmd/*_test.go`
- **Command defaults**: Verified through CLI flag tests
- **Failover behavior**: Verified through mock HTTP server tests in provider packages
- **Output formats**: Verified through output package tests

Network-dependent live results (actual Exa/Grok API responses) are not evaluated for correctness, as they change over time and depend on external services.
