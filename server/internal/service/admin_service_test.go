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

type fakeAdminProjectRepo struct{}

func (r fakeAdminProjectRepo) ListByUser(ctx context.Context, userID string) ([]types.ProjectSummary, error) {
	return nil, nil
}

func (r fakeAdminProjectRepo) ListAll(ctx context.Context) ([]types.ProjectSummary, error) {
	return nil, nil
}

func (r fakeAdminProjectRepo) DailySummary(ctx context.Context, from, to string) ([]types.DailyPoint, error) {
	return nil, nil
}

func TestCreateUserGeneratesAPIKeyAndStoresHash(t *testing.T) {
	users := &fakeAdminUserRepo{}
	service := NewAdminService(users, fakeAdminProjectRepo{})

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
	service := NewAdminService(&fakeAdminUserRepo{}, fakeAdminProjectRepo{})

	if _, err := service.CreateUser(context.Background(), &domain.User{ID: "user-1"}, types.CreateUserRequest{
		Email: "dev@example.com",
		Name:  "Dev User",
	}); err == nil {
		t.Fatal("CreateUser() error = nil, want forbidden")
	}
}
