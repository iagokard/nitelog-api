package services

import (
	"cloud.google.com/go/firestore"
)

var firestoreClient *firestore.Client

func SetFirestoreClient(client *firestore.Client) {
	firestoreClient = client
}

func GetCollection(collectionName string) *firestore.CollectionRef {
	return firestoreClient.Collection(collectionName)
}
