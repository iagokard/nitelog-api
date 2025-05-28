package models

import (
	"time"
)

// @model Attendance
type Attendance struct {
	Registration string     `firestore:"registration" json:"registration" example:"8854652123"`
	StartTime    time.Time  `firestore:"startTime" json:"start_time" example:"2025-05-14T20:19:02.1Z"`
	EndTime      *time.Time `firestore:"endTime,omitempty" json:"end_time,omitempty" example:"2025-05-14T22:12:34.1Z"`
}
