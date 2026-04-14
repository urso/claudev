---
title: User Store
type: reference
tags: [auth, database, users]
scope: Schema and access patterns for the user data store
updates-when:
  - user table schema changes
watches:
  - migrations/
---

# User Store

## Schema

The `users` table stores credentials and profile data.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| email | TEXT | Unique, indexed |
| password_hash | TEXT | bcrypt hash |
| created_at | TIMESTAMP | Registration time |

## Access Patterns

- Lookup by email (login flow) — see [auth design](auth-design.md)
- Lookup by ID (session validation)
- List with pagination (admin panel)
