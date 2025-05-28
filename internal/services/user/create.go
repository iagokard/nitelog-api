package services

import (
	"context"
	"fmt"
	"time"

	"nitelog/internal/models"
)

func (s *UserService) Create(ctx context.Context, registration, email string, pswdHash []byte) (*models.User, error) {
	emailTaken, err := s.isFieldTaken(ctx, "email", email, "")
	if err != nil {
		return nil, fmt.Errorf("email check failed: %w", err)
	}
	if emailTaken {
		return nil, ErrEmailTaken
	}

	registrationTaken, err := s.isFieldTaken(ctx, "registration", registration, "")
	if err != nil {
		return nil, fmt.Errorf("registration check failed: %w", err)
	}
	if registrationTaken {
		return nil, ErrRegistrationTaken
	}

	user := models.User{
		Registration: registration,
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
