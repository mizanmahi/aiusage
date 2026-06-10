package service

import (
	"context"
	"errors"
	"testing"

	"github.com/mizanmahi/aiusage/types"
)

type fakeEventRepo struct {
	gotUserID string
	gotEvents []types.UsageEvent
	err       error
}

func (r *fakeEventRepo) Upsert(ctx context.Context, userID string, events []types.UsageEvent) (int, int, error) {
	r.gotUserID = userID
	r.gotEvents = events
	if r.err != nil {
		return 0, 0, r.err
	}
	return len(events), 0, nil
}

func TestIngestStoresServerFields(t *testing.T) {
	repo := &fakeEventRepo{}
	service := NewIngestService(repo)

	result, err := service.Ingest(context.Background(), "user-1", types.PushPayload{
		Events: []types.UsageEvent{{
			SessionID:       "session-1",
			Model:           "gpt-5.5",
			InputTokens:     1_000_000,
			OutputTokens:    1_000_000,
			CacheReadTokens: 1_000_000,
		}},
	})
	if err != nil {
		t.Fatalf("Ingest() error = %v", err)
	}

	if result.Accepted != 1 || result.Skipped != 0 {
		t.Fatalf("result = %+v, want accepted 1 skipped 0", result)
	}
	if repo.gotUserID != "user-1" {
		t.Fatalf("userID = %q, want user-1", repo.gotUserID)
	}
	if repo.gotEvents[0].UserID != "user-1" {
		t.Fatalf("event userID = %q, want user-1", repo.gotEvents[0].UserID)
	}
	if repo.gotEvents[0].CostUSD == 0 {
		t.Fatal("CostUSD = 0, want server-calculated cost")
	}
}

func TestIngestRejectsEmptyEvents(t *testing.T) {
	service := NewIngestService(&fakeEventRepo{})

	if _, err := service.Ingest(context.Background(), "user-1", types.PushPayload{}); err == nil {
		t.Fatal("Ingest() error = nil, want validation error")
	}
}

func TestIngestWrapsRepositoryErrors(t *testing.T) {
	service := NewIngestService(&fakeEventRepo{err: errors.New("db failed")})

	_, err := service.Ingest(context.Background(), "user-1", types.PushPayload{
		Events: []types.UsageEvent{{SessionID: "session-1"}},
	})
	if err == nil {
		t.Fatal("Ingest() error = nil, want repository error")
	}
}
