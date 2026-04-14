---
title: Deployment Runbook
type: runbook
tags: [ops, deployment, infrastructure]
scope: Step-by-step production deployment procedure
updates-when:
  - CI/CD pipeline changes
  - new service added
watches:
  - docker-compose.yml
  - .github/workflows/
---

# Deployment Runbook

## Prerequisites

- Docker and docker-compose installed
- Access to production secrets vault
- VPN connection active

## Steps

1. Pull latest from main
2. Run `make build`
3. Run integration tests: `make test-integration`
4. Deploy with `docker-compose up -d`
5. Verify health checks pass

## Rollback

If health checks fail, run `docker-compose rollback` to revert to the previous image.
