# Scrum Agents

GitHub Issues + Projects 기반 Scrum 프로세스를 운영하고, Claude Code Agent Teams로 개발 워크플로우를 자동화하는 프로젝트.

## Tech Stack

- **Backend**: Go
- **Frontend**: TypeScript
- **Sprint Cycle**: 1주

## Project Structure

```
scrum-agents/
├── backend/          # Go backend
├── frontend/         # TypeScript frontend
├── .claude/agents/   # Claude Code Agent configurations
├── .github/          # Issue templates, PR template, workflows
└── CLAUDE.md         # Project conventions
```

## Getting Started

```bash
# Backend
cd backend
go mod tidy
go run ./cmd/server

# Frontend
cd frontend
npm install
npm run dev
```

## Scrum Workflow

1. **Backlog Grooming**: Product Owner가 User Story 작성 및 우선순위 설정
2. **Sprint Planning**: Sprint에 포함할 Issue 선정, Story Point 할당
3. **Daily**: GitHub Project Board로 진행상황 트래킹
4. **Sprint Review**: PR 리뷰 및 Done 처리
5. **Retrospective**: Sprint 종료 후 회고

## Labels

| Category | Labels |
|----------|--------|
| Type | `type:story`, `type:task`, `type:bug`, `type:spike` |
| Priority | `priority:critical`, `priority:high`, `priority:medium`, `priority:low` |
| Sprint | `sprint:current`, `sprint:next`, `sprint:backlog` |
| Component | `component:backend`, `component:frontend` |
| Status | `status:blocked`, `status:review` |
