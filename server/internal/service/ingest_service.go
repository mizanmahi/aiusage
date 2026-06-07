package service

import (
	"context"

	"github.com/mizanmahi/aiusage/server/internal/apperror"
	"github.com/mizanmahi/aiusage/server/internal/repository"
	"github.com/mizanmahi/aiusage/types"
)

type IngestService struct {
	events repository.EventRepository
}

func NewIngestService(events repository.EventRepository) *IngestService {
	return &IngestService{events: events}
}

func (s *IngestService) Ingest(ctx context.Context, userID string, payload types.PushPayload) (*types.PushResponse, error) {
	if userID == "" {
		return nil, apperror.Unauthorized("invalid API key")
	}
	if len(payload.Events) == 0 {
		return nil, apperror.BadRequest("events are required")
	}

	events := withServerFields(userID, payload.Events)
	accepted, skipped, err := s.events.Upsert(ctx, userID, events)
	if err != nil {
		return nil, apperror.Internal("failed to ingest events", err)
	}

	return &types.PushResponse{Accepted: accepted, Skipped: skipped, Message: "ok"}, nil
}

func withServerFields(userID string, events []types.UsageEvent) []types.UsageEvent {
	result := make([]types.UsageEvent, len(events))
	for i, event := range events {
		event.UserID = userID
		event.CostUSD = calculateCost(event)
		result[i] = event
	}
	return result
}
