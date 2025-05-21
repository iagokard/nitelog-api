package meeting

import (
	"context"
	"errors"
	"net/http"

	"nitelog/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetMeetingByID godoc
// @Summary      Procura reunião por id
// @Description  Procura reunião por id especifico
// @Tags         meeting
// @Accept       json
// @Produce      json
// @Param        meeting_id   path     string true "ID da reunião BSON primitive.ObjectID"
// @Success      200         {object}  models.Meeting
// @Failure      400         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @Router       /meetings/:id [get]
func (h *MeetingController) GetMeetingByID(c *gin.Context) {
	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	ctx := context.Background()

	meetingService := services.NewMeetingService()
	meeting, err := meetingService.GetByID(ctx, objID)

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
