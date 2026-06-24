package main

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/mizanmahi/aiusage/types"
)

func TestSendBatchesSplitsEventsAndCombinesResults(t *testing.T) {
	oldSender := sendUsageEvents
	defer func() { sendUsageEvents = oldSender }()

	events := make([]types.UsageEvent, 21)
	var batchSizes []int
	sendUsageEvents = func(_ string, _ string, _ string, batch []types.UsageEvent) (*types.PushResponse, error) {
		batchSizes = append(batchSizes, len(batch))
		return &types.PushResponse{Accepted: len(batch), Message: "ok"}, nil
	}

	result, err := sendBatches("http://localhost:8080", "ak_secret_value", events)
	if err != nil {
		t.Fatalf("sendBatches() error = %v", err)
	}

	if got, want := batchSizes, []int{10, 10, 1}; !slices.Equal(got, want) {
		t.Fatalf("batch sizes = %v, want %v", got, want)
	}
	if result.Accepted != 21 || result.Skipped != 0 || result.Message != "ok" {
		t.Fatalf("result = %+v, want 21 accepted and ok message", result)
	}
}

func TestSendBatchesStopsOnFailure(t *testing.T) {
	oldSender := sendUsageEvents
	defer func() { sendUsageEvents = oldSender }()

	events := make([]types.UsageEvent, 21)
	calls := 0
	sendUsageEvents = func(_ string, _ string, _ string, _ []types.UsageEvent) (*types.PushResponse, error) {
		calls++
		if calls == 2 {
			return nil, errors.New("server unavailable")
		}
		return &types.PushResponse{Accepted: 10}, nil
	}

	_, err := sendBatches("http://localhost:8080", "ak_secret_value", events)
	if err == nil || !strings.Contains(err.Error(), "batch 2 of 3: server unavailable") {
		t.Fatalf("sendBatches() error = %v, want second batch error", err)
	}
	if calls != 2 {
		t.Fatalf("send calls = %d, want 2", calls)
	}
}
