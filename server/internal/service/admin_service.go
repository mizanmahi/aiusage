package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/mizanmahi/aiusage/server/internal/apperror"
	"github.com/mizanmahi/aiusage/server/internal/domain"
	"github.com/mizanmahi/aiusage/server/internal/repository"
	"github.com/mizanmahi/aiusage/types"
)

type AdminService struct {
	users    repository.UserRepository
	projects repository.ProjectRepository
}

func NewAdminService(users repository.UserRepository, projects repository.ProjectRepository) *AdminService {
	return &AdminService{users: users, projects: projects}
}

func (s *AdminService) Users(ctx context.Context, actor *domain.User) ([]types.UserSummary, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}

	users, err := s.users.ListWithTotals(ctx)
	if err != nil {
		return nil, apperror.Internal("failed to list users", err)
	}
	return users, nil
}

func (s *AdminService) CreateUser(ctx context.Context, actor *domain.User, request types.CreateUserRequest) (*types.CreateUserResponse, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}

	email := strings.TrimSpace(request.Email)
	name := strings.TrimSpace(request.Name)
	if email == "" {
		return nil, apperror.BadRequest("email is required")
	}
	if name == "" {
		return nil, apperror.BadRequest("name is required")
	}

	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, apperror.Internal("failed to generate API key", err)
	}

	user, err := s.users.Create(ctx, domain.User{
		Email:      email,
		Name:       name,
		APIKeyHash: hashAPIKey(apiKey),
		IsAdmin:    request.IsAdmin,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return nil, apperror.BadRequest("user email already exists")
		}
		return nil, apperror.Internal("failed to create user", err)
	}

	return &types.CreateUserResponse{User: *user, APIKey: apiKey}, nil
}

func (s *AdminService) UserProjects(ctx context.Context, actor *domain.User, userID string) ([]types.ProjectSummary, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, apperror.BadRequest("user id is required")
	}

	projects, err := s.projects.ListByUser(ctx, userID)
	if err != nil {
		return nil, apperror.Internal("failed to list user projects", err)
	}
	return projects, nil
}

func (s *AdminService) Summary(ctx context.Context, actor *domain.User, from, to string) ([]types.DailyPoint, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	from, to = defaultDateRange(from, to)

	points, err := s.projects.DailySummary(ctx, from, to)
	if err != nil {
		return nil, apperror.Internal("failed to load daily summary", err)
	}
	return points, nil
}

func (s *AdminService) UserBreakdown(ctx context.Context, actor *domain.User, userID, groupBy, from, to string) ([]types.UsageBreakdownRow, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, apperror.BadRequest("user id is required")
	}
	groupBy = strings.TrimSpace(groupBy)
	if groupBy == "" {
		groupBy = "day"
	}
	if groupBy != "day" && groupBy != "month" && groupBy != "project" {
		return nil, apperror.BadRequest("group_by must be day, month, or project")
	}
	from, to = defaultDateRange(from, to)

	rows, err := s.projects.UserBreakdown(ctx, userID, groupBy, from, to)
	if err != nil {
		return nil, apperror.Internal("failed to load user breakdown", err)
	}
	return rows, nil
}

func (s *AdminService) UserUsageSummary(ctx context.Context, actor *domain.User, userID, provider, from, to string) (*types.UsageSummaryStats, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, apperror.BadRequest("user id is required")
	}
	provider = strings.TrimSpace(provider)
	if provider == "" {
		provider = "all"
	}
	if provider != "all" && provider != "codex" && provider != "claude" {
		return nil, apperror.BadRequest("provider must be all, codex, or claude")
	}
	from, to = defaultDateRange(from, to)

	stats, err := s.projects.UserUsageSummary(ctx, userID, provider, from, to)
	if err != nil {
		return nil, apperror.Internal("failed to load user usage summary", err)
	}
	return stats, nil
}

func requireAdmin(actor *domain.User) error {
	if actor == nil {
		return apperror.Unauthorized("invalid API key")
	}
	if !actor.IsAdmin {
		return apperror.Forbidden("admin access required")
	}
	return nil
}

func defaultDateRange(from, to string) (string, string) {
	if from == "" {
		from = "2000-01-01"
	}
	if to == "" {
		to = "2099-12-31"
	}
	return from, to
}

func generateAPIKey() (string, error) {
	var token [24]byte
	if _, err := rand.Read(token[:]); err != nil {
		return "", err
	}
	return "ak_" + hex.EncodeToString(token[:]), nil
}

func hashAPIKey(apiKey string) string {
	sum := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(sum[:])
}
