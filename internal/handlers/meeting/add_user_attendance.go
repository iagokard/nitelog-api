package meeting

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"nitelog/internal/models"
)

type AddUserAttendanceRequest struct {
	UserID      string `json:"userId" binding:"required"`
	Date        string `json:"date" binding:"required"`
	MeetingCode string `json:"meetingCode" binding:"required"`
}

// AddUserAttendance godoc
// @Summary      Registra presença em reunião
// @Description  Adiciona usuário à lista de presença
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        attendance   body     AddUserAttendanceRequest true "Dados da presença"
// @Success      200         {object}  util.MessageResponse
// @Failure      400         {object}  util.ErrorResponse
// @Failure      403         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @Router       /meetings/add-attendance [post]
func (h *MeetingController) AddUserAttendance(c *gin.Context) {
	var req AddUserAttendanceRequest

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

	c.JSON(http.StatusOK, "User added to attendance")
}
