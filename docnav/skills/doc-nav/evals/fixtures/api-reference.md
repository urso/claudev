---
title: API Reference
type: reference
tags: [api, http, endpoints]
scope: Complete REST API endpoint documentation
---

# API Reference

## Authentication Endpoints

### POST /api/login

Authenticates a user and returns tokens.

### POST /api/refresh

Refreshes an expired access token.

### POST /api/logout

Invalidates the current session.

## User Endpoints

### GET /api/users/:id

Returns user profile. Requires authentication — see [auth design](auth-design.md).

### PUT /api/users/:id

Updates user profile.

## Deployment

See the [deployment runbook](deployment-runbook.md) for production setup.
