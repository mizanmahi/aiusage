package repository

import (
	"context"
	"errors"

	"github.com/mizanmahi/aiusage/server/internal/domain"
	"github.com/mizanmahi/aiusage/types"
)

var ErrUserExists = errors.New("user already exists")

type EventRepository interface {
	Upsert(ctx context.Context, userID string, events []types.UsageEvent) (accepted, skipped int, err error)
}

type UserRepository interface {
	FindByAPIKeyHash(ctx context.Context, apiKeyHash string) (*domain.User, error)
	ListWithTotals(ctx context.Context) ([]types.UserSummary, error)
	Create(ctx context.Context, user domain.User) (*types.UserSummary, error)
}

type ProjectRepository interface {
	ListByUser(ctx context.Context, userID string) ([]types.ProjectSummary, error)
	ListAll(ctx context.Context) ([]types.ProjectSummary, error)
	DailySummary(ctx context.Context, from, to string) ([]types.DailyPoint, error)
	UserBreakdown(ctx context.Context, userID, groupBy, provider, from, to string) ([]types.UsageBreakdownRow, error)
	UserUsageSummary(ctx context.Context, userID, provider, from, to string) (*types.UsageSummaryStats, error)
}
