package services

import "go.mongodb.org/mongo-driver/mongo"

var db *mongo.Database

func SetServicesDatabase(database *mongo.Database) { db = database }
