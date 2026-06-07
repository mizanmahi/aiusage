package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mizanmahi/aiusage/server/internal/domain"
	"github.com/mizanmahi/aiusage/server/internal/middleware"
	"github.com/mizanmahi/aiusage/types"
)

type fakeIngestService struct {
	gotUserID  string
	gotPayload types.PushPayload
	err        error
}

func (s *fakeIngestService) Ingest(ctx context.Context, userID string, payload types.PushPayload) (*types.PushResponse, error) {
	s.gotUserID = userID
	s.gotPayload = payload
	if s.err != nil {
		return nil, s.err
	}
	return &types.PushResponse{Accepted: len(payload.Events), Skipped: 0, Message: "ok"}, nil
}

type handlerUserRepo struct {
	user *domain.User
}

func (r handlerUserRepo) FindByAPIKeyHash(ctx context.Context, apiKeyHash string) (*domain.User, error) {
	return r.user, nil
}

func TestIngestHandlerAcceptsValidPayload(t *testing.T) {
	service := &fakeIngestService{}
	handler := authedHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/ingest", strings.NewReader(`{
		"events": [{
			"session_id": "session-1",
			"date": "2026-06-01",
			"tool": "codex",
			"project": "aiusage"
		}]
	}`))
	req.Header.Set("Authorization", "Bearer ak_secret")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", res.Code, res.Body.String())
	}
	if service.gotUserID != "user-1" {
		t.Fatalf("userID = %q, want user-1", service.gotUserID)
	}
	if !strings.Contains(res.Body.String(), `"accepted":1`) {
		t.Fatalf("body missing accepted count: %s", res.Body.String())
	}
}

func TestIngestHandlerRejectsInvalidJSON(t *testing.T) {
	handler := authedHandler(&fakeIngestService{})
	req := httptest.NewRequest(http.MethodPost, "/ingest", strings.NewReader(`{`))
	req.Header.Set("Authorization", "Bearer ak_secret")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.Code)
	}
	if !strings.Contains(res.Body.String(), `"code":"INVALID_PAYLOAD"`) {
		t.Fatalf("body missing invalid payload code: %s", res.Body.String())
	}
}

func TestIngestHandlerRejectsMissingSessionID(t *testing.T) {
	handler := authedHandler(&fakeIngestService{})
	req := httptest.NewRequest(http.MethodPost, "/ingest", strings.NewReader(`{
		"events": [{"date": "2026-06-01", "tool": "codex", "project": "aiusage"}]
	}`))
	req.Header.Set("Authorization", "Bearer ak_secret")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.Code)
	}
	if !strings.Contains(res.Body.String(), "session_id is required") {
		t.Fatalf("body missing validation message: %s", res.Body.String())
	}
}

func TestIngestHandlerWritesServiceErrors(t *testing.T) {
	handler := authedHandler(&fakeIngestService{err: errors.New("unexpected")})
	req := httptest.NewRequest(http.MethodPost, "/ingest", strings.NewReader(`{
		"events": [{"session_id": "s1", "date": "2026-06-01", "tool": "codex", "project": "aiusage"}]
	}`))
	req.Header.Set("Authorization", "Bearer ak_secret")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", res.Code)
	}
}

func authedHandler(service *fakeIngestService) http.Handler {
	ingest := NewIngestHandler(service)
	return middleware.Auth(handlerUserRepo{user: &domain.User{ID: "user-1"}})(http.HandlerFunc(ingest.Create))
}
