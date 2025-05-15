package meeting

import (
	"context"
	"fmt"
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

type CreateMeetingRequest struct {
	Date string `json:"date" example:"2025-10-26" binding:"required"`
}

// CreateMeeting godoc
// @Summary      Cria uma nova reunião
// @Description  Registra uma nova reunião com código único
// @Tags         meeting
// @Accept       json
// @Produce      json
// @Param        meeting  body      CreateMeetingRequest  true  "Data da Reunião"
// @Success      201      {object}  models.Meeting
// @Failure      400      {object}  util.ErrorResponse
// @Failure      409      {object}  util.ErrorResponse
// @Failure      500      {object}  util.ErrorResponse
// @Router       /meetings [post]
func (h *MeetingController) CreateMeeting(c *gin.Context) {
	var req CreateMeetingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		details := fmt.Sprintf("examle: 2025-04-02, received: %s", req.Date)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid date format",
			"details": details,
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

	var existing models.Meeting
	err = h.collection.FindOne(ctx, bson.M{
		"date": normalizedDate,
	}).Decode(&existing)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Meeting already exists for this date",
			"details": existing,
		})
		return
	}

	if err != mongo.ErrNoDocuments {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database operation failed",
			"details": err.Error(), // Return actual error details
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
	}

	// Insert with timeout
	insertCtx, insertCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer insertCancel()

	_, err = h.collection.InsertOne(insertCtx, newMeeting)
	if err != nil {
		log.Printf("Insert error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create meeting",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, newMeeting)
}
