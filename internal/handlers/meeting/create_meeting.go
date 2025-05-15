package meeting

import (
	"context"
	"log"
	"net/http"
	"time"

	"nitelog/internal/models"
	"nitelog/internal/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *MeetingController) CreateMeeting(c *gin.Context) {
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
