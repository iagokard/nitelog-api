package meeting

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
// @Failure      403         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @Router       /meetings/:id [delete]
func (h *MeetingController) DeleteMeeting(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.collection.DeleteOne(
		context.Background(),
		bson.M{"_id": objID},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meeting deleted successfully"})
}
