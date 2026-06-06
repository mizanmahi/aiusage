
## Project Overview

**aiusage** is an internal developer analytics platform that tracks AI coding tool usage
across a team. Developers run a lightweight CLI (`aiusage`) on their machines that reads
local session files from Claude Code (`~/.claude/projects/`) and Codex CLI
(`~/.codex/sessions/`) and pushes usage data to a central server. An admin dashboard
gives the engineering manager a real-time view of who is using which AI tool, on which
project, and how many tokens they are burning.

### Why this exists

AI coding tools (Claude Code, Codex) store session data locally. There is no built-in
way for a team to see cross-developer, cross-project usage. This tool fills that gap
without relying on third-party analytics services or requiring developers to change
their workflow — they just run one command (or set a cron job) and everything is tracked
automatically.

### What the system does

- **CLI (Go + Cobra)** — runs on each developer's machine. Reads JSONL session files,
  extracts project name (from `cwd`), token counts, model, and date. Pushes only new
  sessions since the last push (cursor-based). Supports `push`, `push --dry-run`,
  `status`, and `init` commands.

- **Server (Go + Chi)** — central HTTP API. Authenticates developers via API key.
  Accepts batch session payloads, deduplicates by session ID, stores in Postgres.
  Exposes admin endpoints for the dashboard and a self-serve `/me` endpoint for
  individual developers.

- **Dashboard (React)** — admin-only web UI served as embedded static files from the
  server binary. Shows per-user usage analytics.

### Key data points tracked per session

| Field | Source |
|---|---|
| Developer (user) | API key → user lookup |
| Tool | `claude` or `codex` |
| Project | `basename(cwd)` from session file |
| Date | session timestamp or directory path |
| Model | from session JSONL (`gpt-5.5`, `claude-sonnet-4-6` etc.) |
| Input tokens | accumulated across all turns in the session |
| Output tokens | accumulated across all turns |
| Cache tokens | cache read + cache creation tokens |
| Reasoning tokens | where available (Codex) |
| Cost USD | calculated server-side from token counts + model pricing |

---

## Architecture

### Repository layout

```
aiusage/                        ← monorepo root
├── go.work                     ← Go workspace linking all modules
├── Makefile
├── version.txt                 ← single source of version truth
├── AGENTS.md
├── BUILDING.md
│
├── types/                      ← shared Go structs (imported by cli + server)
│   ├── go.mod
│   ├── event.go                ← UsageEvent, PushPayload, PushResponse
│   └── admin.go                ← UserSummary, ProjectSummary, DailyPoint
│
├── cli/                        ← aiusage binary
│   ├── go.mod
│   ├── main.go                 ← cobra root, version flag
│   ├── cmd/
│   │   ├── root.go
│   │   ├── init.go             ← aiusage init
│   │   ├── push.go             ← aiusage push [--dry-run]
│   │   └── status.go           ← aiusage status
│   └── internal/
│       ├── claude/
│       │   ├── reader.go       ← parses ~/.claude/projects/**/*.jsonl
│       │   └── reader_test.go
│       ├── codex/
│       │   ├── reader.go       ← parses ~/.codex/sessions/**/*.jsonl, extracts cwd
│       │   └── reader_test.go
│       ├── push/
│       │   ├── client.go       ← HTTP POST /ingest, compat header check
│       │   └── client_test.go
│       ├── config/
│       │   └── config.go       ← loads ~/.aiusage/config.toml
│       └── state/
│           └── state.go        ← last_pushed cursor in ~/.aiusage/state.json
│
├── server/                     ← central API server
│   ├── go.mod
│   ├── main.go                 ← wires everything, starts http.Server
│   ├── migrations/
│   │   └── 001_initial.sql
│   └── internal/
│       ├── config/
│       │   └── config.go       ← loads from env vars
│       ├── domain/             ← pure business types (no DB, no HTTP)
│       │   ├── event.go
│       │   ├── user.go
│       │   └── project.go
│       ├── repository/         ← data access layer (repository pattern)
│       │   ├── repository.go   ← interfaces
│       │   ├── event_repo.go   ← EventRepository implementation
│       │   ├── user_repo.go    ← UserRepository implementation
│       │   └── project_repo.go ← ProjectRepository implementation
│       ├── service/            ← business logic (calls repository, never HTTP)
│       │   ├── ingest_service.go
│       │   ├── user_service.go
│       │   └── project_service.go
│       ├── handler/            ← thin HTTP handlers (calls service, never DB)
│       │   ├── ingest.go
│       │   ├── admin.go
│       │   └── me.go
│       ├── middleware/
│       │   ├── auth.go         ← API key validation, sets user in context
│       │   ├── logger.go       ← structured request logging
│       │   └── recovery.go     ← panic recovery
│       └── apperror/
│           └── errors.go       ← typed application errors → consistent HTTP codes
│
└── ui/                         ← React dashboard (built output embedded in server)
    ├── package.json
    └── src/
```

