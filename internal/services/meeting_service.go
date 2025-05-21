package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"nitelog/internal/models"
	"nitelog/internal/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrMeetingNotFound      = errors.New("meeting not found")
	ErrDuplicateMeeting     = errors.New("meeting already exists for this date")
	ErrNoAttendanceToFinish = errors.New("attendance not started or already finished for this date")
)

type MeetingService struct {
	collection *mongo.Collection
}

func NewMeetingService() *MeetingService {
	return &MeetingService{
		collection: db.Collection("meetings"),
	}
}

func (s *MeetingService) Create(ctx context.Context, date time.Time) (*models.Meeting, error) {
	isCodeExists := func(code string) (bool, error) {
		count, err := s.collection.CountDocuments(ctx, bson.M{
			"meeting_code": code,
		})
		return count > 0, err
	}

	var meetingCode string
	for range 10 {
		meetingCode = util.GenerateMeetingCode()
		exists, err := isCodeExists(meetingCode)

		if err != nil {
			return nil, fmt.Errorf("code check failed: %w", err)
		}

		if !exists {
			break
		}

		log.Println("Existing meeting code was generated, retrying...")
		// TODO properly return error
		meetingCode = "" // Reset if all attempts fail
	}

	meeting := models.Meeting{
		ID:          primitive.NewObjectID(),
		Date:        date,
		MeetingCode: meetingCode,
		Attendance:  []models.Attendance{},
		CreatedAt:   time.Now().UTC(),
	}

	filter := bson.M{
		"date": meeting.Date,
		"$or": []bson.M{
			{"deleted_at": nil},
			{"deleted_at": bson.M{"$exists": false}},
		},
	}

	var existing models.Meeting
	err := s.collection.FindOne(ctx, filter).Decode(&existing)

	if err == nil {
		return &existing, ErrDuplicateMeeting
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	result, err := s.collection.InsertOne(ctx, meeting)
	if err != nil {
		return nil, err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to get valid ObjectID from insertion result")
	}

	meeting.ID = oid
	return &meeting, nil
}

func (s *MeetingService) GetByDate(ctx context.Context, date time.Time) (*models.Meeting, error) {
	var meeting models.Meeting
	err := s.collection.FindOne(ctx, bson.M{
		"date":       date,
		"deleted_at": nil,
	}).Decode(&meeting)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrMeetingNotFound
	}
	return &meeting, err
}

func (s *MeetingService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Meeting, error) {
	var meeting models.Meeting
	err := s.collection.FindOne(ctx, bson.M{
		"_id":        id,
		"deleted_at": nil,
	}).Decode(&meeting)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrMeetingNotFound
	}
	return &meeting, err
}

func (s *MeetingService) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := s.collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	)
	return err
}

func (s *MeetingService) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		}},
	)
	return err
}

func (s *MeetingService) FindAttendance(ctx context.Context, date time.Time, userID primitive.ObjectID) (*bson.M, error) {
	filter := bson.M{
		"date": date,
		"attendance": bson.M{
			"$elemMatch": bson.M{
				"user_id":  userID,
				"end_time": bson.M{"$ne": nil},
			},
		},
	}

	var existingRecord bson.M
	err := s.collection.FindOne(ctx, filter).Decode(&existingRecord)

	return &existingRecord, err
}

func (s *MeetingService) FinishAttendance(ctx context.Context, date time.Time, userID primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"attendance.$[elem].end_time": time.Now(),
		},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []any{
			bson.M{"elem.user_id": userID, "elem.end_time": bson.M{"$eq": nil}},
		},
	})

	result, err := s.collection.UpdateOne(
		ctx,
		bson.M{"date": date},
		update,
		arrayFilters,
	)

	if result.MatchedCount == 0 {
		return ErrNoAttendanceToFinish
	}

	return err
}
