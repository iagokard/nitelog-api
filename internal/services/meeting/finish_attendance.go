package services

import (
	"context"
	"fmt"
	"time"

	"nitelog/internal/models"

	"cloud.google.com/go/firestore"
)

func (s *MeetingService) FinishAttendance(ctx context.Context, date time.Time, userID string) error {
	query := s.collection.
		Where("date", "==", date).
		Where("deletedAt", "==", nil).
		Limit(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to find meeting: %w", err)
	}
	if len(docs) == 0 {
		return ErrMeetingNotFound
	}

	docRef := docs[0].Ref
	var meeting models.Meeting
	if err := docs[0].DataTo(&meeting); err != nil {
		return err
	}

	updated := false
	now := time.Now()
	for i := range meeting.Attendance {
		if meeting.Attendance[i].UserID == userID && meeting.Attendance[i].EndTime == nil {
			meeting.Attendance[i].EndTime = &now
			updated = true
			break
		}
	}

	if !updated {
		return ErrNoAttendanceToFinish
	}

	_, err = docRef.Update(ctx, []firestore.Update{
		{
			Path:  "attendance",
			Value: meeting.Attendance,
		},
		{
			Path:  "updatedAt",
			Value: firestore.ServerTimestamp,
		},
	})

	return err
}
