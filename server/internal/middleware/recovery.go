package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mizanmahi/aiusage/types"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error("panic recovered",
						"error", recovered,
						"method", r.Method,
						"path", r.URL.Path,
					)
					writeRecoveryError(w)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func writeRecoveryError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(types.APIErrorResponse{
		Error: types.APIError{
			Code:    "INTERNAL",
			Message: "internal server error",
		},
	})
}
