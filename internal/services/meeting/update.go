package services

import (
	"context"
	"fmt"
	"time"

	"nitelog/internal/models"

	"cloud.google.com/go/firestore"
)

func (s *MeetingService) Update(ctx context.Context, id string, updatedMeeting models.Meeting) error {
	existingMeeting, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	update := make(map[string]interface{})
	changes := false

	if updatedMeeting.MeetingCode != "" && updatedMeeting.MeetingCode != existingMeeting.MeetingCode {
		exists, err := s.isMeetingCodeTaken(ctx, updatedMeeting.MeetingCode, id)
		if err != nil {
			return fmt.Errorf("meeting code check failed: %w", err)
		}
		if exists {
			return ErrMeetingCodeTaken
		}

		update["meetingCode"] = updatedMeeting.MeetingCode
		changes = true
	}

	if !updatedMeeting.Date.IsZero() && !updatedMeeting.Date.Equal(existingMeeting.Date) {
		exists, err := s.isDateTaken(ctx, updatedMeeting.Date, id)
		if err != nil {
			return fmt.Errorf("date check failed: %w", err)
		}
		if exists {
			return ErrDateTaken
		}

		update["date"] = updatedMeeting.Date
		changes = true
	}

	// Check attendance changes
	if updatedMeeting.Attendance != nil && !equalAttendance(updatedMeeting.Attendance, existingMeeting.Attendance) {
		update["attendance"] = updatedMeeting.Attendance
		changes = true
	}

	if !changes {
		return ErrNoChangesDetected
	}

	// Add timestamp and perform update
	update["updatedAt"] = firestore.ServerTimestamp
	_, err = s.collection.Doc(id).Update(ctx, []firestore.Update{
		{Path: "updatedAt", Value: firestore.ServerTimestamp},
	})

	return err
}

func (s *MeetingService) isMeetingCodeTaken(ctx context.Context, code string, excludeID string) (bool, error) {
	query := s.collection.
		Where("meetingCode", "==", code).
		Where("deletedAt", "==", nil).
		Limit(1)

	if excludeID != "" {
		query = query.Where(firestore.DocumentID, "!=", excludeID)
	}

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, fmt.Errorf("firestore query failed: %w", err)
	}
	return len(docs) > 0, nil
}

func (s *MeetingService) isDateTaken(ctx context.Context, date time.Time, excludeID string) (bool, error) {
	query := s.collection.
		Where("date", "==", date).
		Where("deletedAt", "==", nil).
		Limit(1)

	if excludeID != "" {
		query = query.Where(firestore.DocumentID, "!=", excludeID)
	}

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, fmt.Errorf("firestore query failed: %w", err)
	}
	return len(docs) > 0, nil
}

func equalAttendance(a, b []models.Attendance) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[string]models.Attendance)
	for _, item := range a {
		aMap[item.UserID+item.StartTime.String()] = item
	}

	for _, item := range b {
		key := item.UserID + item.StartTime.String()
		if _, exists := aMap[key]; !exists {
			return false
		}
	}
	return true
}
