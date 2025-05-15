package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// @model Meeting
type Meeting struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"6824f98cb453ef098596dc92"`
	Date        time.Time          `bson:"date" json:"date" example:"2024-10-26"`
	MeetingCode string             `bson:"meeting_code" json:"meeting_code" example:"qE522Af8"`
	Attendance  []Attendance       `bson:"attendance" json:"attendance"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at" example:"2025-05-14T20:14:04.245Z"`
}

// @model Attendance
type Attendance struct {
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id" example:"68253a5154c3608b34c81d79"`
	StartTime time.Time          `bson:"start_time" json:"start_time" example:"2025-05-14T20:19:02.1Z"`
	EndTime   *time.Time         `bson:"end_time,omitempty" json:"end_time" example:"2025-05-14T22:12:34.1Z"`
}

// @model User
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"68253a5154c3608b34c81d79"`
	Username     string             `bson:"username" json:"username" example:"username01"`
	Email        string             `bson:"email" json:"email" example:"sample@email.com"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Roles        []string           `bson:"roles,omitempty" json:"roles"`
}
