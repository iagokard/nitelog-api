package meeting

import (
	"context"
	"errors"
	"net/http"
	"time"

	"nitelog/internal/services/meeting"
	"nitelog/internal/util"

	"github.com/gin-gonic/gin"
)

// GetMeetingByDate godoc
// @Summary      Procura reunião por data
// @Description  Procura reunião por data especifica
// @Tags         meeting
// @Accept       json
// @Produce      json
// @Param        date        path      string true "Data no estilo: 2024-10-26"
// @Success      200         {object}  models.Meeting
// @Failure      400         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @Router       /meetings/by-date/:date [get]
func GetMeetingByDate(c *gin.Context) {
	dateParam := c.Param("date")

	date, err := time.Parse("2006-01-02", dateParam)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "No meeting for this date"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meeting)
}
