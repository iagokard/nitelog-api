package meeting

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"nitelog/internal/models"
	"nitelog/internal/services"
	"nitelog/internal/util"
)

type AddUserAttendanceRequest struct {
	UserID string `json:"user_id" example:"68253a5154c3608b34c81d79" binding:"required"`
	Date   string `json:"date" example:"2025-10-26" binding:"required"`
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
// @Failure      404         {object}  util.ErrorResponse
// @Failure      409      {object}  util.ErrorResponse
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

	normalizedDate, err := util.NormalizeDate(date)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Unable to normalize date",
			"details": err.Error(),
		})
		return
	}

	ctx := context.Background()
	meetingService := services.NewMeetingService()

	meeting, err := meetingService.GetByDate(ctx, *normalizedDate)

	if errors.Is(err, services.ErrMeetingNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userService := services.NewUserService()
	_, err = userService.GetByID(ctx, userID)

	if errors.Is(err, services.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, attendance := range meeting.Attendance {
		if attendance.UserID == userID {
			c.JSON(http.StatusConflict, gin.H{"error": "User already in attendance"})
			return
		}
	}

	update := bson.M{
		"$push": bson.M{
			"attendance": models.Attendance{
				UserID:    userID,
				StartTime: time.Now(),
			},
		},
	}

	err = meetingService.Update(ctx, meeting.ID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "User added to attendance")
}
