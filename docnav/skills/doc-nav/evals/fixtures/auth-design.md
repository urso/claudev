---
title: Authentication Architecture
type: design
tags: [auth, security, jwt]
scope: How the authentication system handles login, token refresh, and session management
updates-when:
  - auth middleware changes
watches:
  - internal/auth/
---

# Authentication Architecture

## Overview

The system uses JWT-based authentication with refresh tokens. Sessions are stored in Redis with a 24-hour TTL.

## Token Flow

1. User submits credentials to `/api/login`
2. Server validates against the user store
3. Returns access token (15min) and refresh token (7d)
4. Client includes access token in Authorization header

## Session Management

Sessions are tracked in Redis. See [user store](user-store.md) for the backing data model.

The [API reference](api-reference.md) documents all auth endpoints.
