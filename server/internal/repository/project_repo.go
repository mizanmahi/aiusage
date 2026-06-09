package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
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
			SUM(input_tokens + output_tokens + cache_creation_tokens + cache_read_tokens + reasoning_tokens)::bigint,
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

func (r *ProjectRepo) UserBreakdown(ctx context.Context, userID, groupBy, from, to string) ([]types.UsageBreakdownRow, error) {
	groupExpr, err := breakdownGroupExpr(groupBy)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM (
			SELECT
				%s AS usage_group,
				'all' AS agent,
				array_remove(array_agg(DISTINCT model ORDER BY model), '') AS models,
				SUM(input_tokens)::bigint,
				SUM(output_tokens)::bigint,
				SUM(cache_creation_tokens)::bigint,
				SUM(cache_read_tokens)::bigint,
				SUM(reasoning_tokens)::bigint,
				SUM(input_tokens + output_tokens + cache_creation_tokens + cache_read_tokens + reasoning_tokens)::bigint,
				SUM(cost_usd)::float8,
				MAX(date)::text
			FROM usage_events
			WHERE user_id = $1 AND date >= $2::date AND date <= $3::date
			GROUP BY usage_group
			UNION ALL
			SELECT
				%s AS usage_group,
				tool AS agent,
				array_remove(array_agg(DISTINCT model ORDER BY model), '') AS models,
				SUM(input_tokens)::bigint,
				SUM(output_tokens)::bigint,
				SUM(cache_creation_tokens)::bigint,
				SUM(cache_read_tokens)::bigint,
				SUM(reasoning_tokens)::bigint,
				SUM(input_tokens + output_tokens + cache_creation_tokens + cache_read_tokens + reasoning_tokens)::bigint,
				SUM(cost_usd)::float8,
				MAX(date)::text
			FROM usage_events
			WHERE user_id = $1 AND date >= $2::date AND date <= $3::date
			GROUP BY usage_group, tool
		) rows
		ORDER BY usage_group ASC, CASE WHEN agent = 'all' THEN 0 ELSE 1 END, agent ASC
	`, groupExpr, groupExpr)

	rows, err := r.db.QueryContext(ctx, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []types.UsageBreakdownRow
	for rows.Next() {
		var row types.UsageBreakdownRow
		if err := rows.Scan(
			&row.Group,
			&row.Agent,
			pq.Array(&row.Models),
			&row.InputTokens,
			&row.OutputTokens,
			&row.CacheCreateTokens,
			&row.CacheReadTokens,
			&row.ReasoningTokens,
			&row.TotalTokens,
			&row.TotalCost,
			&row.LastActive,
		); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ProjectRepo) UserUsageSummary(ctx context.Context, userID, provider, from, to string) (*types.UsageSummaryStats, error) {
	providerFilter := ""
	args := []any{userID, from, to}
	if provider != "all" {
		providerFilter = " AND tool = $4"
		args = append(args, provider)
	}

	query := `
		SELECT
			COUNT(DISTINCT project)::bigint,
			COALESCE(SUM(input_tokens), 0)::bigint,
			COALESCE(SUM(output_tokens), 0)::bigint,
			COALESCE(SUM(cache_creation_tokens + cache_read_tokens), 0)::bigint,
			COALESCE(SUM(input_tokens + output_tokens + cache_creation_tokens + cache_read_tokens + reasoning_tokens), 0)::bigint,
			COALESCE(SUM(cost_usd), 0)::float8
		FROM usage_events
		WHERE user_id = $1 AND date >= $2::date AND date <= $3::date` + providerFilter

	var stats types.UsageSummaryStats
	stats.Provider = provider
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&stats.TotalProjects,
		&stats.TotalInput,
		&stats.TotalOutput,
		&stats.TotalCached,
		&stats.TotalTokens,
		&stats.TotalCost,
	); err != nil {
		return nil, err
	}

	return &stats, nil
}

func breakdownGroupExpr(groupBy string) (string, error) {
	switch groupBy {
	case "day":
		return "date::text", nil
	case "month":
		return "to_char(date, 'YYYY-MM')", nil
	case "project":
		return "project", nil
	default:
		return "", fmt.Errorf("unsupported breakdown group: %s", groupBy)
	}
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
