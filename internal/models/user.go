package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// @model User
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"68253a5154c3608b34c81d79"`
	Username     string             `bson:"username" json:"username" example:"username01"`
	Email        string             `bson:"email" json:"email" example:"sample@email.com"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Roles        []string           `bson:"roles,omitempty" json:"roles"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at" example:"2025-05-14T20:14:04.245Z"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at" example:"2026-05-14T12:18:34.245Z"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty" example:"2025-05-15T09:45:00Z"`
}