### Request lifecycle (server)

```
HTTP request
    ↓
Chi router
    ↓
Middleware chain: Recovery → Logger → Auth
    ↓
Handler         ← thin: parse input, call service, write response
    ↓
Service         ← business logic: validation, orchestration
    ↓
Repository      ← data access: SQL only, returns domain types
    ↓
Postgres
```

Handlers never touch the DB. Services never touch `http.Request`. Repositories never
contain business logic. This separation makes every layer independently testable.

### Repository pattern

Every data access is behind an interface:

```go
// server/internal/repository/repository.go
type EventRepository interface {
    Upsert(ctx context.Context, events []domain.Event) (accepted, skipped int, err error)
    ListByUser(ctx context.Context, userID string, from, to time.Time) ([]domain.Event, error)
}

type UserRepository interface {
    FindByAPIKey(ctx context.Context, apiKey string) (*domain.User, error)
    ListWithTotals(ctx context.Context) ([]domain.UserSummary, error)
}

type ProjectRepository interface {
    ListByUser(ctx context.Context, userID string) ([]domain.ProjectSummary, error)
    ListAll(ctx context.Context) ([]domain.ProjectSummary, error)
}
```

In tests, swap the real Postgres implementation for an in-memory mock. No test database
needed for unit tests.

### Error handling

Use typed application errors that map cleanly to HTTP status codes:

```go
// server/internal/apperror/errors.go
type AppError struct {
    Code    string // machine-readable e.g. "UNAUTHORIZED", "INVALID_PAYLOAD"
    Message string // human-readable
    Status  int    // HTTP status code
    Err     error  // underlying error for logging (never sent to client)
}

// standard constructors
apperror.Unauthorized("invalid API key")         // 401
apperror.BadRequest("session_id is required")    // 400
apperror.NotFound("user not found")              // 404
apperror.Internal("db query failed", err)        // 500
```

All handlers use a single `writeError(w, err)` helper that:
- Checks if err is `*AppError` → writes `{"error": {"code": ..., "message": ...}}`
- Logs the underlying `Err` at ERROR level (not sent to client)
- Falls back to 500 for unknown errors

### Logging

Use `log/slog` (stdlib, Go 1.21+). No external logging library needed.

```go
slog.Info("session ingested",
    "user_id", userID,
    "accepted", accepted,
    "skipped", skipped,
    "tool", tool,
)

slog.Error("db upsert failed",
    "error", err,
    "user_id", userID,
)
```

JSON format in production (`slog.NewJSONHandler`), text format in development.
Never log API keys, tokens, or any secret values.

### Configuration

Server reads all config from environment variables — no config files on the server.

```go
// server/internal/config/config.go
type Config struct {
    DatabaseURL string // DATABASE_URL
    Port        string // PORT, default "8080"
    Env         string // ENV, "development" or "production"
    MinCLIVersion string // MIN_CLI_VERSION, for compat header
}
```

CLI reads from `~/.aiusage/config.toml` (written by `aiusage init`).

---

## Project Rules

- Inspect related files before changing any code.
- Keep changes minimal and focused on the task.
- Do not refactor unrelated code in the same PR.
- Follow existing patterns and naming in the file you are editing.
- Do not add dependencies unless clearly necessary and justified.
- Never hardcode secrets, API keys, or credentials anywhere.
- Preserve existing API response shapes unless explicitly asked to change them.
- A file must not exceed 250 lines (300 for UI files). If it does, split it.

