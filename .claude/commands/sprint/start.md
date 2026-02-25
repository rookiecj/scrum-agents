You are the Scrum Master starting a planned sprint. Activate the sprint, enqueue tickets to the DEV queue, and dispatch agents.

## Input
Sprint start arguments: $ARGUMENTS

## Instructions

### 1. Identify the Sprint

- If sprint name/number is given in arguments, use it
- Otherwise, find the planned (next) sprint:
  ```bash
  gh issue list -R rookiecj/scrum-agents -l "sprint:next" --state open
  ```
- If no `sprint:next` tickets exist, inform the user they need to run `/sprint plan` first

### 2. Activate the Sprint & Enqueue to DEV Queue

First, validate that no `sprint:next` tickets have pre-existing `status:*` labels (they should not have been assigned during planning):
```bash
# Check for incorrectly labeled tickets
gh issue list -R rookiecj/scrum-agents -l "sprint:next" --state open --json number,labels | \
  jq '[.[] | select(.labels[].name | startswith("status:"))]'
```
If any tickets have `status:*` labels, remove them before proceeding.

Move all `sprint:next` tickets to `sprint:current` and enqueue them as `status:planned`:
```bash
# For each sprint ticket
gh issue edit <number> -R rookiecj/scrum-agents \
  --remove-label "sprint:next" \
  --add-label "sprint:current,status:planned"
```

### 3. Display Sprint Board (Queue View)

Show all current sprint tickets organized by queue:
```bash
# DEV Queue
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:planned" --state open --json number,title,labels

# In Progress
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:in-progress" --state open --json number,title,labels

# QA Queue
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:dev-complete" --state open --json number,title,labels

# Verified
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:verified" --json number,title,labels
```

Present as a queue-based sprint board:
```
## Sprint Board

### DEV Queue (status:planned)
- #1 [Story] URL 타입 감지 및 웹 아티클 콘텐츠 추출 (5pts, Backend)
- #2 [Task] API 엔드포인트 설계 (3pts, Backend)

### In Progress (status:in-progress)
(none yet)

### QA Queue (status:dev-complete)
(none yet)

### Verified (status:verified)
(none yet)
```

### 4. Dispatch Agents

Ask the user which mode to use, or detect from arguments:

#### Sequential Mode (single agent)
Process tickets one by one, acting as both Dev and QA:

For each ticket in the DEV queue (priority order: critical → high → medium → low):

**DEV Phase:**
1. Claim: `status:planned` → `status:in-progress`
2. Read the full ticket: `gh issue view <number> -R rookiecj/scrum-agents`
3. Implement based on component label (see component-specific instructions below)
4. Complete: `status:in-progress` → `status:dev-complete`

**QA Phase:**
5. Claim: `status:dev-complete` → `status:in-review`
6. Verify each acceptance criterion
7. If pass: `status:in-review` → `status:verified`, close the issue
8. If fail: `status:in-review` → `status:planned` with failure comment, then re-process later

**Component-specific implementation:**

For `component:backend` tickets:
- Create feature branch: `git checkout -b feature/<number>-<short-description>`
- Implement in `backend/` following Go conventions from CLAUDE.md
- Write tests (table-driven tests)
- Run `cd backend && go build ./... && go test ./... -v`
- Commit with conventional commit: `feat: <description> (#<number>)`

For `component:frontend` tickets:
- Create feature branch: `git checkout -b feature/<number>-<short-description>`
- Implement in `frontend/` following TypeScript conventions from CLAUDE.md
- Write tests
- Run `cd frontend && npm run build && npm test`
- Commit with conventional commit: `feat: <description> (#<number>)`

For tickets with both components:
- Implement backend first, then frontend
- Same branch for both

#### Parallel Mode (multi-agent team)

Use the Task tool to spawn agents in two phases. Each agent runs in an **isolated git worktree** (`isolation: "worktree"`) to avoid working directory conflicts.

