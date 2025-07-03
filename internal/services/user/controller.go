package services

import (
	"context"
	"errors"
	"fmt"

	"nitelog/internal/models"
	"nitelog/internal/services"
	"nitelog/internal/util"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
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

func GetAuthJWTWithUser(ginContext *gin.Context) (*models.User, error) {
	userID, err := util.GetAuthJWT(ginContext)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	userService := NewUserService()

	user, err := userService.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
