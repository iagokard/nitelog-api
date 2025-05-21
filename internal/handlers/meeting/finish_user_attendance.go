package meeting

import (
	"context"
	"errors"
	"net/http"
	"nitelog/internal/services"
	"nitelog/internal/util"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FinishUserAttendanceRequest struct {
	UserID string `json:"user_id" binding:"required" example:"68253a5154c3608b34c81d79"`
	Date   string `json:"date" binding:"required" example:"2025-10-26"`
}

// FinishUserAttendance godoc
// @Summary      Finaliza presença em reunião
// @Description  Finaliza a presença usuário do usuário
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        attendance     body   FinishUserAttendanceRequest true "Dados da presença"
// @Success      200         {object}  util.MessageResponse
// @Failure      400         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @Router       /meetings/finish-attendance [post]
func (h *MeetingController) FinishUserAttendance(c *gin.Context) {
	var req FinishUserAttendanceRequest

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

	userID, err := primitive.ObjectIDFromHex(req.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := context.Background()

	userService := services.NewUserService()
	_, err = userService.GetByID(ctx, userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	meetingService := services.NewMeetingService()
	_, err = meetingService.FindAttendance(ctx, *normalizedDate, userID)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Attendance already finalized for this user on the specified date",
		})
		return
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error checking attendance status",
		})
		return
	}

	err = meetingService.FinishAttendance(ctx, *normalizedDate, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance finalized successfully"})
}
