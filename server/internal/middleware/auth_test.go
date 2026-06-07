package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mizanmahi/aiusage/server/internal/domain"
)

type fakeUserRepo struct {
	user    *domain.User
	gotHash string
	err     error
}

func (r *fakeUserRepo) FindByAPIKeyHash(ctx context.Context, apiKeyHash string) (*domain.User, error) {
	r.gotHash = apiKeyHash
	return r.user, r.err
}

func TestAuthSetsUserContext(t *testing.T) {
	users := &fakeUserRepo{user: &domain.User{ID: "user-1"}}
	handler := Auth(users)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if User(r).ID != "user-1" {
			t.Fatalf("User().ID = %q, want user-1", User(r).ID)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/ingest", nil)
	req.Header.Set("Authorization", "Bearer ak_secret")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", res.Code)
	}
	if users.gotHash != hashAPIKey("ak_secret") {
		t.Fatalf("hash = %q, want hashed API key", users.gotHash)
	}
}

func TestAuthRejectsMissingBearerToken(t *testing.T) {
	users := &fakeUserRepo{user: &domain.User{ID: "user-1"}}
	handler := Auth(users)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler was called")
	}))

	req := httptest.NewRequest(http.MethodPost, "/ingest", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", res.Code)
	}
	if !strings.Contains(res.Body.String(), `"code":"UNAUTHORIZED"`) {
		t.Fatalf("body missing unauthorized code: %s", res.Body.String())
	}
}

func TestAuthRejectsLookupErrors(t *testing.T) {
	users := &fakeUserRepo{err: errors.New("db failed")}
	handler := Auth(users)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler was called")
	}))

	req := httptest.NewRequest(http.MethodPost, "/ingest", nil)
	req.Header.Set("Authorization", "Bearer ak_secret")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", res.Code)
	}
}
