package services

import (
	"context"
	"fmt"

	"nitelog/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	docRef := s.collection.Doc(id)
	doc, err := docRef.Get(ctx)

	if status.Code(err) == codes.NotFound {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	if user.DeletedAt != nil {
		return nil, ErrUserNotFound
	}

	user.ID = doc.Ref.ID

	return &user, nil
}
