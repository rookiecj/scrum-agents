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
| Type | `type:story`, `type:task`, `type:bug`, `type:spike` |
| Priority | `priority:critical`, `priority:high`, `priority:medium`, `priority:low` |
| Sprint | `sprint:current`, `sprint:next`, `sprint:backlog` |
| Component | `component:backend`, `component:frontend` |
| Status | `status:blocked`, `status:review` |

### Issue Workflow
1. Product Owner creates issues with appropriate labels
2. Sprint Planning: issues labeled `sprint:current`
3. Developer picks issue → creates branch `feature/<issue-number>-<description>`
4. Development → PR → Code Review → Merge
5. Scrum Master updates Project Board status

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
| `reviewer` | Code Reviewer | PR review, quality assurance |

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

## Commands

```bash
# Backend
cd backend && go build ./...
cd backend && go test ./... -v -cover
cd backend && go run ./cmd/server

# Frontend
cd frontend && npm install
cd frontend && npm run dev
cd frontend && npm run build
cd frontend && npm test
```
