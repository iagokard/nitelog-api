package services

import (
	"context"
	"fmt"

	"nitelog/internal/models"
)

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := s.collection.
		Where("email", "==", email).
		Where("deletedAt", "==", nil).
		Limit(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if len(docs) == 0 {
		return nil, ErrUserNotFound
	}

	var user models.User
	if err := docs[0].DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	user.ID = docs[0].Ref.ID

	return &user, nil
}
