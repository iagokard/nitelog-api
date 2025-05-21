package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// @model Attendance
type Attendance struct {
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id" example:"68253a5154c3608b34c81d79"`
	StartTime time.Time          `bson:"start_time" json:"start_time" example:"2025-05-14T20:19:02.1Z"`
	EndTime   *time.Time         `bson:"end_time,omitempty" json:"end_time" example:"2025-05-14T22:12:34.1Z"`
}
