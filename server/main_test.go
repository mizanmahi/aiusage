package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouterHealth(t *testing.T) {
	router := newRouter(slog.New(slog.NewTextHandler(io.Discard, nil)), nil, nil, nil, "")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.Code)
	}
	if strings.TrimSpace(res.Body.String()) != "ok" {
		t.Fatalf("body = %q, want ok", res.Body.String())
	}
}
