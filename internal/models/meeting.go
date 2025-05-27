package models

import (
	"time"
)

// @model Meeting
type Meeting struct {
	ID          string       `firestore:"-" json:"id" example:"a1b2c3d4e5f6g7h8i9j0k1"`
	Date        time.Time    `firestore:"date" json:"date" example:"2024-10-26"`
	MeetingCode string       `firestore:"meetingCode" json:"meeting_code" example:"qE522Af8"`
	Attendance  []Attendance `firestore:"attendance" json:"attendance"`
	CreatedAt   time.Time    `firestore:"createdAt" json:"created_at" example:"2025-05-14T20:14:04.245Z"`
	DeletedAt   *time.Time   `firestore:"deletedAt,omitempty" json:"deleted_at,omitempty" example:"2025-05-15T09:45:00Z"`
}
