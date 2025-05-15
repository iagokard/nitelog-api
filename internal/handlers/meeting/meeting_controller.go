package meeting

import "go.mongodb.org/mongo-driver/mongo"

type MeetingController struct {
	collection *mongo.Collection
}

func NewMeetingController(coll *mongo.Collection) *MeetingController {
	return &MeetingController{collection: coll}
}
