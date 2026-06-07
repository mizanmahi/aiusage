package repository

import (
	"context"
	"database/sql"

	"github.com/mizanmahi/aiusage/server/internal/domain"
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
