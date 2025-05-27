package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"nitelog/internal/models"
	"nitelog/internal/util"

	"cloud.google.com/go/firestore"
)

func (s *MeetingService) Create(ctx context.Context, date time.Time) (*models.Meeting, error) {
	meetingCode, err := s.generateUniqueMeetingCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate meeting code: %w", err)
	}

	exists, err := s.dateMeetingExists(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("date check failed: %w", err)
	}

	if exists {
		return nil, ErrDuplicateMeeting
	}

	meetingRef, _, err := s.collection.Add(ctx, map[string]any{
		"date":        date,
		"meetingCode": meetingCode,
		"attendance":  []models.Attendance{},
		"createdAt":   firestore.ServerTimestamp,
		"deletedAt":   nil,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create meeting: %w", err)
	}

	doc, err := meetingRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to verify creation: %w", err)
	}

	var meeting models.Meeting
	if err := doc.DataTo(&meeting); err != nil {
		return nil, err
	}

	// TODO check if need this line
	meeting.ID = doc.Ref.ID

	return &meeting, nil
}

func (s *MeetingService) generateUniqueMeetingCode(ctx context.Context) (string, error) {
	const maxAttempts = 10
	for range maxAttempts {
		code := util.GenerateMeetingCode()

		exists, err := isMeetingCodeExists(ctx, s.collection, code)
		if err != nil {
			return "", fmt.Errorf("code check failed: %w", err)
		}

		if !exists {
			return code, nil
		}

		log.Printf("Duplicate meeting code generated: %s, retrying...", code)
	}
	return "", errors.New("failed to generate unique code after 10 attempts")
}

func isMeetingCodeExists(ctx context.Context, coll *firestore.CollectionRef, code string) (bool, error) {
	query := coll.Where("meetingCode", "==", code).Limit(1)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, fmt.Errorf("firestore query failed: %w", err)
	}
	return len(docs) > 0, nil
}

func (s *MeetingService) dateMeetingExists(ctx context.Context, date time.Time) (bool, error) {
	query := s.collection.
		Where("date", "==", date).
		Where("deletedAt", "==", nil).
		Limit(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, fmt.Errorf("firestore query failed: %w", err)
	}

	return len(docs) > 0, nil
}
