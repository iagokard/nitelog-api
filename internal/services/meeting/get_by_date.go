package services

import (
	"context"
	"fmt"
	"time"

	"nitelog/internal/models"
)

func (s *MeetingService) GetByDate(ctx context.Context, date time.Time) (*models.Meeting, error) {
	query := s.collection.
		Where("date", "==", date).
		Where("deletedAt", "==", nil).
		Limit(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to query meetings: %w", err)
	}

	if len(docs) == 0 {
		return nil, ErrMeetingNotFound
	}

	doc := docs[0]
	var meeting models.Meeting
	if err := doc.DataTo(&meeting); err != nil {
		return nil, fmt.Errorf("failed to decode meeting: %w", err)
	}

	meeting.ID = doc.Ref.ID

	return &meeting, nil
}
