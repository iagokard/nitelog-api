package services

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MeetingService) SoftDelete(ctx context.Context, id string) error {
	_, err := s.collection.Doc(id).Update(ctx, []firestore.Update{
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
		return ErrMeetingNotFound
	}

	return err
}
