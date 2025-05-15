package user

import "go.mongodb.org/mongo-driver/mongo"

type UserController struct {
	collection *mongo.Collection
}

func NewUserController(coll *mongo.Collection) *UserController {
	return &UserController{collection: coll}
}
