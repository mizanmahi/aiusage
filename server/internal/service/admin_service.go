package service

import (
	"context"

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
	if from == "" {
		from = "2000-01-01"
	}
	if to == "" {
		to = "2099-12-31"
	}

	points, err := s.projects.DailySummary(ctx, from, to)
	if err != nil {
		return nil, apperror.Internal("failed to load daily summary", err)
	}
	return points, nil
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
