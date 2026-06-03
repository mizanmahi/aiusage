# aiusage

Self-hosted analytics platform for AI coding tools. Collects Claude Code and Codex CLI usage from developer machines, stores it centrally, and surfaces per-user, per-project token and cost breakdowns in an admin dashboard.

## Why

Neither Claude Code nor Codex CLI expose per-project usage across a team. `ccusage` solves this for a single developer locally but has no push, no server, and no multi-user view. `aiusage` fills that gap.

---

## Architecture

```
Developer machines                     Your server
──────────────────────────────────     ──────────────────────────────────

~/.claude/projects/**/  ──┐
                           ├──▶  aiusage push  ──▶  POST /ingest
~/.codex/sessions/**/  ───┘      (Go CLI)            │
                                                      ▼
                                               aiusage-server
                                               (Go + Postgres)
                                                      │
                                                      ▼
                                               Admin dashboard
                                               (React, served by server)
```

### Data flow

1. Developer runs `aiusage push` (manually, cron, or git hook)
2. CLI reads local JSONL files from Claude Code and Codex
3. Extracts session ID, working directory, token counts, model, date
4. Derives project name from `basename(cwd)`
5. POSTs only new sessions (cursor-based, deduped by session ID) to the server
6. Server validates API key, upserts into Postgres
7. Admin views the dashboard — per user, per project, per day

---

## Repository layout

```
aiusage/
├── go.work              # Go workspace — links all three modules
├── Makefile             # build + dev commands
│
├── types/               # shared Go structs (UsageEvent, Project, User)
│   └── event.go
│
├── cli/                 # aiusage binary — runs on developer machines
│   ├── go.mod
│   ├── main.go
│   └── internal/
│       ├── claude/      # reads ~/.claude/projects/**/*.jsonl
│       ├── codex/       # reads ~/.codex/sessions/**/*.jsonl, parses cwd
│       ├── push/        # HTTP POST to server
│       └── state/       # tracks last_pushed cursor
│
├── server/              # central server
│   ├── go.mod
│   ├── main.go
│   └── internal/
│       ├── api/         # /ingest, /admin/*, /me/usage
│       ├── db/          # sqlc-generated Postgres queries
│       └── auth/        # API key middleware
│
└── ui/                  # admin dashboard (Vite + React)
    └── src/
```

---

## Database schema

```sql
users            -- id, email, api_key_hash, role (admin|dev)
usage_events     -- session_id (unique), user_id, date, tool, project,
                 -- cwd, model, input/output/cache/reasoning tokens, cost_usd
projects         -- materialized: project, user_id, total_tokens, total_cost
daily_snapshots  -- user_id, date, project, tool, tokens, cost_usd
```

---

## CLI commands

```bash
aiusage init          # first-time setup — saves server URL + API key
aiusage push          # push new sessions since last run
aiusage push --dry-run  # preview what would be pushed
aiusage status        # show last push time and pending session count
```

Config lives at `~/.aiusage/config.toml`:

```toml
server_url = "https://your-server.com"
api_key    = "ak_xxxxx"
```

---

## API endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/ingest` | Receive batch of usage events from CLI |
| GET | `/admin/users` | All users with token + cost totals |
| GET | `/admin/users/:id` | Per-project breakdown for one user |
| GET | `/admin/projects` | All projects across all users |
| GET | `/admin/summary` | Org-wide daily chart data |
| GET | `/me/usage` | Developer's own usage (non-admin) |

---

## Tech stack

| Layer | Tech |
|-------|------|
| CLI | Go, Cobra, TOML config |
| Server | Go stdlib `net/http`, sqlc |
| Database | Postgres |
| Dashboard | React, Vite |
| Monorepo | Go workspaces (`go.work`) |
| Dev infra | Docker Compose (Postgres) |

---

## Getting started

### Server

```bash
# start Postgres
docker-compose up -d

# run migrations
make migrate

# start server
make dev-server
```

### CLI (developer machine)

```bash
# install
go install github.com/yourusername/aiusage/cli@latest

# configure
aiusage init

# first push
aiusage push
```

### Dashboard

```bash
cd ui && npm install && npm run dev
```

---

## Roadmap

- [x] Claude Code parser
- [x] Codex CLI parser (cwd-based project grouping)
- [x] Push CLI with cursor tracking
- [x] Ingest API + Postgres storage
- [x] Admin dashboard — per user, per project
- [ ] Token efficiency ratio — correlate token spend with git output (commits, LOC, PRs)
- [ ] Webhook / Slack alert on budget threshold
- [ ] Self-serve developer view