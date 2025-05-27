package meeting

import (
	"context"
	"errors"
	"net/http"
	"nitelog/internal/services/meeting"

	"github.com/gin-gonic/gin"
)

// DeleteMeeting godoc
// @Summary      Deleta uma reunião
// @Description  Deleta uma reunião do banco de dados
// @Tags         meeting
// @Accept       json
// @Produce      json
// @Param        meeting_id   path   string true "Id da reunião (BSON primitive.ObjectID)"
// @Success      200         {object}  util.MessageResponse
// @Failure      400         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @Router       /meetings/:id [delete]
func DeleteMeeting(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()

	meetingService := services.NewMeetingService()
	err := meetingService.SoftDelete(ctx, id)

	if errors.Is(err, services.ErrMeetingNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Meeting marked as deleted successfully",
		"details": "The meeting has been soft deleted and can be recovered",
	})
}
