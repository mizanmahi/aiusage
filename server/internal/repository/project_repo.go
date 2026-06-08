package repository

import (
	"context"
	"database/sql"

	"github.com/mizanmahi/aiusage/types"
)

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) ListByUser(ctx context.Context, userID string) ([]types.ProjectSummary, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT project, user_id::text, tool, total_tokens, total_cost_usd::float8, last_active::text
		FROM project_totals
		WHERE user_id = $1
		ORDER BY total_tokens DESC, project ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProjectSummaries(rows)
}

func (r *ProjectRepo) ListAll(ctx context.Context) ([]types.ProjectSummary, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT project, user_id::text, tool, total_tokens, total_cost_usd::float8, last_active::text
		FROM project_totals
		ORDER BY total_tokens DESC, project ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProjectSummaries(rows)
}

func (r *ProjectRepo) DailySummary(ctx context.Context, from, to string) ([]types.DailyPoint, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			date::text,
			SUM(input_tokens + output_tokens + cache_tokens + reasoning_tokens)::bigint,
			SUM(cost_usd)::float8
		FROM usage_events
		WHERE date >= $1::date AND date <= $2::date
		GROUP BY date
		ORDER BY date ASC
	`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []types.DailyPoint
	for rows.Next() {
		var point types.DailyPoint
		if err := rows.Scan(&point.Date, &point.TotalTokens, &point.TotalCost); err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return points, nil
}

func scanProjectSummaries(rows *sql.Rows) ([]types.ProjectSummary, error) {
	var projects []types.ProjectSummary
	for rows.Next() {
		var project types.ProjectSummary
		if err := rows.Scan(
			&project.Project,
			&project.UserID,
			&project.Tool,
			&project.TotalTokens,
			&project.TotalCost,
			&project.LastActive,
		); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
