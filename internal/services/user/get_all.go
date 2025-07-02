package services

import (
	"context"
	"fmt"

	"nitelog/internal/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func (s *UserService) GetAllUsers(ctx context.Context) (*[]models.User, error) {
	query := s.collection.
		Where("deletedAt", "==", nil).
		OrderBy("name", firestore.Asc)

	usersIter := query.Documents(ctx)
	defer usersIter.Stop()

	users := make([]models.User, 0)
	for {
		doc, err := usersIter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to iterate users: %w", err)
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			return nil, fmt.Errorf("failed to parse user %s: %w", doc.Ref.ID, err)
		}

		user.ID = doc.Ref.ID
		users = append(users, user)
	}

	return &users, nil
}
