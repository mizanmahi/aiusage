package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggerRecordsRequest(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	handler := Logger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/ingest", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	output := logs.String()
	if !strings.Contains(output, "msg=\"http request\"") {
		t.Fatalf("log output missing message: %s", output)
	}
	if !strings.Contains(output, "method=POST") {
		t.Fatalf("log output missing method: %s", output)
	}
	if !strings.Contains(output, "path=/ingest") {
		t.Fatalf("log output missing path: %s", output)
	}
	if !strings.Contains(output, "status=201") {
		t.Fatalf("log output missing status: %s", output)
	}
}
