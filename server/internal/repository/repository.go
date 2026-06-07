package repository

import (
	"context"

	"github.com/mizanmahi/aiusage/server/internal/domain"
	"github.com/mizanmahi/aiusage/types"
)

type EventRepository interface {
	Upsert(ctx context.Context, userID string, events []types.UsageEvent) (accepted, skipped int, err error)
}

type UserRepository interface {
	FindByAPIKeyHash(ctx context.Context, apiKeyHash string) (*domain.User, error)
}
