package services

import (
	"context"
	"fmt"
	"time"

	"nitelog/internal/models"
)

func (s *MeetingService) FindCompletedAttendance(ctx context.Context, date time.Time, userID string) (*models.Meeting, error) {
	query := s.collection.
		Where("date", "==", date).
		Where("deletedAt", "==", nil).
		Where("attendance.userId", "==", userID).
		Where("attendance.endTime", "!=", nil).
		Limit(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("attendance query failed: %w", err)
	}

	if len(docs) == 0 {
		return nil, ErrAttendanceNotFound
	}

	var meeting models.Meeting
	if err := docs[0].DataTo(&meeting); err != nil {
		return nil, err
	}
	return &meeting, nil
}
