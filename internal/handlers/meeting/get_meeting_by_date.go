package meeting

import (
	"context"
	"errors"
	"net/http"
	"time"

	"nitelog/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
// @Failure      403         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @Router       /meetings/by-date/:date [get]
func (h *MeetingController) GetMeetingByDate(c *gin.Context) {
	dateParam := c.Param("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	var meeting models.Meeting
	ctx := context.Background()
	err = h.collection.FindOne(ctx, bson.M{"date": date}).Decode(&meeting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No meeting for this date"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meeting)
}
