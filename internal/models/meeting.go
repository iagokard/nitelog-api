package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Meeting struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Date        time.Time          `bson:"date"`
	MeetingCode string             `bson:"meeting_code"`
	Attendance  []Attendance       `bson:"attendance"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type Attendance struct {
	UserID    primitive.ObjectID `bson:"user_id"`
	StartTime time.Time          `bson:"start_time"`
	EndTime   *time.Time         `bson:"end_time,omitempty"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`  // MongoDB's _id
	Username     string             `bson:"username" json:"username"` // Required field
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`       // Never JSON encoded
	Roles        []string           `bson:"roles,omitempty" json:"roles"` // Optional array
}
