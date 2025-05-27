package services

import (
	"errors"

	"nitelog/internal/services"

	"cloud.google.com/go/firestore"
)

var (
	ErrMeetingNotFound   = errors.New("meeting not found")
	ErrDuplicateMeeting  = errors.New("meeting already exists for this date")
	ErrMeetingCodeTaken  = errors.New("meeting code already taken")
	ErrDateTaken         = errors.New("meeting date already taken")
	ErrNoChangesDetected = errors.New("no changes detected on meeting update")

	ErrNoAttendanceToFinish   = errors.New("attendance not started or already finished for this date")
	ErrActiveAttendanceExists = errors.New("attendance already started")
	ErrAttendanceNotFound     = errors.New("attendance not found")
)

type MeetingService struct {
	collection *firestore.CollectionRef
}

func NewMeetingService() *MeetingService {
	return &MeetingService{
		collection: services.GetCollection("meetings"),
	}
}