**Pre-flight: Handle dual-component tickets**
Before dispatching, check for tickets with both `component:backend` and `component:frontend`. These must be assigned to one agent explicitly (default: backend-dev handles the full ticket). Remove the other component label temporarily, or split the ticket.

**Phase 1 — Dev Agents (parallel)**
Read the agent definition files and include their full content in the Task tool prompt. Spawn Backend Dev and Frontend Dev simultaneously:

```
Task tool call #1:
  description: "Backend dev sprint work"
  subagent_type: "general-purpose"
  isolation: "worktree"
  prompt: |
    <paste full contents of .claude/agents/backend-dev.md>

    ## Sprint Context
    You are working on sprint tickets. Process ALL `component:backend` + `status:planned`
    tickets in priority order (critical → high → medium → low).
    For each ticket: claim → implement → mark dev-complete.
    Stop when the DEV queue for backend is empty.

Task tool call #2:
  description: "Frontend dev sprint work"
  subagent_type: "general-purpose"
  isolation: "worktree"
  prompt: |
    <paste full contents of .claude/agents/frontend-dev.md>

    ## Sprint Context
    You are working on sprint tickets. Process ALL `component:frontend` + `status:planned`
    tickets in priority order (critical → high → medium → low).
    For each ticket: claim → implement → mark dev-complete.
    Stop when the DEV queue for frontend is empty.
```

**Phase 2 — QA Agent (after dev agents complete or when QA queue has items)**
Wait for at least one Dev Agent to complete, or check that `status:dev-complete` tickets exist, then spawn QA:

```
Task tool call #3:
  description: "QA verification sprint work"
  subagent_type: "general-purpose"
  isolation: "worktree"
  prompt: |
    <paste full contents of .claude/agents/qa.md>

    ## Sprint Context
    You are verifying sprint tickets. Process ALL `status:dev-complete` tickets.
    For each ticket: claim → verify AC → pass or fail.
    Stop when the QA queue is empty AND no tickets are in `status:in-progress`
    (meaning no more dev work will produce new QA items).
```

**Error handling:**
- If a Dev Agent task fails mid-execution, its claimed tickets (`status:in-progress`) will be stale. The Scrum Master should detect these via Queue Health Check and transition them back to `status:planned`.
- If the QA Agent task fails, `status:in-review` tickets should be transitioned back to `status:dev-complete` by the Scrum Master.
- After all agents complete, run a final Queue Health Check to detect any orphaned tickets.

### 5. Handle Blockers

If a ticket is blocked during implementation:
```bash
# Mark as blocked (remove current status first)
gh issue edit <number> -R rookiecj/scrum-agents \
  --remove-label "status:planned" \
  --add-label "status:blocked"
gh issue comment <number> -R rookiecj/scrum-agents \
  --body "🚫 **Blocked**: <reason for block>. Previous state: status:planned."
```
- Skip to the next unblocked ticket

When the blocker is resolved, unblock the ticket:
```bash
gh issue edit <number> -R rookiecj/scrum-agents \
  --remove-label "status:blocked" \
  --add-label "status:planned"
gh issue comment <number> -R rookiecj/scrum-agents \
  --body "🔓 **Unblocked**: Blocker resolved. Returning to DEV queue."
```

### 6. Sprint Progress Updates

After each ticket completion, show queue-based progress:
```
## Sprint Progress
- DEV Queue:    X tickets
- In Progress:  X tickets
- QA Queue:     X tickets
- In Review:    X tickets
- Verified:     X/Y tickets (XX%)
- Blocked:      X tickets
- Points Done:  XX/XX pts
```

## Important
- **Status labels are mutually exclusive**: always remove the previous status label before adding the new one
- Always create feature branches — never commit directly to main
- Run tests before marking tickets as dev-complete
- If implementation requires design decisions not in the ticket, ask the user
- Commit messages must follow conventional commit format and reference the issue number
- If a ticket turns out to be larger than estimated, inform the user and discuss splitting
- QA rework tickets (returned to `status:planned` with failure comments) should be prioritized
