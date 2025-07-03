package meeting

import (
	"context"
	"errors"
	"net/http"

	"nitelog/internal/services/meeting"

	"github.com/gin-gonic/gin"
)

// GetMeetingByID godoc
// @Summary      Procura reunião por id
// @Description  Procura reunião por id especifico
// @Tags         meeting
// @Accept       json
// @Produce      json
// @Param        meeting_id   path     string true "ID da reunião"
// @Success      200         {object}  models.Meeting
// @Failure      400         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @securityDefinitions.apikey  BearerAuth
// @Router       /meetings/:id [get]
func GetMeetingByID(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()

	meetingService := services.NewMeetingService()
	meeting, err := meetingService.GetByID(ctx, id)

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
