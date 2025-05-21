package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nitelog/internal/models"
	"nitelog/internal/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailTaken        = errors.New("email already taken")
	ErrUsernameTaken     = errors.New("username already taken")
	ErrNoChangesDetected = errors.New("no changes detected on update")
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService() *UserService {
	return &UserService{
		collection: db.Collection("users"),
	}
}

func (s *UserService) Create(ctx context.Context, username, email string, pswdHash []byte) (*models.User, error) {
	var existing models.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&existing)

	if err == nil {
		return nil, ErrEmailTaken
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	err = s.collection.FindOne(ctx, bson.M{"username": username}).Decode(&existing)

	if err == nil {
		return nil, ErrUsernameTaken
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	user := models.User{
		ID:           primitive.NewObjectID(),
		Username:     username,
		Email:        email,
		PasswordHash: string(pswdHash),
		Roles:        []string{},
	}

	result, err := s.collection.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to get valid ObjectID from insertion result")
	}

	user.ID = oid
	return &user, nil
}

func (s *UserService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{
		"_id":        id,
		"deleted_at": nil,
	}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{
		"email":      email,
		"deleted_at": nil,
	}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (s *UserService) Update(ctx context.Context, id string, updatedUser models.User) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrUserNotFound
	}

	existingUser, err := s.GetByID(ctx, objID)
	if err != nil {
		return err
	}

	update := bson.M{}
	changes := false

	if updatedUser.Username != "" && updatedUser.Username != existingUser.Username {
		var userWithSameUsername models.User
		err := s.collection.FindOne(ctx, bson.M{"username": updatedUser.Username}).Decode(&userWithSameUsername)
		if err == nil {
			return ErrUsernameTaken
		}
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}

		update["username"] = updatedUser.Username
		changes = true
	}

	if updatedUser.Email != "" && updatedUser.Email != existingUser.Email {
		var userWithSameEmail models.User
		err := s.collection.FindOne(ctx, bson.M{"email": updatedUser.Email}).Decode(&userWithSameEmail)
		if err == nil {
			return ErrEmailTaken
		}
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}

		update["email"] = updatedUser.Email
		changes = true
	}

	if updatedUser.Roles != nil && !equalRoles(updatedUser.Roles, existingUser.Roles) {
		update["roles"] = updatedUser.Roles
		changes = true
	}

	if updatedUser.PasswordHash != "" {
		err := util.CheckPassword(existingUser.PasswordHash, updatedUser.PasswordHash)
		if err != nil {
			hashedPassword, err := util.HashPassword(updatedUser.PasswordHash)
			if err != nil {
				return fmt.Errorf("erro ao gerar hash da senha: %w", err)
			}
			update["password_hash"] = hashedPassword
			changes = true
		}
	}

	if !changes {
		return ErrNoChangesDetected
	}

	update["updated_at"] = time.Now()

	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": update},
	)

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

func (s *UserService) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		}},
	)
	return err
}
