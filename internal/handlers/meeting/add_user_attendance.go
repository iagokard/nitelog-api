package meeting

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	meetingServices "nitelog/internal/services/meeting"
	userServices "nitelog/internal/services/user"
	"nitelog/internal/util"
)

type AddUserAttendanceRequest struct {
	Registration string `firestore:"registration" json:"registration" example:"8854652123" binding:"required"`
	Date         string `json:"date" example:"2025-10-26" binding:"required"`
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
func AddUserAttendance(c *gin.Context) {
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
	meetingService := meetingServices.NewMeetingService()

	_, err = meetingService.GetByDate(ctx, *normalizedDate)

	if errors.Is(err, meetingServices.ErrMeetingNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userService := userServices.NewUserService()
	_, err = userService.GetByRegistration(ctx, req.Registration)

	if errors.Is(err, userServices.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = meetingService.AddAttendance(ctx, *normalizedDate, req.Registration)
	if errors.Is(err, meetingServices.ErrActiveAttendanceExists) {
		c.JSON(http.StatusConflict, gin.H{"error": "User already in attendance"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "User added to attendance")
}
