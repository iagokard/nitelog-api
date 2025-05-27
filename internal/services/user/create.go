package services

import (
	"context"
	"fmt"
	"time"

	"nitelog/internal/models"
)

func (s *UserService) Create(ctx context.Context, username, email string, pswdHash []byte) (*models.User, error) {
	emailTaken, err := s.isFieldTaken(ctx, "email", email)
	if err != nil {
		return nil, fmt.Errorf("email check failed: %w", err)
	}
	if emailTaken {
		return nil, ErrEmailTaken
	}

	usernameTaken, err := s.isFieldTaken(ctx, "username", username)
	if err != nil {
		return nil, fmt.Errorf("username check failed: %w", err)
	}
	if usernameTaken {
		return nil, ErrUsernameTaken
	}

	user := models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(pswdHash),
		Roles:        []string{},
		CreatedAt:    time.Now(),
	}

	docRef, _, err := s.collection.Add(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = docRef.ID

	doc, err := docRef.Get(ctx)
	if err == nil {
		doc.DataTo(&user)
	}

	return &user, nil
}
