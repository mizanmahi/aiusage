-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    api_key_hash TEXT NOT NULL UNIQUE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    tool TEXT NOT NULL CHECK (tool IN ('claude', 'codex')),
    cwd TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, tool, name)
);

CREATE TABLE usage_events (
    session_id TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
    date DATE NOT NULL,
    tool TEXT NOT NULL CHECK (tool IN ('claude', 'codex')),
    project TEXT NOT NULL,
    cwd TEXT NOT NULL DEFAULT '',
    model TEXT NOT NULL DEFAULT '',
    input_tokens BIGINT NOT NULL DEFAULT 0 CHECK (input_tokens >= 0),
    output_tokens BIGINT NOT NULL DEFAULT 0 CHECK (output_tokens >= 0),
    cache_tokens BIGINT NOT NULL DEFAULT 0 CHECK (cache_tokens >= 0),
    reasoning_tokens BIGINT NOT NULL DEFAULT 0 CHECK (reasoning_tokens >= 0),
    cost_usd NUMERIC(12, 6) NOT NULL DEFAULT 0 CHECK (cost_usd >= 0),
    pushed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_usage_events_user_date ON usage_events(user_id, date DESC);
CREATE INDEX idx_usage_events_project_date ON usage_events(project, date DESC);
CREATE INDEX idx_usage_events_tool_date ON usage_events(tool, date DESC);
CREATE INDEX idx_usage_events_created_at ON usage_events(created_at DESC);
CREATE INDEX idx_projects_user_name ON projects(user_id, name);

CREATE MATERIALIZED VIEW project_totals AS
SELECT
    user_id,
    project,
    tool,
    SUM(input_tokens + output_tokens + cache_tokens + reasoning_tokens)::BIGINT AS total_tokens,
    SUM(cost_usd)::NUMERIC(12, 6) AS total_cost_usd,
    MAX(date)::DATE AS last_active
FROM usage_events
GROUP BY user_id, project, tool;

CREATE UNIQUE INDEX idx_project_totals_unique ON project_totals(user_id, project, tool);
CREATE INDEX idx_project_totals_last_active ON project_totals(last_active DESC);

-- +goose Down
DROP MATERIALIZED VIEW IF EXISTS project_totals;
DROP TABLE IF EXISTS usage_events;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;
