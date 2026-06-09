package repository

import (
	"context"
	"database/sql"

	"github.com/mizanmahi/aiusage/types"
)

type EventRepo struct {
	db *sql.DB
}

func NewEventRepo(db *sql.DB) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) Upsert(ctx context.Context, userID string, events []types.UsageEvent) (int, int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	var accepted int
	for _, event := range events {
		projectID, err := upsertProject(ctx, tx, userID, event)
		if err != nil {
			return 0, 0, err
		}

		inserted, err := insertEvent(ctx, tx, userID, projectID, event)
		if err != nil {
			return 0, 0, err
		}
		if inserted {
			accepted++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}

	go refreshProjectTotals(context.Background(), r.db)
	return accepted, len(events) - accepted, nil
}

func upsertProject(ctx context.Context, tx *sql.Tx, userID string, event types.UsageEvent) (string, error) {
	row := tx.QueryRowContext(ctx, `
		INSERT INTO projects (user_id, name, tool, cwd)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, tool, name)
		DO UPDATE SET cwd = EXCLUDED.cwd, updated_at = now()
		RETURNING id::text
	`, userID, event.Project, event.Tool, event.Cwd)

	var projectID string
	if err := row.Scan(&projectID); err != nil {
		return "", err
	}
	return projectID, nil
}

func insertEvent(ctx context.Context, tx *sql.Tx, userID, projectID string, event types.UsageEvent) (bool, error) {
	result, err := tx.ExecContext(ctx, `
		INSERT INTO usage_events (
			session_id, user_id, project_id, date, tool, project, cwd, model,
			input_tokens, output_tokens, cache_creation_tokens, cache_read_tokens, reasoning_tokens, cost_usd, pushed_at
		)
		VALUES ($1, $2, $3, $4::date, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (session_id) DO NOTHING
	`, event.SessionID, userID, projectID, event.Date, event.Tool, event.Project, event.Cwd, event.Model,
		event.InputTokens, event.OutputTokens, event.CacheCreateTokens, event.CacheReadTokens, event.ReasoningTokens, event.CostUSD, event.PushedAt)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows == 1, nil
}

func refreshProjectTotals(ctx context.Context, db *sql.DB) {
	_, _ = db.ExecContext(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY project_totals")
}
