package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/mizanmahi/aiusage/server/internal/apperror"
	"github.com/mizanmahi/aiusage/types"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		if appErr.Err != nil {
			slog.Error("request failed", "code", appErr.Code, "error", appErr.Err)
		}
		writeJSON(w, appErr.Status, types.APIErrorResponse{
			Error: types.APIError{Code: appErr.Code, Message: appErr.Message},
		})
		return
	}

	slog.Error("request failed", "error", err)
	writeJSON(w, http.StatusInternalServerError, types.APIErrorResponse{
		Error: types.APIError{Code: "INTERNAL_ERROR", Message: "internal server error"},
	})
}
