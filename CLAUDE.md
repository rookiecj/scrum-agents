# Scrum Agents - Project Configuration

## Project Overview

GitHub Issues + Projects 기반 Scrum 개발 환경. Go 백엔드 + TypeScript 프론트엔드, 1주 스프린트.

## Repository

- **Owner**: rookiecj
- **Repo**: scrum-agents
- **Main branch**: main

## Tech Stack

- **Backend**: Go (in `backend/`)
- **Frontend**: TypeScript (in `frontend/`)

## Scrum Workflow

### Sprint Cycle
- Duration: 1 week (Monday to Friday)
- Sprint Planning: Monday
- Sprint Review/Retro: Friday

### Issue Labels
| Category | Labels |
|----------|--------|
| Type | `type:epic`, `type:story`, `type:task`, `type:bug`, `type:spike` |
| Priority | `priority:critical`, `priority:high`, `priority:medium`, `priority:low` |
| Sprint | `sprint:current`, `sprint:next`, `sprint:backlog` |
| Component | `component:backend`, `component:frontend` |
| Status | `status:planned`, `status:in-progress`, `status:dev-complete`, `status:in-review`, `status:verified`, `status:blocked` |

### Issue Workflow (Queue-Based)

Issues flow through a label-based state machine. Status labels are **mutually exclusive** — always remove the previous status before adding the next.

```
status:planned → status:in-progress → status:dev-complete → status:in-review → status:verified → CLOSED
  (DEV queue)     (DEV working)         (QA queue)            (QA working)       (QA passed)      (Dev merges)
       ↑                                                            |
       └────────────────────────────────────────────────────────────┘
                                                           (QA failed → rework)
```

**State Transitions:**

| From | To | Actor | Action |
|------|----|-------|--------|
| (new) | `status:planned` | Sprint Start | Enqueue to DEV queue |
| `status:planned` | `status:in-progress` | Dev Agent | Claim ticket |
| `status:in-progress` | `status:dev-complete` | Dev Agent | Implementation done, tests pass, PR created |
| `status:dev-complete` | `status:in-review` | QA Agent | Claim ticket for verification |
| `status:in-review` | `status:verified` | QA Agent | All AC verified (do NOT close issue) |
| `status:verified` | CLOSED | Scrum Master | Merge PR to main (`gh pr merge --squash`), close issue |
| `status:in-review` | `status:planned` | QA Agent | Verification failed, post results, rework needed |
| (any) | `status:blocked` | Any Agent | Blocker identified |
| `status:blocked` | `status:planned` | Scrum Master | Blocker resolved, return to DEV queue |
| `status:in-progress` | `status:planned` | Scrum Master | Agent crashed/abandoned, recover to DEV queue |
| `status:in-review` | `status:dev-complete` | Scrum Master | QA agent crashed, recover to QA queue |

**Workflow Steps:**
1. Product Owner creates issues with appropriate labels
2. Sprint Planning: issues labeled `sprint:next` (no `status:*` labels yet)
3. Sprint Start: `sprint:next` → `sprint:current` + `status:planned`
4. Dev Agent claims from queue → creates branch → implements → creates PR → marks `status:dev-complete`
5. QA Agent claims from queue → verifies AC → marks `status:verified` (or rework → `status:planned`)
6. Scrum Master merges verified PRs (`gh pr merge --squash`) → issue auto-closed via `Closes #N`
7. Scrum Master monitors queue health

### Branch Naming
- Feature: `feature/<issue-number>-<short-description>`
- Bug fix: `fix/<issue-number>-<short-description>`
- Spike: `spike/<issue-number>-<short-description>`

### Commit Messages
- Use conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`
- Reference issue number: `feat: add user login (#123)`

### PR Requirements
- Link related issue with `Closes #<number>`
- Fill out PR template
- Pass all CI checks
- At least one approval required

## Agent Roles

| Agent | Role | Key Responsibilities |
|-------|------|---------------------|
| `scrum-master` | Scrum Master | Sprint management, blocker resolution, progress tracking |
| `product-owner` | Product Owner | Backlog management, prioritization, story writing |
| `backend-dev` | Backend Dev | Go development, testing, PRs |
| `frontend-dev` | Frontend Dev | TypeScript development, testing, PRs |
| `qa` | QA Agent | Verification of dev-complete tickets, AC validation, rework decisions |
| `reviewer` | Code Reviewer | PR review, quality assurance (operates on PRs independently, not part of queue state machine) |

## Code Conventions

### Go Backend
- Follow standard Go project layout
- Use `internal/` for private packages
- Table-driven tests
- Handle all errors with wrapped messages
- Structured logging

### TypeScript Frontend
- Strict mode enabled
- Component-based architecture
- ESLint + Prettier
- Meaningful type definitions (avoid `any`)

## Queue Monitoring

```bash
# DEV Queue — tickets ready for development
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:planned" --state open

# In Progress — developers actively working
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:in-progress" --state open

# QA Queue — tickets awaiting verification
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:dev-complete" --state open

# In Review — QA actively verifying
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:in-review" --state open

# Verified — QA passed, ready to close
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:verified"

# Blocked
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:blocked" --state open
```

## Versioning

Each component has a `VERSION` file as the single source of truth:

| Component | Version File | How It's Consumed |
|-----------|-------------|-------------------|
| Backend | `backend/VERSION` | Injected at build time via `-ldflags "-X main.Version=$(cat VERSION)"` |
| Frontend | `frontend/VERSION` | Read by `vite.config.ts` → `__APP_VERSION__` global constant |

Both VERSION files must always contain the same semver value. Use `/sprint:release` to bump versions.

## Commands

```bash
# Backend
cd backend && go build -ldflags "-X main.Version=$(cat VERSION)" ./cmd/server
cd backend && go test ./... -v -cover
cd backend && go run -ldflags "-X main.Version=$(cat VERSION)" ./cmd/server

# Frontend
cd frontend && npm install
cd frontend && npm run dev
cd frontend && npm run build
cd frontend && npm test
```
