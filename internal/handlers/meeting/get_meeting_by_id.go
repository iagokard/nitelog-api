package meeting

import (
	"context"
	"errors"
	"net/http"

	"nitelog/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
// @Failure      403         {object}  util.ErrorResponse
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

	var meeting models.Meeting
	ctx := context.Background()
	err = h.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&meeting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meeting)
}
