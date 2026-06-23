---
id: "<ID>"
title: "<TITLE>"
status: ready
created: <DATE>
start-git-sha: ""
blocked-by: []
description: "<ONE-LINE SUMMARY>"
---

# <TITLE>

## Tasks

### Task 1: <Task Name>

#### Acceptance Criteria
<!-- Outcomes verifiable by looking at the code — name types, signatures, behaviors -->
- [ ] <Verifiable outcome 1>
- [ ] <Verifiable outcome 2>

#### Subtasks
<!-- Specific work steps — name files, functions, operations -->
- [ ] <Sub-task 1>
- [ ] <Sub-task 2>
- [ ] <Sub-task 3>

<!-- Example:
#### Acceptance Criteria
- [ ] `AuthMiddleware` type exists with `Wrap(http.Handler) http.Handler` method
- [ ] Unauthenticated requests return 401 with JSON error body
- [ ] Context contains `UserID` after successful auth

#### Subtasks
- [ ] Define `AuthMiddleware` struct in `pkg/auth/middleware.go`
- [ ] Add `Wrap` method that extracts Bearer token from Authorization header
- [ ] Call `TokenValidator.Validate(token)` and set `UserID` in context
- [ ] Return `{"error": "unauthorized"}` with 401 if validation fails
- [ ] Register middleware in `cmd/server/main.go` router setup
-->

### Task 2: <Task Name>

#### Acceptance Criteria
- [ ] <Verifiable outcome>

#### Subtasks
- [ ] <Sub-task 1>
- [ ] <Sub-task 2>

## Technical Notes
<Implementation guidance, architecture context>

## Developer Logs

### Decision Log
<Key technical decisions and rationale>

### Blockers Encountered
<Issues faced and resolutions>

### Deviations from Design
<Changes from original design and why>

### Lessons Learned
<Insights for future stories>
