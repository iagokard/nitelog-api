package services

import (
	"context"
	"fmt"

	"nitelog/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MeetingService) GetByID(ctx context.Context, id string) (*models.Meeting, error) {
	docRef := s.collection.Doc(id)
	doc, err := docRef.Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrMeetingNotFound
		}
		return nil, fmt.Errorf("failed to get meeting: %w", err)
	}

	deletedAt, err := doc.DataAt("deletedAt")
	if err != nil || deletedAt != nil {
		return nil, ErrMeetingNotFound
	}

	var meeting models.Meeting
	if err := doc.DataTo(&meeting); err != nil {
		return nil, fmt.Errorf("failed to decode meeting: %w", err)
	}

	meeting.ID = doc.Ref.ID

	return &meeting, nil
}
