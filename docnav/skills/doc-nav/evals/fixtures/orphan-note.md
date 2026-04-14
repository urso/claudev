---
title: Performance Tuning Notes
type: note
tags: [performance, database, optimization]
scope: Notes from the Q3 performance investigation
---

# Performance Tuning Notes

## Findings

- Connection pool size was too small (default 10, increased to 50)
- Missing index on `users.email` caused full table scans during login
- Redis TTL was set to 1 hour, increased to 24 hours to reduce auth round-trips

## Recommendations

- Add database query logging in staging
- Set up Grafana dashboard for connection pool utilization
- Consider read replicas for user lookups
