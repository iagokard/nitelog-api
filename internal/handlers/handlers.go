package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"nitelog/internal/config"
	"nitelog/internal/models"
	"nitelog/internal/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataHandler struct {
	collection *mongo.Collection
}

func NewDataHandler(db *mongo.Database, collection string) *DataHandler {
	return &DataHandler{
		collection: db.Collection(collection),
	}
}

func (h *DataHandler) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	checkCtx, checkCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer checkCancel()

	var existing models.User
	err := h.collection.FindOne(checkCtx, bson.M{
		"$or": []bson.M{
			{"email": req.Email},
			{"username": req.Username},
		},
	}).Decode(&existing)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already exists",
			"user":  existing,
		})
		return
	}

	hash, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error hashing password",
		})
		return
	}

	newUser := models.User{
		ID:           primitive.NewObjectID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		Roles:        []string{},
	}

	insertCtx, insertCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer insertCancel()

	_, err = h.collection.InsertOne(insertCtx, newUser)
	if err != nil {
		log.Printf("Insert error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to create meeting",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func (h *DataHandler) LoginUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := h.collection.FindOne(ctx, bson.M{
		"email": req.Email,
	}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credentials invalid"})
		return
	}

	if err := util.CheckPassword(user.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credentials invalid"})
		return
	}

	cfg := config.Load()
	token, err := util.GenerateJWT(user.ID.Hex(), cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *DataHandler) UpdateUser(c *gin.Context) {
	// Extrai userID do token
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Verifica se path param coincide com token
	idParam := c.Param("id")
	if idParam != uid.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot update other user"})
		return
	}

	var req struct {
		Username *string `json:"username" binding:"omitempty,min=3"`
		Email    *string `json:"email" binding:"omitempty,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"username": req.Username,
			"email":    req.Email,
		},
	}

	var updated models.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.collection.FindOneAndUpdate(ctx, bson.M{"_id": oid}, update, nil).Decode(&updated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *DataHandler) CreateMeeting(c *gin.Context) {
	var req struct {
		Date string `json:"date" binding:"required"`
	}

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Parse date with strict validation
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Invalid date format",
			"example":  "2023-10-05",
			"received": req.Date,
		})
		return
	}

	// Normalize to UTC midnight
	normalizedDate := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check for existing meeting with proper error handling
	var existing models.Meeting
	err = h.collection.FindOne(ctx, bson.M{
		"date": normalizedDate,
	}).Decode(&existing)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Meeting already exists for this date",
			"meeting": existing,
		})
		return
	}

	if err != mongo.ErrNoDocuments {
		// Log detailed error for debugging
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Database operation failed",
			"detail": err.Error(), // Return actual error details
		})
		return
	}

	// Create new meeting document
	newMeeting := models.Meeting{
		ID:          primitive.NewObjectID(),
		Date:        normalizedDate,
		MeetingCode: util.GenerateMeetingCode(),
		Attendance:  []models.Attendance{},
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Insert with timeout
	insertCtx, insertCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer insertCancel()

	_, err = h.collection.InsertOne(insertCtx, newMeeting)
	if err != nil {
		log.Printf("Insert error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to create meeting",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, newMeeting)
}

func (h *DataHandler) AddUserAttendance(c *gin.Context) {
	var req struct {
		UserID      string `json:"userId" binding:"required"`
		Date        string `json:"date" binding:"required"`
		MeetingCode string `json:"meetingCode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Find meeting
	var meeting models.Meeting
	ctx := context.Background()
	err = h.collection.FindOne(ctx, bson.M{"date": date}).Decode(&meeting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if meeting.MeetingCode != req.MeetingCode {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid meeting code"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check existing attendance
	for _, attendance := range meeting.Attendance {
		if attendance.UserID == userID {
			c.JSON(http.StatusConflict, gin.H{"error": "User already in attendance"})
			return
		}
	}

	// Add attendance
	update := bson.M{
		"$push": bson.M{
			"attendance": models.Attendance{
				UserID:    userID,
				StartTime: time.Now(),
			},
		},
	}

	_, err = h.collection.UpdateByID(ctx, meeting.ID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to attendance"})
}

func (h *DataHandler) FinishUserAttendance(c *gin.Context) {
	var req struct {
		UserID string `json:"userId" binding:"required"`
		Date   string `json:"date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Step 1: Check if the attendance has already been finalized
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
	err = h.collection.FindOne(ctx, filter).Decode(&existingRecord)
	if err == nil {
		// Attendance already finalized
		c.JSON(http.StatusBadRequest, gin.H{"error": "Attendance already finalized for this user on the specified date"})
		return
	} else if err != mongo.ErrNoDocuments {
		// An error occurred during the query
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking attendance status"})
		return
	}

	// Step 2: Finalize attendance by setting end_time
	update := bson.M{
		"$set": bson.M{
			"attendance.$[elem].end_time": time.Now(),
		},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.user_id": userID, "elem.end_time": bson.M{"$eq": nil}},
		},
	})

	result, err := h.collection.UpdateOne(
		ctx,
		bson.M{"date": date},
		update,
		arrayFilters,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating attendance"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No matching attendance record found or attendance already finalized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance finalized successfully"})
}

func (h *DataHandler) GetMeetingByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var meeting models.Meeting
	ctx := context.Background()
	err = h.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&meeting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meeting)
}

func (h *DataHandler) GetMeetingByDate(c *gin.Context) {
	dateParam := c.Param("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	var meeting models.Meeting
	ctx := context.Background()
	err = h.collection.FindOne(ctx, bson.M{"date": date}).Decode(&meeting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No meeting for this date"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meeting)
}

// func (h *DataHandler) UpdateMeeting(c *gin.Context) {
// 	id := c.Param("id")
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
// 		return
// 	}
//
// 	var updateFields map[string]interface{}
// 	if err := c.ShouldBindJSON(&updateFields); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	update := bson.M{"$set": updateFields}
// 	update["$set"].(bson.M)["updated_at"] = time.Now()
//
// 	result, err := h.collection.UpdateByID(
// 		context.Background(),
// 		objID,
// 		update,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	if result.MatchedCount == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{"message": "Meeting updated successfully"})
// }

func (h *DataHandler) UpdateMeetingCode(c *gin.Context) {
	dateParam := c.Param("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	var meeting models.Meeting
	ctx := context.Background()
	err = h.collection.FindOne(ctx, bson.M{"date": date}).Decode(&meeting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No meeting for this date"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newCode := util.GenerateMeetingCode()
	_, err = h.collection.UpdateByID(
		ctx,
		meeting.ID,
		bson.M{"$set": bson.M{
			"meeting_code": newCode,
			"updated_at":   time.Now(),
		}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"new_meeting_code": newCode})
}

func (h *DataHandler) DeleteMeeting(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.collection.DeleteOne(
		context.Background(),
		bson.M{"_id": objID},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meeting deleted successfully"})
}
