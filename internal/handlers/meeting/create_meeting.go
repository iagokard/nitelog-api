package meeting

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"nitelog/internal/models"
	"nitelog/internal/services"
	"nitelog/internal/util"

	"github.com/gin-gonic/gin"
)

type CreateMeetingRequest struct {
	Date string `json:"date" example:"2025-10-26" binding:"required"`
}

type DuplicatedMeetingErrorResponse struct {
	Error           string          `json:"error" example:"Meeting already exists for this date"`
	ExistingMeeting *models.Meeting `json:"existing_meeting"`
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
// @Failure      409      {object}  DuplicatedMeetingErrorResponse
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
	meeting, err := meetingService.Create(ctx, *normalizedDate)

	if errors.Is(err, services.ErrDuplicateMeeting) {
		res := DuplicatedMeetingErrorResponse{
			Error:           "Meeting already exists for this date",
			ExistingMeeting: meeting,
		}

		c.JSON(http.StatusConflict, res)
		return
	}

	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database operation failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, meeting)
}
