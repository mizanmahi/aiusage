//go:build integration

package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/mizanmahi/aiusage/types"
)

func TestRepositoriesIntegration(t *testing.T) {
	db := openIntegrationDB(t)
	ctx := context.Background()

	apiKey := "integration-api-key"
	apiKeyHash := hashIntegrationAPIKey(apiKey)
	userID := createIntegrationUser(t, ctx, db, apiKeyHash)
	defer cleanupIntegrationUser(t, ctx, db, userID)

	user, err := NewUserRepo(db).FindByAPIKeyHash(ctx, apiKeyHash)
	if err != nil {
		t.Fatalf("FindByAPIKeyHash() error = %v", err)
	}
	if user.ID != userID || !user.IsAdmin {
		t.Fatalf("user = %+v, want ID %s and admin", user, userID)
	}

	events := []types.UsageEvent{
		integrationEvent("integration-session-001"),
		integrationEvent("integration-session-001"),
	}
	accepted, skipped, err := NewEventRepo(db).Upsert(ctx, userID, events)
	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}
	if accepted != 1 || skipped != 1 {
		t.Fatalf("accepted/skipped = %d/%d, want 1/1", accepted, skipped)
	}

	var stored int
	if err := db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM usage_events
		WHERE user_id = $1 AND session_id = $2
	`, userID, "integration-session-001").Scan(&stored); err != nil {
		t.Fatalf("count usage_events error = %v", err)
	}
	if stored != 1 {
		t.Fatalf("stored events = %d, want 1", stored)
	}
}

func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is required for integration tests")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	if err := db.Ping(); err != nil {
		t.Fatalf("db.Ping() error = %v", err)
	}

	return db
}

func createIntegrationUser(t *testing.T, ctx context.Context, db *sql.DB, apiKeyHash string) string {
	t.Helper()

	email := fmt.Sprintf("integration-%d@example.com", time.Now().UnixNano())
	var userID string
	if err := db.QueryRowContext(ctx, `
		INSERT INTO users (email, name, api_key_hash, is_admin)
		VALUES ($1, 'Integration Test', $2, TRUE)
		RETURNING id::text
	`, email, apiKeyHash).Scan(&userID); err != nil {
		t.Fatalf("insert integration user error = %v", err)
	}

	return userID
}

func cleanupIntegrationUser(t *testing.T, ctx context.Context, db *sql.DB, userID string) {
	t.Helper()

	if _, err := db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID); err != nil {
		t.Fatalf("cleanup integration user error = %v", err)
	}
}

func integrationEvent(sessionID string) types.UsageEvent {
	return types.UsageEvent{
		SessionID:       sessionID,
		Date:            "2026-06-08",
		Tool:            types.ToolCodex,
		Project:         "aiusage",
		Cwd:             "/home/mizanmahi/Work/aiusage",
		Model:           "gpt-5.5",
		InputTokens:     1000,
		OutputTokens:    200,
		CacheTokens:     50,
		ReasoningTokens: 25,
		CostUSD:         0.01,
		PushedAt:        time.Now().UTC(),
	}
}

func hashIntegrationAPIKey(apiKey string) string {
	sum := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(sum[:])
}
