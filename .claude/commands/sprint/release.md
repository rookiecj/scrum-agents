You are the Scrum Master releasing the sprint deliverables. Bump versions, run quality checks, tag, and push.

## Input
Release arguments: $ARGUMENTS

## Instructions

### 1. Pre-flight Checks

Ensure the working tree is clean and on the main branch:
```bash
git status --porcelain
git branch --show-current
```

If there are uncommitted changes, warn the user and stop.
If not on `main`, warn and ask the user to confirm.

### 2. Read Current Versions

Read versions from both components:

**Backend** — `backend/cmd/server/version.go`:
```bash
grep 'const Version' backend/cmd/server/version.go
```
Extract the semver string (e.g., `"0.1.0"` → `0.1.0`).

**Frontend** — `frontend/package.json`:
```bash
node -e "console.log(require('./frontend/package.json').version)"
```

Display:
```
Current Versions:
  Backend:  X.Y.Z
  Frontend: X.Y.Z
```

### 3. Determine New Version

Compare the two versions and pick the **higher** one as the base.

Determine bump type from arguments:
- If arguments contain `major` → bump major (X+1.0.0)
- If arguments contain `minor` → bump minor (X.Y+1.0)
- If arguments contain `patch` or no bump type specified → bump patch (X.Y.Z+1)
- If arguments contain a specific version like `1.2.3` → use that exact version

Calculate the new version and display:
```
Version Bump:
  Base (higher of two): X.Y.Z
  Bump type: patch
  New version: X.Y.Z+1
```

Ask the user to confirm the new version before proceeding (skip if arguments contain `--yes` or `-y`).

### 4. Update Version Files

**Backend** — update `backend/cmd/server/version.go`:
Use the Edit tool to replace the version string:
```go
const Version = "<new-version>"
```

**Frontend** — update `frontend/package.json`:
```bash
cd frontend && npm version <new-version> --no-git-tag-version --allow-same-version
```

Verify both updates:
```bash
grep 'const Version' backend/cmd/server/version.go
node -e "console.log(require('./frontend/package.json').version)"
```

### 5. Run Quality Checks

Run all checks sequentially. If any step fails, stop and report the error.

**Backend**:
```bash
# Format check
cd backend && gofmt -l . | head -20

# Vet (lint)
cd backend && go vet ./...

# Build
cd backend && go build ./...

# Test
cd backend && go test ./... -v -cover
```

**Frontend**:
```bash
# Lint (type check)
cd frontend && npm run lint

# Build
cd frontend && npm run build

# Test
cd frontend && npm test
```

Display results as a checklist:
```
Quality Checks:
  Backend:
    ✅ gofmt — no formatting issues
    ✅ go vet — no lint issues
    ✅ go build — build successful
    ✅ go test — all tests pass (XX% coverage)
  Frontend:
    ✅ tsc --noEmit — no type errors
    ✅ vite build — build successful
    ✅ vitest — all tests pass
```

If any check fails, show `❌` with the error and stop. Do NOT proceed to tagging.

### 6. Commit Version Bump

```bash
git add backend/cmd/server/version.go frontend/package.json
git commit -m "$(cat <<'EOF'
chore: bump version to v<new-version>

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

### 7. Create Tag and Push

```bash
# Create annotated tag
git tag -a "v<new-version>" -m "Release v<new-version>"

# Push commit and tag
git push origin main
git push origin "v<new-version>"
```

### 8. Summary

Display the release summary:
```
## Release v<new-version>

### Versions Updated
- Backend:  X.Y.Z → <new-version>
- Frontend: X.Y.Z → <new-version>

### Quality Checks
- All passed ✅

### Git
- Commit: <short-hash>
- Tag: v<new-version>
- Pushed to: origin/main

### Next Steps
- Create GitHub Release from tag: gh release create v<new-version> --generate-notes
```

## Important
- Never tag or push if any quality check fails
- Both backend and frontend MUST have the same version after release
- Use annotated tags (`git tag -a`) not lightweight tags
- The version format is strict semver: `MAJOR.MINOR.PATCH`
- If arguments contain `--dry-run`, run all checks but skip commit/tag/push steps
