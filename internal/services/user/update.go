package services

import (
	"context"
	"fmt"

	"nitelog/internal/models"
	"nitelog/internal/util"

	"cloud.google.com/go/firestore"
)

func (s *UserService) Update(ctx context.Context, id string, updatedUser models.User) error {
	existingUser, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	var updates []firestore.Update
	changes := false

	if updatedUser.Registration != "" && updatedUser.Registration != existingUser.Registration {
		taken, err := s.isFieldTaken(ctx, "registration", updatedUser.Registration, id)
		if err != nil {
			return fmt.Errorf("registration check failed: %w", err)
		}

		if taken {
			return ErrRegistrationTaken
		}

		updates = append(updates, firestore.Update{
			Path:  "registration",
			Value: updatedUser.Registration,
		})
		changes = true
	}

	if updatedUser.Email != "" && updatedUser.Email != existingUser.Email {
		taken, err := s.isFieldTaken(ctx, "email", updatedUser.Email, id)
		if err != nil {
			return fmt.Errorf("email check failed: %w", err)
		}
		if taken {
			return ErrEmailTaken
		}
		updates = append(updates, firestore.Update{
			Path:  "email",
			Value: updatedUser.Email,
		})
		changes = true
	}

	if updatedUser.Roles != nil && !equalRoles(updatedUser.Roles, existingUser.Roles) {
		updates = append(updates, firestore.Update{
			Path:  "roles",
			Value: updatedUser.Roles,
		})
		changes = true
	}

	if updatedUser.PasswordHash != "" {
		err := util.CheckPassword(existingUser.PasswordHash, updatedUser.PasswordHash)
		if err != nil {
			hashedPassword, err := util.HashPassword(updatedUser.PasswordHash)
			if err != nil {
				return fmt.Errorf("failed to hash password: %w", err)
			}
			updates = append(updates, firestore.Update{
				Path:  "passwordHash",
				Value: hashedPassword,
			})
			changes = true
		}
	}

	if updatedUser.Name != "" && updatedUser.Name != existingUser.Name {
		updates = append(updates, firestore.Update{
			Path:  "name",
			Value: updatedUser.Name,
		})
		changes = true
	}

	if !changes {
		return ErrNoChangesDetected
	}

	updates = append(updates, firestore.Update{
		Path:  "updatedAt",
		Value: firestore.ServerTimestamp,
	})

	_, err = s.collection.Doc(id).Update(ctx, updates)
	return err
}

func equalRoles(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aMap := make(map[string]bool)
	for _, v := range a {
		aMap[v] = true
	}
	for _, v := range b {
		if !aMap[v] {
			return false
		}
	}
	return true
}
