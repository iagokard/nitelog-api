package models

import (
	"time"
)

// @model User
type User struct {
	ID           string     `firestore:"-" json:"id" example:"d4e5f6a7b8c9d0e1f2a3b4c5"`
	Registration string     `firestore:"registration" json:"registration" example:"42682087032"`
	Username     string     `firestore:"username" json:"username" example:"username01"`
	Name         string     `firestore:"name" json:"name" example:"John Testes"`
	Email        string     `firestore:"email" json:"email" example:"sample@email.com"`
	PasswordHash string     `firestore:"passwordHash" json:"-"`
	Roles        []string   `firestore:"roles,omitempty" json:"roles"`
	CreatedAt    time.Time  `firestore:"createdAt" json:"created_at" example:"2025-05-14T20:14:04.245Z"`
	UpdatedAt    time.Time  `firestore:"updatedAt" json:"updated_at" example:"2026-05-14T12:18:34.245Z"`
	DeletedAt    *time.Time `firestore:"deletedAt" json:"deleted_at,omitempty" example:"2025-05-15T09:45:00Z"`
}
