You are the Scrum Master checking the current sprint status. Query all queue stages and present a comprehensive progress report.

## Input
Sprint status arguments: $ARGUMENTS

## Instructions

### 1. Check for Active Sprint

```bash
gh issue list -R rookiecj/scrum-agents -l "sprint:current" --state all --json number,title,labels,state
```

If no `sprint:current` tickets exist, inform the user there is no active sprint. Suggest running `/sprint plan` then `/sprint start`.

### 2. Query All Queue Stages

Run all queries to gather the full picture:
```bash
# DEV Queue — waiting for development
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:planned" --state open --json number,title,labels

# In Progress — developer actively working
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:in-progress" --state open --json number,title,labels

# QA Queue — awaiting verification
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:dev-complete" --state open --json number,title,labels

# In Review — QA actively verifying
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:in-review" --state open --json number,title,labels

# Verified — QA passed (may still be open)
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:verified" --json number,title,labels

# Blocked
gh issue list -R rookiecj/scrum-agents -l "sprint:current" -l "status:blocked" --state open --json number,title,labels

# Done — closed issues
gh issue list -R rookiecj/scrum-agents -l "sprint:current" --state closed --json number,title,labels
```

### 3. Display Sprint Board

Present the queue-based sprint board:

```
## Sprint Board

### DEV Queue (status:planned) — X tickets
- #<number> [<type>] <title> (<points>pts, <component>)

### In Progress (status:in-progress) — X tickets
- #<number> [<type>] <title> (<points>pts, <component>)

### QA Queue (status:dev-complete) — X tickets
- #<number> [<type>] <title> (<points>pts, <component>)

### In Review (status:in-review) — X tickets
- #<number> [<type>] <title> (<points>pts, <component>)

### Done (verified/closed) — X tickets
- #<number> [<type>] <title> (<points>pts, <component>)

### Blocked (status:blocked) — X tickets
- #<number> [<type>] <title> (<points>pts, <component>) — Blocker: <reason>
```

### 4. Show Progress Metrics

Calculate and display key metrics:

```
## Sprint Progress

### Completion
- Total:      XX tickets (XX pts)
- Done:       XX tickets (XX pts) — XX%
- Remaining:  XX tickets (XX pts)

### Pipeline
| Stage         | Count | Points |
|---------------|-------|--------|
| DEV Queue     | X     | XX     |
| In Progress   | X     | XX     |
| QA Queue      | X     | XX     |
| In Review     | X     | XX     |
| Done          | X     | XX     |
| Blocked       | X     | XX     |

### Progress Bar
[████████░░░░░░░░░░░░] 40% (4/10 tickets)
```

### 5. Health Check & Alerts

Flag issues that need attention:

- **Blocked tickets**: List each with the blocker reason (read the latest `🚫 **Blocked**` comment)
- **Stale in-progress**: Tickets in `status:in-progress` that may be abandoned (check for recent comments)
- **QA bottleneck**: If QA Queue has more tickets than DEV Queue, flag it
- **DEV bottleneck**: If DEV Queue has more tickets than QA Queue and QA is idle, flag it
- **Rework tickets**: Check for tickets returned from QA (have `❌ **QA Failed**` comments) — these need priority attention

Present alerts as:
```
## Alerts

⚠️ **Blocked**: #5 — waiting on external API access
⚠️ **QA Bottleneck**: 3 tickets in QA queue, 0 in DEV queue — consider dispatching QA agent
⚠️ **Rework**: #7 returned from QA (rework 1/3) — needs priority fix
```

If no alerts, show: `✅ No issues detected. Sprint is on track.`

## Important
- This command is **read-only** — it does not modify any labels or tickets
- Extract story points from issue labels or body (look for patterns like `(Xpts)` or story point labels)
- For blocked tickets, read the most recent blocker comment to show the reason
- If arguments contain "brief" or "short", show only the Progress Metrics section (skip the full board)
