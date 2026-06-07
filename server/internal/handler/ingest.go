package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mizanmahi/aiusage/server/internal/apperror"
	"github.com/mizanmahi/aiusage/server/internal/middleware"
	"github.com/mizanmahi/aiusage/types"
)

type IngestService interface {
	Ingest(ctx context.Context, userID string, payload types.PushPayload) (*types.PushResponse, error)
}

type IngestHandler struct {
	service IngestService
}

func NewIngestHandler(service IngestService) *IngestHandler {
	return &IngestHandler{service: service}
}

func (h *IngestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload types.PushPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, apperror.BadRequest("invalid JSON payload"))
		return
	}
	if err := validatePushPayload(payload); err != nil {
		writeError(w, err)
		return
	}

	user := middleware.User(r)
	if user == nil {
		writeError(w, apperror.Unauthorized("invalid API key"))
		return
	}

	result, err := h.service.Ingest(r.Context(), user.ID, payload)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, types.APIResponse[*types.PushResponse]{Data: result})
}

func validatePushPayload(payload types.PushPayload) error {
	if len(payload.Events) == 0 {
		return apperror.BadRequest("events are required")
	}
	for _, event := range payload.Events {
		if event.SessionID == "" {
			return apperror.BadRequest("session_id is required")
		}
		if event.Date == "" {
			return apperror.BadRequest("date is required")
		}
		if event.Tool != types.ToolClaude && event.Tool != types.ToolCodex {
			return apperror.BadRequest("tool must be claude or codex")
		}
		if event.Project == "" {
			return apperror.BadRequest("project is required")
		}
	}
	return nil
}
