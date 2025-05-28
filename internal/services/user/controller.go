package services

import (
	"context"
	"errors"
	"fmt"

	"nitelog/internal/services"

	"cloud.google.com/go/firestore"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailTaken        = errors.New("email already taken")
	ErrRegistrationTaken = errors.New("username already taken")
	ErrNoChangesDetected = errors.New("no changes detected on update")
)

type UserService struct {
	collection *firestore.CollectionRef
}

func NewUserService() *UserService {
	return &UserService{
		collection: services.GetCollection("users"),
	}
}

func (s *UserService) isFieldTaken(ctx context.Context, field, value string, excludeID string) (bool, error) {
	query := s.collection.
		Where(field, "==", value).
		Where("deletedAt", "==", nil).
		Limit(1)

	if excludeID != "" {
		docRef := s.collection.Doc(excludeID)
		query = query.Where(firestore.DocumentID, "!=", docRef)
	}

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, fmt.Errorf("firestore query failed: %w", err)
	}
	return len(docs) > 0, nil
}

// func (s *UserService) isAdmin(ctx context.Context, id string) (bool, error) {
// 	user, err := s.GetByID(ctx, id)
// 	Vk
