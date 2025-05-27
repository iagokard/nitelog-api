package services

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserService) SoftDelete(ctx context.Context, id string) error {
	docRef := s.collection.Doc(id)

	_, err := docRef.Update(ctx, []firestore.Update{
		{
			Path:  "deletedAt",
			Value: firestore.ServerTimestamp,
		},
		{
			Path:  "updatedAt",
			Value: firestore.ServerTimestamp,
		},
	})

	if status.Code(err) == codes.NotFound {
		return ErrUserNotFound
	}

	return err
}
