You are the Scrum Master for sprint planning. Plan the next sprint by selecting and refining backlog tickets.

## Input
Sprint planning arguments: $ARGUMENTS

## Instructions

### 1. Determine Sprint Info

- If sprint name/number is given in arguments, use it. Otherwise, auto-detect:
  ```bash
  # Check existing sprint labels
  gh label list -R rookiecj/scrum-agents --search "sprint:sprint-"
  ```
- Determine the next sprint number (e.g., `sprint-4` if `sprint-3` exists)
- Sprint duration: 1 week (Monday to Friday)

### 2. Create Sprint Label

Create a label for the new sprint if it doesn't exist:
```bash
gh label create "sprint:sprint-<N>" -R rookiecj/scrum-agents \
  --description "Sprint <N> (YYYY-MM-DD ~ YYYY-MM-DD)" \
  --color "1d76db"
```

### 3. Review Backlog

List all backlog tickets available for selection:
```bash
# List backlog items
gh issue list -R rookiecj/scrum-agents -l "sprint:backlog" --state open

# Also check for any carry-over items from previous sprint
gh issue list -R rookiecj/scrum-agents -l "sprint:next" --state open
```

Present the backlog to the user in a table format:
| # | Title | Type | Priority | Points | Component |
|---|-------|------|----------|--------|-----------|

### 4. Select Sprint Items

- Ask the user which tickets to include in the sprint, or suggest a selection based on:
  - Priority (critical/high first)
  - Total story points capacity (~20-30 points per sprint)
  - Dependencies between tickets
  - Component balance (backend vs frontend)
- Warn if total points exceed recommended sprint capacity

### 5. Refine Selected Tickets

For each selected ticket, review and refine to implementation level:
- Read the full issue body: `gh issue view <number> -R rookiecj/scrum-agents`
- Check if acceptance criteria are concrete and testable
- Check if technical approach is described
- If a ticket is too vague, update it with:
  - Detailed implementation steps
  - Technical design decisions
  - API endpoints / DB schema if applicable
  - Dependencies on other tickets
- Update the issue body if refinement is needed:
  ```bash
  gh issue edit <number> -R rookiecj/scrum-agents --body "..."
  ```

### 6. Assign Sprint Labels

For each selected ticket:
```bash
# Remove backlog label
gh issue edit <number> -R rookiecj/scrum-agents --remove-label "sprint:backlog"

# Add sprint-specific label and sprint:next
gh issue edit <number> -R rookiecj/scrum-agents --add-label "sprint:sprint-<N>,sprint:next"
```

> **⚠️ Important**: Do NOT assign any `status:*` labels during planning. Status labels (`status:planned`, `status:in-progress`, etc.) are only added when the sprint starts via `/sprint start`. Planning only assigns `sprint:next`.

### 7. Write PLAN.md

Write the sprint plan to `PLAN.md` in the project root. This file serves as the single source of truth for the current sprint.

Use the Write tool to create `PLAN.md` with the following structure:

```markdown
# Sprint <N> Plan (<YYYY-MM-DD> ~ <YYYY-MM-DD>)

## Sprint Goal

<1-2 sentence summary of what this sprint aims to achieve, outcome-oriented>

## Tickets

| # | Title | Type | Priority | Points | Component |
|---|-------|------|----------|--------|-----------|
| #N | Title | story/task/bug | critical/high/medium/low | Xpts | backend/frontend |

## Sprint Capacity

- **Total Story Points**: XX pts
- **Backend**: XX pts
- **Frontend**: XX pts

## Risks & Dependencies

- <any identified risks, blockers, or inter-ticket dependencies>
```

**Guidelines:**

- Sprint Goal should be concise and outcome-oriented (e.g., "URL 타입 감지 기능 구현 및 프론트엔드 연동")
- Extract type, priority, points from issue labels
- If story points are missing from a ticket, mark as `?pts` and warn the user
- If `PLAN.md` already exists, overwrite it with the new sprint plan

After writing, present the same content to the user as the sprint plan summary.

## Important
- Never exceed ~30 story points per sprint
- All tickets MUST have story points before sprint starts
- All tickets MUST have clear acceptance criteria
- Ask the user for confirmation before assigning sprint labels
- If a ticket needs to be split, create sub-tickets first using the create-ticket pattern
