package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mizanmahi/aiusage/server/internal/domain"
	"github.com/mizanmahi/aiusage/server/internal/middleware"
	"github.com/mizanmahi/aiusage/types"
)

type AdminService interface {
	Users(ctx context.Context, actor *domain.User) ([]types.UserSummary, error)
	UserProjects(ctx context.Context, actor *domain.User, userID string) ([]types.ProjectSummary, error)
	Summary(ctx context.Context, actor *domain.User, from, to string) ([]types.DailyPoint, error)
}

type AdminHandler struct {
	service AdminService
}

func NewAdminHandler(service AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

func (h *AdminHandler) Users(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.Users(r.Context(), middleware.User(r))
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, types.APIResponse[[]types.UserSummary]{Data: users})
}

func (h *AdminHandler) UserProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.UserProjects(r.Context(), middleware.User(r), chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, types.APIResponse[[]types.ProjectSummary]{Data: projects})
}

func (h *AdminHandler) Summary(w http.ResponseWriter, r *http.Request) {
	points, err := h.service.Summary(
		r.Context(),
		middleware.User(r),
		r.URL.Query().Get("from"),
		r.URL.Query().Get("to"),
	)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, types.APIResponse[[]types.DailyPoint]{Data: points})
}
