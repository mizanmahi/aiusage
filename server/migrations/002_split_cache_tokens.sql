-- +goose Up
DROP MATERIALIZED VIEW IF EXISTS project_totals;

ALTER TABLE usage_events
    ADD COLUMN cache_creation_tokens BIGINT NOT NULL DEFAULT 0 CHECK (cache_creation_tokens >= 0),
    ADD COLUMN cache_read_tokens BIGINT NOT NULL DEFAULT 0 CHECK (cache_read_tokens >= 0);

ALTER TABLE usage_events
    DROP COLUMN cache_tokens;

CREATE MATERIALIZED VIEW project_totals AS
SELECT
    user_id,
    project,
    tool,
    SUM(input_tokens + output_tokens + cache_creation_tokens + cache_read_tokens + reasoning_tokens)::BIGINT AS total_tokens,
    SUM(cost_usd)::NUMERIC(12, 6) AS total_cost_usd,
    MAX(date)::DATE AS last_active
FROM usage_events
GROUP BY user_id, project, tool;

CREATE UNIQUE INDEX idx_project_totals_unique ON project_totals(user_id, project, tool);
CREATE INDEX idx_project_totals_last_active ON project_totals(last_active DESC);

-- +goose Down
DROP MATERIALIZED VIEW IF EXISTS project_totals;

ALTER TABLE usage_events
    ADD COLUMN cache_tokens BIGINT NOT NULL DEFAULT 0 CHECK (cache_tokens >= 0);

ALTER TABLE usage_events
    DROP COLUMN cache_creation_tokens,
    DROP COLUMN cache_read_tokens;

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
