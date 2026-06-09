package service

import (
	"context"
	"testing"

	"github.com/mizanmahi/aiusage/server/internal/domain"
	"github.com/mizanmahi/aiusage/types"
)

type fakeAdminUserRepo struct {
	created domain.User
}

func (r *fakeAdminUserRepo) FindByAPIKeyHash(ctx context.Context, apiKeyHash string) (*domain.User, error) {
	return nil, nil
}

func (r *fakeAdminUserRepo) ListWithTotals(ctx context.Context) ([]types.UserSummary, error) {
	return nil, nil
}

func (r *fakeAdminUserRepo) Create(ctx context.Context, user domain.User) (*types.UserSummary, error) {
	r.created = user
	return &types.UserSummary{
		ID:      "user-1",
		Email:   user.Email,
		Name:    user.Name,
		IsAdmin: user.IsAdmin,
	}, nil
}

type fakeAdminProjectRepo struct {
	breakdownUserID string
	breakdownGroup  string
	breakdownFrom   string
	breakdownTo     string
	summaryUserID   string
	summaryProvider string
	summaryFrom     string
	summaryTo       string
}

func (r fakeAdminProjectRepo) ListByUser(ctx context.Context, userID string) ([]types.ProjectSummary, error) {
	return nil, nil
}

func (r fakeAdminProjectRepo) ListAll(ctx context.Context) ([]types.ProjectSummary, error) {
	return nil, nil
}

func (r fakeAdminProjectRepo) DailySummary(ctx context.Context, from, to string) ([]types.DailyPoint, error) {
	return nil, nil
}

func (r *fakeAdminProjectRepo) UserBreakdown(ctx context.Context, userID, groupBy, from, to string) ([]types.UsageBreakdownRow, error) {
	r.breakdownUserID = userID
	r.breakdownGroup = groupBy
	r.breakdownFrom = from
	r.breakdownTo = to
	return []types.UsageBreakdownRow{{Group: "2026-06-09", Agent: "all"}}, nil
}

func (r *fakeAdminProjectRepo) UserUsageSummary(ctx context.Context, userID, provider, from, to string) (*types.UsageSummaryStats, error) {
	r.summaryUserID = userID
	r.summaryProvider = provider
	r.summaryFrom = from
	r.summaryTo = to
	return &types.UsageSummaryStats{Provider: provider}, nil
}

func TestCreateUserGeneratesAPIKeyAndStoresHash(t *testing.T) {
	users := &fakeAdminUserRepo{}
	service := NewAdminService(users, &fakeAdminProjectRepo{})

	result, err := service.CreateUser(context.Background(), &domain.User{ID: "admin-1", IsAdmin: true}, types.CreateUserRequest{
		Email:   " dev@example.com ",
		Name:    " Dev User ",
		IsAdmin: false,
	})
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if result.APIKey == "" {
		t.Fatal("APIKey = empty, want generated key")
	}
	if users.created.APIKeyHash != hashAPIKey(result.APIKey) {
		t.Fatal("stored API key hash does not match returned API key")
	}
	if users.created.Email != "dev@example.com" || users.created.Name != "Dev User" {
		t.Fatalf("created user = %+v, want trimmed email and name", users.created)
	}
	if result.User.ID != "user-1" {
		t.Fatalf("created ID = %q, want user-1", result.User.ID)
	}
}

func TestCreateUserRequiresAdmin(t *testing.T) {
	service := NewAdminService(&fakeAdminUserRepo{}, &fakeAdminProjectRepo{})

	if _, err := service.CreateUser(context.Background(), &domain.User{ID: "user-1"}, types.CreateUserRequest{
		Email: "dev@example.com",
		Name:  "Dev User",
	}); err == nil {
		t.Fatal("CreateUser() error = nil, want forbidden")
	}
}

func TestUserBreakdownDefaultsToDayAndWideDateRange(t *testing.T) {
	projects := &fakeAdminProjectRepo{}
	service := NewAdminService(&fakeAdminUserRepo{}, projects)

	rows, err := service.UserBreakdown(context.Background(), &domain.User{ID: "admin-1", IsAdmin: true}, "user-1", "", "", "")
	if err != nil {
		t.Fatalf("UserBreakdown() error = %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows len = %d, want 1", len(rows))
	}
	if projects.breakdownUserID != "user-1" || projects.breakdownGroup != "day" {
		t.Fatalf("breakdown args user=%q group=%q, want user-1/day", projects.breakdownUserID, projects.breakdownGroup)
	}
	if projects.breakdownFrom != "2000-01-01" || projects.breakdownTo != "2099-12-31" {
		t.Fatalf("date range = %s..%s, want defaults", projects.breakdownFrom, projects.breakdownTo)
	}
}

func TestUserBreakdownRejectsInvalidGroup(t *testing.T) {
	service := NewAdminService(&fakeAdminUserRepo{}, &fakeAdminProjectRepo{})

	if _, err := service.UserBreakdown(context.Background(), &domain.User{ID: "admin-1", IsAdmin: true}, "user-1", "week", "", ""); err == nil {
		t.Fatal("UserBreakdown() error = nil, want invalid group error")
	}
}

func TestUserUsageSummaryDefaultsToAllProvider(t *testing.T) {
	projects := &fakeAdminProjectRepo{}
	service := NewAdminService(&fakeAdminUserRepo{}, projects)

	stats, err := service.UserUsageSummary(context.Background(), &domain.User{ID: "admin-1", IsAdmin: true}, "user-1", "", "", "")
	if err != nil {
		t.Fatalf("UserUsageSummary() error = %v", err)
	}
	if stats.Provider != "all" {
		t.Fatalf("Provider = %q, want all", stats.Provider)
	}
	if projects.summaryUserID != "user-1" || projects.summaryProvider != "all" {
		t.Fatalf("summary args user=%q provider=%q, want user-1/all", projects.summaryUserID, projects.summaryProvider)
	}
}

func TestUserUsageSummaryRejectsInvalidProvider(t *testing.T) {
	service := NewAdminService(&fakeAdminUserRepo{}, &fakeAdminProjectRepo{})

	if _, err := service.UserUsageSummary(context.Background(), &domain.User{ID: "admin-1", IsAdmin: true}, "user-1", "openai", "", ""); err == nil {
		t.Fatal("UserUsageSummary() error = nil, want invalid provider error")
	}
}
