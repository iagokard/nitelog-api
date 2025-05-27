package models

import (
	"time"
)

// @model Attendance
type Attendance struct {
	UserID    string     `firestore:"userId" json:"user_id" example:"d4e5f6a7b8c9d0e1f2a3b4c5"`
	StartTime time.Time  `firestore:"startTime" json:"start_time" example:"2025-05-14T20:19:02.1Z"`
	EndTime   *time.Time `firestore:"endTime,omitempty" json:"end_time,omitempty" example:"2025-05-14T22:12:34.1Z"`
}
