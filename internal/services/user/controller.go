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
	ErrUsernameTaken     = errors.New("username already taken")
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

func (s *UserService) isFieldTaken(ctx context.Context, field, value string, excludeIDs ...string) (bool, error) {
	query := s.collection.
		Where(field, "==", value).
		Where("deletedAt", "==", nil).
		Limit(1)

	if len(excludeIDs) > 0 {
		docRef := s.collection.Doc(excludeIDs[0])
		query = query.Where(firestore.DocumentID, "!=", docRef)
	}

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, fmt.Errorf("firestore query failed: %w", err)
	}
	return len(docs) > 0, nil
}
