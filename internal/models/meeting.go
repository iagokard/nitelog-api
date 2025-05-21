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
	DeletedAt   *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty" example:"2025-05-15T09:45:00Z"`
}