---

## Code Style

### General
- Prefer readable code over clever one-liners.
- Use early returns to reduce nesting.
- Keep functions small and focused on one thing.
- Reuse existing utilities, services, and helpers.
- Avoid duplicate logic — extract to a shared helper if used twice.
- Comment non-obvious decisions, not what the code obviously does.

### Go specific
- Return errors, never panic in library code.
- Use `context.Context` as the first argument in every function that touches I/O.
- Name interfaces by behavior: `EventRepository`, not `IEventRepository`.
- Keep struct fields exported only when necessary.
- Use table-driven tests with `t.Run` for subtests.
- Use `t.TempDir()` for any test that writes files — never write to real home dir.
- Mock at the repository layer, not the database level.

---

## Backend / API

- Validate all inputs at the handler level before passing to service.
- Keep response shape consistent across all endpoints:
  ```json
  // success
  { "data": { ... } }

  // error
  { "error": { "code": "UNAUTHORIZED", "message": "invalid API key" } }
  ```
- Handle errors using `apperror` — never write raw `http.Error` strings in handlers.
- Keep handlers thin: parse → call service → write response. No business logic.
- Put business logic in services. Put SQL in repositories.
- Use `chi.URLParam(r, "id")` for path params, never split `r.URL.Path` manually.
- Set `Content-Type: application/json` on every JSON response.
- Use structured logging (`slog`) — never `fmt.Println` or `log.Println` in server code.

---

## CLI

- Every command must have a `--dry-run` equivalent where relevant.
- Print human-readable output to stdout, errors to stderr.
- Exit code 0 = success, non-zero = failure (so scripts can chain with `&&`).
- Config and state files go under `~/.aiusage/` with mode `0600`.
- Never print the API key in any output, including `--verbose` mode.
- Keep command handlers thin: parse flags → call internal package → print result.

---

## Database

- Avoid N+1 queries — fetch related data in one query where possible.
- Use `ON CONFLICT (session_id) DO NOTHING` for idempotent session inserts.
- Wrap multi-step writes in a transaction.
- Never run `DROP` or `TRUNCATE` without an explicit migration file.
- Refresh the `project_totals` materialized view asynchronously after ingest
  (`go db.Exec("REFRESH MATERIALIZED VIEW CONCURRENTLY project_totals")`).
- Think about race conditions on concurrent pushes from the same user.
- Index columns used in `WHERE` and `ORDER BY` clauses.

---

## Testing

### Unit tests
- Test each parser (claude, codex) with fixture JSONL files checked into the repo
  under `cli/internal/testdata/`.
- Test services with mock repositories — no real DB, no real HTTP.
- Test handlers with `httptest.NewRecorder` — no real server needed.

### Integration tests
- Use a real Postgres instance spun up via `docker-compose` or `testcontainers-go`.
- Tag integration tests with `//go:build integration` so they don't run on unit test pass.

### Running tests
```bash
# unit tests only (fast, no DB)
make test

# integration tests (requires docker)
make test-integration

# single package
cd cli && go test ./internal/codex/... -v

# with race detector
cd server && go test -race ./...
```

---

## Verification

Before finishing any task, run:

```bash
make test       # all unit tests
make lint       # golangci-lint (if configured)
cd ui && npm run typecheck  # frontend type check
```

Then summarize:

- **Files changed** — list each file
- **What changed** — one line per file
- **Verification done** — which checks passed
- **Risks or side effects** — anything that could affect other parts of the system

---

## Versioning and Release

- Single semver version for the whole repo (`version.txt` at root).
- Commit messages must follow Conventional Commits:
  - `feat:` → minor bump
  - `fix:` → patch bump
  - `feat!:` or `BREAKING CHANGE:` → major bump
  - `chore:`, `docs:`, `refactor:` → no bump
- Release Please automates the changelog and tag on merge to `main`.
- Build workflow triggers on `v*` tags and uploads CLI + server binaries to GitHub Releases.
- CLI and server are versioned together. Server sets `X-Aiusage-Min-CLI-Version` header
  to enforce minimum CLI version when the ingest contract changes.