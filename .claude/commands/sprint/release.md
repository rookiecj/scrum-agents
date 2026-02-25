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

Read versions from the VERSION files:

**Backend** — `backend/VERSION`:
```bash
cat backend/VERSION
```

**Frontend** — `frontend/VERSION`:
```bash
cat frontend/VERSION
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

### 4. Update VERSION Files

**Backend** — write to `backend/VERSION`:
Use the Write tool to write the new version (with trailing newline) to `backend/VERSION`.

**Frontend** — write to `frontend/VERSION`:
Use the Write tool to write the new version (with trailing newline) to `frontend/VERSION`.

**Frontend package.json** — keep in sync:
```bash
cd frontend && npm version <new-version> --no-git-tag-version --allow-same-version
```

Verify all updates:
```bash
cat backend/VERSION
cat frontend/VERSION
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

# Build (with version injection from VERSION file)
cd backend && go build -ldflags "-X main.Version=$(cat VERSION)" ./cmd/server

# Test
cd backend && go test ./... -v -cover
```

**Frontend**:
```bash
# Lint (type check)
cd frontend && npm run lint

# Build (VERSION file is read by vite.config.ts at build time)
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
git add backend/VERSION frontend/VERSION frontend/package.json
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
- Backend:  X.Y.Z → <new-version> (backend/VERSION)
- Frontend: X.Y.Z → <new-version> (frontend/VERSION)

### Quality Checks
- All passed ✅

### Git
- Commit: <short-hash>
- Tag: v<new-version>
- Pushed to: origin/main

### Next Steps
- Create GitHub Release from tag: gh release create v<new-version> --generate-notes
```

## Version File Locations

| Component | Source of Truth | Consumed By |
|-----------|----------------|-------------|
| Backend | `backend/VERSION` | `go build -ldflags "-X main.Version=$(cat VERSION)"` injects into binary |
| Frontend | `frontend/VERSION` | `vite.config.ts` reads file and defines `__APP_VERSION__` at build time |
| Frontend | `frontend/package.json` | Kept in sync for npm ecosystem compatibility |

## Important
- Never tag or push if any quality check fails
- Both `backend/VERSION` and `frontend/VERSION` MUST have the same version after release
- `frontend/package.json` version is also synced but VERSION file is the source of truth
- Use annotated tags (`git tag -a`) not lightweight tags
- The version format is strict semver: `MAJOR.MINOR.PATCH`
- VERSION files contain only the version string with a trailing newline (e.g., `0.2.0\n`)
- If arguments contain `--dry-run`, run all checks but skip commit/tag/push steps
