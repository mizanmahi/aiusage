package repository

import (
	"context"
	"database/sql"

	"github.com/mizanmahi/aiusage/server/internal/domain"
	"github.com/mizanmahi/aiusage/types"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByAPIKeyHash(ctx context.Context, apiKeyHash string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id::text, email, name, api_key_hash, is_admin
		FROM users
		WHERE api_key_hash = $1
	`, apiKeyHash)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Email, &user.Name, &user.APIKeyHash, &user.IsAdmin); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) ListWithTotals(ctx context.Context) ([]types.UserSummary, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			u.id::text,
			u.email,
			u.name,
			COALESCE(SUM(e.input_tokens + e.output_tokens + e.cache_tokens + e.reasoning_tokens), 0)::bigint,
			COALESCE(SUM(e.cost_usd), 0)::float8,
			COALESCE(to_char(MAX(e.date), 'YYYY-MM-DD'), '')
		FROM users u
		LEFT JOIN usage_events e ON e.user_id = u.id
		GROUP BY u.id, u.email, u.name
		ORDER BY
			COALESCE(SUM(e.input_tokens + e.output_tokens + e.cache_tokens + e.reasoning_tokens), 0) DESC,
			u.email ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.UserSummary
	for rows.Next() {
		var user types.UserSummary
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.TotalTokens,
			&user.TotalCost,
			&user.LastActive,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
