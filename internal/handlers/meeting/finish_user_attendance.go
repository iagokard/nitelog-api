package meeting

import (
	"context"
	"errors"
	"net/http"
	meetingServices "nitelog/internal/services/meeting"
	userServices "nitelog/internal/services/user"
	"nitelog/internal/util"
	"time"

	"github.com/gin-gonic/gin"
)

type FinishUserAttendanceRequest struct {
	Registration string `firestore:"registration" json:"registration" example:"8854652123" binding:"required"`
	Date         string `json:"date" example:"2025-10-26" binding:"required"`
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
// @securityDefinitions.apikey  BearerAuth
// @Router       /meetings/finish-attendance [post]
func FinishUserAttendance(c *gin.Context) {
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

	ctx := context.Background()

	userService := userServices.NewUserService()
	_, err = userService.GetByRegistration(ctx, req.Registration)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	meetingService := meetingServices.NewMeetingService()
	err = meetingService.FinishAttendance(ctx, *normalizedDate, req.Registration)

	if errors.Is(err, meetingServices.ErrMeetingNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if errors.Is(err, meetingServices.ErrNoAttendanceToFinish) {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance finalized successfully"})
}
